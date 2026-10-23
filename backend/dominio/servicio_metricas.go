package dominio

import "context"

// ServicioMetricas es el caso de uso de la spec 10: consultar las métricas de avance de un proyecto.
type ServicioMetricas struct {
	proyectos RepositorioProyectos
	historias RepositorioHistorias
	sprints   RepositorioSprints
}

func NuevoServicioMetricas(proyectos RepositorioProyectos, historias RepositorioHistorias, sprints RepositorioSprints) *ServicioMetricas {
	return &ServicioMetricas{proyectos: proyectos, historias: historias, sprints: sprints}
}

// Consultar calcula las métricas de cada Sprint del proyecto (por fecha de inicio) y las del
// proyecto, en el momento de la consulta (RN-8). Si el proyecto no existe, ErrProyectoNoExiste (CA-7).
func (s *ServicioMetricas) Consultar(ctx context.Context, proyectoID int64) (MetricasProyecto, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return MetricasProyecto{}, err
	}
	sprints, err := s.sprints.Listar(ctx, proyectoID)
	if err != nil {
		return MetricasProyecto{}, err
	}
	actuales, err := s.puntosActuales(ctx, proyectoID)
	if err != nil {
		return MetricasProyecto{}, err
	}
	metricas := []MetricasSprint{}
	for _, sp := range ordenarPorInicio(sprints) {
		puntos := actuales // Sprint abierto: Story Points de hoy (RN-3)
		if sp.Estado == SprintCerrado {
			puntos = sp.StoryPointsAlCerrar // Sprint cerrado: los congelados al cerrarlo (RN-4)
		}
		metricas = append(metricas, CalcularMetricasSprint(sp, puntos))
	}
	return CalcularMetricasProyecto(metricas), nil
}

// puntosActuales devuelve los Story Points de hoy de cada historia del backlog (nil = sin estimar).
func (s *ServicioMetricas) puntosActuales(ctx context.Context, proyectoID int64) (map[int64]*int, error) {
	historias, err := s.historias.Listar(ctx, proyectoID)
	if err != nil {
		return nil, err
	}
	return puntosDe(historias), nil
}

// puntosDe arma el mapa historia → Story Points (nil = sin estimar).
func puntosDe(historias []Historia) map[int64]*int {
	puntos := make(map[int64]*int, len(historias))
	for _, h := range historias {
		puntos[h.ID] = h.StoryPoints
	}
	return puntos
}
