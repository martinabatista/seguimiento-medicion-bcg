package dominio

import (
	"context"
	"strconv"
)

// ParsearHorasEstimadas convierte el texto de las horas estimadas en centésimas (spec 12 RN-2):
// mayor a 0, con hasta 2 decimales y sin máximo. Vacío es "faltan las horas estimadas".
func ParsearHorasEstimadas(texto string) (int, error) {
	if texto == "" {
		return 0, ErrHorasEstimadasObligatorias
	}
	partes := formatoHoras.FindStringSubmatch(texto)
	if partes == nil {
		return 0, ErrHorasEstimadasInvalidas
	}
	enteras, err := strconv.Atoi(partes[1])
	if err != nil || enteras > (1<<31)/centesimasPorHora { // no entra en un entero de la base
		return 0, ErrHorasEstimadasInvalidas
	}
	decimales := 0
	if partes[2] != "" {
		decimales, _ = strconv.Atoi((partes[2] + "0")[:2])
	}
	centesimas := enteras*centesimasPorHora + decimales
	if centesimas <= 0 {
		return 0, ErrHorasEstimadasInvalidas
	}
	return centesimas, nil
}

// EstimarHoras carga o cambia las horas estimadas de una historia del backlog (spec 12 CA-1).
// Si algo falla, no cambia nada.
func (s *ServicioBacklog) EstimarHoras(ctx context.Context, proyectoID, id int64, texto string) (Historia, error) {
	centesimas, err := ParsearHorasEstimadas(texto)
	if err != nil {
		return Historia{}, err
	}
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return Historia{}, err
	}
	h, err := s.historias.Obtener(ctx, proyectoID, id)
	if err != nil {
		return Historia{}, err
	}
	if h.Estado == EstadoCompletada {
		return Historia{}, ErrHistoriaCompletada
	}
	h.HorasEstimadas = &centesimas
	if err := s.historias.Actualizar(ctx, h); err != nil {
		return Historia{}, err
	}
	return h, nil
}

// QuitarHorasEstimadas deja la historia sin horas estimadas (spec 12 CA-1, RN-1).
func (s *ServicioBacklog) QuitarHorasEstimadas(ctx context.Context, proyectoID, id int64) (Historia, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return Historia{}, err
	}
	h, err := s.historias.Obtener(ctx, proyectoID, id)
	if err != nil {
		return Historia{}, err
	}
	if h.Estado == EstadoCompletada {
		return Historia{}, ErrHistoriaCompletada
	}
	h.HorasEstimadas = nil
	if err := s.historias.Actualizar(ctx, h); err != nil {
		return Historia{}, err
	}
	return h, nil
}
