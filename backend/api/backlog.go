package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// historiaJSON es una historia del backlog tal como sale por la API (spec 04).
// story_points es null mientras la historia está sin estimar; horas_estimadas, mientras no tiene
// horas estimadas (spec 12).
type historiaJSON struct {
	ID             int64        `json:"id"`
	ProyectoID     int64        `json:"proyecto_id"`
	Titulo         string       `json:"titulo"`
	Descripcion    string       `json:"descripcion"`
	Prioridad      int          `json:"prioridad"`
	Criterios      string       `json:"criterios_aceptacion"`
	Estado         string       `json:"estado"`
	StoryPoints    *int         `json:"story_points"`
	HorasEstimadas *json.Number `json:"horas_estimadas"`
}

// historiaEntrada es lo que llega al crear o modificar. La prioridad se lee como json.Number para
// que el dominio distinga "falta" de "no es un entero" (spec 04 RN-4). El estado solo se usa al
// modificar: al crear, la historia siempre arranca Pendiente (RN-6).
type historiaEntrada struct {
	Titulo      string      `json:"titulo"`
	Descripcion string      `json:"descripcion"`
	Prioridad   json.Number `json:"prioridad"`
	Criterios   string      `json:"criterios_aceptacion"`
	Estado      string      `json:"estado"`
}

type backlogHandler struct {
	servicio *dominio.ServicioBacklog
}

// crear: POST /api/proyectos/{id}/backlog (CA-1).
func (h backlogHandler) crear(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	_, datos, err := leerHistoria(r) // al crear se ignora el estado (RN-6)
	if err != nil {
		responderError(w, err)
		return
	}
	hist, err := h.servicio.Crear(r.Context(), proyectoID, datos)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, historiaAJSON(hist))
}

// listar: GET /api/proyectos/{id}/backlog (CA-4), ordenado por prioridad.
func (h backlogHandler) listar(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	historias, err := h.servicio.Listar(r.Context(), proyectoID)
	if err != nil {
		responderError(w, err)
		return
	}
	lista := make([]historiaJSON, 0, len(historias))
	for _, hist := range historias {
		lista = append(lista, historiaAJSON(hist))
	}
	writeJSON(w, http.StatusOK, lista)
}

// ver: GET /api/proyectos/{id}/backlog/{historiaID} (CA-4, CA-7).
func (h backlogHandler) ver(w http.ResponseWriter, r *http.Request) {
	proyectoID, historiaID, err := leerIDs(r)
	if err != nil {
		responderError(w, err)
		return
	}
	hist, err := h.servicio.Obtener(r.Context(), proyectoID, historiaID)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, historiaAJSON(hist))
}

// modificar: PUT /api/proyectos/{id}/backlog/{historiaID} (CA-2, CA-7). Reemplaza los datos y el estado.
func (h backlogHandler) modificar(w http.ResponseWriter, r *http.Request) {
	proyectoID, historiaID, err := leerIDs(r)
	if err != nil {
		responderError(w, err)
		return
	}
	entrada, datos, err := leerHistoria(r)
	if err != nil {
		responderError(w, err)
		return
	}
	hist, err := h.servicio.Modificar(r.Context(), proyectoID, historiaID, datos, entrada.Estado)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, historiaAJSON(hist))
}

// eliminar: DELETE /api/proyectos/{id}/backlog/{historiaID} (CA-3). Responde 204 sin cuerpo.
func (h backlogHandler) eliminar(w http.ResponseWriter, r *http.Request) {
	proyectoID, historiaID, err := leerIDs(r)
	if err != nil {
		responderError(w, err)
		return
	}
	if err := h.servicio.Eliminar(r.Context(), proyectoID, historiaID); err != nil {
		responderError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// leerHistoria decodifica el cuerpo y convierte la prioridad con la regla del dominio.
func leerHistoria(r *http.Request) (historiaEntrada, dominio.DatosHistoria, error) {
	var entrada historiaEntrada
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		return entrada, dominio.DatosHistoria{}, errCuerpoInvalido
	}
	prioridad, err := dominio.ParsearPrioridad(entrada.Prioridad.String())
	if err != nil {
		return entrada, dominio.DatosHistoria{}, err
	}
	return entrada, dominio.DatosHistoria{
		Titulo:      entrada.Titulo,
		Descripcion: entrada.Descripcion,
		Prioridad:   prioridad,
		Criterios:   entrada.Criterios,
	}, nil
}

// leerIDs lee el {id} del proyecto y el {historiaID} de la ruta. Un id que no es un número no
// puede existir (CA-7).
func leerIDs(r *http.Request) (proyectoID, historiaID int64, err error) {
	proyectoID, err = leerID(r)
	if err != nil {
		return 0, 0, err
	}
	historiaID, err = strconv.ParseInt(r.PathValue("historiaID"), 10, 64)
	if err != nil {
		return 0, 0, dominio.ErrHistoriaNoExiste
	}
	return proyectoID, historiaID, nil
}

func historiaAJSON(h dominio.Historia) historiaJSON {
	return historiaJSON{
		ID:             h.ID,
		ProyectoID:     h.ProyectoID,
		Titulo:         h.Titulo,
		Descripcion:    h.Descripcion,
		Prioridad:      h.Prioridad,
		Criterios:      h.Criterios,
		Estado:         string(h.Estado),
		StoryPoints:    h.StoryPoints,
		HorasEstimadas: horasOpcionales(h.HorasEstimadas),
	}
}
