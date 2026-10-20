// Package memoria implementa los repositorios del dominio en memoria, para los tests
// (unitarios, de la API y BDD) que no necesitan una base de datos real.
package memoria

import (
	"context"
	"sync"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// RepositorioProyectos guarda los proyectos en un slice, en orden de creación.
type RepositorioProyectos struct {
	mu        sync.Mutex
	proyectos []dominio.Proyecto
	ultimoID  int64
}

func NuevoRepositorioProyectos() *RepositorioProyectos {
	return &RepositorioProyectos{}
}

func (r *RepositorioProyectos) Guardar(_ context.Context, p dominio.Proyecto) (dominio.Proyecto, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ultimoID++
	p.ID = r.ultimoID
	r.proyectos = append(r.proyectos, p)
	return p, nil
}

func (r *RepositorioProyectos) Actualizar(_ context.Context, p dominio.Proyecto) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.proyectos {
		if r.proyectos[i].ID == p.ID {
			r.proyectos[i] = p
			return nil
		}
	}
	return dominio.ErrProyectoNoExiste
}

func (r *RepositorioProyectos) Listar(_ context.Context) ([]dominio.Proyecto, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]dominio.Proyecto{}, r.proyectos...), nil
}

func (r *RepositorioProyectos) Obtener(_ context.Context, id int64) (dominio.Proyecto, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.proyectos {
		if p.ID == id {
			return p, nil
		}
	}
	return dominio.Proyecto{}, dominio.ErrProyectoNoExiste
}
