package dominio_test

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// conSprintYHistorias crea las historias "Registrar usuario" (1) y "Ver catálogo" (2) y un Sprint
// abierto con las dos asignadas.
func (s sprintsDePrueba) conSprintYHistorias(t *testing.T) (sp dominio.Sprint, registrar, catalogo dominio.Historia) {
	t.Helper()
	ctx := context.Background()
	registrar, _ = s.backlog.Crear(ctx, s.tienda, datos("Registrar usuario", 1))
	catalogo, _ = s.backlog.Crear(ctx, s.tienda, datos("Ver catálogo", 2))
	sp = s.crearSprint(t, s.tienda, "2026-11-01", "2026-11-14")
	for _, h := range []dominio.Historia{registrar, catalogo} {
		if _, err := s.servicio.AsignarHistoria(ctx, s.tienda, sp.ID, h.ID); err != nil {
			t.Fatalf("asignar historia de partida: %v", err)
		}
	}
	return sp, registrar, catalogo
}

// spec 06 CA-1, RN-1 a RN-3: completar y deshacer una historia del Sprint.
func TestSprintsCompletarYDeshacer(t *testing.T) {
	s := nuevosSprints(t)
	ctx := context.Background()
	sp, registrar, catalogo := s.conSprintYHistorias(t)

	d, err := s.servicio.CompletarHistoria(ctx, s.tienda, sp.ID, registrar.ID)
	if err != nil {
		t.Fatalf("Completar: %v", err)
	}
	if !d.Sprint.EstaCompletada(registrar.ID) || d.Sprint.EstaCompletada(catalogo.ID) {
		t.Errorf("completadas = %v, se esperaba solo %d", d.Sprint.Completadas, registrar.ID)
	}
	if h, _ := s.backlog.Obtener(ctx, s.tienda, registrar.ID); h.Estado != dominio.EstadoCompletada {
		t.Errorf("estado en el backlog = %q, se esperaba Completada (RN-2)", h.Estado)
	}

	errores := []struct {
		nombre string
		err    error
		accion func() error
	}{
		{"completar dos veces (RN-2)", dominio.ErrHistoriaCompletada, func() error {
			_, err := s.servicio.CompletarHistoria(ctx, s.tienda, sp.ID, registrar.ID)
			return err
		}},
		{"deshacer una no completada (RN-3)", dominio.ErrHistoriaNoCompletada, func() error {
			_, err := s.servicio.DeshacerHistoria(ctx, s.tienda, sp.ID, catalogo.ID)
			return err
		}},
		{"completar una historia que no existe (CA-6)", dominio.ErrHistoriaNoExiste, func() error {
			_, err := s.servicio.CompletarHistoria(ctx, s.tienda, sp.ID, 999)
			return err
		}},
		{"completar en un Sprint que no existe (CA-6)", dominio.ErrSprintNoExiste, func() error {
			_, err := s.servicio.CompletarHistoria(ctx, s.tienda, 999, registrar.ID)
			return err
		}},
	}
	for _, c := range errores {
		t.Run(c.nombre, func(t *testing.T) {
			if err := c.accion(); !errors.Is(err, c.err) {
				t.Errorf("error = %v, se esperaba %v", err, c.err)
			}
		})
	}

	d, err = s.servicio.DeshacerHistoria(ctx, s.tienda, sp.ID, registrar.ID)
	if err != nil || len(d.Sprint.Completadas) != 0 {
		t.Fatalf("Deshacer = %v, %v; se esperaba ninguna completada", d.Sprint.Completadas, err)
	}
	if h, _ := s.backlog.Obtener(ctx, s.tienda, registrar.ID); h.Estado != dominio.EstadoEnProgreso {
		t.Errorf("estado en el backlog = %q, se esperaba En progreso (RN-3)", h.Estado)
	}
}

