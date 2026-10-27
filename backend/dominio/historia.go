package dominio

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Largos máximos de los textos de una historia (spec 04 RN-2, RN-3, RN-5).
const (
	MaxTituloHistoria      = 150
	MaxDescripcionHistoria = 1000
	MaxCriteriosHistoria   = 2000
)

// EstadoHistoria es el estado de una historia del backlog (spec 04 RN-7).
type EstadoHistoria string

const (
	EstadoPendiente  EstadoHistoria = "Pendiente"
	EstadoEnProgreso EstadoHistoria = "En progreso"
	EstadoCompletada EstadoHistoria = "Completada"
)

// ParsearEstado convierte un texto en estado. Se escribe exactamente como en la spec
// (con mayúscula): cualquier otro texto es inválido (RN-7).
func ParsearEstado(texto string) (EstadoHistoria, error) {
	switch e := EstadoHistoria(texto); e {
	case EstadoPendiente, EstadoEnProgreso, EstadoCompletada:
		return e, nil
	}
	return "", ErrEstadoInvalido
}

// ParsearPrioridad convierte el texto de la prioridad en número (RN-4).
// Vacío es "falta la prioridad"; con decimales, cero o negativo es inválido.
func ParsearPrioridad(texto string) (int, error) {
	if texto == "" {
		return 0, ErrPrioridadObligatoria
	}
	p, err := strconv.Atoi(texto)
	if err != nil || p < 1 {
		return 0, ErrPrioridadInvalida
	}
	return p, nil
}

// DatosHistoria son los datos que carga la usuaria al crear o modificar una historia.
type DatosHistoria struct {
	Titulo      string
	Descripcion string
	Prioridad   int // 1 es la más importante (RN-4)
	Criterios   string
}

// Historia es un elemento del Product Backlog de un proyecto (spec 04).
type Historia struct {
	ID          int64
	ProyectoID  int64
	Titulo      string
	Descripcion string
	Prioridad   int
	Criterios   string
	Estado      EstadoHistoria
	StoryPoints *int // nil = sin estimar; se asignan en las specs 07 y 08
	// HorasEstimadas en centésimas (8 h = 800); nil = sin horas estimadas. Se cargan en la spec 12.
	HorasEstimadas *int
}

// NuevaHistoria arma una historia del backlog validando las reglas de la spec 04.
// El ID queda en cero: lo asigna el sistema al guardarla (RN-10).
func NuevaHistoria(proyectoID int64, d DatosHistoria) (Historia, error) {
	// RN-2: título obligatorio, hasta 150 caracteres (no bytes: cuenta bien las tildes y la ñ).
	if strings.TrimSpace(d.Titulo) == "" {
		return Historia{}, ErrTituloObligatorio
	}
	if utf8.RuneCountInString(d.Titulo) > MaxTituloHistoria {
		return Historia{}, ErrTituloMuyLargo
	}
	// RN-3: descripción obligatoria, hasta 1000 caracteres.
	if strings.TrimSpace(d.Descripcion) == "" {
		return Historia{}, ErrDescripcionObligatoria
	}
	if utf8.RuneCountInString(d.Descripcion) > MaxDescripcionHistoria {
		return Historia{}, ErrDescripcionHistoriaMuyLarga
	}
	// RN-4: prioridad entera mayor o igual a 1, sin máximo.
	if d.Prioridad < 1 {
		return Historia{}, ErrPrioridadInvalida
	}
	// RN-5: criterios de aceptación opcionales, hasta 2000 caracteres.
	if utf8.RuneCountInString(d.Criterios) > MaxCriteriosHistoria {
		return Historia{}, ErrCriteriosMuyLargos
	}
	// RN-6: arranca Pendiente y con los Story Points sin estimar (nil).
	return Historia{
		ProyectoID:  proyectoID,
		Titulo:      d.Titulo,
		Descripcion: d.Descripcion,
		Prioridad:   d.Prioridad,
		Criterios:   d.Criterios,
		Estado:      EstadoPendiente,
	}, nil
}

// OrdenarBacklog devuelve las historias ordenadas como se ve el backlog (RN-8): por prioridad
// ascendente (1 primero) y, a igual prioridad, por orden de creación (el ID es creciente).
// El desempate es explícito para que el orden no dependa de cómo lleguen de la base.
// No modifica el slice recibido.
func OrdenarBacklog(historias []Historia) []Historia {
	ordenadas := slices.Clone(historias)
	slices.SortFunc(ordenadas, func(a, b Historia) int {
		return cmp.Or(cmp.Compare(a.Prioridad, b.Prioridad), cmp.Compare(a.ID, b.ID))
	})
	return ordenadas
}
