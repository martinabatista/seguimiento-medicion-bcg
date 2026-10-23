package features

// Pasos de godog para los escenarios de Story Points (spec 07, story-points.feature).
// Reutiliza los pasos del backlog (spec 04) y de los Sprints (specs 05 y 06) para armar el estado.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

func (m *mundo) registrarPasosStoryPoints(sc *godog.ScenarioContext) {
	// Acciones
	sc.Step(`^estimo la historia "([^"]*)" con (\d+) Story Points$`, m.estimoLaHistoria)
	sc.Step(`^que estimé la historia "([^"]*)" con (\d+) Story Points$`, m.queEstimeLaHistoria)
	sc.Step(`^estimo la historia "([^"]*)" con "([^"]*)" Story Points$`, m.estimoLaHistoriaConValorEscrito)
	sc.Step(`^quito la estimación de la historia "([^"]*)"$`, m.quitoLaEstimacion)
	sc.Step(`^consulto la escala de Story Points$`, m.consultoLaEscala)
	sc.Step(`^estimo la historia con identificador (\d+) del proyecto con identificador (\d+) con (\d+) Story Points$`, m.estimoEnElProyectoConID)
	sc.Step(`^estimo la historia con identificador (\d+) con (\d+) Story Points$`, m.estimoLaHistoriaConID)
	sc.Step(`^estimo la historia "([^"]*)" desde el proyecto "([^"]*)" con (\d+) Story Points$`, m.estimoDesdeElProyecto)

	// Verificaciones
	sc.Step(`^la historia "([^"]*)" tiene (\d+) Story Points$`, m.laHistoriaTieneStoryPoints)
	sc.Step(`^la historia "([^"]*)" está sin estimar$`, m.laHistoriaEstaSinEstimar)
	sc.Step(`^la escala es ((?:\d+, )*\d+ y \d+)$`, m.laEscalaEs)
}

// ---------- helpers ----------

func rutaStoryPoints(proyectoID, historiaID int64) string {
	return rutaHistoria(proyectoID, historiaID) + "/story-points"
}

// estimar pide asignar los Story Points tal cual se escribieron en el escenario (pueden ser "2.5"
// o faltar), para que los valide el sistema y no el paso.
func (m *mundo) estimar(proyectoID, historiaID int64, valor string) error {
	cuerpo := map[string]any{}
	if valor != "" {
		cuerpo["story_points"] = json.Number(valor)
	}
	return m.pedir(http.MethodPut, rutaStoryPoints(proyectoID, historiaID), cuerpo)
}

func (m *mundo) estimarPorTitulo(titulo, valor string) error {
	id, err := m.historiaLlamada(titulo)
	if err != nil {
		return err
	}
	return m.estimar(m.backlog.principal, id, valor)
}

// storyPointsDe lee la historia desde el sistema: nil = sin estimar.
func (m *mundo) storyPointsDe(titulo string) (*int, error) {
	id, err := m.historiaLlamada(titulo)
	if err != nil {
		return nil, err
	}
	h, err := m.obtenerHistoria(m.backlog.principal, id)
	if err != nil {
		return nil, err
	}
	return h.StoryPoints, nil
}

// ---------- acciones ----------

func (m *mundo) estimoLaHistoria(titulo string, valor int) error {
	return m.estimarPorTitulo(titulo, strconv.Itoa(valor))
}

func (m *mundo) queEstimeLaHistoria(titulo string, valor int) error {
	if err := m.estimoLaHistoria(titulo, valor); err != nil {
		return err
	}
	if m.respuesta.Code != http.StatusOK {
		return fmt.Errorf("no se pudo estimar la historia de partida: %s", m.respuesta.Body.String())
	}
	return nil
}

func (m *mundo) estimoLaHistoriaConValorEscrito(titulo, valor string) error {
	return m.estimarPorTitulo(titulo, texto(valor))
}

func (m *mundo) quitoLaEstimacion(titulo string) error {
	id, err := m.historiaLlamada(titulo)
	if err != nil {
		return err
	}
	return m.pedir(http.MethodDelete, rutaStoryPoints(m.backlog.principal, id), nil)
}

func (m *mundo) consultoLaEscala() error {
	return m.pedir(http.MethodGet, "/api/story-points/escala", nil)
}

func (m *mundo) estimoEnElProyectoConID(historiaID, proyectoID int64, valor int) error {
	return m.estimar(proyectoID, historiaID, strconv.Itoa(valor))
}

func (m *mundo) estimoLaHistoriaConID(historiaID int64, valor int) error {
	return m.estimar(m.backlog.principal, historiaID, strconv.Itoa(valor))
}

func (m *mundo) estimoDesdeElProyecto(titulo, proyecto string, valor int) error {
	historiaID, err := m.historiaLlamada(titulo)
	if err != nil {
		return err
	}
	proyectoID, err := m.proyectoLlamado(proyecto)
	if err != nil {
		return err
	}
	return m.estimar(proyectoID, historiaID, strconv.Itoa(valor))
}

// ---------- verificaciones ----------

func (m *mundo) laHistoriaTieneStoryPoints(titulo string, esperado int) error {
	sp, err := m.storyPointsDe(titulo)
	if err != nil {
		return err
	}
	if sp == nil || *sp != esperado {
		return fmt.Errorf("Story Points de %q = %v, se esperaba %d", titulo, sp, esperado)
	}
	return nil
}

func (m *mundo) laHistoriaEstaSinEstimar(titulo string) error {
	sp, err := m.storyPointsDe(titulo)
	if err != nil {
		return err
	}
	if sp != nil {
		return fmt.Errorf("Story Points de %q = %d, se esperaba sin estimar", titulo, *sp)
	}
	return nil
}

// laEscalaEs compara la respuesta con una lista escrita como "1, 2, 3, 5, 8, 13 y 21".
func (m *mundo) laEscalaEs(lista string) error {
	var esperada []int
	for _, v := range strings.FieldsFunc(lista, func(r rune) bool { return r == ',' || r == ' ' || r == 'y' }) {
		n, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("valor de la escala inválido en el escenario: %q", v)
		}
		esperada = append(esperada, n)
	}
	var respuesta struct {
		Valores []int `json:"valores"`
	}
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &respuesta); err != nil {
		return fmt.Errorf("la respuesta no es JSON válido (%s): %w", m.respuesta.Body.String(), err)
	}
	if !slices.Equal(respuesta.Valores, esperada) {
		return fmt.Errorf("escala = %v, se esperaba %v", respuesta.Valores, esperada)
	}
	return nil
}
