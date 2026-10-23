package dominio

import (
	"slices"
	"strconv"
	"strings"
)

// escalaStoryPoints es la escala de Fibonacci hasta 21 (spec 07 RN-1): cada valor es la suma de
// los dos anteriores, así que cuanto más grande es una historia, más separados están los valores.
var escalaStoryPoints = []int{1, 2, 3, 5, 8, 13, 21}

// EscalaStoryPoints devuelve los valores permitidos de menor a mayor (spec 07 CA-5). Es una copia:
// quien la consulta no puede cambiar la escala del sistema.
func EscalaStoryPoints() []int {
	return slices.Clone(escalaStoryPoints)
}

// EnEscala indica si el valor es uno de los de la escala (spec 07 RN-2).
func EnEscala(storyPoints int) bool {
	return slices.Contains(escalaStoryPoints, storyPoints)
}

// ParsearStoryPoints convierte el texto de los Story Points en número (spec 07 RN-2).
// Vacío es "falta el valor"; cualquier cosa que no sea un valor de la escala es inválida.
func ParsearStoryPoints(texto string) (int, error) {
	if texto == "" {
		return 0, ErrStoryPointsObligatorios
	}
	sp, err := strconv.Atoi(texto)
	if err != nil || !EnEscala(sp) {
		return 0, ErrStoryPointsFueraDeEscala
	}
	return sp, nil
}

// textoEscala arma "1, 2, 3, 5, 8, 13 o 21" para el mensaje de error, a partir de la escala.
func textoEscala() string {
	valores := make([]string, len(escalaStoryPoints))
	for i, v := range escalaStoryPoints {
		valores[i] = strconv.Itoa(v)
	}
	ultimo := len(valores) - 1
	return strings.Join(valores[:ultimo], ", ") + " o " + valores[ultimo]
}
