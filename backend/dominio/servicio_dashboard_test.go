package dominio_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/memoria"
)

// spec 13 RN-2 (CA-1 a CA-6): el Dashboard da los mismos números que las consultas de estado
// (spec 03) y de métricas (spec 10), porque las reutiliza.
func TestDashboardReutilizaEstadoYMetricas(t *testing.T) {
	ctx := context.Background()
	hoy := func() time.Time { return f("2026-11-20") }
	proyectos := memoria.NuevoRepositorioProyectos()
	historias := memoria.NuevoRepositorioHistorias()
	sprints := memoria.NuevoRepositorioSprints()
	estado := dominio.NuevoServicioEstado(proyectos, historias, sprints, hoy)
	metricas := dominio.NuevoServicioMetricas(proyectos, historias, sprints)
	servicioSprints := dominio.NuevoServicioSprints(proyectos, historias, sprints)
	dashboard := dominio.NuevoServicioDashboard(estado, metricas, hoy)

	p, _ := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "Tienda Online", Descripcion: "d", FechaInicio: f("2026-11-01"), FechaFin: f("2026-12-31")})
	cinco, ocho := 5, 8
	registrar, _ := historias.Guardar(ctx, dominio.Historia{ProyectoID: p.ID, Titulo: "Registrar usuario", Descripcion: "d", Prioridad: 1, Estado: dominio.EstadoPendiente, StoryPoints: &cinco})
	catalogo, _ := historias.Guardar(ctx, dominio.Historia{ProyectoID: p.ID, Titulo: "Ver catálogo", Descripcion: "d", Prioridad: 2, Estado: dominio.EstadoPendiente, StoryPoints: &ocho})
	historias.Guardar(ctx, dominio.Historia{ProyectoID: p.ID, Titulo: "Pagar pedido", Descripcion: "d", Prioridad: 3, Estado: dominio.EstadoPendiente})

	// Sprint 1 cerrado con "Registrar usuario" completada; Sprint 2 en curso con "Ver catálogo".
	s1, err := servicioSprints.Crear(ctx, p.ID, f("2026-11-01"), f("2026-11-14"), "Registro")
	if err != nil {
		t.Fatalf("crear Sprint 1: %v", err)
	}
	servicioSprints.AsignarHistoria(ctx, p.ID, s1.ID, registrar.ID)
	servicioSprints.AsignarHistoria(ctx, p.ID, s1.ID, catalogo.ID)
	servicioSprints.CompletarHistoria(ctx, p.ID, s1.ID, registrar.ID)
	servicioSprints.Cerrar(ctx, p.ID, s1.ID)
	s2, _ := servicioSprints.Crear(ctx, p.ID, f("2026-11-15"), f("2026-11-28"), "Catálogo")
	servicioSprints.AsignarHistoria(ctx, p.ID, s2.ID, catalogo.ID)

	d, err := dashboard.Consultar(ctx, p.ID)
	if err != nil {
		t.Fatalf("Consultar: %v", err)
	}
	e, _ := estado.Consultar(ctx, p.ID)
	m, _ := metricas.Consultar(ctx, p.ID)

	if d.Proyecto != e.Proyecto || d.Estado != e.Estado || d.Backlog != e.Backlog {
		t.Errorf("indicadores = %+v / %q / %+v, se esperaba lo mismo que el estado (spec 03)", d.Proyecto, d.Estado, d.Backlog)
	}
	if d.DiasParaTerminar == nil || *d.DiasParaTerminar != 41 {
		t.Errorf("días para terminar = %v, se esperaba 41", d.DiasParaTerminar)
	}
	if d.SprintEnCurso == nil || d.SprintEnCurso.Sprint.ID != s2.ID || d.SprintEnCurso.DiasRestantes != 8 ||
		d.SprintEnCurso.StoryPointsPlanificados != 8 || d.SprintEnCurso.StoryPointsCompletados != 0 {
		t.Errorf("Sprint en curso = %+v, se esperaba el Sprint 2 con 8 días y 8/0 Story Points", d.SprintEnCurso)
	}
	if d.HayVelocidad != m.HayVelocidad || d.VelocidadDecimas != m.VelocidadDecimas || d.VelocidadDecimas != 50 {
		t.Errorf("velocidad = %v %d, se esperaba la de las métricas (spec 10): 5,0", d.HayVelocidad, d.VelocidadDecimas)
	}
	if len(d.StoryPointsPorSprint) != 2 || d.StoryPointsPorSprint[0].Planificados != 13 || d.StoryPointsPorSprint[0].Completados != 5 ||
		!d.StoryPointsPorSprint[1].Parcial || len(d.VelocidadPorSprint) != 1 {
		t.Errorf("gráficos = %+v / %+v, no coinciden con las métricas de los Sprints", d.StoryPointsPorSprint, d.VelocidadPorSprint)
	}
}

// spec 13 RN-1 (CA-8): el Dashboard de un proyecto que no existe informa que no existe.
func TestDashboardProyectoInexistente(t *testing.T) {
	hoy := func() time.Time { return f("2026-11-20") }
	proyectos := memoria.NuevoRepositorioProyectos()
	historias := memoria.NuevoRepositorioHistorias()
	sprints := memoria.NuevoRepositorioSprints()
	dashboard := dominio.NuevoServicioDashboard(
		dominio.NuevoServicioEstado(proyectos, historias, sprints, hoy),
		dominio.NuevoServicioMetricas(proyectos, historias, sprints),
		hoy,
	)
	if _, err := dashboard.Consultar(context.Background(), 999); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("error = %v, se esperaba %v", err, dominio.ErrProyectoNoExiste)
	}
}
