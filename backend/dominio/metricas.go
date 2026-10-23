package dominio

// MetricasSprint son las métricas de avance de un Sprint (spec 10 CA-1, CA-2, CA-4, CA-6).
type MetricasSprint struct {
	SprintID                int64
	Numero                  int
	Estado                  EstadoSprint
	Parcial                 bool // el Sprint está abierto: sus números todavía pueden cambiar (RN-3)
	HistoriasPlanificadas   int
	HistoriasCompletadas    int
	PorcentajeCompletadas   int // entero redondeado (RN-5)
	SinEstimar              int // historias planificadas sin Story Points (RN-2)
	StoryPointsPlanificados int
	StoryPointsCompletados  int
}

// TotalesSprints suma las métricas de los Sprints cerrados de un proyecto (spec 10 RN-7).
type TotalesSprints struct {
	HistoriasPlanificadas   int
	HistoriasCompletadas    int
	PorcentajeCompletadas   int
	StoryPointsPlanificados int
	StoryPointsCompletados  int
}

// MetricasProyecto son las métricas de avance de un proyecto (spec 10).
type MetricasProyecto struct {
	Sprints          []MetricasSprint // ordenados por fecha de inicio
	SprintsCerrados  int
	HayVelocidad     bool // false si no hay Sprints cerrados (RN-6)
	VelocidadDecimas int  // velocidad × 10: 65 es 6,5 Story Points por Sprint
	Totales          TotalesSprints
}

// CalcularMetricasSprint calcula las métricas de un Sprint. puntos tiene los Story Points de cada
// historia (nil = sin estimar): los actuales si el Sprint está abierto y los congelados al cerrarlo
// si está cerrado (RN-3, RN-4). Una historia que no está en puntos cuenta como sin estimar.
func CalcularMetricasSprint(s Sprint, puntos map[int64]*int) MetricasSprint {
	m := MetricasSprint{
		SprintID:              s.ID,
		Numero:                s.Numero,
		Estado:                s.Estado,
		Parcial:               s.Estado == SprintAbierto,
		HistoriasPlanificadas: len(s.Historias),
		HistoriasCompletadas:  len(s.Completadas),
	}
	for _, id := range s.Historias {
		p := puntos[id]
		if p == nil {
			m.SinEstimar++
			continue
		}
		m.StoryPointsPlanificados += *p
		if s.EstaCompletada(id) {
			m.StoryPointsCompletados += *p
		}
	}
	m.PorcentajeCompletadas = porcentajeRedondeado(m.HistoriasCompletadas, m.HistoriasPlanificadas)
	return m
}

// CalcularVelocidad es el promedio de los Story Points completados en cada Sprint cerrado, en
// décimas y con la mitad hacia arriba (RN-6): 10 y 15 → 125 (12,5). Solo usa enteros, sin errores de
// coma flotante. Sin Sprints cerrados no hay velocidad.
func CalcularVelocidad(completadosPorSprint []int) (decimas int, hay bool) {
	n := len(completadosPorSprint)
	if n == 0 {
		return 0, false
	}
	total := 0
	for _, c := range completadosPorSprint {
		total += c
	}
	return (total*20 + n) / (2 * n), true
}

// CalcularMetricasProyecto arma las métricas del proyecto a partir de las de cada Sprint: la
// velocidad y los totales usan solo los Sprints cerrados (RN-6, RN-7).
func CalcularMetricasProyecto(sprints []MetricasSprint) MetricasProyecto {
	m := MetricasProyecto{Sprints: append([]MetricasSprint{}, sprints...)}
	completados := []int{}
	for _, s := range sprints {
		if s.Estado != SprintCerrado {
			continue
		}
		m.SprintsCerrados++
		completados = append(completados, s.StoryPointsCompletados)
		m.Totales.HistoriasPlanificadas += s.HistoriasPlanificadas
		m.Totales.HistoriasCompletadas += s.HistoriasCompletadas
		m.Totales.StoryPointsPlanificados += s.StoryPointsPlanificados
		m.Totales.StoryPointsCompletados += s.StoryPointsCompletados
	}
	m.Totales.PorcentajeCompletadas = porcentajeRedondeado(m.Totales.HistoriasCompletadas, m.Totales.HistoriasPlanificadas)
	m.VelocidadDecimas, m.HayVelocidad = CalcularVelocidad(completados)
	return m
}
