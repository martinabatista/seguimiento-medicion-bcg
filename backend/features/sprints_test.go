package features

// Pasos de godog para los escenarios de Sprints (spec 05, sprints.feature).

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cucumber/godog"
)

// estadoSprints guarda lo que los pasos de Sprints recuerdan dentro de un escenario.
type estadoSprints struct {
	creados    []int64 // Sprints creados en el escenario, en orden
	altaPedida bool    // el último "creo un Sprint" del escenario (para "el resultado es")
}

func (m *mundo) registrarPasosSprints(sc *godog.ScenarioContext) {
	// Contexto
	sc.Step(`^que existe un Sprint del "([^"]*)" al "([^"]*)"$`, m.queExisteUnSprint)
	sc.Step(`^que existe un Sprint del "([^"]*)" al "([^"]*)" con la historia "([^"]*)"$`, m.queExisteUnSprintConLaHistoria)
	sc.Step(`^que el backlog tiene la historia "([^"]*)" con prioridad (\d+) en estado "([^"]*)"$`, m.queElBacklogTieneLaHistoriaConPrioridadEnEstado)

	// Acciones
	sc.Step(`^creo un Sprint del "([^"]*)" al "([^"]*)" con objetivo "([^"]*)"$`, m.creoUnSprint)
	sc.Step(`^creo un Sprint del "([^"]*)" al "([^"]*)" con un objetivo de (\d+) caracteres$`, m.creoUnSprintConGoalDeLargo)
	sc.Step(`^creo un Sprint en el proyecto con identificador (\d+)$`, m.creoUnSprintEnElProyectoConID)
	sc.Step(`^cambio su objetivo a "([^"]*)"$`, m.cambioSuObjetivo)
	sc.Step(`^asigno al Sprint las historias "([^"]*)" y "([^"]*)"$`, m.asignoAlSprintLasHistorias)
	sc.Step(`^asigno al Sprint la historia "([^"]*)"$`, m.asignoAlSprintLaHistoria)
	sc.Step(`^asigno al segundo Sprint la historia "([^"]*)"$`, m.asignoAlSegundoSprintLaHistoria)
	sc.Step(`^asigno al Sprint la historia con identificador (\d+)$`, m.asignoAlSprintLaHistoriaConID)
	sc.Step(`^quito del Sprint la historia "([^"]*)"$`, m.quitoDelSprintLaHistoria)
	sc.Step(`^listo los Sprints del proyecto$`, m.listoLosSprints)
	sc.Step(`^veo el Sprint con identificador (\d+)$`, m.veoElSprintConID)

	// Verificaciones
	sc.Step(`^el Sprint queda registrado como el Sprint (\d+), abierto y sin historias$`, m.elSprintQuedaRegistradoComo)
	sc.Step(`^el Sprint tiene objetivo "([^"]*)"$`, m.elSprintTieneObjetivo)
	sc.Step(`^el Sprint tiene las historias "([^"]*)" y "([^"]*)" en ese orden$`, m.elSprintTieneLasHistorias)
	sc.Step(`^el Sprint no tiene historias$`, m.elSprintNoTieneHistorias)
	sc.Step(`^veo los Sprints que empiezan el "([^"]*)" y el "([^"]*)" en ese orden$`, m.veoLosSprintsEnOrden)
	sc.Step(`^el Sprint en curso es "([^"]*)"$`, m.elSprintEnCursoEs)
}

// ---------- helpers ----------

type sprint struct {
	ID          int64      `json:"id"`
	Numero      int        `json:"numero"`
	FechaInicio string     `json:"fecha_inicio"`
	FechaFin    string     `json:"fecha_fin"`
	Goal        string     `json:"sprint_goal"`
	Estado      string     `json:"estado"`
	Historias   []historia `json:"historias"`
}

func rutaSprints(proyectoID int64) string {
	return fmt.Sprintf("/api/proyectos/%d/sprints", proyectoID)
}

func rutaSprint(proyectoID, sprintID int64) string {
	return fmt.Sprintf("/api/proyectos/%d/sprints/%d", proyectoID, sprintID)
}

// crearSprint pide el alta en el proyecto de los Antecedentes y recuerda el Sprint si se creó.
func (m *mundo) crearSprint(proyectoID int64, inicio, fin, goal string) error {
	m.sprints.altaPedida = true
	if err := m.pedir(http.MethodPost, rutaSprints(proyectoID), map[string]string{
		"fecha_inicio": texto(inicio), "fecha_fin": texto(fin), "sprint_goal": goal,
	}); err != nil {
		return err
	}
	if m.respuesta.Code == http.StatusCreated {
		var s sprint
		if err := json.Unmarshal(m.respuesta.Body.Bytes(), &s); err != nil {
			return err
		}
		m.sprints.creados = append(m.sprints.creados, s.ID)
	}
	return nil
}

