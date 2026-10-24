package api

import (
	"net/http"
	"strings"
	"testing"
)

const rutaEsfuerzo = "/api/proyectos/1/backlog/1/esfuerzo"

// conHistoriaEIntegrante devuelve un router con el proyecto 1, su historia 1 y su integrante 1
// (Ana Pérez). Hoy es hoyDePrueba (2026-11-20).
func conHistoriaEIntegrante(t *testing.T) http.Handler {
	t.Helper()
	h := conProyecto(t)
	if rec := pedir(t, h, http.MethodPost, "/api/proyectos/1/backlog", historiaJSONDePrueba); rec.Code != http.StatusCreated {
		t.Fatalf("crear historia de partida: %d %s", rec.Code, rec.Body.String())
	}
	if rec := pedir(t, h, http.MethodPost, "/api/proyectos/1/integrantes", anaJSON); rec.Code != http.StatusCreated {
		t.Fatalf("crear integrante de partida: %d %s", rec.Code, rec.Body.String())
	}
	return h
}

// spec 09 CA-1, CA-2: 201 con el registro, el nombre del integrante y las horas como número.
func TestRegistrarEsfuerzo(t *testing.T) {
	rec := pedir(t, conHistoriaEIntegrante(t), http.MethodPost, rutaEsfuerzo,
		`{"integrante_id":1,"fecha":"2026-11-18","actividad":"Diseño del formulario","horas":2.5}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("código = %d, se esperaba 201 (%s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"horas":2.5`) {
		t.Errorf("cuerpo = %s, se esperaba \"horas\":2.5 como número", rec.Body.String())
	}
	var r registroEsfuerzoJSON
	leerJSON(t, rec, &r)
	esperado := registroEsfuerzoJSON{ID: 1, ProyectoID: 1, HistoriaID: 1, IntegranteID: 1, IntegranteNombre: "Ana Pérez",
		Fecha: "2026-11-18", Actividad: "Diseño del formulario", Horas: "2.5"}
	if r != esperado {
		t.Errorf("registro = %+v, se esperaba %+v", r, esperado)
	}
}

