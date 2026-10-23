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
		respuesta.Sprints = append(respuesta.Sprints, metricasSprintJSON{
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
		})
	}
	if m.HayVelocidad {
		// El dominio calcula en décimas, con enteros; recién acá se pasa a decimal para el JSON.
		v := float64(m.VelocidadDecimas) / 10
		respuesta.Velocidad = &v
	}
	writeJSON(w, http.StatusOK, respuesta)
}
