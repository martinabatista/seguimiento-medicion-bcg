package postgres

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// spec 05 CA-1, CA-3, CA-7 contra PostgreSQL: guardar, asignar historias, cerrar, listar y obtener.
func TestRepositorioSprints(t *testing.T) {
	historias, tienda, turnos := historiasDePrueba(t)
	r := NuevoRepositorioSprints(historias.db)
	ctx := context.Background()
	h1, _ := historias.Guardar(ctx, dominio.Historia{ProyectoID: tienda, Titulo: "Registrar usuario", Descripcion: "Alta", Prioridad: 1, Estado: dominio.EstadoPendiente})
	h2, _ := historias.Guardar(ctx, dominio.Historia{ProyectoID: tienda, Titulo: "Ver catálogo", Descripcion: "Listado", Prioridad: 2, Estado: dominio.EstadoPendiente})

	s, err := r.Guardar(ctx, dominio.Sprint{ProyectoID: tienda, Numero: 1, FechaInicio: fecha("2026-11-01"), FechaFin: fecha("2026-11-14"), Goal: "Registrar clientes", Estado: dominio.SprintAbierto, Historias: []int64{}})
	if err != nil || s.ID == 0 {
		t.Fatalf("Guardar = %+v, %v; se esperaba un ID asignado", s, err)
	}

	// Actualizar: cambia el Goal, las historias y el estado.
	s.Goal = "Registrar y loguear"
	s.Historias = []int64{h1.ID, h2.ID}
	if err := r.Actualizar(ctx, s); err != nil {
		t.Fatalf("Actualizar: %v", err)
	}
	visto, err := r.Obtener(ctx, tienda, s.ID)
	slices.Sort(visto.Historias)
	if err != nil || visto.Goal != "Registrar y loguear" || !slices.Equal(visto.Historias, []int64{h1.ID, h2.ID}) ||
		!visto.FechaInicio.Equal(fecha("2026-11-01")) || visto.Numero != 1 || visto.Estado != dominio.SprintAbierto {
		t.Errorf("Obtener = %+v, %v", visto, err)
	}
	s.Historias = []int64{h2.ID}
	s.Estado = dominio.SprintCerrado
	r.Actualizar(ctx, s)
	visto, _ = r.Obtener(ctx, tienda, s.ID)
	if !slices.Equal(visto.Historias, []int64{h2.ID}) || visto.Estado != dominio.SprintCerrado {
		t.Errorf("después de quitar y cerrar = %+v", visto)
	}

	// Listar: solo los del proyecto; sin Sprints, lista vacía.
	lista, err := r.Listar(ctx, tienda)
	if err != nil || len(lista) != 1 || lista[0].ID != s.ID {
		t.Errorf("Listar = %+v, %v", lista, err)
	}
	if vacia, err := r.Listar(ctx, turnos); err != nil || vacia == nil || len(vacia) != 0 {
		t.Errorf("Listar sin Sprints = %#v, %v; se esperaba lista vacía", vacia, err)
	}

	// Un Sprint de otro proyecto o inexistente no existe.
	if _, err := r.Obtener(ctx, turnos, s.ID); !errors.Is(err, dominio.ErrSprintNoExiste) {
		t.Errorf("Obtener de otro proyecto: error = %v", err)
	}
	if err := r.Actualizar(ctx, dominio.Sprint{ID: 999, ProyectoID: tienda, Goal: "x", Estado: dominio.SprintAbierto}); !errors.Is(err, dominio.ErrSprintNoExiste) {
		t.Errorf("Actualizar inexistente: error = %v", err)
	}
}

