// Package api expone el dominio por HTTP como una API REST.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// NewRouter arma las rutas de la API.
func NewRouter(proyectos *dominio.ServicioProyectos, backlog *dominio.ServicioBacklog, integrantes *dominio.ServicioIntegrantes, estado *dominio.ServicioEstado, sprints *dominio.ServicioSprints, poker *dominio.ServicioPlanningPoker, metricas *dominio.ServicioMetricas, esfuerzo *dominio.ServicioEsfuerzo, dashboard *dominio.ServicioDashboard, metricasEsfuerzo *dominio.ServicioMetricasEsfuerzo) http.Handler {
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
	mux.HandleFunc("PUT /api/proyectos/{id}/backlog/{historiaID}/story-points", b.estimar)
	mux.HandleFunc("DELETE /api/proyectos/{id}/backlog/{historiaID}/story-points", b.quitarEstimacion)
	mux.HandleFunc("GET /api/story-points/escala", escalaStoryPoints)
	mux.HandleFunc("PUT /api/proyectos/{id}/backlog/{historiaID}/horas-estimadas", b.estimarHoras)
	mux.HandleFunc("DELETE /api/proyectos/{id}/backlog/{historiaID}/horas-estimadas", b.quitarHorasEstimadas)

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
	mux.HandleFunc("POST /api/proyectos/{id}/sprints/{sprintID}/historias/{historiaID}/completar", s.completar)
	mux.HandleFunc("DELETE /api/proyectos/{id}/sprints/{sprintID}/historias/{historiaID}/completar", s.deshacer)
	mux.HandleFunc("POST /api/proyectos/{id}/sprints/{sprintID}/cerrar", s.cerrar)
	mux.HandleFunc("GET /api/proyectos/{id}/sprints/anteriores", s.anteriores)

	pp := pokerHandler{servicio: poker}
	mux.HandleFunc("POST /api/proyectos/{id}/votaciones", pp.abrir)
	mux.HandleFunc("GET /api/proyectos/{id}/votaciones", pp.abiertas)
	mux.HandleFunc("GET /api/proyectos/{id}/votaciones/{votacionID}", pp.ver)
	mux.HandleFunc("POST /api/proyectos/{id}/votaciones/{votacionID}/votos", pp.votar)
	mux.HandleFunc("POST /api/proyectos/{id}/votaciones/{votacionID}/revelar", pp.revelar)
	mux.HandleFunc("POST /api/proyectos/{id}/votaciones/{votacionID}/rondas", pp.nuevaRonda)
	mux.HandleFunc("POST /api/proyectos/{id}/votaciones/{votacionID}/acordar", pp.acordar)

	me := metricasHandler{servicio: metricas}
	mux.HandleFunc("GET /api/proyectos/{id}/metricas", me.consultar)

	mesf := metricasEsfuerzoHandler{servicio: metricasEsfuerzo}
	mux.HandleFunc("GET /api/proyectos/{id}/metricas/esfuerzo", mesf.consultar)

	es := esfuerzoHandler{servicio: esfuerzo}
	mux.HandleFunc("POST /api/proyectos/{id}/backlog/{historiaID}/esfuerzo", es.registrar)
	mux.HandleFunc("GET /api/proyectos/{id}/backlog/{historiaID}/esfuerzo", es.consultar)
	mux.HandleFunc("DELETE /api/proyectos/{id}/backlog/{historiaID}/esfuerzo/{registroID}", es.eliminar)

	da := dashboardHandler{servicio: dashboard}
	mux.HandleFunc("GET /api/proyectos/{id}/dashboard", da.consultar)
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
