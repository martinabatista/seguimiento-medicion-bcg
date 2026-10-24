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

// Errores de la spec 04 (Product Backlog). La descripción vacía usa ErrDescripcionObligatoria
// de la spec 01: el mensaje es el mismo.
var (
	ErrTituloObligatorio           = ErrDatoInvalido{"El título es obligatorio"}
	ErrTituloMuyLargo              = ErrDatoInvalido{"El título no puede superar los 150 caracteres"}
	ErrDescripcionHistoriaMuyLarga = ErrDatoInvalido{"La descripción no puede superar los 1000 caracteres"}
	ErrPrioridadObligatoria        = ErrDatoInvalido{"La prioridad es obligatoria"}
	ErrPrioridadInvalida           = ErrDatoInvalido{"La prioridad debe ser un número entero mayor o igual a 1"}
	ErrCriteriosMuyLargos          = ErrDatoInvalido{"Los criterios de aceptación no pueden superar los 2000 caracteres"}
	ErrEstadoInvalido              = ErrDatoInvalido{"El estado no es válido"}

	// ErrHistoriaNoExiste no es un dato inválido: la API lo responde como 404.
	ErrHistoriaNoExiste = errors.New("La historia no existe")
)

// Errores de la spec 02 (integrantes). El nombre vacío usa ErrNombreObligatorio de la spec 01:
// el mensaje es el mismo.
var (
	ErrNombreIntegranteMuyLargo = ErrDatoInvalido{"El nombre no puede superar los 100 caracteres"}
	ErrEmailObligatorio         = ErrDatoInvalido{"El email es obligatorio"}
	ErrEmailInvalido            = ErrDatoInvalido{"El email no es válido"}
	ErrEmailMuyLargo            = ErrDatoInvalido{"El email no puede superar los 254 caracteres"}
	ErrRolObligatorio           = ErrDatoInvalido{"El rol es obligatorio"}
	ErrRolInvalido              = ErrDatoInvalido{"El rol no es válido"}

	// ErrIntegranteNoExiste no es un dato inválido: la API lo responde como 404.
	ErrIntegranteNoExiste = errors.New("El integrante no existe")
	// ErrIntegranteYaRegistrado es un conflicto con lo que ya hay: la API lo responde como 409.
	ErrIntegranteYaRegistrado = errors.New("El integrante ya está registrado en el proyecto")
)

// ErrConflicto es una operación que choca con lo que ya está registrado (no con el formato de un
// dato). La API lo responde como 409. El mensaje es el que ve la usuaria, tal cual lo define la spec.
type ErrConflicto struct{ Mensaje string }

func (e ErrConflicto) Error() string { return e.Mensaje }

// Errores de la spec 05 (Sprints). Las fechas faltantes, inválidas o invertidas usan los errores
// de la spec 01: los mensajes son los mismos.
var (
	ErrSprintFueraDelProyecto = ErrDatoInvalido{"El Sprint debe estar dentro de las fechas del proyecto"}
	ErrSprintGoalObligatorio  = ErrDatoInvalido{"El Sprint Goal es obligatorio"}
	ErrSprintGoalMuyLargo     = ErrDatoInvalido{"El Sprint Goal no puede superar los 300 caracteres"}

	ErrSprintSuperpuesto       = ErrConflicto{"El Sprint se superpone con otro Sprint del proyecto"}
	ErrSprintCerrado           = ErrConflicto{"El Sprint está cerrado"}
	ErrHistoriaCompletada      = ErrConflicto{"La historia ya está completada"}
	ErrHistoriaYaAsignada      = ErrConflicto{"La historia ya está asignada a un Sprint abierto"}
	ErrHistoriaNoAsignada      = ErrConflicto{"La historia no está asignada a este Sprint"}
	ErrHistoriaEnSprintAbierto = ErrConflicto{"La historia está asignada a un Sprint abierto"}

	// ErrSprintNoExiste no es un dato inválido: la API lo responde como 404.
	ErrSprintNoExiste = errors.New("El Sprint no existe")
)

// Errores de la spec 06 (historias completadas y cierre de Sprints). Completar dos veces, el Sprint
// cerrado y la historia no asignada usan los errores de la spec 05: los mensajes son los mismos.
var (
	ErrHistoriaNoCompletada = ErrConflicto{"La historia no está completada"}
	ErrCompletarDesdeSprint = ErrConflicto{"La historia está en un Sprint abierto: se completa desde el Sprint"}
)

// Errores de la spec 07 (Story Points). Estimar una historia completada usa ErrHistoriaCompletada
// de la spec 05: el mensaje es el mismo. El mensaje de la escala se arma con sus valores.
var (
	ErrStoryPointsObligatorios  = ErrDatoInvalido{"Los Story Points son obligatorios"}
	ErrStoryPointsFueraDeEscala = ErrDatoInvalido{"Los Story Points deben ser un valor de la escala: " + textoEscala()}
)

// Errores de la spec 08 (Planning Poker). Abrir o acordar con una historia completada usa
// ErrHistoriaCompletada (spec 05) y los valores de los votos usan los errores de la spec 07.
var (
	ErrVotacionAbierta = ErrConflicto{"La historia ya tiene una votación abierta"}
	ErrSinDevelopers   = ErrConflicto{"El proyecto no tiene Developers para votar"}
	ErrSoloDevelopers  = ErrConflicto{"Solo los Developers pueden votar"}
	ErrRondaRevelada   = ErrConflicto{"La ronda ya fue revelada"}
	ErrFaltanVotos     = ErrConflicto{"Faltan votos: todavía no votaron todos los Developers"}
	ErrRondaSinRevelar = ErrConflicto{"Primero hay que revelar la ronda actual"}
	ErrVotacionCerrada = ErrConflicto{"La votación está cerrada"}

	// ErrVotacionNoExiste no es un conflicto: la API lo responde como 404.
	ErrVotacionNoExiste = errors.New("La votación no existe")
)

// Errores de la spec 09 (registro de esfuerzo). La fecha con formato inválido usa ErrFechaInvalida
// de la spec 01 y el integrante inexistente, ErrIntegranteNoExiste de la spec 02.
var (
	ErrIntegranteObligatorio = ErrDatoInvalido{"El integrante es obligatorio"}
	ErrFechaObligatoria      = ErrDatoInvalido{"La fecha es obligatoria"}
	ErrFechaFueraDelProyecto = ErrDatoInvalido{"La fecha debe estar dentro de las fechas del proyecto"}
	ErrFechaFutura           = ErrDatoInvalido{"La fecha no puede ser posterior a hoy"}
	ErrActividadObligatoria  = ErrDatoInvalido{"La actividad es obligatoria"}
	ErrActividadMuyLarga     = ErrDatoInvalido{"La actividad no puede superar los 300 caracteres"}
	ErrHorasObligatorias     = ErrDatoInvalido{"Las horas son obligatorias"}
	ErrHorasInvalidas        = ErrDatoInvalido{"Las horas deben ser un número mayor a 0 y hasta 24, con hasta 2 decimales"}

	ErrHistoriaConEsfuerzo = ErrConflicto{"La historia tiene esfuerzo registrado"}

	// ErrRegistroEsfuerzoNoExiste no es un conflicto: la API lo responde como 404.
	ErrRegistroEsfuerzoNoExiste = errors.New("El registro de esfuerzo no existe")
)
