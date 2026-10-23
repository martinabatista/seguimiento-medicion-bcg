package dominio

import "testing"

func sp(n int) *int { return &n }

// spec 10 RN-1, RN-2, RN-3, RN-5 (CA-1, CA-2, CA-4, CA-6): métricas de un Sprint.
func TestCalcularMetricasSprint(t *testing.T) {
	cerrado := Sprint{ID: 7, Numero: 2, Estado: SprintCerrado, Historias: []int64{1, 2, 3}, Completadas: []int64{1, 3}}
	casos := []struct {
		nombre   string
		sprint   Sprint
		puntos   map[int64]*int
		esperado MetricasSprint
	}{
		{"Sprint cerrado con todo estimado", cerrado, map[int64]*int{1: sp(5), 2: sp(8), 3: sp(3)},
			MetricasSprint{SprintID: 7, Numero: 2, Estado: SprintCerrado, HistoriasPlanificadas: 3, HistoriasCompletadas: 2,
				PorcentajeCompletadas: 67, StoryPointsPlanificados: 16, StoryPointsCompletados: 8}},
		{"las sin estimar no suman SP pero cuentan como historias (RN-2)", cerrado, map[int64]*int{1: sp(5), 2: nil},
			MetricasSprint{SprintID: 7, Numero: 2, Estado: SprintCerrado, HistoriasPlanificadas: 3, HistoriasCompletadas: 2,
				PorcentajeCompletadas: 67, SinEstimar: 2, StoryPointsPlanificados: 5, StoryPointsCompletados: 5}},
		{"Sprint abierto: parcial (RN-3)", Sprint{ID: 8, Numero: 3, Estado: SprintAbierto, Historias: []int64{1, 2}, Completadas: []int64{2}},
			map[int64]*int{1: sp(5), 2: sp(8)},
			MetricasSprint{SprintID: 8, Numero: 3, Estado: SprintAbierto, Parcial: true, HistoriasPlanificadas: 2, HistoriasCompletadas: 1,
				PorcentajeCompletadas: 50, StoryPointsPlanificados: 13, StoryPointsCompletados: 8}},
		{"Sprint sin historias: 0 % (RN-5)", Sprint{ID: 9, Numero: 4, Estado: SprintCerrado}, nil,
			MetricasSprint{SprintID: 9, Numero: 4, Estado: SprintCerrado}},
		{"porcentaje con la mitad hacia arriba: 1 de 8 = 12,5 → 13 (RN-5)",
			Sprint{ID: 1, Numero: 1, Estado: SprintCerrado, Historias: []int64{1, 2, 3, 4, 5, 6, 7, 8}, Completadas: []int64{1}}, nil,
			MetricasSprint{SprintID: 1, Numero: 1, Estado: SprintCerrado, HistoriasPlanificadas: 8, HistoriasCompletadas: 1,
				PorcentajeCompletadas: 13, SinEstimar: 8}},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if m := CalcularMetricasSprint(c.sprint, c.puntos); m != c.esperado {
				t.Errorf("métricas = %+v\nse esperaba  %+v", m, c.esperado)
			}
		})
	}
}

// spec 10 RN-6 (CA-3): velocidad = promedio de SP completados de los Sprints cerrados, en décimas.
func TestCalcularVelocidad(t *testing.T) {
	casos := []struct {
		nombre      string
		completados []int
		decimas     int
		hay         bool
	}{
		{"sin Sprints cerrados no hay velocidad", nil, 0, false},
		{"un Sprint", []int{8}, 80, true},
		{"10 y 15 → 12,5", []int{10, 15}, 125, true},
		{"10, 15 y 12 → 12,33… → 12,3", []int{10, 15, 12}, 123, true},
		{"5 y 8 → 6,5", []int{5, 8}, 65, true},
		{"mitad hacia arriba: 1, 2, 2, 2 → 1,75 → 1,8", []int{1, 2, 2, 2}, 18, true},
		{"un Sprint sin completar nada baja la velocidad: 5 y 0 → 2,5", []int{5, 0}, 25, true},
		{"todo en 0", []int{0, 0}, 0, true},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			decimas, hay := CalcularVelocidad(c.completados)
			if decimas != c.decimas || hay != c.hay {
				t.Errorf("velocidad = %d décimas (hay = %v), se esperaba %d (hay = %v)", decimas, hay, c.decimas, c.hay)
			}
		})
	}
}

// spec 10 RN-6, RN-7 (CA-1 a CA-4): totales y velocidad del proyecto, solo con los Sprints cerrados.
func TestCalcularMetricasProyecto(t *testing.T) {
	sprints := []MetricasSprint{
		{Numero: 1, Estado: SprintCerrado, HistoriasPlanificadas: 2, HistoriasCompletadas: 1, StoryPointsPlanificados: 13, StoryPointsCompletados: 5},
		{Numero: 2, Estado: SprintCerrado, HistoriasPlanificadas: 1, HistoriasCompletadas: 1, StoryPointsPlanificados: 8, StoryPointsCompletados: 8},
		{Numero: 3, Estado: SprintAbierto, Parcial: true, HistoriasPlanificadas: 4, HistoriasCompletadas: 4, StoryPointsPlanificados: 20, StoryPointsCompletados: 20},
	}
	m := CalcularMetricasProyecto(sprints)
	esperado := TotalesSprints{HistoriasPlanificadas: 3, HistoriasCompletadas: 2, PorcentajeCompletadas: 67, StoryPointsPlanificados: 21, StoryPointsCompletados: 13}
	if m.SprintsCerrados != 2 || !m.HayVelocidad || m.VelocidadDecimas != 65 || m.Totales != esperado || len(m.Sprints) != 3 {
		t.Errorf("proyecto = %+v\nse esperaban 2 cerrados, velocidad 6,5 y totales %+v", m, esperado)
	}

	vacio := CalcularMetricasProyecto(nil)
	if vacio.HayVelocidad || vacio.SprintsCerrados != 0 || vacio.Totales != (TotalesSprints{}) || vacio.Sprints == nil {
		t.Errorf("sin Sprints = %+v, se esperaba sin velocidad, totales en 0 y lista vacía", vacio)
	}
}
