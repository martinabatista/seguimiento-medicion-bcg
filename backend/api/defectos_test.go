package api

import (
	"net/http"
	"strings"
	"testing"
)

const (
	rutaDefectos = "/api/proyectos/1/defectos"
	rutaDefecto1 = rutaDefectos + "/1"
	// defectoJSONDePrueba es un defecto válido en la historia 1, detectado en el Sprint 1.
	defectoJSONDePrueba = `{"descripcion":"El alta acepta emails sin @","severidad":"Alta","historia_id":1,"sprint_deteccion_id":1}`
)

// conDosSprints devuelve un router con el proyecto 1, la historia 1 y dos Sprints: el 1 (del
// 2026-11-01 al 14) y el 2 (del 16 al 28).
func conDosSprints(t *testing.T) http.Handler {
	t.Helper()
	h := conProyecto(t)
	if rec := pedir(t, h, http.MethodPost, "/api/proyectos/1/backlog", historiaJSONDePrueba); rec.Code != http.StatusCreated {
		t.Fatalf("crear historia de partida: %d %s", rec.Code, rec.Body.String())
	}
	for _, s := range []string{
		`{"fecha_inicio":"2026-11-01","fecha_fin":"2026-11-14","sprint_goal":"Alta"}`,
		`{"fecha_inicio":"2026-11-16","fecha_fin":"2026-11-28","sprint_goal":"Compra"}`,
	} {
		if rec := pedir(t, h, http.MethodPost, "/api/proyectos/1/sprints", s); rec.Code != http.StatusCreated {
			t.Fatalf("crear Sprint de partida: %d %s", rec.Code, rec.Body.String())
		}
	}
	return h
}

