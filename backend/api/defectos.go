package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// sprintDeDefectoJSON identifica un Sprint del defecto: su identificador y su número.
type sprintDeDefectoJSON struct {
	ID     int64 `json:"id"`
	Numero int   `json:"numero"`
}

// defectoJSON es un defecto tal como sale por la API (spec 11). sprint_resolucion es null
// mientras no está Resuelto.
type defectoJSON struct {
	ID               int64                `json:"id"`
	ProyectoID       int64                `json:"proyecto_id"`
	Descripcion      string               `json:"descripcion"`
	Severidad        string               `json:"severidad"`
	Estado           string               `json:"estado"`
	HistoriaID       int64                `json:"historia_id"`
	HistoriaTitulo   string               `json:"historia_titulo"`
	SprintDeteccion  sprintDeDefectoJSON  `json:"sprint_deteccion"`
	SprintResolucion *sprintDeDefectoJSON `json:"sprint_resolucion"`
}

// defectoEntrada es lo que llega al registrar. Los identificadores que faltan llegan en 0.
type defectoEntrada struct {
	Descripcion       string `json:"descripcion"`
	Severidad         string `json:"severidad"`
	HistoriaID        int64  `json:"historia_id"`
	SprintDeteccionID int64  `json:"sprint_deteccion_id"`
}

// estadoDefectoEntrada es lo que llega al cambiar el estado; el Sprint solo cuenta al resolver.
type estadoDefectoEntrada struct {
	Estado             string `json:"estado"`
	SprintResolucionID int64  `json:"sprint_resolucion_id"`
}

// metricasDefectosJSON son los detectados y resueltos del proyecto y de cada Sprint (CA-4).
type metricasDefectosJSON struct {
	Detectados int                    `json:"detectados"`
	Resueltos  int                    `json:"resueltos"`
	Sprints    []defectosDeSprintJSON `json:"sprints"`
}

type defectosDeSprintJSON struct {
	SprintID   int64 `json:"sprint_id"`
	Numero     int   `json:"numero"`
	Detectados int   `json:"detectados"`
	Resueltos  int   `json:"resueltos"`
}

type defectosHandler struct {
	servicio *dominio.ServicioDefectos
}

// registrar: POST /api/proyectos/{id}/defectos (CA-1, CA-5).
func (h defectosHandler) registrar(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	var entrada defectoEntrada
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, errCuerpoInvalido)
		return
	}
	d, err := h.servicio.Registrar(r.Context(), proyectoID, dominio.DatosDefecto{
		Descripcion:       entrada.Descripcion,
		Severidad:         entrada.Severidad,
		HistoriaID:        entrada.HistoriaID,
		SprintDeteccionID: entrada.SprintDeteccionID,
	})
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, defectoAJSON(d))
}

// listar: GET /api/proyectos/{id}/defectos (CA-5). En el orden en que se registraron.
func (h defectosHandler) listar(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	defectos, err := h.servicio.Listar(r.Context(), proyectoID)
	if err != nil {
		responderError(w, err)
		return
	}
	respuesta := make([]defectoJSON, 0, len(defectos))
	for _, d := range defectos {
		respuesta = append(respuesta, defectoAJSON(d))
	}
	writeJSON(w, http.StatusOK, respuesta)
}

// ver: GET /api/proyectos/{id}/defectos/{defectoID} (CA-5).
func (h defectosHandler) ver(w http.ResponseWriter, r *http.Request) {
	proyectoID, defectoID, err := leerIDsDefecto(r)
	if err != nil {
		responderError(w, err)
		return
	}
	d, err := h.servicio.Obtener(r.Context(), proyectoID, defectoID)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, defectoAJSON(d))
}

// cambiarEstado: PUT /api/proyectos/{id}/defectos/{defectoID}/estado (CA-2, CA-3).
func (h defectosHandler) cambiarEstado(w http.ResponseWriter, r *http.Request) {
	proyectoID, defectoID, err := leerIDsDefecto(r)
	if err != nil {
		responderError(w, err)
		return
	}
	var entrada estadoDefectoEntrada
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, errCuerpoInvalido)
		return
	}
	d, err := h.servicio.CambiarEstado(r.Context(), proyectoID, defectoID, entrada.Estado, entrada.SprintResolucionID)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, defectoAJSON(d))
}

// metricas: GET /api/proyectos/{id}/metricas/defectos (CA-4).
func (h defectosHandler) metricas(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	m, err := h.servicio.Metricas(r.Context(), proyectoID)
	if err != nil {
		responderError(w, err)
		return
	}
	respuesta := metricasDefectosJSON{
		Detectados: m.Detectados,
		Resueltos:  m.Resueltos,
		Sprints:    make([]defectosDeSprintJSON, 0, len(m.Sprints)),
	}
	for _, s := range m.Sprints {
		respuesta.Sprints = append(respuesta.Sprints, defectosDeSprintJSON{SprintID: s.SprintID, Numero: s.Numero, Detectados: s.Detectados, Resueltos: s.Resueltos})
	}
	writeJSON(w, http.StatusOK, respuesta)
}

// leerIDsDefecto lee el proyecto y el defecto de la ruta. Un id que no es un número no puede existir (CA-7).
func leerIDsDefecto(r *http.Request) (proyectoID, defectoID int64, err error) {
	proyectoID, err = leerID(r)
	if err != nil {
		return 0, 0, err
	}
	defectoID, err = strconv.ParseInt(r.PathValue("defectoID"), 10, 64)
	if err != nil {
		return 0, 0, dominio.ErrDefectoNoExiste
	}
	return proyectoID, defectoID, nil
}

func defectoAJSON(d dominio.Defecto) defectoJSON {
	out := defectoJSON{
		ID:              d.ID,
		ProyectoID:      d.ProyectoID,
		Descripcion:     d.Descripcion,
		Severidad:       string(d.Severidad),
		Estado:          string(d.Estado),
		HistoriaID:      d.HistoriaID,
		HistoriaTitulo:  d.HistoriaTitulo,
		SprintDeteccion: sprintDeDefectoJSON{ID: d.SprintDeteccionID, Numero: d.SprintDeteccionNumero},
	}
	if d.SprintResolucionID != 0 {
		out.SprintResolucion = &sprintDeDefectoJSON{ID: d.SprintResolucionID, Numero: d.SprintResolucionNumero}
	}
	return out
}
