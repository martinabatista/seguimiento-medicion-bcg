package dominio

import (
	"testing"
	"time"
)

// Proyecto de prueba: del 1/11/2026 al 31/12/2026.
var tiendaDashboard = Proyecto{ID: 1, Nombre: "Tienda Online", FechaInicio: dia("2026-11-01"), FechaFin: dia("2026-12-31")}

func dia(s string) time.Time {
	d, _ := time.Parse(FormatoFecha, s)
	return d
}

// estadoEl arma el estado de la spec 03 del proyecto de prueba para la fecha de hoy.
func estadoEl(hoy time.Time) EstadoDeProyecto {
	return EstadoDeProyecto{Proyecto: tiendaDashboard, Estado: CalcularEstado(tiendaDashboard, hoy)}
}

func valor(p *int) any {
	if p == nil {
		return "ninguno"
	}
	return *p
}

// spec 13 RN-3 (CA-1): días que faltan para empezar o para terminar, según el estado.
func TestDiasQueFaltanDelProyecto(t *testing.T) {
	casos := []struct {
		nombre            string
		hoy               time.Time
		estado            EstadoProyecto
		empezar, terminar *int
	}{
		{"en curso", dia("2026-11-20"), EstadoEnCurso, nil, ptr(41)},
		{"termina hoy", dia("2026-12-31"), EstadoEnCurso, nil, ptr(0)},
		{"empieza hoy", dia("2026-11-01"), EstadoEnCurso, nil, ptr(60)},
		{"todavía no empezó", dia("2026-10-25"), EstadoNoIniciado, ptr(7), nil},
		{"empieza mañana", dia("2026-10-31"), EstadoNoIniciado, ptr(1), nil},
		{"finalizado: no se informan días", dia("2027-01-05"), EstadoFinalizado, nil, nil},
		{"cuenta solo el día, no la hora", time.Date(2026, 11, 20, 23, 30, 0, 0, time.UTC), EstadoEnCurso, nil, ptr(41)},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			d := ArmarDashboard(estadoEl(c.hoy), MetricasProyecto{}, c.hoy)
			if d.Estado != c.estado || !mismoPuntero(d.DiasParaEmpezar, c.empezar) || !mismoPuntero(d.DiasParaTerminar, c.terminar) {
				t.Errorf("estado %q, para empezar %v, para terminar %v; se esperaba %q, %v y %v",
					d.Estado, valor(d.DiasParaEmpezar), valor(d.DiasParaTerminar), c.estado, valor(c.empezar), valor(c.terminar))
			}
		})
	}
}

// spec 13 RN-4 (CA-2): el Sprint en curso con los días que le quedan y sus Story Points parciales.
func TestSprintEnCursoDelDashboard(t *testing.T) {
	hoy := dia("2026-11-20")
	sprint := Sprint{ID: 5, Numero: 2, Goal: "Vender online", FechaInicio: dia("2026-11-15"), FechaFin: dia("2026-11-28"), Estado: SprintAbierto}
	metricas := MetricasProyecto{Sprints: []MetricasSprint{
		{SprintID: 4, Numero: 1, Estado: SprintCerrado, StoryPointsPlanificados: 8, StoryPointsCompletados: 8},
		{SprintID: 5, Numero: 2, Estado: SprintAbierto, Parcial: true, StoryPointsPlanificados: 13, StoryPointsCompletados: 5},
	}}

	t.Run("hay Sprint en curso", func(t *testing.T) {
		e := estadoEl(hoy)
		e.SprintEnCurso = &sprint
		d := ArmarDashboard(e, metricas, hoy)
		s := d.SprintEnCurso
		if s == nil || s.Sprint.ID != 5 || s.DiasRestantes != 8 || s.StoryPointsPlanificados != 13 || s.StoryPointsCompletados != 5 {
			t.Errorf("Sprint en curso = %+v, se esperaba el 2 con 8 días y 13/5 Story Points", s)
		}
	})
	t.Run("el Sprint en curso termina hoy", func(t *testing.T) {
		e := estadoEl(dia("2026-11-28"))
		e.SprintEnCurso = &sprint
		if d := ArmarDashboard(e, metricas, dia("2026-11-28")); d.SprintEnCurso == nil || d.SprintEnCurso.DiasRestantes != 0 {
			t.Errorf("Sprint en curso = %+v, se esperaba 0 días", d.SprintEnCurso)
		}
	})
	t.Run("no hay Sprint en curso", func(t *testing.T) {
		if d := ArmarDashboard(estadoEl(hoy), metricas, hoy); d.SprintEnCurso != nil {
			t.Errorf("Sprint en curso = %+v, se esperaba ninguno", d.SprintEnCurso)
		}
	})
}

