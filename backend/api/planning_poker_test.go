package api

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// conEquipo devuelve un router con el proyecto 1, sus integrantes Ana (1) y Beto (2), Developers,
// y Caro (3), Product Owner, y la historia 1 sin estimar.
func conEquipo(t *testing.T) http.Handler {
	t.Helper()
	h := conHistoria(t)
	for _, i := range []struct{ nombre, rol string }{{"Ana", "Developer"}, {"Beto", "Developer"}, {"Caro", "Product Owner"}} {
		cuerpo := fmt.Sprintf(`{"nombre":%q,"email":"%s@mail.com","rol":%q}`, i.nombre, strings.ToLower(i.nombre), i.rol)
		if rec := pedir(t, h, http.MethodPost, "/api/proyectos/1/integrantes", cuerpo); rec.Code != http.StatusCreated {
			t.Fatalf("agregar integrante de partida: %d %s", rec.Code, rec.Body.String())
		}
	}
	return h
}

const rutaVotacion = "/api/proyectos/1/votaciones/1"

// conVotacion devuelve el router con la votación 1 abierta para la historia 1.
func conVotacion(t *testing.T) http.Handler {
	t.Helper()
	h := conEquipo(t)
	if rec := pedir(t, h, http.MethodPost, "/api/proyectos/1/votaciones", `{"historia_id":1}`); rec.Code != http.StatusCreated {
		t.Fatalf("abrir votación de partida: %d %s", rec.Code, rec.Body.String())
	}
	return h
}

func votarEn(t *testing.T, h http.Handler, integranteID, valor int) {
	t.Helper()
	cuerpo := fmt.Sprintf(`{"integrante_id":%d,"story_points":%d}`, integranteID, valor)
	if rec := pedir(t, h, http.MethodPost, rutaVotacion+"/votos", cuerpo); rec.Code != http.StatusOK {
		t.Fatalf("votar %d: %d %s", valor, rec.Code, rec.Body.String())
	}
}

func nombresJSON(is []integranteVotacionJSON) []string {
	var ns []string
	for _, i := range is {
		ns = append(ns, i.Nombre)
	}
	return ns
}

