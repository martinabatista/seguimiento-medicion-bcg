package api

import (
	"encoding/json"
	"net/http"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// estimacionEntrada es lo que llega al estimar. Se lee como json.Number para que el dominio
// distinga "falta el valor" de "no es un valor de la escala" (spec 07 RN-2).
type estimacionEntrada struct {
	StoryPoints json.Number `json:"story_points"`
}

// escalaJSON es la escala de Story Points tal como sale por la API (spec 07 CA-5).
type escalaJSON struct {
	Valores []int `json:"valores"`
}

// estimar: PUT /api/proyectos/{id}/backlog/{historiaID}/story-points (spec 07 CA-1, CA-2).
func (h backlogHandler) estimar(w http.ResponseWriter, r *http.Request) {
	proyectoID, historiaID, err := leerIDs(r)
	if err != nil {
		responderError(w, err)
		return
	}
	var entrada estimacionEntrada
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, errCuerpoInvalido)
		return
	}
	storyPoints, err := dominio.ParsearStoryPoints(entrada.StoryPoints.String())
	if err != nil {
		responderError(w, err)
		return
	}
	hist, err := h.servicio.Estimar(r.Context(), proyectoID, historiaID, storyPoints)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, historiaAJSON(hist))
}

// quitarEstimacion: DELETE /api/proyectos/{id}/backlog/{historiaID}/story-points (spec 07 CA-2).
// Devuelve la historia sin estimar.
func (h backlogHandler) quitarEstimacion(w http.ResponseWriter, r *http.Request) {
	proyectoID, historiaID, err := leerIDs(r)
	if err != nil {
		responderError(w, err)
		return
	}
	hist, err := h.servicio.QuitarEstimacion(r.Context(), proyectoID, historiaID)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, historiaAJSON(hist))
}

// escalaStoryPoints: GET /api/story-points/escala (spec 07 CA-5). La escala la define el dominio:
// el frontend y el Planning Poker (spec 08) la consultan acá en vez de repetirla.
func escalaStoryPoints(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, escalaJSON{Valores: dominio.EscalaStoryPoints()})
}
