package features

// Pasos de godog para los escenarios del Product Backlog (spec 04, product-backlog.feature).

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

// estadoBacklog guarda lo que los pasos del backlog recuerdan dentro de un escenario.
type estadoBacklog struct {
	proyectos    map[string]int64 // proyectos creados en el escenario, por nombre
	principal    int64            // proyecto de los Antecedentes (el primero creado): "el backlog"
	historias    map[string]int64 // historias creadas, por título
	historiaID   int64            // la última historia creada
	backlogVisto int64            // proyecto del último backlog listado (si no, el principal)
}

func (m *mundo) registrarPasosBacklog(sc *godog.ScenarioContext) {
	// Acciones
	sc.Step(`^creo en el backlog la historia "([^"]*)" con descripción "([^"]*)", prioridad (\d+) y criterios "([^"]*)"$`, m.creoLaHistoriaConCriterios)
	sc.Step(`^creo en el backlog la historia "([^"]*)" con descripción "([^"]*)", prioridad (\d+) y sin criterios$`, m.creoLaHistoriaSinCriterios)
	sc.Step(`^creo en el backlog la historia "([^"]*)" con descripción "([^"]*)", prioridad "([^"]*)" y sin criterios$`, m.creoLaHistoriaConPrioridadEscrita)
	sc.Step(`^creo en el backlog una historia con un título de (\d+) caracteres, una descripción de (\d+) caracteres y criterios de (\d+) caracteres$`, m.creoUnaHistoriaConLargos)
	sc.Step(`^creo una historia en el backlog del proyecto con identificador (\d+)$`, m.creoUnaHistoriaEnElProyectoConID)
	sc.Step(`^que el backlog tiene la historia "([^"]*)" con prioridad (\d+)$`, m.queElBacklogTieneLaHistoria)
	sc.Step(`^cambio su prioridad a (\d+) y su estado a "([^"]*)"$`, m.cambioSuPrioridadYEstado)
	sc.Step(`^cambio su prioridad a (\d+)$`, m.cambioSuPrioridad)
	sc.Step(`^cambio su estado a "([^"]*)"$`, m.cambioSuEstado)
	sc.Step(`^elimino la historia "([^"]*)"$`, m.eliminoLaHistoria)
	sc.Step(`^elimino otra vez la historia "([^"]*)"$`, m.eliminoLaHistoria)
	sc.Step(`^listo el backlog$`, m.listoElBacklog)
	sc.Step(`^listo el backlog del proyecto "([^"]*)"$`, m.listoElBacklogDelProyecto)
	sc.Step(`^veo esa historia por su identificador$`, m.veoEsaHistoria)
	sc.Step(`^veo la historia con identificador (\d+)$`, m.veoLaHistoriaConID)
	sc.Step(`^veo esa historia desde el proyecto "([^"]*)"$`, m.veoEsaHistoriaDesdeElProyecto)

	// Verificaciones
	sc.Step(`^la historia queda registrada con un identificador$`, m.laHistoriaQuedaRegistrada)
	sc.Step(`^la historia tiene título "([^"]*)", descripción "([^"]*)", prioridad (\d+) y criterios "([^"]*)"$`, m.laHistoriaTieneTodo)
	sc.Step(`^la historia está en estado "([^"]*)" y sin estimar$`, m.laHistoriaEstaEnEstadoYSinEstimar)
	sc.Step(`^la historia no tiene criterios de aceptación$`, m.laHistoriaNoTieneCriterios)
	sc.Step(`^la historia tiene prioridad (\d+) y estado "([^"]*)"$`, m.laHistoriaTienePrioridadYEstado)
	sc.Step(`^la historia tiene título "([^"]*)" y prioridad (\d+)$`, m.laHistoriaTieneTituloYPrioridad)
	sc.Step(`^el backlog tiene solo la historia "([^"]*)"$`, m.elBacklogTieneSoloLaHistoria)
	sc.Step(`^veo las historias "([^"]*)", "([^"]*)" y "([^"]*)" en ese orden$`, m.veoLasHistoriasEnOrden)
	sc.Step(`^el backlog está vacío$`, m.elBacklogEstaVacio)
}