// spec 09 CA-3, CA-5, CA-8: datos inválidos → 400 con el mensaje de la spec.
func TestRegistrarEsfuerzoInvalido(t *testing.T) {
	casos := []struct{ nombre, cuerpo, mensaje string }{
		{"sin integrante (CA-8)", `{"fecha":"2026-11-18","actividad":"A","horas":1}`, "El integrante es obligatorio"},
		{"sin fecha (CA-5)", `{"integrante_id":1,"actividad":"A","horas":1}`, "La fecha es obligatoria"},
		{"fecha futura (CA-5)", `{"integrante_id":1,"fecha":"2026-11-21","actividad":"A","horas":1}`, "La fecha no puede ser posterior a hoy"},
		{"fecha fuera del proyecto (CA-5)", `{"integrante_id":1,"fecha":"2026-10-31","actividad":"A","horas":1}`, "La fecha debe estar dentro de las fechas del proyecto"},
		{"sin actividad (CA-5)", `{"integrante_id":1,"fecha":"2026-11-18","horas":1}`, "La actividad es obligatoria"},
		{"sin horas (CA-3)", `{"integrante_id":1,"fecha":"2026-11-18","actividad":"A"}`, "Las horas son obligatorias"},
		{"horas en cero (CA-3)", `{"integrante_id":1,"fecha":"2026-11-18","actividad":"A","horas":0}`, "Las horas deben ser un número mayor a 0 y hasta 24, con hasta 2 decimales"},
		{"tres decimales (CA-3)", `{"integrante_id":1,"fecha":"2026-11-18","actividad":"A","horas":1.125}`, "Las horas deben ser un número mayor a 0 y hasta 24, con hasta 2 decimales"},
		{"horas como texto (CA-3)", `{"integrante_id":1,"fecha":"2026-11-18","actividad":"A","horas":"dos"}`, "El cuerpo de la solicitud no es un JSON válido"},
		{"JSON mal formado", `{"integrante_id":`, "El cuerpo de la solicitud no es un JSON válido"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedir(t, conHistoriaEIntegrante(t), http.MethodPost, rutaEsfuerzo, c.cuerpo)
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

// spec 09 CA-4, CA-6: consultar por fecha con total exacto; eliminar (204) y deja de sumar.
func TestConsultarYEliminarEsfuerzo(t *testing.T) {
	h := conHistoriaEIntegrante(t)

	rec := pedir(t, h, http.MethodGet, rutaEsfuerzo, "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"registros":[]`) || !strings.Contains(rec.Body.String(), `"total_horas":0`) {
		t.Fatalf("sin registros: código = %d, cuerpo = %s; se esperaba [] y total 0", rec.Code, rec.Body.String())
	}

	pedir(t, h, http.MethodPost, rutaEsfuerzo, `{"integrante_id":1,"fecha":"2026-11-19","actividad":"Backend","horas":0.1}`)
	pedir(t, h, http.MethodPost, rutaEsfuerzo, `{"integrante_id":1,"fecha":"2026-11-18","actividad":"Diseño","horas":0.2}`)

	rec = pedir(t, h, http.MethodGet, rutaEsfuerzo, "")
	var e esfuerzoJSON
	leerJSON(t, rec, &e)
	if len(e.Registros) != 2 || e.Registros[0].Actividad != "Diseño" || e.Registros[1].Actividad != "Backend" {
		t.Errorf("registros = %+v, se esperaba Diseño y Backend en ese orden", e.Registros)
	}
	if e.TotalHoras != "0.3" {
		t.Errorf("total = %s, se esperaba exactamente 0.3", e.TotalHoras)
	}

	if rec := pedir(t, h, http.MethodDelete, rutaEsfuerzo+"/1", ""); rec.Code != http.StatusNoContent {
		t.Fatalf("DELETE: código = %d, se esperaba 204 (%s)", rec.Code, rec.Body.String())
	}
	rec = pedir(t, h, http.MethodGet, rutaEsfuerzo, "")
	leerJSON(t, rec, &e)
	if len(e.Registros) != 1 || e.TotalHoras != "0.2" {
		t.Errorf("después de eliminar = %+v, se esperaba solo Diseño con 0.2", e)
	}
}

// spec 09 CA-7: eliminar del backlog una historia con esfuerzo → 409; sin esfuerzo, 204.
func TestEliminarHistoriaConEsfuerzo(t *testing.T) {
	h := conHistoriaEIntegrante(t)
	pedir(t, h, http.MethodPost, rutaEsfuerzo, `{"integrante_id":1,"fecha":"2026-11-18","actividad":"Backend","horas":2}`)

	rec := pedir(t, h, http.MethodDelete, "/api/proyectos/1/backlog/1", "")
	if rec.Code != http.StatusConflict {
		t.Fatalf("código = %d, se esperaba 409 (%s)", rec.Code, rec.Body.String())
	}
	var e errorJSON
	leerJSON(t, rec, &e)
	if e.Error != "La historia tiene esfuerzo registrado" {
		t.Errorf("error = %q", e.Error)
	}

	pedir(t, h, http.MethodDelete, rutaEsfuerzo+"/1", "")
	if rec := pedir(t, h, http.MethodDelete, "/api/proyectos/1/backlog/1", ""); rec.Code != http.StatusNoContent {
		t.Errorf("sin esfuerzo: código = %d, se esperaba 204 (%s)", rec.Code, rec.Body.String())
	}
}

// spec 09 CA-8: proyecto, historia, integrante o registro inexistente → 404 con el mensaje de cuál.
func TestEsfuerzoNoExiste(t *testing.T) {
	h := conHistoriaEIntegrante(t)
	valido := `{"integrante_id":1,"fecha":"2026-11-18","actividad":"A","horas":1}`
	casos := []struct{ nombre, metodo, ruta, cuerpo, mensaje string }{
		{"proyecto inexistente", http.MethodGet, "/api/proyectos/999/backlog/1/esfuerzo", "", "El proyecto no existe"},
		{"historia inexistente", http.MethodPost, "/api/proyectos/1/backlog/999/esfuerzo", valido, "La historia no existe"},
		{"historia con id que no es un número", http.MethodGet, "/api/proyectos/1/backlog/abc/esfuerzo", "", "La historia no existe"},
		{"integrante inexistente", http.MethodPost, rutaEsfuerzo, `{"integrante_id":999,"fecha":"2026-11-18","actividad":"A","horas":1}`, "El integrante no existe"},
		{"registro inexistente", http.MethodDelete, rutaEsfuerzo + "/999", "", "El registro de esfuerzo no existe"},
		{"registro con id que no es un número", http.MethodDelete, rutaEsfuerzo + "/abc", "", "El registro de esfuerzo no existe"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedir(t, h, c.metodo, c.ruta, c.cuerpo)
			if rec.Code != http.StatusNotFound {
				t.Fatalf("código = %d, se esperaba 404 (%s)", rec.Code, rec.Body.String())
			}
			var e errorJSON
			leerJSON(t, rec, &e)
			if e.Error != c.mensaje {
				t.Errorf("error = %q, se esperaba %q", e.Error, c.mensaje)
			}
		})
	}
}
