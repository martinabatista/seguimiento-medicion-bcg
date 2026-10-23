package api

import (
	"net/http"
	"strings"
	"testing"
)

// El proyecto de conProyecto va del 2026-11-01 al 2026-12-15 y "hoy" es el 2026-11-20.
const sprintJSONDePrueba = `{"fecha_inicio":"2026-11-15","fecha_fin":"2026-11-28","sprint_goal":"Registrar clientes"}`

// conSprint devuelve un router con el proyecto 1, la historia 1 en su backlog y el Sprint 1 creado.
func conSprint(t *testing.T) http.Handler {
	t.Helper()
	h := conProyecto(t)
	if rec := pedir(t, h, http.MethodPost, "/api/proyectos/1/backlog", historiaJSONDePrueba); rec.Code != http.StatusCreated {
		t.Fatalf("crear historia de partida: %d %s", rec.Code, rec.Body.String())
	}
	if rec := pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints", sprintJSONDePrueba); rec.Code != http.StatusCreated {
		t.Fatalf("crear Sprint de partida: %d %s", rec.Code, rec.Body.String())
	}
	return h
}

// spec 05 CA-1: 201 con el Sprint 1 abierto y sin historias.
func TestCrearSprint(t *testing.T) {
	rec := pedir(t, conProyecto(t), http.MethodPost, "/api/proyectos/1/sprints", sprintJSONDePrueba)

	if rec.Code != http.StatusCreated {
		t.Fatalf("código = %d, se esperaba 201 (%s)", rec.Code, rec.Body.String())
	}
	var s detalleSprintJSON
	leerJSON(t, rec, &s)
	if s.ID != 1 || s.ProyectoID != 1 || s.Numero != 1 || s.FechaInicio != "2026-11-15" || s.FechaFin != "2026-11-28" ||
		s.Goal != "Registrar clientes" || s.Estado != "Abierto" || s.Historias == nil || len(s.Historias) != 0 {
		t.Errorf("sprint = %+v, se esperaba el Sprint 1 abierto con historias []", s)
	}
}

// spec 05 CA-1, CA-2: datos inválidos → 400; superposición → 409.
func TestCrearSprintInvalido(t *testing.T) {
	casos := []struct {
		nombre, cuerpo, mensaje string
		codigo                  int
	}{
		{"fuera del proyecto", `{"fecha_inicio":"2026-12-10","fecha_fin":"2026-12-20","sprint_goal":"x"}`, "El Sprint debe estar dentro de las fechas del proyecto", 400},
		{"fin igual al inicio", `{"fecha_inicio":"2026-12-01","fecha_fin":"2026-12-01","sprint_goal":"x"}`, "La fecha de finalización debe ser posterior a la de inicio", 400},
		{"sin fecha de fin", `{"fecha_inicio":"2026-12-01","sprint_goal":"x"}`, "Las fechas de inicio y finalización son obligatorias", 400},
		{"fecha inexistente", `{"fecha_inicio":"2026-11-31","fecha_fin":"2026-12-05","sprint_goal":"x"}`, "La fecha no es válida", 400},
		{"sin Sprint Goal", `{"fecha_inicio":"2026-12-01","fecha_fin":"2026-12-05"}`, "El Sprint Goal es obligatorio", 400},
		{"se superpone", `{"fecha_inicio":"2026-11-28","fecha_fin":"2026-12-05","sprint_goal":"x"}`, "El Sprint se superpone con otro Sprint del proyecto", 409},
		{"JSON mal formado", `{"fecha_inicio":`, "El cuerpo de la solicitud no es un JSON válido", 400},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedir(t, conSprint(t), http.MethodPost, "/api/proyectos/1/sprints", c.cuerpo)
			var e errorJSON
			leerJSON(t, rec, &e)
			if rec.Code != c.codigo || e.Error != c.mensaje {
				t.Errorf("código = %d, error = %q; se esperaba %d y %q", rec.Code, e.Error, c.codigo, c.mensaje)
			}
		})
	}
}

// spec 05 CA-2: modificar el Sprint Goal.
func TestModificarSprintGoal(t *testing.T) {
	h := conSprint(t)
	rec := pedir(t, h, http.MethodPut, "/api/proyectos/1/sprints/1", `{"sprint_goal":"Registrar y loguear clientes"}`)
	var s detalleSprintJSON
	leerJSON(t, rec, &s)
	if rec.Code != http.StatusOK || s.Goal != "Registrar y loguear clientes" {
		t.Errorf("código = %d, sprint = %+v", rec.Code, s)
	}
}