// ---------- helpers ----------

type historia struct {
	ID          int64  `json:"id"`
	ProyectoID  int64  `json:"proyecto_id"`
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	Prioridad   int    `json:"prioridad"`
	Criterios   string `json:"criterios_aceptacion"`
	Estado      string `json:"estado"`
	StoryPoints *int   `json:"story_points"`
}

// recordarProyecto anota un proyecto creado; el primero es el de los Antecedentes.
func (m *mundo) recordarProyecto(nombre string, id int64) {
	if m.backlog.proyectos == nil {
		m.backlog.proyectos = map[string]int64{}
	}
	m.backlog.proyectos[nombre] = id
	if m.backlog.principal == 0 {
		m.backlog.principal = id
	}
}

func (m *mundo) proyectoLlamado(nombre string) (int64, error) {
	id, ok := m.backlog.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el escenario no creó el proyecto %q", nombre)
	}
	return id, nil
}

func rutaBacklog(proyectoID int64) string {
	return fmt.Sprintf("/api/proyectos/%d/backlog", proyectoID)
}

func rutaHistoria(proyectoID, historiaID int64) string {
	return fmt.Sprintf("/api/proyectos/%d/backlog/%d", proyectoID, historiaID)
}

// crearHistoria pide crear la historia. La prioridad va tal cual se escribió en el escenario
// (puede ser "2.5" o faltar), para que la valide el sistema y no el paso.
func (m *mundo) crearHistoria(proyectoID int64, titulo, descripcion, prioridad, criterios string) error {
	cuerpo := map[string]any{"titulo": titulo, "descripcion": descripcion, "criterios_aceptacion": criterios}
	if prioridad != "" {
		cuerpo["prioridad"] = json.Number(prioridad)
	}
	if err := m.pedir(http.MethodPost, rutaBacklog(proyectoID), cuerpo); err != nil {
		return err
	}
	if m.respuesta.Code == http.StatusCreated {
		var creada historia
		if err := json.Unmarshal(m.respuesta.Body.Bytes(), &creada); err != nil {
			return err
		}
		if m.backlog.historias == nil {
			m.backlog.historias = map[string]int64{}
		}
		m.backlog.historias[titulo] = creada.ID
		m.backlog.historiaID = creada.ID
	}
	return nil
}

func (m *mundo) obtenerHistoria(proyectoID, historiaID int64) (historia, error) {
	rec := httptest.NewRecorder()
	m.api.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, rutaHistoria(proyectoID, historiaID), nil))
	var h historia
	if rec.Code != http.StatusOK {
		return h, fmt.Errorf("no se pudo ver la historia %d: código %d", historiaID, rec.Code)
	}
	return h, json.Unmarshal(rec.Body.Bytes(), &h)
}

func (m *mundo) listarBacklog(proyectoID int64) ([]historia, error) {
	rec := httptest.NewRecorder()
	m.api.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, rutaBacklog(proyectoID), nil))
	var lista []historia
	if rec.Code != http.StatusOK {
		return nil, fmt.Errorf("no se pudo listar el backlog: código %d", rec.Code)
	}
	return lista, json.Unmarshal(rec.Body.Bytes(), &lista)
}

// modificarHistoria lee la última historia, le aplica el cambio y la manda completa (PUT).
func (m *mundo) modificarHistoria(cambiar func(*historia)) error {
	h, err := m.obtenerHistoria(m.backlog.principal, m.backlog.historiaID)
	if err != nil {
		return err
	}
	cambiar(&h)
	return m.pedir(http.MethodPut, rutaHistoria(m.backlog.principal, h.ID), map[string]any{
		"titulo": h.Titulo, "descripcion": h.Descripcion, "prioridad": h.Prioridad,
		"criterios_aceptacion": h.Criterios, "estado": h.Estado,
	})
}

