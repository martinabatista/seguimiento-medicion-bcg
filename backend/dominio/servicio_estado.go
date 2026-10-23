package dominio

import (
	"context"
	"time"
)

// EstadoDeProyecto es la respuesta de la consulta de la spec 03: los datos del proyecto, su estado
// y el resumen de su backlog. El Sprint en curso (CA-2) se agrega cuando exista la spec 05.
type EstadoDeProyecto struct {
	Proyecto Proyecto
	Estado   EstadoProyecto
	Backlog  ResumenBacklog
}

// ServicioEstado es el caso de uso de la spec 03. Lee proyectos e historias, y la fecha de hoy
// de un reloj que recibe: en producción es la fecha del sistema y en los tests, una fija.
type ServicioEstado struct {
	proyectos RepositorioProyectos
	historias RepositorioHistorias
	hoy       func() time.Time
}

func NuevoServicioEstado(proyectos RepositorioProyectos, historias RepositorioHistorias, hoy func() time.Time) *ServicioEstado {
	return &ServicioEstado{proyectos: proyectos, historias: historias, hoy: hoy}
}

// Consultar arma el estado del proyecto en el momento de la consulta (RN-7): sus datos, su estado
// según la fecha de hoy (CA-1, RN-2) y el resumen de su backlog (CA-3). Si el proyecto no existe,
// devuelve ErrProyectoNoExiste (CA-4, RN-1).
func (s *ServicioEstado) Consultar(ctx context.Context, proyectoID int64) (EstadoDeProyecto, error) {
	p, err := s.proyectos.Obtener(ctx, proyectoID)
	if err != nil {
		return EstadoDeProyecto{}, err
	}
	historias, err := s.historias.Listar(ctx, proyectoID)
	if err != nil {
		return EstadoDeProyecto{}, err
	}
	return EstadoDeProyecto{
		Proyecto: p,
		Estado:   CalcularEstado(p, s.hoy()),
		Backlog:  ResumirBacklog(historias),
	}, nil
}
