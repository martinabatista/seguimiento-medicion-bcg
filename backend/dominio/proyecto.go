package dominio

import (
	"strings"
	"time"
	"unicode/utf8"
)

// MaxDescripcion es el largo máximo de la descripción de un proyecto (spec 01 RN-2).
const MaxDescripcion = 500

// FormatoFecha es el formato en que se reciben y devuelven las fechas: AAAA-MM-DD.
const FormatoFecha = "2006-01-02"

// Proyecto es un proyecto de software al que se le hace seguimiento (spec 01).
type Proyecto struct {
	ID          int64
	Nombre      string
	Descripcion string
	FechaInicio time.Time
	FechaFin    time.Time
}

// NuevoProyecto arma un proyecto validando las reglas de la spec 01.
// El ID queda en cero: lo asigna el sistema al guardarlo (RN-6).
func NuevoProyecto(nombre, descripcion string, inicio, fin time.Time) (Proyecto, error) {
	// RN-1: nombre obligatorio, sin largo máximo, puede repetirse.
	if strings.TrimSpace(nombre) == "" {
		return Proyecto{}, ErrNombreObligatorio
	}
	// RN-2: descripción obligatoria, hasta 500 caracteres (no bytes: cuenta bien las tildes).
	if strings.TrimSpace(descripcion) == "" {
		return Proyecto{}, ErrDescripcionObligatoria
	}
	if utf8.RuneCountInString(descripcion) > MaxDescripcion {
		return Proyecto{}, ErrDescripcionMuyLarga
	}
	// RN-3: ambas fechas obligatorias.
	if inicio.IsZero() || fin.IsZero() {
		return Proyecto{}, ErrFechasObligatorias
	}
	// RN-4: la de fin, estrictamente posterior a la de inicio.
	if !fin.After(inicio) {
		return Proyecto{}, ErrFinNoPosteriorAInicio
	}
	return Proyecto{Nombre: nombre, Descripcion: descripcion, FechaInicio: inicio, FechaFin: fin}, nil
}

// ParsearFecha convierte un texto AAAA-MM-DD en fecha (spec 01 RN-3).
// Un texto vacío devuelve la fecha cero, que NuevoProyecto trata como fecha faltante.
func ParsearFecha(texto string) (time.Time, error) {
	if texto == "" {
		return time.Time{}, nil
	}
	f, err := time.Parse(FormatoFecha, texto)
	if err != nil {
		return time.Time{}, ErrFechaInvalida
	}
	return f, nil
}
