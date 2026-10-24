package features

// Pasos de godog para los escenarios del registro de esfuerzo (spec 09, registro-esfuerzo.feature).

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cucumber/godog"
)

// estadoEsfuerzo guarda lo que los pasos de esfuerzo recuerdan dentro de un escenario.
type estadoEsfuerzo struct {
	registroPedido bool                        // se pidió registrar esfuerzo: "registrado" se refiere a eso
	registroID     int64                       // el último registro guardado
	porActividad   map[string]registroGuardado // registros guardados, por actividad
}

type registroGuardado struct{ historiaID, id int64 }

func (m *mundo) registrarPasosEsfuerzo(sc *godog.ScenarioContext) {
	// Contexto
	sc.Step(`^que el proyecto "([^"]*)" tiene al integrante "([^"]*)" con email "([^"]*)" y rol "([^"]*)"$`, m.queElProyectoLlamadoTieneAlIntegrante)
	sc.Step(`^que la historia "([^"]*)" pasa a estado "([^"]*)"$`, m.queLaHistoriaPasaAEstado)
	sc.Step(`^que registré (\S+) horas? de "([^"]*)" el "([^"]*)" en "([^"]*)" con la actividad "([^"]*)"$`, m.queRegistreEsfuerzo)
	sc.Step(`^que eliminé el registro de la actividad "([^"]*)"$`, m.queElimineElRegistro)

	// Acciones
	sc.Step(`^registro (\S+) horas? de "([^"]*)" el "([^"]*)" en "([^"]*)" con la actividad "([^"]*)"$`, m.registroEsfuerzo)
	sc.Step(`^registro (\S+) horas? de "([^"]*)" el "([^"]*)" en "([^"]*)" con una actividad de (\d+) caracteres$`, m.registroEsfuerzoConActividadDeLargo)
	sc.Step(`^registro (\S+) horas? sin integrante el "([^"]*)" en "([^"]*)" con la actividad "([^"]*)"$`, m.registroEsfuerzoSinIntegrante)
	sc.Step(`^registro (\S+) horas? de "([^"]*)" el "([^"]*)" en la historia con identificador (\d+) con la actividad "([^"]*)"$`, m.registroEsfuerzoEnHistoriaConID)
	sc.Step(`^consulto el esfuerzo de "([^"]*)"$`, m.consultoElEsfuerzo)
	sc.Step(`^consulto el esfuerzo de la historia (\d+) del proyecto con identificador (\d+)$`, m.consultoElEsfuerzoEnProyectoConID)
	sc.Step(`^elimino (?:otra vez )?el registro de la actividad "([^"]*)"$`, m.eliminoElRegistro)

	// Verificaciones
	sc.Step(`^el registro de esfuerzo queda guardado con un identificador$`, m.elRegistroQuedaGuardado)
	sc.Step(`^el registro tiene integrante "([^"]*)", fecha "([^"]*)", actividad "([^"]*)" y (\S+) horas$`, m.elRegistroTiene)
	sc.Step(`^veo las actividades "([^"]*)", "([^"]*)" y "([^"]*)" en ese orden$`, m.veoLasActividadesEnOrden)
	sc.Step(`^veo solo la actividad "([^"]*)"$`, m.veoSoloLaActividad)
	sc.Step(`^el total de horas es (\S+)$`, m.elTotalDeHorasEs)
	sc.Step(`^no hay registros de esfuerzo$`, m.noHayRegistrosDeEsfuerzo)
	sc.Step(`^no hay esfuerzo registrado en "([^"]*)"$`, m.noHayEsfuerzoRegistradoEn)
}

// ---------- helpers ----------

type registroEsfuerzo struct {
	ID               int64       `json:"id"`
	HistoriaID       int64       `json:"historia_id"`
	IntegranteID     int64       `json:"integrante_id"`
	IntegranteNombre string      `json:"integrante_nombre"`
	Fecha            string      `json:"fecha"`
	Actividad        string      `json:"actividad"`
	Horas            json.Number `json:"horas"`
}

type esfuerzoDeHistoria struct {
	Registros  []registroEsfuerzo `json:"registros"`
	TotalHoras json.Number        `json:"total_horas"`
}

func rutaEsfuerzo(proyectoID, historiaID int64) string {
	return fmt.Sprintf("/api/proyectos/%d/backlog/%d/esfuerzo", proyectoID, historiaID)
}

// horasDelEscenario pasa "2,5" (como se lee) a "2.5" (como viaja en JSON).
func horasDelEscenario(texto string) string {
	return strings.ReplaceAll(texto, ",", ".")
}

