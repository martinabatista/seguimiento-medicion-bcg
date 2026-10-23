package dominio

import "context"

// Estimar asigna o cambia los Story Points de una historia del backlog (spec 07 CA-1, CA-2).
// Primero verifica que el proyecto y la historia existan (CA-6); después, que el valor sea de la
// escala (RN-2) y que la historia no esté completada (RN-4). Si algo falla, no cambia nada.
// Solo cambian los Story Points: el resto de la historia queda igual (RN-5).
func (s *ServicioBacklog) Estimar(ctx context.Context, proyectoID, id int64, storyPoints int) (Historia, error) {
	h, err := s.historiaDelProyecto(ctx, proyectoID, id)
	if err != nil {
		return Historia{}, err
	}
	if !EnEscala(storyPoints) {
		return Historia{}, ErrStoryPointsFueraDeEscala
	}
	if h.Estado == EstadoCompletada {
		return Historia{}, ErrHistoriaCompletada
	}
	h.StoryPoints = &storyPoints
	if err := s.historias.Actualizar(ctx, h); err != nil {
		return Historia{}, err
	}
	return h, nil
}

// QuitarEstimacion deja la historia sin estimar (spec 07 CA-2, RN-3). Si ya estaba sin estimar, no
// es un error. Una historia completada no se puede cambiar (RN-4).
func (s *ServicioBacklog) QuitarEstimacion(ctx context.Context, proyectoID, id int64) (Historia, error) {
	h, err := s.historiaDelProyecto(ctx, proyectoID, id)
	if err != nil {
		return Historia{}, err
	}
	if h.Estado == EstadoCompletada {
		return Historia{}, ErrHistoriaCompletada
	}
	h.StoryPoints = nil
	if err := s.historias.Actualizar(ctx, h); err != nil {
		return Historia{}, err
	}
	return h, nil
}

// historiaDelProyecto devuelve la historia si el proyecto existe y la historia es suya (spec 07 RN-6).
func (s *ServicioBacklog) historiaDelProyecto(ctx context.Context, proyectoID, id int64) (Historia, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return Historia{}, err
	}
	return s.historias.Obtener(ctx, proyectoID, id)
}
