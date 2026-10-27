package features

// Pasos de godog para los escenarios de defectos (spec 11, defectos.feature).

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cucumber/godog"
)

// estadoDefectos guarda lo que los pasos de defectos recuerdan dentro de un escenario.
type estadoDefectos struct {
	registroPedido bool             // se pidió registrar un defecto: "registrado" se refiere a eso
	ultimo         int64            // el último defecto registrado o cambiado: "el defecto" de los pasos
	porDescripcion map[string]int64 // defectos registrados, por descripción
	sprintsAjenos  map[string]int64 // el Sprint de cada proyecto que no es el de los Antecedentes
}

func (m *mundo) registrarPasosDefectos(sc *godog.ScenarioContext) {
	// Contexto
	sc.Step(`^que cerré el Sprint (\d+)$`, m.queCerreElSprintNumero)
	sc.Step(`^que el proyecto "([^"]*)" tiene un Sprint del "([^"]*)" al "([^"]*)"$`, m.queElProyectoTieneUnSprint)
	sc.Step(`^que registré el defecto "([^"]*)" de severidad "([^"]*)" en "([^"]*)" detectado en el Sprint (\d+)$`, m.queRegistreElDefecto)
	sc.Step(`^que resolví el defecto "([^"]*)" en el Sprint (\d+)$`, m.queResolviElDefecto)
	sc.Step(`^que cambié el estado del defecto "([^"]*)" a "([^"]*)"$`, m.queCambieElEstadoDelDefecto)

	// Acciones
	sc.Step(`^registro el defecto "([^"]*)" de severidad "([^"]*)" en "([^"]*)" detectado en el Sprint (\d+)$`, m.registroElDefecto)
	sc.Step(`^registro un defecto con una descripción de (\d+) caracteres de severidad "([^"]*)" en "([^"]*)" detectado en el Sprint (\d+)$`, m.registroUnDefectoConDescripcionDeLargo)
	sc.Step(`^registro el defecto "([^"]*)" de severidad "([^"]*)" sin historia detectado en el Sprint (\d+)$`, m.registroElDefectoSinHistoria)
	sc.Step(`^registro el defecto "([^"]*)" de severidad "([^"]*)" en "([^"]*)" sin Sprint de detección$`, m.registroElDefectoSinSprint)
	sc.Step(`^registro el defecto "([^"]*)" de severidad "([^"]*)" en la historia con identificador (\d+) detectado en el Sprint (\d+)$`, m.registroElDefectoEnHistoriaConID)
	sc.Step(`^registro el defecto "([^"]*)" de severidad "([^"]*)" en "([^"]*)" detectado en el Sprint con identificador (\d+)$`, m.registroElDefectoEnSprintConID)
	sc.Step(`^cambio el estado del defecto "([^"]*)" a "([^"]*)"$`, m.cambioElEstadoDelDefecto)
	sc.Step(`^resuelvo el defecto "([^"]*)" en el Sprint (\d+)$`, m.resuelvoElDefecto)
	sc.Step(`^resuelvo el defecto "([^"]*)" en el Sprint de "([^"]*)"$`, m.resuelvoElDefectoEnElSprintDe)
	sc.Step(`^consulto las métricas de defectos$`, m.consultoLasMetricasDeDefectos)
	sc.Step(`^listo los defectos$`, m.listoLosDefectos)
	sc.Step(`^veo el defecto con identificador (\d+)$`, m.veoElDefectoConID)

	// Verificaciones
	sc.Step(`^el defecto queda registrado con un identificador$`, m.elDefectoQuedaRegistrado)
	sc.Step(`^el defecto tiene descripción "([^"]*)", severidad "([^"]*)", historia "([^"]*)" y Sprint de detección (\d+)$`, m.elDefectoTiene)
	sc.Step(`^el defecto está "([^"]*)" y sin Sprint de resolución$`, m.elDefectoEstaSinResolucion)
	sc.Step(`^el defecto está "([^"]*)" con Sprint de resolución (\d+)$`, m.elDefectoEstaConResolucion)
	sc.Step(`^el proyecto tiene (\d+) defectos? detectados? y (\d+) resueltos$`, m.elProyectoTieneDefectos)
	sc.Step(`^el Sprint (\d+) tiene (\d+) defectos? detectados? y (\d+) resueltos$`, m.elSprintTieneDefectos)
	sc.Step(`^veo los defectos "([^"]*)" y "([^"]*)" en ese orden$`, m.veoLosDefectosEnOrden)
	sc.Step(`^no hay defectos registrados$`, m.noHayDefectosRegistrados)
}