// sprintActual es "el Sprint" de los pasos: el último creado en el escenario.
func (m *mundo) sprintActual() (int64, error) {
	if len(m.sprints.creados) == 0 {
		return 0, fmt.Errorf("el escenario no creó ningún Sprint")
	}
	return m.sprints.creados[len(m.sprints.creados)-1], nil
}

func (m *mundo) historiaLlamada(titulo string) (int64, error) {
	id, ok := m.backlog.historias[titulo]
	if !ok {
		return 0, fmt.Errorf("el escenario no creó la historia %q", titulo)
	}
	return id, nil
}

func (m *mundo) asignar(sprintID, historiaID int64) error {
	return m.pedir(http.MethodPost, rutaSprint(m.backlog.principal, sprintID)+"/historias", map[string]int64{"historia_id": historiaID})
}

// sprintVisto lee la última respuesta como un Sprint, o lo pide si la última no fue un Sprint.
func (m *mundo) sprintVisto() (sprint, error) {
	id, err := m.sprintActual()
	if err != nil {
		return sprint{}, err
	}
	if err := m.pedir(http.MethodGet, rutaSprint(m.backlog.principal, id), nil); err != nil {
		return sprint{}, err
	}
	var s sprint
	if m.respuesta.Code != http.StatusOK {
		return s, fmt.Errorf("ver el Sprint: código %d, %s", m.respuesta.Code, m.respuesta.Body.String())
	}
	return s, json.Unmarshal(m.respuesta.Body.Bytes(), &s)
}

// ---------- contexto ----------

func (m *mundo) queExisteUnSprint(inicio, fin string) error {
	if err := m.crearSprint(m.backlog.principal, inicio, fin, "Objetivo"); err != nil {
		return err
	}
	m.sprints.altaPedida = false
	if m.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("no se pudo crear el Sprint de partida: %s", m.respuesta.Body.String())
	}
	return nil
}

func (m *mundo) queExisteUnSprintConLaHistoria(inicio, fin, titulo string) error {
	if err := m.queExisteUnSprint(inicio, fin); err != nil {
		return err
	}
	if err := m.asignoAlSprintLaHistoria(titulo); err != nil {
		return err
	}
	if m.respuesta.Code != http.StatusOK {
		return fmt.Errorf("no se pudo asignar la historia de partida: %s", m.respuesta.Body.String())
	}
	return nil
}

func (m *mundo) queElBacklogTieneLaHistoriaConPrioridadEnEstado(titulo string, prioridad int, estado string) error {
	if err := m.queElBacklogTieneLaHistoria(titulo, prioridad); err != nil {
		return err
	}
	return m.cambiarEstadoDeHistoria(m.backlog.principal, m.backlog.historiaID, titulo, estado)
}

// ---------- acciones ----------

func (m *mundo) creoUnSprint(inicio, fin, goal string) error {
	return m.crearSprint(m.backlog.principal, inicio, fin, goal)
}

func (m *mundo) creoUnSprintConGoalDeLargo(inicio, fin string, largo int) error {
	return m.crearSprint(m.backlog.principal, inicio, fin, strings.Repeat("g", largo))
}

func (m *mundo) creoUnSprintEnElProyectoConID(proyectoID int64) error {
	return m.crearSprint(proyectoID, "2026-11-01", "2026-11-14", "Objetivo")
}

func (m *mundo) cambioSuObjetivo(goal string) error {
	id, err := m.sprintActual()
	if err != nil {
		return err
	}
	return m.pedir(http.MethodPut, rutaSprint(m.backlog.principal, id), map[string]string{"sprint_goal": goal})
}

func (m *mundo) asignoAlSprintLasHistorias(primera, segunda string) error {
	for _, titulo := range []string{primera, segunda} {
		if err := m.asignoAlSprintLaHistoria(titulo); err != nil {
			return err
		}
		if m.respuesta.Code != http.StatusOK {
			return fmt.Errorf("no se pudo asignar %q: %s", titulo, m.respuesta.Body.String())
		}
	}
	return nil
}

func (m *mundo) asignoAlSprintLaHistoria(titulo string) error {
	sprintID, err := m.sprintActual()
	if err != nil {
		return err
	}
	historiaID, err := m.historiaLlamada(titulo)
	if err != nil {
		return err
	}
	return m.asignar(sprintID, historiaID)
}

