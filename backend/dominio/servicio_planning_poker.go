package dominio

import "context"

// RepositorioVotaciones guarda y lee las votaciones de Planning Poker con sus rondas y votos. Lo
// implementan postgres/ (base real) y memoria/ (tests). Las búsquedas reciben también el proyecto:
// una votación de otro proyecto se trata como inexistente (spec 08 RN-11).
type RepositorioVotaciones interface {
	// Guardar registra una votación nueva y la devuelve con el ID asignado.
	Guardar(ctx context.Context, v Votacion) (Votacion, error)
	// Actualizar reemplaza el estado, la acordada y las rondas con sus votos.
	Actualizar(ctx context.Context, v Votacion) error
	// Listar devuelve las votaciones del proyecto (en cualquier orden).
	Listar(ctx context.Context, proyectoID int64) ([]Votacion, error)
	// Obtener devuelve una votación del proyecto, o ErrVotacionNoExiste.
	Obtener(ctx context.Context, proyectoID, id int64) (Votacion, error)
}

// ServicioPlanningPoker son los casos de uso de la spec 08. Lee los integrantes para saber quién
// es Developer (RN-2) y registra la estimación acordada con el backlog (spec 07).
type ServicioPlanningPoker struct {
	proyectos   RepositorioProyectos
	historias   RepositorioHistorias
	integrantes RepositorioIntegrantes
	votaciones  RepositorioVotaciones
	backlog     *ServicioBacklog
}

func NuevoServicioPlanningPoker(proyectos RepositorioProyectos, historias RepositorioHistorias, integrantes RepositorioIntegrantes, votaciones RepositorioVotaciones, backlog *ServicioBacklog) *ServicioPlanningPoker {
	return &ServicioPlanningPoker{proyectos: proyectos, historias: historias, integrantes: integrantes, votaciones: votaciones, backlog: backlog}
}

// Abrir abre una votación para una historia del backlog (CA-1).
func (s *ServicioPlanningPoker) Abrir(ctx context.Context, proyectoID, historiaID int64) (VistaVotacion, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return VistaVotacion{}, err
	}
	h, err := s.historias.Obtener(ctx, proyectoID, historiaID)
	if err != nil {
		return VistaVotacion{}, err
	}
	integrantes, err := s.integrantes.Listar(ctx, proyectoID)
	if err != nil {
		return VistaVotacion{}, err
	}
	existentes, err := s.votaciones.Listar(ctx, proyectoID)
	if err != nil {
		return VistaVotacion{}, err
	}
	v, err := NuevaVotacion(h, integrantes, existentes)
	if err != nil {
		return VistaVotacion{}, err
	}
	if v, err = s.votaciones.Guardar(ctx, v); err != nil {
		return VistaVotacion{}, err
	}
	return v.Vista(integrantes), nil
}

// Votar registra el voto de un integrante del proyecto en la ronda actual (CA-2). Primero se
// verifica que existan la votación y el integrante (CA-8); después, las reglas de la votación.
func (s *ServicioPlanningPoker) Votar(ctx context.Context, proyectoID, votacionID, integranteID int64, valor int) (VistaVotacion, error) {
	v, integrantes, err := s.votacionConEquipo(ctx, proyectoID, votacionID)
	if err != nil {
		return VistaVotacion{}, err
	}
	i, err := s.integrantes.Obtener(ctx, proyectoID, integranteID)
	if err != nil {
		return VistaVotacion{}, err
	}
	if err := v.Votar(i, valor); err != nil {
		return VistaVotacion{}, err
	}
	return s.guardar(ctx, v, integrantes)
}

// Ver devuelve la votación sin los valores de la ronda actual mientras no se revele (CA-3).
func (s *ServicioPlanningPoker) Ver(ctx context.Context, proyectoID, votacionID int64) (VistaVotacion, error) {
	v, integrantes, err := s.votacionConEquipo(ctx, proyectoID, votacionID)
	if err != nil {
		return VistaVotacion{}, err
	}
	return v.Vista(integrantes), nil
}

// Revelar muestra los votos de la ronda actual cuando votaron todos los Developers (CA-4, CA-5).
func (s *ServicioPlanningPoker) Revelar(ctx context.Context, proyectoID, votacionID int64) (VistaVotacion, error) {
	v, integrantes, err := s.votacionConEquipo(ctx, proyectoID, votacionID)
	if err != nil {
		return VistaVotacion{}, err
	}
	if err := v.Revelar(integrantes); err != nil {
		return VistaVotacion{}, err
	}
	return s.guardar(ctx, v, integrantes)
}

// NuevaRonda empieza la ronda siguiente después de revelar (CA-6).
func (s *ServicioPlanningPoker) NuevaRonda(ctx context.Context, proyectoID, votacionID int64) (VistaVotacion, error) {
	v, integrantes, err := s.votacionConEquipo(ctx, proyectoID, votacionID)
	if err != nil {
		return VistaVotacion{}, err
	}
	if err := v.NuevaRonda(); err != nil {
		return VistaVotacion{}, err
	}
	return s.guardar(ctx, v, integrantes)
}

// Acordar registra la estimación acordada como Story Points de la historia y cierra la votación
// (CA-7, RN-9). Se estima con la operación de la spec 07, con sus reglas: si la historia se
// completó mientras tanto, se rechaza y la votación queda abierta.
func (s *ServicioPlanningPoker) Acordar(ctx context.Context, proyectoID, votacionID int64, valor int) (VistaVotacion, error) {
	v, integrantes, err := s.votacionConEquipo(ctx, proyectoID, votacionID)
	if err != nil {
		return VistaVotacion{}, err
	}
	if err := v.Acordar(valor); err != nil {
		return VistaVotacion{}, err
	}
	if _, err := s.backlog.Estimar(ctx, proyectoID, v.HistoriaID, valor); err != nil {
		return VistaVotacion{}, err
	}
	return s.guardar(ctx, v, integrantes)
}

// votacionConEquipo devuelve la votación del proyecto y sus integrantes (para saber quién es
// Developer y nombrar a quienes votaron).
func (s *ServicioPlanningPoker) votacionConEquipo(ctx context.Context, proyectoID, votacionID int64) (Votacion, []Integrante, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return Votacion{}, nil, err
	}
	v, err := s.votaciones.Obtener(ctx, proyectoID, votacionID)
	if err != nil {
		return Votacion{}, nil, err
	}
	integrantes, err := s.integrantes.Listar(ctx, proyectoID)
	if err != nil {
		return Votacion{}, nil, err
	}
	return v, integrantes, nil
}

func (s *ServicioPlanningPoker) guardar(ctx context.Context, v Votacion, integrantes []Integrante) (VistaVotacion, error) {
	if err := s.votaciones.Actualizar(ctx, v); err != nil {
		return VistaVotacion{}, err
	}
	return v.Vista(integrantes), nil
}
