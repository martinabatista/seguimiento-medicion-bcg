// Package features ejecuta los escenarios BDD (*.feature) con godog como parte de go test.
// Los pasos usan la API REST con un repositorio en memoria: no hace falta una base real.
package features

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/api"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/memoria"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: inicializarEscenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"."},
			Strict:   true,
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("fallaron escenarios BDD")
	}
}

// mundo guarda el estado de un escenario entre sus pasos.
type mundo struct {
	api         http.Handler
	respuesta   *httptest.ResponseRecorder
	proyectoID  int64             // el último proyecto creado o consultado
	backlog     estadoBacklog     // pasos del backlog (spec 04), en backlog_test.go
	integrantes estadoIntegrantes // pasos de integrantes (spec 02), en integrantes_test.go
	sprints     estadoSprints     // pasos de Sprints (spec 05), en sprints_test.go
	poker       estadoPoker       // pasos de Planning Poker (spec 08), en planning_poker_test.go
	esfuerzo    estadoEsfuerzo    // pasos del registro de esfuerzo (spec 09), en esfuerzo_test.go
	defectos    estadoDefectos    // pasos de defectos (spec 11), en defectos_test.go
	hoy         time.Time         // fecha de hoy fijada por un escenario (spec 03); vacía = la del sistema
}

// fechaDeHoy es el reloj que usa la API en los escenarios: la fecha fijada o, si no, la del sistema.
func (m *mundo) fechaDeHoy() time.Time {
	if m.hoy.IsZero() {
		return time.Now()
	}
	return m.hoy
}