// ---------- helpers ----------

type sprintDeDefecto struct {
	ID     int64 `json:"id"`
	Numero int   `json:"numero"`
}

type defecto struct {
	ID               int64            `json:"id"`
	Descripcion      string           `json:"descripcion"`
	Severidad        string           `json:"severidad"`
	Estado           string           `json:"estado"`
	HistoriaID       int64            `json:"historia_id"`
	HistoriaTitulo   string           `json:"historia_titulo"`
	SprintDeteccion  sprintDeDefecto  `json:"sprint_deteccion"`
	SprintResolucion *sprintDeDefecto `json:"sprint_resolucion"`
}

type metricasDefectos struct {
	Detectados int `json:"detectados"`
	Resueltos  int `json:"resueltos"`
	Sprints    []struct {
		Numero     int `json:"numero"`
		Detectados int `json:"detectados"`
		Resueltos  int `json:"resueltos"`
	} `json:"sprints"`
}

func rutaDefectos(proyectoID int64) string {
	return fmt.Sprintf("/api/proyectos/%d/defectos", proyectoID)
}

// sprintNumero es el identificador del Sprint N de los Antecedentes (se crean en orden).
func (m *mundo) sprintNumero(numero int) (int64, error) {
	if numero < 1 || numero > len(m.sprints.creados) {
		return 0, fmt.Errorf("el escenario no creó el Sprint %d", numero)
	}
	return m.sprints.creados[numero-1], nil
}

func (m *mundo) defectoLlamado(descripcion string) (int64, error) {
	id, ok := m.defectos.porDescripcion[descripcion]
	if !ok {
		return 0, fmt.Errorf("el escenario no registró el defecto %q", descripcion)
	}
	return id, nil
}

// registrarDefecto pide el registro con los datos tal cual se escribieron (identificadores en 0 =
// sin indicar), para que los valide el sistema y no el paso.
func (m *mundo) registrarDefecto(descripcion, severidad string, historiaID, sprintID int64) error {
	m.defectos.registroPedido = true
	cuerpo := map[string]any{"descripcion": texto(descripcion), "severidad": texto(severidad)}
	if historiaID != 0 {
		cuerpo["historia_id"] = historiaID
	}
	if sprintID != 0 {
		cuerpo["sprint_deteccion_id"] = sprintID
	}
	if err := m.pedir(http.MethodPost, rutaDefectos(m.backlog.principal), cuerpo); err != nil {
		return err
	}
	if m.respuesta.Code == http.StatusCreated {
		var d defecto
		if err := json.Unmarshal(m.respuesta.Body.Bytes(), &d); err != nil {
			return err
		}
		m.defectos.ultimo = d.ID
		if m.defectos.porDescripcion == nil {
			m.defectos.porDescripcion = map[string]int64{}
		}
		m.defectos.porDescripcion[descripcion] = d.ID
	}
	return nil
}

func (m *mundo) registrarDefectoPorNombres(descripcion, severidad, historia string, sprint int) error {
	historiaID, err := m.historiaLlamada(historia)
	if err != nil {
		return err
	}
	sprintID, err := m.sprintNumero(sprint)
	if err != nil {
		return err
	}
	return m.registrarDefecto(descripcion, severidad, historiaID, sprintID)
}

func (m *mundo) cambiarEstadoDefecto(descripcion, estado string, sprintID int64) error {
	id, err := m.defectoLlamado(descripcion)
	if err != nil {
		return err
	}
	m.defectos.ultimo = id
	cuerpo := map[string]any{"estado": estado}
	if sprintID != 0 {
		cuerpo["sprint_resolucion_id"] = sprintID
	}
	return m.pedir(http.MethodPut, fmt.Sprintf("%s/%d/estado", rutaDefectos(m.backlog.principal), id), cuerpo)
}

// defectoActual pide "el defecto" de los pasos, para verificar lo que quedó guardado.
func (m *mundo) defectoActual() (defecto, error) {
	var d defecto
	if m.defectos.ultimo == 0 {
		return d, fmt.Errorf("el escenario no registró ningún defecto")
	}
	if err := m.pedir(http.MethodGet, fmt.Sprintf("%s/%d", rutaDefectos(m.backlog.principal), m.defectos.ultimo), nil); err != nil {
		return d, err
	}
	if m.respuesta.Code != http.StatusOK {
		return d, fmt.Errorf("no se pudo ver el defecto: código %d, %s", m.respuesta.Code, m.respuesta.Body.String())
	}
	return d, json.Unmarshal(m.respuesta.Body.Bytes(), &d)
}

