package postgres

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// spec 11 CA-1, CA-2, CA-3, CA-5, CA-7 contra PostgreSQL. Como repoDePrueba, se saltea si no está
// DATABASE_URL.
func TestRepositorioDefectos(t *testing.T) {
	historias, tienda, turnos := historiasDePrueba(t)
	ctx := context.Background()
	h, err := historias.Guardar(ctx, dominio.Historia{ProyectoID: tienda, Titulo: "Registrar usuario", Descripcion: "Alta", Prioridad: 1, Estado: dominio.EstadoPendiente})
	if err != nil {
		t.Fatalf("crear historia de partida: %v", err)
	}
	sprints := NuevoRepositorioSprints(historias.db)
	s1, err := sprints.Guardar(ctx, dominio.Sprint{ProyectoID: tienda, Numero: 1, FechaInicio: fecha("2026-11-01"), FechaFin: fecha("2026-11-14"), Goal: "Alta", Estado: dominio.SprintAbierto})
	if err != nil {
		t.Fatalf("crear Sprint de partida: %v", err)
	}
	s2, _ := sprints.Guardar(ctx, dominio.Sprint{ProyectoID: tienda, Numero: 2, FechaInicio: fecha("2026-11-16"), FechaFin: fecha("2026-11-28"), Goal: "Compra", Estado: dominio.SprintAbierto})
	r := NuevoRepositorioDefectos(historias.db)

	// Guardar: el ID lo genera PostgreSQL; arranca sin Sprint de resolución. Descripción en el
	// límite, con eñes (VARCHAR cuenta caracteres) y severidad con tilde.
	a, err := r.Guardar(ctx, dominio.Defecto{ProyectoID: tienda, HistoriaID: h.ID, Descripcion: strings.Repeat("ñ", 1000), Severidad: dominio.SeveridadCritica, Estado: dominio.DefectoAbierto, SprintDeteccionID: s1.ID})
	if err != nil || a.ID == 0 {
		t.Fatalf("Guardar = %+v, %v; se esperaba un ID asignado", a, err)
	}
	b, _ := r.Guardar(ctx, dominio.Defecto{ProyectoID: tienda, HistoriaID: h.ID, Descripcion: "Otra", Severidad: dominio.SeveridadBaja, Estado: dominio.DefectoAbierto, SprintDeteccionID: s2.ID})

	visto, err := r.Obtener(ctx, tienda, a.ID)
	if err != nil || visto != a {
		t.Errorf("Obtener = %+v, %v; se esperaba %+v", visto, err, a)
	}

	// Actualizar: resolver guarda el Sprint de resolución; reabrir lo vuelve a NULL.
	a.Estado, a.SprintResolucionID = dominio.DefectoResuelto, s2.ID
	if err := r.Actualizar(ctx, a); err != nil {
		t.Fatalf("Actualizar a Resuelto: %v", err)
	}
	if visto, _ := r.Obtener(ctx, tienda, a.ID); visto.Estado != dominio.DefectoResuelto || visto.SprintResolucionID != s2.ID {
		t.Errorf("después de resolver = %+v", visto)
	}
	a.Estado, a.SprintResolucionID = dominio.DefectoEnCorreccion, 0
	if err := r.Actualizar(ctx, a); err != nil {
		t.Fatalf("Actualizar a En corrección: %v", err)
	}
	if visto, _ := r.Obtener(ctx, tienda, a.ID); visto.Estado != dominio.DefectoEnCorreccion || visto.SprintResolucionID != 0 {
		t.Errorf("después de reabrir = %+v", visto)
	}

	// Listar en orden de registro; cada proyecto ve los suyos (RN-12).
	lista, err := r.Listar(ctx, tienda)
	if err != nil || len(lista) != 2 || lista[0].ID != a.ID || lista[1].ID != b.ID {
		t.Fatalf("Listar = %+v, %v; se esperaban %d y %d en ese orden", lista, err, a.ID, b.ID)
	}
	if lista, _ := r.Listar(ctx, turnos); len(lista) != 0 {
		t.Errorf("Listar de otro proyecto = %+v, se esperaba vacía", lista)
	}

	// Un defecto de otro proyecto o inexistente no existe.
	if _, err := r.Obtener(ctx, turnos, a.ID); !errors.Is(err, dominio.ErrDefectoNoExiste) {
		t.Errorf("Obtener desde otro proyecto: error = %v", err)
	}
	if err := r.Actualizar(ctx, dominio.Defecto{ID: 999, ProyectoID: tienda, Estado: dominio.DefectoAbierto}); !errors.Is(err, dominio.ErrDefectoNoExiste) {
		t.Errorf("Actualizar inexistente: error = %v", err)
	}

	// RN-14: la base tampoco deja borrar una historia con defectos (sin ON DELETE CASCADE).
	if err := historias.Eliminar(ctx, tienda, h.ID); err == nil {
		t.Error("se eliminó una historia con defectos")
	}
}
