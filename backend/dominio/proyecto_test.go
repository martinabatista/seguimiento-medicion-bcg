package dominio

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func fecha(t *testing.T, s string) time.Time {
	t.Helper()
	f, err := time.Parse("2006-01-02", s)
	if err != nil {
		t.Fatalf("fecha de prueba inválida %q: %v", s, err)
	}
	return f
}

// spec 01 RN-1 a RN-4 (CA-1, CA-3, CA-4, CA-5)
func TestNuevoProyecto(t *testing.T) {
	inicio := fecha(t, "2026-11-01")
	casos := []struct {
		nombre      string
		proyNombre  string
		descripcion string
		inicio      time.Time
		fin         time.Time
		err         error
	}{
		{"datos válidos (CA-1)", "Tienda Online", "Venta de ropa", inicio, fecha(t, "2026-12-15"), nil},
		{"dura un solo día (CA-1)", "Hackatón", "Evento de un día", inicio, fecha(t, "2026-11-02"), nil},
		{"nombre vacío (spec 01 RN-1, CA-4)", "", "Venta de ropa", inicio, fecha(t, "2026-12-15"), ErrNombreObligatorio},
		{"nombre con solo espacios (spec 01 RN-1, CA-4)", "   ", "Venta de ropa", inicio, fecha(t, "2026-12-15"), ErrNombreObligatorio},
		{"nombre muy largo (spec 01 RN-1, CA-4)", strings.Repeat("a", 300), "Venta de ropa", inicio, fecha(t, "2026-12-15"), nil},
		{"descripción vacía (spec 01 RN-2, CA-5)", "Tienda Online", "", inicio, fecha(t, "2026-12-15"), ErrDescripcionObligatoria},
		{"descripción con solo espacios (spec 01 RN-2, CA-5)", "Tienda Online", "  ", inicio, fecha(t, "2026-12-15"), ErrDescripcionObligatoria},
		{"descripción de 500 caracteres (spec 01 RN-2, CA-5)", "Tienda Online", strings.Repeat("d", 500), inicio, fecha(t, "2026-12-15"), nil},
		{"descripción de 501 caracteres (spec 01 RN-2, CA-5)", "Tienda Online", strings.Repeat("d", 501), inicio, fecha(t, "2026-12-15"), ErrDescripcionMuyLarga},
		{"500 caracteres con tildes (spec 01 RN-2, CA-5)", "Tienda Online", strings.Repeat("á", 500), inicio, fecha(t, "2026-12-15"), nil},
		{"falta la fecha de inicio (spec 01 RN-3, CA-3)", "Tienda Online", "Venta de ropa", time.Time{}, fecha(t, "2026-12-15"), ErrFechasObligatorias},
		{"falta la fecha de fin (spec 01 RN-3, CA-3)", "Tienda Online", "Venta de ropa", inicio, time.Time{}, ErrFechasObligatorias},
		{"fin igual a inicio (spec 01 RN-4, CA-3)", "Tienda Online", "Venta de ropa", inicio, inicio, ErrFinNoPosteriorAInicio},
		{"fin anterior a inicio (spec 01 RN-4, CA-3)", "Tienda Online", "Venta de ropa", fecha(t, "2026-11-10"), inicio, ErrFinNoPosteriorAInicio},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			p, err := NuevoProyecto(c.proyNombre, c.descripcion, c.inicio, c.fin)
			if !errors.Is(err, c.err) {
				t.Fatalf("error = %v, se esperaba %v", err, c.err)
			}
			if err != nil {
				return
			}
			if p.Nombre != c.proyNombre || p.Descripcion != c.descripcion ||
				!p.FechaInicio.Equal(c.inicio) || !p.FechaFin.Equal(c.fin) {
				t.Errorf("proyecto = %+v, no conserva los datos recibidos", p)
			}
			if p.ID != 0 {
				t.Errorf("ID = %d, el identificador lo asigna el sistema al guardar (spec 01 RN-6)", p.ID)
			}
		})
	}
}

// spec 01 RN-3 (CA-3): las fechas llegan como texto AAAA-MM-DD
func TestParsearFecha(t *testing.T) {
	casos := []struct {
		nombre string
		texto  string
		vacia  bool
		err    error
	}{
		{"fecha válida", "2026-11-01", false, nil},
		{"vacía: se trata como fecha faltante", "", true, nil},
		{"día inexistente", "2026-02-30", false, ErrFechaInvalida},
		{"formato día/mes/año", "01/11/2026", false, ErrFechaInvalida},
		{"texto cualquiera", "mañana", false, ErrFechaInvalida},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			f, err := ParsearFecha(c.texto)
			if !errors.Is(err, c.err) {
				t.Fatalf("error = %v, se esperaba %v", err, c.err)
			}
			if err == nil && f.IsZero() != c.vacia {
				t.Errorf("fecha = %v, vacía = %v; se esperaba vacía = %v", f, f.IsZero(), c.vacia)
			}
		})
	}
}

// Los errores de validación se distinguen de "no existe" (400 vs 404 en la API).
func TestErroresDeValidacionSonDatoInvalido(t *testing.T) {
	invalidos := []error{
		ErrNombreObligatorio, ErrDescripcionObligatoria, ErrDescripcionMuyLarga,
		ErrFechasObligatorias, ErrFechaInvalida, ErrFinNoPosteriorAInicio,
	}
	for _, e := range invalidos {
		var d ErrDatoInvalido
		if !errors.As(e, &d) {
			t.Errorf("%q debería ser un ErrDatoInvalido", e)
		}
	}
	var d ErrDatoInvalido
	if errors.As(ErrProyectoNoExiste, &d) {
		t.Error("ErrProyectoNoExiste no es un dato inválido")
	}
}
