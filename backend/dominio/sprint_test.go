package dominio

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// proyecto de prueba: del 2026-11-01 al 2026-12-31.
func proyectoDeSprints(t *testing.T) Proyecto {
	return Proyecto{ID: 1, Nombre: "Tienda Online", Descripcion: "Venta de ropa",
		FechaInicio: fecha(t, "2026-11-01"), FechaFin: fecha(t, "2026-12-31")}
}

// spec 05 RN-2 a RN-6 (CA-1, CA-2)
func TestNuevoSprint(t *testing.T) {
	p := proyectoDeSprints(t)
	existente := Sprint{ID: 7, ProyectoID: 1, Numero: 1, FechaInicio: fecha(t, "2026-11-15"), FechaFin: fecha(t, "2026-11-28"), Goal: "Catálogo", Estado: SprintAbierto}
	casos := []struct {
		nombre     string
		inicio     string
		fin        string
		goal       string
		existentes []Sprint
		err        error
		numero     int
	}{
		{"primer Sprint (CA-1)", "2026-11-01", "2026-11-14", "Registrar clientes", nil, nil, 1},
		{"número siguiente al último (RN-6)", "2026-11-01", "2026-11-14", "Registrar clientes", []Sprint{existente}, nil, 2},
		{"empieza el día siguiente al fin de otro (RN-4)", "2026-11-29", "2026-12-05", "Pagos", []Sprint{existente}, nil, 2},
		{"termina el día anterior al inicio de otro (RN-4)", "2026-11-01", "2026-11-14", "Pagos", []Sprint{existente}, nil, 2},
		{"igual al período del proyecto (RN-3)", "2026-11-01", "2026-12-31", "Todo", nil, nil, 1},
		{"dura un solo día (RN-2)", "2026-11-01", "2026-11-02", "Rápido", nil, nil, 1},
		{"falta la fecha de inicio (RN-2)", "", "2026-11-14", "Objetivo", nil, ErrFechasObligatorias, 0},
		{"falta la fecha de fin (RN-2)", "2026-11-01", "", "Objetivo", nil, ErrFechasObligatorias, 0},
		{"fin igual al inicio (RN-2)", "2026-12-01", "2026-12-01", "Objetivo", nil, ErrFinNoPosteriorAInicio, 0},
		{"fin anterior al inicio (RN-2)", "2026-12-10", "2026-12-01", "Objetivo", nil, ErrFinNoPosteriorAInicio, 0},
		{"empieza antes del proyecto (RN-3)", "2026-10-31", "2026-11-10", "Objetivo", nil, ErrSprintFueraDelProyecto, 0},
		{"termina después del proyecto (RN-3)", "2026-12-20", "2027-01-02", "Objetivo", nil, ErrSprintFueraDelProyecto, 0},
		{"Sprint Goal vacío (RN-5)", "2026-11-01", "2026-11-14", "", nil, ErrSprintGoalObligatorio, 0},
		{"Sprint Goal con solo espacios (RN-5)", "2026-11-01", "2026-11-14", "   ", nil, ErrSprintGoalObligatorio, 0},
		{"Sprint Goal de 300 caracteres (RN-5)", "2026-11-01", "2026-11-14", strings.Repeat("g", 300), nil, nil, 1},
		{"Sprint Goal de 301 caracteres (RN-5)", "2026-11-01", "2026-11-14", strings.Repeat("g", 301), nil, ErrSprintGoalMuyLargo, 0},
		{"se superpone al principio de otro (RN-4)", "2026-11-10", "2026-11-15", "Objetivo", []Sprint{existente}, ErrSprintSuperpuesto, 0},
		{"empieza el día en que termina otro (RN-4)", "2026-11-28", "2026-12-05", "Objetivo", []Sprint{existente}, ErrSprintSuperpuesto, 0},
		{"contiene a otro (RN-4)", "2026-11-10", "2026-12-05", "Objetivo", []Sprint{existente}, ErrSprintSuperpuesto, 0},
		{"se superpone con uno cerrado (RN-4)", "2026-11-20", "2026-12-05", "Objetivo", []Sprint{{Numero: 1, FechaInicio: fecha(t, "2026-11-15"), FechaFin: fecha(t, "2026-11-28"), Estado: SprintCerrado}}, ErrSprintSuperpuesto, 0},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			var inicio, fin time.Time
			if c.inicio != "" {
				inicio = fecha(t, c.inicio)
			}
			if c.fin != "" {
				fin = fecha(t, c.fin)
			}
			s, err := NuevoSprint(p, inicio, fin, c.goal, c.existentes)
			if !errors.Is(err, c.err) {
				t.Fatalf("error = %v, se esperaba %v", err, c.err)
			}
			if err != nil {
				return
			}
			if s.ProyectoID != p.ID || s.Numero != c.numero || s.Estado != SprintAbierto || len(s.Historias) != 0 ||
				!s.FechaInicio.Equal(inicio) || !s.FechaFin.Equal(fin) || s.Goal != c.goal || s.ID != 0 {
				t.Errorf("sprint = %+v, se esperaba el Sprint %d abierto, sin historias y con los datos recibidos", s, c.numero)
			}
		})
	}
}

