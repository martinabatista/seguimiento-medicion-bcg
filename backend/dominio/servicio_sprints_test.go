package dominio_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/memoria"
)

// sprintsDePrueba arma el servicio con repositorios en memoria y dos proyectos ya creados.
type sprintsDePrueba struct {
	servicio *dominio.ServicioSprints
	backlog  *dominio.ServicioBacklog
	sprints  *memoria.RepositorioSprints
	tienda   int64 // "Tienda Online", del 2026-11-01 al 2026-12-31
	turnos   int64 // "App Turnos"
}

func nuevosSprints(t *testing.T) sprintsDePrueba {
	t.Helper()
	ctx := context.Background()
	proyectos := memoria.NuevoRepositorioProyectos()
	historias := memoria.NuevoRepositorioHistorias()
	sprints := memoria.NuevoRepositorioSprints()
	tienda, err := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "Tienda Online", Descripcion: "Venta de ropa", FechaInicio: f("2026-11-01"), FechaFin: f("2026-12-31")})
	if err != nil {
		t.Fatalf("crear proyecto de partida: %v", err)
	}
	turnos, _ := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "App Turnos", Descripcion: "Turnos médicos", FechaInicio: f("2026-11-01"), FechaFin: f("2027-01-20")})
	return sprintsDePrueba{
		servicio: dominio.NuevoServicioSprints(proyectos, historias, sprints),
		backlog:  dominio.NuevoServicioBacklog(proyectos, historias, sprints, memoria.NuevoRepositorioEsfuerzo()),
		sprints:  sprints,
		tienda:   tienda.ID,
		turnos:   turnos.ID,
	}
}

// crearSprint crea un Sprint de partida y falla el test si no puede.
func (s sprintsDePrueba) crearSprint(t *testing.T, proyectoID int64, inicio, fin string) dominio.Sprint {
	t.Helper()
	sp, err := s.servicio.Crear(context.Background(), proyectoID, f(inicio), f(fin), "Objetivo")
	if err != nil {
		t.Fatalf("crear Sprint de partida: %v", err)
	}
	return sp
}

// cerrar simula el cierre del Sprint (spec 06) directamente en el repositorio.
func (s sprintsDePrueba) cerrar(t *testing.T, sp dominio.Sprint) {
	t.Helper()
	sp, _ = s.sprints.Obtener(context.Background(), sp.ProyectoID, sp.ID)
	sp.Estado = dominio.SprintCerrado
	if err := s.sprints.Actualizar(context.Background(), sp); err != nil {
		t.Fatalf("cerrar Sprint: %v", err)
	}
}

// spec 05 CA-1, RN-6, RN-13: crea el Sprint abierto, numerado y con identificador.
func TestSprintsCrear(t *testing.T) {
	s := nuevosSprints(t)
	ctx := context.Background()

	a, err := s.servicio.Crear(ctx, s.tienda, f("2026-11-01"), f("2026-11-14"), "Registrar clientes")
	if err != nil {
		t.Fatalf("Crear: %v", err)
	}
	b, err := s.servicio.Crear(ctx, s.tienda, f("2026-11-15"), f("2026-11-28"), "Catálogo")
	if err != nil {
		t.Fatalf("Crear el segundo: %v", err)
	}
	c, _ := s.servicio.Crear(ctx, s.turnos, f("2026-11-01"), f("2026-11-14"), "Turnos") // otro proyecto: numera desde 1
	if a.ID == 0 || a.ID == b.ID || a.Numero != 1 || b.Numero != 2 || c.Numero != 1 {
		t.Errorf("sprints = %+v, %+v, %+v; se esperaban IDs distintos y números 1, 2 y 1", a, b, c)
	}
	if a.Estado != dominio.SprintAbierto || a.ProyectoID != s.tienda {
		t.Errorf("sprint = %+v, se esperaba abierto y del proyecto %d", a, s.tienda)
	}

	// RN-4: no se superpone con otro del mismo proyecto; no se registra nada.
	if _, err := s.servicio.Crear(ctx, s.tienda, f("2026-11-10"), f("2026-11-20"), "Pisa"); !errors.Is(err, dominio.ErrSprintSuperpuesto) {
		t.Errorf("error = %v, se esperaba %v", err, dominio.ErrSprintSuperpuesto)
	}
	lista, _ := s.servicio.Listar(ctx, s.tienda)
	if len(lista) != 2 {
		t.Errorf("hay %d Sprints, se esperaban 2", len(lista))
	}
}

// spec 05 CA-7: crear un Sprint en un proyecto que no existe.
func TestSprintsProyectoInexistente(t *testing.T) {
	s := nuevosSprints(t)
	if _, err := s.servicio.Crear(context.Background(), 999, f("2026-11-01"), f("2026-11-14"), "Objetivo"); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("error = %v, se esperaba %v", err, dominio.ErrProyectoNoExiste)
	}
	if _, err := s.servicio.Listar(context.Background(), 999); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Listar: error = %v, se esperaba %v", err, dominio.ErrProyectoNoExiste)
	}
}

