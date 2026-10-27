package api

import (
	"net/http"
	"strings"
	"testing"
)

const rutaHorasEstimadas = "/api/proyectos/1/backlog/1/horas-estimadas"

// spec 12 CA-1: cargar, cambiar y quitar las horas estimadas; salen como número en la historia.
func TestEstimarHoras(t *testing.T) {
	h := conHistoriaEIntegrante(t)

	rec := pedir(t, h, http.MethodPut, rutaHorasEstimadas, `{"horas_estimadas":8}`)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"horas_estimadas":8`) {
		t.Fatalf("cargar: código = %d, cuerpo = %s; se esperaba 200 con \"horas_estimadas\":8", rec.Code, rec.Body.String())
	}
	rec = pedir(t, h, http.MethodPut, rutaHorasEstimadas, `{"horas_estimadas":6.5}`)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"horas_estimadas":6.5`) {
		t.Fatalf("cambiar: código = %d, cuerpo = %s; se esperaba 200 con \"horas_estimadas\":6.5", rec.Code, rec.Body.String())
	}
	// El backlog también las muestra.
	if rec := pedir(t, h, http.MethodGet, "/api/proyectos/1/backlog/1", ""); !strings.Contains(rec.Body.String(), `"horas_estimadas":6.5`) {
		t.Errorf("ver la historia: cuerpo = %s, se esperaba \"horas_estimadas\":6.5", rec.Body.String())
	}

	rec = pedir(t, h, http.MethodDelete, rutaHorasEstimadas, "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"horas_estimadas":null`) {
		t.Errorf("quitar: código = %d, cuerpo = %s; se esperaba 200 con \"horas_estimadas\":null", rec.Code, rec.Body.String())
	}
}

// spec 12 CA-1, CA-2, CA-7: errores con su código y su mensaje.
func TestEstimarHorasErrores(t *testing.T) {
	h := conHistoriaEIntegrante(t)
	pedir(t, h, http.MethodPost, "/api/proyectos/1/backlog", `{"titulo":"Ver catálogo","descripcion":"d","prioridad":2,"criterios_aceptacion":""}`)
	pedir(t, h, http.MethodPut, "/api/proyectos/1/backlog/2", `{"titulo":"Ver catálogo","descripcion":"d","prioridad":2,"criterios_aceptacion":"","estado":"Completada"}`)

	casos := []struct {
		nombre, metodo, ruta, cuerpo string
		codigo                       int
		mensaje                      string
	}{
		{"faltan las horas", http.MethodPut, rutaHorasEstimadas, `{}`, http.StatusBadRequest, "Las horas estimadas son obligatorias"},
		{"cero", http.MethodPut, rutaHorasEstimadas, `{"horas_estimadas":0}`, http.StatusBadRequest, "Las horas estimadas deben ser un número mayor a 0, con hasta 2 decimales"},
		{"tres decimales", http.MethodPut, rutaHorasEstimadas, `{"horas_estimadas":1.125}`, http.StatusBadRequest, "Las horas estimadas deben ser un número mayor a 0, con hasta 2 decimales"},
		{"cuerpo inválido", http.MethodPut, rutaHorasEstimadas, `{`, http.StatusBadRequest, "El cuerpo de la solicitud no es un JSON válido"},
		{"historia completada", http.MethodPut, "/api/proyectos/1/backlog/2/horas-estimadas", `{"horas_estimadas":8}`, http.StatusConflict, "La historia ya está completada"},
		{"quitar en una completada", http.MethodDelete, "/api/proyectos/1/backlog/2/horas-estimadas", "", http.StatusConflict, "La historia ya está completada"},
		{"historia inexistente", http.MethodPut, "/api/proyectos/1/backlog/999/horas-estimadas", `{"horas_estimadas":8}`, http.StatusNotFound, "La historia no existe"},
		{"proyecto inexistente", http.MethodPut, "/api/proyectos/999/backlog/1/horas-estimadas", `{"horas_estimadas":8}`, http.StatusNotFound, "El proyecto no existe"},
	}
	for _, c := range casos {
		rec := pedir(t, h, c.metodo, c.ruta, c.cuerpo)
		var e errorJSON
		leerJSON(t, rec, &e)
		if rec.Code != c.codigo || e.Error != c.mensaje {
			t.Errorf("%s: código = %d, error = %q; se esperaba %d y %q", c.nombre, rec.Code, e.Error, c.codigo, c.mensaje)
		}
	}
}

