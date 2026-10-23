package api

import (
	"net/http"
	"testing"
)

const anaJSON = `{"nombre":"Ana Pérez","email":"ana@mail.com","rol":"Developer"}`

// spec 02 CA-1: 201 con el integrante y el email normalizado.
func TestAgregarIntegrante(t *testing.T) {
	rec := pedir(t, conProyecto(t), http.MethodPost, "/api/proyectos/1/integrantes",
		`{"nombre":"Ana Pérez","email":" Ana@Mail.COM ","rol":"Developer"}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("código = %d, se esperaba 201 (%s)", rec.Code, rec.Body.String())
	}
	var i integranteJSON
	leerJSON(t, rec, &i)
	esperado := integranteJSON{ID: 1, ProyectoID: 1, Nombre: "Ana Pérez", Email: "ana@mail.com", Rol: "Developer"}
	if i != esperado {
		t.Errorf("integrante = %+v, se esperaba %+v", i, esperado)
	}
}

// spec 02 CA-4, CA-5: datos inválidos → 400 con el mensaje de la spec.
func TestAgregarIntegranteInvalido(t *testing.T) {
	casos := []struct{ nombre, cuerpo, mensaje string }{
		{"nombre vacío (CA-4)", `{"nombre":"","email":"ana@mail.com","rol":"Developer"}`, "El nombre es obligatorio"},
		{"sin email (CA-4)", `{"nombre":"Ana","rol":"Developer"}`, "El email es obligatorio"},
		{"email inválido (CA-4)", `{"nombre":"Ana","email":"ana.mail.com","rol":"Developer"}`, "El email no es válido"},
		{"sin rol (CA-5)", `{"nombre":"Ana","email":"ana@mail.com"}`, "El rol es obligatorio"},
		{"rol inválido (CA-5)", `{"nombre":"Ana","email":"ana@mail.com","rol":"Tester"}`, "El rol no es válido"},
		{"JSON mal formado", `{"nombre":`, "El cuerpo de la solicitud no es un JSON válido"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedir(t, conProyecto(t), http.MethodPost, "/api/proyectos/1/integrantes", c.cuerpo)
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

// spec 02 CA-3: el mismo email en el proyecto → 409 Conflicto.
func TestAgregarIntegranteRepetido(t *testing.T) {
	h := conProyecto(t)
	pedir(t, h, http.MethodPost, "/api/proyectos/1/integrantes", anaJSON)

	rec := pedir(t, h, http.MethodPost, "/api/proyectos/1/integrantes",
		`{"nombre":"Ana P.","email":"ANA@mail.com","rol":"Scrum Master"}`)
	if rec.Code != http.StatusConflict {
		t.Fatalf("código = %d, se esperaba 409 (%s)", rec.Code, rec.Body.String())
	}
	var e errorJSON
	leerJSON(t, rec, &e)
	if e.Error != "El integrante ya está registrado en el proyecto" {
		t.Errorf("error = %q", e.Error)
	}
}

// spec 02 CA-2: listar en orden de alta ([] si no hay) y ver uno por su identificador.
func TestListarYVerIntegrantes(t *testing.T) {
	h := conProyecto(t)

	rec := pedir(t, h, http.MethodGet, "/api/proyectos/1/integrantes", "")
	if rec.Code != http.StatusOK || rec.Body.String() != "[]\n" {
		t.Fatalf("lista vacía: código = %d, cuerpo = %q; se esperaba 200 y []", rec.Code, rec.Body.String())
	}

	pedir(t, h, http.MethodPost, "/api/proyectos/1/integrantes", anaJSON)
	pedir(t, h, http.MethodPost, "/api/proyectos/1/integrantes", `{"nombre":"Bruno Díaz","email":"bruno@mail.com","rol":"Scrum Master"}`)

	rec = pedir(t, h, http.MethodGet, "/api/proyectos/1/integrantes", "")
	var lista []integranteJSON
	leerJSON(t, rec, &lista)
	if len(lista) != 2 || lista[0].Nombre != "Ana Pérez" || lista[1].Nombre != "Bruno Díaz" {
		t.Errorf("lista = %+v, se esperaban Ana Pérez y Bruno Díaz en ese orden", lista)
	}

	rec = pedir(t, h, http.MethodGet, "/api/proyectos/1/integrantes/2", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("ver: código = %d, se esperaba 200 (%s)", rec.Code, rec.Body.String())
	}
	var bruno integranteJSON
	leerJSON(t, rec, &bruno)
	esperado := integranteJSON{ID: 2, ProyectoID: 1, Nombre: "Bruno Díaz", Email: "bruno@mail.com", Rol: "Scrum Master"}
	if bruno != esperado {
		t.Errorf("integrante = %+v, se esperaba %+v", bruno, esperado)
	}
}

// spec 02 CA-6: proyecto o integrante inexistente → 404 con el mensaje de cuál no existe.
func TestIntegranteNoExiste(t *testing.T) {
	h := conProyecto(t)
	pedir(t, h, http.MethodPost, "/api/proyectos", `{"nombre":"App Turnos","descripcion":"Turnos","fecha_inicio":"2026-11-05","fecha_fin":"2027-01-20"}`)
	pedir(t, h, http.MethodPost, "/api/proyectos/1/integrantes", anaJSON)

	casos := []struct{ nombre, metodo, ruta, cuerpo, mensaje string }{
		{"agregar en proyecto inexistente", http.MethodPost, "/api/proyectos/999/integrantes", anaJSON, "El proyecto no existe"},
		{"listar proyecto inexistente", http.MethodGet, "/api/proyectos/999/integrantes", "", "El proyecto no existe"},
		{"ver en proyecto inexistente", http.MethodGet, "/api/proyectos/999/integrantes/1", "", "El proyecto no existe"},
		{"proyecto con id que no es un número", http.MethodGet, "/api/proyectos/abc/integrantes", "", "El proyecto no existe"},
		{"integrante inexistente", http.MethodGet, "/api/proyectos/1/integrantes/999", "", "El integrante no existe"},
		{"integrante con id que no es un número", http.MethodGet, "/api/proyectos/1/integrantes/abc", "", "El integrante no existe"},
		{"integrante de otro proyecto", http.MethodGet, "/api/proyectos/2/integrantes/1", "", "El integrante no existe"},
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
