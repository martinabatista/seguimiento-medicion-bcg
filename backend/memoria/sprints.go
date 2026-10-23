package memoria

import (
	"context"
	"slices"
	"sync"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// RepositorioSprints guarda los Sprints de todos los proyectos en un slice, en orden de creación.
type RepositorioSprints struct {
	mu       sync.Mutex
	sprints  []dominio.Sprint
	ultimoID int64
}

func NuevoRepositorioSprints() *RepositorioSprints {
	return &RepositorioSprints{}
}

func (r *RepositorioSprints) Guardar(_ context.Context, s dominio.Sprint) (dominio.Sprint, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ultimoID++
	s.ID = r.ultimoID
	s.Historias = slices.Clone(s.Historias)
	r.sprints = append(r.sprints, s)
	return copia(s), nil
}

func (r *RepositorioSprints) Actualizar(_ context.Context, s dominio.Sprint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.sprints {
		if r.sprints[i].ID == s.ID && r.sprints[i].ProyectoID == s.ProyectoID {
			r.sprints[i] = copia(s)
			return nil
		}
	}
	return dominio.ErrSprintNoExiste
}

func (r *RepositorioSprints) Listar(_ context.Context, proyectoID int64) ([]dominio.Sprint, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	lista := []dominio.Sprint{}
	for _, s := range r.sprints {
		if s.ProyectoID == proyectoID {
			lista = append(lista, copia(s))
		}
	}
	return lista, nil
}

func (r *RepositorioSprints) Obtener(_ context.Context, proyectoID, id int64) (dominio.Sprint, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, s := range r.sprints {
		if s.ID == id && s.ProyectoID == proyectoID {
			return copia(s), nil
		}
	}
	return dominio.Sprint{}, dominio.ErrSprintNoExiste
}

// copia evita que quien recibe un Sprint modifique, sin querer, la lista de historias guardada.
func copia(s dominio.Sprint) dominio.Sprint {
	s.Historias = slices.Clone(s.Historias)
	if s.Historias == nil {
		s.Historias = []int64{}
	}
	return s
}
