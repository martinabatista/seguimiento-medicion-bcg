package features

// Pasos de godog para los escenarios del estado del proyecto (spec 03, estado-proyecto.feature).

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cucumber/godog"
)

func (m *mundo) registrarPasosEstado(sc *godog.ScenarioContext) {
	// Contexto y acciones
	sc.Step(`^que hoy es "([^"]*)"$`, m.queHoyEs)
	sc.Step(`^que el backlog tiene la historia "([^"]*)" en estado "([^"]*)" sin estimar$`, m.queElBacklogTieneLaHistoriaEnEstado)
	sc.Step(`^que el backlog tiene la historia "([^"]*)" en estado "([^"]*)" con (\d+) Story Points$`, m.queElBacklogTieneLaHistoriaEstimada)
	sc.Step(`^que el backlog del proyecto "([^"]*)" tiene la historia "([^"]*)" en estado "([^"]*)" sin estimar$`, m.queElBacklogDelProyectoTieneLaHistoriaEnEstado)
	sc.Step(`^que el backlog tiene (\d+) historias de las cuales (\d+) están completadas$`, m.queElBacklogTieneHistoriasCompletadas)
	sc.Step(`^la historia "([^"]*)" pasa a estado "([^"]*)"$`, m.laHistoriaPasaAEstado)
	sc.Step(`^(?:que )?consulto el estado del proyecto$`, m.consultoElEstadoDelProyecto)
	sc.Step(`^consulto el estado del proyecto "([^"]*)"$`, m.consultoElEstadoDelProyectoLlamado)
	sc.Step(`^consulto el estado del proyecto con identificador (\d+)$`, m.consultoElEstadoDelProyectoConID)

	// Verificaciones
	sc.Step(`^veo el proyecto "([^"]*)" con inicio "([^"]*)" y fin "([^"]*)"$`, m.veoElProyectoConFechas)
	sc.Step(`^el estado del proyecto es "([^"]*)"$`, m.elEstadoDelProyectoEs)
	sc.Step(`^el resumen del backlog tiene (\d+) historias?: (\d+) pendientes?, (\d+) en progreso y (\d+) completadas?$`, m.elResumenDelBacklogTiene)
	sc.Step(`^el porcentaje de historias completadas es (\d+)$`, m.elPorcentajeDeCompletadasEs)
	sc.Step(`^los Story Points son (\d+) en total y (\d+) completados, con (\d+) historias? sin estimar$`, m.losStoryPointsSon)
}

// ---------- helpers ----------

type estadoDelProyecto struct {
	Proyecto proyecto `json:"proyecto"`
	Estado   string   `json:"estado"`
	Backlog  struct {
		Total                  int `json:"total"`
		Pendientes             int `json:"pendientes"`
		EnProgreso             int `json:"en_progreso"`
		Completadas            int `json:"completadas"`
		PorcentajeCompletadas  int `json:"porcentaje_completadas"`
		StoryPointsTotales     int `json:"story_points_totales"`
		StoryPointsCompletados int `json:"story_points_completados"`
		SinEstimar             int `json:"sin_estimar"`
	} `json:"backlog"`
}

// estadoConsultado lee la respuesta de la última consulta del estado.
func (m *mundo) estadoConsultado() (estadoDelProyecto, error) {
	var e estadoDelProyecto
	if m.respuesta.Code != http.StatusOK {
		return e, fmt.Errorf("código = %d, se esperaba 200 (%s)", m.respuesta.Code, m.respuesta.Body.String())
	}
	return e, json.Unmarshal(m.respuesta.Body.Bytes(), &e)
}

// historiaEnEstado crea la historia en el backlog del proyecto y, si hace falta, la pasa al estado
// pedido (una historia nueva siempre arranca Pendiente, spec 04 RN-6).
func (m *mundo) historiaEnEstado(proyectoID int64, titulo, estado string) error {
	if err := m.crearHistoria(proyectoID, titulo, "Alta de clientes", "1", ""); err != nil {
		return err
	}
	if m.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("no se pudo crear la historia de partida: %s", m.respuesta.Body.String())
	}
	if estado == "Pendiente" {
		return nil
	}
	return m.cambiarEstadoDeHistoria(proyectoID, m.backlog.historiaID, titulo, estado)
}

func (m *mundo) cambiarEstadoDeHistoria(proyectoID, historiaID int64, titulo, estado string) error {
	if err := m.pedir(http.MethodPut, rutaHistoria(proyectoID, historiaID), map[string]any{
		"titulo": titulo, "descripcion": "Alta de clientes", "prioridad": 1, "estado": estado,
	}); err != nil {
		return err
	}
	if m.respuesta.Code != http.StatusOK {
		return fmt.Errorf("no se pudo pasar la historia %q a %q: %s", titulo, estado, m.respuesta.Body.String())
	}
	return nil
}

// ---------- contexto y acciones ----------

func (m *mundo) queHoyEs(texto string) error {
	hoy, err := time.Parse("2006-01-02", texto)
	if err != nil {
		return fmt.Errorf("fecha de hoy inválida %q: %w", texto, err)
	}
	m.hoy = hoy
	return nil
}

