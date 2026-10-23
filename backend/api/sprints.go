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
	Historias []historiaEnSprintJSON `json:"historias"`
}

// historiaEnSprintJSON es una historia del Sprint: la del backlog más si se completó en ese Sprint
// (spec 06). En un Sprint cerrado, completada_en_sprint no cambia aunque cambie la historia (RN-8).
type historiaEnSprintJSON struct {
	historiaJSON
	CompletadaEnSprint bool `json:"completada_en_sprint"`
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
	historias := make([]historiaEnSprintJSON, 0, len(d.Historias))
	for _, h := range d.Historias {
		historias = append(historias, historiaEnSprintJSON{historiaAJSON(h), d.Sprint.EstaCompletada(h.ID)})
	}
	return detalleSprintJSON{sprintJSON: sprintAJSON(d.Sprint), Historias: historias}
}

// completar: POST /api/proyectos/{id}/sprints/{sprintID}/historias/{historiaID}/completar (spec 06 CA-1).
func (h sprintsHandler) completar(w http.ResponseWriter, r *http.Request) {
	proyectoID, sprintID, historiaID, err := leerIDsHistoriaDelSprint(r)
	if err != nil {
		responderError(w, err)
		return
	}
	d, err := h.servicio.CompletarHistoria(r.Context(), proyectoID, sprintID, historiaID)
	h.responder(w, d, err)
}

// deshacer: DELETE /api/proyectos/{id}/sprints/{sprintID}/historias/{historiaID}/completar (spec 06 CA-1).
func (h sprintsHandler) deshacer(w http.ResponseWriter, r *http.Request) {
	proyectoID, sprintID, historiaID, err := leerIDsHistoriaDelSprint(r)
	if err != nil {
		responderError(w, err)
		return
	}
	d, err := h.servicio.DeshacerHistoria(r.Context(), proyectoID, sprintID, historiaID)
	h.responder(w, d, err)
}

// cerrar: POST /api/proyectos/{id}/sprints/{sprintID}/cerrar (spec 06 CA-2, CA-3).
func (h sprintsHandler) cerrar(w http.ResponseWriter, r *http.Request) {
	proyectoID, sprintID, err := leerIDsSprint(r)
	if err != nil {
		responderError(w, err)
		return
	}
	d, err := h.servicio.Cerrar(r.Context(), proyectoID, sprintID)
	h.responder(w, d, err)
}

// anteriores: GET /api/proyectos/{id}/sprints/anteriores (spec 06 CA-4): los cerrados, con sus
// historias, del más reciente al más antiguo.
func (h sprintsHandler) anteriores(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	detalles, err := h.servicio.Anteriores(r.Context(), proyectoID)
	if err != nil {
		responderError(w, err)
		return
	}
	lista := make([]detalleSprintJSON, 0, len(detalles))
	for _, d := range detalles {
		lista = append(lista, detalleSprintAJSON(d))
	}
	writeJSON(w, http.StatusOK, lista)
}

// leerIDsHistoriaDelSprint lee el proyecto, el Sprint y la historia de la ruta.
func leerIDsHistoriaDelSprint(r *http.Request) (proyectoID, sprintID, historiaID int64, err error) {
	proyectoID, sprintID, err = leerIDsSprint(r)
	if err != nil {
		return 0, 0, 0, err
	}
	historiaID, err = strconv.ParseInt(r.PathValue("historiaID"), 10, 64)
	if err != nil {
		return 0, 0, 0, dominio.ErrHistoriaNoExiste
	}
	return proyectoID, sprintID, historiaID, nil
}
