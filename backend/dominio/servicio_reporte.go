package dominio

import (
	"context"
	"time"
)

// ServicioReporte son los casos de uso de la spec 14. Llama a los servicios del estado (spec 03),
// de las métricas de avance (spec 10) y de esfuerzo (spec 12) y de los defectos (spec 11): no
// repite sus cálculos (RN-2), así el reporte cambia solo si cambia una de esas reglas.
type ServicioReporte struct {
	estado           *ServicioEstado
	metricas         *ServicioMetricas
	metricasEsfuerzo *ServicioMetricasEsfuerzo
	defectos         *ServicioDefectos
	integrantes      RepositorioIntegrantes
	historias        RepositorioHistorias
	sprints          RepositorioSprints
	hoy              func() time.Time
}

func NuevoServicioReporte(estado *ServicioEstado, metricas *ServicioMetricas, metricasEsfuerzo *ServicioMetricasEsfuerzo, defectos *ServicioDefectos, integrantes RepositorioIntegrantes, historias RepositorioHistorias, sprints RepositorioSprints, hoy func() time.Time) *ServicioReporte {
	return &ServicioReporte{estado: estado, metricas: metricas, metricasEsfuerzo: metricasEsfuerzo, defectos: defectos,
		integrantes: integrantes, historias: historias, sprints: sprints, hoy: hoy}
}

// ReporteProyecto arma el reporte del proyecto (CA-1 a CA-5). Si no existe, ErrProyectoNoExiste.
func (s *ServicioReporte) ReporteProyecto(ctx context.Context, proyectoID int64) (ReporteProyecto, error) {
	d, err := s.datos(ctx, proyectoID)
	if err != nil {
		return ReporteProyecto{}, err
	}
	return ArmarReporteProyecto(d), nil
}

// ReporteSprint arma el reporte de un Sprint del proyecto (CA-1 a CA-5). Si el proyecto no existe,
// ErrProyectoNoExiste; si el Sprint no es del proyecto, ErrSprintNoExiste (RN-1).
func (s *ServicioReporte) ReporteSprint(ctx context.Context, proyectoID, sprintID int64) (ReporteSprint, error) {
	d, err := s.datos(ctx, proyectoID)
	if err != nil {
		return ReporteSprint{}, err
	}
	return ArmarReporteSprint(d, sprintID)
}

// datos junta todo lo que usa un reporte, con la fecha de hoy (RN-9). La consulta del estado
// verifica que el proyecto exista.
func (s *ServicioReporte) datos(ctx context.Context, proyectoID int64) (DatosReporte, error) {
	estado, err := s.estado.Consultar(ctx, proyectoID)
	if err != nil {
		return DatosReporte{}, err
	}
	integrantes, err := s.integrantes.Listar(ctx, proyectoID)
	if err != nil {
		return DatosReporte{}, err
	}
	historias, err := s.historias.Listar(ctx, proyectoID)
	if err != nil {
		return DatosReporte{}, err
	}
	sprints, err := s.sprints.Listar(ctx, proyectoID)
	if err != nil {
		return DatosReporte{}, err
	}
	avance, err := s.metricas.Consultar(ctx, proyectoID)
	if err != nil {
		return DatosReporte{}, err
	}
	esfuerzo, err := s.metricasEsfuerzo.Consultar(ctx, proyectoID)
	if err != nil {
		return DatosReporte{}, err
	}
	defectos, err := s.defectos.Listar(ctx, proyectoID)
	if err != nil {
		return DatosReporte{}, err
	}
	metricasDefectos, err := s.defectos.Metricas(ctx, proyectoID)
	if err != nil {
		return DatosReporte{}, err
	}
	return DatosReporte{
		Estado:           estado,
		Integrantes:      integrantes,
		Historias:        historias,
		Sprints:          sprints,
		Avance:           avance,
		Esfuerzo:         esfuerzo,
		Defectos:         defectos,
		MetricasDefectos: metricasDefectos,
		Hoy:              s.hoy(),
	}, nil
}
