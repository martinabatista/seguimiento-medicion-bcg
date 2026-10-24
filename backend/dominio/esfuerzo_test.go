package dominio

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// spec 09 RN-5 (CA-3): las horas llegan como texto (el número del JSON) y se guardan en centésimas
// enteras, para que las sumas sean exactas.
func TestParsearHoras(t *testing.T) {
	casos := []struct {
		texto      string
		centesimas int
		err        error
	}{
		{"1", 100, nil},
		{"2.5", 250, nil},
		{"0.25", 25, nil},
		{"0.01", 1, nil}, // el mínimo
		{"24", 2400, nil},
		{"24.00", 2400, nil},
		{"7.10", 710, nil},
		{"", 0, ErrHorasObligatorias},
		{"0", 0, ErrHorasInvalidas},
		{"0.00", 0, ErrHorasInvalidas},
		{"24.01", 0, ErrHorasInvalidas},
		{"25", 0, ErrHorasInvalidas},
		{"-2", 0, ErrHorasInvalidas},
		{"1.125", 0, ErrHorasInvalidas}, // más de 2 decimales
		{"1e1", 0, ErrHorasInvalidas},
		{"1.", 0, ErrHorasInvalidas},
		{".5", 0, ErrHorasInvalidas},
		{"dos", 0, ErrHorasInvalidas},
		{"100000000000000000000", 0, ErrHorasInvalidas},
	}
	for _, c := range casos {
		t.Run(c.texto, func(t *testing.T) {
			got, err := ParsearHoras(c.texto)
			if !errors.Is(err, c.err) {
				t.Fatalf("error = %v, se esperaba %v", err, c.err)
			}
			if got != c.centesimas {
				t.Errorf("centésimas = %d, se esperaba %d", got, c.centesimas)
			}
		})
	}
}

// spec 09 RN-6: las horas se muestran con hasta 2 decimales, sin ceros de más.
func TestHorasATexto(t *testing.T) {
	casos := map[int]string{0: "0", 1: "0.01", 30: "0.3", 100: "1", 150: "1.5", 375: "3.75", 2400: "24", 4810: "48.1"}
	for centesimas, esperado := range casos {
		if got := HorasATexto(centesimas); got != esperado {
			t.Errorf("HorasATexto(%d) = %q, se esperaba %q", centesimas, got, esperado)
		}
	}
}

// spec 09 RN-2 a RN-5, RN-9 (CA-1, CA-3, CA-5)
func TestNuevoRegistroEsfuerzo(t *testing.T) {
	p := Proyecto{ID: 7, FechaInicio: fecha(t, "2026-11-01"), FechaFin: fecha(t, "2026-12-15")}
	hoy := fecha(t, "2026-11-20")
	valido := DatosEsfuerzo{IntegranteID: 3, Fecha: "2026-11-18", Actividad: "Diseño del formulario", Horas: "2.5"}
	con := func(cambiar func(*DatosEsfuerzo)) DatosEsfuerzo {
		d := valido
		cambiar(&d)
		return d
	}
	casos := []struct {
		nombre string
		datos  DatosEsfuerzo
		err    error
	}{
		{"datos válidos (CA-1)", valido, nil},
		{"sin integrante (spec 09 RN-2, CA-8)", con(func(d *DatosEsfuerzo) { d.IntegranteID = 0 }), ErrIntegranteObligatorio},
		{"sin fecha (spec 09 RN-3, CA-5)", con(func(d *DatosEsfuerzo) { d.Fecha = "" }), ErrFechaObligatoria},
		{"fecha inexistente (spec 09 RN-3, CA-5)", con(func(d *DatosEsfuerzo) { d.Fecha = "2026-02-30" }), ErrFechaInvalida},
		{"el día de inicio del proyecto (spec 09 RN-3, CA-5)", con(func(d *DatosEsfuerzo) { d.Fecha = "2026-11-01" }), nil},
		{"el día anterior al inicio (spec 09 RN-3, CA-5)", con(func(d *DatosEsfuerzo) { d.Fecha = "2026-10-31" }), ErrFechaFueraDelProyecto},
		{"hoy (spec 09 RN-3, CA-5)", con(func(d *DatosEsfuerzo) { d.Fecha = "2026-11-20" }), nil},
		{"mañana (spec 09 RN-3, CA-5)", con(func(d *DatosEsfuerzo) { d.Fecha = "2026-11-21" }), ErrFechaFutura},
		{"después del fin y en el futuro: primero el proyecto (spec 09 RN-3)", con(func(d *DatosEsfuerzo) { d.Fecha = "2027-01-10" }), ErrFechaFueraDelProyecto},
		{"actividad vacía (spec 09 RN-4, CA-5)", con(func(d *DatosEsfuerzo) { d.Actividad = "" }), ErrActividadObligatoria},
		{"actividad con solo espacios (spec 09 RN-4, CA-5)", con(func(d *DatosEsfuerzo) { d.Actividad = "   " }), ErrActividadObligatoria},
		{"actividad de 300 caracteres con eñes (spec 09 RN-4, CA-5)", con(func(d *DatosEsfuerzo) { d.Actividad = strings.Repeat("ñ", 300) }), nil},
		{"actividad de 301 caracteres (spec 09 RN-4, CA-5)", con(func(d *DatosEsfuerzo) { d.Actividad = strings.Repeat("a", 301) }), ErrActividadMuyLarga},
		{"sin horas (spec 09 RN-5, CA-3)", con(func(d *DatosEsfuerzo) { d.Horas = "" }), ErrHorasObligatorias},
		{"horas en cero (spec 09 RN-5, CA-3)", con(func(d *DatosEsfuerzo) { d.Horas = "0" }), ErrHorasInvalidas},
		{"varios inválidos: primero la fecha (spec 09 RN-9)", con(func(d *DatosEsfuerzo) { d.Fecha, d.Actividad, d.Horas = "2026-11-21", "", "0" }), ErrFechaFutura},
		{"actividad y horas inválidas: primero la actividad (spec 09 RN-9)", con(func(d *DatosEsfuerzo) { d.Actividad, d.Horas = "", "0" }), ErrActividadObligatoria},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			r, err := NuevoRegistroEsfuerzo(p, 5, c.datos, hoy)
			if !errors.Is(err, c.err) {
				t.Fatalf("error = %v, se esperaba %v", err, c.err)
			}
			if err != nil {
				return
			}
			if r.ID != 0 || r.ProyectoID != 7 || r.HistoriaID != 5 || r.IntegranteID != 3 ||
				r.Actividad != c.datos.Actividad || r.Fecha.Format(FormatoFecha) != c.datos.Fecha {
				t.Errorf("registro = %+v, no conserva los datos recibidos (el ID lo asigna el sistema)", r)
			}
		})
	}
}

