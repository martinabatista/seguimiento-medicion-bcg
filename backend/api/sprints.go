package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// sprintJSON es un Sprint tal como sale por la API (spec 05), sin sus historias: así se lista y
// así aparece como Sprint en curso en el estado del proyecto.
type sprintJSON struct {
	ID          int64  `json:"id"`
	ProyectoID  int64  `json:"proyecto_id"`
	Numero      int    `json:"numero"`
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
	Goal        string `json:"sprint_goal"`
	Estado      string `json:"estado"`
}

// detalleSprintJSON es un Sprint con sus historias, ordenadas como el backlog (spec 05 CA-3).
type detalleSprintJSON struct {
	sprintJSON
	Historias []historiaJSON `json:"historias"`
}

// sprintEntrada es lo que llega al crear (fechas y Sprint Goal) o al modificar (solo Sprint Goal).
type sprintEntrada struct {
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
	Goal        string `json:"sprint_goal"`
}

type sprintsHandler struct {
	servicio *dominio.ServicioSprints
}

// crear: POST /api/proyectos/{id}/sprints (CA-1). 409 si se superpone con otro Sprint (RN-4).
func (h sprintsHandler) crear(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	var entrada sprintEntrada
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, errCuerpoInvalido)
		return
	}
	inicio, err := dominio.ParsearFecha(entrada.FechaInicio)
	if err != nil {
		responderError(w, err)
		return
	}
	fin, err := dominio.ParsearFecha(entrada.FechaFin)
	if err != nil {
		responderError(w, err)
		return
	}
	s, err := h.servicio.Crear(r.Context(), proyectoID, inicio, fin, entrada.Goal)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, detalleSprintAJSON(dominio.DetalleSprint{Sprint: s}))
}

// listar: GET /api/proyectos/{id}/sprints (CA-3), ordenados por fecha de inicio.
func (h sprintsHandler) listar(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	sprints, err := h.servicio.Listar(r.Context(), proyectoID)
	if err != nil {
		responderError(w, err)
		return
	}
	lista := make([]sprintJSON, 0, len(sprints))
	for _, s := range sprints {
		lista = append(lista, sprintAJSON(s))
	}
	writeJSON(w, http.StatusOK, lista)
}

// ver: GET /api/proyectos/{id}/sprints/{sprintID} (CA-3, CA-7).
func (h sprintsHandler) ver(w http.ResponseWriter, r *http.Request) {
	proyectoID, sprintID, err := leerIDsSprint(r)
	if err != nil {
		responderError(w, err)
		return
	}
	d, err := h.servicio.Obtener(r.Context(), proyectoID, sprintID)
	h.responder(w, d, err)
}

// modificar: PUT /api/proyectos/{id}/sprints/{sprintID} (CA-2). Solo cambia el Sprint Goal.
func (h sprintsHandler) modificar(w http.ResponseWriter, r *http.Request) {
	proyectoID, sprintID, err := leerIDsSprint(r)
	if err != nil {
		responderError(w, err)
		return
	}
	var entrada sprintEntrada
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, errCuerpoInvalido)
		return
	}
	d, err := h.servicio.ModificarGoal(r.Context(), proyectoID, sprintID, entrada.Goal)
	h.responder(w, d, err)
}

// asignar: POST /api/proyectos/{id}/sprints/{sprintID}/historias con {"historia_id": N} (CA-3, CA-4).
func (h sprintsHandler) asignar(w http.ResponseWriter, r *http.Request) {
	proyectoID, sprintID, err := leerIDsSprint(r)
	if err != nil {
		responderError(w, err)
		return
	}
	var entrada struct {
		HistoriaID int64 `json:"historia_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, errCuerpoInvalido)
		return
	}
	d, err := h.servicio.AsignarHistoria(r.Context(), proyectoID, sprintID, entrada.HistoriaID)
	h.responder(w, d, err)
}

// quitar: DELETE /api/proyectos/{id}/sprints/{sprintID}/historias/{historiaID} (CA-3).
func (h sprintsHandler) quitar(w http.ResponseWriter, r *http.Request) {
	proyectoID, sprintID, err := leerIDsSprint(r)
	if err != nil {
		responderError(w, err)
		return
	}
	historiaID, err := strconv.ParseInt(r.PathValue("historiaID"), 10, 64)
	if err != nil {
		responderError(w, dominio.ErrHistoriaNoExiste)
		return
	}
	d, err := h.servicio.QuitarHistoria(r.Context(), proyectoID, sprintID, historiaID)
	h.responder(w, d, err)
}

func (h sprintsHandler) responder(w http.ResponseWriter, d dominio.DetalleSprint, err error) {
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, detalleSprintAJSON(d))
}

// leerIDsSprint lee el {id} del proyecto y el {sprintID} de la ruta. Un id que no es un número no
// puede existir (CA-7).
func leerIDsSprint(r *http.Request) (proyectoID, sprintID int64, err error) {
	proyectoID, err = leerID(r)
	if err != nil {
		return 0, 0, err
	}
	sprintID, err = strconv.ParseInt(r.PathValue("sprintID"), 10, 64)
	if err != nil {
		return 0, 0, dominio.ErrSprintNoExiste
	}
	return proyectoID, sprintID, nil
}

func sprintAJSON(s dominio.Sprint) sprintJSON {
	return sprintJSON{
		ID:          s.ID,
		ProyectoID:  s.ProyectoID,
		Numero:      s.Numero,
		FechaInicio: s.FechaInicio.Format(dominio.FormatoFecha),
		FechaFin:    s.FechaFin.Format(dominio.FormatoFecha),
		Goal:        s.Goal,
		Estado:      string(s.Estado),
	}
}

func detalleSprintAJSON(d dominio.DetalleSprint) detalleSprintJSON {
	historias := make([]historiaJSON, 0, len(d.Historias))
	for _, h := range d.Historias {
		historias = append(historias, historiaAJSON(h))
	}
	return detalleSprintJSON{sprintJSON: sprintAJSON(d.Sprint), Historias: historias}
}
