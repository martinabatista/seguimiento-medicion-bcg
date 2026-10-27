package features

// Pasos de godog para los escenarios de métricas de esfuerzo (spec 12, metricas-esfuerzo.feature).

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/cucumber/godog"
)

func (m *mundo) registrarPasosMetricasEsfuerzo(sc *godog.ScenarioContext) {
	// Antecedentes
	sc.Step(`^que estimé la historia "([^"]*)" en (\S+) horas$`, m.queEstimeLasHoras)

	// Acciones
	sc.Step(`^estimo la historia "([^"]*)" en (\S+) horas$`, m.estimoLasHoras)
	sc.Step(`^estimo la historia "([^"]*)" sin indicar las horas$`, m.estimoSinIndicarLasHoras)
	sc.Step(`^estimo la historia con identificador (\d+) en (\S+) horas$`, m.estimoLasHorasDeLaHistoriaConID)
	sc.Step(`^quito las horas estimadas de la historia "([^"]*)"$`, m.quitoLasHorasEstimadas)
	sc.Step(`^consulto las métricas de esfuerzo del proyecto$`, m.consultoLasMetricasDeEsfuerzo)
	sc.Step(`^consulto las métricas de esfuerzo del proyecto con identificador (\d+)$`, m.consultoLasMetricasDeEsfuerzoConID)

	// Verificaciones
	sc.Step(`^la historia "([^"]*)" tiene (\S+) horas estimadas$`, m.laHistoriaTieneHorasEstimadas)
	sc.Step(`^la historia "([^"]*)" no tiene horas estimadas$`, m.laHistoriaNoTieneHorasEstimadas)
	sc.Step(`^la historia "([^"]*)" tiene (\S+) horas estimadas y (\S+) horas reales$`, m.laHistoriaTieneHorasEstimadasYReales)
	sc.Step(`^la historia "([^"]*)" no tiene horas estimadas y tiene (\S+) horas reales$`, m.laHistoriaNoTieneEstimadasYTieneReales)
	sc.Step(`^la historia "([^"]*)" no tiene desviación$`, m.laHistoriaNoTieneDesviacion)
	sc.Step(`^la historia "([^"]*)" tiene una desviación de (\S+) horas y (\S+) %$`, m.laHistoriaTieneDesviacion)
	sc.Step(`^el proyecto tiene (\S+) horas estimadas y (\S+) horas reales$`, m.elProyectoTieneHoras)
	sc.Step(`^la desviación del proyecto es "sin datos"$`, m.laDesviacionDelProyectoEsSinDatos)
	sc.Step(`^el Sprint (\d+) tiene (\S+) horas estimadas, (\S+) horas reales y una desviación de (\S+) horas y (\S+) %$`, m.elSprintTieneHorasYDesviacion)
	sc.Step(`^las métricas de esfuerzo del Sprint (\d+) son parciales$`, m.lasMetricasDeEsfuerzoDelSprintSonParciales)
	sc.Step(`^el Sprint (\d+) tiene (\d+) historias? completadas? sin horas estimadas$`, m.elSprintTieneCompletadasSinEstimar)
	sc.Step(`^la desviación del Sprint (\d+) es "sin datos"$`, m.laDesviacionDelSprintEsSinDatos)
	sc.Step(`^el proyecto tiene (\d+) historias? completadas? sin horas estimadas$`, m.elProyectoTieneCompletadasSinEstimar)
}

// ---------- helpers ----------

// Las horas se comparan como texto ("6.5", "-2"): así salen de la API, exactas.
type desviacionEsfuerzo struct {
	Horas      json.Number `json:"horas"`
	Porcentaje json.Number `json:"porcentaje"`
}

type historiaConEsfuerzo struct {
	HistoriaID     int64               `json:"historia_id"`
	HorasEstimadas *json.Number        `json:"horas_estimadas"`
	HorasReales    json.Number         `json:"horas_reales"`
	Desviacion     *desviacionEsfuerzo `json:"desviacion"`
}

type sprintConEsfuerzo struct {
	Numero                       int                 `json:"numero"`
	Parcial                      bool                `json:"parcial"`
	HorasEstimadas               json.Number         `json:"horas_estimadas"`
	HorasReales                  json.Number         `json:"horas_reales"`
	Desviacion                   *desviacionEsfuerzo `json:"desviacion"`
	CompletadasSinHorasEstimadas int                 `json:"completadas_sin_horas_estimadas"`
}

type metricasDeEsfuerzo struct {
	Historias                    []historiaConEsfuerzo `json:"historias"`
	Sprints                      []sprintConEsfuerzo   `json:"sprints"`
	HorasEstimadas               json.Number           `json:"horas_estimadas"`
	HorasReales                  json.Number           `json:"horas_reales"`
	Desviacion                   *desviacionEsfuerzo   `json:"desviacion"`
	CompletadasSinHorasEstimadas int                   `json:"completadas_sin_horas_estimadas"`
}

func rutaHorasEstimadas(proyectoID, historiaID int64) string {
	return rutaHistoria(proyectoID, historiaID) + "/horas-estimadas"
}