// spec 13 RN-2, RN-5, RN-6 (CA-3 a CA-6): los indicadores y los gráficos salen del estado y de las
// métricas tal cual (sin recalcular), con las barras en el orden de los Sprints.
func TestGraficosDelDashboard(t *testing.T) {
	hoy := dia("2026-11-20")
	e := estadoEl(hoy)
	e.Backlog = ResumenBacklog{Total: 5, Pendientes: 3, EnProgreso: 1, Completadas: 1, PorcentajeCompletadas: 20, StoryPointsTotales: 13, StoryPointsCompletados: 5, SinEstimar: 1}
	metricas := MetricasProyecto{
		Sprints: []MetricasSprint{
			{SprintID: 4, Numero: 1, Estado: SprintCerrado, StoryPointsPlanificados: 13, StoryPointsCompletados: 5},
			{SprintID: 5, Numero: 2, Estado: SprintAbierto, Parcial: true, StoryPointsPlanificados: 8},
		},
		SprintsCerrados: 1, HayVelocidad: true, VelocidadDecimas: 50,
	}

	d := ArmarDashboard(e, metricas, hoy)

	if d.Backlog != e.Backlog {
		t.Errorf("backlog = %+v, se esperaba el resumen de la spec 03 sin cambios", d.Backlog)
	}
	if !d.HayVelocidad || d.VelocidadDecimas != 50 {
		t.Errorf("velocidad = %v %d, se esperaba la de la spec 10 (5,0)", d.HayVelocidad, d.VelocidadDecimas)
	}
	esperadas := []BarraSprint{
		{Nombre: "Sprint 1", Numero: 1, Planificados: 13, Completados: 5},
		{Nombre: "Sprint 2", Numero: 2, Planificados: 8, Completados: 0, Parcial: true},
	}
	if len(d.StoryPointsPorSprint) != 2 || d.StoryPointsPorSprint[0] != esperadas[0] || d.StoryPointsPorSprint[1] != esperadas[1] {
		t.Errorf("Story Points por Sprint = %+v, se esperaba %+v", d.StoryPointsPorSprint, esperadas)
	}
	if len(d.VelocidadPorSprint) != 1 || d.VelocidadPorSprint[0] != (BarraVelocidad{Nombre: "Sprint 1", Numero: 1, Completados: 5}) {
		t.Errorf("velocidad por Sprint = %+v, se esperaba solo el Sprint 1 cerrado con 5", d.VelocidadPorSprint)
	}
	porEstado := []PorcionBacklog{{EstadoPendiente, 3}, {EstadoEnProgreso, 1}, {EstadoCompletada, 1}}
	if len(d.BacklogPorEstado) != 3 || d.BacklogPorEstado[0] != porEstado[0] || d.BacklogPorEstado[1] != porEstado[1] || d.BacklogPorEstado[2] != porEstado[2] {
		t.Errorf("backlog por estado = %+v, se esperaba %+v", d.BacklogPorEstado, porEstado)
	}
}

// spec 13 RN-8 (CA-8): un proyecto sin datos igual tiene Dashboard, con listas vacías (no nil) y ceros.
func TestDashboardVacio(t *testing.T) {
	hoy := dia("2026-11-20")
	d := ArmarDashboard(estadoEl(hoy), MetricasProyecto{}, hoy)

	if d.StoryPointsPorSprint == nil || len(d.StoryPointsPorSprint) != 0 || d.VelocidadPorSprint == nil || len(d.VelocidadPorSprint) != 0 {
		t.Errorf("gráficos de Sprints = %#v y %#v, se esperaban listas vacías", d.StoryPointsPorSprint, d.VelocidadPorSprint)
	}
	if d.HayVelocidad || d.SprintEnCurso != nil {
		t.Errorf("velocidad %v, Sprint en curso %+v; se esperaba que no haya ninguno", d.HayVelocidad, d.SprintEnCurso)
	}
	porEstado := []PorcionBacklog{{EstadoPendiente, 0}, {EstadoEnProgreso, 0}, {EstadoCompletada, 0}}
	if len(d.BacklogPorEstado) != 3 || d.BacklogPorEstado[0] != porEstado[0] || d.BacklogPorEstado[2] != porEstado[2] {
		t.Errorf("backlog por estado = %+v, se esperaba %+v", d.BacklogPorEstado, porEstado)
	}
}

func ptr(n int) *int { return &n }

func mismoPuntero(a, b *int) bool {
	return (a == nil && b == nil) || (a != nil && b != nil && *a == *b)
}