// spec 05 RN-12 (CA-5): el Sprint en curso es el abierto cuyo período incluye hoy.
func TestSprintEnCurso(t *testing.T) {
	s1 := Sprint{ID: 1, Numero: 1, FechaInicio: fecha(t, "2026-11-01"), FechaFin: fecha(t, "2026-11-14"), Estado: SprintAbierto}
	s2 := Sprint{ID: 2, Numero: 2, FechaInicio: fecha(t, "2026-11-16"), FechaFin: fecha(t, "2026-11-28"), Estado: SprintAbierto}
	cerrado := Sprint{ID: 3, Numero: 3, FechaInicio: fecha(t, "2026-12-01"), FechaFin: fecha(t, "2026-12-14"), Estado: SprintCerrado}
	sprints := []Sprint{s2, s1, cerrado}
	casos := []struct {
		hoy    string
		numero int // 0 = no hay Sprint en curso
	}{
		{"2026-11-01", 1}, // primer día
		{"2026-11-14", 1}, // último día: sigue en curso
		{"2026-11-15", 0}, // entre dos Sprints
		{"2026-11-20", 2},
		{"2026-12-05", 0}, // el Sprint de esas fechas está cerrado
		{"2027-01-10", 0},
	}
	for _, c := range casos {
		t.Run(c.hoy, func(t *testing.T) {
			// la hora del día no importa: solo cuenta la fecha
			hoy := fecha(t, c.hoy).Add(23 * time.Hour)
			s, ok := SprintEnCurso(sprints, hoy)
			if c.numero == 0 {
				if ok {
					t.Errorf("Sprint en curso = %d, se esperaba ninguno", s.Numero)
				}
				return
			}
			if !ok || s.Numero != c.numero {
				t.Errorf("Sprint en curso = %d (hay = %v), se esperaba %d", s.Numero, ok, c.numero)
			}
		})
	}
}

// Los conflictos se distinguen de los datos inválidos (409 vs 400 en la API).
func TestErroresDeSprintSonConflictos(t *testing.T) {
	conflictos := []error{
		ErrSprintSuperpuesto, ErrSprintCerrado, ErrHistoriaCompletada, ErrHistoriaYaAsignada,
		ErrHistoriaNoAsignada, ErrHistoriaEnSprintAbierto,
	}
	for _, e := range conflictos {
		var c ErrConflicto
		if !errors.As(e, &c) {
			t.Errorf("%q debería ser un ErrConflicto", e)
		}
	}
	var d ErrDatoInvalido
	for _, e := range []error{ErrSprintFueraDelProyecto, ErrSprintGoalObligatorio, ErrSprintGoalMuyLargo} {
		if !errors.As(e, &d) {
			t.Errorf("%q debería ser un ErrDatoInvalido", e)
		}
	}
}
