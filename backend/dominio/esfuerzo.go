package dominio

import (
	"cmp"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Límites del registro de esfuerzo (spec 09 RN-4, RN-5). Las horas se manejan en centésimas
// enteras (1,5 h = 150): así las sumas son exactas (0,1 + 0,2 = 0,3, spec 09 Restricciones).
const (
	MaxActividad           = 300
	MaxCentesimasRegistro  = 24 * 100
	centesimasPorHora      = 100
	maxDigitosEnterosHoras = 2 // "24" es el máximo: más dígitos ya es inválido
)

// formatoHoras acepta un número sin signo con hasta 2 decimales: "1", "2.5", "0.25" (RN-5).
var formatoHoras = regexp.MustCompile(`^(\d+)(?:\.(\d{1,2}))?$`)

// DatosEsfuerzo son los datos que carga la usuaria al registrar esfuerzo. La fecha y las horas
// llegan como texto para que el dominio distinga "falta" de "no es válido".
type DatosEsfuerzo struct {
	IntegranteID int64 // 0 = no se indicó
	Fecha        string
	Actividad    string
	Horas        string
}

// RegistroEsfuerzo es un registro de horas trabajadas sobre una historia (spec 09).
type RegistroEsfuerzo struct {
	ID               int64
	ProyectoID       int64
	HistoriaID       int64
	IntegranteID     int64
	IntegranteNombre string // se completa al consultar; no se guarda con el registro
	Fecha            time.Time
	Actividad        string
	Centesimas       int // horas × 100
}

// EsfuerzoDeHistoria es lo que devuelve la consulta por historia (CA-4).
type EsfuerzoDeHistoria struct {
	Registros       []RegistroEsfuerzo
	TotalCentesimas int
}

// NuevoRegistroEsfuerzo valida los datos en el orden de RN-9: integrante, fecha, actividad y horas.
// La existencia del integrante en el proyecto la verifica el servicio. El ID queda en cero: lo
// asigna el sistema al guardarlo (RN-8).
func NuevoRegistroEsfuerzo(p Proyecto, historiaID int64, d DatosEsfuerzo, hoy time.Time) (RegistroEsfuerzo, error) {
	if d.IntegranteID <= 0 {
		return RegistroEsfuerzo{}, ErrIntegranteObligatorio
	}
	fecha, err := validarFechaEsfuerzo(p, d.Fecha, hoy)
	if err != nil {
		return RegistroEsfuerzo{}, err
	}
	// RN-4: actividad obligatoria, hasta 300 caracteres (no bytes).
	if strings.TrimSpace(d.Actividad) == "" {
		return RegistroEsfuerzo{}, ErrActividadObligatoria
	}
	if utf8.RuneCountInString(d.Actividad) > MaxActividad {
		return RegistroEsfuerzo{}, ErrActividadMuyLarga
	}
	centesimas, err := ParsearHoras(d.Horas)
	if err != nil {
		return RegistroEsfuerzo{}, err
	}
	return RegistroEsfuerzo{
		ProyectoID:   p.ID,
		HistoriaID:   historiaID,
		IntegranteID: d.IntegranteID,
		Fecha:        fecha,
		Actividad:    d.Actividad,
		Centesimas:   centesimas,
	}, nil
}

// validarFechaEsfuerzo aplica RN-3: obligatoria, válida, dentro de las fechas del proyecto (ambas
// incluidas) y no posterior a hoy. Si está fuera del proyecto y además en el futuro, se informa
// primero que está fuera del proyecto.
func validarFechaEsfuerzo(p Proyecto, texto string, hoy time.Time) (time.Time, error) {
	if texto == "" {
		return time.Time{}, ErrFechaObligatoria
	}
	fecha, err := ParsearFecha(texto)
	if err != nil {
		return time.Time{}, err
	}
	dia := soloDia(fecha)
	if dia.Before(soloDia(p.FechaInicio)) || dia.After(soloDia(p.FechaFin)) {
		return time.Time{}, ErrFechaFueraDelProyecto
	}
	if dia.After(soloDia(hoy)) {
		return time.Time{}, ErrFechaFutura
	}
	return fecha, nil
}

// ParsearHoras convierte el texto de las horas en centésimas (RN-5): mayor a 0, hasta 24 y con
// hasta 2 decimales. Vacío es "faltan las horas"; cualquier otro texto inválido es ErrHorasInvalidas.
func ParsearHoras(texto string) (int, error) {
	if texto == "" {
		return 0, ErrHorasObligatorias
	}
	centesimas, ok := centesimasDeTexto(texto, maxDigitosEnterosHoras)
	if !ok || centesimas <= 0 || centesimas > MaxCentesimasRegistro {
		return 0, ErrHorasInvalidas
	}
	return centesimas, nil
}

// centesimasDeTexto convierte un número sin signo con hasta 2 decimales en centésimas ("2.5" → 250).
// Devuelve false si no tiene ese formato o si la parte entera tiene más de maxDigitos dígitos. Lo
// usan las horas del esfuerzo (spec 09) y las horas estimadas (spec 12).
func centesimasDeTexto(texto string, maxDigitos int) (int, bool) {
	partes := formatoHoras.FindStringSubmatch(texto)
	if partes == nil || len(partes[1]) > maxDigitos {
		return 0, false
	}
	enteras, _ := strconv.Atoi(partes[1]) // el formato y el largo garantizan un entero válido
	decimales := 0
	if partes[2] != "" {
		decimales, _ = strconv.Atoi((partes[2] + "0")[:2]) // "5" son 50 centésimas
	}
	return enteras*centesimasPorHora + decimales, true
}

// HorasATexto muestra las centésimas como horas con hasta 2 decimales y sin ceros de más:
// 150 → "1.5", 375 → "3.75", 100 → "1" (RN-6).
func HorasATexto(centesimas int) string {
	texto := strconv.Itoa(centesimas/centesimasPorHora) + "." +
		strconv.Itoa(100 + centesimas%centesimasPorHora)[1:]
	return strings.TrimSuffix(strings.TrimRight(texto, "0"), ".")
}

// ResumirEsfuerzo ordena los registros por fecha y, a igual fecha, por orden de carga (el ID es
// creciente) (RN-7), y suma el total exacto en centésimas (RN-6). No modifica el slice recibido.
func ResumirEsfuerzo(registros []RegistroEsfuerzo) EsfuerzoDeHistoria {
	ordenados := slices.Clone(registros)
	if ordenados == nil {
		ordenados = []RegistroEsfuerzo{}
	}
	slices.SortFunc(ordenados, func(a, b RegistroEsfuerzo) int {
		return cmp.Or(a.Fecha.Compare(b.Fecha), cmp.Compare(a.ID, b.ID))
	})
	total := 0
	for _, r := range ordenados {
		total += r.Centesimas
	}
	return EsfuerzoDeHistoria{Registros: ordenados, TotalCentesimas: total}
}