func (m *mundo) asignoAlSegundoSprintLaHistoria(titulo string) error {
	if len(m.sprints.creados) < 2 {
		return fmt.Errorf("el escenario creó %d Sprints, se necesitan 2", len(m.sprints.creados))
	}
	historiaID, err := m.historiaLlamada(titulo)
	if err != nil {
		return err
	}
	return m.asignar(m.sprints.creados[1], historiaID)
}

func (m *mundo) asignoAlSprintLaHistoriaConID(historiaID int64) error {
	sprintID, err := m.sprintActual()
	if err != nil {
		return err
	}
	return m.asignar(sprintID, historiaID)
}

func (m *mundo) quitoDelSprintLaHistoria(titulo string) error {
	sprintID, err := m.sprintActual()
	if err != nil {
		return err
	}
	historiaID, err := m.historiaLlamada(titulo)
	if err != nil {
		return err
	}
	return m.pedir(http.MethodDelete, fmt.Sprintf("%s/historias/%d", rutaSprint(m.backlog.principal, sprintID), historiaID), nil)
}

func (m *mundo) listoLosSprints() error {
	return m.pedir(http.MethodGet, rutaSprints(m.backlog.principal), nil)
}

func (m *mundo) veoElSprintConID(id int64) error {
	return m.pedir(http.MethodGet, rutaSprint(m.backlog.principal, id), nil)
}

// ---------- verificaciones ----------

func (m *mundo) elSprintQuedaRegistradoComo(numero int) error {
	if m.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("no se registró: código %d, %s", m.respuesta.Code, m.respuesta.Body.String())
	}
	var s sprint
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &s); err != nil {
		return err
	}
	if s.Numero != numero || s.Estado != "Abierto" || len(s.Historias) != 0 {
		return fmt.Errorf("sprint = %+v, se esperaba el Sprint %d abierto y sin historias", s, numero)
	}
	return nil
}

func (m *mundo) elSprintQuedaRegistrado() error {
	if m.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("no se registró: código %d, %s", m.respuesta.Code, m.respuesta.Body.String())
	}
	return nil
}

func (m *mundo) elSprintTieneObjetivo(goal string) error {
	s, err := m.sprintVisto()
	if err != nil {
		return err
	}
	if s.Goal != goal {
		return fmt.Errorf("objetivo = %q, se esperaba %q", s.Goal, goal)
	}
	return nil
}

func (m *mundo) elSprintTieneLasHistorias(primera, segunda string) error {
	s, err := m.sprintVisto()
	if err != nil {
		return err
	}
	if len(s.Historias) != 2 || s.Historias[0].Titulo != primera || s.Historias[1].Titulo != segunda {
		return fmt.Errorf("historias del Sprint = %+v, se esperaba %q y %q en ese orden", s.Historias, primera, segunda)
	}
	return nil
}

func (m *mundo) elSprintNoTieneHistorias() error {
	s, err := m.sprintVisto()
	if err != nil {
		return err
	}
	if len(s.Historias) != 0 {
		return fmt.Errorf("historias del Sprint = %+v, se esperaba ninguna", s.Historias)
	}
	return nil
}

func (m *mundo) veoLosSprintsEnOrden(primero, segundo string) error {
	var lista []sprint
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &lista); err != nil {
		return err
	}
	if len(lista) != 2 || lista[0].FechaInicio != primero || lista[1].FechaInicio != segundo {
		return fmt.Errorf("Sprints = %+v, se esperaba que empiecen el %s y el %s en ese orden", lista, primero, segundo)
	}
	return nil
}

// elSprintEnCursoEs lee el estado consultado: "Sprint N" o "ninguno".
func (m *mundo) elSprintEnCursoEs(esperado string) error {
	var e struct {
		SprintEnCurso *sprint `json:"sprint_en_curso"`
	}
	if m.respuesta.Code != http.StatusOK {
		return fmt.Errorf("código = %d, se esperaba 200 (%s)", m.respuesta.Code, m.respuesta.Body.String())
	}
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &e); err != nil {
		return err
	}
	if esperado == "ninguno" {
		if e.SprintEnCurso != nil {
			return fmt.Errorf("Sprint en curso = Sprint %d, se esperaba ninguno", e.SprintEnCurso.Numero)
		}
		return nil
	}
	if e.SprintEnCurso == nil || fmt.Sprintf("Sprint %d", e.SprintEnCurso.Numero) != esperado {
		return fmt.Errorf("Sprint en curso = %+v, se esperaba %s", e.SprintEnCurso, esperado)
	}
	return nil
}