// estimarHoras pide cargar las horas tal cual se escribieron en el escenario (pueden ser "1.125" o
// faltar), para que las valide el sistema y no el paso.
func (m *mundo) estimarHoras(historiaID int64, horas string) error {
	cuerpo := map[string]any{}
	if horas != "" {
		cuerpo["horas_estimadas"] = json.Number(horas)
	}
	return m.pedir(http.MethodPut, rutaHorasEstimadas(m.backlog.principal, historiaID), cuerpo)
}

func (m *mundo) estimarHorasPorTitulo(titulo, horas string) error {
	id, err := m.historiaLlamada(titulo)
	if err != nil {
		return err
	}
	return m.estimarHoras(id, horas)
}

// horasEstimadasDe lee la historia desde el sistema: nil = sin horas estimadas.
func (m *mundo) horasEstimadasDe(titulo string) (*json.Number, error) {
	id, err := m.historiaLlamada(titulo)
	if err != nil {
		return nil, err
	}
	rec := httptest.NewRecorder()
	m.api.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, rutaHistoria(m.backlog.principal, id), nil))
	if rec.Code != http.StatusOK {
		return nil, fmt.Errorf("no se pudo ver la historia %q: código %d", titulo, rec.Code)
	}
	var h struct {
		HorasEstimadas *json.Number `json:"horas_estimadas"`
	}
	return h.HorasEstimadas, json.Unmarshal(rec.Body.Bytes(), &h)
}

func (m *mundo) metricasDeEsfuerzo() (metricasDeEsfuerzo, error) {
	var r metricasDeEsfuerzo
	if err := esperarOK(m, "consultar métricas de esfuerzo"); err != nil {
		return r, err
	}
	return r, json.Unmarshal(m.respuesta.Body.Bytes(), &r)
}

func (m *mundo) esfuerzoDeLaHistoria(titulo string) (historiaConEsfuerzo, error) {
	id, err := m.historiaLlamada(titulo)
	if err != nil {
		return historiaConEsfuerzo{}, err
	}
	r, err := m.metricasDeEsfuerzo()
	if err != nil {
		return historiaConEsfuerzo{}, err
	}
	for _, h := range r.Historias {
		if h.HistoriaID == id {
			return h, nil
		}
	}
	return historiaConEsfuerzo{}, fmt.Errorf("no está la historia %q en las métricas de esfuerzo", titulo)
}

func (m *mundo) esfuerzoDelSprint(numero int) (sprintConEsfuerzo, error) {
	r, err := m.metricasDeEsfuerzo()
	if err != nil {
		return sprintConEsfuerzo{}, err
	}
	for _, s := range r.Sprints {
		if s.Numero == numero {
			return s, nil
		}
	}
	return sprintConEsfuerzo{}, fmt.Errorf("no está el Sprint %d en las métricas de esfuerzo: %+v", numero, r.Sprints)
}

func compararDesviacion(que string, d *desviacionEsfuerzo, horas, porcentaje string) error {
	if d == nil || d.Horas.String() != horas || d.Porcentaje.String() != porcentaje {
		return fmt.Errorf("%s: desviación = %+v, se esperaba %s horas y %s %%", que, d, horas, porcentaje)
	}
	return nil
}

// ---------- acciones ----------

func (m *mundo) queEstimeLasHoras(titulo, horas string) error {
	if err := m.estimarHorasPorTitulo(titulo, horas); err != nil {
		return err
	}
	return esperarOK(m, "estimar las horas de la historia de partida")
}

func (m *mundo) estimoLasHoras(titulo, horas string) error {
	return m.estimarHorasPorTitulo(titulo, horas)
}

func (m *mundo) estimoSinIndicarLasHoras(titulo string) error {
	return m.estimarHorasPorTitulo(titulo, "")
}

func (m *mundo) estimoLasHorasDeLaHistoriaConID(historiaID int64, horas string) error {
	return m.estimarHoras(historiaID, horas)
}

func (m *mundo) quitoLasHorasEstimadas(titulo string) error {
	id, err := m.historiaLlamada(titulo)
	if err != nil {
		return err
	}
	return m.pedir(http.MethodDelete, rutaHorasEstimadas(m.backlog.principal, id), nil)
}

func (m *mundo) consultoLasMetricasDeEsfuerzo() error {
	return m.pedir(http.MethodGet, fmt.Sprintf("/api/proyectos/%d/metricas/esfuerzo", m.backlog.principal), nil)
}

func (m *mundo) consultoLasMetricasDeEsfuerzoConID(id int64) error {
	return m.pedir(http.MethodGet, fmt.Sprintf("/api/proyectos/%d/metricas/esfuerzo", id), nil)
}

// ---------- verificaciones ----------

func (m *mundo) laHistoriaTieneHorasEstimadas(titulo, horas string) error {
	h, err := m.horasEstimadasDe(titulo)
	if err != nil {
		return err
	}
	if h == nil || h.String() != horas {
		return fmt.Errorf("horas estimadas de %q = %v, se esperaban %s", titulo, h, horas)
	}
	return nil
}

