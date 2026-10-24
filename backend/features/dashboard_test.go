package features

// Pasos de godog para los escenarios del Dashboard (spec 13, dashboard.feature). Reutiliza los
// pasos de proyectos, backlog, Story Points, Sprints y fecha de hoy para armar el estado, y
// "todavía no hay velocidad del equipo" de las métricas (spec 10): lee el mismo campo "velocidad".

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
)

func (m *mundo) registrarPasosDashboard(sc *godog.ScenarioContext) {
	// Acciones
	sc.Step(`^consulto el Dashboard$`, m.consultoElDashboard)
	sc.Step(`^consulto el Dashboard del proyecto "([^"]*)"$`, m.consultoElDashboardDelProyecto)
	sc.Step(`^consulto el Dashboard del proyecto con identificador (\d+)$`, m.consultoElDashboardDelProyectoConID)

	// Verificaciones
	sc.Step(`^el proyecto está "([^"]*)" y faltan (\d+) días para (empezar|terminar)$`, m.elProyectoEstaYFaltanDias)
	sc.Step(`^el proyecto está "([^"]*)" sin días que falten$`, m.elProyectoEstaSinDias)
	sc.Step(`^el Sprint en curso es el Sprint (\d+) y le quedan (\d+) días$`, m.elSprintEnCursoDelDashboardEs)
	sc.Step(`^el Sprint en curso tiene (\d+) Story Points planificados y (\d+) completados$`, m.elSprintEnCursoTieneStoryPoints)
	sc.Step(`^no hay Sprint en curso$`, m.noHaySprintEnCursoEnElDashboard)
	sc.Step(`^el backlog tiene (\d+) de (\d+) historias completadas, el (\d+) %$`, m.elBacklogTieneHistoriasCompletadas)
	sc.Step(`^el backlog tiene (\d+) de (\d+) Story Points completados y (\d+) historias? sin estimar$`, m.elBacklogTieneStoryPointsCompletados)
	sc.Step(`^la velocidad del equipo es ([\d.]+)$`, m.laVelocidadDelEquipoEs)
	sc.Step(`^el gráfico de Story Points muestra el "([^"]*)" con (\d+) planificados y (\d+) completados(, parcial)?$`, m.elGraficoDeStoryPointsMuestra)
	sc.Step(`^el gráfico de velocidad muestra solo el "([^"]*)" con (\d+) Story Points completados$`, m.elGraficoDeVelocidadMuestraSolo)
	sc.Step(`^el promedio del gráfico de velocidad es ([\d.]+)$`, m.elPromedioDelGraficoDeVelocidadEs)
	sc.Step(`^el gráfico del backlog muestra (\d+) pendientes?, (\d+) en progreso y (\d+) completadas?$`, m.elGraficoDelBacklogMuestra)
	sc.Step(`^los gráficos de Sprints están vacíos$`, m.losGraficosDeSprintsEstanVacios)
}

// ---------- helpers ----------

type dashboard struct {
	Estado           string `json:"estado"`
	DiasParaEmpezar  *int   `json:"dias_para_empezar"`
	DiasParaTerminar *int   `json:"dias_para_terminar"`
	SprintEnCurso    *struct {
		Numero                  int `json:"numero"`
		DiasRestantes           int `json:"dias_restantes"`
		StoryPointsPlanificados int `json:"story_points_planificados"`
		StoryPointsCompletados  int `json:"story_points_completados"`
	} `json:"sprint_en_curso"`
	Backlog struct {
		Total                  int `json:"total"`
		Completadas            int `json:"completadas"`
		PorcentajeCompletadas  int `json:"porcentaje_completadas"`
		StoryPointsTotales     int `json:"story_points_totales"`
		StoryPointsCompletados int `json:"story_points_completados"`
		SinEstimar             int `json:"sin_estimar"`
	} `json:"backlog"`
	Velocidad *float64 `json:"velocidad"`
	Graficos  struct {
		StoryPointsPorSprint []struct {
			Sprint       string `json:"sprint"`
			Planificados int    `json:"planificados"`
			Completados  int    `json:"completados"`
			Parcial      bool   `json:"parcial"`
		} `json:"story_points_por_sprint"`
		Velocidad struct {
			Sprints []struct {
				Sprint      string `json:"sprint"`
				Completados int    `json:"completados"`
			} `json:"sprints"`
			Promedio *float64 `json:"promedio"`
		} `json:"velocidad"`
		BacklogPorEstado []struct {
			Estado   string `json:"estado"`
			Cantidad int    `json:"cantidad"`
		} `json:"backlog_por_estado"`
	} `json:"graficos"`
}