func (m *mundo) metricasDeDefectosVistas() (metricasDefectos, error) {
	var md metricasDefectos
	if m.respuesta.Code != http.StatusOK {
		return md, fmt.Errorf("no se pudieron consultar las métricas: código %d, %s", m.respuesta.Code, m.respuesta.Body.String())
	}
	return md, json.Unmarshal(m.respuesta.Body.Bytes(), &md)
}

func (m *mundo) listaDeDefectos() ([]defecto, error) {
	if err := m.pedir(http.MethodGet, rutaDefectos(m.backlog.principal), nil); err != nil {
		return nil, err
	}
	var lista []defecto
	return lista, json.Unmarshal(m.respuesta.Body.Bytes(), &lista)
}

// ---------- contexto ----------

func (m *mundo) queCerreElSprintNumero(numero int) error {
	id, err := m.sprintNumero(numero)
	if err != nil {
		return err
	}
	if err := m.cierroElSprintConID(id); err != nil {
		return err
	}
	return esperarOK(m, "cerrar el Sprint")
}

func (m *mundo) queElProyectoTieneUnSprint(nombre, inicio, fin string) error {
	proyectoID, err := m.proyectoLlamado(nombre)
	if err != nil {
		return err
	}
	if err := m.pedir(http.MethodPost, rutaSprints(proyectoID), map[string]string{
		"fecha_inicio": inicio, "fecha_fin": fin, "sprint_goal": "Objetivo",
	}); err != nil {
		return err
	}
	if m.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("no se pudo crear el Sprint de partida: %s", m.respuesta.Body.String())
	}
	var s sprint
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &s); err != nil {
		return err
	}
	if m.defectos.sprintsAjenos == nil {
		m.defectos.sprintsAjenos = map[string]int64{}
	}
	m.defectos.sprintsAjenos[nombre] = s.ID
	return nil
}

func (m *mundo) queRegistreElDefecto(descripcion, severidad, historia string, sprint int) error {
	if err := m.registrarDefectoPorNombres(descripcion, severidad, historia, sprint); err != nil {
		return err
	}
	m.defectos.registroPedido = false
	if m.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("no se pudo registrar el defecto de partida: %s", m.respuesta.Body.String())
	}
	return nil
}

func (m *mundo) queResolviElDefecto(descripcion string, sprint int) error {
	if err := m.resuelvoElDefecto(descripcion, sprint); err != nil {
		return err
	}
	return esperarOK(m, "resolver el defecto")
}

func (m *mundo) queCambieElEstadoDelDefecto(descripcion, estado string) error {
	if err := m.cambioElEstadoDelDefecto(descripcion, estado); err != nil {
		return err
	}
	return esperarOK(m, "cambiar el estado del defecto")
}

// ---------- acciones ----------

func (m *mundo) registroElDefecto(descripcion, severidad, historia string, sprint int) error {
	return m.registrarDefectoPorNombres(descripcion, severidad, historia, sprint)
}

func (m *mundo) registroUnDefectoConDescripcionDeLargo(largo int, severidad, historia string, sprint int) error {
	return m.registrarDefectoPorNombres(strings.Repeat("d", largo), severidad, historia, sprint)
}

func (m *mundo) registroElDefectoSinHistoria(descripcion, severidad string, sprint int) error {
	sprintID, err := m.sprintNumero(sprint)
	if err != nil {
		return err
	}
	return m.registrarDefecto(descripcion, severidad, 0, sprintID)
}

func (m *mundo) registroElDefectoSinSprint(descripcion, severidad, historia string) error {
	historiaID, err := m.historiaLlamada(historia)
	if err != nil {
		return err
	}
	return m.registrarDefecto(descripcion, severidad, historiaID, 0)
}

func (m *mundo) registroElDefectoEnHistoriaConID(descripcion, severidad string, historiaID int64, sprint int) error {
	sprintID, err := m.sprintNumero(sprint)
	if err != nil {
		return err
	}
	return m.registrarDefecto(descripcion, severidad, historiaID, sprintID)
}

func (m *mundo) registroElDefectoEnSprintConID(descripcion, severidad, historia string, sprintID int64) error {
	historiaID, err := m.historiaLlamada(historia)
	if err != nil {
		return err
	}
	return m.registrarDefecto(descripcion, severidad, historiaID, sprintID)
}

func (m *mundo) cambioElEstadoDelDefecto(descripcion, estado string) error {
	return m.cambiarEstadoDefecto(descripcion, estado, 0)
}