// spec 06 RN-2, RN-6, RN-8: el Sprint guarda cuáles de sus historias se completaron en él.
func TestRepositorioSprintsCompletadas(t *testing.T) {
	historias, tienda, _ := historiasDePrueba(t)
	r := NuevoRepositorioSprints(historias.db)
	ctx := context.Background()
	h1, _ := historias.Guardar(ctx, dominio.Historia{ProyectoID: tienda, Titulo: "Registrar usuario", Descripcion: "Alta", Prioridad: 1, Estado: dominio.EstadoPendiente})
	h2, _ := historias.Guardar(ctx, dominio.Historia{ProyectoID: tienda, Titulo: "Ver catálogo", Descripcion: "Listado", Prioridad: 2, Estado: dominio.EstadoPendiente})
	s, _ := r.Guardar(ctx, dominio.Sprint{ProyectoID: tienda, Numero: 1, FechaInicio: fecha("2026-11-01"), FechaFin: fecha("2026-11-14"), Goal: "Objetivo", Estado: dominio.SprintAbierto})
	if s.Completadas == nil || len(s.Completadas) != 0 {
		t.Errorf("Guardar: completadas = %#v, se esperaba []", s.Completadas)
	}

	s.Historias = []int64{h1.ID, h2.ID}
	s.Completadas = []int64{h1.ID}
	s.Estado = dominio.SprintCerrado
	if err := r.Actualizar(ctx, s); err != nil {
		t.Fatalf("Actualizar: %v", err)
	}
	visto, err := r.Obtener(ctx, tienda, s.ID)
	slices.Sort(visto.Historias)
	if err != nil || !slices.Equal(visto.Historias, []int64{h1.ID, h2.ID}) || !slices.Equal(visto.Completadas, []int64{h1.ID}) {
		t.Errorf("Obtener = historias %v, completadas %v, %v; se esperaba [%d %d] y [%d]", visto.Historias, visto.Completadas, err, h1.ID, h2.ID, h1.ID)
	}
	lista, _ := r.Listar(ctx, tienda)
	if len(lista) != 1 || !slices.Equal(lista[0].Completadas, []int64{h1.ID}) {
		t.Errorf("Listar = %+v, se esperaban las completadas", lista)
	}
}

// spec 10 RN-4: el Sprint cerrado guarda los Story Points congelados de sus historias (NULL = sin estimar).
func TestRepositorioSprintsStoryPointsAlCerrar(t *testing.T) {
	historias, tienda, _ := historiasDePrueba(t)
	r := NuevoRepositorioSprints(historias.db)
	ctx := context.Background()
	h1, _ := historias.Guardar(ctx, dominio.Historia{ProyectoID: tienda, Titulo: "Registrar usuario", Descripcion: "Alta", Prioridad: 1, Estado: dominio.EstadoPendiente})
	h2, _ := historias.Guardar(ctx, dominio.Historia{ProyectoID: tienda, Titulo: "Ver catálogo", Descripcion: "Listado", Prioridad: 2, Estado: dominio.EstadoPendiente})
	s, _ := r.Guardar(ctx, dominio.Sprint{ProyectoID: tienda, Numero: 1, FechaInicio: fecha("2026-11-01"), FechaFin: fecha("2026-11-14"), Goal: "Objetivo", Estado: dominio.SprintAbierto})

	ocho := 8
	s.Historias = []int64{h1.ID, h2.ID}
	s.Estado = dominio.SprintCerrado
	s.StoryPointsAlCerrar = map[int64]*int{h1.ID: &ocho, h2.ID: nil}
	if err := r.Actualizar(ctx, s); err != nil {
		t.Fatalf("Actualizar: %v", err)
	}
	visto, err := r.Obtener(ctx, tienda, s.ID)
	if err != nil {
		t.Fatalf("Obtener: %v", err)
	}
	p1, ok1 := visto.StoryPointsAlCerrar[h1.ID]
	p2, ok2 := visto.StoryPointsAlCerrar[h2.ID]
	if !ok1 || p1 == nil || *p1 != 8 || !ok2 || p2 != nil {
		t.Errorf("congelados = %v, se esperaba %d → 8 y %d → sin estimar", visto.StoryPointsAlCerrar, h1.ID, h2.ID)
	}
	if lista, _ := r.Listar(ctx, tienda); len(lista) != 1 || lista[0].StoryPointsAlCerrar[h1.ID] == nil {
		t.Errorf("Listar = %+v, se esperaban los congelados", lista)
	}

	// Un Sprint abierto no tiene congelados.
	abierto, _ := r.Guardar(ctx, dominio.Sprint{ProyectoID: tienda, Numero: 2, FechaInicio: fecha("2026-11-15"), FechaFin: fecha("2026-11-28"), Goal: "Otro", Estado: dominio.SprintAbierto})
	abierto.Historias = []int64{h2.ID}
	r.Actualizar(ctx, abierto)
	if visto, _ := r.Obtener(ctx, tienda, abierto.ID); visto.StoryPointsAlCerrar != nil {
		t.Errorf("Sprint abierto: congelados = %v, se esperaba nil", visto.StoryPointsAlCerrar)
	}
}
