package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// Test de integración: necesita la base de docker compose y DATABASE_URL.
// Sin DATABASE_URL (por ejemplo, en la CI) se saltea.
func repoDePrueba(t *testing.T) *RepositorioProyectos {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL no está definida: se saltea el test de integración con PostgreSQL")
	}
	db, err := Abrir(url)
	if err != nil {
		t.Fatalf("Abrir: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	// CASCADE: también vacía las tablas que referencian a proyectos, directa o indirectamente
	// (historias, integrantes, Sprints, votaciones y esfuerzo; specs 02 a 09).
	if _, err := db.Exec("TRUNCATE proyectos RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("limpiar la tabla: %v", err)
	}
	return NuevoRepositorioProyectos(db)
}

func fecha(s string) time.Time {
	f, _ := time.Parse(dominio.FormatoFecha, s)
	return f
}

// spec 01 CA-1, CA-2, CA-6, CA-7 contra PostgreSQL.
func TestRepositorioProyectos(t *testing.T) {
	r := repoDePrueba(t)
	ctx := context.Background()

	a, err := r.Guardar(ctx, dominio.Proyecto{Nombre: "Tienda Online", Descripcion: "Venta de ropa", FechaInicio: fecha("2026-11-01"), FechaFin: fecha("2026-12-15")})
	if err != nil || a.ID == 0 {
		t.Fatalf("Guardar = %+v, %v; se esperaba un ID asignado", a, err)
	}
	b, _ := r.Guardar(ctx, dominio.Proyecto{Nombre: "App Turnos", Descripcion: "Turnos médicos", FechaInicio: fecha("2026-11-05"), FechaFin: fecha("2027-01-20")})

	visto, err := r.Obtener(ctx, a.ID)
	if err != nil || visto != a {
		t.Errorf("Obtener = %+v, %v; se esperaba %+v", visto, err, a)
	}

	a.Nombre = "Tienda Online v2"
	a.FechaFin = fecha("2027-01-31")
	if err := r.Actualizar(ctx, a); err != nil {
		t.Fatalf("Actualizar: %v", err)
	}
	lista, err := r.Listar(ctx)
	if err != nil || len(lista) != 2 || lista[0] != a || lista[1] != b {
		t.Errorf("Listar = %+v, %v; se esperaba [%+v %+v]", lista, err, a, b)
	}

	if _, err := r.Obtener(ctx, 999); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Obtener inexistente: error = %v", err)
	}
	if err := r.Actualizar(ctx, dominio.Proyecto{ID: 999, Nombre: "x", Descripcion: "x", FechaInicio: fecha("2026-11-01"), FechaFin: fecha("2026-11-02")}); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Actualizar inexistente: error = %v", err)
	}
}

// spec 01 CA-6: sin proyectos, la lista está vacía (no nil ni error).
func TestListarVacio(t *testing.T) {
	r := repoDePrueba(t)
	lista, err := r.Listar(context.Background())
	if err != nil || lista == nil || len(lista) != 0 {
		t.Errorf("Listar = %#v, %v; se esperaba una lista vacía", lista, err)
	}
}