func inicializarEscenario(sc *godog.ScenarioContext) {
	m := &mundo{}
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		// Cada escenario arranca con una API y un repositorio vacíos.
		proyectos := memoria.NuevoRepositorioProyectos()
		historias := memoria.NuevoRepositorioHistorias()
		sprints := memoria.NuevoRepositorioSprints()
		integrantes := memoria.NuevoRepositorioIntegrantes()
		esfuerzo := memoria.NuevoRepositorioEsfuerzo()
		defectos := memoria.NuevoRepositorioDefectos()
		backlog := dominio.NuevoServicioBacklog(proyectos, historias, sprints, esfuerzo, defectos)
		estado := dominio.NuevoServicioEstado(proyectos, historias, sprints, m.fechaDeHoy)
		metricas := dominio.NuevoServicioMetricas(proyectos, historias, sprints)
		*m = mundo{api: api.NewRouter(
			dominio.NuevoServicioProyectos(proyectos),
			backlog,
			dominio.NuevoServicioIntegrantes(proyectos, integrantes),
			estado,
			dominio.NuevoServicioSprints(proyectos, historias, sprints),
			dominio.NuevoServicioPlanningPoker(proyectos, historias, integrantes, memoria.NuevoRepositorioVotaciones(), backlog),
			metricas,
			dominio.NuevoServicioEsfuerzo(proyectos, historias, integrantes, esfuerzo, m.fechaDeHoy),
			dominio.NuevoServicioDashboard(estado, metricas, m.fechaDeHoy),
			dominio.NuevoServicioMetricasEsfuerzo(proyectos, historias, sprints, esfuerzo),
			dominio.NuevoServicioDefectos(proyectos, historias, sprints, defectos),
		)}
		return ctx, nil
	})

	// Estado del servicio (spec 00)
	sc.Step(`^consulto el estado del servicio$`, m.consultoElEstadoDelServicio)
	sc.Step(`^la respuesta tiene código (\d+)$`, m.laRespuestaTieneCodigo)
	sc.Step(`^el estado informado es "([^"]*)"$`, m.elEstadoInformadoEs)

	// Proyectos (spec 01)
	sc.Step(`^creo un proyecto con nombre "([^"]*)", descripción "([^"]*)", inicio "([^"]*)" y fin "([^"]*)"$`, m.creoUnProyecto)
	sc.Step(`^creo un proyecto con un nombre de (\d+) caracteres, descripción "([^"]*)", inicio "([^"]*)" y fin "([^"]*)"$`, m.creoUnProyectoConNombreDeLargo)
	sc.Step(`^creo un proyecto con nombre "([^"]*)", una descripción de (\d+) caracteres, inicio "([^"]*)" y fin "([^"]*)"$`, m.creoUnProyectoConDescripcionDeLargo)
	sc.Step(`^que existe el proyecto "([^"]*)" con inicio "([^"]*)" y fin "([^"]*)"$`, m.queExisteElProyecto)
	sc.Step(`^cambio su nombre a "([^"]*)"$`, m.cambioSuNombre)
	sc.Step(`^cambio sus fechas a inicio "([^"]*)" y fin "([^"]*)"$`, m.cambioSusFechas)
	sc.Step(`^cambio el nombre del proyecto con identificador (\d+) a "([^"]*)"$`, m.cambioElNombreDelProyectoConID)
	sc.Step(`^listo los proyectos$`, m.listoLosProyectos)
	sc.Step(`^veo ese proyecto por su identificador$`, m.veoEseProyecto)
	sc.Step(`^veo el proyecto con identificador (\d+)$`, m.veoElProyectoConID)
	sc.Step(`^el proyecto queda registrado con un identificador$`, m.elProyectoQuedaRegistrado)
	sc.Step(`^el proyecto tiene nombre "([^"]*)", descripción "([^"]*)", inicio "([^"]*)" y fin "([^"]*)"$`, m.elProyectoTieneTodo)
	sc.Step(`^el proyecto tiene nombre "([^"]*)"$`, m.elProyectoTieneNombre)
	sc.Step(`^el proyecto tiene inicio "([^"]*)" y fin "([^"]*)"$`, m.elProyectoTieneFechas)
	sc.Step(`^se informa el error "([^"]*)"$`, m.seInformaElError)
	sc.Step(`^no se registra ningún proyecto$`, m.noSeRegistraNingunProyecto)
	sc.Step(`^hay (\d+) proyectos registrados$`, m.hayProyectosRegistrados)
	sc.Step(`^el resultado es "([^"]*)"$`, m.elResultadoEs)
	sc.Step(`^veo los proyectos "([^"]*)" y "([^"]*)" en ese orden$`, m.veoLosProyectosEnOrden)
	sc.Step(`^la lista de proyectos está vacía$`, m.laListaEstaVacia)

	// Product Backlog (spec 04)
	m.registrarPasosBacklog(sc)

	// Integrantes (spec 02)
	m.registrarPasosIntegrantes(sc)

	// Estado del proyecto (spec 03)
	m.registrarPasosEstado(sc)

	// Sprints (spec 05)
	m.registrarPasosSprints(sc)

	// Cerrar Sprint (spec 06)
	m.registrarPasosCerrarSprint(sc)

	// Métricas de avance (spec 10)
	m.registrarPasosMetricas(sc)

	// Story Points (spec 07)
	m.registrarPasosStoryPoints(sc)

	// Planning Poker (spec 08)
	m.registrarPasosPlanningPoker(sc)

	// Registro de esfuerzo (spec 09)
	m.registrarPasosEsfuerzo(sc)

	// Dashboard (spec 13)
	m.registrarPasosDashboard(sc)

	// Métricas de esfuerzo (spec 12)
	m.registrarPasosMetricasEsfuerzo(sc)

	// Defectos (spec 11)
	m.registrarPasosDefectos(sc)
}

// ---------- helpers ----------