// spec 08 CA-1: abrir una votación → 201, ronda 1, faltan los Developers y todavía no hay rondas reveladas.
func TestAbrirVotacion(t *testing.T) {
	rec := pedir(t, conEquipo(t), http.MethodPost, "/api/proyectos/1/votaciones", `{"historia_id":1}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("código = %d, se esperaba 201 (%s)", rec.Code, rec.Body.String())
	}
	var v votacionJSON
	leerJSON(t, rec, &v)
	if v.ID != 1 || v.HistoriaID != 1 || v.Estado != "Abierta" || v.Ronda != 1 || v.Acordada != nil ||
		fmt.Sprint(nombresJSON(v.Faltan)) != "[Ana Beto]" || len(v.Votaron) != 0 {
		t.Errorf("votación = %+v, se esperaba abierta en la ronda 1 con Ana y Beto faltando", v)
	}
	if !strings.Contains(rec.Body.String(), `"rondas":[]`) || !strings.Contains(rec.Body.String(), `"votaron":[]`) {
		t.Errorf("cuerpo = %s, las listas vacías van como [] (no null)", rec.Body.String())
	}
}

// spec 08 CA-3, RN-5: mientras la ronda no se revela, ninguna respuesta trae el valor de un voto.
func TestVotosOcultosEnLaAPI(t *testing.T) {
	h := conVotacion(t)
	respuestas := []string{
		pedir(t, h, http.MethodPost, rutaVotacion+"/votos", `{"integrante_id":1,"story_points":13}`).Body.String(),
		pedir(t, h, http.MethodGet, rutaVotacion, "").Body.String(),
	}
	for _, cuerpo := range respuestas {
		if strings.Contains(cuerpo, `"story_points"`) || strings.Contains(cuerpo, "13") {
			t.Errorf("cuerpo = %s, no debe mostrar el valor de un voto sin revelar", cuerpo)
		}
		if !strings.Contains(cuerpo, `"votaron":[{"id":1,"nombre":"Ana"}]`) {
			t.Errorf("cuerpo = %s, se esperaba que muestre que votó Ana", cuerpo)
		}
	}
}

// spec 08 CA-4, CA-5: revelar muestra los votos y las diferencias.
func TestRevelarVotacion(t *testing.T) {
	h := conVotacion(t)
	votarEn(t, h, 1, 3)
	votarEn(t, h, 2, 13)

	rec := pedir(t, h, http.MethodPost, rutaVotacion+"/revelar", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200 (%s)", rec.Code, rec.Body.String())
	}
	var v votacionJSON
	leerJSON(t, rec, &v)
	if len(v.Rondas) != 1 {
		t.Fatalf("rondas = %+v, se esperaba la 1 revelada", v.Rondas)
	}
	r := v.Rondas[0]
	if r.Numero != 1 || len(r.Votos) != 2 || r.Votos[0].Nombre != "Ana" || r.Votos[0].StoryPoints != 3 || r.Votos[1].StoryPoints != 13 {
		t.Errorf("votos = %+v, se esperaba Ana 3 y Beto 13", r.Votos)
	}
	if !r.HayDiferencias || r.MasBajo.StoryPoints != 3 || fmt.Sprint(nombresJSON(r.MasBajo.Integrantes)) != "[Ana]" ||
		r.MasAlto.StoryPoints != 13 || fmt.Sprint(nombresJSON(r.MasAlto.Integrantes)) != "[Beto]" {
		t.Errorf("resultado = %+v, se esperaban diferencias entre 3 (Ana) y 13 (Beto)", r)
	}
}

// spec 08 CA-6, CA-7: nueva ronda, acuerdo y la historia queda estimada.
func TestNuevaRondaYAcordar(t *testing.T) {
	h := conVotacion(t)
	votarEn(t, h, 1, 3)
	votarEn(t, h, 2, 13)
	pedir(t, h, http.MethodPost, rutaVotacion+"/revelar", "")

	rec := pedir(t, h, http.MethodPost, rutaVotacion+"/rondas", "")
	var v votacionJSON
	leerJSON(t, rec, &v)
	if rec.Code != http.StatusOK || v.Ronda != 2 || len(v.Rondas) != 1 || len(v.Faltan) != 2 {
		t.Fatalf("nueva ronda: código = %d, votación = %+v", rec.Code, v)
	}
	votarEn(t, h, 1, 5)
	votarEn(t, h, 2, 5)
	pedir(t, h, http.MethodPost, rutaVotacion+"/revelar", "")

	rec = pedir(t, h, http.MethodPost, rutaVotacion+"/acordar", `{"story_points":5}`)
	leerJSON(t, rec, &v)
	if rec.Code != http.StatusOK || v.Estado != "Cerrada" || v.Acordada == nil || *v.Acordada != 5 {
		t.Fatalf("acordar: código = %d, votación = %+v", rec.Code, v)
	}
	var hist historiaJSON
	leerJSON(t, pedir(t, h, http.MethodGet, "/api/proyectos/1/backlog/1", ""), &hist)
	if hist.StoryPoints == nil || *hist.StoryPoints != 5 {
		t.Errorf("story_points de la historia = %v, se esperaba 5", hist.StoryPoints)
	}
}

// spec 08: errores con su código y su mensaje.
func TestPlanningPokerErrores(t *testing.T) {
	const fueraDeEscala = "Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21"
	casos := []struct {
		nombre, metodo, ruta, cuerpo string
		preparar                     func(*testing.T, http.Handler)
		codigo                       int
		mensaje                      string
	}{
		{"abrir: historia que no existe", http.MethodPost, "/api/proyectos/1/votaciones", `{"historia_id":999}`, nil, http.StatusNotFound, "La historia no existe"},
		{"abrir: proyecto que no existe", http.MethodPost, "/api/proyectos/999/votaciones", `{"historia_id":1}`, nil, http.StatusNotFound, "El proyecto no existe"},
		{"abrir: JSON mal formado", http.MethodPost, "/api/proyectos/1/votaciones", `{"historia_id":`, nil, http.StatusBadRequest, "El cuerpo de la solicitud no es un JSON válido"},
		{"abrir dos veces", http.MethodPost, "/api/proyectos/1/votaciones", `{"historia_id":1}`, nil, http.StatusConflict, "La historia ya tiene una votación abierta"},
		{"votar fuera de la escala", http.MethodPost, rutaVotacion + "/votos", `{"integrante_id":1,"story_points":4}`, nil, http.StatusBadRequest, fueraDeEscala},
		{"votar sin valor", http.MethodPost, rutaVotacion + "/votos", `{"integrante_id":1}`, nil, http.StatusBadRequest, "Los Story Points son obligatorios"},
		{"vota el Product Owner", http.MethodPost, rutaVotacion + "/votos", `{"integrante_id":3,"story_points":5}`, nil, http.StatusConflict, "Solo los Developers pueden votar"},
		{"vota un integrante que no existe", http.MethodPost, rutaVotacion + "/votos", `{"integrante_id":999,"story_points":5}`, nil, http.StatusNotFound, "El integrante no existe"},
		{"revelar faltando votos", http.MethodPost, rutaVotacion + "/revelar", "", func(t *testing.T, h http.Handler) { votarEn(t, h, 1, 5) }, http.StatusConflict, "Faltan votos: todavía no votaron todos los Developers"},
		{"nueva ronda sin revelar", http.MethodPost, rutaVotacion + "/rondas", "", nil, http.StatusConflict, "Primero hay que revelar la ronda actual"},
		{"acordar sin revelar", http.MethodPost, rutaVotacion + "/acordar", `{"story_points":5}`, nil, http.StatusConflict, "Primero hay que revelar la ronda actual"},
		{"votación que no existe", http.MethodGet, "/api/proyectos/1/votaciones/999", "", nil, http.StatusNotFound, "La votación no existe"},
		{"votación que no es un número", http.MethodGet, "/api/proyectos/1/votaciones/abc", "", nil, http.StatusNotFound, "La votación no existe"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			h := conVotacion(t)
			if c.preparar != nil {
				c.preparar(t, h)
			}
			rec := pedir(t, h, c.metodo, c.ruta, c.cuerpo)
			var e errorJSON
			leerJSON(t, rec, &e)
			if rec.Code != c.codigo || e.Error != c.mensaje {
				t.Errorf("código = %d, error = %q; se esperaba %d y %q", rec.Code, e.Error, c.codigo, c.mensaje)
			}
		})
	}
}

// spec 08 CA-9: las votaciones abiertas del proyecto → 200 lista, sin valores de la ronda sin revelar.
func TestVotacionesAbiertasEnLaAPI(t *testing.T) {
	if rec := pedir(t, conEquipo(t), http.MethodGet, "/api/proyectos/1/votaciones", ""); rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Errorf("sin votaciones: código = %d, cuerpo = %q; se esperaba 200 y []", rec.Code, rec.Body.String())
	}

	h := conVotacion(t)
	votarEn(t, h, 1, 13)
	rec := pedir(t, h, http.MethodGet, "/api/proyectos/1/votaciones", "")
	var lista []votacionJSON
	leerJSON(t, rec, &lista)
	if rec.Code != http.StatusOK || len(lista) != 1 || lista[0].ID != 1 || lista[0].HistoriaID != 1 || fmt.Sprint(nombresJSON(lista[0].Votaron)) != "[Ana]" {
		t.Errorf("código = %d, lista = %+v; se esperaba la votación 1 con Ana votando", rec.Code, lista)
	}
	if cuerpo := rec.Body.String(); strings.Contains(cuerpo, `"story_points"`) || strings.Contains(cuerpo, "13") {
		t.Errorf("cuerpo = %s, no debe mostrar el valor de un voto sin revelar (RN-5)", cuerpo)
	}

	rec = pedir(t, h, http.MethodGet, "/api/proyectos/999/votaciones", "")
	var e errorJSON
	leerJSON(t, rec, &e)
	if rec.Code != http.StatusNotFound || e.Error != "El proyecto no existe" {
		t.Errorf("proyecto inexistente: código = %d, error = %q", rec.Code, e.Error)
	}
}