// spec 12 CA-3, CA-4, CA-5: las métricas de esfuerzo por historia, por Sprint y del proyecto.
func TestConsultarMetricasEsfuerzo(t *testing.T) {
	h := conHistoriaEIntegrante(t) // hoy es 2026-11-20
	pedir(t, h, http.MethodPost, "/api/proyectos/1/backlog", `{"titulo":"Ver catálogo","descripcion":"d","prioridad":2,"criterios_aceptacion":""}`)
	pedir(t, h, http.MethodPut, rutaHorasEstimadas, `{"horas_estimadas":8}`)
	pedir(t, h, http.MethodPost, rutaEsfuerzo, `{"integrante_id":1,"fecha":"2026-11-18","actividad":"Programar","horas":10}`)
	pedir(t, h, http.MethodPost, "/api/proyectos/1/backlog/2/esfuerzo", `{"integrante_id":1,"fecha":"2026-11-18","actividad":"Diseñar","horas":2.5}`)
	pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints", sprintJSONDePrueba)
	pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints/1/historias", `{"historia_id":1}`)
	pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints/1/historias/1/completar", "")

	rec := pedir(t, h, http.MethodGet, "/api/proyectos/1/metricas/esfuerzo", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200 (%s)", rec.Code, rec.Body.String())
	}
	var m metricasEsfuerzoJSON
	leerJSON(t, rec, &m)

	// Historia 1: completada, 8 h estimadas y 10 h reales → +2 h, +25 %.
	h1 := m.Historias[0]
	if h1.HistoriaID != 1 || h1.Estado != "Completada" || h1.HorasEstimadas == nil || *h1.HorasEstimadas != "8" ||
		h1.HorasReales != "10" || h1.Desviacion == nil || h1.Desviacion.Horas != "2" || h1.Desviacion.Porcentaje != 25 {
		t.Errorf("historia 1 = %+v, se esperaba 8 h estimadas, 10 h reales y +2 h / +25 %%", h1)
	}
	// Historia 2: sin horas estimadas y sin terminar.
	if h2 := m.Historias[1]; h2.HorasEstimadas != nil || h2.HorasReales != "2.5" || h2.Desviacion != nil {
		t.Errorf("historia 2 = %+v, se esperaba sin horas estimadas, 2.5 h reales y sin desviación", h2)
	}
	if s := m.Sprints[0]; !s.Parcial || s.HorasEstimadas != "8" || s.HorasReales != "10" || s.Desviacion == nil || s.Desviacion.Porcentaje != 25 {
		t.Errorf("Sprint 1 = %+v, se esperaba parcial con 8 h estimadas, 10 h reales y +25 %%", s)
	}
	if m.HorasEstimadas != "8" || m.HorasReales != "12.5" || m.Desviacion == nil || m.Desviacion.Horas != "2" || m.CompletadasSinHorasEstimadas != 0 {
		t.Errorf("proyecto = %+v, se esperaba 8 h estimadas, 12.5 h reales y +2 h", m)
	}
}

// spec 12 RN-6, RN-10: la desviación negativa sale con signo, y sin datos sale null.
func TestMetricasEsfuerzoNegativaYSinDatos(t *testing.T) {
	h := conHistoriaEIntegrante(t)
	rec := pedir(t, h, http.MethodGet, "/api/proyectos/1/metricas/esfuerzo", "")
	if !strings.Contains(rec.Body.String(), `"desviacion":null`) || !strings.Contains(rec.Body.String(), `"sprints":[]`) {
		t.Errorf("sin datos: cuerpo = %s, se esperaba \"desviacion\":null y \"sprints\":[]", rec.Body.String())
	}

	pedir(t, h, http.MethodPut, rutaHorasEstimadas, `{"horas_estimadas":8}`)
	pedir(t, h, http.MethodPost, rutaEsfuerzo, `{"integrante_id":1,"fecha":"2026-11-18","actividad":"Programar","horas":7.5}`)
	pedir(t, h, http.MethodPut, "/api/proyectos/1/backlog/1", `{"titulo":"Registrar usuario","descripcion":"Alta de clientes","prioridad":1,"criterios_aceptacion":"","estado":"Completada"}`)
	rec = pedir(t, h, http.MethodGet, "/api/proyectos/1/metricas/esfuerzo", "")
	if !strings.Contains(rec.Body.String(), `"desviacion":{"horas":-0.5,"porcentaje":-6}`) {
		t.Errorf("negativa: cuerpo = %s, se esperaba \"desviacion\":{\"horas\":-0.5,\"porcentaje\":-6}", rec.Body.String())
	}
}

// spec 12 CA-7: proyecto inexistente → 404.
func TestMetricasEsfuerzoProyectoInexistente(t *testing.T) {
	rec := pedir(t, routerDePrueba(), http.MethodGet, "/api/proyectos/999/metricas/esfuerzo", "")
	var e errorJSON
	leerJSON(t, rec, &e)
	if rec.Code != http.StatusNotFound || e.Error != "El proyecto no existe" {
		t.Errorf("código = %d, error = %q; se esperaba 404", rec.Code, e.Error)
	}
}