// spec 05 CA-3, CA-4: asignar, ver, listar y quitar historias.
func TestAsignarYQuitarHistoria(t *testing.T) {
	h := conSprint(t)

	rec := pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints/1/historias", `{"historia_id":1}`)
	var s detalleSprintJSON
	leerJSON(t, rec, &s)
	if rec.Code != http.StatusOK || len(s.Historias) != 1 || s.Historias[0].Titulo != "Registrar usuario" {
		t.Fatalf("asignar: código = %d, sprint = %+v", rec.Code, s)
	}

	rec = pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints/1/historias", `{"historia_id":1}`)
	if rec.Code != http.StatusConflict || !strings.Contains(rec.Body.String(), "La historia ya está asignada a un Sprint abierto") {
		t.Errorf("asignar dos veces: código = %d, %s; se esperaba 409", rec.Code, rec.Body.String())
	}

	rec = pedir(t, h, http.MethodGet, "/api/proyectos/1/sprints", "")
	var lista []detalleSprintJSON
	leerJSON(t, rec, &lista)
	if rec.Code != http.StatusOK || len(lista) != 1 || lista[0].Numero != 1 {
		t.Errorf("listar: código = %d, lista = %+v", rec.Code, lista)
	}

	rec = pedir(t, h, http.MethodDelete, "/api/proyectos/1/sprints/1/historias/1", "")
	leerJSON(t, rec, &s)
	if rec.Code != http.StatusOK || len(s.Historias) != 0 {
		t.Errorf("quitar: código = %d, sprint = %+v", rec.Code, s)
	}

	rec = pedir(t, h, http.MethodDelete, "/api/proyectos/1/sprints/1/historias/1", "")
	if rec.Code != http.StatusConflict || !strings.Contains(rec.Body.String(), "La historia no está asignada a este Sprint") {
		t.Errorf("quitar dos veces: código = %d, %s; se esperaba 409", rec.Code, rec.Body.String())
	}
}

// spec 05 CA-6: eliminar del backlog una historia de un Sprint abierto → 409.
func TestEliminarHistoriaDeSprintAbierto(t *testing.T) {
	h := conSprint(t)
	pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints/1/historias", `{"historia_id":1}`)

	rec := pedir(t, h, http.MethodDelete, "/api/proyectos/1/backlog/1", "")
	if rec.Code != http.StatusConflict || !strings.Contains(rec.Body.String(), "La historia está asignada a un Sprint abierto") {
		t.Errorf("código = %d, %s; se esperaba 409", rec.Code, rec.Body.String())
	}
}

// spec 05 CA-5: el estado del proyecto incluye el Sprint en curso (hoy = 2026-11-20).
func TestEstadoConSprintEnCurso(t *testing.T) {
	sin := pedir(t, conProyecto(t), http.MethodGet, "/api/proyectos/1/estado", "")
	if !strings.Contains(sin.Body.String(), `"sprint_en_curso":null`) {
		t.Errorf("sin Sprints: cuerpo = %s, se esperaba \"sprint_en_curso\":null", sin.Body.String())
	}

	rec := pedir(t, conSprint(t), http.MethodGet, "/api/proyectos/1/estado", "")
	var e estadoJSON
	leerJSON(t, rec, &e)
	if rec.Code != http.StatusOK || e.SprintEnCurso == nil || e.SprintEnCurso.Numero != 1 || e.SprintEnCurso.Goal != "Registrar clientes" {
		t.Errorf("código = %d, sprint en curso = %+v; se esperaba el Sprint 1", rec.Code, e.SprintEnCurso)
	}
}

// spec 05 CA-7: proyecto, Sprint o historia inexistentes → 404.
func TestSprintInexistente(t *testing.T) {
	casos := []struct{ nombre, metodo, ruta, cuerpo, mensaje string }{
		{"crear en proyecto inexistente", http.MethodPost, "/api/proyectos/999/sprints", sprintJSONDePrueba, "El proyecto no existe"},
		{"ver Sprint inexistente", http.MethodGet, "/api/proyectos/1/sprints/999", "", "El Sprint no existe"},
		{"id de Sprint que no es un número", http.MethodGet, "/api/proyectos/1/sprints/abc", "", "El Sprint no existe"},
		{"asignar historia inexistente", http.MethodPost, "/api/proyectos/1/sprints/1/historias", `{"historia_id":999}`, "La historia no existe"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedir(t, conSprint(t), c.metodo, c.ruta, c.cuerpo)
			var e errorJSON
			leerJSON(t, rec, &e)
			if rec.Code != http.StatusNotFound || e.Error != c.mensaje {
				t.Errorf("código = %d, error = %q; se esperaba 404 y %q", rec.Code, e.Error, c.mensaje)
			}
		})
	}
}