func (m *mundo) resuelvoElDefecto(descripcion string, sprint int) error {
	sprintID, err := m.sprintNumero(sprint)
	if err != nil {
		return err
	}
	return m.cambiarEstadoDefecto(descripcion, "Resuelto", sprintID)
}

func (m *mundo) resuelvoElDefectoEnElSprintDe(descripcion, proyecto string) error {
	sprintID, ok := m.defectos.sprintsAjenos[proyecto]
	if !ok {
		return fmt.Errorf("el escenario no creó un Sprint en %q", proyecto)
	}
	return m.cambiarEstadoDefecto(descripcion, "Resuelto", sprintID)
}

func (m *mundo) consultoLasMetricasDeDefectos() error {
	return m.pedir(http.MethodGet, fmt.Sprintf("/api/proyectos/%d/metricas/defectos", m.backlog.principal), nil)
}

func (m *mundo) listoLosDefectos() error {
	return m.pedir(http.MethodGet, rutaDefectos(m.backlog.principal), nil)
}

func (m *mundo) veoElDefectoConID(id int64) error {
	return m.pedir(http.MethodGet, fmt.Sprintf("%s/%d", rutaDefectos(m.backlog.principal), id), nil)
}

// ---------- verificaciones ----------

func (m *mundo) elDefectoQuedaRegistrado() error {
	if m.respuesta.Code != http.StatusCreated || m.defectos.ultimo == 0 {
		return fmt.Errorf("no se registró: código %d, %s", m.respuesta.Code, m.respuesta.Body.String())
	}
	return nil
}

func (m *mundo) elDefectoTiene(descripcion, severidad, historia string, sprint int) error {
	d, err := m.defectoActual()
	if err != nil {
		return err
	}
	if d.Descripcion != descripcion || d.Severidad != severidad || d.HistoriaTitulo != historia || d.SprintDeteccion.Numero != sprint {
		return fmt.Errorf("defecto = %+v, se esperaba %q, %q, %q y Sprint de detección %d", d, descripcion, severidad, historia, sprint)
	}
	return nil
}

func (m *mundo) elDefectoEstaSinResolucion(estado string) error {
	d, err := m.defectoActual()
	if err != nil {
		return err
	}
	if d.Estado != estado || d.SprintResolucion != nil {
		return fmt.Errorf("defecto = %+v, se esperaba %q y sin Sprint de resolución", d, estado)
	}
	return nil
}

func (m *mundo) elDefectoEstaConResolucion(estado string, sprint int) error {
	d, err := m.defectoActual()
	if err != nil {
		return err
	}
	if d.Estado != estado || d.SprintResolucion == nil || d.SprintResolucion.Numero != sprint {
		return fmt.Errorf("defecto = %+v, se esperaba %q con Sprint de resolución %d", d, estado, sprint)
	}
	return nil
}

func (m *mundo) elProyectoTieneDefectos(detectados, resueltos int) error {
	md, err := m.metricasDeDefectosVistas()
	if err != nil {
		return err
	}
	if md.Detectados != detectados || md.Resueltos != resueltos {
		return fmt.Errorf("proyecto: %d detectados y %d resueltos, se esperaban %d y %d", md.Detectados, md.Resueltos, detectados, resueltos)
	}
	return nil
}

func (m *mundo) elSprintTieneDefectos(numero, detectados, resueltos int) error {
	md, err := m.metricasDeDefectosVistas()
	if err != nil {
		return err
	}
	for _, s := range md.Sprints {
		if s.Numero == numero {
			if s.Detectados != detectados || s.Resueltos != resueltos {
				return fmt.Errorf("Sprint %d: %d detectados y %d resueltos, se esperaban %d y %d", numero, s.Detectados, s.Resueltos, detectados, resueltos)
			}
			return nil
		}
	}
	return fmt.Errorf("las métricas no tienen el Sprint %d: %+v", numero, md.Sprints)
}

func (m *mundo) veoLosDefectosEnOrden(primero, segundo string) error {
	var lista []defecto
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &lista); err != nil {
		return err
	}
	if len(lista) != 2 || lista[0].Descripcion != primero || lista[1].Descripcion != segundo {
		return fmt.Errorf("lista = %+v, se esperaba %q y %q en ese orden", lista, primero, segundo)
	}
	return nil
}

func (m *mundo) noHayDefectosRegistrados() error {
	lista, err := m.listaDeDefectos()
	if err != nil {
		return err
	}
	if len(lista) != 0 {
		return fmt.Errorf("hay %d defectos, no se debía registrar ninguno", len(lista))
	}
	return nil
}
