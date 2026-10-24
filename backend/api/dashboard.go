package api

import (
	"net/http"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// dashboardJSON es el Dashboard de un proyecto tal como sale por la API (spec 13): los indicadores y
// los datos de los tres gráficos, ya calculados. Donde no hay dato va null; las listas vacías, [].
type dashboardJSON struct {
	Proyecto         proyectoJSON       `json:"proyecto"`
	Estado           string             `json:"estado"`
	DiasParaEmpezar  *int               `json:"dias_para_empezar"`
	DiasParaTerminar *int               `json:"dias_para_terminar"`
	SprintEnCurso    *sprintEnCursoJSON `json:"sprint_en_curso"`
	Backlog          avanceBacklogJSON  `json:"backlog"`
	Velocidad        *float64           `json:"velocidad"`
	Graficos         graficosJSON       `json:"graficos"`
}

type sprintEnCursoJSON struct {
	ID                      int64  `json:"id"`
	Numero                  int    `json:"numero"`
	SprintGoal              string `json:"sprint_goal"`
	FechaFin                string `json:"fecha_fin"`
	DiasRestantes           int    `json:"dias_restantes"`
	StoryPointsPlanificados int    `json:"story_points_planificados"`
	StoryPointsCompletados  int    `json:"story_points_completados"`
}

type avanceBacklogJSON struct {
	Total                  int `json:"total"`
	Completadas            int `json:"completadas"`
	PorcentajeCompletadas  int `json:"porcentaje_completadas"`
	StoryPointsTotales     int `json:"story_points_totales"`
	StoryPointsCompletados int `json:"story_points_completados"`
	SinEstimar             int `json:"sin_estimar"`
}

type graficosJSON struct {
	StoryPointsPorSprint []barraSprintJSON    `json:"story_points_por_sprint"`
	Velocidad            graficoVelocidadJSON `json:"velocidad"`
	BacklogPorEstado     []porcionBacklogJSON `json:"backlog_por_estado"`
}

type barraSprintJSON struct {
	Sprint       string `json:"sprint"`
	Numero       int    `json:"numero"`
	Planificados int    `json:"planificados"`
	Completados  int    `json:"completados"`
	Parcial      bool   `json:"parcial"`
}

type graficoVelocidadJSON struct {
	Sprints  []barraVelocidadJSON `json:"sprints"`
	Promedio *float64             `json:"promedio"`
}

type barraVelocidadJSON struct {
	Sprint      string `json:"sprint"`
	Numero      int    `json:"numero"`
	Completados int    `json:"completados"`
}

type porcionBacklogJSON struct {
	Estado   string `json:"estado"`
	Cantidad int    `json:"cantidad"`
}

type dashboardHandler struct {
	servicio *dominio.ServicioDashboard
}

// consultar: GET /api/proyectos/{id}/dashboard (spec 13 CA-1 a CA-8).
func (h dashboardHandler) consultar(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	d, err := h.servicio.Consultar(r.Context(), proyectoID)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, dashboardAJSON(d))
}

func dashboardAJSON(d dominio.Dashboard) dashboardJSON {
	var velocidad *float64
	if d.HayVelocidad {
		// El dominio calcula en décimas, con enteros; recién acá se pasa a decimal para el JSON.
		v := float64(d.VelocidadDecimas) / 10
		velocidad = &v
	}
	out := dashboardJSON{
		Proyecto:         aJSON(d.Proyecto),
		Estado:           string(d.Estado),
		DiasParaEmpezar:  d.DiasParaEmpezar,
		DiasParaTerminar: d.DiasParaTerminar,
		Backlog: avanceBacklogJSON{
			Total:                  d.Backlog.Total,
			Completadas:            d.Backlog.Completadas,
			PorcentajeCompletadas:  d.Backlog.PorcentajeCompletadas,
			StoryPointsTotales:     d.Backlog.StoryPointsTotales,
			StoryPointsCompletados: d.Backlog.StoryPointsCompletados,
			SinEstimar:             d.Backlog.SinEstimar,
		},
		Velocidad: velocidad,
		Graficos: graficosJSON{
			StoryPointsPorSprint: make([]barraSprintJSON, 0, len(d.StoryPointsPorSprint)),
			Velocidad:            graficoVelocidadJSON{Sprints: make([]barraVelocidadJSON, 0, len(d.VelocidadPorSprint)), Promedio: velocidad},
			BacklogPorEstado:     make([]porcionBacklogJSON, 0, len(d.BacklogPorEstado)),
		},
	}
	if s := d.SprintEnCurso; s != nil {
		out.SprintEnCurso = &sprintEnCursoJSON{
			ID:                      s.Sprint.ID,
			Numero:                  s.Sprint.Numero,
			SprintGoal:              s.Sprint.Goal,
			FechaFin:                s.Sprint.FechaFin.Format(dominio.FormatoFecha),
			DiasRestantes:           s.DiasRestantes,
			StoryPointsPlanificados: s.StoryPointsPlanificados,
			StoryPointsCompletados:  s.StoryPointsCompletados,
		}
	}
	for _, b := range d.StoryPointsPorSprint {
		out.Graficos.StoryPointsPorSprint = append(out.Graficos.StoryPointsPorSprint,
			barraSprintJSON{Sprint: b.Nombre, Numero: b.Numero, Planificados: b.Planificados, Completados: b.Completados, Parcial: b.Parcial})
	}
	for _, b := range d.VelocidadPorSprint {
		out.Graficos.Velocidad.Sprints = append(out.Graficos.Velocidad.Sprints,
			barraVelocidadJSON{Sprint: b.Nombre, Numero: b.Numero, Completados: b.Completados})
	}
	for _, p := range d.BacklogPorEstado {
		out.Graficos.BacklogPorEstado = append(out.Graficos.BacklogPorEstado, porcionBacklogJSON{Estado: string(p.Estado), Cantidad: p.Cantidad})
	}
	return out
}
