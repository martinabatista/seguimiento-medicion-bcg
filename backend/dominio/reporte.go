package dominio

import (
	"slices"
	"time"
)

// DatosReporte son los datos que junta el servicio para armar un reporte (spec 14): los del
// proyecto y los cálculos ya hechos por las specs 03, 10, 11 y 12, sin cambios (RN-2).
type DatosReporte struct {
	Estado           EstadoDeProyecto // el proyecto y su estado (spec 03)
	Integrantes      []Integrante
	Historias        []Historia
	Sprints          []Sprint
	Avance           MetricasProyecto // spec 10
	Esfuerzo         MetricasEsfuerzo // spec 12
	Defectos         []Defecto        // spec 11, ya completados con la historia y los Sprints
	MetricasDefectos MetricasDefectos // spec 11 CA-4
	Hoy              time.Time
}

// HistoriaDeReporte es una historia con sus estimaciones y su esfuerzo (CA-1 a CA-3).
type HistoriaDeReporte struct {
	HistoriaID  int64
	Titulo      string
	Estado      EstadoHistoria
	StoryPoints *int        // nil = sin estimar; en un Sprint cerrado, los congelados al cerrarlo (RN-3)
	Estimadas   *int        // centésimas; nil = sin horas estimadas
	Reales      int         // centésimas: el total de su esfuerzo registrado (RN-4)
	Desviacion  *Desviacion // nil = "sin datos" (spec 12 RN-7)
	// CompletadaEnSprint solo se usa en el reporte de un Sprint: se completó en ese Sprint.
	CompletadaEnSprint bool
}

// SprintDeReporte es un Sprint con sus métricas de avance (spec 10) y de esfuerzo (spec 12).
type SprintDeReporte struct {
	Sprint   Sprint
	Avance   MetricasSprint
	Esfuerzo EsfuerzoSprint
}

// DefectosDeReporte son las cantidades y la lista de defectos (CA-5).
type DefectosDeReporte struct {
	Detectados int
	Resueltos  int
	Lista      []Defecto // en el orden en que se registraron
}

// ReporteProyecto es el reporte de un proyecto completo (spec 14 Salidas esperadas).
type ReporteProyecto struct {
	Proyecto    Proyecto
	Estado      EstadoProyecto
	Integrantes []Integrante
	GeneradoEl  time.Time // la fecha de hoy (RN-9)
	Historias   []HistoriaDeReporte
	Sprints     []SprintDeReporte // por fecha de inicio
	// Métricas del proyecto (spec 10 RN-6, RN-7 y spec 12 RN-9, RN-10).
	SprintsCerrados       int
	HayVelocidad          bool
	VelocidadDecimas      int
	Totales               TotalesSprints
	HorasEstimadas        int
	HorasReales           int
	Desviacion            *Desviacion
	CompletadasSinEstimar int
	Defectos              DefectosDeReporte
}

// ReporteSprint es el reporte de un Sprint (spec 14 Salidas esperadas).
type ReporteSprint struct {
	Proyecto   Proyecto
	GeneradoEl time.Time
	Sprint     SprintDeReporte // Avance.Parcial indica que está abierto (RN-6)
	Historias  []HistoriaDeReporte
	Defectos   DefectosDeReporte
}

// ArmarReporteProyecto junta los datos sin recalcularlos (RN-2): todo el backlog en su orden con
// los Story Points actuales y su esfuerzo, los Sprints por fecha de inicio, las métricas del
// proyecto y todos los defectos.
func ArmarReporteProyecto(d DatosReporte) ReporteProyecto {
	r := ReporteProyecto{
		Proyecto:              d.Estado.Proyecto,
		Estado:                d.Estado.Estado,
		Integrantes:           noNulo(d.Integrantes),
		GeneradoEl:            d.Hoy,
		Historias:             []HistoriaDeReporte{},
		Sprints:               []SprintDeReporte{},
		SprintsCerrados:       d.Avance.SprintsCerrados,
		HayVelocidad:          d.Avance.HayVelocidad,
		VelocidadDecimas:      d.Avance.VelocidadDecimas,
		Totales:               d.Avance.Totales,
		HorasEstimadas:        d.Esfuerzo.Estimadas,
		HorasReales:           d.Esfuerzo.Reales,
		Desviacion:            d.Esfuerzo.Desviacion,
		CompletadasSinEstimar: d.Esfuerzo.CompletadasSinEstimar,
		Defectos: DefectosDeReporte{
			Detectados: d.MetricasDefectos.Detectados,
			Resueltos:  d.MetricasDefectos.Resueltos,
			Lista:      noNulo(d.Defectos),
		},
	}
	esfuerzo := esfuerzoPorHistoria(d.Esfuerzo)
	for _, h := range OrdenarBacklog(d.Historias) {
		r.Historias = append(r.Historias, historiaDeReporte(h, h.StoryPoints, esfuerzo[h.ID]))
	}
	for _, s := range ordenarPorInicio(d.Sprints) {
		r.Sprints = append(r.Sprints, sprintDeReporte(d, s))
	}
	return r
}