// spec 11 CA-1, CA-5: 201 con el defecto Abierto, su historia, su Sprint de detección y sin
// resolución (null); después se lista y se ve.
func TestRegistrarDefecto(t *testing.T) {
	h := conDosSprints(t)

	rec := pedir(t, h, http.MethodPost, rutaDefectos, defectoJSONDePrueba)
	if rec.Code != http.StatusCreated {
		t.Fatalf("código = %d, se esperaba 201 (%s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"sprint_resolucion":null`) {
		t.Errorf("cuerpo = %s, se esperaba \"sprint_resolucion\":null", rec.Body.String())
	}
	var d defectoJSON
	leerJSON(t, rec, &d)
	esperado := defectoJSON{ID: 1, ProyectoID: 1, Descripcion: "El alta acepta emails sin @", Severidad: "Alta", Estado: "Abierto",
		HistoriaID: 1, HistoriaTitulo: "Registrar usuario", SprintDeteccion: sprintDeDefectoJSON{ID: 1, Numero: 1}}
	if d.ID != esperado.ID || d.ProyectoID != esperado.ProyectoID || d.Descripcion != esperado.Descripcion || d.Severidad != esperado.Severidad ||
		d.Estado != esperado.Estado || d.HistoriaID != esperado.HistoriaID || d.HistoriaTitulo != esperado.HistoriaTitulo ||
		d.SprintDeteccion != esperado.SprintDeteccion || d.SprintResolucion != nil {
		t.Errorf("defecto = %+v, se esperaba %+v", d, esperado)
	}

	rec = pedir(t, h, http.MethodGet, rutaDefectos, "")
	var lista []defectoJSON
	leerJSON(t, rec, &lista)
	if rec.Code != http.StatusOK || len(lista) != 1 || lista[0].ID != 1 {
		t.Errorf("GET lista: código = %d, cuerpo = %s; se esperaba el defecto 1", rec.Code, rec.Body.String())
	}
	if rec := pedir(t, h, http.MethodGet, rutaDefecto1, ""); rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"historia_titulo":"Registrar usuario"`) {
		t.Errorf("GET uno: código = %d, cuerpo = %s", rec.Code, rec.Body.String())
	}
}

// spec 11 CA-5: sin defectos, la lista es [] (no null).
func TestListarDefectosVacio(t *testing.T) {
	rec := pedir(t, conDosSprints(t), http.MethodGet, rutaDefectos, "")
	if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Errorf("código = %d, cuerpo = %s; se esperaba []", rec.Code, rec.Body.String())
	}
}

// spec 11 CA-5: datos inválidos → 400 con el mensaje de la spec.
func TestRegistrarDefectoInvalido(t *testing.T) {
	casos := []struct{ nombre, cuerpo, mensaje string }{
		{"sin descripción", `{"severidad":"Alta","historia_id":1,"sprint_deteccion_id":1}`, "La descripción es obligatoria"},
		{"descripción larga", `{"descripcion":"` + strings.Repeat("d", 1001) + `","severidad":"Alta","historia_id":1,"sprint_deteccion_id":1}`, "La descripción no puede superar los 1000 caracteres"},
		{"sin severidad", `{"descripcion":"Falla","historia_id":1,"sprint_deteccion_id":1}`, "La severidad es obligatoria"},
		{"severidad inválida", `{"descripcion":"Falla","severidad":"Critica","historia_id":1,"sprint_deteccion_id":1}`, "La severidad no es válida"},
		{"sin historia", `{"descripcion":"Falla","severidad":"Alta","sprint_deteccion_id":1}`, "La historia relacionada es obligatoria"},
		{"sin Sprint de detección", `{"descripcion":"Falla","severidad":"Alta","historia_id":1}`, "El Sprint de detección es obligatorio"},
		{"JSON mal formado", `{"descripcion":`, "El cuerpo de la solicitud no es un JSON válido"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedir(t, conDosSprints(t), http.MethodPost, rutaDefectos, c.cuerpo)
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

// spec 11 CA-2, CA-3: resolver (200, con el Sprint de resolución) y reabrir (vuelve a null).
func TestCambiarEstadoDefecto(t *testing.T) {
	h := conDosSprints(t)
	pedir(t, h, http.MethodPost, rutaDefectos, defectoJSONDePrueba)

	rec := pedir(t, h, http.MethodPut, rutaDefecto1+"/estado", `{"estado":"Resuelto","sprint_resolucion_id":2}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("resolver: código = %d, se esperaba 200 (%s)", rec.Code, rec.Body.String())
	}
	var d defectoJSON
	leerJSON(t, rec, &d)
	if d.Estado != "Resuelto" || d.SprintResolucion == nil || *d.SprintResolucion != (sprintDeDefectoJSON{ID: 2, Numero: 2}) {
		t.Errorf("resuelto = %+v; se esperaba Resuelto en el Sprint 2", d)
	}

	rec = pedir(t, h, http.MethodPut, rutaDefecto1+"/estado", `{"estado":"Abierto"}`)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"estado":"Abierto"`) || !strings.Contains(rec.Body.String(), `"sprint_resolucion":null`) {
		t.Errorf("reabrir: código = %d, cuerpo = %s; se esperaba Abierto y sin resolución", rec.Code, rec.Body.String())
	}
}

// spec 11 CA-2, CA-3: errores al cambiar el estado → 400 con el mensaje de la spec.
func TestCambiarEstadoDefectoInvalido(t *testing.T) {
	casos := []struct{ nombre, cuerpo, mensaje string }{
		{"estado inválido", `{"estado":"Cerrado"}`, "El estado no es válido"},
		{"resolver sin Sprint", `{"estado":"Resuelto"}`, "Para resolver el defecto hay que indicar el Sprint de resolución"},
		{"JSON mal formado", `{"estado":`, "El cuerpo de la solicitud no es un JSON válido"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			h := conDosSprints(t)
			pedir(t, h, http.MethodPost, rutaDefectos, defectoJSONDePrueba)
			rec := pedir(t, h, http.MethodPut, rutaDefecto1+"/estado", c.cuerpo)
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

	// Detectado en el Sprint 2, no se puede resolver en el 1.
	h := conDosSprints(t)
	pedir(t, h, http.MethodPost, rutaDefectos, `{"descripcion":"Falla","severidad":"Alta","historia_id":1,"sprint_deteccion_id":2}`)
	rec := pedir(t, h, http.MethodPut, rutaDefecto1+"/estado", `{"estado":"Resuelto","sprint_resolucion_id":1}`)
	var e errorJSON
	leerJSON(t, rec, &e)
	if rec.Code != http.StatusBadRequest || e.Error != "El Sprint de resolución no puede ser anterior al de detección" {
		t.Errorf("resolver antes de la detección: código = %d, error = %q", rec.Code, e.Error)
	}
}

// spec 11 CA-4: detectados y resueltos del proyecto y de cada Sprint.
func TestMetricasDefectos(t *testing.T) {
	h := conDosSprints(t)

	rec := pedir(t, h, http.MethodGet, "/api/proyectos/1/metricas/defectos", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"detectados":0,"resueltos":0`) {
		t.Fatalf("sin defectos: código = %d, cuerpo = %s", rec.Code, rec.Body.String())
	}

	pedir(t, h, http.MethodPost, rutaDefectos, defectoJSONDePrueba)
	pedir(t, h, http.MethodPost, rutaDefectos, `{"descripcion":"Otra","severidad":"Baja","historia_id":1,"sprint_deteccion_id":2}`)
	pedir(t, h, http.MethodPut, rutaDefecto1+"/estado", `{"estado":"Resuelto","sprint_resolucion_id":2}`)

	rec = pedir(t, h, http.MethodGet, "/api/proyectos/1/metricas/defectos", "")
	var m metricasDefectosJSON
	leerJSON(t, rec, &m)
	esperado := metricasDefectosJSON{Detectados: 2, Resueltos: 1, Sprints: []defectosDeSprintJSON{
		{SprintID: 1, Numero: 1, Detectados: 1, Resueltos: 0},
		{SprintID: 2, Numero: 2, Detectados: 1, Resueltos: 1},
	}}
	if m.Detectados != esperado.Detectados || m.Resueltos != esperado.Resueltos || len(m.Sprints) != 2 ||
		m.Sprints[0] != esperado.Sprints[0] || m.Sprints[1] != esperado.Sprints[1] {
		t.Errorf("métricas = %+v, se esperaba %+v", m, esperado)
	}
}

// spec 11 CA-6: eliminar del backlog una historia con defectos → 409.
func TestEliminarHistoriaConDefectos(t *testing.T) {
	h := conDosSprints(t)
	pedir(t, h, http.MethodPost, rutaDefectos, defectoJSONDePrueba)

	rec := pedir(t, h, http.MethodDelete, "/api/proyectos/1/backlog/1", "")
	if rec.Code != http.StatusConflict {
		t.Fatalf("código = %d, se esperaba 409 (%s)", rec.Code, rec.Body.String())
	}
	var e errorJSON
	leerJSON(t, rec, &e)
	if e.Error != "La historia tiene defectos registrados" {
		t.Errorf("error = %q", e.Error)
	}
}

// spec 11 CA-7: proyecto, historia, Sprint o defecto inexistente → 404 con el mensaje de cuál.
func TestDefectosNoExiste(t *testing.T) {
	casos := []struct{ nombre, metodo, ruta, cuerpo, mensaje string }{
		{"proyecto inexistente", http.MethodGet, "/api/proyectos/999/defectos", "", "El proyecto no existe"},
		{"proyecto inexistente en métricas", http.MethodGet, "/api/proyectos/999/metricas/defectos", "", "El proyecto no existe"},
		{"historia inexistente", http.MethodPost, rutaDefectos, `{"descripcion":"Falla","severidad":"Alta","historia_id":999,"sprint_deteccion_id":1}`, "La historia no existe"},
		{"Sprint de detección inexistente", http.MethodPost, rutaDefectos, `{"descripcion":"Falla","severidad":"Alta","historia_id":1,"sprint_deteccion_id":999}`, "El Sprint no existe"},
		{"Sprint de resolución inexistente", http.MethodPut, rutaDefecto1 + "/estado", `{"estado":"Resuelto","sprint_resolucion_id":999}`, "El Sprint no existe"},
		{"defecto inexistente", http.MethodGet, rutaDefectos + "/999", "", "El defecto no existe"},
		{"defecto con id que no es un número", http.MethodGet, rutaDefectos + "/abc", "", "El defecto no existe"},
		{"cambiar el estado de un defecto inexistente", http.MethodPut, rutaDefectos + "/999/estado", `{"estado":"Abierto"}`, "El defecto no existe"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			h := conDosSprints(t)
			pedir(t, h, http.MethodPost, rutaDefectos, defectoJSONDePrueba)
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
