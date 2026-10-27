package dominio

import (
	"slices"
	"strings"
	"unicode/utf8"
)

// Severidad indica qué tan grave es un defecto (spec 11 RN-3).
type Severidad string

const (
	SeveridadBaja    Severidad = "Baja"
	SeveridadMedia   Severidad = "Media"
	SeveridadAlta    Severidad = "Alta"
	SeveridadCritica Severidad = "Crítica"
)

var severidades = []Severidad{SeveridadBaja, SeveridadMedia, SeveridadAlta, SeveridadCritica}

// EstadoDefecto es el seguimiento de un defecto hasta resolverlo (spec 11 RN-7).
type EstadoDefecto string

const (
	DefectoAbierto      EstadoDefecto = "Abierto"
	DefectoEnCorreccion EstadoDefecto = "En corrección"
	DefectoResuelto     EstadoDefecto = "Resuelto"
)

var estadosDefecto = []EstadoDefecto{DefectoAbierto, DefectoEnCorreccion, DefectoResuelto}

// DatosDefecto son los datos que carga la usuaria al registrar un defecto. Los identificadores en
// cero significan que no se indicaron.
type DatosDefecto struct {
	Descripcion       string
	Severidad         string
	HistoriaID        int64
	SprintDeteccionID int64
}

// Defecto es un defecto encontrado en una historia del backlog (spec 11).
type Defecto struct {
	ID                 int64
	ProyectoID         int64
	HistoriaID         int64
	Descripcion        string
	Severidad          Severidad
	Estado             EstadoDefecto
	SprintDeteccionID  int64
	SprintResolucionID int64 // 0 = sin resolver (RN-6, RN-9)
	// Se completan al consultar; no se guardan con el defecto.
	HistoriaTitulo         string
	SprintDeteccionNumero  int
	SprintResolucionNumero int // 0 = sin resolver
}

// DefectosDeSprint son los defectos detectados y resueltos en un Sprint (CA-4).
type DefectosDeSprint struct {
	SprintID   int64
	Numero     int
	Detectados int
	Resueltos  int
}

// MetricasDefectos son las cantidades del proyecto y de cada Sprint (CA-4).
type MetricasDefectos struct {
	Detectados int
	Resueltos  int
	Sprints    []DefectosDeSprint
}

// NuevoDefecto valida los datos en el orden de RN-13: descripción, severidad, historia y Sprint de
// detección. Que la historia y el Sprint existan en el proyecto lo verifica el servicio. Arranca
// Abierto y sin Sprint de resolución (RN-6); el ID lo asigna el sistema al guardarlo (RN-12).
func NuevoDefecto(proyectoID int64, d DatosDefecto) (Defecto, error) {
	// RN-2: obligatoria y hasta 1000 caracteres (no bytes), como la descripción de una historia.
	if strings.TrimSpace(d.Descripcion) == "" {
		return Defecto{}, ErrDescripcionObligatoria
	}
	if utf8.RuneCountInString(d.Descripcion) > MaxDescripcionHistoria {
		return Defecto{}, ErrDescripcionHistoriaMuyLarga
	}
	// RN-3: escrita exactamente como una de las cuatro.
	if d.Severidad == "" {
		return Defecto{}, ErrSeveridadObligatoria
	}
	severidad := Severidad(d.Severidad)
	if !slices.Contains(severidades, severidad) {
		return Defecto{}, ErrSeveridadInvalida
	}
	if d.HistoriaID <= 0 {
		return Defecto{}, ErrHistoriaDefectoObligatoria
	}
	if d.SprintDeteccionID <= 0 {
		return Defecto{}, ErrSprintDeteccionObligatorio
	}
	return Defecto{
		ProyectoID:        proyectoID,
		HistoriaID:        d.HistoriaID,
		Descripcion:       d.Descripcion,
		Severidad:         severidad,
		Estado:            DefectoAbierto,
		SprintDeteccionID: d.SprintDeteccionID,
	}, nil
}

// ParsearEstadoDefecto acepta solo los tres estados, escritos exactamente así (RN-7).
func ParsearEstadoDefecto(texto string) (EstadoDefecto, error) {
	estado := EstadoDefecto(texto)
	if !slices.Contains(estadosDefecto, estado) {
		return "", ErrEstadoInvalido
	}
	return estado, nil
}

// CambiarEstadoDefecto pasa el defecto al estado nuevo; se puede pasar de cualquiera a cualquiera
// (RN-7). Para Resuelto hace falta el Sprint de resolución, que no puede empezar antes que el de
// detección (RN-8: se compara por fecha de inicio). Para los otros dos estados el defecto queda
// sin Sprint de resolución y el indicado no se tiene en cuenta (RN-9). No modifica el recibido.
func CambiarEstadoDefecto(d Defecto, estado EstadoDefecto, deteccion Sprint, resolucion *Sprint) (Defecto, error) {
	if estado != DefectoResuelto {
		d.Estado = estado
		d.SprintResolucionID = 0
		return d, nil
	}
	if resolucion == nil {
		return Defecto{}, ErrSprintResolucionObligatorio
	}
	if resolucion.FechaInicio.Before(deteccion.FechaInicio) {
		return Defecto{}, ErrResolucionAnteriorADeteccion
	}
	d.Estado = DefectoResuelto
	d.SprintResolucionID = resolucion.ID
	return d, nil
}

// ContarDefectos calcula los detectados y resueltos del proyecto y de cada Sprint (RN-10, RN-11):
// un defecto cuenta como detectado en su Sprint de detección cualquiera sea su estado, y como
// resuelto solo si está Resuelto, en su Sprint de resolución. Incluye todos los Sprints, ordenados
// por fecha de inicio, con 0 en los que no tienen defectos. No modifica los slices recibidos.
func ContarDefectos(sprints []Sprint, defectos []Defecto) MetricasDefectos {
	ordenados := slices.Clone(sprints)
	slices.SortFunc(ordenados, func(a, b Sprint) int { return a.FechaInicio.Compare(b.FechaInicio) })

	m := MetricasDefectos{Sprints: make([]DefectosDeSprint, 0, len(ordenados))}
	indice := make(map[int64]int, len(ordenados))
	for i, s := range ordenados {
		indice[s.ID] = i
		m.Sprints = append(m.Sprints, DefectosDeSprint{SprintID: s.ID, Numero: s.Numero})
	}
	for _, d := range defectos {
		m.Detectados++
		if i, ok := indice[d.SprintDeteccionID]; ok {
			m.Sprints[i].Detectados++
		}
		if d.Estado != DefectoResuelto {
			continue
		}
		m.Resueltos++
		if i, ok := indice[d.SprintResolucionID]; ok {
			m.Sprints[i].Resueltos++
		}
	}
	return m
}
