package dominio

import (
	"cmp"
	"context"
	"slices"
	"time"
)

// RepositorioSprints guarda y lee los Sprints. Lo implementan postgres/ (base real) y memoria/
// (tests). Las búsquedas por Sprint reciben también el proyecto: un Sprint de otro proyecto se
// trata como inexistente (spec 05 RN-13).
type RepositorioSprints interface {
	// Guardar registra un Sprint nuevo y lo devuelve con el ID asignado.
	Guardar(ctx context.Context, s Sprint) (Sprint, error)
	// Actualizar reemplaza el Sprint Goal, el estado y las historias asignadas.
	Actualizar(ctx context.Context, s Sprint) error
	// Listar devuelve los Sprints del proyecto, con sus historias (en cualquier orden).
	Listar(ctx context.Context, proyectoID int64) ([]Sprint, error)
	// Obtener devuelve un Sprint del proyecto, o ErrSprintNoExiste.
	Obtener(ctx context.Context, proyectoID, id int64) (Sprint, error)
}

// DetalleSprint es un Sprint con sus historias, ordenadas como el backlog (spec 05 CA-3).
type DetalleSprint struct {
	Sprint    Sprint
	Historias []Historia
}

// ServicioSprints son los casos de uso de la spec 05.
type ServicioSprints struct {
	proyectos RepositorioProyectos
	historias RepositorioHistorias
	sprints   RepositorioSprints
}

func NuevoServicioSprints(proyectos RepositorioProyectos, historias RepositorioHistorias, sprints RepositorioSprints) *ServicioSprints {
	return &ServicioSprints{proyectos: proyectos, historias: historias, sprints: sprints}
}

// Crear valida el Sprint contra el proyecto y sus otros Sprints, y lo registra abierto (CA-1).
func (s *ServicioSprints) Crear(ctx context.Context, proyectoID int64, inicio, fin time.Time, goal string) (Sprint, error) {
	p, err := s.proyectos.Obtener(ctx, proyectoID)
	if err != nil {
		return Sprint{}, err
	}
	existentes, err := s.sprints.Listar(ctx, proyectoID)
	if err != nil {
		return Sprint{}, err
	}
	nuevo, err := NuevoSprint(p, inicio, fin, goal, existentes)
	if err != nil {
		return Sprint{}, err
	}
	return s.sprints.Guardar(ctx, nuevo)
}

// ModificarGoal cambia el Sprint Goal de un Sprint abierto (CA-2, RN-5, RN-7).
func (s *ServicioSprints) ModificarGoal(ctx context.Context, proyectoID, sprintID int64, goal string) (DetalleSprint, error) {
	sp, err := s.sprintAbierto(ctx, proyectoID, sprintID)
	if err != nil {
		return DetalleSprint{}, err
	}
	if err := ValidarSprintGoal(goal); err != nil {
		return DetalleSprint{}, err
	}
	sp.Goal = goal
	return s.guardarYDetallar(ctx, sp)
}

// AsignarHistoria suma una historia del backlog al Sprint (CA-3, CA-4, RN-8, RN-9).
func (s *ServicioSprints) AsignarHistoria(ctx context.Context, proyectoID, sprintID, historiaID int64) (DetalleSprint, error) {
	sp, err := s.sprintAbierto(ctx, proyectoID, sprintID)
	if err != nil {
		return DetalleSprint{}, err
	}
	h, err := s.historias.Obtener(ctx, proyectoID, historiaID)
	if err != nil {
		return DetalleSprint{}, err
	}
	if h.Estado == EstadoCompletada {
		return DetalleSprint{}, ErrHistoriaCompletada
	}
	enOtro, err := enSprintAbierto(ctx, s.sprints, proyectoID, historiaID)
	if err != nil {
		return DetalleSprint{}, err
	}
	if enOtro {
		return DetalleSprint{}, ErrHistoriaYaAsignada
	}
	sp.Historias = append(sp.Historias, historiaID)
	return s.guardarYDetallar(ctx, sp)
}

