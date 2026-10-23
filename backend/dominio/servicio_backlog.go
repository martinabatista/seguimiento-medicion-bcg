package dominio

import "context"

// RepositorioHistorias guarda y lee las historias del backlog. Lo implementan postgres/ (base
// real) y memoria/ (tests). Las búsquedas por historia reciben también el proyecto: una
// historia de otro proyecto se trata como inexistente (spec 04 RN-10).
type RepositorioHistorias interface {
	// Guardar registra una historia nueva y la devuelve con el ID asignado.
	Guardar(ctx context.Context, h Historia) (Historia, error)
	// Actualizar reemplaza los datos de una historia existente, o devuelve ErrHistoriaNoExiste.
	Actualizar(ctx context.Context, h Historia) error
	// Eliminar borra la historia, o devuelve ErrHistoriaNoExiste.
	Eliminar(ctx context.Context, proyectoID, id int64) error
	// Listar devuelve las historias del proyecto (en cualquier orden: lo ordena el dominio).
	Listar(ctx context.Context, proyectoID int64) ([]Historia, error)
	// Obtener devuelve una historia del proyecto, o ErrHistoriaNoExiste.
	Obtener(ctx context.Context, proyectoID, id int64) (Historia, error)
}

// ServicioBacklog son los casos de uso de la spec 04. Usa el repositorio de proyectos solo para
// verificar que el proyecto exista (RN-1) y el de Sprints para no eliminar una historia de un
// Sprint abierto (spec 05 RN-11).
type ServicioBacklog struct {
	proyectos RepositorioProyectos
	historias RepositorioHistorias
	sprints   RepositorioSprints
}

func NuevoServicioBacklog(proyectos RepositorioProyectos, historias RepositorioHistorias, sprints RepositorioSprints) *ServicioBacklog {
	return &ServicioBacklog{proyectos: proyectos, historias: historias, sprints: sprints}
}

// Crear valida los datos y registra la historia en el backlog del proyecto (CA-1).
// Si el proyecto no existe o algún dato es inválido, no registra nada.
func (s *ServicioBacklog) Crear(ctx context.Context, proyectoID int64, d DatosHistoria) (Historia, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return Historia{}, err
	}
	h, err := NuevaHistoria(proyectoID, d)
	if err != nil {
		return Historia{}, err
	}
	return s.historias.Guardar(ctx, h)
}

// Modificar reemplaza los datos y el estado de una historia (CA-2). Primero verifica que el
// proyecto y la historia existan (CA-7); después valida con las mismas reglas que al crear y el
// estado (RN-7): si algo falla, la historia queda como estaba (RN-9).
func (s *ServicioBacklog) Modificar(ctx context.Context, proyectoID, id int64, d DatosHistoria, estado string) (Historia, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return Historia{}, err
	}
	actual, err := s.historias.Obtener(ctx, proyectoID, id)
	if err != nil {
		return Historia{}, err
	}
	h, err := NuevaHistoria(proyectoID, d)
	if err != nil {
		return Historia{}, err
	}
	if h.Estado, err = ParsearEstado(estado); err != nil {
		return Historia{}, err
	}
	// spec 06 RN-4: en un Sprint abierto, "Completada" solo se pone o se saca desde el Sprint.
	if cambiaCompletada(actual.Estado, h.Estado) {
		enSprint, err := enSprintAbierto(ctx, s.sprints, proyectoID, id)
		if err != nil {
			return Historia{}, err
		}
		if enSprint {
			return Historia{}, ErrCompletarDesdeSprint
		}
	}
	h.ID = id                          // RN-10: el identificador no cambia.
	h.StoryPoints = actual.StoryPoints // Los Story Points no se modifican acá (specs 07 y 08).
	if err := s.historias.Actualizar(ctx, h); err != nil {
		return Historia{}, err
	}
	return h, nil
}

// Eliminar saca la historia del backlog (CA-3). No se puede si está asignada a un Sprint abierto:
// primero hay que quitarla del Sprint (spec 05 RN-11, CA-6).
func (s *ServicioBacklog) Eliminar(ctx context.Context, proyectoID, id int64) error {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return err
	}
	enSprint, err := enSprintAbierto(ctx, s.sprints, proyectoID, id)
	if err != nil {
		return err
	}
	if enSprint {
		return ErrHistoriaEnSprintAbierto
	}
	return s.historias.Eliminar(ctx, proyectoID, id)
}

// Listar devuelve el backlog del proyecto ordenado por prioridad y creación (CA-4, RN-8).
func (s *ServicioBacklog) Listar(ctx context.Context, proyectoID int64) ([]Historia, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return nil, err
	}
	historias, err := s.historias.Listar(ctx, proyectoID)
	if err != nil {
		return nil, err
	}
	return OrdenarBacklog(historias), nil
}

// Obtener devuelve una historia del proyecto (CA-4) o ErrHistoriaNoExiste (CA-7, RN-10).
func (s *ServicioBacklog) Obtener(ctx context.Context, proyectoID, id int64) (Historia, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return Historia{}, err
	}
	return s.historias.Obtener(ctx, proyectoID, id)
}

// cambiaCompletada indica si el cambio de estado pone o saca a la historia de "Completada"
// (spec 06 RN-4). Pasar de Pendiente a En progreso, o al revés, no cuenta.
func cambiaCompletada(antes, despues EstadoHistoria) bool {
	return antes != despues && (antes == EstadoCompletada || despues == EstadoCompletada)
}
