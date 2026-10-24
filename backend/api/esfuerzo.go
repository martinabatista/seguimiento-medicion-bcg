package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// registroEsfuerzoJSON es un registro de esfuerzo tal como sale por la API (spec 09). Las horas
// son un número con hasta 2 decimales (2.5), armado desde las centésimas para que sea exacto.
type registroEsfuerzoJSON struct {
	ID               int64       `json:"id"`
	ProyectoID       int64       `json:"proyecto_id"`
	HistoriaID       int64       `json:"historia_id"`
	IntegranteID     int64       `json:"integrante_id"`
	IntegranteNombre string      `json:"integrante_nombre"`
	Fecha            string      `json:"fecha"`
	Actividad        string      `json:"actividad"`
	Horas            json.Number `json:"horas"`
}

// esfuerzoJSON es el esfuerzo de una historia: sus registros por fecha y el total (CA-4).
type esfuerzoJSON struct {
	Registros  []registroEsfuerzoJSON `json:"registros"`
	TotalHoras json.Number            `json:"total_horas"`
}

// esfuerzoEntrada es lo que llega al registrar. Las horas se leen como json.Number para que el
// dominio distinga "faltan" de "no son válidas" (spec 09 RN-5).
type esfuerzoEntrada struct {
	IntegranteID int64       `json:"integrante_id"`
	Fecha        string      `json:"fecha"`
	Actividad    string      `json:"actividad"`
	Horas        json.Number `json:"horas"`
}

type esfuerzoHandler struct {
	servicio *dominio.ServicioEsfuerzo
}

// registrar: POST /api/proyectos/{id}/backlog/{historiaID}/esfuerzo (CA-1, CA-2).
func (h esfuerzoHandler) registrar(w http.ResponseWriter, r *http.Request) {
	proyectoID, historiaID, err := leerIDs(r)
	if err != nil {
		responderError(w, err)
		return
	}
	var entrada esfuerzoEntrada
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, errCuerpoInvalido)
		return
	}
	registro, err := h.servicio.Registrar(r.Context(), proyectoID, historiaID, dominio.DatosEsfuerzo{
		IntegranteID: entrada.IntegranteID,
		Fecha:        entrada.Fecha,
		Actividad:    entrada.Actividad,
		Horas:        entrada.Horas.String(),
	})
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, registroEsfuerzoAJSON(registro))
}

// consultar: GET /api/proyectos/{id}/backlog/{historiaID}/esfuerzo (CA-4).
func (h esfuerzoHandler) consultar(w http.ResponseWriter, r *http.Request) {
	proyectoID, historiaID, err := leerIDs(r)
	if err != nil {
		responderError(w, err)
		return
	}
	e, err := h.servicio.Consultar(r.Context(), proyectoID, historiaID)
	if err != nil {
		responderError(w, err)
		return
	}
	respuesta := esfuerzoJSON{
		Registros:  make([]registroEsfuerzoJSON, 0, len(e.Registros)),
		TotalHoras: json.Number(dominio.HorasATexto(e.TotalCentesimas)),
	}
	for _, registro := range e.Registros {
		respuesta.Registros = append(respuesta.Registros, registroEsfuerzoAJSON(registro))
	}
	writeJSON(w, http.StatusOK, respuesta)
}

// eliminar: DELETE /api/proyectos/{id}/backlog/{historiaID}/esfuerzo/{registroID} (CA-6).
// Responde 204 sin cuerpo.
func (h esfuerzoHandler) eliminar(w http.ResponseWriter, r *http.Request) {
	proyectoID, historiaID, err := leerIDs(r)
	if err != nil {
		responderError(w, err)
		return
	}
	// Un id que no es un número no puede existir (CA-8).
	registroID, err := strconv.ParseInt(r.PathValue("registroID"), 10, 64)
	if err != nil {
		responderError(w, dominio.ErrRegistroEsfuerzoNoExiste)
		return
	}
	if err := h.servicio.Eliminar(r.Context(), proyectoID, historiaID, registroID); err != nil {
		responderError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func registroEsfuerzoAJSON(e dominio.RegistroEsfuerzo) registroEsfuerzoJSON {
	return registroEsfuerzoJSON{
		ID:               e.ID,
		ProyectoID:       e.ProyectoID,
		HistoriaID:       e.HistoriaID,
		IntegranteID:     e.IntegranteID,
		IntegranteNombre: e.IntegranteNombre,
		Fecha:            e.Fecha.Format(dominio.FormatoFecha),
		Actividad:        e.Actividad,
		Horas:            json.Number(dominio.HorasATexto(e.Centesimas)),
	}
}
