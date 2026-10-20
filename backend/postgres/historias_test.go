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

// mismaHistoria compara dos historias, incluido el valor de los Story Points (que es un puntero).
func mismaHistoria(a, b dominio.Historia) bool {
	spIguales := (a.StoryPoints == nil && b.StoryPoints == nil) ||
		(a.StoryPoints != nil && b.StoryPoints != nil && *a.StoryPoints == *b.StoryPoints)
	a.StoryPoints, b.StoryPoints = nil, nil
	return a == b && spIguales
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