// QuitarHistoria saca una historia del Sprint (CA-3, RN-10).
func (s *ServicioSprints) QuitarHistoria(ctx context.Context, proyectoID, sprintID, historiaID int64) (DetalleSprint, error) {
	sp, err := s.sprintAbierto(ctx, proyectoID, sprintID)
	if err != nil {
		return DetalleSprint{}, err
	}
	if !sp.Tiene(historiaID) {
		return DetalleSprint{}, ErrHistoriaNoAsignada
	}
	sp.Historias = slices.DeleteFunc(sp.Historias, func(id int64) bool { return id == historiaID })
	sp.Completadas = slices.DeleteFunc(sp.Completadas, func(id int64) bool { return id == historiaID })
	return s.guardarYDetallar(ctx, sp)
}

// CompletarHistoria marca como completada una historia del Sprint abierto: pasa a Completada en el
// backlog y el Sprint registra que se completó en él (spec 06 CA-1, RN-1, RN-2).
func (s *ServicioSprints) CompletarHistoria(ctx context.Context, proyectoID, sprintID, historiaID int64) (DetalleSprint, error) {
	sp, h, err := s.historiaDelSprint(ctx, proyectoID, sprintID, historiaID)
	if err != nil {
		return DetalleSprint{}, err
	}
	if sp.EstaCompletada(historiaID) {
		return DetalleSprint{}, ErrHistoriaCompletada
	}
	h.Estado = EstadoCompletada
	sp.Completadas = append(sp.Completadas, historiaID)
	return s.guardarHistoriaYSprint(ctx, h, sp)
}

// DeshacerHistoria vuelve a "En progreso" una historia completada en el Sprint abierto
// (spec 06 CA-1, RN-3).
func (s *ServicioSprints) DeshacerHistoria(ctx context.Context, proyectoID, sprintID, historiaID int64) (DetalleSprint, error) {
	sp, h, err := s.historiaDelSprint(ctx, proyectoID, sprintID, historiaID)
	if err != nil {
		return DetalleSprint{}, err
	}
	if !sp.EstaCompletada(historiaID) {
		return DetalleSprint{}, ErrHistoriaNoCompletada
	}
	h.Estado = EstadoEnProgreso
	sp.Completadas = slices.DeleteFunc(sp.Completadas, func(id int64) bool { return id == historiaID })
	return s.guardarHistoriaYSprint(ctx, h, sp)
}

// Cerrar pasa el Sprint a Cerrado en cualquier momento (spec 06 CA-2, RN-5). Conserva sus historias
// y cuáles se completaron; las no completadas mantienen su estado y quedan libres (RN-6).
func (s *ServicioSprints) Cerrar(ctx context.Context, proyectoID, sprintID int64) (DetalleSprint, error) {
	sp, err := s.sprintAbierto(ctx, proyectoID, sprintID)
	if err != nil {
		return DetalleSprint{}, err
	}
	sp.Estado = SprintCerrado
	return s.guardarYDetallar(ctx, sp)
}

// Anteriores devuelve los Sprints cerrados del proyecto con sus historias, del más reciente al más
// antiguo por fecha de inicio (spec 06 CA-4, RN-9).
func (s *ServicioSprints) Anteriores(ctx context.Context, proyectoID int64) ([]DetalleSprint, error) {
	sprints, err := s.Listar(ctx, proyectoID)
	if err != nil {
		return nil, err
	}
	anteriores := []DetalleSprint{}
	for _, sp := range slices.Backward(sprints) {
		if sp.Estado != SprintCerrado {
			continue
		}
		d, err := s.detallar(ctx, sp)
		if err != nil {
			return nil, err
		}
		anteriores = append(anteriores, d)
	}
	return anteriores, nil
}

