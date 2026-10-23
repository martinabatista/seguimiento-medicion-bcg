package dominio

import (
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

// MaxSprintGoal es el largo máximo del Sprint Goal (spec 05 RN-5).
const MaxSprintGoal = 300

// EstadoSprint indica si un Sprint admite cambios. Se cierra en la spec 06.
type EstadoSprint string

const (
	SprintAbierto EstadoSprint = "Abierto"
	SprintCerrado EstadoSprint = "Cerrado"
)

// Sprint es una iteración de un proyecto (spec 05): sus fechas, su objetivo y las historias del
// backlog que entran.
type Sprint struct {
	ID          int64
	ProyectoID  int64
	Numero      int // correlativo dentro del proyecto: Sprint 1, Sprint 2… (RN-6)
	FechaInicio time.Time
	FechaFin    time.Time
	Goal        string
	Estado      EstadoSprint
	Historias   []int64 // identificadores de las historias asignadas (lo planificado)
	Completadas []int64 // las que se completaron en este Sprint (spec 06 RN-2); no cambia al cerrarlo (RN-8)
}

// NuevoSprint arma un Sprint del proyecto validando las reglas de la spec 05 contra los Sprints
// que el proyecto ya tiene. El ID queda en cero: lo asigna el sistema al guardarlo (RN-13).
func NuevoSprint(p Proyecto, inicio, fin time.Time, goal string, existentes []Sprint) (Sprint, error) {
	// RN-2: fechas obligatorias y fin posterior al inicio (como en la spec 01).
	if inicio.IsZero() || fin.IsZero() {
		return Sprint{}, ErrFechasObligatorias
	}
	if !fin.After(inicio) {
		return Sprint{}, ErrFinNoPosteriorAInicio
	}
	// RN-3: dentro del período del proyecto (los extremos están permitidos).
	if inicio.Before(p.FechaInicio) || fin.After(p.FechaFin) {
		return Sprint{}, ErrSprintFueraDelProyecto
	}
	if err := ValidarSprintGoal(goal); err != nil {
		return Sprint{}, err
	}
	numero := 1
	for _, e := range existentes {
		// RN-4: ningún día puede pertenecer a dos Sprints, estén abiertos o cerrados.
		if seSuperponen(inicio, fin, e.FechaInicio, e.FechaFin) {
			return Sprint{}, ErrSprintSuperpuesto
		}
		// RN-6: el número siguiente al del último Sprint del proyecto.
		if e.Numero >= numero {
			numero = e.Numero + 1
		}
	}
	return Sprint{
		ProyectoID:  p.ID,
		Numero:      numero,
		FechaInicio: inicio,
		FechaFin:    fin,
		Goal:        goal,
		Estado:      SprintAbierto,
		Historias:   []int64{},
		Completadas: []int64{},
	}, nil
}

// ValidarSprintGoal aplica la RN-5: obligatorio y hasta 300 caracteres (no bytes).
func ValidarSprintGoal(goal string) error {
	if strings.TrimSpace(goal) == "" {
		return ErrSprintGoalObligatorio
	}
	if utf8.RuneCountInString(goal) > MaxSprintGoal {
		return ErrSprintGoalMuyLargo
	}
	return nil
}

// seSuperponen indica si dos períodos comparten al menos un día (los extremos cuentan).
func seSuperponen(inicioA, finA, inicioB, finB time.Time) bool {
	return !finA.Before(inicioB) && !finB.Before(inicioA)
}

// Tiene indica si la historia está asignada al Sprint.
func (s Sprint) Tiene(historiaID int64) bool {
	return slices.Contains(s.Historias, historiaID)
}

// EstaCompletada indica si la historia se completó en este Sprint (spec 06 RN-2).
func (s Sprint) EstaCompletada(historiaID int64) bool {
	return slices.Contains(s.Completadas, historiaID)
}

// SprintEnCurso devuelve el Sprint abierto cuyo período incluye la fecha de hoy (RN-12). Por la
// RN-4 hay como máximo uno. Solo cuenta el día de hoy, no la hora.
func SprintEnCurso(sprints []Sprint, hoy time.Time) (Sprint, bool) {
	dia := soloDia(hoy)
	for _, s := range sprints {
		if s.Estado == SprintAbierto && !dia.Before(soloDia(s.FechaInicio)) && !dia.After(soloDia(s.FechaFin)) {
			return s, true
		}
	}
	return Sprint{}, false
}