func (m *mundo) queElBacklogTieneLaHistoriaEnEstado(titulo, estado string) error {
	return m.historiaEnEstado(m.backlog.principal, titulo, estado)
}

// queElBacklogTieneLaHistoriaEstimada la crea, la estima y después la pasa al estado pedido: una
// historia completada ya no se puede estimar (spec 07).
func (m *mundo) queElBacklogTieneLaHistoriaEstimada(titulo, estado string, storyPoints int) error {
	if err := m.historiaEnEstado(m.backlog.principal, titulo, "Pendiente"); err != nil {
		return err
	}
	if err := m.queEstimeLaHistoria(titulo, storyPoints); err != nil {
		return err
	}
	if estado == "Pendiente" {
		return nil
	}
	id, err := m.historiaLlamada(titulo)
	if err != nil {
		return err
	}
	return m.cambiarEstadoDeHistoria(m.backlog.principal, id, titulo, estado)
}

func (m *mundo) queElBacklogDelProyectoTieneLaHistoriaEnEstado(proyecto, titulo, estado string) error {
	id, err := m.proyectoLlamado(proyecto)
	if err != nil {
		return err
	}
	return m.historiaEnEstado(id, titulo, estado)
}

func (m *mundo) queElBacklogTieneHistoriasCompletadas(total, completadas int) error {
	for i := 1; i <= total; i++ {
		estado := "Pendiente"
		if i <= completadas {
			estado = "Completada"
		}
		if err := m.historiaEnEstado(m.backlog.principal, fmt.Sprintf("Historia %d", i), estado); err != nil {
			return err
		}
	}
	return nil
}

// laHistoriaPasaAEstado pide el cambio de estado sin exigir que salga bien: los pasos siguientes
// verifican el resultado (por ejemplo, el error 409 de la spec 06 CA-5).
func (m *mundo) laHistoriaPasaAEstado(titulo, estado string) error {
	id, ok := m.backlog.historias[titulo]
	if !ok {
		return fmt.Errorf("el escenario no creó la historia %q", titulo)
	}
	return m.pedir(http.MethodPut, rutaHistoria(m.backlog.principal, id), map[string]any{
		"titulo": titulo, "descripcion": "Alta de clientes", "prioridad": 1, "estado": estado,
	})
}

func (m *mundo) consultoElEstadoDelProyecto() error {
	return m.consultoElEstadoDelProyectoConID(m.backlog.principal)
}

func (m *mundo) consultoElEstadoDelProyectoLlamado(nombre string) error {
	id, err := m.proyectoLlamado(nombre)
	if err != nil {
		return err
	}
	return m.consultoElEstadoDelProyectoConID(id)
}

func (m *mundo) consultoElEstadoDelProyectoConID(id int64) error {
	return m.pedir(http.MethodGet, fmt.Sprintf("/api/proyectos/%d/estado", id), nil)
}

// ---------- verificaciones ----------

func (m *mundo) veoElProyectoConFechas(nombre, inicio, fin string) error {
	e, err := m.estadoConsultado()
	if err != nil {
		return err
	}
	if e.Proyecto.Nombre != nombre || e.Proyecto.FechaInicio != inicio || e.Proyecto.FechaFin != fin {
		return fmt.Errorf("proyecto = %+v, se esperaba %q del %s al %s", e.Proyecto, nombre, inicio, fin)
	}
	return nil
}

func (m *mundo) elEstadoDelProyectoEs(estado string) error {
	e, err := m.estadoConsultado()
	if err != nil {
		return err
	}
	if e.Estado != estado {
		return fmt.Errorf("estado = %q, se esperaba %q", e.Estado, estado)
	}
	return nil
}

func (m *mundo) elResumenDelBacklogTiene(total, pendientes, enProgreso, completadas int) error {
	e, err := m.estadoConsultado()
	if err != nil {
		return err
	}
	b := e.Backlog
	if b.Total != total || b.Pendientes != pendientes || b.EnProgreso != enProgreso || b.Completadas != completadas {
		return fmt.Errorf("resumen = %+v, se esperaba %d historias: %d pendientes, %d en progreso y %d completadas",
			b, total, pendientes, enProgreso, completadas)
	}
	return nil
}

func (m *mundo) elPorcentajeDeCompletadasEs(porcentaje int) error {
	e, err := m.estadoConsultado()
	if err != nil {
		return err
	}
	if e.Backlog.PorcentajeCompletadas != porcentaje {
		return fmt.Errorf("porcentaje = %d, se esperaba %d", e.Backlog.PorcentajeCompletadas, porcentaje)
	}
	return nil
}

func (m *mundo) losStoryPointsSon(totales, completados, sinEstimar int) error {
	e, err := m.estadoConsultado()
	if err != nil {
		return err
	}
	b := e.Backlog
	if b.StoryPointsTotales != totales || b.StoryPointsCompletados != completados || b.SinEstimar != sinEstimar {
		return fmt.Errorf("Story Points = %d en total y %d completados, %d sin estimar; se esperaba %d, %d y %d",
			b.StoryPointsTotales, b.StoryPointsCompletados, b.SinEstimar, totales, completados, sinEstimar)
	}
	return nil
}
