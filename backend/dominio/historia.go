package dominio

import (
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
		return Historia{}, ErrDescripcionHistoriaObligatoria
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
