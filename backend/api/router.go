// Package api expone el dominio por HTTP como una API REST.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// NewRouter arma las rutas de la API.
func NewRouter(proyectos *dominio.ServicioProyectos, backlog *dominio.ServicioBacklog, integrantes *dominio.ServicioIntegrantes, estado *dominio.ServicioEstado, sprints *dominio.ServicioSprints) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", health)

	p := proyectosHandler{servicio: proyectos}
	mux.HandleFunc("POST /api/proyectos", p.crear)
	mux.HandleFunc("GET /api/proyectos", p.listar)
	mux.HandleFunc("GET /api/proyectos/{id}", p.ver)
	mux.HandleFunc("PUT /api/proyectos/{id}", p.modificar)

	b := backlogHandler{servicio: backlog}
	mux.HandleFunc("POST /api/proyectos/{id}/backlog", b.crear)
	mux.HandleFunc("GET /api/proyectos/{id}/backlog", b.listar)
	mux.HandleFunc("GET /api/proyectos/{id}/backlog/{historiaID}", b.ver)
	mux.HandleFunc("PUT /api/proyectos/{id}/backlog/{historiaID}", b.modificar)
	mux.HandleFunc("DELETE /api/proyectos/{id}/backlog/{historiaID}", b.eliminar)

	i := integrantesHandler{servicio: integrantes}
	mux.HandleFunc("POST /api/proyectos/{id}/integrantes", i.agregar)
	mux.HandleFunc("GET /api/proyectos/{id}/integrantes", i.listar)
	mux.HandleFunc("GET /api/proyectos/{id}/integrantes/{integranteID}", i.ver)

	e := estadoHandler{servicio: estado}
	mux.HandleFunc("GET /api/proyectos/{id}/estado", e.consultar)

	s := sprintsHandler{servicio: sprints}
	mux.HandleFunc("POST /api/proyectos/{id}/sprints", s.crear)
	mux.HandleFunc("GET /api/proyectos/{id}/sprints", s.listar)
	mux.HandleFunc("GET /api/proyectos/{id}/sprints/{sprintID}", s.ver)
	mux.HandleFunc("PUT /api/proyectos/{id}/sprints/{sprintID}", s.modificar)
	mux.HandleFunc("POST /api/proyectos/{id}/sprints/{sprintID}/historias", s.asignar)
	mux.HandleFunc("DELETE /api/proyectos/{id}/sprints/{sprintID}/historias/{historiaID}", s.quitar)
	return mux
}

func health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
