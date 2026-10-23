package api

import (
	"net/http"
	"strings"
	"testing"
)

// conHistoriaEnSprint devuelve un router con el Sprint 1 y la historia 1 asignada.
func conHistoriaEnSprint(t *testing.T) http.Handler {
	t.Helper()
	h := conSprint(t)
	if rec := pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints/1/historias", `{"historia_id":1}`); rec.Code != http.StatusOK {
		t.Fatalf("asignar historia de partida: %d %s", rec.Code, rec.Body.String())
	}
	return h
}

// spec 06 CA-1: completar y deshacer.
func TestCompletarYDeshacerHistoria(t *testing.T) {
	h := conHistoriaEnSprint(t)

	rec := pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints/1/historias/1/completar", "")
	var s detalleSprintJSON
	leerJSON(t, rec, &s)
	if rec.Code != http.StatusOK || len(s.Historias) != 1 || !s.Historias[0].CompletadaEnSprint || s.Historias[0].Estado != "Completada" {
		t.Fatalf("completar: código = %d, sprint = %+v", rec.Code, s)
	}

	rec = pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints/1/historias/1/completar", "")
	if rec.Code != http.StatusConflict || !strings.Contains(rec.Body.String(), "La historia ya está completada") {
		t.Errorf("completar dos veces: código = %d, %s; se esperaba 409", rec.Code, rec.Body.String())
	}

	rec = pedir(t, h, http.MethodDelete, "/api/proyectos/1/sprints/1/historias/1/completar", "")
	leerJSON(t, rec, &s)
	if rec.Code != http.StatusOK || s.Historias[0].CompletadaEnSprint || s.Historias[0].Estado != "En progreso" {
		t.Errorf("deshacer: código = %d, sprint = %+v", rec.Code, s)
	}
}

// spec 06 CA-2, CA-3: cerrar y, después, no admite cambios.
func TestCerrarSprint(t *testing.T) {
	h := conHistoriaEnSprint(t)

	rec := pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints/1/cerrar", "")
	var s detalleSprintJSON
	leerJSON(t, rec, &s)
	if rec.Code != http.StatusOK || s.Estado != "Cerrado" || len(s.Historias) != 1 {
		t.Fatalf("cerrar: código = %d, sprint = %+v", rec.Code, s)
	}

	for _, ruta := range []string{"/api/proyectos/1/sprints/1/cerrar", "/api/proyectos/1/sprints/1/historias/1/completar"} {
		rec = pedir(t, h, http.MethodPost, ruta, "")
		if rec.Code != http.StatusConflict || !strings.Contains(rec.Body.String(), "El Sprint está cerrado") {
			t.Errorf("POST %s: código = %d, %s; se esperaba 409", ruta, rec.Code, rec.Body.String())
		}
	}
}

// spec 06 CA-4: Sprints anteriores (cerrados), del más reciente al más antiguo.
func TestSprintsAnteriores(t *testing.T) {
	h := conSprint(t) // Sprint 1: del 2026-11-15 al 2026-11-28
	pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints", `{"fecha_inicio":"2026-11-01","fecha_fin":"2026-11-14","sprint_goal":"Primero"}`)
	if rec := pedir(t, h, http.MethodGet, "/api/proyectos/1/sprints/anteriores", ""); rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Errorf("sin cerrados: código = %d, cuerpo = %s; se esperaba []", rec.Code, rec.Body.String())
	}
	pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints/1/cerrar", "")
	pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints/2/cerrar", "")

	rec := pedir(t, h, http.MethodGet, "/api/proyectos/1/sprints/anteriores", "")
	var lista []detalleSprintJSON
	leerJSON(t, rec, &lista)
	if rec.Code != http.StatusOK || len(lista) != 2 || lista[0].FechaInicio != "2026-11-15" || lista[1].FechaInicio != "2026-11-01" {
		t.Errorf("anteriores: código = %d, lista = %+v; se esperaba el del 15/11 y después el del 01/11", rec.Code, lista)
	}
}

// spec 06 CA-5: el backlog no pasa a Completada una historia de un Sprint abierto (409).
func TestBacklogNoCompletaHistoriaDeSprint(t *testing.T) {
	h := conHistoriaEnSprint(t)
	rec := pedir(t, h, http.MethodPut, "/api/proyectos/1/backlog/1",
		`{"titulo":"Registrar usuario","descripcion":"Alta de clientes","prioridad":1,"estado":"Completada"}`)
	if rec.Code != http.StatusConflict || !strings.Contains(rec.Body.String(), "se completa desde el Sprint") {
		t.Errorf("código = %d, %s; se esperaba 409", rec.Code, rec.Body.String())
	}
}

// spec 06 CA-6: Sprint o historia inexistentes → 404.
func TestCerrarSprintInexistente(t *testing.T) {
	h := conHistoriaEnSprint(t)
	casos := []struct{ ruta, mensaje string }{
		{"/api/proyectos/1/sprints/999/cerrar", "El Sprint no existe"},
		{"/api/proyectos/1/sprints/1/historias/999/completar", "La historia no existe"},
		{"/api/proyectos/1/sprints/1/historias/abc/completar", "La historia no existe"},
	}
	for _, c := range casos {
		rec := pedir(t, h, http.MethodPost, c.ruta, "")
		var e errorJSON
		leerJSON(t, rec, &e)
		if rec.Code != http.StatusNotFound || e.Error != c.mensaje {
			t.Errorf("POST %s: código = %d, error = %q; se esperaba 404 y %q", c.ruta, rec.Code, e.Error, c.mensaje)
		}
	}
}
