package features

// Pasos de godog para los escenarios de métricas de avance (spec 10, metricas-avance.feature).

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cucumber/godog"
)

func (m *mundo) registrarPasosMetricas(sc *godog.ScenarioContext) {
	// Acciones
	sc.Step(`^consulto las métricas del proyecto$`, m.consultoLasMetricas)
	sc.Step(`^consulto las métricas del proyecto con identificador (\d+)$`, m.consultoLasMetricasDelProyectoConID)

	// Verificaciones
	sc.Step(`^el Sprint (\d+) tiene (\d+) Story Points planificados y (\d+) completados$`, m.elSprintTieneStoryPoints)
	sc.Step(`^el Sprint (\d+) tiene (\d+) historias? planificadas?, (\d+) completadas? y (\d+) % completadas$`, m.elSprintTieneHistorias)
	sc.Step(`^el Sprint (\d+) tiene (\d+) historias? sin estimar$`, m.elSprintTieneSinEstimar)
	sc.Step(`^el Sprint (\d+) está abierto y sus métricas son parciales$`, m.elSprintEsParcial)
	sc.Step(`^la velocidad del equipo es ([\d.]+) Story Points por Sprint$`, m.laVelocidadEs)
	sc.Step(`^todavía no hay velocidad del equipo$`, m.noHayVelocidad)
	sc.Step(`^el proyecto tiene (\d+) Sprints? cerrados? con (\d+) Story Points planificados y (\d+) completados$`, m.elProyectoTieneStoryPoints)
	sc.Step(`^el proyecto tiene (\d+) historias? planificadas?, (\d+) completadas? y (\d+) % completadas$`, m.elProyectoTieneHistorias)
	sc.Step(`^no hay Sprints en las métricas$`, m.noHaySprintsEnLasMetricas)
}

// ---------- helpers ----------

type metricasSprint struct {
	Numero                  int    `json:"numero"`
	Estado                  string `json:"estado"`
	Parcial                 bool   `json:"parcial"`
	HistoriasPlanificadas   int    `json:"historias_planificadas"`
	HistoriasCompletadas    int    `json:"historias_completadas"`
	PorcentajeCompletadas   int    `json:"porcentaje_completadas"`
	SinEstimar              int    `json:"sin_estimar"`
	StoryPointsPlanificados int    `json:"story_points_planificados"`
	StoryPointsCompletados  int    `json:"story_points_completados"`
}

type metricas struct {
	Sprints         []metricasSprint `json:"sprints"`
	SprintsCerrados int              `json:"sprints_cerrados"`
	Velocidad       *float64         `json:"velocidad"`
	Totales         struct {
		HistoriasPlanificadas   int `json:"historias_planificadas"`
		HistoriasCompletadas    int `json:"historias_completadas"`
		PorcentajeCompletadas   int `json:"porcentaje_completadas"`
		StoryPointsPlanificados int `json:"story_points_planificados"`
		StoryPointsCompletados  int `json:"story_points_completados"`
	} `json:"totales"`
}

// metricasConsultadas lee la respuesta de la última consulta de métricas.
func (m *mundo) metricasConsultadas() (metricas, error) {
	var r metricas
	if err := esperarOK(m, "consultar métricas"); err != nil {
		return r, err
	}
	return r, json.Unmarshal(m.respuesta.Body.Bytes(), &r)
}

func (m *mundo) metricasDelSprint(numero int) (metricasSprint, error) {
	r, err := m.metricasConsultadas()
	if err != nil {
		return metricasSprint{}, err
	}
	for _, s := range r.Sprints {
		if s.Numero == numero {
			return s, nil
		}
	}
	return metricasSprint{}, fmt.Errorf("no está el Sprint %d en las métricas: %+v", numero, r.Sprints)
}

// ---------- acciones ----------

func (m *mundo) consultoLasMetricas() error {
	return m.pedir(http.MethodGet, fmt.Sprintf("/api/proyectos/%d/metricas", m.backlog.principal), nil)
}

func (m *mundo) consultoLasMetricasDelProyectoConID(id int64) error {
	return m.pedir(http.MethodGet, fmt.Sprintf("/api/proyectos/%d/metricas", id), nil)
}