// spec 05 CA-2, RN-5, RN-7: el Sprint Goal se modifica solo con el Sprint abierto.
func TestSprintsModificarGoal(t *testing.T) {
	s := nuevosSprints(t)
	ctx := context.Background()
	sp := s.crearSprint(t, s.tienda, "2026-11-01", "2026-11-14")

	m, err := s.servicio.ModificarGoal(ctx, s.tienda, sp.ID, "Registrar y loguear clientes")
	if err != nil || m.Sprint.Goal != "Registrar y loguear clientes" {
		t.Fatalf("ModificarGoal = %+v, %v", m.Sprint, err)
	}
	if _, err := s.servicio.ModificarGoal(ctx, s.tienda, sp.ID, strings.Repeat("g", 301)); !errors.Is(err, dominio.ErrSprintGoalMuyLargo) {
		t.Errorf("error = %v, se esperaba %v", err, dominio.ErrSprintGoalMuyLargo)
	}
	visto, _ := s.servicio.Obtener(ctx, s.tienda, sp.ID)
	if visto.Sprint.Goal != "Registrar y loguear clientes" {
		t.Errorf("Goal = %q, un cambio inválido no debía modificarlo", visto.Sprint.Goal)
	}

	s.cerrar(t, sp)
	if _, err := s.servicio.ModificarGoal(ctx, s.tienda, sp.ID, "Otro"); !errors.Is(err, dominio.ErrSprintCerrado) {
		t.Errorf("Sprint cerrado: error = %v, se esperaba %v", err, dominio.ErrSprintCerrado)
	}
}

// spec 05 CA-3, CA-4, RN-8 a RN-10: asignar y quitar historias.
func TestSprintsAsignarYQuitar(t *testing.T) {
	s := nuevosSprints(t)
	ctx := context.Background()
	catalogo, _ := s.backlog.Crear(ctx, s.tienda, datos("Ver catálogo", 2))
	registrar, _ := s.backlog.Crear(ctx, s.tienda, datos("Registrar usuario", 1))
	terminada, _ := s.backlog.Crear(ctx, s.tienda, datos("Pagar", 3))
	s.backlog.Modificar(ctx, s.tienda, terminada.ID, datos("Pagar", 3), "Completada")
	ajena, _ := s.backlog.Crear(ctx, s.turnos, datos("Pedir turno", 1))
	s1 := s.crearSprint(t, s.tienda, "2026-11-01", "2026-11-14")
	s2 := s.crearSprint(t, s.tienda, "2026-11-15", "2026-11-28")

	if _, err := s.servicio.AsignarHistoria(ctx, s.tienda, s1.ID, catalogo.ID); err != nil {
		t.Fatalf("Asignar: %v", err)
	}
	d, err := s.servicio.AsignarHistoria(ctx, s.tienda, s1.ID, registrar.ID)
	if err != nil {
		t.Fatalf("Asignar: %v", err)
	}
	// Ver un Sprint: sus historias ordenadas como el backlog (prioridad 1 primero).
	if len(d.Historias) != 2 || d.Historias[0].ID != registrar.ID || d.Historias[1].ID != catalogo.ID {
		t.Errorf("historias = %+v, se esperaba [Registrar usuario, Ver catálogo]", d.Historias)
	}

	errores := []struct {
		nombre             string
		sprintID, historia int64
		err                error
	}{
		{"historia completada (RN-8)", s1.ID, terminada.ID, dominio.ErrHistoriaCompletada},
		{"ya está en este Sprint (RN-9, CA-4)", s1.ID, registrar.ID, dominio.ErrHistoriaYaAsignada},
		{"ya está en otro Sprint abierto (RN-9, CA-4)", s2.ID, registrar.ID, dominio.ErrHistoriaYaAsignada},
		{"historia de otro proyecto (RN-8, CA-7)", s1.ID, ajena.ID, dominio.ErrHistoriaNoExiste},
		{"historia inexistente (CA-7)", s1.ID, 999, dominio.ErrHistoriaNoExiste},
		{"Sprint inexistente (CA-7)", 999, catalogo.ID, dominio.ErrSprintNoExiste},
	}
	for _, c := range errores {
		t.Run(c.nombre, func(t *testing.T) {
			if _, err := s.servicio.AsignarHistoria(ctx, s.tienda, c.sprintID, c.historia); !errors.Is(err, c.err) {
				t.Errorf("error = %v, se esperaba %v", err, c.err)
			}
		})
	}

	// Quitar (RN-10): solo una historia que está en ese Sprint.
	if _, err := s.servicio.QuitarHistoria(ctx, s.tienda, s2.ID, registrar.ID); !errors.Is(err, dominio.ErrHistoriaNoAsignada) {
		t.Errorf("quitar de otro Sprint: error = %v, se esperaba %v", err, dominio.ErrHistoriaNoAsignada)
	}
	d, err = s.servicio.QuitarHistoria(ctx, s.tienda, s1.ID, registrar.ID)
	if err != nil || len(d.Historias) != 1 || d.Historias[0].ID != catalogo.ID {
		t.Fatalf("Quitar = %+v, %v; se esperaba solo Ver catálogo", d.Historias, err)
	}
	// Ya fuera del Sprint 1, se puede asignar al Sprint 2.
	if _, err := s.servicio.AsignarHistoria(ctx, s.tienda, s2.ID, registrar.ID); err != nil {
		t.Errorf("asignar después de quitar: %v", err)
	}
}

