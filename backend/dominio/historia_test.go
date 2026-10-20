package dominio

import (
	"errors"
	"strings"
	"testing"
)

// spec 04 RN-2 a RN-6 (CA-1, CA-5, CA-6)
func TestNuevaHistoria(t *testing.T) {
	valida := DatosHistoria{Titulo: "Registrar usuario", Descripcion: "Alta de clientes", Prioridad: 1, Criterios: "Valida el email"}
	con := func(cambiar func(*DatosHistoria)) DatosHistoria {
		d := valida
		cambiar(&d)
		return d
	}
	casos := []struct {
		nombre string
		datos  DatosHistoria
		err    error
	}{
		{"datos válidos (CA-1)", valida, nil},
		{"sin criterios de aceptación (spec 04 RN-5, CA-1)", con(func(d *DatosHistoria) { d.Criterios = "" }), nil},
		{"título vacío (spec 04 RN-2, CA-5)", con(func(d *DatosHistoria) { d.Titulo = "" }), ErrTituloObligatorio},
		{"título con solo espacios (spec 04 RN-2, CA-5)", con(func(d *DatosHistoria) { d.Titulo = "   " }), ErrTituloObligatorio},
		{"título de 150 caracteres (spec 04 RN-2, CA-5)", con(func(d *DatosHistoria) { d.Titulo = strings.Repeat("t", 150) }), nil},
		{"título de 151 caracteres (spec 04 RN-2, CA-5)", con(func(d *DatosHistoria) { d.Titulo = strings.Repeat("t", 151) }), ErrTituloMuyLargo},
		{"título de 150 caracteres con eñes (spec 04 RN-2, CA-5)", con(func(d *DatosHistoria) { d.Titulo = strings.Repeat("ñ", 150) }), nil},
		{"descripción vacía (spec 04 RN-3, CA-5)", con(func(d *DatosHistoria) { d.Descripcion = "" }), ErrDescripcionObligatoria},
		{"descripción con solo espacios (spec 04 RN-3, CA-5)", con(func(d *DatosHistoria) { d.Descripcion = "  " }), ErrDescripcionObligatoria},
		{"descripción de 1000 caracteres (spec 04 RN-3, CA-5)", con(func(d *DatosHistoria) { d.Descripcion = strings.Repeat("d", 1000) }), nil},
		{"descripción de 1001 caracteres (spec 04 RN-3, CA-5)", con(func(d *DatosHistoria) { d.Descripcion = strings.Repeat("d", 1001) }), ErrDescripcionHistoriaMuyLarga},
		{"prioridad 1, la más alta (spec 04 RN-4, CA-6)", con(func(d *DatosHistoria) { d.Prioridad = 1 }), nil},
		{"prioridad 1000, sin máximo (spec 04 RN-4, CA-6)", con(func(d *DatosHistoria) { d.Prioridad = 1000 }), nil},
		{"prioridad 0 (spec 04 RN-4, CA-6)", con(func(d *DatosHistoria) { d.Prioridad = 0 }), ErrPrioridadInvalida},
		{"prioridad negativa (spec 04 RN-4, CA-6)", con(func(d *DatosHistoria) { d.Prioridad = -2 }), ErrPrioridadInvalida},
		{"criterios de 2000 caracteres (spec 04 RN-5, CA-5)", con(func(d *DatosHistoria) { d.Criterios = strings.Repeat("c", 2000) }), nil},
		{"criterios de 2001 caracteres (spec 04 RN-5, CA-5)", con(func(d *DatosHistoria) { d.Criterios = strings.Repeat("c", 2001) }), ErrCriteriosMuyLargos},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			h, err := NuevaHistoria(7, c.datos)
			if !errors.Is(err, c.err) {
				t.Fatalf("error = %v, se esperaba %v", err, c.err)
			}
			if err != nil {
				return
			}
			if h.ProyectoID != 7 || h.Titulo != c.datos.Titulo || h.Descripcion != c.datos.Descripcion ||
				h.Prioridad != c.datos.Prioridad || h.Criterios != c.datos.Criterios {
				t.Errorf("historia = %+v, no conserva los datos recibidos", h)
			}
			if h.ID != 0 {
				t.Errorf("ID = %d, el identificador lo asigna el sistema al guardar (spec 04 RN-10)", h.ID)
			}
			// RN-6: una historia nueva arranca Pendiente y sin estimar.
			if h.Estado != EstadoPendiente {
				t.Errorf("estado = %q, se esperaba %q (spec 04 RN-6)", h.Estado, EstadoPendiente)
			}
			if h.StoryPoints != nil {
				t.Errorf("Story Points = %d, se esperaba sin estimar (spec 04 RN-6)", *h.StoryPoints)
			}
		})
	}
}

// Los errores de validación de la historia son datos inválidos (400); "no existe" no (404).
func TestErroresDeHistoria(t *testing.T) {
	invalidos := []error{
		ErrTituloObligatorio, ErrTituloMuyLargo, ErrDescripcionObligatoria,
		ErrDescripcionHistoriaMuyLarga, ErrPrioridadInvalida, ErrCriteriosMuyLargos,
	}
	for _, e := range invalidos {
		var d ErrDatoInvalido
		if !errors.As(e, &d) {
			t.Errorf("%q debería ser un ErrDatoInvalido", e)
		}
	}
	var d ErrDatoInvalido
	if errors.As(ErrHistoriaNoExiste, &d) {
		t.Error("ErrHistoriaNoExiste no es un dato inválido")
	}
}