// ---------- verificaciones ----------

func (m *mundo) elSprintTieneStoryPoints(numero, planificados, completados int) error {
	s, err := m.metricasDelSprint(numero)
	if err != nil {
		return err
	}
	if s.StoryPointsPlanificados != planificados || s.StoryPointsCompletados != completados {
		return fmt.Errorf("Sprint %d: %d SP planificados y %d completados, se esperaba %d y %d",
			numero, s.StoryPointsPlanificados, s.StoryPointsCompletados, planificados, completados)
	}
	return nil
}

func (m *mundo) elSprintTieneHistorias(numero, planificadas, completadas, porcentaje int) error {
	s, err := m.metricasDelSprint(numero)
	if err != nil {
		return err
	}
	if s.HistoriasPlanificadas != planificadas || s.HistoriasCompletadas != completadas || s.PorcentajeCompletadas != porcentaje {
		return fmt.Errorf("Sprint %d: %+v, se esperaban %d planificadas, %d completadas y %d %%", numero, s, planificadas, completadas, porcentaje)
	}
	return nil
}

func (m *mundo) elSprintTieneSinEstimar(numero, sinEstimar int) error {
	s, err := m.metricasDelSprint(numero)
	if err != nil {
		return err
	}
	if s.SinEstimar != sinEstimar {
		return fmt.Errorf("Sprint %d: %d sin estimar, se esperaban %d", numero, s.SinEstimar, sinEstimar)
	}
	return nil
}

func (m *mundo) elSprintEsParcial(numero int) error {
	s, err := m.metricasDelSprint(numero)
	if err != nil {
		return err
	}
	if s.Estado != "Abierto" || !s.Parcial {
		return fmt.Errorf("Sprint %d: estado %q, parcial %v; se esperaba abierto y parcial", numero, s.Estado, s.Parcial)
	}
	return nil
}

func (m *mundo) laVelocidadEs(texto string) error {
	esperada, err := strconv.ParseFloat(texto, 64)
	if err != nil {
		return err
	}
	r, err := m.metricasConsultadas()
	if err != nil {
		return err
	}
	if r.Velocidad == nil || *r.Velocidad != esperada {
		return fmt.Errorf("velocidad = %v, se esperaba %v", r.Velocidad, esperada)
	}
	return nil
}

func (m *mundo) noHayVelocidad() error {
	r, err := m.metricasConsultadas()
	if err != nil {
		return err
	}
	if r.Velocidad != nil {
		return fmt.Errorf("velocidad = %v, se esperaba que todavía no haya", *r.Velocidad)
	}
	return nil
}

func (m *mundo) elProyectoTieneStoryPoints(cerrados, planificados, completados int) error {
	r, err := m.metricasConsultadas()
	if err != nil {
		return err
	}
	if r.SprintsCerrados != cerrados || r.Totales.StoryPointsPlanificados != planificados || r.Totales.StoryPointsCompletados != completados {
		return fmt.Errorf("proyecto: %d cerrados, %d SP planificados y %d completados; se esperaba %d, %d y %d",
			r.SprintsCerrados, r.Totales.StoryPointsPlanificados, r.Totales.StoryPointsCompletados, cerrados, planificados, completados)
	}
	return nil
}

func (m *mundo) elProyectoTieneHistorias(planificadas, completadas, porcentaje int) error {
	r, err := m.metricasConsultadas()
	if err != nil {
		return err
	}
	t := r.Totales
	if t.HistoriasPlanificadas != planificadas || t.HistoriasCompletadas != completadas || t.PorcentajeCompletadas != porcentaje {
		return fmt.Errorf("proyecto: %+v, se esperaban %d planificadas, %d completadas y %d %%", t, planificadas, completadas, porcentaje)
	}
	return nil
}

func (m *mundo) noHaySprintsEnLasMetricas() error {
	r, err := m.metricasConsultadas()
	if err != nil {
		return err
	}
	if len(r.Sprints) != 0 {
		return fmt.Errorf("sprints = %+v, se esperaba ninguno", r.Sprints)
	}
	return nil
}
