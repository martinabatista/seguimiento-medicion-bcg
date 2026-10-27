package dominio

import "context"

// maxDigitosEnterosEstimadas no es un máximo de negocio (spec 12 RN-2: sin máximo): 7 dígitos
// (9.999.999 h) es lo que entra en centésimas en un entero de la base.
const maxDigitosEnterosEstimadas = 7

// ParsearHorasEstimadas convierte el texto de las horas estimadas en centésimas (spec 12 RN-2):
// mayor a 0, con hasta 2 decimales y sin máximo. Vacío es "faltan las horas estimadas".
func ParsearHorasEstimadas(texto string) (int, error) {
	if texto == "" {
		return 0, ErrHorasEstimadasObligatorias
	}
	centesimas, ok := centesimasDeTexto(texto, maxDigitosEnterosEstimadas)
	if !ok || centesimas <= 0 {
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
	return s.cambiarSiNoEstaCompletada(ctx, proyectoID, id, func(h *Historia) { h.HorasEstimadas = &centesimas })
}

// QuitarHorasEstimadas deja la historia sin horas estimadas (spec 12 CA-1, RN-1).
func (s *ServicioBacklog) QuitarHorasEstimadas(ctx context.Context, proyectoID, id int64) (Historia, error) {
	return s.cambiarSiNoEstaCompletada(ctx, proyectoID, id, func(h *Historia) { h.HorasEstimadas = nil })
}
