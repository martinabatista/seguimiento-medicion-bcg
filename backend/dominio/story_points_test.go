package dominio

import (
	"errors"
	"slices"
	"testing"
)

// spec 07 RN-1 (CA-5): la escala es Fibonacci hasta 21, de menor a mayor.
func TestEscalaStoryPoints(t *testing.T) {
	esperada := []int{1, 2, 3, 5, 8, 13, 21}
	escala := EscalaStoryPoints()
	if !slices.Equal(escala, esperada) {
		t.Fatalf("escala = %v, se esperaba %v", escala, esperada)
	}
	// Quien la consulta recibe una copia: no puede cambiar la escala del sistema.
	escala[0] = 99
	if EscalaStoryPoints()[0] != 1 {
		t.Error("modificar la escala devuelta cambió la escala del sistema")
	}
}

// spec 07 RN-2 (CA-3): solo se aceptan valores de la escala.
func TestEnEscala(t *testing.T) {
	for _, v := range []int{1, 2, 3, 5, 8, 13, 21} {
		if !EnEscala(v) {
			t.Errorf("EnEscala(%d) = false, es un valor de la escala", v)
		}
	}
	for _, v := range []int{0, 4, 6, 10, 14, 22, 34, -3} {
		if EnEscala(v) {
			t.Errorf("EnEscala(%d) = true, no es un valor de la escala", v)
		}
	}
}

// spec 07 RN-2 (CA-3): los Story Points llegan como texto (el número del JSON).
func TestParsearStoryPoints(t *testing.T) {
	casos := []struct {
		nombre string
		texto  string
		valor  int
		err    error
	}{
		{"el menor de la escala", "1", 1, nil},
		{"el mayor de la escala", "21", 21, nil},
		{"un valor del medio", "5", 5, nil},
		{"falta el valor", "", 0, ErrStoryPointsObligatorios},
		{"cero", "0", 0, ErrStoryPointsFueraDeEscala},
		{"entre valores de la escala", "4", 0, ErrStoryPointsFueraDeEscala},
		{"mayor que 21", "22", 0, ErrStoryPointsFueraDeEscala},
		{"negativo", "-3", 0, ErrStoryPointsFueraDeEscala},
		{"con decimales", "2.5", 0, ErrStoryPointsFueraDeEscala},
		{"texto que no es un número", "mucho", 0, ErrStoryPointsFueraDeEscala},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			v, err := ParsearStoryPoints(c.texto)
			if !errors.Is(err, c.err) {
				t.Fatalf("error = %v, se esperaba %v", err, c.err)
			}
			if v != c.valor {
				t.Errorf("Story Points = %d, se esperaba %d", v, c.valor)
			}
		})
	}
}

// spec 07: el mensaje lista la escala y los errores de valor son datos inválidos (400).
func TestErroresDeStoryPoints(t *testing.T) {
	esperado := "Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21"
	if ErrStoryPointsFueraDeEscala.Error() != esperado {
		t.Errorf("mensaje = %q, se esperaba %q", ErrStoryPointsFueraDeEscala.Error(), esperado)
	}
	for _, e := range []error{ErrStoryPointsObligatorios, ErrStoryPointsFueraDeEscala} {
		var d ErrDatoInvalido
		if !errors.As(e, &d) {
			t.Errorf("%q debería ser un ErrDatoInvalido", e)
		}
	}
}
