package dominio_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/memoria"
)

func f(s string) time.Time {
	t, _ := time.Parse(dominio.FormatoFecha, s)
	return t
}

func nuevoServicio() *dominio.ServicioProyectos {
	return dominio.NuevoServicioProyectos(memoria.NuevoRepositorioProyectos())
}

// spec 01 CA-1, RN-6: al crear, el sistema asigna el identificador.
func TestCrearAsignaIdentificador(t *testing.T) {
	s := nuevoServicio()
	ctx := context.Background()

	p1, err := s.Crear(ctx, "Tienda Online", "Venta de ropa", f("2026-11-01"), f("2026-12-15"))
	if err != nil {
		t.Fatalf("Crear: %v", err)
	}
	p2, err := s.Crear(ctx, "Tienda Online", "Otra tienda", f("2027-01-01"), f("2027-02-01"))
	if err != nil {
		t.Fatalf("Crear con nombre repetido (spec 01 CA-4): %v", err)
	}
	if p1.ID == 0 || p2.ID == 0 || p1.ID == p2.ID {
		t.Errorf("IDs = %d y %d, se esperaban distintos de cero y entre sí", p1.ID, p2.ID)
	}
}

// spec 01 CA-1: si un dato es inválido, no se registra nada.
func TestCrearInvalidoNoRegistra(t *testing.T) {
	s := nuevoServicio()
	ctx := context.Background()

	_, err := s.Crear(ctx, "", "Venta de ropa", f("2026-11-01"), f("2026-12-15"))
	if !errors.Is(err, dominio.ErrNombreObligatorio) {
		t.Fatalf("error = %v, se esperaba %v", err, dominio.ErrNombreObligatorio)
	}
	lista, _ := s.Listar(ctx)
	if len(lista) != 0 {
		t.Errorf("hay %d proyectos, no se debía registrar ninguno", len(lista))
	}
}

// spec 01 CA-2: modificar los datos de un proyecto existente.
func TestModificar(t *testing.T) {
	s := nuevoServicio()
	ctx := context.Background()
	p, _ := s.Crear(ctx, "Tienda Online", "Venta de ropa", f("2026-11-01"), f("2026-12-15"))

	m, err := s.Modificar(ctx, p.ID, "Tienda Online v2", "Venta de ropa", f("2026-11-10"), f("2027-01-31"))
	if err != nil {
		t.Fatalf("Modificar: %v", err)
	}
	visto, _ := s.Obtener(ctx, p.ID)
	if visto != m || visto.Nombre != "Tienda Online v2" || !visto.FechaFin.Equal(f("2027-01-31")) || visto.ID != p.ID {
		t.Errorf("proyecto = %+v, no quedó modificado", visto)
	}
}

// spec 01 RN-5, CA-2: una modificación inválida deja el proyecto como estaba.
func TestModificarInvalidoNoCambia(t *testing.T) {
	s := nuevoServicio()
	ctx := context.Background()
	p, _ := s.Crear(ctx, "Tienda Online", "Venta de ropa", f("2026-11-01"), f("2026-12-15"))

	_, err := s.Modificar(ctx, p.ID, "Tienda Online", "Venta de ropa", f("2026-11-10"), f("2026-11-05"))
	if !errors.Is(err, dominio.ErrFinNoPosteriorAInicio) {
		t.Fatalf("error = %v, se esperaba %v", err, dominio.ErrFinNoPosteriorAInicio)
	}
	visto, _ := s.Obtener(ctx, p.ID)
	if visto != p {
		t.Errorf("proyecto = %+v, se esperaba sin cambios %+v", visto, p)
	}
}

// spec 01 CA-6: listar en orden de creación y ver uno.
func TestListarYObtener(t *testing.T) {
	s := nuevoServicio()
	ctx := context.Background()

	vacia, err := s.Listar(ctx)
	if err != nil || len(vacia) != 0 {
		t.Fatalf("Listar sin proyectos = %v, %v; se esperaba lista vacía", vacia, err)
	}
	a, _ := s.Crear(ctx, "Tienda Online", "Venta de ropa", f("2026-11-01"), f("2026-12-15"))
	b, _ := s.Crear(ctx, "App Turnos", "Turnos médicos", f("2026-11-05"), f("2027-01-20"))

	lista, _ := s.Listar(ctx)
	if len(lista) != 2 || lista[0] != a || lista[1] != b {
		t.Errorf("lista = %+v, se esperaba [%+v %+v] en orden de creación", lista, a, b)
	}
	visto, err := s.Obtener(ctx, b.ID)
	if err != nil || visto != b {
		t.Errorf("Obtener = %+v, %v; se esperaba %+v", visto, err, b)
	}
}

// spec 01 CA-7: ver o modificar un proyecto que no existe.
func TestProyectoInexistente(t *testing.T) {
	s := nuevoServicio()
	ctx := context.Background()

	if _, err := s.Obtener(ctx, 999); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Obtener: error = %v, se esperaba %v", err, dominio.ErrProyectoNoExiste)
	}
	_, err := s.Modificar(ctx, 999, "Fantasma", "Venta de ropa", f("2026-11-01"), f("2026-12-15"))
	if !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Modificar: error = %v, se esperaba %v", err, dominio.ErrProyectoNoExiste)
	}
}
