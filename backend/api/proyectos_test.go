package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/memoria"
)

func routerDePrueba() http.Handler {
	return NewRouter(dominio.NuevoServicioProyectos(memoria.NuevoRepositorioProyectos()))
}

func pedir(t *testing.T, h http.Handler, metodo, ruta, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(metodo, ruta, strings.NewReader(cuerpo))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func leerJSON(t *testing.T, rec *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), v); err != nil {
		t.Fatalf("respuesta no es JSON válido (%s): %v", rec.Body.String(), err)
	}
}

const tiendaJSON = `{"nombre":"Tienda Online","descripcion":"Venta de ropa","fecha_inicio":"2026-11-01","fecha_fin":"2026-12-15"}`

// spec 01 CA-1
func TestCrearProyecto(t *testing.T) {
	rec := pedir(t, routerDePrueba(), http.MethodPost, "/api/proyectos", tiendaJSON)

	if rec.Code != http.StatusCreated {
		t.Fatalf("código = %d, se esperaba 201 (%s)", rec.Code, rec.Body.String())
	}
	var p proyectoJSON
	leerJSON(t, rec, &p)
	esperado := proyectoJSON{ID: 1, Nombre: "Tienda Online", Descripcion: "Venta de ropa", FechaInicio: "2026-11-01", FechaFin: "2026-12-15"}
	if p != esperado {
		t.Errorf("proyecto = %+v, se esperaba %+v", p, esperado)
	}
}

// spec 01 CA-3, CA-4, CA-5: datos inválidos → 400 con el mensaje de la spec.
func TestCrearProyectoInvalido(t *testing.T) {
	casos := []struct{ nombre, cuerpo, mensaje string }{
		{"nombre vacío (CA-4)", `{"nombre":"","descripcion":"d","fecha_inicio":"2026-11-01","fecha_fin":"2026-12-15"}`, "El nombre es obligatorio"},
		{"sin descripción (CA-5)", `{"nombre":"n","fecha_inicio":"2026-11-01","fecha_fin":"2026-12-15"}`, "La descripción es obligatoria"},
		{"sin fecha de fin (CA-3)", `{"nombre":"n","descripcion":"d","fecha_inicio":"2026-11-01"}`, "Las fechas de inicio y finalización son obligatorias"},
		{"fecha inexistente (CA-3)", `{"nombre":"n","descripcion":"d","fecha_inicio":"2026-02-30","fecha_fin":"2026-03-10"}`, "La fecha no es válida"},
		{"fin igual a inicio (CA-3)", `{"nombre":"n","descripcion":"d","fecha_inicio":"2026-11-10","fecha_fin":"2026-11-10"}`, "La fecha de finalización debe ser posterior a la de inicio"},
		{"JSON mal formado", `{"nombre":`, "El cuerpo de la solicitud no es un JSON válido"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedir(t, routerDePrueba(), http.MethodPost, "/api/proyectos", c.cuerpo)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("código = %d, se esperaba 400", rec.Code)
			}
			var e errorJSON
			leerJSON(t, rec, &e)
			if e.Error != c.mensaje {
				t.Errorf("error = %q, se esperaba %q", e.Error, c.mensaje)
			}
		})
	}
}

// spec 01 CA-2, CA-6
func TestModificarListarYVer(t *testing.T) {
	h := routerDePrueba()
	pedir(t, h, http.MethodPost, "/api/proyectos", tiendaJSON)

	rec := pedir(t, h, http.MethodPut, "/api/proyectos/1",
		`{"nombre":"Tienda Online v2","descripcion":"Venta de ropa","fecha_inicio":"2026-11-10","fecha_fin":"2027-01-31"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT: código = %d, se esperaba 200 (%s)", rec.Code, rec.Body.String())
	}

	rec = pedir(t, h, http.MethodGet, "/api/proyectos/1", "")
	var p proyectoJSON
	leerJSON(t, rec, &p)
	if rec.Code != http.StatusOK || p.Nombre != "Tienda Online v2" || p.FechaFin != "2027-01-31" {
		t.Errorf("GET: código = %d, proyecto = %+v", rec.Code, p)
	}

	rec = pedir(t, h, http.MethodGet, "/api/proyectos", "")
	var lista []proyectoJSON
	leerJSON(t, rec, &lista)
	if rec.Code != http.StatusOK || len(lista) != 1 || lista[0] != p {
		t.Errorf("GET lista: código = %d, lista = %+v", rec.Code, lista)
	}
}

// spec 01 CA-6: sin proyectos, la respuesta es una lista vacía [] (no null).
func TestListarVacio(t *testing.T) {
	rec := pedir(t, routerDePrueba(), http.MethodGet, "/api/proyectos", "")
	if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Errorf("código = %d, cuerpo = %q; se esperaba 200 y []", rec.Code, rec.Body.String())
	}
}

// spec 01 CA-7: ver o modificar un proyecto que no existe → 404.
func TestProyectoInexistente(t *testing.T) {
	h := routerDePrueba()
	casos := []struct{ nombre, metodo, ruta, cuerpo string }{
		{"ver", http.MethodGet, "/api/proyectos/999", ""},
		{"modificar", http.MethodPut, "/api/proyectos/999", tiendaJSON},
		{"id que no es un número", http.MethodGet, "/api/proyectos/abc", ""},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedir(t, h, c.metodo, c.ruta, c.cuerpo)
			var e errorJSON
			leerJSON(t, rec, &e)
			if rec.Code != http.StatusNotFound || e.Error != "El proyecto no existe" {
				t.Errorf("código = %d, error = %q; se esperaba 404 y \"El proyecto no existe\"", rec.Code, e.Error)
			}
		})
	}
}
