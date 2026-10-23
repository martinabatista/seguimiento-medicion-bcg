package features

// Pasos de godog para los escenarios de cierre de Sprints (spec 06, cerrar-sprint.feature).

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
)

func (m *mundo) registrarPasosCerrarSprint(sc *godog.ScenarioContext) {
	// Contexto
	sc.Step(`^que existe un Sprint del "([^"]*)" al "([^"]*)" con las historias "([^"]*)" y "([^"]*)"$`, m.queExisteUnSprintConLasHistorias)
	sc.Step(`^que completé en el Sprint la historia "([^"]*)"$`, m.queCompleteEnElSprint)
	sc.Step(`^que cerré el Sprint$`, m.queCerreElSprint)

	// Acciones
	sc.Step(`^completo en el Sprint la historia "([^"]*)"$`, m.completoEnElSprint)
	sc.Step(`^completo en el Sprint la historia con identificador (\d+)$`, m.completoEnElSprintLaHistoriaConID)
	sc.Step(`^deshago en el Sprint la historia "([^"]*)"$`, m.deshagoEnElSprint)
	sc.Step(`^cierro el Sprint$`, m.cierroElSprint)
	sc.Step(`^cierro el Sprint con identificador (\d+)$`, m.cierroElSprintConID)
	sc.Step(`^consulto los Sprints anteriores$`, m.consultoLosSprintsAnteriores)

	// Verificaciones
	sc.Step(`^en el Sprint la historia "([^"]*)" está completada y "([^"]*)" no$`, m.enElSprintEstaCompletadaYLaOtraNo)
	sc.Step(`^en el Sprint la historia "([^"]*)" está completada$`, m.enElSprintEstaCompletada)
	sc.Step(`^en el Sprint ninguna historia está completada$`, m.enElSprintNingunaEstaCompletada)
	sc.Step(`^la historia "([^"]*)" del backlog está en estado "([^"]*)"$`, m.laHistoriaDelBacklogEstaEnEstado)
	sc.Step(`^el Sprint está cerrado con las historias "([^"]*)" y "([^"]*)"$`, m.elSprintEstaCerradoConLasHistorias)
	sc.Step(`^el Sprint está cerrado sin historias$`, m.elSprintEstaCerradoSinHistorias)
	sc.Step(`^el Sprint tiene la historia "([^"]*)"$`, m.elSprintTieneLaHistoria)
	sc.Step(`^veo los Sprints anteriores que empiezan el "([^"]*)" y el "([^"]*)" en ese orden$`, m.veoLosSprintsAnterioresEnOrden)
	sc.Step(`^en el Sprint anterior que empieza el "([^"]*)" la historia "([^"]*)" está completada$`, m.enElSprintAnteriorEstaCompletada)
	sc.Step(`^no hay Sprints anteriores$`, m.noHaySprintsAnteriores)
}

// ---------- helpers ----------

func (m *mundo) rutaHistoriaDelSprint(titulo string) (string, error) {
	sprintID, err := m.sprintActual()
	if err != nil {
		return "", err
	}
	historiaID, err := m.historiaLlamada(titulo)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/historias/%d/completar", rutaSprint(m.backlog.principal, sprintID), historiaID), nil
}

// historiaDelSprint busca, en el Sprint visto, la historia con ese título.
func historiaDelSprint(s sprint, titulo string) (historia, error) {
	for _, h := range s.Historias {
		if h.Titulo == titulo {
			return h, nil
		}
	}
	return historia{}, fmt.Errorf("el Sprint %d no tiene la historia %q: %+v", s.Numero, titulo, s.Historias)
}

func esperarOK(m *mundo, que string) error {
	if m.respuesta.Code != http.StatusOK {
		return fmt.Errorf("%s: código %d, %s", que, m.respuesta.Code, m.respuesta.Body.String())
	}
	return nil
}

func (m *mundo) sprintsAnteriores() ([]sprint, error) {
	var lista []sprint
	if err := esperarOK(m, "consultar Sprints anteriores"); err != nil {
		return nil, err
	}
	return lista, json.Unmarshal(m.respuesta.Body.Bytes(), &lista)
}

// ---------- contexto ----------

func (m *mundo) queExisteUnSprintConLasHistorias(inicio, fin, primera, segunda string) error {
	if err := m.queExisteUnSprint(inicio, fin); err != nil {
		return err
	}
	return m.asignoAlSprintLasHistorias(primera, segunda)
}

func (m *mundo) queCompleteEnElSprint(titulo string) error {
	if err := m.completoEnElSprint(titulo); err != nil {
		return err
	}
	return esperarOK(m, fmt.Sprintf("completar %q", titulo))
}

func (m *mundo) queCerreElSprint() error {
	if err := m.cierroElSprint(); err != nil {
		return err
	}
	return esperarOK(m, "cerrar el Sprint")
}

// ---------- acciones ----------

func (m *mundo) completoEnElSprint(titulo string) error {
	ruta, err := m.rutaHistoriaDelSprint(titulo)
	if err != nil {
		return err
	}
	return m.pedir(http.MethodPost, ruta, nil)
}