// spec 09 RN-3: la hora de "hoy" no importa, solo el día.
func TestNuevoRegistroEsfuerzoHoyConHora(t *testing.T) {
	p := Proyecto{FechaInicio: fecha(t, "2026-11-01"), FechaFin: fecha(t, "2026-12-15")}
	hoy := time.Date(2026, 11, 20, 0, 5, 0, 0, time.FixedZone("ART", -3*60*60))
	d := DatosEsfuerzo{IntegranteID: 1, Fecha: "2026-11-20", Actividad: "Backend", Horas: "1"}
	if _, err := NuevoRegistroEsfuerzo(p, 1, d, hoy); err != nil {
		t.Errorf("registrar el mismo día de hoy: %v", err)
	}
}

// spec 09 RN-6, RN-7 (CA-4): orden por fecha (a igual fecha, por orden de carga) y total exacto.
func TestResumirEsfuerzo(t *testing.T) {
	reg := func(id int64, dia string, centesimas int) RegistroEsfuerzo {
		return RegistroEsfuerzo{ID: id, Fecha: fecha(t, dia), Actividad: dia, Centesimas: centesimas}
	}
	casos := []struct {
		nombre    string
		registros []RegistroEsfuerzo
		ids       []int64
		total     int
	}{
		{"sin registros", nil, nil, 0},
		{"por fecha y, a igual fecha, por carga", []RegistroEsfuerzo{
			reg(1, "2026-11-19", 200), reg(2, "2026-11-18", 150), reg(3, "2026-11-19", 25),
		}, []int64{2, 1, 3}, 375},
		{"aunque lleguen desordenados", []RegistroEsfuerzo{
			reg(9, "2026-11-19", 100), reg(4, "2026-11-19", 100),
		}, []int64{4, 9}, 200},
		{"0,1 + 0,2 es exactamente 0,3", []RegistroEsfuerzo{
			reg(1, "2026-11-18", 10), reg(2, "2026-11-18", 20),
		}, []int64{1, 2}, 30},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			e := ResumirEsfuerzo(c.registros)
			var ids []int64
			for _, r := range e.Registros {
				ids = append(ids, r.ID)
			}
			if len(ids) != len(c.ids) {
				t.Fatalf("registros = %v, se esperaba %v", ids, c.ids)
			}
			for i := range ids {
				if ids[i] != c.ids[i] {
					t.Fatalf("orden = %v, se esperaba %v", ids, c.ids)
				}
			}
			if e.TotalCentesimas != c.total {
				t.Errorf("total = %d centésimas, se esperaba %d", e.TotalCentesimas, c.total)
			}
			if e.Registros == nil {
				t.Error("la lista de registros no debe ser nil (la API la devuelve como [])")
			}
		})
	}
}

// Los errores de validación del esfuerzo son datos inválidos (400); "no existe" (404) y
// "historia con esfuerzo" (409) no.
func TestErroresDeEsfuerzo(t *testing.T) {
	invalidos := []error{
		ErrIntegranteObligatorio, ErrFechaObligatoria, ErrFechaFueraDelProyecto, ErrFechaFutura,
		ErrActividadObligatoria, ErrActividadMuyLarga, ErrHorasObligatorias, ErrHorasInvalidas,
	}
	for _, e := range invalidos {
		var d ErrDatoInvalido
		if !errors.As(e, &d) {
			t.Errorf("%q debería ser un ErrDatoInvalido", e)
		}
	}
	var c ErrConflicto
	if !errors.As(ErrHistoriaConEsfuerzo, &c) {
		t.Error("ErrHistoriaConEsfuerzo debería ser un ErrConflicto (409)")
	}
	var d ErrDatoInvalido
	if errors.As(ErrRegistroEsfuerzoNoExiste, &d) || errors.As(ErrRegistroEsfuerzoNoExiste, &c) {
		t.Error("ErrRegistroEsfuerzoNoExiste es un 'no existe' (404)")
	}
}
