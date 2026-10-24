package postgres

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// esfuerzoDePrueba arma el repositorio de esfuerzo sobre la base vacía, con una historia y una
// integrante en la tienda, y una historia en App Turnos. Como repoDePrueba, se saltea si no está
// DATABASE_URL.
func esfuerzoDePrueba(t *testing.T) (r *RepositorioEsfuerzo, historias *RepositorioHistorias, tienda, historia, otraHistoria, ana int64) {
	t.Helper()
	historias, tienda, turnos := historiasDePrueba(t)
	ctx := context.Background()
	h, err := historias.Guardar(ctx, dominio.Historia{ProyectoID: tienda, Titulo: "Registrar usuario", Descripcion: "Alta", Prioridad: 1, Estado: dominio.EstadoPendiente})
	if err != nil {
		t.Fatalf("crear historia de partida: %v", err)
	}
	otra, _ := historias.Guardar(ctx, dominio.Historia{ProyectoID: turnos, Titulo: "Pedir turno", Descripcion: "Reserva", Prioridad: 1, Estado: dominio.EstadoPendiente})
	i, err := NuevoRepositorioIntegrantes(historias.db).Guardar(ctx, dominio.Integrante{ProyectoID: tienda, Nombre: "Ana Pérez", Email: "ana@mail.com", Rol: dominio.RolDeveloper})
	if err != nil {
		t.Fatalf("crear integrante de partida: %v", err)
	}
	return NuevoRepositorioEsfuerzo(historias.db), historias, tienda, h.ID, otra.ID, i.ID
}

// spec 09 CA-1, CA-4, CA-6, CA-8 contra PostgreSQL.
func TestRepositorioEsfuerzo(t *testing.T) {
	r, _, tienda, historia, otraHistoria, ana := esfuerzoDePrueba(t)
	ctx := context.Background()

	// Guardar: el ID lo genera PostgreSQL; las horas se guardan en centésimas exactas.
	a, err := r.Guardar(ctx, dominio.RegistroEsfuerzo{ProyectoID: tienda, HistoriaID: historia, IntegranteID: ana, Fecha: fecha("2026-11-18"), Actividad: "Diseño", Centesimas: 10})
	if err != nil || a.ID == 0 {
		t.Fatalf("Guardar = %+v, %v; se esperaba un ID asignado", a, err)
	}
	// Actividad en el límite, con eñes: VARCHAR cuenta caracteres, no bytes. 24 horas.
	b, err := r.Guardar(ctx, dominio.RegistroEsfuerzo{ProyectoID: tienda, HistoriaID: historia, IntegranteID: ana, Fecha: fecha("2026-11-19"), Actividad: strings.Repeat("ñ", 300), Centesimas: 2400})
	if err != nil {
		t.Fatalf("Guardar con datos en el límite: %v", err)
	}

	lista, err := r.Listar(ctx, tienda, historia)
	if err != nil || len(lista) != 2 {
		t.Fatalf("Listar = %+v, %v; se esperaban 2 registros", lista, err)
	}
	for _, got := range lista {
		esperado := a
		if got.ID == b.ID {
			esperado = b
		}
		if got.ID != esperado.ID || got.ProyectoID != tienda || got.HistoriaID != historia || got.IntegranteID != ana ||
			!got.Fecha.Equal(esperado.Fecha) || got.Actividad != esperado.Actividad || got.Centesimas != esperado.Centesimas {
			t.Errorf("Listar devolvió %+v, se esperaba %+v", got, esperado)
		}
	}

	// Un registro de otra historia o de otro proyecto no existe (RN-8).
	if err := r.Eliminar(ctx, tienda, otraHistoria, a.ID); !errors.Is(err, dominio.ErrRegistroEsfuerzoNoExiste) {
		t.Errorf("Eliminar desde otra historia: error = %v", err)
	}
	if err := r.Eliminar(ctx, tienda, historia, 999); !errors.Is(err, dominio.ErrRegistroEsfuerzoNoExiste) {
		t.Errorf("Eliminar inexistente: error = %v", err)
	}

	// Eliminar: deja de estar en la lista; otra vez, no existe.
	if err := r.Eliminar(ctx, tienda, historia, a.ID); err != nil {
		t.Fatalf("Eliminar: %v", err)
	}
	if lista, _ := r.Listar(ctx, tienda, historia); len(lista) != 1 || lista[0].ID != b.ID {
		t.Errorf("después de eliminar = %+v, se esperaba solo el registro %d", lista, b.ID)
	}
	if err := r.Eliminar(ctx, tienda, historia, a.ID); !errors.Is(err, dominio.ErrRegistroEsfuerzoNoExiste) {
		t.Errorf("Eliminar otra vez: error = %v", err)
	}
}

// spec 09 RN-10: la base tampoco deja eliminar una historia con esfuerzo registrado (aunque no pase
// por el dominio): las horas reales no se pierden.
func TestHistoriaConEsfuerzoNoSeBorraEnLaBase(t *testing.T) {
	r, historias, tienda, historia, _, ana := esfuerzoDePrueba(t)
	ctx := context.Background()
	if _, err := r.Guardar(ctx, dominio.RegistroEsfuerzo{ProyectoID: tienda, HistoriaID: historia, IntegranteID: ana, Fecha: fecha("2026-11-18"), Actividad: "Backend", Centesimas: 200}); err != nil {
		t.Fatalf("Guardar de partida: %v", err)
	}
	if err := historias.Eliminar(ctx, tienda, historia); err == nil {
		t.Error("se eliminó una historia con esfuerzo registrado; la base debía impedirlo")
	}
}

// spec 09 CA-4: una historia sin registros tiene la lista vacía (no nil ni error).
func TestListarEsfuerzoVacio(t *testing.T) {
	r, _, tienda, historia, _, _ := esfuerzoDePrueba(t)
	lista, err := r.Listar(context.Background(), tienda, historia)
	if err != nil || lista == nil || len(lista) != 0 {
		t.Errorf("Listar = %#v, %v; se esperaba una lista vacía", lista, err)
	}
}
