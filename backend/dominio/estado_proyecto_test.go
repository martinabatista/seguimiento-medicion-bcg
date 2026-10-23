package dominio

import (
	"testing"
	"time"
)

// spec 03 RN-2 (CA-1): el estado se calcula comparando la fecha de hoy con las del proyecto.
func TestCalcularEstado(t *testing.T) {
	p := Proyecto{Nombre: "Tienda Online", FechaInicio: fecha(t, "2026-11-01"), FechaFin: fecha(t, "2026-12-15")}
	casos := []struct {
		nombre string
		hoy    time.Time
		estado EstadoProyecto
	}{
		{"mucho antes del inicio", fecha(t, "2026-10-15"), EstadoNoIniciado},
		{"el día anterior al inicio", fecha(t, "2026-10-31"), EstadoNoIniciado},
		{"el mismo día del inicio", fecha(t, "2026-11-01"), EstadoEnCurso},
		{"entre inicio y fin", fecha(t, "2026-11-20"), EstadoEnCurso},
		{"el mismo día del fin", fecha(t, "2026-12-15"), EstadoEnCurso},
		{"el día siguiente al fin", fecha(t, "2026-12-16"), EstadoFinalizado},
		{"mucho después del fin", fecha(t, "2027-03-01"), EstadoFinalizado},
		// La hora no importa: solo cuenta el día (spec 03 Restricciones).
		{"último minuto del día de fin", time.Date(2026, 12, 15, 23, 59, 0, 0, time.UTC), EstadoEnCurso},
		{"día de fin en otra zona horaria", time.Date(2026, 12, 15, 22, 0, 0, 0, time.FixedZone("ART", -3*60*60)), EstadoEnCurso},
		{"primer minuto del día de inicio", time.Date(2026, 11, 1, 0, 1, 0, 0, time.UTC), EstadoEnCurso},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := CalcularEstado(p, c.hoy); got != c.estado {
				t.Errorf("estado = %q, se esperaba %q", got, c.estado)
			}
		})
	}
}

// historiaCon arma una historia del backlog en un estado, con Story Points o sin estimar (sp < 0).
func historiaCon(estado EstadoHistoria, sp int) Historia {
	h := Historia{Titulo: "Historia", Estado: estado}
	if sp >= 0 {
		h.StoryPoints = &sp
	}
	return h
}

// spec 03 RN-4 a RN-6 (CA-3): cantidades por estado, porcentaje redondeado y Story Points.
func TestResumirBacklog(t *testing.T) {
	const sinEstimar = -1
	repetir := func(n int, h Historia) []Historia {
		hs := make([]Historia, n)
		for i := range hs {
			hs[i] = h
		}
		return hs
	}
	completada := historiaCon(EstadoCompletada, sinEstimar)
	pendiente := historiaCon(EstadoPendiente, sinEstimar)

	casos := []struct {
		nombre    string
		historias []Historia
		esperado  ResumenBacklog
	}{
		{"backlog vacío (spec 03 RN-5)", nil, ResumenBacklog{}},
		{"historias en todos los estados (spec 03 RN-4 a RN-6)", []Historia{
			historiaCon(EstadoCompletada, 5),
			historiaCon(EstadoCompletada, 3),
			historiaCon(EstadoEnProgreso, 8),
			historiaCon(EstadoPendiente, sinEstimar),
		}, ResumenBacklog{Total: 4, Pendientes: 1, EnProgreso: 1, Completadas: 2, PorcentajeCompletadas: 50,
			StoryPointsTotales: 16, StoryPointsCompletados: 8, SinEstimar: 1}},
		{"sin estimar no suman Story Points (spec 03 RN-6)", []Historia{completada, pendiente},
			ResumenBacklog{Total: 2, Pendientes: 1, Completadas: 1, PorcentajeCompletadas: 50, SinEstimar: 2}},
		{"una completada sin estimar no suma a los completados (spec 03 RN-6)", []Historia{
			historiaCon(EstadoCompletada, sinEstimar),
			historiaCon(EstadoPendiente, 5),
		}, ResumenBacklog{Total: 2, Pendientes: 1, Completadas: 1, PorcentajeCompletadas: 50,
			StoryPointsTotales: 5, SinEstimar: 1}},
		{"Story Points en cero cuentan como estimada (spec 03 RN-6)", []Historia{historiaCon(EstadoCompletada, 0)},
			ResumenBacklog{Total: 1, Completadas: 1, PorcentajeCompletadas: 100}},
		{"1 de 3 completadas: 33 (spec 03 RN-5)", append(repetir(2, pendiente), completada),
			ResumenBacklog{Total: 3, Pendientes: 2, Completadas: 1, PorcentajeCompletadas: 33, SinEstimar: 3}},
		{"2 de 3 completadas: 67 (spec 03 RN-5)", append(repetir(2, completada), pendiente),
			ResumenBacklog{Total: 3, Pendientes: 1, Completadas: 2, PorcentajeCompletadas: 67, SinEstimar: 3}},
		{"1 de 8 completadas: 12,5 redondea a 13 (spec 03 RN-5)", append(repetir(7, pendiente), completada),
			ResumenBacklog{Total: 8, Pendientes: 7, Completadas: 1, PorcentajeCompletadas: 13, SinEstimar: 8}},
		{"todas completadas: 100 (spec 03 RN-5)", repetir(4, completada),
			ResumenBacklog{Total: 4, Completadas: 4, PorcentajeCompletadas: 100, SinEstimar: 4}},
		{"ninguna completada: 0 (spec 03 RN-5)", repetir(4, pendiente),
			ResumenBacklog{Total: 4, Pendientes: 4, SinEstimar: 4}},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := ResumirBacklog(c.historias); got != c.esperado {
				t.Errorf("resumen = %+v, se esperaba %+v", got, c.esperado)
			}
		})
	}
}
