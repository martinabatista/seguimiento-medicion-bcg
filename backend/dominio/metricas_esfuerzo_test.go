package dominio_test

import (
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// spec 12 RN-6: desviación en centésimas (reales − estimadas) y en % con la mitad alejándose de cero.
func TestCalcularDesviacion(t *testing.T) {
	casos := []struct {
		estimadas, reales int
		centesimas        int
		porcentaje        int
	}{
		{800, 1000, 200, 25},  // se trabajó más de lo estimado
		{800, 600, -200, -25}, // menos
		{800, 800, 0, 0},      // justo
		{300, 350, 50, 17},    // 16,67 % → 17
		{800, 0, -800, -100},  // completada sin esfuerzo registrado
		{800, 900, 100, 13},   // +12,5 % → +13
		{800, 700, -100, -13}, // −12,5 % → −13: la mitad se aleja de cero
		{300, 301, 1, 0},      // 0,33 % → 0
		{100, 50050, 49950, 49950},
	}
	for _, c := range casos {
		d := dominio.CalcularDesviacion(c.estimadas, c.reales)
		if d.Centesimas != c.centesimas || d.Porcentaje != c.porcentaje {
			t.Errorf("CalcularDesviacion(%d, %d) = %+v; se esperaba %d centésimas y %d %%", c.estimadas, c.reales, d, c.centesimas, c.porcentaje)
		}
	}
}

func horas(centesimas int) *int { return &centesimas }

// spec 12 CA-3, CA-4, CA-5, RN-5 a RN-10: por historia, por Sprint y del proyecto.
func TestCalcularMetricasEsfuerzo(t *testing.T) {
	// En desorden, para verificar que se ordenan como el backlog (spec 04 RN-8).
	historias := []dominio.Historia{
		{ID: 3, Titulo: "Pagar pedido", Prioridad: 3, Estado: dominio.EstadoEnProgreso, HorasEstimadas: horas(500)},
		{ID: 1, Titulo: "Registrar usuario", Prioridad: 1, Estado: dominio.EstadoCompletada, HorasEstimadas: horas(800)},
		{ID: 2, Titulo: "Ver catálogo", Prioridad: 2, Estado: dominio.EstadoCompletada, HorasEstimadas: horas(400)},
		{ID: 4, Titulo: "Ver ofertas", Prioridad: 4, Estado: dominio.EstadoCompletada},
	}
	sprints := []dominio.Sprint{
		{ID: 20, Numero: 2, Estado: dominio.SprintAbierto, FechaInicio: f("2026-11-15"), Historias: []int64{3, 4}, Completadas: []int64{4}},
		{ID: 10, Numero: 1, Estado: dominio.SprintCerrado, FechaInicio: f("2026-11-01"), Historias: []int64{1, 2, 3}, Completadas: []int64{1, 2}},
	}
	reales := map[int64]int{1: 1000, 2: 200, 3: 700} // la 4 no tiene esfuerzo registrado

	m := dominio.CalcularMetricasEsfuerzo(historias, sprints, reales)

	// Por historia, en el orden del backlog.
	if len(m.Historias) != 4 {
		t.Fatalf("historias = %+v, se esperaban 4", m.Historias)
	}
	h := m.Historias[0]
	if h.HistoriaID != 1 || h.Titulo != "Registrar usuario" || !iguales(h.Estimadas, 800) || h.Reales != 1000 ||
		h.Desviacion == nil || *h.Desviacion != (dominio.Desviacion{Centesimas: 200, Porcentaje: 25}) {
		t.Errorf("historia 1 = %+v, se esperaba 8 h estimadas, 10 h reales y +2 h / +25 %%", h)
	}
	// RN-7: una historia sin terminar muestra sus horas, pero no tiene desviación.
	if h := m.Historias[2]; h.HistoriaID != 3 || !iguales(h.Estimadas, 500) || h.Reales != 700 || h.Desviacion != nil {
		t.Errorf("historia 3 = %+v, se esperaba 5 h estimadas, 7 h reales y sin desviación", h)
	}
	// Completada sin horas estimadas: sin desviación y sin esfuerzo, 0 horas reales.
	if h := m.Historias[3]; h.HistoriaID != 4 || h.Estimadas != nil || h.Reales != 0 || h.Desviacion != nil {
		t.Errorf("historia 4 = %+v, se esperaba sin horas estimadas, 0 h reales y sin desviación", h)
	}

	// Por Sprint, por fecha de inicio. RN-8: sobre las sumas de las completadas en el Sprint.
	if len(m.Sprints) != 2 || m.Sprints[0].SprintID != 10 || m.Sprints[1].SprintID != 20 {
		t.Fatalf("sprints = %+v, se esperaban el 10 y el 20 en ese orden", m.Sprints)
	}
	s1 := m.Sprints[0]
	if s1.Numero != 1 || s1.Parcial || s1.Estimadas != 1200 || s1.Reales != 1200 || s1.CompletadasSinEstimar != 0 ||
		s1.Desviacion == nil || *s1.Desviacion != (dominio.Desviacion{}) {
		t.Errorf("Sprint 1 = %+v, se esperaba 12 h estimadas, 12 h reales y 0 h / 0 %% (se compensan)", s1)
	}
	// RN-10: el Sprint 2 solo completó una historia sin horas estimadas: "sin datos" y se informa aparte.
	s2 := m.Sprints[1]
	if !s2.Parcial || s2.Estimadas != 0 || s2.Reales != 0 || s2.Desviacion != nil || s2.CompletadasSinEstimar != 1 {
		t.Errorf("Sprint 2 = %+v, se esperaba parcial, sin datos y 1 completada sin horas estimadas", s2)
	}

	// RN-9: del proyecto, estimadas de todo el backlog y reales de todo el esfuerzo; la desviación,
	// de las completadas con horas estimadas (1 y 2: 12 h estimadas, 12 h reales).
	if m.Estimadas != 1700 || m.Reales != 1900 || m.CompletadasSinEstimar != 1 ||
		m.Desviacion == nil || *m.Desviacion != (dominio.Desviacion{}) {
		t.Errorf("proyecto = estimadas %d, reales %d, desviación %+v, completadas sin estimar %d; se esperaba 1700, 1900, 0 y 1",
			m.Estimadas, m.Reales, m.Desviacion, m.CompletadasSinEstimar)
	}
}

// spec 12 RN-10: sin historias ni Sprints, todo en 0 y la desviación "sin datos" (nil, no 0).
func TestCalcularMetricasEsfuerzoVacio(t *testing.T) {
	m := dominio.CalcularMetricasEsfuerzo(nil, nil, nil)
	if len(m.Historias) != 0 || len(m.Sprints) != 0 || m.Estimadas != 0 || m.Reales != 0 || m.Desviacion != nil {
		t.Errorf("métricas = %+v, se esperaban vacías y sin desviación", m)
	}
	if m.Historias == nil || m.Sprints == nil {
		t.Errorf("las listas vacías tienen que ser [] y no nil, para que la API responda listas")
	}
}

func iguales(p *int, v int) bool { return p != nil && *p == v }
