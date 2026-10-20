package api

import (
	"net/http"
	"strings"
	"testing"
)

const historiaJSONDePrueba = `{"titulo":"Registrar usuario","descripcion":"Alta de clientes","prioridad":1,"criterios_aceptacion":"Valida el email"}`

// conProyecto devuelve un router con el proyecto 1 ("Tienda Online") ya creado.
func conProyecto(t *testing.T) http.Handler {
	t.Helper()
	h := routerDePrueba()
	if rec := pedir(t, h, http.MethodPost, "/api/proyectos", tiendaJSON); rec.Code != http.StatusCreated {
		t.Fatalf("crear proyecto de partida: %d %s", rec.Code, rec.Body.String())
	}
	return h
}

// spec 04 CA-1: 201 con la historia Pendiente y sin estimar (story_points null).
func TestCrearHistoria(t *testing.T) {
	rec := pedir(t, conProyecto(t), http.MethodPost, "/api/proyectos/1/backlog", historiaJSONDePrueba)

	if rec.Code != http.StatusCreated {
		t.Fatalf("código = %d, se esperaba 201 (%s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"story_points":null`) {
		t.Errorf("cuerpo = %s, se esperaba \"story_points\":null (sin estimar)", rec.Body.String())
	}
	var h historiaJSON
	leerJSON(t, rec, &h)
	esperado := historiaJSON{ID: 1, ProyectoID: 1, Titulo: "Registrar usuario", Descripcion: "Alta de clientes", Prioridad: 1, Criterios: "Valida el email", Estado: "Pendiente"}
	if h.StoryPoints != nil {
		t.Errorf("story_points = %d, se esperaba sin estimar", *h.StoryPoints)
	}
	h.StoryPoints = nil
	if h != esperado {
		t.Errorf("historia = %+v, se esperaba %+v", h, esperado)
	}
}

// spec 04 CA-5, CA-6: datos inválidos → 400 con el mensaje de la spec.
func TestCrearHistoriaInvalida(t *testing.T) {
	casos := []struct{ nombre, cuerpo, mensaje string }{
		{"título vacío (CA-5)", `{"titulo":"","descripcion":"d","prioridad":1}`, "El título es obligatorio"},
		{"sin descripción (CA-5)", `{"titulo":"t","prioridad":1}`, "La descripción es obligatoria"},
		{"sin prioridad (CA-6)", `{"titulo":"t","descripcion":"d"}`, "La prioridad es obligatoria"},
		{"prioridad 0 (CA-6)", `{"titulo":"t","descripcion":"d","prioridad":0}`, "La prioridad debe ser un número entero mayor o igual a 1"},
		{"prioridad con decimales (CA-6)", `{"titulo":"t","descripcion":"d","prioridad":2.5}`, "La prioridad debe ser un número entero mayor o igual a 1"},
		{"JSON mal formado", `{"titulo":`, "El cuerpo de la solicitud no es un JSON válido"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedir(t, conProyecto(t), http.MethodPost, "/api/proyectos/1/backlog", c.cuerpo)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("código = %d, se esperaba 400 (%s)", rec.Code, rec.Body.String())
			}
			var e errorJSON
			leerJSON(t, rec, &e)
			if e.Error != c.mensaje {
				t.Errorf("error = %q, se esperaba %q", e.Error, c.mensaje)
			}
		})
	}
}