func (m *mundo) laHistoriaNoTieneHorasEstimadas(titulo string) error {
	h, err := m.horasEstimadasDe(titulo)
	if err != nil {
		return err
	}
	if h != nil {
		return fmt.Errorf("horas estimadas de %q = %s, se esperaba sin horas estimadas", titulo, *h)
	}
	return nil
}

func (m *mundo) laHistoriaTieneHorasEstimadasYReales(titulo, estimadas, reales string) error {
	h, err := m.esfuerzoDeLaHistoria(titulo)
	if err != nil {
		return err
	}
	if h.HorasEstimadas == nil || h.HorasEstimadas.String() != estimadas || h.HorasReales.String() != reales {
		return fmt.Errorf("%q: %+v, se esperaban %s horas estimadas y %s reales", titulo, h, estimadas, reales)
	}
	return nil
}

func (m *mundo) laHistoriaNoTieneEstimadasYTieneReales(titulo, reales string) error {
	h, err := m.esfuerzoDeLaHistoria(titulo)
	if err != nil {
		return err
	}
	if h.HorasEstimadas != nil || h.HorasReales.String() != reales {
		return fmt.Errorf("%q: %+v, se esperaba sin horas estimadas y %s reales", titulo, h, reales)
	}
	return nil
}

func (m *mundo) laHistoriaNoTieneDesviacion(titulo string) error {
	h, err := m.esfuerzoDeLaHistoria(titulo)
	if err != nil {
		return err
	}
	if h.Desviacion != nil {
		return fmt.Errorf("%q: desviación = %+v, se esperaba ninguna", titulo, *h.Desviacion)
	}
	return nil
}

func (m *mundo) laHistoriaTieneDesviacion(titulo, horas, porcentaje string) error {
	h, err := m.esfuerzoDeLaHistoria(titulo)
	if err != nil {
		return err
	}
	return compararDesviacion(titulo, h.Desviacion, horas, porcentaje)
}

func (m *mundo) elProyectoTieneHoras(estimadas, reales string) error {
	r, err := m.metricasDeEsfuerzo()
	if err != nil {
		return err
	}
	if r.HorasEstimadas.String() != estimadas || r.HorasReales.String() != reales {
		return fmt.Errorf("proyecto: %s horas estimadas y %s reales, se esperaban %s y %s", r.HorasEstimadas, r.HorasReales, estimadas, reales)
	}
	return nil
}

func (m *mundo) laDesviacionDelProyectoEsSinDatos() error {
	r, err := m.metricasDeEsfuerzo()
	if err != nil {
		return err
	}
	if r.Desviacion != nil {
		return fmt.Errorf("desviación del proyecto = %+v, se esperaba \"sin datos\"", *r.Desviacion)
	}
	return nil
}

func (m *mundo) elSprintTieneHorasYDesviacion(numero int, estimadas, reales, horas, porcentaje string) error {
	s, err := m.esfuerzoDelSprint(numero)
	if err != nil {
		return err
	}
	if s.HorasEstimadas.String() != estimadas || s.HorasReales.String() != reales {
		return fmt.Errorf("Sprint %d: %s horas estimadas y %s reales, se esperaban %s y %s", numero, s.HorasEstimadas, s.HorasReales, estimadas, reales)
	}
	return compararDesviacion(fmt.Sprintf("Sprint %d", numero), s.Desviacion, horas, porcentaje)
}

func (m *mundo) lasMetricasDeEsfuerzoDelSprintSonParciales(numero int) error {
	s, err := m.esfuerzoDelSprint(numero)
	if err != nil {
		return err
	}
	if !s.Parcial {
		return fmt.Errorf("Sprint %d: se esperaban métricas parciales", numero)
	}
	return nil
}

func (m *mundo) elSprintTieneCompletadasSinEstimar(numero, cantidad int) error {
	s, err := m.esfuerzoDelSprint(numero)
	if err != nil {
		return err
	}
	if s.CompletadasSinHorasEstimadas != cantidad {
		return fmt.Errorf("Sprint %d: %d completadas sin horas estimadas, se esperaban %d", numero, s.CompletadasSinHorasEstimadas, cantidad)
	}
	return nil
}

func (m *mundo) laDesviacionDelSprintEsSinDatos(numero int) error {
	s, err := m.esfuerzoDelSprint(numero)
	if err != nil {
		return err
	}
	if s.Desviacion != nil {
		return fmt.Errorf("Sprint %d: desviación = %+v, se esperaba \"sin datos\"", numero, *s.Desviacion)
	}
	return nil
}

func (m *mundo) elProyectoTieneCompletadasSinEstimar(cantidad int) error {
	r, err := m.metricasDeEsfuerzo()
	if err != nil {
		return err
	}
	if r.CompletadasSinHorasEstimadas != cantidad {
		return fmt.Errorf("proyecto: %d completadas sin horas estimadas, se esperaban %d", r.CompletadasSinHorasEstimadas, cantidad)
	}
	return nil
}