// ---------- acciones ----------

func (m *mundo) creoLaHistoriaConCriterios(titulo, descripcion string, prioridad int, criterios string) error {
	return m.crearHistoria(m.backlog.principal, texto(titulo), descripcion, strconv.Itoa(prioridad), criterios)
}

func (m *mundo) creoLaHistoriaSinCriterios(titulo, descripcion string, prioridad int) error {
	return m.crearHistoria(m.backlog.principal, texto(titulo), descripcion, strconv.Itoa(prioridad), "")
}

func (m *mundo) creoLaHistoriaConPrioridadEscrita(titulo, descripcion, prioridad string) error {
	return m.crearHistoria(m.backlog.principal, titulo, descripcion, texto(prioridad), "")
}

func (m *mundo) creoUnaHistoriaConLargos(largoTitulo, largoDescripcion, largoCriterios int) error {
	return m.crearHistoria(m.backlog.principal, strings.Repeat("t", largoTitulo),
		strings.Repeat("d", largoDescripcion), "1", strings.Repeat("c", largoCriterios))
}

func (m *mundo) creoUnaHistoriaEnElProyectoConID(proyectoID int64) error {
	return m.crearHistoria(proyectoID, "Registrar usuario", "Alta de clientes", "1", "")
}

func (m *mundo) queElBacklogTieneLaHistoria(titulo string, prioridad int) error {
	if err := m.crearHistoria(m.backlog.principal, titulo, "Alta de clientes", strconv.Itoa(prioridad), ""); err != nil {
		return err
	}
	if m.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("no se pudo crear la historia de partida: %s", m.respuesta.Body.String())
	}
	return nil
}

func (m *mundo) cambioSuPrioridadYEstado(prioridad int, estado string) error {
	return m.modificarHistoria(func(h *historia) { h.Prioridad, h.Estado = prioridad, estado })
}

func (m *mundo) cambioSuPrioridad(prioridad int) error {
	return m.modificarHistoria(func(h *historia) { h.Prioridad = prioridad })
}

func (m *mundo) cambioSuEstado(estado string) error {
	return m.modificarHistoria(func(h *historia) { h.Estado = estado })
}

func (m *mundo) eliminoLaHistoria(titulo string) error {
	id, ok := m.backlog.historias[titulo]
	if !ok {
		return fmt.Errorf("el escenario no creó la historia %q", titulo)
	}
	return m.pedir(http.MethodDelete, rutaHistoria(m.backlog.principal, id), nil)
}

func (m *mundo) listoElBacklog() error {
	m.backlog.backlogVisto = m.backlog.principal
	return m.pedir(http.MethodGet, rutaBacklog(m.backlog.principal), nil)
}

func (m *mundo) listoElBacklogDelProyecto(nombre string) error {
	id, err := m.proyectoLlamado(nombre)
	if err != nil {
		return err
	}
	m.backlog.backlogVisto = id
	return m.pedir(http.MethodGet, rutaBacklog(id), nil)
}

func (m *mundo) veoEsaHistoria() error {
	return m.pedir(http.MethodGet, rutaHistoria(m.backlog.principal, m.backlog.historiaID), nil)
}

func (m *mundo) veoLaHistoriaConID(id int64) error {
	return m.pedir(http.MethodGet, rutaHistoria(m.backlog.principal, id), nil)
}

func (m *mundo) veoEsaHistoriaDesdeElProyecto(nombre string) error {
	id, err := m.proyectoLlamado(nombre)
	if err != nil {
		return err
	}
	return m.pedir(http.MethodGet, rutaHistoria(id, m.backlog.historiaID), nil)
}

// ---------- verificaciones ----------

func (m *mundo) laHistoriaQuedaRegistrada() error {
	if m.respuesta.Code != http.StatusCreated || m.backlog.historiaID == 0 {
		return fmt.Errorf("no se registró: código %d, %s", m.respuesta.Code, m.respuesta.Body.String())
	}
	return nil
}