// dashboardRespondido lee el Dashboard de la última respuesta (la de "consulto el Dashboard").
func (m *mundo) dashboardRespondido() (dashboard, error) {
	var d dashboard
	if m.respuesta.Code != http.StatusOK {
		return d, fmt.Errorf("no se pudo consultar el Dashboard: código %d, %s", m.respuesta.Code, m.respuesta.Body.String())
	}
	return d, json.Unmarshal(m.respuesta.Body.Bytes(), &d)
}

func textoDeVelocidad(v *float64) string {
	if v == nil {
		return "ninguna"
	}
	return fmt.Sprint(*v)
}

// ---------- acciones ----------

func (m *mundo) consultoElDashboard() error {
	return m.consultoElDashboardDelProyectoConID(m.backlog.principal)
}

func (m *mundo) consultoElDashboardDelProyecto(nombre string) error {
	id, err := m.proyectoLlamado(nombre)
	if err != nil {
		return err
	}
	return m.consultoElDashboardDelProyectoConID(id)
}

func (m *mundo) consultoElDashboardDelProyectoConID(id int64) error {
	return m.pedir(http.MethodGet, fmt.Sprintf("/api/proyectos/%d/dashboard", id), nil)
}

// ---------- verificaciones ----------

func (m *mundo) elProyectoEstaYFaltanDias(estado string, dias int, para string) error {
	d, err := m.dashboardRespondido()
	if err != nil {
		return err
	}
	esperado, otro := d.DiasParaTerminar, d.DiasParaEmpezar
	if para == "empezar" {
		esperado, otro = d.DiasParaEmpezar, d.DiasParaTerminar
	}
	if d.Estado != estado || esperado == nil || *esperado != dias || otro != nil {
		return fmt.Errorf("estado %q, para empezar %v, para terminar %v; se esperaba %q y %d días para %s",
			d.Estado, d.DiasParaEmpezar, d.DiasParaTerminar, estado, dias, para)
	}
	return nil
}

func (m *mundo) elProyectoEstaSinDias(estado string) error {
	d, err := m.dashboardRespondido()
	if err != nil {
		return err
	}
	if d.Estado != estado || d.DiasParaEmpezar != nil || d.DiasParaTerminar != nil {
		return fmt.Errorf("estado %q, para empezar %v, para terminar %v; se esperaba %q sin días", d.Estado, d.DiasParaEmpezar, d.DiasParaTerminar, estado)
	}
	return nil
}

func (m *mundo) elSprintEnCursoDelDashboardEs(numero, dias int) error {
	d, err := m.dashboardRespondido()
	if err != nil {
		return err
	}
	if d.SprintEnCurso == nil || d.SprintEnCurso.Numero != numero || d.SprintEnCurso.DiasRestantes != dias {
		return fmt.Errorf("Sprint en curso = %+v, se esperaba el %d con %d días", d.SprintEnCurso, numero, dias)
	}
	return nil
}

func (m *mundo) elSprintEnCursoTieneStoryPoints(planificados, completados int) error {
	d, err := m.dashboardRespondido()
	if err != nil {
		return err
	}
	if d.SprintEnCurso == nil || d.SprintEnCurso.StoryPointsPlanificados != planificados || d.SprintEnCurso.StoryPointsCompletados != completados {
		return fmt.Errorf("Sprint en curso = %+v, se esperaban %d planificados y %d completados", d.SprintEnCurso, planificados, completados)
	}
	return nil
}

func (m *mundo) noHaySprintEnCursoEnElDashboard() error {
	d, err := m.dashboardRespondido()
	if err != nil {
		return err
	}
	if d.SprintEnCurso != nil {
		return fmt.Errorf("Sprint en curso = %+v, se esperaba ninguno", d.SprintEnCurso)
	}
	return nil
}

