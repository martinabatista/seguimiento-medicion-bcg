package dominio

import "context"

// ServicioMetricasEsfuerzo es el caso de uso de la spec 12: consultar las métricas de esfuerzo de un
// proyecto. Lee las historias, los Sprints y el esfuerzo registrado; no guarda nada.
type ServicioMetricasEsfuerzo struct {
	proyectos RepositorioProyectos
	historias RepositorioHistorias
	sprints   RepositorioSprints
	esfuerzo  RepositorioEsfuerzo
}

func NuevoServicioMetricasEsfuerzo(proyectos RepositorioProyectos, historias RepositorioHistorias, sprints RepositorioSprints, esfuerzo RepositorioEsfuerzo) *ServicioMetricasEsfuerzo {
	return &ServicioMetricasEsfuerzo{proyectos: proyectos, historias: historias, sprints: sprints, esfuerzo: esfuerzo}
}

// Consultar calcula las métricas de esfuerzo en el momento de la consulta. Si el proyecto no existe,
// ErrProyectoNoExiste (CA-7).
func (s *ServicioMetricasEsfuerzo) Consultar(ctx context.Context, proyectoID int64) (MetricasEsfuerzo, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return MetricasEsfuerzo{}, err
	}
	historias, err := s.historias.Listar(ctx, proyectoID)
	if err != nil {
		return MetricasEsfuerzo{}, err
	}
	sprints, err := s.sprints.Listar(ctx, proyectoID)
	if err != nil {
		return MetricasEsfuerzo{}, err
	}
	// Las horas reales de cada historia son el total de sus registros (RN-5, spec 09 RN-6).
	reales := make(map[int64]int, len(historias))
	for _, h := range historias {
		registros, err := s.esfuerzo.Listar(ctx, proyectoID, h.ID)
		if err != nil {
			return MetricasEsfuerzo{}, err
		}
		reales[h.ID] = ResumirEsfuerzo(registros).TotalCentesimas
	}
	return CalcularMetricasEsfuerzo(historias, sprints, reales), nil
}