// spec 06 CA-1, RN-1: solo se completa una historia asignada a ese Sprint.
func TestSprintsCompletarHistoriaNoAsignada(t *testing.T) {
	s := nuevosSprints(t)
	ctx := context.Background()
	h, _ := s.backlog.Crear(ctx, s.tienda, datos("Registrar usuario", 1))
	sp := s.crearSprint(t, s.tienda, "2026-11-01", "2026-11-14")
	if _, err := s.servicio.CompletarHistoria(ctx, s.tienda, sp.ID, h.ID); !errors.Is(err, dominio.ErrHistoriaNoAsignada) {
		t.Errorf("error = %v, se esperaba %v", err, dominio.ErrHistoriaNoAsignada)
	}
}

// spec 06 CA-2, CA-3, RN-5 a RN-7: cerrar conserva lo planificado y lo completado; después no
// admite cambios y las no completadas quedan libres.
func TestSprintsCerrar(t *testing.T) {
	s := nuevosSprints(t)
	ctx := context.Background()
	sp, registrar, catalogo := s.conSprintYHistorias(t)
	s.servicio.CompletarHistoria(ctx, s.tienda, sp.ID, registrar.ID)

	d, err := s.servicio.Cerrar(ctx, s.tienda, sp.ID)
	if err != nil {
		t.Fatalf("Cerrar: %v", err)
	}
	if d.Sprint.Estado != dominio.SprintCerrado || len(d.Historias) != 2 ||
		!d.Sprint.EstaCompletada(registrar.ID) || d.Sprint.EstaCompletada(catalogo.ID) {
		t.Errorf("sprint cerrado = %+v, se esperaba cerrado con las 2 historias y solo una completada", d.Sprint)
	}
	if h, _ := s.backlog.Obtener(ctx, s.tienda, catalogo.ID); h.Estado != dominio.EstadoPendiente {
		t.Errorf("la no completada quedó %q, se esperaba que mantenga su estado (RN-6)", h.Estado)
	}

	cerrado := []struct {
		nombre string
		accion func() error
	}{
		{"volver a cerrar", func() error { _, err := s.servicio.Cerrar(ctx, s.tienda, sp.ID); return err }},
		{"completar", func() error { _, err := s.servicio.CompletarHistoria(ctx, s.tienda, sp.ID, catalogo.ID); return err }},
		{"deshacer", func() error { _, err := s.servicio.DeshacerHistoria(ctx, s.tienda, sp.ID, registrar.ID); return err }},
		{"cambiar el Sprint Goal", func() error { _, err := s.servicio.ModificarGoal(ctx, s.tienda, sp.ID, "Otro"); return err }},
	}
	for _, c := range cerrado {
		t.Run(c.nombre, func(t *testing.T) {
			if err := c.accion(); !errors.Is(err, dominio.ErrSprintCerrado) {
				t.Errorf("error = %v, se esperaba %v (RN-7)", err, dominio.ErrSprintCerrado)
			}
		})
	}

	// RN-6: la no completada se puede asignar a un Sprint nuevo.
	nuevo := s.crearSprint(t, s.tienda, "2026-11-15", "2026-11-28")
	if _, err := s.servicio.AsignarHistoria(ctx, s.tienda, nuevo.ID, catalogo.ID); err != nil {
		t.Errorf("asignar la no completada a un Sprint nuevo: %v", err)
	}
	if _, err := s.servicio.Cerrar(ctx, s.tienda, 999); !errors.Is(err, dominio.ErrSprintNoExiste) {
		t.Errorf("cerrar inexistente: error = %v, se esperaba %v", err, dominio.ErrSprintNoExiste)
	}
}

