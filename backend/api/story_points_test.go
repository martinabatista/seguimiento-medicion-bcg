package api

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// conHistoria devuelve un router con el proyecto 1 y su historia 1 ("Registrar usuario") sin estimar.
func conHistoria(t *testing.T) http.Handler {
	t.Helper()
	h := conProyecto(t)
	if rec := pedir(t, h, http.MethodPost, "/api/proyectos/1/backlog", historiaJSONDePrueba); rec.Code != http.StatusCreated {
		t.Fatalf("crear historia de partida: %d %s", rec.Code, rec.Body.String())
	}
	return h
}

const rutaStoryPoints = "/api/proyectos/1/backlog/1/story-points"

// spec 07 CA-1, CA-2: estimar devuelve la historia con sus Story Points; cambiar los reemplaza.
func TestEstimarHistoria(t *testing.T) {
	h := conHistoria(t)

	for _, sp := range []int{5, 13} {
		rec := pedir(t, h, http.MethodPut, rutaStoryPoints, `{"story_points":`+strconv.Itoa(sp)+`}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("PUT %d: código = %d, se esperaba 200 (%s)", sp, rec.Code, rec.Body.String())
		}
		var hist historiaJSON
		leerJSON(t, rec, &hist)
		if hist.StoryPoints == nil || *hist.StoryPoints != sp || hist.Titulo != "Registrar usuario" || hist.Estado != "Pendiente" {
			t.Errorf("historia = %+v, se esperaba %d Story Points y el resto sin cambios", hist, sp)
		}
	}
	rec := pedir(t, h, http.MethodGet, "/api/proyectos/1/backlog/1", "")
	var vista historiaJSON
	leerJSON(t, rec, &vista)
	if vista.StoryPoints == nil || *vista.StoryPoints != 13 {
		t.Errorf("GET: story_points = %v, se esperaba 13", vista.StoryPoints)
	}
}

// spec 07 CA-2: quitar la estimación devuelve la historia con story_points null.
func TestQuitarEstimacionHistoria(t *testing.T) {
	h := conHistoria(t)
	pedir(t, h, http.MethodPut, rutaStoryPoints, `{"story_points":3}`)

	rec := pedir(t, h, http.MethodDelete, rutaStoryPoints, "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"story_points":null`) {
		t.Errorf("DELETE: código = %d, cuerpo = %s; se esperaba 200 y story_points null", rec.Code, rec.Body.String())
	}
}

// spec 07 CA-3: valores fuera de la escala → 400 con el mensaje de la spec.
func TestEstimarHistoriaInvalida(t *testing.T) {
	const fueraDeEscala = "Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21"
	casos := []struct{ nombre, cuerpo, mensaje string }{
		{"falta el valor", `{}`, "Los Story Points son obligatorios"},
		{"cero", `{"story_points":0}`, fueraDeEscala},
		{"entre valores de la escala", `{"story_points":4}`, fueraDeEscala},
		{"mayor que 21", `{"story_points":22}`, fueraDeEscala},
		{"con decimales", `{"story_points":2.5}`, fueraDeEscala},
		{"JSON mal formado", `{"story_points":`, "El cuerpo de la solicitud no es un JSON válido"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedir(t, conHistoria(t), http.MethodPut, rutaStoryPoints, c.cuerpo)
			var e errorJSON
			leerJSON(t, rec, &e)
			if rec.Code != http.StatusBadRequest || e.Error != c.mensaje {
				t.Errorf("código = %d, error = %q; se esperaba 400 y %q", rec.Code, e.Error, c.mensaje)
			}
		})
	}
}

// spec 07 CA-4: una historia completada no se estima ni se le quita la estimación → 409.
func TestEstimarHistoriaCompletada(t *testing.T) {
	h := conHistoria(t)
	pedir(t, h, http.MethodPut, "/api/proyectos/1/backlog/1",
		`{"titulo":"Registrar usuario","descripcion":"Alta de clientes","prioridad":1,"estado":"Completada"}`)

	for _, metodo := range []string{http.MethodPut, http.MethodDelete} {
		rec := pedir(t, h, metodo, rutaStoryPoints, `{"story_points":5}`)
		var e errorJSON
		leerJSON(t, rec, &e)
		if rec.Code != http.StatusConflict || e.Error != "La historia ya está completada" {
			t.Errorf("%s: código = %d, error = %q; se esperaba 409", metodo, rec.Code, e.Error)
		}
	}
}

// spec 07 CA-6: proyecto o historia que no existe → 404.
func TestEstimarHistoriaInexistente(t *testing.T) {
	h := conHistoria(t)
	casos := []struct{ nombre, metodo, ruta, mensaje string }{
		{"proyecto que no existe", http.MethodPut, "/api/proyectos/999/backlog/1/story-points", "El proyecto no existe"},
		{"historia que no existe", http.MethodPut, "/api/proyectos/1/backlog/999/story-points", "La historia no existe"},
		{"quitar en historia que no existe", http.MethodDelete, "/api/proyectos/1/backlog/999/story-points", "La historia no existe"},
		{"historia que no es un número", http.MethodPut, "/api/proyectos/1/backlog/abc/story-points", "La historia no existe"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedir(t, h, c.metodo, c.ruta, `{"story_points":3}`)
			var e errorJSON
			leerJSON(t, rec, &e)
			if rec.Code != http.StatusNotFound || e.Error != c.mensaje {
				t.Errorf("código = %d, error = %q; se esperaba 404 y %q", rec.Code, e.Error, c.mensaje)
			}
		})
	}
}

// spec 07 CA-5: la escala se consulta de menor a mayor.
func TestConsultarEscalaStoryPoints(t *testing.T) {
	rec := pedir(t, routerDePrueba(), http.MethodGet, "/api/story-points/escala", "")
	if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != `{"valores":[1,2,3,5,8,13,21]}` {
		t.Errorf("código = %d, cuerpo = %s; se esperaba 200 y la escala", rec.Code, rec.Body.String())
	}
}