// ArmarReporteSprint arma el reporte de un Sprint del proyecto (RN-3, RN-5, RN-6): sus historias
// planificadas en el orden del backlog, cuáles se completaron en él, sus Story Points (congelados
// si está cerrado), sus métricas y los defectos detectados o resueltos en él. Si el Sprint no es
// del proyecto, devuelve ErrSprintNoExiste (RN-1).
func ArmarReporteSprint(d DatosReporte, sprintID int64) (ReporteSprint, error) {
	i := slices.IndexFunc(d.Sprints, func(s Sprint) bool { return s.ID == sprintID })
	if i < 0 {
		return ReporteSprint{}, ErrSprintNoExiste
	}
	s := d.Sprints[i]
	r := ReporteSprint{
		Proyecto:   d.Estado.Proyecto,
		GeneradoEl: d.Hoy,
		Sprint:     sprintDeReporte(d, s),
		Historias:  []HistoriaDeReporte{},
		Defectos:   DefectosDeReporte{Lista: []Defecto{}},
	}
	esfuerzo := esfuerzoPorHistoria(d.Esfuerzo)
	for _, h := range OrdenarBacklog(d.Historias) {
		if !slices.Contains(s.Historias, h.ID) {
			continue
		}
		storyPoints := h.StoryPoints // Sprint abierto: los actuales (spec 10 RN-3)
		if s.Estado == SprintCerrado {
			storyPoints = s.StoryPointsAlCerrar[h.ID] // cerrado: los congelados (spec 10 RN-4)
		}
		hr := historiaDeReporte(h, storyPoints, esfuerzo[h.ID])
		hr.CompletadaEnSprint = slices.Contains(s.Completadas, h.ID)
		if !hr.CompletadaEnSprint {
			hr.Desviacion = nil // solo las completadas en el Sprint entran en la desviación (spec 12 RN-8)
		}
		r.Historias = append(r.Historias, hr)
	}
	for _, c := range d.MetricasDefectos.Sprints {
		if c.SprintID == s.ID {
			r.Defectos.Detectados, r.Defectos.Resueltos = c.Detectados, c.Resueltos
		}
	}
	for _, def := range d.Defectos {
		// SprintResolucionID solo está cargado si el defecto está Resuelto (spec 11 RN-9).
		if def.SprintDeteccionID == s.ID || def.SprintResolucionID == s.ID {
			r.Defectos.Lista = append(r.Defectos.Lista, def)
		}
	}
	return r, nil
}

// sprintDeReporte junta el Sprint con su fila de las métricas de avance y de esfuerzo.
func sprintDeReporte(d DatosReporte, s Sprint) SprintDeReporte {
	sr := SprintDeReporte{Sprint: s}
	for _, a := range d.Avance.Sprints {
		if a.SprintID == s.ID {
			sr.Avance = a
		}
	}
	for _, e := range d.Esfuerzo.Sprints {
		if e.SprintID == s.ID {
			sr.Esfuerzo = e
		}
	}
	return sr
}

func historiaDeReporte(h Historia, storyPoints *int, e EsfuerzoHistoria) HistoriaDeReporte {
	return HistoriaDeReporte{
		HistoriaID:  h.ID,
		Titulo:      h.Titulo,
		Estado:      h.Estado,
		StoryPoints: storyPoints,
		Estimadas:   h.HorasEstimadas,
		Reales:      e.Reales,
		Desviacion:  e.Desviacion,
	}
}

func esfuerzoPorHistoria(m MetricasEsfuerzo) map[int64]EsfuerzoHistoria {
	porID := make(map[int64]EsfuerzoHistoria, len(m.Historias))
	for _, e := range m.Historias {
		porID[e.HistoriaID] = e
	}
	return porID
}

// noNulo devuelve una lista vacía en lugar de nil, para que el JSON muestre [] (RN-7).
func noNulo[T any](lista []T) []T {
	if lista == nil {
		return []T{}
	}
	return lista
}
