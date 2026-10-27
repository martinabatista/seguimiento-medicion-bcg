package api

import (
	"net/http"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// metricasSprintJSON son las métricas de un Sprint tal como salen por la API (spec 10).
type metricasSprintJSON struct {
	SprintID                int64  `json:"sprint_id"`
	Numero                  int    `json:"numero"`
	Estado                  string `json:"estado"`
	Parcial                 bool   `json:"parcial"`
	HistoriasPlanificadas   int    `json:"historias_planificadas"`
	HistoriasCompletadas    int    `json:"historias_completadas"`
	PorcentajeCompletadas   int    `json:"porcentaje_completadas"`
	SinEstimar              int    `json:"sin_estimar"`
	StoryPointsPlanificados int    `json:"story_points_planificados"`
	StoryPointsCompletados  int    `json:"story_points_completados"`
}

// totalesJSON son los totales de los Sprints cerrados (spec 10 RN-7).
type totalesJSON struct {
	HistoriasPlanificadas   int `json:"historias_planificadas"`
	HistoriasCompletadas    int `json:"historias_completadas"`
	PorcentajeCompletadas   int `json:"porcentaje_completadas"`
	StoryPointsPlanificados int `json:"story_points_planificados"`
	StoryPointsCompletados  int `json:"story_points_completados"`
}

// metricasJSON son las métricas de un proyecto. velocidad es null sin Sprints cerrados (RN-6).
type metricasJSON struct {
	Sprints         []metricasSprintJSON `json:"sprints"`
	SprintsCerrados int                  `json:"sprints_cerrados"`
	Velocidad       *float64             `json:"velocidad"`
	Totales         totalesJSON          `json:"totales"`
}

type metricasHandler struct {
	servicio *dominio.ServicioMetricas
}

// consultar: GET /api/proyectos/{id}/metricas (spec 10 CA-1 a CA-4, CA-6, CA-7).
func (h metricasHandler) consultar(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	m, err := h.servicio.Consultar(r.Context(), proyectoID)
	if err != nil {
		responderError(w, err)
		return
	}
	respuesta := metricasJSON{
		Sprints:         make([]metricasSprintJSON, 0, len(m.Sprints)),
		SprintsCerrados: m.SprintsCerrados,
		Totales:         totalesJSON(m.Totales),
	}
	for _, s := range m.Sprints {
		respuesta.Sprints = append(respuesta.Sprints, metricasSprintAJSON(s))
	}
	respuesta.Velocidad = velocidadAJSON(m.HayVelocidad, m.VelocidadDecimas)
	writeJSON(w, http.StatusOK, respuesta)
}

// velocidadAJSON pasa la velocidad a decimal para el JSON (6.5), o nil si no hay (spec 10 RN-6). El
// dominio la calcula en décimas, con enteros; recién acá se convierte. La usan las métricas y el
// Dashboard (spec 13), así la velocidad sale igual en las dos rutas.
func velocidadAJSON(hay bool, decimas int) *float64 {
	if !hay {
		return nil
	}
	v := float64(decimas) / 10
	return &v
}

// metricasSprintAJSON pasa las métricas de avance de un Sprint al JSON. La usan las métricas y el
// reporte (spec 14), así un Sprint sale igual en las dos rutas.
func metricasSprintAJSON(s dominio.MetricasSprint) metricasSprintJSON {
	return metricasSprintJSON{
		SprintID:                s.SprintID,
		Numero:                  s.Numero,
		Estado:                  string(s.Estado),
		Parcial:                 s.Parcial,
		HistoriasPlanificadas:   s.HistoriasPlanificadas,
		HistoriasCompletadas:    s.HistoriasCompletadas,
		PorcentajeCompletadas:   s.PorcentajeCompletadas,
		SinEstimar:              s.SinEstimar,
		StoryPointsPlanificados: s.StoryPointsPlanificados,
		StoryPointsCompletados:  s.StoryPointsCompletados,
	}
}
