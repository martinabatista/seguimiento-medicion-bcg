package dominio

import "context"

// Estimar asigna o cambia los Story Points de una historia del backlog (spec 07 CA-1, CA-2).
// El valor tiene que ser de la escala (RN-2). Si algo falla, no cambia nada.
func (s *ServicioBacklog) Estimar(ctx context.Context, proyectoID, id int64, storyPoints int) (Historia, error) {
	if !EnEscala(storyPoints) {
		return Historia{}, ErrStoryPointsFueraDeEscala
	}
	return s.cambiarSiNoEstaCompletada(ctx, proyectoID, id, func(h *Historia) { h.StoryPoints = &storyPoints })
}

// QuitarEstimacion deja la historia sin estimar (spec 07 CA-2, RN-3). Si ya estaba sin estimar, no
// es un error.
func (s *ServicioBacklog) QuitarEstimacion(ctx context.Context, proyectoID, id int64) (Historia, error) {
	return s.cambiarSiNoEstaCompletada(ctx, proyectoID, id, func(h *Historia) { h.StoryPoints = nil })
}

// cambiarSiNoEstaCompletada aplica una estimación (Story Points, spec 07, u horas estimadas, spec 12)
// a una historia del proyecto. Verifica que el proyecto y la historia existan (spec 07 CA-6, RN-6) y
// que la historia no esté completada (spec 07 RN-4, spec 12 RN-3). Solo cambia lo que toca cambio:
// el resto de la historia queda igual (spec 07 RN-5, spec 12 RN-4).
func (s *ServicioBacklog) cambiarSiNoEstaCompletada(ctx context.Context, proyectoID, id int64, cambio func(*Historia)) (Historia, error) {
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
	cambio(&h)
	if err := s.historias.Actualizar(ctx, h); err != nil {
		return Historia{}, err
	}
	return h, nil
}
