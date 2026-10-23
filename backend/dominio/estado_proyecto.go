package dominio

import "time"

// EstadoProyecto es la etapa en la que está un proyecto según sus fechas (spec 03 RN-2).
// No se guarda: se calcula en cada consulta.
type EstadoProyecto string

const (
	EstadoNoIniciado EstadoProyecto = "No iniciado"
	EstadoEnCurso    EstadoProyecto = "En curso"
	EstadoFinalizado EstadoProyecto = "Finalizado"
)

// CalcularEstado compara la fecha de hoy con las del proyecto (RN-2). Inicio y fin están incluidos
// en "En curso". Solo cuenta el día de hoy, no la hora.
func CalcularEstado(p Proyecto, hoy time.Time) EstadoProyecto {
	dia := soloDia(hoy)
	switch {
	case dia.Before(soloDia(p.FechaInicio)):
		return EstadoNoIniciado
	case dia.After(soloDia(p.FechaFin)):
		return EstadoFinalizado
	}
	return EstadoEnCurso
}

// soloDia deja el año, mes y día de t (en su propia zona horaria) para comparar fechas sin hora.
func soloDia(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// ResumenBacklog es el avance del Product Backlog de un proyecto (spec 03 CA-3).
type ResumenBacklog struct {
	Total                  int
	Pendientes             int
	EnProgreso             int
	Completadas            int
	PorcentajeCompletadas  int // entero redondeado (RN-5)
	StoryPointsTotales     int // suma de las historias estimadas (RN-6)
	StoryPointsCompletados int // suma de las historias Completada estimadas (RN-6)
	SinEstimar             int
}

// ResumirBacklog cuenta las historias por estado (RN-4), calcula el porcentaje de completadas
// (RN-5) y suma los Story Points de las estimadas (RN-6).
func ResumirBacklog(historias []Historia) ResumenBacklog {
	r := ResumenBacklog{Total: len(historias)}
	for _, h := range historias {
		switch h.Estado {
		case EstadoPendiente:
			r.Pendientes++
		case EstadoEnProgreso:
			r.EnProgreso++
		case EstadoCompletada:
			r.Completadas++
		}
		if h.StoryPoints == nil {
			r.SinEstimar++
			continue
		}
		r.StoryPointsTotales += *h.StoryPoints
		if h.Estado == EstadoCompletada {
			r.StoryPointsCompletados += *h.StoryPoints
		}
	}
	r.PorcentajeCompletadas = porcentajeRedondeado(r.Completadas, r.Total)
	return r
}

// porcentajeRedondeado devuelve parte/total × 100 redondeado al entero más cercano, con la mitad
// hacia arriba (12,5 → 13), usando solo enteros. Con total 0 devuelve 0 (RN-5).
func porcentajeRedondeado(parte, total int) int {
	if total == 0 {
		return 0
	}
	return (parte*200 + total) / (2 * total)
}