func (m *mundo) elBacklogTieneHistoriasCompletadas(completadas, total, porcentaje int) error {
	d, err := m.dashboardRespondido()
	if err != nil {
		return err
	}
	if b := d.Backlog; b.Completadas != completadas || b.Total != total || b.PorcentajeCompletadas != porcentaje {
		return fmt.Errorf("backlog = %d de %d (%d %%), se esperaba %d de %d (%d %%)", b.Completadas, b.Total, b.PorcentajeCompletadas, completadas, total, porcentaje)
	}
	return nil
}

func (m *mundo) elBacklogTieneStoryPointsCompletados(completados, totales, sinEstimar int) error {
	d, err := m.dashboardRespondido()
	if err != nil {
		return err
	}
	if b := d.Backlog; b.StoryPointsCompletados != completados || b.StoryPointsTotales != totales || b.SinEstimar != sinEstimar {
		return fmt.Errorf("backlog = %d de %d SP y %d sin estimar, se esperaba %d de %d y %d", b.StoryPointsCompletados, b.StoryPointsTotales, b.SinEstimar, completados, totales, sinEstimar)
	}
	return nil
}

func (m *mundo) laVelocidadDelEquipoEs(velocidad float64) error {
	d, err := m.dashboardRespondido()
	if err != nil {
		return err
	}
	if d.Velocidad == nil || *d.Velocidad != velocidad {
		return fmt.Errorf("velocidad = %s, se esperaba %v", textoDeVelocidad(d.Velocidad), velocidad)
	}
	return nil
}

func (m *mundo) elGraficoDeStoryPointsMuestra(sprint string, planificados, completados int, parcial string) error {
	d, err := m.dashboardRespondido()
	if err != nil {
		return err
	}
	for _, b := range d.Graficos.StoryPointsPorSprint {
		if b.Sprint == sprint {
			if b.Planificados != planificados || b.Completados != completados || b.Parcial != (parcial != "") {
				return fmt.Errorf("barra de %s = %+v, se esperaban %d planificados, %d completados y parcial %v", sprint, b, planificados, completados, parcial != "")
			}
			return nil
		}
	}
	return fmt.Errorf("el gráfico de Story Points no muestra el %s: %+v", sprint, d.Graficos.StoryPointsPorSprint)
}

func (m *mundo) elGraficoDeVelocidadMuestraSolo(sprint string, completados int) error {
	d, err := m.dashboardRespondido()
	if err != nil {
		return err
	}
	if s := d.Graficos.Velocidad.Sprints; len(s) != 1 || s[0].Sprint != sprint || s[0].Completados != completados {
		return fmt.Errorf("gráfico de velocidad = %+v, se esperaba solo el %s con %d", s, sprint, completados)
	}
	return nil
}

func (m *mundo) elPromedioDelGraficoDeVelocidadEs(promedio float64) error {
	d, err := m.dashboardRespondido()
	if err != nil {
		return err
	}
	if p := d.Graficos.Velocidad.Promedio; p == nil || *p != promedio {
		return fmt.Errorf("promedio = %s, se esperaba %v", textoDeVelocidad(p), promedio)
	}
	return nil
}

func (m *mundo) elGraficoDelBacklogMuestra(pendientes, enProgreso, completadas int) error {
	d, err := m.dashboardRespondido()
	if err != nil {
		return err
	}
	p := d.Graficos.BacklogPorEstado
	if len(p) != 3 || p[0].Estado != "Pendiente" || p[0].Cantidad != pendientes || p[1].Estado != "En progreso" ||
		p[1].Cantidad != enProgreso || p[2].Estado != "Completada" || p[2].Cantidad != completadas {
		return fmt.Errorf("gráfico del backlog = %+v, se esperaba %d pendientes, %d en progreso y %d completadas", p, pendientes, enProgreso, completadas)
	}
	return nil
}

func (m *mundo) losGraficosDeSprintsEstanVacios() error {
	d, err := m.dashboardRespondido()
	if err != nil {
		return err
	}
	if len(d.Graficos.StoryPointsPorSprint) != 0 || len(d.Graficos.Velocidad.Sprints) != 0 {
		return fmt.Errorf("gráficos de Sprints = %+v y %+v, se esperaban vacíos", d.Graficos.StoryPointsPorSprint, d.Graficos.Velocidad.Sprints)
	}
	return nil
}
