package dominio

import "context"

// Estimar asigna o cambia los Story Points de una historia del backlog (spec 07 CA-1, CA-2).
// El valor tiene que ser de la escala (RN-2). Si algo falla, no cambia nada.
func (s *ServicioBacklog) Estimar(ctx context.Context, proyectoID, id int64, storyPoints int) (Historia, error) {
	if !EnEscala(storyPoints) {
		return Historia{}, ErrStoryPointsFueraDeEscala
	}
	return s.cambiarEstimacion(ctx, proyectoID, id, &storyPoints)
}

// QuitarEstimacion deja la historia sin estimar (spec 07 CA-2, RN-3). Si ya estaba sin estimar, no
// es un error.
func (s *ServicioBacklog) QuitarEstimacion(ctx context.Context, proyectoID, id int64) (Historia, error) {
	return s.cambiarEstimacion(ctx, proyectoID, id, nil)
}

// cambiarEstimacion guarda los Story Points (nil = sin estimar) de una historia del proyecto.
// Verifica que el proyecto y la historia existan (CA-6, RN-6) y que la historia no esté completada
// (RN-4). Solo cambian los Story Points: el resto de la historia queda igual (RN-5).
func (s *ServicioBacklog) cambiarEstimacion(ctx context.Context, proyectoID, id int64, storyPoints *int) (Historia, error) {
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
	h.StoryPoints = storyPoints
	if err := s.historias.Actualizar(ctx, h); err != nil {
		return Historia{}, err
	}
	return h, nil
}
