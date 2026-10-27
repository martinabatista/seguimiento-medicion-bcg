package dominio_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/memoria"
)

// spec 12 CA-3, CA-4, CA-5: el servicio junta las historias, los Sprints y el esfuerzo registrado.
func TestMetricasEsfuerzoConsultar(t *testing.T) {
	ctx := context.Background()
	proyectos := memoria.NuevoRepositorioProyectos()
	historias := memoria.NuevoRepositorioHistorias()
	sprints := memoria.NuevoRepositorioSprints()
	integrantes := memoria.NuevoRepositorioIntegrantes()
	esfuerzo := memoria.NuevoRepositorioEsfuerzo()
	backlog := dominio.NuevoServicioBacklog(proyectos, historias, sprints, esfuerzo)
	servicioSprints := dominio.NuevoServicioSprints(proyectos, historias, sprints)
	hoy := func() time.Time { return f("2026-11-20") }
	registro := dominio.NuevoServicioEsfuerzo(proyectos, historias, integrantes, esfuerzo, hoy)
	metricas := dominio.NuevoServicioMetricasEsfuerzo(proyectos, historias, sprints, esfuerzo)

	tienda, _ := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "Tienda Online", Descripcion: "Venta de ropa", FechaInicio: f("2026-11-01"), FechaFin: f("2026-12-31")})
	ana, _ := integrantes.Guardar(ctx, dominio.Integrante{ProyectoID: tienda.ID, Nombre: "Ana", Email: "ana@mail.com", Rol: dominio.RolDeveloper})
	registrar, _ := backlog.Crear(ctx, tienda.ID, datos("Registrar usuario", 1))
	catalogo, _ := backlog.Crear(ctx, tienda.ID, datos("Ver catálogo", 2))
	backlog.EstimarHoras(ctx, tienda.ID, registrar.ID, "8")
	backlog.EstimarHoras(ctx, tienda.ID, catalogo.ID, "4")
	for _, r := range []struct {
		historia int64
		horas    string
	}{{registrar.ID, "6"}, {registrar.ID, "4"}, {catalogo.ID, "2"}} {
		if _, err := registro.Registrar(ctx, tienda.ID, r.historia, dominio.DatosEsfuerzo{IntegranteID: ana.ID, Fecha: "2026-11-10", Actividad: "Programar", Horas: r.horas}); err != nil {
			t.Fatalf("registrar esfuerzo de partida: %v", err)
		}
	}
	s1, err := servicioSprints.Crear(ctx, tienda.ID, f("2026-11-01"), f("2026-11-14"), "Primer Sprint")
	if err != nil {
		t.Fatalf("crear Sprint de partida: %v", err)
	}
	servicioSprints.AsignarHistoria(ctx, tienda.ID, s1.ID, registrar.ID)
	servicioSprints.AsignarHistoria(ctx, tienda.ID, s1.ID, catalogo.ID)
	servicioSprints.CompletarHistoria(ctx, tienda.ID, s1.ID, registrar.ID)

	m, err := metricas.Consultar(ctx, tienda.ID)
	if err != nil {
		t.Fatalf("Consultar: %v", err)
	}
	// Registrar usuario: completada, 8 h estimadas y 10 h reales → +2 h, +25 %.
	if h := m.Historias[0]; h.HistoriaID != registrar.ID || h.Reales != 1000 || h.Desviacion == nil || *h.Desviacion != (dominio.Desviacion{Centesimas: 200, Porcentaje: 25}) {
		t.Errorf("Registrar usuario = %+v, se esperaban 10 h reales y +2 h / +25 %%", h)
	}
	// Ver catálogo: sin terminar, sin desviación.
	if h := m.Historias[1]; h.Reales != 200 || h.Desviacion != nil {
		t.Errorf("Ver catálogo = %+v, se esperaban 2 h reales y sin desviación", h)
	}
	// El Sprint abierto es parcial y solo cuenta la completada en él.
	if len(m.Sprints) != 1 || !m.Sprints[0].Parcial || m.Sprints[0].Estimadas != 800 || m.Sprints[0].Reales != 1000 {
		t.Errorf("sprints = %+v, se esperaba uno parcial con 8 h estimadas y 10 h reales", m.Sprints)
	}
	if m.Estimadas != 1200 || m.Reales != 1200 {
		t.Errorf("proyecto: %d estimadas y %d reales, se esperaban 1200 y 1200", m.Estimadas, m.Reales)
	}

	// CA-7: el proyecto tiene que existir.
	if _, err := metricas.Consultar(ctx, 999); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("proyecto inexistente: err = %v, se esperaba %v", err, dominio.ErrProyectoNoExiste)
	}
}
