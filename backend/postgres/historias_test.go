package postgres

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// historiasDePrueba arma el repositorio de historias sobre la base vacía, con dos proyectos creados.
// Como repoDePrueba, se saltea si no está DATABASE_URL.
func historiasDePrueba(t *testing.T) (r *RepositorioHistorias, tienda, turnos int64) {
	t.Helper()
	proyectos := repoDePrueba(t)
	ctx := context.Background()
	a, err := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "Tienda Online", Descripcion: "Venta de ropa", FechaInicio: fecha("2026-11-01"), FechaFin: fecha("2026-12-15")})
	if err != nil {
		t.Fatalf("crear proyecto de partida: %v", err)
	}
	b, _ := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "App Turnos", Descripcion: "Turnos médicos", FechaInicio: fecha("2026-11-05"), FechaFin: fecha("2027-01-20")})
	return NuevoRepositorioHistorias(proyectos.db), a.ID, b.ID
}

// mismaHistoria compara dos historias, incluido el valor de los Story Points y de las horas
// estimadas (que son punteros).
func mismaHistoria(a, b dominio.Historia) bool {
	iguales := mismoValor(a.StoryPoints, b.StoryPoints) && mismoValor(a.HorasEstimadas, b.HorasEstimadas)
	a.StoryPoints, b.StoryPoints = nil, nil
	a.HorasEstimadas, b.HorasEstimadas = nil, nil
	return a == b && iguales
}

func mismoValor(a, b *int) bool {
	return (a == nil && b == nil) || (a != nil && b != nil && *a == *b)
}

// spec 04 CA-1 a CA-4, CA-7 contra PostgreSQL.
func TestRepositorioHistorias(t *testing.T) {
	r, tienda, turnos := historiasDePrueba(t)
	ctx := context.Background()

	// Guardar: el ID lo genera PostgreSQL y los Story Points quedan sin estimar (NULL).
	a, err := r.Guardar(ctx, dominio.Historia{ProyectoID: tienda, Titulo: "Registrar usuario", Descripcion: "Alta de clientes", Prioridad: 1, Criterios: "Valida el email", Estado: dominio.EstadoPendiente})
	if err != nil || a.ID == 0 {
		t.Fatalf("Guardar = %+v, %v; se esperaba un ID asignado", a, err)
	}
	// Textos en el límite, con eñes: VARCHAR cuenta caracteres, no bytes. Criterios vacíos.
	b, err := r.Guardar(ctx, dominio.Historia{ProyectoID: tienda, Titulo: strings.Repeat("ñ", 150), Descripcion: strings.Repeat("á", 1000), Prioridad: 2, Estado: dominio.EstadoPendiente})
	if err != nil {
		t.Fatalf("Guardar con textos en el límite: %v", err)
	}
	c, _ := r.Guardar(ctx, dominio.Historia{ProyectoID: turnos, Titulo: "Pedir turno", Descripcion: "Reserva", Prioridad: 1, Estado: dominio.EstadoPendiente})

	visto, err := r.Obtener(ctx, tienda, a.ID)
	if err != nil || !mismaHistoria(visto, a) || visto.StoryPoints != nil {
		t.Errorf("Obtener = %+v, %v; se esperaba %+v sin estimar", visto, err, a)
	}

	// Actualizar: reemplaza los datos, el estado y guarda los Story Points.
	sp := 5
	a.Titulo, a.Prioridad, a.Estado, a.StoryPoints = "Registrar cliente", 3, dominio.EstadoEnProgreso, &sp
	if err := r.Actualizar(ctx, a); err != nil {
		t.Fatalf("Actualizar: %v", err)
	}
	visto, _ = r.Obtener(ctx, tienda, a.ID)
	if !mismaHistoria(visto, a) {
		t.Errorf("después de Actualizar = %+v, se esperaba %+v", visto, a)
	}

	// Listar: solo las del proyecto.
	lista, err := r.Listar(ctx, tienda)
	if err != nil || len(lista) != 2 {
		t.Fatalf("Listar = %+v, %v; se esperaban las 2 historias de la tienda", lista, err)
	}
	for _, h := range lista {
		if h.ProyectoID != tienda || (h.ID != a.ID && h.ID != b.ID) {
			t.Errorf("Listar devolvió %+v, que no es del proyecto %d", h, tienda)
		}
	}

	// Una historia de otro proyecto se trata como inexistente (RN-10).
	if _, err := r.Obtener(ctx, tienda, c.ID); !errors.Is(err, dominio.ErrHistoriaNoExiste) {
		t.Errorf("Obtener de otro proyecto: error = %v", err)
	}
	otra := c
	otra.ProyectoID = tienda
	if err := r.Actualizar(ctx, otra); !errors.Is(err, dominio.ErrHistoriaNoExiste) {
		t.Errorf("Actualizar de otro proyecto: error = %v", err)
	}
	if err := r.Eliminar(ctx, tienda, c.ID); !errors.Is(err, dominio.ErrHistoriaNoExiste) {
		t.Errorf("Eliminar de otro proyecto: error = %v", err)
	}

	// Eliminar: deja de existir; otra vez, no existe.
	if err := r.Eliminar(ctx, tienda, a.ID); err != nil {
		t.Fatalf("Eliminar: %v", err)
	}
	if _, err := r.Obtener(ctx, tienda, a.ID); !errors.Is(err, dominio.ErrHistoriaNoExiste) {
		t.Errorf("Obtener eliminada: error = %v", err)
	}
	if err := r.Eliminar(ctx, tienda, a.ID); !errors.Is(err, dominio.ErrHistoriaNoExiste) {
		t.Errorf("Eliminar otra vez: error = %v", err)
	}
}

