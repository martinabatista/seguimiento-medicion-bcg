package api

import (
	"bytes"
	"mime"
	"net/http"
	"strings"
	"testing"
)

// conReporte devuelve un router con el proyecto 1 ("Tienda Online"), su historia 1, el Sprint 1
// (con la historia completada) y un defecto detectado en él.
func conReporte(t *testing.T) http.Handler {
	t.Helper()
	h := conDosSprints(t)
	for _, paso := range []struct{ metodo, ruta, cuerpo string }{
		{http.MethodPost, "/api/proyectos/1/sprints/1/historias", `{"historia_id":1}`},
		{http.MethodPost, "/api/proyectos/1/sprints/1/historias/1/completar", ""},
		{http.MethodPost, rutaDefectos, defectoJSONDePrueba},
	} {
		if rec := pedir(t, h, paso.metodo, paso.ruta, paso.cuerpo); rec.Code >= 300 {
			t.Fatalf("%s %s: %d %s", paso.metodo, paso.ruta, rec.Code, rec.Body.String())
		}
	}
	return h
}

// spec 14 CA-1 a CA-5: la vista previa del reporte del proyecto, con todas sus secciones.
func TestReporteProyecto(t *testing.T) {
	rec := pedir(t, conReporte(t), http.MethodGet, "/api/proyectos/1/reporte", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200 (%s)", rec.Code, rec.Body.String())
	}
	var r reporteProyectoJSON
	leerJSON(t, rec, &r)
	if r.Proyecto.Nombre != "Tienda Online" || r.Estado != "En curso" || r.GeneradoEl != "2026-11-20" {
		t.Errorf("encabezado = %+v %q %q", r.Proyecto, r.Estado, r.GeneradoEl)
	}
	if len(r.Historias) != 1 || r.Historias[0].Titulo != "Registrar usuario" || r.Historias[0].Estado != "Completada" || r.Historias[0].CompletadaEnSprint != nil {
		t.Errorf("historias = %+v", r.Historias)
	}
	if len(r.Sprints) != 2 || r.Sprints[0].Numero != 1 || r.Sprints[0].Avance.HistoriasCompletadas != 1 || !r.Sprints[0].Parcial {
		t.Errorf("Sprints = %+v", r.Sprints)
	}
	if r.Metricas.Velocidad != nil || r.Metricas.Desviacion != nil {
		t.Errorf("métricas = %+v; se esperaba sin velocidad y desviación sin datos", r.Metricas)
	}
	if r.Defectos.Detectados != 1 || len(r.Defectos.Lista) != 1 || r.Defectos.Lista[0].HistoriaTitulo != "Registrar usuario" {
		t.Errorf("defectos = %+v", r.Defectos)
	}
	for _, campo := range []string{`"integrantes":[]`, `"story_points":null`, `"horas_estimadas":null`, `"velocidad":null`} {
		if !strings.Contains(rec.Body.String(), campo) {
			t.Errorf("cuerpo sin %s: %s", campo, rec.Body.String())
		}
	}
}

// spec 14 CA-1 a CA-5: la vista previa del reporte de un Sprint marca las completadas en él.
func TestReporteSprint(t *testing.T) {
	rec := pedir(t, conReporte(t), http.MethodGet, "/api/proyectos/1/sprints/1/reporte", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200 (%s)", rec.Code, rec.Body.String())
	}
	var r reporteSprintJSON
	leerJSON(t, rec, &r)
	if r.Sprint.Numero != 1 || r.Sprint.SprintGoal != "Alta" || !r.Sprint.Parcial || len(r.Historias) != 1 ||
		r.Historias[0].CompletadaEnSprint == nil || !*r.Historias[0].CompletadaEnSprint || r.Defectos.Detectados != 1 {
		t.Errorf("reporte del Sprint = %+v", r)
	}
}

// spec 14 CA-6: los PDF se descargan con su nombre (con tildes, codificado como pide HTTP).
func TestReportePDF(t *testing.T) {
	h := conReporte(t)
	casos := []struct{ ruta, nombre string }{
		{"/api/proyectos/1/reporte/pdf", "Reporte Tienda Online.pdf"},
		{"/api/proyectos/1/sprints/1/reporte/pdf", "Reporte Tienda Online - Sprint 1.pdf"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedir(t, h, http.MethodGet, c.ruta, "")
			if rec.Code != http.StatusOK || rec.Header().Get("Content-Type") != "application/pdf" {
				t.Fatalf("código = %d, tipo = %q (%s)", rec.Code, rec.Header().Get("Content-Type"), rec.Body.String())
			}
			if !bytes.HasPrefix(rec.Body.Bytes(), []byte("%PDF-")) {
				t.Errorf("el cuerpo no es un PDF")
			}
			_, params, err := mime.ParseMediaType(rec.Header().Get("Content-Disposition"))
			if err != nil || params["filename"] != c.nombre {
				t.Errorf("Content-Disposition = %q, se esperaba el archivo %q", rec.Header().Get("Content-Disposition"), c.nombre)
			}
		})
	}

	h = conProyecto(t)
	pedir(t, h, http.MethodPost, "/api/proyectos", `{"nombre":"Gestión de Turnos","descripcion":"d","fecha_inicio":"2026-11-01","fecha_fin":"2026-12-31"}`)
	rec := pedir(t, h, http.MethodGet, "/api/proyectos/2/reporte/pdf", "")
	if _, params, _ := mime.ParseMediaType(rec.Header().Get("Content-Disposition")); params["filename"] != "Reporte Gestión de Turnos.pdf" {
		t.Errorf("Content-Disposition = %q, se esperaba el nombre con tilde", rec.Header().Get("Content-Disposition"))
	}
}

// spec 14 CA-7: proyecto o Sprint inexistente → 404 con el mensaje de cuál, también para el PDF.
func TestReporteNoExiste(t *testing.T) {
	casos := []struct{ nombre, ruta, mensaje string }{
		{"proyecto", "/api/proyectos/999/reporte", "El proyecto no existe"},
		{"proyecto en PDF", "/api/proyectos/999/reporte/pdf", "El proyecto no existe"},
		{"Sprint", "/api/proyectos/1/sprints/999/reporte", "El Sprint no existe"},
		{"Sprint en PDF", "/api/proyectos/1/sprints/999/reporte/pdf", "El Sprint no existe"},
		{"Sprint con id que no es un número", "/api/proyectos/1/sprints/abc/reporte", "El Sprint no existe"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedir(t, conProyecto(t), http.MethodGet, c.ruta, "")
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