func (m *mundo) completoEnElSprintLaHistoriaConID(historiaID int64) error {
	sprintID, err := m.sprintActual()
	if err != nil {
		return err
	}
	return m.pedir(http.MethodPost, fmt.Sprintf("%s/historias/%d/completar", rutaSprint(m.backlog.principal, sprintID), historiaID), nil)
}

func (m *mundo) deshagoEnElSprint(titulo string) error {
	ruta, err := m.rutaHistoriaDelSprint(titulo)
	if err != nil {
		return err
	}
	return m.pedir(http.MethodDelete, ruta, nil)
}

func (m *mundo) cierroElSprint() error {
	sprintID, err := m.sprintActual()
	if err != nil {
		return err
	}
	return m.pedir(http.MethodPost, rutaSprint(m.backlog.principal, sprintID)+"/cerrar", nil)
}

func (m *mundo) cierroElSprintConID(sprintID int64) error {
	return m.pedir(http.MethodPost, rutaSprint(m.backlog.principal, sprintID)+"/cerrar", nil)
}

func (m *mundo) consultoLosSprintsAnteriores() error {
	return m.pedir(http.MethodGet, rutaSprints(m.backlog.principal)+"/anteriores", nil)
}

// ---------- verificaciones ----------

func (m *mundo) enElSprintEstaCompletadaYLaOtraNo(completada, otra string) error {
	if err := m.enElSprintEstaCompletada(completada); err != nil {
		return err
	}
	s, err := m.sprintVisto()
	if err != nil {
		return err
	}
	h, err := historiaDelSprint(s, otra)
	if err != nil {
		return err
	}
	if h.CompletadaEnSprint {
		return fmt.Errorf("la historia %q figura completada en el Sprint", otra)
	}
	return nil
}

func (m *mundo) enElSprintEstaCompletada(titulo string) error {
	s, err := m.sprintVisto()
	if err != nil {
		return err
	}
	h, err := historiaDelSprint(s, titulo)
	if err != nil {
		return err
	}
	if !h.CompletadaEnSprint {
		return fmt.Errorf("la historia %q no figura completada en el Sprint", titulo)
	}
	return nil
}

func (m *mundo) enElSprintNingunaEstaCompletada() error {
	s, err := m.sprintVisto()
	if err != nil {
		return err
	}
	for _, h := range s.Historias {
		if h.CompletadaEnSprint {
			return fmt.Errorf("la historia %q figura completada en el Sprint", h.Titulo)
		}
	}
	return nil
}

func (m *mundo) laHistoriaDelBacklogEstaEnEstado(titulo, estado string) error {
	id, err := m.historiaLlamada(titulo)
	if err != nil {
		return err
	}
	h, err := m.obtenerHistoria(m.backlog.principal, id)
	if err != nil {
		return err
	}
	if h.Estado != estado {
		return fmt.Errorf("estado de %q en el backlog = %q, se esperaba %q", titulo, h.Estado, estado)
	}
	return nil
}

func (m *mundo) elSprintEstaCerradoConLasHistorias(primera, segunda string) error {
	s, err := m.sprintVisto()
	if err != nil {
		return err
	}
	if s.Estado != "Cerrado" || len(s.Historias) != 2 {
		return fmt.Errorf("sprint = %+v, se esperaba cerrado con 2 historias", s)
	}
	for _, titulo := range []string{primera, segunda} {
		if _, err := historiaDelSprint(s, titulo); err != nil {
			return err
		}
	}
	return nil
}

func (m *mundo) elSprintEstaCerradoSinHistorias() error {
	s, err := m.sprintVisto()
	if err != nil {
		return err
	}
	if s.Estado != "Cerrado" || len(s.Historias) != 0 {
		return fmt.Errorf("sprint = %+v, se esperaba cerrado y sin historias", s)
	}
	return nil
}

func (m *mundo) elSprintTieneLaHistoria(titulo string) error {
	s, err := m.sprintVisto()
	if err != nil {
		return err
	}
	_, err = historiaDelSprint(s, titulo)
	return err
}

func (m *mundo) veoLosSprintsAnterioresEnOrden(primero, segundo string) error {
	lista, err := m.sprintsAnteriores()
	if err != nil {
		return err
	}
	if len(lista) != 2 || lista[0].FechaInicio != primero || lista[1].FechaInicio != segundo {
		return fmt.Errorf("anteriores = %+v, se esperaba que empiecen el %s y el %s en ese orden", lista, primero, segundo)
	}
	return nil
}

func (m *mundo) enElSprintAnteriorEstaCompletada(inicio, titulo string) error {
	lista, err := m.sprintsAnteriores()
	if err != nil {
		return err
	}
	for _, s := range lista {
		if s.FechaInicio != inicio {
			continue
		}
		h, err := historiaDelSprint(s, titulo)
		if err != nil {
			return err
		}
		if !h.CompletadaEnSprint {
			return fmt.Errorf("en el Sprint del %s la historia %q no figura completada", inicio, titulo)
		}
		return nil
	}
	return fmt.Errorf("no hay un Sprint anterior que empiece el %s", inicio)
}

func (m *mundo) noHaySprintsAnteriores() error {
	lista, err := m.sprintsAnteriores()
	if err != nil {
		return err
	}
	if len(lista) != 0 {
		return fmt.Errorf("anteriores = %+v, se esperaba ninguno", lista)
	}
	return nil
}