// spec 06 CA-4, RN-8, RN-9: Sprints anteriores (cerrados), del más reciente al más antiguo; su
// registro no cambia aunque la historia cambie en el backlog.
func TestSprintsAnteriores(t *testing.T) {
	s := nuevosSprints(t)
	ctx := context.Background()
	if lista, err := s.servicio.Anteriores(ctx, s.tienda); err != nil || len(lista) != 0 {
		t.Fatalf("Anteriores sin cerrados = %v, %v; se esperaba lista vacía", lista, err)
	}
	primero, registrar, _ := s.conSprintYHistorias(t)
	s.servicio.CompletarHistoria(ctx, s.tienda, primero.ID, registrar.ID)
	s.servicio.Cerrar(ctx, s.tienda, primero.ID)
	segundo := s.crearSprint(t, s.tienda, "2026-11-15", "2026-11-28")
	s.servicio.Cerrar(ctx, s.tienda, segundo.ID)
	s.crearSprint(t, s.tienda, "2026-11-29", "2026-12-12") // abierto: no es anterior

	// RN-8: la historia vuelve a Pendiente en el backlog, pero en el Sprint 1 sigue completada.
	s.backlog.Modificar(ctx, s.tienda, registrar.ID, datos("Registrar usuario", 1), "Pendiente")

	lista, err := s.servicio.Anteriores(ctx, s.tienda)
	if err != nil {
		t.Fatalf("Anteriores: %v", err)
	}
	ids := []int64{}
	for _, d := range lista {
		ids = append(ids, d.Sprint.ID)
	}
	if !slices.Equal(ids, []int64{segundo.ID, primero.ID}) {
		t.Errorf("anteriores = %v, se esperaba [%d %d] (más reciente primero)", ids, segundo.ID, primero.ID)
	}
	if !lista[1].Sprint.EstaCompletada(registrar.ID) || len(lista[1].Historias) != 2 {
		t.Errorf("Sprint 1 = %+v, se esperaba con sus 2 historias y Registrar usuario completada", lista[1].Sprint)
	}
	if _, err := s.servicio.Anteriores(ctx, 999); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("proyecto inexistente: error = %v", err)
	}
}

// spec 06 CA-5, RN-4: el backlog no cambia "Completada" de una historia en un Sprint abierto.
func TestBacklogNoCompletaHistoriaDeSprintAbierto(t *testing.T) {
	s := nuevosSprints(t)
	ctx := context.Background()
	sp, registrar, catalogo := s.conSprintYHistorias(t)

	if _, err := s.backlog.Modificar(ctx, s.tienda, registrar.ID, datos("Registrar usuario", 1), "Completada"); !errors.Is(err, dominio.ErrCompletarDesdeSprint) {
		t.Errorf("pasar a Completada: error = %v, se esperaba %v", err, dominio.ErrCompletarDesdeSprint)
	}
	if _, err := s.backlog.Modificar(ctx, s.tienda, registrar.ID, datos("Registrar usuario", 1), "En progreso"); err != nil {
		t.Errorf("pasar a En progreso: %v (se puede desde el backlog)", err)
	}
	// Completada desde el Sprint, el backlog tampoco la puede sacar de Completada…
	s.servicio.CompletarHistoria(ctx, s.tienda, sp.ID, catalogo.ID)
	if _, err := s.backlog.Modificar(ctx, s.tienda, catalogo.ID, datos("Ver catálogo", 2), "Pendiente"); !errors.Is(err, dominio.ErrCompletarDesdeSprint) {
		t.Errorf("sacar de Completada: error = %v, se esperaba %v", err, dominio.ErrCompletarDesdeSprint)
	}
	// …pero sí cambiar otros datos dejándola Completada.
	if _, err := s.backlog.Modificar(ctx, s.tienda, catalogo.ID, datos("Ver catálogo de ropa", 2), "Completada"); err != nil {
		t.Errorf("cambiar el título sin tocar el estado: %v", err)
	}
	// Fuera de un Sprint abierto, el backlog sí la completa (spec 04).
	libre, _ := s.backlog.Crear(ctx, s.tienda, datos("Pagar", 3))
	if _, err := s.backlog.Modificar(ctx, s.tienda, libre.ID, datos("Pagar", 3), "Completada"); err != nil {
		t.Errorf("completar desde el backlog fuera de un Sprint: %v", err)
	}
}
