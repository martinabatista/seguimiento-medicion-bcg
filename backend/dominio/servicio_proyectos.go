package dominio

import (
	"context"
	"time"
)

// RepositorioProyectos guarda y lee proyectos. Lo implementan postgres/ (base real) y
// memoria/ (tests). El dominio no sabe cuál se usa.
type RepositorioProyectos interface {
	// Guardar registra un proyecto nuevo y lo devuelve con el ID asignado.
	Guardar(ctx context.Context, p Proyecto) (Proyecto, error)
	// Actualizar reemplaza los datos de un proyecto existente.
	Actualizar(ctx context.Context, p Proyecto) error
	// Listar devuelve todos los proyectos en orden de creación.
	Listar(ctx context.Context) ([]Proyecto, error)
	// Obtener devuelve un proyecto o ErrProyectoNoExiste.
	Obtener(ctx context.Context, id int64) (Proyecto, error)
}

// ServicioProyectos son los casos de uso de la spec 01.
type ServicioProyectos struct {
	repo RepositorioProyectos
}

func NuevoServicioProyectos(repo RepositorioProyectos) *ServicioProyectos {
	return &ServicioProyectos{repo: repo}
}

// Crear valida los datos y registra el proyecto (CA-1). Si algo es inválido, no registra nada.
func (s *ServicioProyectos) Crear(ctx context.Context, nombre, descripcion string, inicio, fin time.Time) (Proyecto, error) {
	p, err := NuevoProyecto(nombre, descripcion, inicio, fin)
	if err != nil {
		return Proyecto{}, err
	}
	return s.repo.Guardar(ctx, p)
}

// Modificar reemplaza los datos de un proyecto existente (CA-2). Primero verifica que exista
// (CA-7) y después valida con las mismas reglas que al crear: si falla, no cambia nada (RN-5).
func (s *ServicioProyectos) Modificar(ctx context.Context, id int64, nombre, descripcion string, inicio, fin time.Time) (Proyecto, error) {
	if _, err := s.repo.Obtener(ctx, id); err != nil {
		return Proyecto{}, err
	}
	p, err := NuevoProyecto(nombre, descripcion, inicio, fin)
	if err != nil {
		return Proyecto{}, err
	}
	p.ID = id // RN-6: el identificador no cambia.
	if err := s.repo.Actualizar(ctx, p); err != nil {
		return Proyecto{}, err
	}
	return p, nil
}

// Listar devuelve todos los proyectos en orden de creación (CA-6).
func (s *ServicioProyectos) Listar(ctx context.Context) ([]Proyecto, error) {
	return s.repo.Listar(ctx)
}

// Obtener devuelve un proyecto por su identificador (CA-6) o ErrProyectoNoExiste (CA-7).
func (s *ServicioProyectos) Obtener(ctx context.Context, id int64) (Proyecto, error) {
	return s.repo.Obtener(ctx, id)
}