// registrar pide registrar el esfuerzo con los datos tal cual se escribieron (integranteID 0 = sin
// integrante, horas "(vacío)" = sin horas), para que los valide el sistema y no el paso.
func (m *mundo) registrar(historiaID, integranteID int64, fecha, actividad, horas string) error {
	m.esfuerzo.registroPedido = true
	cuerpo := map[string]any{"fecha": texto(fecha), "actividad": actividad}
	if integranteID != 0 {
		cuerpo["integrante_id"] = integranteID
	}
	if h := texto(horas); h != "" {
		cuerpo["horas"] = json.Number(horasDelEscenario(h))
	}
	if err := m.pedir(http.MethodPost, rutaEsfuerzo(m.backlog.principal, historiaID), cuerpo); err != nil {
		return err
	}
	if m.respuesta.Code == http.StatusCreated {
		var r registroEsfuerzo
		if err := json.Unmarshal(m.respuesta.Body.Bytes(), &r); err != nil {
			return err
		}
		m.esfuerzo.registroID = r.ID
		if m.esfuerzo.porActividad == nil {
			m.esfuerzo.porActividad = map[string]registroGuardado{}
		}
		m.esfuerzo.porActividad[actividad] = registroGuardado{historiaID: historiaID, id: r.ID}
	}
	return nil
}

func (m *mundo) registrarPorNombres(horas, integrante, fecha, historia, actividad string) error {
	integranteID, err := m.integranteLlamado(integrante)
	if err != nil {
		return err
	}
	historiaID, err := m.historiaLlamada(historia)
	if err != nil {
		return err
	}
	return m.registrar(historiaID, integranteID, fecha, actividad, horas)
}

func (m *mundo) consultarEsfuerzo(historia string) (esfuerzoDeHistoria, error) {
	var e esfuerzoDeHistoria
	id, err := m.historiaLlamada(historia)
	if err != nil {
		return e, err
	}
	if err := m.pedir(http.MethodGet, rutaEsfuerzo(m.backlog.principal, id), nil); err != nil {
		return e, err
	}
	return m.esfuerzoConsultado()
}

func (m *mundo) esfuerzoConsultado() (esfuerzoDeHistoria, error) {
	var e esfuerzoDeHistoria
	if m.respuesta.Code != http.StatusOK {
		return e, fmt.Errorf("código = %d, se esperaba 200 (%s)", m.respuesta.Code, m.respuesta.Body.String())
	}
	return e, json.Unmarshal(m.respuesta.Body.Bytes(), &e)
}

// ---------- contexto ----------

func (m *mundo) queElProyectoLlamadoTieneAlIntegrante(proyecto, nombre, email, rol string) error {
	id, err := m.proyectoLlamado(proyecto)
	if err != nil {
		return err
	}
	if err := m.agregarIntegrante(id, nombre, email, rol); err != nil {
		return err
	}
	if m.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("no se pudo agregar el integrante de partida: %s", m.respuesta.Body.String())
	}
	return nil
}

func (m *mundo) queLaHistoriaPasaAEstado(titulo, estado string) error {
	id, err := m.historiaLlamada(titulo)
	if err != nil {
		return err
	}
	return m.cambiarEstadoDeHistoria(m.backlog.principal, id, titulo, estado)
}

func (m *mundo) queRegistreEsfuerzo(horas, integrante, fecha, historia, actividad string) error {
	if err := m.registrarPorNombres(horas, integrante, fecha, historia, actividad); err != nil {
		return err
	}
	if m.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("no se pudo registrar el esfuerzo de partida: %s", m.respuesta.Body.String())
	}
	return nil
}

func (m *mundo) queElimineElRegistro(actividad string) error {
	if err := m.eliminoElRegistro(actividad); err != nil {
		return err
	}
	if m.respuesta.Code != http.StatusNoContent {
		return fmt.Errorf("no se pudo eliminar el registro de partida: %s", m.respuesta.Body.String())
	}
	return nil
}

// ---------- acciones ----------

func (m *mundo) registroEsfuerzo(horas, integrante, fecha, historia, actividad string) error {
	return m.registrarPorNombres(horas, integrante, fecha, historia, actividad)
}

func (m *mundo) registroEsfuerzoConActividadDeLargo(horas, integrante, fecha, historia string, largo int) error {
	return m.registrarPorNombres(horas, integrante, fecha, historia, strings.Repeat("a", largo))
}

