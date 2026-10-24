package dominio

import (
	"context"
	"strconv"
	"time"
)

// Dashboard es lo que muestra la pantalla del Dashboard de un proyecto (spec 13): los indicadores
// y los datos de los tres gráficos, ya calculados. La pantalla solo dibuja.
type Dashboard struct {
	Proyecto Proyecto
	Estado   EstadoProyecto
	// Días que faltan (RN-3): para empezar si no inició, para terminar si está en curso; nil si no
	// corresponde (Finalizado no informa días).
	DiasParaEmpezar  *int
	DiasParaTerminar *int
	SprintEnCurso    *SprintEnCursoDashboard // nil si no hay (spec 03 RN-3)
	Backlog          ResumenBacklog          // el resumen de la spec 03, sin cambios
	HayVelocidad     bool                    // la velocidad de la spec 10, sin cambios
	VelocidadDecimas int

	StoryPointsPorSprint []BarraSprint    // gráfico 1: todos los Sprints (RN-5)
	VelocidadPorSprint   []BarraVelocidad // gráfico 2: solo los cerrados (RN-6)
	BacklogPorEstado     []PorcionBacklog // gráfico 3: pendientes, en progreso y completadas (RN-7)
}

// SprintEnCursoDashboard es el Sprint en curso con los días que le quedan (RN-4) y sus Story Points
// parciales (spec 10 RN-3).
type SprintEnCursoDashboard struct {
	Sprint                  Sprint
	DiasRestantes           int
	StoryPointsPlanificados int
	StoryPointsCompletados  int
}

// BarraSprint son los Story Points planificados y completados de un Sprint (gráfico 1).
type BarraSprint struct {
	Nombre       string // "Sprint 1"
	Numero       int
	Planificados int
	Completados  int
	Parcial      bool // Sprint abierto: sus números todavía pueden cambiar
}

// BarraVelocidad son los Story Points completados de un Sprint cerrado (gráfico 2).
type BarraVelocidad struct {
	Nombre      string
	Numero      int
	Completados int
}

// PorcionBacklog es la cantidad de historias del backlog en un estado (gráfico 3).
type PorcionBacklog struct {
	Estado   EstadoHistoria
	Cantidad int
}

// ArmarDashboard junta el estado del proyecto (spec 03) y sus métricas (spec 10) sin recalcularlos
// (RN-2), y agrega solo lo propio del Dashboard: los días que faltan (RN-3, RN-4) y los datos de los
// gráficos (RN-5 a RN-7). Con un proyecto sin datos devuelve listas vacías y ceros (RN-8).
func ArmarDashboard(e EstadoDeProyecto, m MetricasProyecto, hoy time.Time) Dashboard {
	d := Dashboard{
		Proyecto:             e.Proyecto,
		Estado:               e.Estado,
		Backlog:              e.Backlog,
		HayVelocidad:         m.HayVelocidad,
		VelocidadDecimas:     m.VelocidadDecimas,
		StoryPointsPorSprint: []BarraSprint{},
		VelocidadPorSprint:   []BarraVelocidad{},
		BacklogPorEstado: []PorcionBacklog{
			{EstadoPendiente, e.Backlog.Pendientes},
			{EstadoEnProgreso, e.Backlog.EnProgreso},
			{EstadoCompletada, e.Backlog.Completadas},
		},
	}
	switch e.Estado {
	case EstadoNoIniciado:
		dias := diasEntre(hoy, e.Proyecto.FechaInicio)
		d.DiasParaEmpezar = &dias
	case EstadoEnCurso:
		dias := diasEntre(hoy, e.Proyecto.FechaFin)
		d.DiasParaTerminar = &dias
	}

	for _, s := range m.Sprints {
		nombre := "Sprint " + strconv.Itoa(s.Numero)
		d.StoryPointsPorSprint = append(d.StoryPointsPorSprint, BarraSprint{
			Nombre: nombre, Numero: s.Numero, Parcial: s.Parcial,
			Planificados: s.StoryPointsPlanificados, Completados: s.StoryPointsCompletados,
		})
		if s.Estado == SprintCerrado {
			d.VelocidadPorSprint = append(d.VelocidadPorSprint, BarraVelocidad{Nombre: nombre, Numero: s.Numero, Completados: s.StoryPointsCompletados})
		}
		if e.SprintEnCurso != nil && s.SprintID == e.SprintEnCurso.ID {
			d.SprintEnCurso = &SprintEnCursoDashboard{
				StoryPointsPlanificados: s.StoryPointsPlanificados,
				StoryPointsCompletados:  s.StoryPointsCompletados,
			}
		}
	}
	if e.SprintEnCurso != nil {
		if d.SprintEnCurso == nil { // sin métricas del Sprint: sus Story Points quedan en 0
			d.SprintEnCurso = &SprintEnCursoDashboard{}
		}
		d.SprintEnCurso.Sprint = *e.SprintEnCurso
		d.SprintEnCurso.DiasRestantes = diasEntre(hoy, e.SprintEnCurso.FechaFin)
	}
	return d
}

// diasEntre cuenta los días corridos de desde a hasta, solo por la fecha (sin la hora): 0 si son el
// mismo día.
func diasEntre(desde, hasta time.Time) int {
	return int(soloDia(hasta).Sub(soloDia(desde)).Hours() / 24)
}

// ServicioDashboard es el caso de uso de la spec 13. Reutiliza los servicios del estado (spec 03) y
// de las métricas (spec 10), así un mismo dato vale lo mismo en todas las consultas (RN-2).
type ServicioDashboard struct {
	estado   *ServicioEstado
	metricas *ServicioMetricas
	hoy      func() time.Time
}

func NuevoServicioDashboard(estado *ServicioEstado, metricas *ServicioMetricas, hoy func() time.Time) *ServicioDashboard {
	return &ServicioDashboard{estado: estado, metricas: metricas, hoy: hoy}
}

// Consultar arma el Dashboard del proyecto en el momento de la consulta (RN-9). Si el proyecto no
// existe, devuelve ErrProyectoNoExiste (RN-1).
func (s *ServicioDashboard) Consultar(ctx context.Context, proyectoID int64) (Dashboard, error) {
	e, err := s.estado.Consultar(ctx, proyectoID)
	if err != nil {
		return Dashboard{}, err
	}
	m, err := s.metricas.Consultar(ctx, proyectoID)
	if err != nil {
		return Dashboard{}, err
	}
	return ArmarDashboard(e, m, s.hoy()), nil
}
