package api

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// En los tests de la API hoy es hoyDePrueba (2026-11-20) y el proyecto 1 va del 2026-11-01 al
// 2026-12-15 (tiendaJSON).

// spec 13 CA-1, CA-8: un proyecto sin historias ni Sprints igual tiene Dashboard, con ceros, listas
// vacías ([] y no null) y null donde no hay dato.
func TestDashboardProyectoVacio(t *testing.T) {
	rec := pedir(t, conProyecto(t), http.MethodGet, "/api/proyectos/1/dashboard", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200 (%s)", rec.Code, rec.Body.String())
	}
	cuerpo := rec.Body.String()
	for _, parte := range []string{
		`"estado":"En curso"`, `"dias_para_empezar":null`, `"dias_para_terminar":25`,
		`"sprint_en_curso":null`, `"velocidad":null`,
		`"story_points_por_sprint":[]`, `"velocidad":{"sprints":[],"promedio":null}`,
		`"backlog_por_estado":[{"estado":"Pendiente","cantidad":0},{"estado":"En progreso","cantidad":0},{"estado":"Completada","cantidad":0}]`,
	} {
		if !strings.Contains(cuerpo, parte) {
			t.Errorf("cuerpo = %s\nno contiene %s", cuerpo, parte)
		}
	}
}

// spec 13 CA-2 a CA-7: indicadores y gráficos con un Sprint cerrado y otro en curso.
func TestDashboardConSprints(t *testing.T) {
	h := conProyecto(t)
	crear := func(metodo, ruta, cuerpo string) {
		t.Helper()
		if rec := pedir(t, h, metodo, ruta, cuerpo); rec.Code >= 300 {
			t.Fatalf("%s %s: %d %s", metodo, ruta, rec.Code, rec.Body.String())
		}
	}
	// Historias 1 (5 SP), 2 (8 SP) y 3 (sin estimar).
	for k, titulo := range []string{"Registrar usuario", "Ver catálogo", "Pagar pedido"} {
		crear(http.MethodPost, "/api/proyectos/1/backlog", fmt.Sprintf(`{"titulo":%q,"descripcion":"d","prioridad":%d}`, titulo, k+1))
	}
	crear(http.MethodPut, "/api/proyectos/1/backlog/1/story-points", `{"story_points":5}`)
	crear(http.MethodPut, "/api/proyectos/1/backlog/2/story-points", `{"story_points":8}`)
	// Sprint 1 (cerrado): historias 1 y 2, completada la 1. Sprint 2 (en curso hoy): historia 2.
	crear(http.MethodPost, "/api/proyectos/1/sprints", `{"fecha_inicio":"2026-11-01","fecha_fin":"2026-11-14","sprint_goal":"Registro"}`)
	crear(http.MethodPost, "/api/proyectos/1/sprints/1/historias", `{"historia_id":1}`)
	crear(http.MethodPost, "/api/proyectos/1/sprints/1/historias", `{"historia_id":2}`)
	crear(http.MethodPost, "/api/proyectos/1/sprints/1/historias/1/completar", "")
	crear(http.MethodPost, "/api/proyectos/1/sprints/1/cerrar", "")
	crear(http.MethodPost, "/api/proyectos/1/sprints", `{"fecha_inicio":"2026-11-15","fecha_fin":"2026-11-28","sprint_goal":"Catálogo"}`)
	crear(http.MethodPost, "/api/proyectos/1/sprints/2/historias", `{"historia_id":2}`)

	rec := pedir(t, h, http.MethodGet, "/api/proyectos/1/dashboard", "")
	var d dashboardJSON
	leerJSON(t, rec, &d)

	if d.SprintEnCurso == nil || d.SprintEnCurso.Numero != 2 || d.SprintEnCurso.SprintGoal != "Catálogo" || d.SprintEnCurso.FechaFin != "2026-11-28" ||
		d.SprintEnCurso.DiasRestantes != 8 || d.SprintEnCurso.StoryPointsPlanificados != 8 || d.SprintEnCurso.StoryPointsCompletados != 0 {
		t.Errorf("Sprint en curso = %+v, se esperaba el 2 (Catálogo) con 8 días y 8/0 Story Points (CA-2)", d.SprintEnCurso)
	}
	if b := d.Backlog; b.Total != 3 || b.Completadas != 1 || b.PorcentajeCompletadas != 33 || b.StoryPointsTotales != 13 || b.StoryPointsCompletados != 5 || b.SinEstimar != 1 {
		t.Errorf("backlog = %+v, se esperaba 1 de 3 (33 %%), 5 de 13 SP y 1 sin estimar (CA-3)", b)
	}
	if d.Velocidad == nil || *d.Velocidad != 5 {
		t.Errorf("velocidad = %v, se esperaba 5 (CA-4)", d.Velocidad)
	}
	g := d.Graficos
	if len(g.StoryPointsPorSprint) != 2 || g.StoryPointsPorSprint[0] != (barraSprintJSON{Sprint: "Sprint 1", Numero: 1, Planificados: 13, Completados: 5}) ||
		g.StoryPointsPorSprint[1] != (barraSprintJSON{Sprint: "Sprint 2", Numero: 2, Planificados: 8, Completados: 0, Parcial: true}) {
		t.Errorf("Story Points por Sprint = %+v (CA-5)", g.StoryPointsPorSprint)
	}
	if len(g.Velocidad.Sprints) != 1 || g.Velocidad.Sprints[0] != (barraVelocidadJSON{Sprint: "Sprint 1", Numero: 1, Completados: 5}) ||
		g.Velocidad.Promedio == nil || *g.Velocidad.Promedio != 5 {
		t.Errorf("gráfico de velocidad = %+v, se esperaba solo el Sprint 1 con 5 y promedio 5 (CA-6)", g.Velocidad)
	}
	// Asignar una historia a un Sprint no le cambia el estado (spec 05): "Ver catálogo" sigue Pendiente.
	if len(g.BacklogPorEstado) != 3 || g.BacklogPorEstado[0] != (porcionBacklogJSON{"Pendiente", 2}) ||
		g.BacklogPorEstado[1] != (porcionBacklogJSON{"En progreso", 0}) || g.BacklogPorEstado[2] != (porcionBacklogJSON{"Completada", 1}) {
		t.Errorf("backlog por estado = %+v (CA-7)", g.BacklogPorEstado)
	}
}

// spec 13 CA-8: el Dashboard de un proyecto que no existe → 404.
func TestDashboardProyectoInexistente(t *testing.T) {
	for _, ruta := range []string{"/api/proyectos/999/dashboard", "/api/proyectos/abc/dashboard"} {
		rec := pedir(t, routerDePrueba(), http.MethodGet, ruta, "")
		var e errorJSON
		leerJSON(t, rec, &e)
		if rec.Code != http.StatusNotFound || e.Error != "El proyecto no existe" {
			t.Errorf("%s: código = %d, error = %q; se esperaba 404", ruta, rec.Code, e.Error)
		}
	}
}
