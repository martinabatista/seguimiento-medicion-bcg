package api

import (
	"net/http"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// estadoJSON es el estado de un proyecto tal como sale por la API (spec 03).
type estadoJSON struct {
	Proyecto proyectoJSON       `json:"proyecto"`
	Estado   string             `json:"estado"`
	Backlog  resumenBacklogJSON `json:"backlog"`
}

// resumenBacklogJSON es el resumen del avance del backlog (spec 03 CA-3).
type resumenBacklogJSON struct {
	Total                  int `json:"total"`
	Pendientes             int `json:"pendientes"`
	EnProgreso             int `json:"en_progreso"`
	Completadas            int `json:"completadas"`
	PorcentajeCompletadas  int `json:"porcentaje_completadas"`
	StoryPointsTotales     int `json:"story_points_totales"`
	StoryPointsCompletados int `json:"story_points_completados"`
	SinEstimar             int `json:"sin_estimar"`
}

type estadoHandler struct {
	servicio *dominio.ServicioEstado
}

// consultar: GET /api/proyectos/{id}/estado (CA-1, CA-3, CA-4).
func (h estadoHandler) consultar(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	e, err := h.servicio.Consultar(r.Context(), proyectoID)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, estadoJSON{
		Proyecto: aJSON(e.Proyecto),
		Estado:   string(e.Estado),
		Backlog: resumenBacklogJSON{
			Total:                  e.Backlog.Total,
			Pendientes:             e.Backlog.Pendientes,
			EnProgreso:             e.Backlog.EnProgreso,
			Completadas:            e.Backlog.Completadas,
			PorcentajeCompletadas:  e.Backlog.PorcentajeCompletadas,
			StoryPointsTotales:     e.Backlog.StoryPointsTotales,
			StoryPointsCompletados: e.Backlog.StoryPointsCompletados,
			SinEstimar:             e.Backlog.SinEstimar,
		},
	})
}
