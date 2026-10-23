package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// proyectoJSON es un proyecto tal como viaja por la API (spec 01).
type proyectoJSON struct {
	ID          int64  `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
}

type errorJSON struct {
	Error string `json:"error"`
}

var errCuerpoInvalido = dominio.ErrDatoInvalido{Mensaje: "El cuerpo de la solicitud no es un JSON válido"}

type proyectosHandler struct {
	servicio *dominio.ServicioProyectos
}

// crear: POST /api/proyectos (CA-1).
func (h proyectosHandler) crear(w http.ResponseWriter, r *http.Request) {
	datos, inicio, fin, err := leerProyecto(r)
	if err != nil {
		responderError(w, err)
		return
	}
	p, err := h.servicio.Crear(r.Context(), datos.Nombre, datos.Descripcion, inicio, fin)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, aJSON(p))
}

// listar: GET /api/proyectos (CA-6).
func (h proyectosHandler) listar(w http.ResponseWriter, r *http.Request) {
	proyectos, err := h.servicio.Listar(r.Context())
	if err != nil {
		responderError(w, err)
		return
	}
	lista := make([]proyectoJSON, 0, len(proyectos))
	for _, p := range proyectos {
		lista = append(lista, aJSON(p))
	}
	writeJSON(w, http.StatusOK, lista)
}

// ver: GET /api/proyectos/{id} (CA-6, CA-7).
func (h proyectosHandler) ver(w http.ResponseWriter, r *http.Request) {
	id, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	p, err := h.servicio.Obtener(r.Context(), id)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, aJSON(p))
}

// modificar: PUT /api/proyectos/{id} (CA-2, CA-7). Reemplaza los cuatro datos.
func (h proyectosHandler) modificar(w http.ResponseWriter, r *http.Request) {
	id, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	datos, inicio, fin, err := leerProyecto(r)
	if err != nil {
		responderError(w, err)
		return
	}
	p, err := h.servicio.Modificar(r.Context(), id, datos.Nombre, datos.Descripcion, inicio, fin)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, aJSON(p))
}

// leerProyecto decodifica el cuerpo y convierte las fechas AAAA-MM-DD.
func leerProyecto(r *http.Request) (proyectoJSON, time.Time, time.Time, error) {
	var datos proyectoJSON
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		return datos, time.Time{}, time.Time{}, errCuerpoInvalido
	}
	inicio, err := dominio.ParsearFecha(datos.FechaInicio)
	if err != nil {
		return datos, time.Time{}, time.Time{}, err
	}
	fin, err := dominio.ParsearFecha(datos.FechaFin)
	if err != nil {
		return datos, time.Time{}, time.Time{}, err
	}
	return datos, inicio, fin, nil
}

// leerID lee el {id} de la ruta. Un id que no es un número no puede existir (CA-7).
func leerID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return 0, dominio.ErrProyectoNoExiste
	}
	return id, nil
}

func aJSON(p dominio.Proyecto) proyectoJSON {
	return proyectoJSON{
		ID:          p.ID,
		Nombre:      p.Nombre,
		Descripcion: p.Descripcion,
		FechaInicio: p.FechaInicio.Format(dominio.FormatoFecha),
		FechaFin:    p.FechaFin.Format(dominio.FormatoFecha),
	}
}

// responderError traduce los errores del dominio al formato común {"error": "..."}:
// dato inválido → 400, no existe → 404, conflicto con lo ya registrado → 409, cualquier otro → 500
// (se registra en el log).
func responderError(w http.ResponseWriter, err error) {
	var invalido dominio.ErrDatoInvalido
	var conflicto dominio.ErrConflicto
	switch {
	case errors.As(err, &invalido):
		writeJSON(w, http.StatusBadRequest, errorJSON{invalido.Mensaje})
	case errors.As(err, &conflicto):
		writeJSON(w, http.StatusConflict, errorJSON{conflicto.Mensaje})
	case errors.Is(err, dominio.ErrProyectoNoExiste), errors.Is(err, dominio.ErrHistoriaNoExiste),
		errors.Is(err, dominio.ErrIntegranteNoExiste), errors.Is(err, dominio.ErrSprintNoExiste):
		writeJSON(w, http.StatusNotFound, errorJSON{err.Error()})
	case errors.Is(err, dominio.ErrIntegranteYaRegistrado):
		writeJSON(w, http.StatusConflict, errorJSON{err.Error()})
	default:
		log.Printf("error interno: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorJSON{"Error interno del servidor"})
	}
}