func (m *mundo) laHistoriaTieneTodo(titulo, descripcion string, prioridad int, criterios string) error {
	h, err := m.obtenerHistoria(m.backlog.principal, m.backlog.historiaID)
	if err != nil {
		return err
	}
	if h.Titulo != titulo || h.Descripcion != descripcion || h.Prioridad != prioridad || h.Criterios != criterios {
		return fmt.Errorf("historia = %+v, se esperaba título %q, descripción %q, prioridad %d y criterios %q",
			h, titulo, descripcion, prioridad, criterios)
	}
	return nil
}

func (m *mundo) laHistoriaEstaEnEstadoYSinEstimar(estado string) error {
	h, err := m.obtenerHistoria(m.backlog.principal, m.backlog.historiaID)
	if err != nil {
		return err
	}
	if h.Estado != estado || h.StoryPoints != nil {
		return fmt.Errorf("estado = %q, Story Points = %v; se esperaba %q y sin estimar", h.Estado, h.StoryPoints, estado)
	}
	return nil
}

func (m *mundo) laHistoriaNoTieneCriterios() error {
	h, err := m.obtenerHistoria(m.backlog.principal, m.backlog.historiaID)
	if err != nil {
		return err
	}
	if h.Criterios != "" {
		return fmt.Errorf("criterios = %q, se esperaba sin criterios", h.Criterios)
	}
	return nil
}

func (m *mundo) laHistoriaTienePrioridadYEstado(prioridad int, estado string) error {
	h, err := m.obtenerHistoria(m.backlog.principal, m.backlog.historiaID)
	if err != nil {
		return err
	}
	if h.Prioridad != prioridad || h.Estado != estado {
		return fmt.Errorf("prioridad = %d, estado = %q; se esperaba %d y %q", h.Prioridad, h.Estado, prioridad, estado)
	}
	return nil
}

// laHistoriaTieneTituloYPrioridad verifica la respuesta de "veo esa historia".
func (m *mundo) laHistoriaTieneTituloYPrioridad(titulo string, prioridad int) error {
	if m.respuesta.Code != http.StatusOK {
		return fmt.Errorf("código = %d, se esperaba 200 (%s)", m.respuesta.Code, m.respuesta.Body.String())
	}
	var h historia
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &h); err != nil {
		return err
	}
	if h.Titulo != titulo || h.Prioridad != prioridad {
		return fmt.Errorf("historia = %+v, se esperaba título %q y prioridad %d", h, titulo, prioridad)
	}
	return nil
}

func (m *mundo) elBacklogTieneSoloLaHistoria(titulo string) error {
	lista, err := m.listarBacklog(m.backlog.principal)
	if err != nil {
		return err
	}
	if len(lista) != 1 || lista[0].Titulo != titulo {
		return fmt.Errorf("backlog = %+v, se esperaba solo %q", lista, titulo)
	}
	return nil
}

func (m *mundo) veoLasHistoriasEnOrden(primera, segunda, tercera string) error {
	var lista []historia
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &lista); err != nil {
		return err
	}
	if len(lista) != 3 || lista[0].Titulo != primera || lista[1].Titulo != segunda || lista[2].Titulo != tercera {
		return fmt.Errorf("backlog = %+v, se esperaba %q, %q y %q en ese orden", lista, primera, segunda, tercera)
	}
	return nil
}

// elBacklogEstaVacio vuelve a listar el último backlog visto (o el principal) y espera [].
func (m *mundo) elBacklogEstaVacio() error {
	id := m.backlog.backlogVisto
	if id == 0 {
		id = m.backlog.principal
	}
	lista, err := m.listarBacklog(id)
	if err != nil {
		return err
	}
	if len(lista) != 0 {
		return fmt.Errorf("backlog = %+v, se esperaba vacío", lista)
	}
	return nil
}
