package api

import (
	"net/http"
	"strings"
	"testing"
)

// spec 10 CA-1 a CA-4, CA-6: métricas de un proyecto con un Sprint cerrado.
func TestConsultarMetricas(t *testing.T) {
	h := conSprint(t) // proyecto 1, historia 1 y Sprint 1 abierto
	pedir(t, h, http.MethodPut, "/api/proyectos/1/backlog/1/story-points", `{"story_points":5}`)
	pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints/1/historias", `{"historia_id":1}`)

	// Sprint abierto: parcial, sin velocidad todavía.
	rec := pedir(t, h, http.MethodGet, "/api/proyectos/1/metricas", "")
	var m metricasJSON
	leerJSON(t, rec, &m)
	if rec.Code != http.StatusOK || len(m.Sprints) != 1 || !m.Sprints[0].Parcial || m.Velocidad != nil ||
		m.Sprints[0].StoryPointsPlanificados != 5 {
		t.Fatalf("abierto: código = %d, métricas = %+v", rec.Code, m)
	}
	if !strings.Contains(rec.Body.String(), `"velocidad":null`) {
		t.Errorf("cuerpo = %s, se esperaba \"velocidad\":null sin Sprints cerrados", rec.Body.String())
	}

	pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints/1/historias/1/completar", "")
	pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints/1/cerrar", "")
	rec = pedir(t, h, http.MethodGet, "/api/proyectos/1/metricas", "")
	leerJSON(t, rec, &m)
	s := m.Sprints[0]
	if s.Parcial || s.Estado != "Cerrado" || s.StoryPointsCompletados != 5 || s.PorcentajeCompletadas != 100 ||
		m.Velocidad == nil || *m.Velocidad != 5.0 || m.SprintsCerrados != 1 || m.Totales.StoryPointsCompletados != 5 {
		t.Errorf("cerrado: métricas = %+v", m)
	}
}

// spec 10 CA-7: proyecto inexistente → 404.
func TestMetricasProyectoInexistente(t *testing.T) {
	rec := pedir(t, routerDePrueba(), http.MethodGet, "/api/proyectos/999/metricas", "")
	var e errorJSON
	leerJSON(t, rec, &e)
	if rec.Code != http.StatusNotFound || e.Error != "El proyecto no existe" {
		t.Errorf("código = %d, error = %q; se esperaba 404", rec.Code, e.Error)
	}
}
