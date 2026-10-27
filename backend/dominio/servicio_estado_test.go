package dominio_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/memoria"
)

// estadoDePrueba arma el servicio con repositorios en memoria, dos proyectos y un reloj que el
// test puede mover (la fecha de hoy).
type estadoDePrueba struct {
	servicio *dominio.ServicioEstado
	backlog  *dominio.ServicioBacklog
	sprints  *dominio.ServicioSprints
	hoy      *time.Time
	tienda   dominio.Proyecto // "Tienda Online", del 2026-11-01 al 2026-12-15
	turnos   dominio.Proyecto // "App Turnos"
}

func nuevoEstado(t *testing.T) estadoDePrueba {
	t.Helper()
	ctx := context.Background()
	proyectos := memoria.NuevoRepositorioProyectos()
	historias := memoria.NuevoRepositorioHistorias()
	tienda, err := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "Tienda Online", Descripcion: "Venta de ropa", FechaInicio: f("2026-11-01"), FechaFin: f("2026-12-15")})
	if err != nil {
		t.Fatalf("crear proyecto de partida: %v", err)
	}
	turnos, _ := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "App Turnos", Descripcion: "Turnos médicos", FechaInicio: f("2026-11-05"), FechaFin: f("2027-01-20")})
	hoy := f("2026-11-20")
	sprints := memoria.NuevoRepositorioSprints()
	return estadoDePrueba{
		servicio: dominio.NuevoServicioEstado(proyectos, historias, sprints, func() time.Time { return hoy }),
		backlog:  dominio.NuevoServicioBacklog(proyectos, historias, sprints, memoria.NuevoRepositorioEsfuerzo(), memoria.NuevoRepositorioDefectos()),
		sprints:  dominio.NuevoServicioSprints(proyectos, historias, sprints),
		hoy:      &hoy,
		tienda:   tienda,
		turnos:   turnos,
	}
}

// spec 03 CA-1, CA-3: devuelve los datos del proyecto, su estado según hoy y el resumen de su backlog.
func TestEstadoConsultar(t *testing.T) {
	e := nuevoEstado(t)
	ctx := context.Background()
	h, _ := e.backlog.Crear(ctx, e.tienda.ID, datos("Registrar usuario", 1))
	e.backlog.Modificar(ctx, e.tienda.ID, h.ID, datos("Registrar usuario", 1), "Completada")
	e.backlog.Crear(ctx, e.tienda.ID, datos("Ver catálogo", 2))
	e.backlog.Crear(ctx, e.turnos.ID, datos("Pedir turno", 1)) // de otro proyecto: no cuenta

	estado, err := e.servicio.Consultar(ctx, e.tienda.ID)
	if err != nil {
		t.Fatalf("Consultar: %v", err)
	}
	if estado.Proyecto != e.tienda {
		t.Errorf("proyecto = %+v, se esperaba %+v", estado.Proyecto, e.tienda)
	}
	if estado.Estado != dominio.EstadoEnCurso {
		t.Errorf("estado = %q, se esperaba %q", estado.Estado, dominio.EstadoEnCurso)
	}
	esperado := dominio.ResumenBacklog{Total: 2, Pendientes: 1, Completadas: 1, PorcentajeCompletadas: 50, SinEstimar: 2}
	if estado.Backlog != esperado {
		t.Errorf("resumen = %+v, se esperaba %+v", estado.Backlog, esperado)
	}
}

// spec 03 RN-2, RN-7: el estado usa la fecha de hoy de cada consulta.
func TestEstadoUsaLaFechaDeHoy(t *testing.T) {
	e := nuevoEstado(t)
	ctx := context.Background()
	casos := []struct {
		hoy    string
		estado dominio.EstadoProyecto
	}{
		{"2026-10-31", dominio.EstadoNoIniciado},
		{"2026-11-01", dominio.EstadoEnCurso},
		{"2026-12-16", dominio.EstadoFinalizado},
	}
	for _, c := range casos {
		*e.hoy = f(c.hoy)
		estado, err := e.servicio.Consultar(ctx, e.tienda.ID)
		if err != nil || estado.Estado != c.estado {
			t.Errorf("hoy %s: estado = %q, %v; se esperaba %q", c.hoy, estado.Estado, err, c.estado)
		}
	}
}

// spec 03 RN-7 (CA-3): el resumen refleja los cambios del backlog entre consultas.
func TestEstadoReflejaCambios(t *testing.T) {
	e := nuevoEstado(t)
	ctx := context.Background()
	h, _ := e.backlog.Crear(ctx, e.tienda.ID, datos("Registrar usuario", 1))

	antes, _ := e.servicio.Consultar(ctx, e.tienda.ID)
	e.backlog.Modificar(ctx, e.tienda.ID, h.ID, datos("Registrar usuario", 1), "Completada")
	despues, _ := e.servicio.Consultar(ctx, e.tienda.ID)

	if antes.Backlog.PorcentajeCompletadas != 0 || despues.Backlog.PorcentajeCompletadas != 100 {
		t.Errorf("porcentaje antes = %d, después = %d; se esperaba 0 y 100",
			antes.Backlog.PorcentajeCompletadas, despues.Backlog.PorcentajeCompletadas)
	}
}

// spec 03 CA-3: un proyecto sin historias tiene el resumen en cero (no es un error).
func TestEstadoBacklogVacio(t *testing.T) {
	e := nuevoEstado(t)
	estado, err := e.servicio.Consultar(context.Background(), e.tienda.ID)
	if err != nil || estado.Backlog != (dominio.ResumenBacklog{}) {
		t.Errorf("Consultar = %+v, %v; se esperaba el resumen en cero", estado.Backlog, err)
	}
}

// spec 03 RN-1, CA-4: consultar un proyecto que no existe informa que no existe.
func TestEstadoProyectoInexistente(t *testing.T) {
	e := nuevoEstado(t)
	if _, err := e.servicio.Consultar(context.Background(), 999); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("error = %v, se esperaba %v", err, dominio.ErrProyectoNoExiste)
	}
}

// spec 05 CA-5, RN-12: el estado informa el Sprint en curso según la fecha de hoy.
func TestEstadoSprintEnCurso(t *testing.T) {
	e := nuevoEstado(t)
	ctx := context.Background()
	s1, _ := e.sprints.Crear(ctx, e.tienda.ID, f("2026-11-01"), f("2026-11-14"), "Registrar clientes")
	s2, _ := e.sprints.Crear(ctx, e.tienda.ID, f("2026-11-16"), f("2026-11-28"), "Catálogo")
	casos := []struct {
		hoy    string
		sprint int64 // 0 = no hay Sprint en curso
	}{
		{"2026-11-01", s1.ID},
		{"2026-11-14", s1.ID},
		{"2026-11-15", 0},
		{"2026-11-20", s2.ID},
		{"2026-12-10", 0},
	}
	for _, c := range casos {
		t.Run(c.hoy, func(t *testing.T) {
			*e.hoy = f(c.hoy)
			estado, err := e.servicio.Consultar(ctx, e.tienda.ID)
			if err != nil {
				t.Fatalf("Consultar: %v", err)
			}
			if c.sprint == 0 {
				if estado.SprintEnCurso != nil {
					t.Errorf("Sprint en curso = %+v, se esperaba ninguno", *estado.SprintEnCurso)
				}
				return
			}
			if estado.SprintEnCurso == nil || estado.SprintEnCurso.ID != c.sprint {
				t.Errorf("Sprint en curso = %+v, se esperaba el %d", estado.SprintEnCurso, c.sprint)
			}
		})
	}
}
