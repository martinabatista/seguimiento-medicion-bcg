package dominio

// Desviacion compara las horas reales con las estimadas (spec 12 RN-6): positiva = se trabajó más
// de lo estimado; negativa = menos.
type Desviacion struct {
	Centesimas int // reales − estimadas, en centésimas de hora
	Porcentaje int // sobre lo estimado, entero redondeado con la mitad alejándose de cero
}

// EsfuerzoHistoria son las métricas de esfuerzo de una historia del backlog (spec 12 CA-3, CA-4).
type EsfuerzoHistoria struct {
	HistoriaID int64
	Titulo     string
	Estado     EstadoHistoria
	Estimadas  *int        // centésimas; nil = sin horas estimadas
	Reales     int         // centésimas: el total de su esfuerzo registrado (RN-5)
	Desviacion *Desviacion // nil = no entra en la desviación (RN-7)
}

// EsfuerzoSprint son las métricas de esfuerzo de las historias completadas en un Sprint (RN-8).
type EsfuerzoSprint struct {
	SprintID              int64
	Numero                int
	Estado                EstadoSprint
	Parcial               bool // Sprint abierto
	Estimadas             int
	Reales                int
	Desviacion            *Desviacion // nil = "sin datos" (RN-10)
	CompletadasSinEstimar int
}

// MetricasEsfuerzo son las métricas de esfuerzo de un proyecto (spec 12).
type MetricasEsfuerzo struct {
	Historias             []EsfuerzoHistoria // en el orden del backlog
	Sprints               []EsfuerzoSprint   // por fecha de inicio
	Estimadas             int                // todas las historias con horas estimadas (RN-9)
	Reales                int                // todo el esfuerzo registrado (RN-9)
	Desviacion            *Desviacion        // de las completadas con horas estimadas; nil = "sin datos"
	CompletadasSinEstimar int
}

// CalcularDesviacion compara las horas reales con las estimadas, en centésimas (RN-6). Las
// estimadas son mayores a 0 (RN-2), así que el porcentaje siempre se puede calcular.
func CalcularDesviacion(estimadas, reales int) Desviacion {
	diferencia := reales - estimadas
	porcentaje := porcentajeRedondeado(diferencia, estimadas)
	if diferencia < 0 {
		porcentaje = -porcentajeRedondeado(-diferencia, estimadas)
	}
	return Desviacion{Centesimas: diferencia, Porcentaje: porcentaje}
}

// CalcularMetricasEsfuerzo arma las métricas de esfuerzo por historia, por Sprint y del proyecto.
// reales tiene las horas reales de cada historia en centésimas (falta = sin esfuerzo registrado).
func CalcularMetricasEsfuerzo(historias []Historia, sprints []Sprint, reales map[int64]int) MetricasEsfuerzo {
	m := MetricasEsfuerzo{Historias: []EsfuerzoHistoria{}, Sprints: []EsfuerzoSprint{}}
	porID := make(map[int64]Historia, len(historias))
	estimadasCompletadas, realesCompletadas := 0, 0
	for _, h := range OrdenarBacklog(historias) {
		porID[h.ID] = h
		e := EsfuerzoHistoria{HistoriaID: h.ID, Titulo: h.Titulo, Estado: h.Estado, Estimadas: h.HorasEstimadas, Reales: reales[h.ID]}
		if h.HorasEstimadas != nil {
			m.Estimadas += *h.HorasEstimadas
		}
		m.Reales += e.Reales
		if h.Estado == EstadoCompletada {
			if h.HorasEstimadas == nil {
				m.CompletadasSinEstimar++
			} else {
				d := CalcularDesviacion(*h.HorasEstimadas, e.Reales)
				e.Desviacion = &d
				estimadasCompletadas += *h.HorasEstimadas
				realesCompletadas += e.Reales
			}
		}
		m.Historias = append(m.Historias, e)
	}
	if estimadasCompletadas > 0 {
		d := CalcularDesviacion(estimadasCompletadas, realesCompletadas)
		m.Desviacion = &d
	}

	for _, s := range ordenarPorInicio(sprints) {
		e := EsfuerzoSprint{SprintID: s.ID, Numero: s.Numero, Estado: s.Estado, Parcial: s.Estado == SprintAbierto}
		for _, id := range s.Completadas {
			h, ok := porID[id]
			if !ok {
				continue
			}
			if h.HorasEstimadas == nil {
				e.CompletadasSinEstimar++
				continue
			}
			e.Estimadas += *h.HorasEstimadas
			e.Reales += reales[id]
		}
		if e.Estimadas > 0 {
			d := CalcularDesviacion(e.Estimadas, e.Reales)
			e.Desviacion = &d
		}
		m.Sprints = append(m.Sprints, e)
	}
	return m
}