// spec 04 CA-2, CA-4, CA-6: modificar (incluido el estado), ver y validar el estado.
func TestModificarYVerHistoria(t *testing.T) {
	h := conProyecto(t)
	pedir(t, h, http.MethodPost, "/api/proyectos/1/backlog", historiaJSONDePrueba)

	rec := pedir(t, h, http.MethodPut, "/api/proyectos/1/backlog/1",
		`{"titulo":"Registrar cliente","descripcion":"Alta de clientes","prioridad":2,"criterios_aceptacion":"","estado":"En progreso"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT: código = %d, se esperaba 200 (%s)", rec.Code, rec.Body.String())
	}

	rec = pedir(t, h, http.MethodGet, "/api/proyectos/1/backlog/1", "")
	var vista historiaJSON
	leerJSON(t, rec, &vista)
	if rec.Code != http.StatusOK || vista.Titulo != "Registrar cliente" || vista.Prioridad != 2 || vista.Estado != "En progreso" {
		t.Errorf("GET: código = %d, historia = %+v", rec.Code, vista)
	}

	rec = pedir(t, h, http.MethodPut, "/api/proyectos/1/backlog/1",
		`{"titulo":"Registrar cliente","descripcion":"Alta de clientes","prioridad":2,"estado":"Terminada"}`)
	var e errorJSON
	leerJSON(t, rec, &e)
	if rec.Code != http.StatusBadRequest || e.Error != "El estado no es válido" {
		t.Errorf("PUT estado inválido: código = %d, error = %q", rec.Code, e.Error)
	}
}

// spec 04 CA-3: eliminar → 204 sin cuerpo; después no existe.
func TestEliminarHistoria(t *testing.T) {
	h := conProyecto(t)
	pedir(t, h, http.MethodPost, "/api/proyectos/1/backlog", historiaJSONDePrueba)

	rec := pedir(t, h, http.MethodDelete, "/api/proyectos/1/backlog/1", "")
	if rec.Code != http.StatusNoContent || rec.Body.Len() != 0 {
		t.Fatalf("DELETE: código = %d, cuerpo = %q; se esperaba 204 sin cuerpo", rec.Code, rec.Body.String())
	}
	for _, metodo := range []string{http.MethodGet, http.MethodDelete} {
		rec = pedir(t, h, metodo, "/api/proyectos/1/backlog/1", "")
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s después de eliminar: código = %d, se esperaba 404", metodo, rec.Code)
		}
	}
}

// spec 04 CA-4: el backlog se lista ordenado por prioridad; vacío es [] (no null).
func TestListarBacklog(t *testing.T) {
	h := conProyecto(t)
	rec := pedir(t, h, http.MethodGet, "/api/proyectos/1/backlog", "")
	if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Fatalf("backlog vacío: código = %d, cuerpo = %q; se esperaba 200 y []", rec.Code, rec.Body.String())
	}

	pedir(t, h, http.MethodPost, "/api/proyectos/1/backlog", `{"titulo":"Pagar pedido","descripcion":"d","prioridad":3}`)
	pedir(t, h, http.MethodPost, "/api/proyectos/1/backlog", `{"titulo":"Registrar usuario","descripcion":"d","prioridad":1}`)
	pedir(t, h, http.MethodPost, "/api/proyectos/1/backlog", `{"titulo":"Ver catálogo","descripcion":"d","prioridad":2}`)

	rec = pedir(t, h, http.MethodGet, "/api/proyectos/1/backlog", "")
	var lista []historiaJSON
	leerJSON(t, rec, &lista)
	if len(lista) != 3 || lista[0].Titulo != "Registrar usuario" || lista[1].Titulo != "Ver catálogo" || lista[2].Titulo != "Pagar pedido" {
		t.Errorf("backlog = %+v, se esperaba ordenado por prioridad", lista)
	}
}

// spec 04 CA-7: proyecto o historia que no existe → 404 con el mensaje que corresponde.
func TestBacklogInexistente(t *testing.T) {
	h := conProyecto(t)
	casos := []struct{ nombre, metodo, ruta, cuerpo, mensaje string }{
		{"crear en proyecto inexistente", http.MethodPost, "/api/proyectos/999/backlog", historiaJSONDePrueba, "El proyecto no existe"},
		{"listar proyecto inexistente", http.MethodGet, "/api/proyectos/999/backlog", "", "El proyecto no existe"},
		{"proyecto que no es un número", http.MethodGet, "/api/proyectos/abc/backlog", "", "El proyecto no existe"},
		{"ver historia inexistente", http.MethodGet, "/api/proyectos/1/backlog/999", "", "La historia no existe"},
		{"modificar historia inexistente", http.MethodPut, "/api/proyectos/1/backlog/999", `{"titulo":"t","descripcion":"d","prioridad":1,"estado":"Pendiente"}`, "La historia no existe"},
		{"eliminar historia inexistente", http.MethodDelete, "/api/proyectos/1/backlog/999", "", "La historia no existe"},
		{"historia que no es un número", http.MethodGet, "/api/proyectos/1/backlog/abc", "", "La historia no existe"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedir(t, h, c.metodo, c.ruta, c.cuerpo)
			var e errorJSON
			leerJSON(t, rec, &e)
			if rec.Code != http.StatusNotFound || e.Error != c.mensaje {
				t.Errorf("código = %d, error = %q; se esperaba 404 y %q", rec.Code, e.Error, c.mensaje)
			}
		})
	}
}