// spec 04 CA-4: un proyecto sin historias tiene el backlog vacío (no nil ni error).
func TestListarHistoriasVacio(t *testing.T) {
	r, tienda, _ := historiasDePrueba(t)
	lista, err := r.Listar(context.Background(), tienda)
	if err != nil || lista == nil || len(lista) != 0 {
		t.Errorf("Listar = %#v, %v; se esperaba una lista vacía", lista, err)
	}
}

// spec 07 CA-1, CA-2: los Story Points se guardan, se cambian y se quitan (NULL = sin estimar).
func TestStoryPointsDeHistoria(t *testing.T) {
	r, tienda, _ := historiasDePrueba(t)
	ctx := context.Background()
	h, err := r.Guardar(ctx, dominio.Historia{ProyectoID: tienda, Titulo: "Registrar usuario", Descripcion: "d", Prioridad: 1, Estado: dominio.EstadoPendiente})
	if err != nil {
		t.Fatalf("Guardar: %v", err)
	}
	for _, sp := range []int{21, 1} {
		h.StoryPoints = &sp
		if err := r.Actualizar(ctx, h); err != nil {
			t.Fatalf("Actualizar con %d: %v", sp, err)
		}
		if visto, _ := r.Obtener(ctx, tienda, h.ID); visto.StoryPoints == nil || *visto.StoryPoints != sp {
			t.Errorf("Story Points = %v, se esperaba %d", visto.StoryPoints, sp)
		}
	}
	h.StoryPoints = nil
	if err := r.Actualizar(ctx, h); err != nil {
		t.Fatalf("Actualizar sin estimar: %v", err)
	}
	if visto, _ := r.Obtener(ctx, tienda, h.ID); visto.StoryPoints != nil {
		t.Errorf("Story Points = %d, se esperaba sin estimar", *visto.StoryPoints)
	}
}

// spec 12 CA-1: las horas estimadas se guardan en centésimas, se cambian y se quitan (NULL = sin
// horas estimadas), también al crear y sin tocar los Story Points.
func TestHorasEstimadasDeHistoria(t *testing.T) {
	r, tienda, _ := historiasDePrueba(t)
	ctx := context.Background()
	ocho, sp := 800, 5
	h, err := r.Guardar(ctx, dominio.Historia{ProyectoID: tienda, Titulo: "Registrar usuario", Descripcion: "d", Prioridad: 1, Estado: dominio.EstadoPendiente, StoryPoints: &sp, HorasEstimadas: &ocho})
	if err != nil {
		t.Fatalf("Guardar: %v", err)
	}
	if visto, _ := r.Obtener(ctx, tienda, h.ID); visto.HorasEstimadas == nil || *visto.HorasEstimadas != 800 || visto.StoryPoints == nil || *visto.StoryPoints != 5 {
		t.Errorf("al guardar: historia = %+v, se esperaban 800 centésimas y 5 Story Points", visto)
	}
	for _, c := range []int{650, 1, 999999999} { // el máximo que acepta el dominio: 9.999.999,99 h
		h.HorasEstimadas = &c
		if err := r.Actualizar(ctx, h); err != nil {
			t.Fatalf("Actualizar con %d: %v", c, err)
		}
		if lista, _ := r.Listar(ctx, tienda); len(lista) != 1 || lista[0].HorasEstimadas == nil || *lista[0].HorasEstimadas != c {
			t.Errorf("Listar = %+v, se esperaban %d centésimas", lista, c)
		}
	}
	h.HorasEstimadas = nil
	if err := r.Actualizar(ctx, h); err != nil {
		t.Fatalf("Actualizar sin horas estimadas: %v", err)
	}
	if visto, _ := r.Obtener(ctx, tienda, h.ID); visto.HorasEstimadas != nil {
		t.Errorf("horas estimadas = %d, se esperaba sin horas estimadas", *visto.HorasEstimadas)
	}
}
