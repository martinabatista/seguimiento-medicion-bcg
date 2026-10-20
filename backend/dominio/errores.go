// Package dominio contiene las entidades, reglas de negocio y cálculos del sistema
// (proyectos, backlog, Sprints, estimación, esfuerzo, defectos y métricas).
// No depende de HTTP ni de la base de datos: se testea con go test puro (TDD).
package dominio

import "errors"

// ErrDatoInvalido es un dato que no cumple una regla de negocio. La API lo responde como 400.
// El mensaje es el que ve la usuaria, tal cual lo define la spec.
type ErrDatoInvalido struct{ Mensaje string }

func (e ErrDatoInvalido) Error() string { return e.Mensaje }

// Errores de la spec 01 (proyectos).
var (
	ErrNombreObligatorio      = ErrDatoInvalido{"El nombre es obligatorio"}
	ErrDescripcionObligatoria = ErrDatoInvalido{"La descripción es obligatoria"}
	ErrDescripcionMuyLarga    = ErrDatoInvalido{"La descripción no puede superar los 500 caracteres"}
	ErrFechasObligatorias     = ErrDatoInvalido{"Las fechas de inicio y finalización son obligatorias"}
	ErrFechaInvalida          = ErrDatoInvalido{"La fecha no es válida"}
	ErrFinNoPosteriorAInicio  = ErrDatoInvalido{"La fecha de finalización debe ser posterior a la de inicio"}

	// ErrProyectoNoExiste no es un dato inválido: la API lo responde como 404.
	ErrProyectoNoExiste = errors.New("El proyecto no existe")
)