// spec 05 RN-7, RN-9: un Sprint cerrado no admite cambios, pero sus historias se pueden
// asignar a un Sprint nuevo.
func TestSprintsCerrado(t *testing.T) {
	s := nuevosSprints(t)
	ctx := context.Background()
	h, _ := s.backlog.Crear(ctx, s.tienda, datos("Registrar usuario", 1))
	otra, _ := s.backlog.Crear(ctx, s.tienda, datos("Ver catálogo", 2))
	s1 := s.crearSprint(t, s.tienda, "2026-11-01", "2026-11-14")
	s.servicio.AsignarHistoria(ctx, s.tienda, s1.ID, h.ID)
	s.cerrar(t, s1)

	if _, err := s.servicio.AsignarHistoria(ctx, s.tienda, s1.ID, otra.ID); !errors.Is(err, dominio.ErrSprintCerrado) {
		t.Errorf("asignar a cerrado: error = %v, se esperaba %v", err, dominio.ErrSprintCerrado)
	}
	if _, err := s.servicio.QuitarHistoria(ctx, s.tienda, s1.ID, h.ID); !errors.Is(err, dominio.ErrSprintCerrado) {
		t.Errorf("quitar de cerrado: error = %v, se esperaba %v", err, dominio.ErrSprintCerrado)
	}
	s2 := s.crearSprint(t, s.tienda, "2026-11-15", "2026-11-28")
	if _, err := s.servicio.AsignarHistoria(ctx, s.tienda, s2.ID, h.ID); err != nil {
		t.Errorf("asignar a un Sprint nuevo una historia de uno cerrado: %v", err)
	}
}

// spec 05 CA-3: listar por fecha de inicio y ver uno; un Sprint de otro proyecto no existe.
func TestSprintsListarYObtener(t *testing.T) {
	s := nuevosSprints(t)
	ctx := context.Background()
	if lista, err := s.servicio.Listar(ctx, s.tienda); err != nil || len(lista) != 0 {
		t.Fatalf("Listar sin Sprints = %v, %v; se esperaba lista vacía", lista, err)
	}
	segundo := s.crearSprint(t, s.tienda, "2026-11-15", "2026-11-28")
	primero := s.crearSprint(t, s.tienda, "2026-11-01", "2026-11-14")
	ajeno := s.crearSprint(t, s.turnos, "2026-11-01", "2026-11-14")

	lista, _ := s.servicio.Listar(ctx, s.tienda)
	if len(lista) != 2 || lista[0].ID != primero.ID || lista[1].ID != segundo.ID {
		t.Errorf("lista = %+v, se esperaba por fecha de inicio", lista)
	}
	if d, err := s.servicio.Obtener(ctx, s.tienda, primero.ID); err != nil || d.Sprint.ID != primero.ID || len(d.Historias) != 0 {
		t.Errorf("Obtener = %+v, %v", d, err)
	}
	if _, err := s.servicio.Obtener(ctx, s.tienda, ajeno.ID); !errors.Is(err, dominio.ErrSprintNoExiste) {
		t.Errorf("Sprint de otro proyecto: error = %v, se esperaba %v", err, dominio.ErrSprintNoExiste)
	}
}

// spec 05 CA-6, RN-11: no se elimina del backlog una historia de un Sprint abierto.
func TestBacklogNoEliminaHistoriaDeSprintAbierto(t *testing.T) {
	s := nuevosSprints(t)
	ctx := context.Background()
	h, _ := s.backlog.Crear(ctx, s.tienda, datos("Registrar usuario", 1))
	sp := s.crearSprint(t, s.tienda, "2026-11-01", "2026-11-14")
	s.servicio.AsignarHistoria(ctx, s.tienda, sp.ID, h.ID)

	if err := s.backlog.Eliminar(ctx, s.tienda, h.ID); !errors.Is(err, dominio.ErrHistoriaEnSprintAbierto) {
		t.Fatalf("error = %v, se esperaba %v", err, dominio.ErrHistoriaEnSprintAbierto)
	}
	if _, err := s.backlog.Obtener(ctx, s.tienda, h.ID); err != nil {
		t.Errorf("la historia no debía eliminarse: %v", err)
	}
	// Quitada del Sprint, se puede eliminar.
	s.servicio.QuitarHistoria(ctx, s.tienda, sp.ID, h.ID)
	if err := s.backlog.Eliminar(ctx, s.tienda, h.ID); err != nil {
		t.Errorf("eliminar después de quitarla del Sprint: %v", err)
	}
}