func (m *mundo) registroEsfuerzoSinIntegrante(horas, fecha, historia, actividad string) error {
	historiaID, err := m.historiaLlamada(historia)
	if err != nil {
		return err
	}
	return m.registrar(historiaID, 0, fecha, actividad, horas)
}

func (m *mundo) registroEsfuerzoEnHistoriaConID(horas, integrante, fecha string, historiaID int64, actividad string) error {
	integranteID, err := m.integranteLlamado(integrante)
	if err != nil {
		return err
	}
	return m.registrar(historiaID, integranteID, fecha, actividad, horas)
}

func (m *mundo) consultoElEsfuerzo(historia string) error {
	id, err := m.historiaLlamada(historia)
	if err != nil {
		return err
	}
	return m.pedir(http.MethodGet, rutaEsfuerzo(m.backlog.principal, id), nil)
}

func (m *mundo) consultoElEsfuerzoEnProyectoConID(historiaID, proyectoID int64) error {
	return m.pedir(http.MethodGet, rutaEsfuerzo(proyectoID, historiaID), nil)
}

func (m *mundo) eliminoElRegistro(actividad string) error {
	r, ok := m.esfuerzo.porActividad[actividad]
	if !ok {
		return fmt.Errorf("el escenario no registró la actividad %q", actividad)
	}
	return m.pedir(http.MethodDelete, fmt.Sprintf("%s/%d", rutaEsfuerzo(m.backlog.principal, r.historiaID), r.id), nil)
}

// ---------- verificaciones ----------

func (m *mundo) elRegistroQuedaGuardado() error {
	if m.respuesta.Code != http.StatusCreated || m.esfuerzo.registroID == 0 {
		return fmt.Errorf("no se registró: código %d, %s", m.respuesta.Code, m.respuesta.Body.String())
	}
	return nil
}

// elRegistroTiene verifica la respuesta del alta.
func (m *mundo) elRegistroTiene(integrante, fecha, actividad, horas string) error {
	var r registroEsfuerzo
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &r); err != nil {
		return err
	}
	if r.IntegranteNombre != integrante || r.Fecha != fecha || r.Actividad != actividad || r.Horas.String() != horasDelEscenario(horas) {
		return fmt.Errorf("registro = %+v, se esperaba %q, %s, %q y %s horas", r, integrante, fecha, actividad, horas)
	}
	return nil
}

func (m *mundo) actividadesConsultadas() ([]string, error) {
	e, err := m.esfuerzoConsultado()
	if err != nil {
		return nil, err
	}
	var actividades []string
	for _, r := range e.Registros {
		actividades = append(actividades, r.Actividad)
	}
	return actividades, nil
}

func (m *mundo) veoLasActividadesEnOrden(primera, segunda, tercera string) error {
	actividades, err := m.actividadesConsultadas()
	if err != nil {
		return err
	}
	esperado := []string{primera, segunda, tercera}
	if strings.Join(actividades, "|") != strings.Join(esperado, "|") {
		return fmt.Errorf("actividades = %v, se esperaba %v", actividades, esperado)
	}
	return nil
}

func (m *mundo) veoSoloLaActividad(actividad string) error {
	actividades, err := m.actividadesConsultadas()
	if err != nil {
		return err
	}
	if len(actividades) != 1 || actividades[0] != actividad {
		return fmt.Errorf("actividades = %v, se esperaba solo %q", actividades, actividad)
	}
	return nil
}

func (m *mundo) elTotalDeHorasEs(total string) error {
	e, err := m.esfuerzoConsultado()
	if err != nil {
		return err
	}
	if e.TotalHoras.String() != horasDelEscenario(total) {
		return fmt.Errorf("total = %s horas, se esperaba %s", e.TotalHoras, total)
	}
	return nil
}

func (m *mundo) noHayRegistrosDeEsfuerzo() error {
	e, err := m.esfuerzoConsultado()
	if err != nil {
		return err
	}
	if len(e.Registros) != 0 {
		return fmt.Errorf("registros = %+v, se esperaba ninguno", e.Registros)
	}
	return nil
}

// noHayEsfuerzoRegistradoEn vuelve a consultar la historia y espera la lista vacía.
func (m *mundo) noHayEsfuerzoRegistradoEn(historia string) error {
	e, err := m.consultarEsfuerzo(historia)
	if err != nil {
		return err
	}
	if len(e.Registros) != 0 {
		return fmt.Errorf("registros = %+v, se esperaba ninguno", e.Registros)
	}
	return nil
}
