package dominio_test

import (
	"context"
	"errors"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/memoria"
)

// metricasDePrueba arma los servicios con repositorios en memoria y el proyecto "Tienda Online"
// (del 2026-11-01 al 2026-12-31) con dos historias estimadas: Registrar usuario (5) y Ver catálogo (8).
type metricasDePrueba struct {
	sprintsDePrueba
	metricas            *dominio.ServicioMetricas
	registrar, catalogo dominio.Historia
}

func nuevasMetricas(t *testing.T) metricasDePrueba {
	t.Helper()
	ctx := context.Background()
	proyectos := memoria.NuevoRepositorioProyectos()
	historias := memoria.NuevoRepositorioHistorias()
	sprints := memoria.NuevoRepositorioSprints()
	tienda, _ := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "Tienda Online", Descripcion: "Venta de ropa", FechaInicio: f("2026-11-01"), FechaFin: f("2026-12-31")})
	m := metricasDePrueba{
		sprintsDePrueba: sprintsDePrueba{
			servicio: dominio.NuevoServicioSprints(proyectos, historias, sprints),
			backlog:  dominio.NuevoServicioBacklog(proyectos, historias, sprints, memoria.NuevoRepositorioEsfuerzo()),
			sprints:  sprints,
			tienda:   tienda.ID,
		},
		metricas: dominio.NuevoServicioMetricas(proyectos, historias, sprints),
	}
	m.registrar, _ = m.backlog.Crear(ctx, tienda.ID, datos("Registrar usuario", 1))
	m.catalogo, _ = m.backlog.Crear(ctx, tienda.ID, datos("Ver catálogo", 2))
	if _, err := m.backlog.Estimar(ctx, tienda.ID, m.registrar.ID, 5); err != nil {
		t.Fatalf("estimar historia de partida: %v", err)
	}
	m.backlog.Estimar(ctx, tienda.ID, m.catalogo.ID, 8)
	return m
}

// spec 10 CA-1 a CA-4, CA-6: métricas por Sprint (cerrado y abierto parcial) y velocidad.
func TestMetricasConsultar(t *testing.T) {
	m := nuevasMetricas(t)
	ctx := context.Background()
	s1 := m.crearSprint(t, m.tienda, "2026-11-01", "2026-11-14")
	m.servicio.AsignarHistoria(ctx, m.tienda, s1.ID, m.registrar.ID)
	m.servicio.AsignarHistoria(ctx, m.tienda, s1.ID, m.catalogo.ID)
	m.servicio.CompletarHistoria(ctx, m.tienda, s1.ID, m.registrar.ID)
	m.servicio.Cerrar(ctx, m.tienda, s1.ID)
	s2 := m.crearSprint(t, m.tienda, "2026-11-15", "2026-11-28")
	m.servicio.AsignarHistoria(ctx, m.tienda, s2.ID, m.catalogo.ID)

	r, err := m.metricas.Consultar(ctx, m.tienda)
	if err != nil {
		t.Fatalf("Consultar: %v", err)
	}
	if len(r.Sprints) != 2 {
		t.Fatalf("sprints = %+v, se esperaban 2", r.Sprints)
	}
	primero := dominio.MetricasSprint{SprintID: s1.ID, Numero: 1, Estado: dominio.SprintCerrado, HistoriasPlanificadas: 2,
		HistoriasCompletadas: 1, PorcentajeCompletadas: 50, StoryPointsPlanificados: 13, StoryPointsCompletados: 5}
	segundo := dominio.MetricasSprint{SprintID: s2.ID, Numero: 2, Estado: dominio.SprintAbierto, Parcial: true,
		HistoriasPlanificadas: 1, StoryPointsPlanificados: 8}
	if r.Sprints[0] != primero || r.Sprints[1] != segundo {
		t.Errorf("sprints =\n%+v\n%+v\nse esperaba\n%+v\n%+v", r.Sprints[0], r.Sprints[1], primero, segundo)
	}
	if !r.HayVelocidad || r.VelocidadDecimas != 50 || r.SprintsCerrados != 1 {
		t.Errorf("velocidad = %d décimas (hay = %v), cerrados = %d; se esperaba 5,0 con 1 cerrado", r.VelocidadDecimas, r.HayVelocidad, r.SprintsCerrados)
	}
}

// spec 10 RN-4, CA-6: al cerrar se congelan los Story Points; re-estimar después no cambia el Sprint.
func TestMetricasStoryPointsCongeladosAlCerrar(t *testing.T) {
	m := nuevasMetricas(t)
	ctx := context.Background()
	s1 := m.crearSprint(t, m.tienda, "2026-11-01", "2026-11-14")
	m.servicio.AsignarHistoria(ctx, m.tienda, s1.ID, m.catalogo.ID)
	m.servicio.Cerrar(ctx, m.tienda, s1.ID)

	if _, err := m.backlog.Estimar(ctx, m.tienda, m.catalogo.ID, 3); err != nil {
		t.Fatalf("re-estimar después del cierre: %v", err)
	}
	r, _ := m.metricas.Consultar(ctx, m.tienda)
	if r.Sprints[0].StoryPointsPlanificados != 8 {
		t.Errorf("SP planificados del Sprint cerrado = %d, se esperaba 8 (congelados al cerrar)", r.Sprints[0].StoryPointsPlanificados)
	}
	cerrado, _ := m.sprints.Obtener(ctx, m.tienda, s1.ID)
	if p := cerrado.StoryPointsAlCerrar[m.catalogo.ID]; p == nil || *p != 8 {
		t.Errorf("Story Points congelados = %v, se esperaba 8", p)
	}
}

// spec 10 casos límite y CA-7.
func TestMetricasSinSprintsYProyectoInexistente(t *testing.T) {
	m := nuevasMetricas(t)
	r, err := m.metricas.Consultar(context.Background(), m.tienda)
	if err != nil || len(r.Sprints) != 0 || r.HayVelocidad || r.SprintsCerrados != 0 {
		t.Errorf("sin Sprints = %+v, %v; se esperaba lista vacía y sin velocidad", r, err)
	}
	if _, err := m.metricas.Consultar(context.Background(), 999); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("proyecto inexistente: error = %v, se esperaba %v", err, dominio.ErrProyectoNoExiste)
	}
}