type proyecto struct {
	ID          int64  `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
}

// texto traduce los valores especiales de las tablas de los .feature.
func texto(s string) string {
	switch s {
	case "(vacío)":
		return ""
	case "(solo espacios)":
		return "   "
	}
	return s
}

func (m *mundo) pedir(metodo, ruta string, cuerpo any) error {
	var lector *strings.Reader
	if cuerpo == nil {
		lector = strings.NewReader("")
	} else {
		b, err := json.Marshal(cuerpo)
		if err != nil {
			return err
		}
		lector = strings.NewReader(string(b))
	}
	m.respuesta = httptest.NewRecorder()
	m.api.ServeHTTP(m.respuesta, httptest.NewRequest(metodo, ruta, lector))
	return nil
}

func (m *mundo) crear(p proyecto) error {
	if err := m.pedir(http.MethodPost, "/api/proyectos", p); err != nil {
		return err
	}
	if m.respuesta.Code == http.StatusCreated {
		var creado proyecto
		if err := json.Unmarshal(m.respuesta.Body.Bytes(), &creado); err != nil {
			return err
		}
		m.proyectoID = creado.ID
	}
	return nil
}

func (m *mundo) obtener(id int64) (proyecto, error) {
	rec := httptest.NewRecorder()
	m.api.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/proyectos/%d", id), nil))
	var p proyecto
	if rec.Code != http.StatusOK {
		return p, fmt.Errorf("no se pudo ver el proyecto %d: código %d", id, rec.Code)
	}
	return p, json.Unmarshal(rec.Body.Bytes(), &p)
}

func (m *mundo) listar() ([]proyecto, error) {
	rec := httptest.NewRecorder()
	m.api.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/proyectos", nil))
	var lista []proyecto
	return lista, json.Unmarshal(rec.Body.Bytes(), &lista)
}

func (m *mundo) modificar(id int64, p proyecto) error {
	return m.pedir(http.MethodPut, fmt.Sprintf("/api/proyectos/%d", id), p)
}

// ---------- estado del servicio (spec 00) ----------

func (m *mundo) consultoElEstadoDelServicio() error {
	return m.pedir(http.MethodGet, "/api/health", nil)
}

func (m *mundo) laRespuestaTieneCodigo(codigo int) error {
	if m.respuesta.Code != codigo {
		return fmt.Errorf("código = %d, se esperaba %d", m.respuesta.Code, codigo)
	}
	return nil
}

func (m *mundo) elEstadoInformadoEs(esperado string) error {
	var body map[string]string
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &body); err != nil {
		return fmt.Errorf("la respuesta no es JSON válido: %w", err)
	}
	if body["status"] != esperado {
		return fmt.Errorf("status = %q, se esperaba %q", body["status"], esperado)
	}
	return nil
}

// ---------- proyectos (spec 01): acciones ----------

func (m *mundo) creoUnProyecto(nombre, descripcion, inicio, fin string) error {
	return m.crear(proyecto{Nombre: texto(nombre), Descripcion: texto(descripcion), FechaInicio: texto(inicio), FechaFin: texto(fin)})
}

func (m *mundo) creoUnProyectoConNombreDeLargo(largo int, descripcion, inicio, fin string) error {
	return m.crear(proyecto{Nombre: strings.Repeat("n", largo), Descripcion: descripcion, FechaInicio: inicio, FechaFin: fin})
}

func (m *mundo) creoUnProyectoConDescripcionDeLargo(nombre string, largo int, inicio, fin string) error {
	return m.crear(proyecto{Nombre: nombre, Descripcion: strings.Repeat("d", largo), FechaInicio: inicio, FechaFin: fin})
}

func (m *mundo) queExisteElProyecto(nombre, inicio, fin string) error {
	if err := m.crear(proyecto{Nombre: nombre, Descripcion: "Venta de ropa", FechaInicio: inicio, FechaFin: fin}); err != nil {
		return err
	}
	if m.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("no se pudo crear el proyecto de partida: %s", m.respuesta.Body.String())
	}
	m.recordarProyecto(nombre, m.proyectoID)
	return nil
}

func (m *mundo) cambioSuNombre(nombre string) error {
	p, err := m.obtener(m.proyectoID)
	if err != nil {
		return err
	}
	p.Nombre = nombre
	return m.modificar(m.proyectoID, p)
}

func (m *mundo) cambioSusFechas(inicio, fin string) error {
	p, err := m.obtener(m.proyectoID)
	if err != nil {
		return err
	}
	p.FechaInicio, p.FechaFin = inicio, fin
	return m.modificar(m.proyectoID, p)
}

func (m *mundo) cambioElNombreDelProyectoConID(id int64, nombre string) error {
	return m.modificar(id, proyecto{Nombre: nombre, Descripcion: "Venta de ropa", FechaInicio: "2026-11-01", FechaFin: "2026-12-15"})
}

func (m *mundo) listoLosProyectos() error {
	return m.pedir(http.MethodGet, "/api/proyectos", nil)
}

func (m *mundo) veoEseProyecto() error {
	return m.pedir(http.MethodGet, fmt.Sprintf("/api/proyectos/%d", m.proyectoID), nil)
}

func (m *mundo) veoElProyectoConID(id int64) error {
	return m.pedir(http.MethodGet, fmt.Sprintf("/api/proyectos/%d", id), nil)
}

// ---------- proyectos (spec 01): verificaciones ----------

func (m *mundo) elProyectoQuedaRegistrado() error {
	if m.respuesta.Code != http.StatusCreated || m.proyectoID == 0 {
		return fmt.Errorf("no se registró: código %d, %s", m.respuesta.Code, m.respuesta.Body.String())
	}
	return nil
}

func (m *mundo) elProyectoTieneTodo(nombre, descripcion, inicio, fin string) error {
	p, err := m.obtener(m.proyectoID)
	if err != nil {
		return err
	}
	esperado := proyecto{ID: m.proyectoID, Nombre: nombre, Descripcion: descripcion, FechaInicio: inicio, FechaFin: fin}
	if p != esperado {
		return fmt.Errorf("proyecto = %+v, se esperaba %+v", p, esperado)
	}
	return nil
}

func (m *mundo) elProyectoTieneNombre(nombre string) error {
	p, err := m.obtener(m.proyectoID)
	if err != nil {
		return err
	}
	if p.Nombre != nombre {
		return fmt.Errorf("nombre = %q, se esperaba %q", p.Nombre, nombre)
	}
	return nil
}

func (m *mundo) elProyectoTieneFechas(inicio, fin string) error {
	p, err := m.obtener(m.proyectoID)
	if err != nil {
		return err
	}
	if p.FechaInicio != inicio || p.FechaFin != fin {
		return fmt.Errorf("fechas = %s → %s, se esperaba %s → %s", p.FechaInicio, p.FechaFin, inicio, fin)
	}
	return nil
}

func (m *mundo) seInformaElError(mensaje string) error {
	var e struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &e); err != nil {
		return fmt.Errorf("la respuesta no es JSON válido: %w", err)
	}
	if e.Error != mensaje {
		return fmt.Errorf("error = %q (código %d), se esperaba %q", e.Error, m.respuesta.Code, mensaje)
	}
	return nil
}

func (m *mundo) noSeRegistraNingunProyecto() error {
	return m.hayProyectosRegistrados(0)
}

func (m *mundo) hayProyectosRegistrados(cantidad int) error {
	lista, err := m.listar()
	if err != nil {
		return err
	}
	if len(lista) != cantidad {
		return fmt.Errorf("hay %d proyectos, se esperaban %d", len(lista), cantidad)
	}
	return nil
}

// elResultadoEs acepta "registrado" (el defecto, spec 11, el registro de esfuerzo, spec 09, el Sprint, spec 05, o el
// integrante, spec 02, si el escenario pidió uno; si no, el proyecto), "registrada" o "modificada"
// (historia, spec 04), "modificado" (estado del defecto, spec 11) o el mensaje de error esperado.
func (m *mundo) elResultadoEs(resultado string) error {
	switch resultado {
	case "registrado":
		if m.defectos.registroPedido {
			return m.elDefectoQuedaRegistrado()
		}
		if m.esfuerzo.registroPedido {
			return m.elRegistroQuedaGuardado()
		}
		if m.sprints.altaPedida {
			return m.elSprintQuedaRegistrado()
		}
		if m.integrantes.altaPedida {
			return m.elIntegranteQuedaRegistrado()
		}
		return m.elProyectoQuedaRegistrado()
	case "registrada":
		return m.laHistoriaQuedaRegistrada()
	case "modificada", "modificado", "estimada", "votado":
		if m.respuesta.Code != http.StatusOK {
			return fmt.Errorf("no se modificó: código %d, %s", m.respuesta.Code, m.respuesta.Body.String())
		}
		return nil
	}
	return m.seInformaElError(resultado)
}

func (m *mundo) veoLosProyectosEnOrden(primero, segundo string) error {
	var lista []proyecto
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &lista); err != nil {
		return err
	}
	if len(lista) != 2 || lista[0].Nombre != primero || lista[1].Nombre != segundo {
		return fmt.Errorf("lista = %+v, se esperaba %q y %q en ese orden", lista, primero, segundo)
	}
	return nil
}

func (m *mundo) laListaEstaVacia() error {
	if strings.TrimSpace(m.respuesta.Body.String()) != "[]" {
		return fmt.Errorf("lista = %s, se esperaba []", m.respuesta.Body.String())
	}
	return nil
}
