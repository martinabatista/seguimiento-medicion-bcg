package api

import (
	"net/http"
	"testing"
)

// spec 03 CA-1, CA-3: 200 con los datos del proyecto, su estado y el resumen del backlog.
func TestConsultarEstado(t *testing.T) {
	h := conProyecto(t) // "Tienda Online", del 2026-11-01 al 2026-12-15; hoy es 2026-11-20
	pedir(t, h, http.MethodPost, "/api/proyectos/1/backlog", historiaJSONDePrueba)
	pedir(t, h, http.MethodPost, "/api/proyectos/1/backlog", `{"titulo":"Ver catálogo","descripcion":"Listado","prioridad":2}`)
	pedir(t, h, http.MethodPut, "/api/proyectos/1/backlog/1",
		`{"titulo":"Registrar usuario","descripcion":"Alta de clientes","prioridad":1,"estado":"Completada"}`)

	rec := pedir(t, h, http.MethodGet, "/api/proyectos/1/estado", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200 (%s)", rec.Code, rec.Body.String())
	}
	var e estadoJSON
	leerJSON(t, rec, &e)
	esperado := estadoJSON{
		Proyecto: proyectoJSON{ID: 1, Nombre: "Tienda Online", Descripcion: "Venta de ropa", FechaInicio: "2026-11-01", FechaFin: "2026-12-15"},
		Estado:   "En curso",
		Backlog: resumenBacklogJSON{Total: 2, Pendientes: 1, Completadas: 1, PorcentajeCompletadas: 50,
			StoryPointsTotales: 0, StoryPointsCompletados: 0, SinEstimar: 2},
	}
	if e != esperado {
		t.Errorf("estado = %+v, se esperaba %+v", e, esperado)
	}
}

// spec 03 CA-3: con el backlog vacío, el resumen viene en cero (no es un error).
func TestConsultarEstadoBacklogVacio(t *testing.T) {
	rec := pedir(t, conProyecto(t), http.MethodGet, "/api/proyectos/1/estado", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200 (%s)", rec.Code, rec.Body.String())
	}
	var e estadoJSON
	leerJSON(t, rec, &e)
	if e.Backlog != (resumenBacklogJSON{}) {
		t.Errorf("resumen = %+v, se esperaba en cero", e.Backlog)
	}
}

// spec 03 CA-4: proyecto inexistente (o con id que no es un número) → 404.
func TestConsultarEstadoProyectoInexistente(t *testing.T) {
	for _, ruta := range []string{"/api/proyectos/999/estado", "/api/proyectos/abc/estado"} {
		t.Run(ruta, func(t *testing.T) {
			rec := pedir(t, routerDePrueba(), http.MethodGet, ruta, "")
			if rec.Code != http.StatusNotFound {
				t.Fatalf("código = %d, se esperaba 404 (%s)", rec.Code, rec.Body.String())
			}
			var e errorJSON
			leerJSON(t, rec, &e)
			if e.Error != "El proyecto no existe" {
				t.Errorf("error = %q", e.Error)
			}
		})
	}
}