// historiaDelSprint busca el Sprint abierto y la historia, y verifica que esté asignada a él
// (spec 06 RN-1).
func (s *ServicioSprints) historiaDelSprint(ctx context.Context, proyectoID, sprintID, historiaID int64) (Sprint, Historia, error) {
	sp, err := s.sprintAbierto(ctx, proyectoID, sprintID)
	if err != nil {
		return Sprint{}, Historia{}, err
	}
	h, err := s.historias.Obtener(ctx, proyectoID, historiaID)
	if err != nil {
		return Sprint{}, Historia{}, err
	}
	if !sp.Tiene(historiaID) {
		return Sprint{}, Historia{}, ErrHistoriaNoAsignada
	}
	return sp, h, nil
}

func (s *ServicioSprints) guardarHistoriaYSprint(ctx context.Context, h Historia, sp Sprint) (DetalleSprint, error) {
	if err := s.historias.Actualizar(ctx, h); err != nil {
		return DetalleSprint{}, err
	}
	return s.guardarYDetallar(ctx, sp)
}

// Listar devuelve los Sprints del proyecto ordenados por fecha de inicio (CA-3).
func (s *ServicioSprints) Listar(ctx context.Context, proyectoID int64) ([]Sprint, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return nil, err
	}
	sprints, err := s.sprints.Listar(ctx, proyectoID)
	if err != nil {
		return nil, err
	}
	return ordenarPorInicio(sprints), nil
}

// Obtener devuelve un Sprint con sus historias (CA-3) o ErrSprintNoExiste (CA-7).
func (s *ServicioSprints) Obtener(ctx context.Context, proyectoID, sprintID int64) (DetalleSprint, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return DetalleSprint{}, err
	}
	sp, err := s.sprints.Obtener(ctx, proyectoID, sprintID)
	if err != nil {
		return DetalleSprint{}, err
	}
	return s.detallar(ctx, sp)
}

// sprintAbierto busca el Sprint (CA-7) y verifica que admita cambios (RN-7).
func (s *ServicioSprints) sprintAbierto(ctx context.Context, proyectoID, sprintID int64) (Sprint, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return Sprint{}, err
	}
	sp, err := s.sprints.Obtener(ctx, proyectoID, sprintID)
	if err != nil {
		return Sprint{}, err
	}
	if sp.Estado == SprintCerrado {
		return Sprint{}, ErrSprintCerrado
	}
	return sp, nil
}

func (s *ServicioSprints) guardarYDetallar(ctx context.Context, sp Sprint) (DetalleSprint, error) {
	if err := s.sprints.Actualizar(ctx, sp); err != nil {
		return DetalleSprint{}, err
	}
	return s.detallar(ctx, sp)
}

// detallar busca las historias del Sprint y las ordena como el backlog (spec 04 RN-8).
func (s *ServicioSprints) detallar(ctx context.Context, sp Sprint) (DetalleSprint, error) {
	historias := make([]Historia, 0, len(sp.Historias))
	for _, id := range sp.Historias {
		h, err := s.historias.Obtener(ctx, sp.ProyectoID, id)
		if err != nil {
			return DetalleSprint{}, err
		}
		historias = append(historias, h)
	}
	return DetalleSprint{Sprint: sp, Historias: OrdenarBacklog(historias)}, nil
}

// enSprintAbierto indica si la historia está asignada a algún Sprint abierto del proyecto
// (RN-9 al asignar y RN-11 al eliminar del backlog).
func enSprintAbierto(ctx context.Context, sprints RepositorioSprints, proyectoID, historiaID int64) (bool, error) {
	lista, err := sprints.Listar(ctx, proyectoID)
	if err != nil {
		return false, err
	}
	for _, sp := range lista {
		if sp.Estado == SprintAbierto && sp.Tiene(historiaID) {
			return true, nil
		}
	}
	return false, nil
}

func ordenarPorInicio(sprints []Sprint) []Sprint {
	ordenados := slices.Clone(sprints)
	slices.SortStableFunc(ordenados, func(a, b Sprint) int { return cmp.Compare(a.FechaInicio.Unix(), b.FechaInicio.Unix()) })
	return ordenados
}
