package memoria

import (
	"context"
	"slices"
	"sync"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// RepositorioIntegrantes guarda los integrantes de todos los proyectos en un slice, en orden de alta.
type RepositorioIntegrantes struct {
	mu          sync.Mutex
	integrantes []dominio.Integrante
	ultimoID    int64
}

func NuevoRepositorioIntegrantes() *RepositorioIntegrantes {
	return &RepositorioIntegrantes{}
}

// Guardar rechaza el email repetido en el proyecto, como la restricción UNIQUE de la base (spec 02 RN-6).
func (r *RepositorioIntegrantes) Guardar(_ context.Context, i dominio.Integrante) (dominio.Integrante, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	repetido := slices.ContainsFunc(r.integrantes, func(otro dominio.Integrante) bool {
		return otro.ProyectoID == i.ProyectoID && otro.Email == i.Email
	})
	if repetido {
		return dominio.Integrante{}, dominio.ErrIntegranteYaRegistrado
	}
	r.ultimoID++
	i.ID = r.ultimoID
	r.integrantes = append(r.integrantes, i)
	return i, nil
}

func (r *RepositorioIntegrantes) Listar(_ context.Context, proyectoID int64) ([]dominio.Integrante, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	lista := []dominio.Integrante{}
	for _, i := range r.integrantes {
		if i.ProyectoID == proyectoID {
			lista = append(lista, i)
		}
	}
	return lista, nil
}

func (r *RepositorioIntegrantes) Obtener(_ context.Context, proyectoID, id int64) (dominio.Integrante, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, i := range r.integrantes {
		if i.ID == id && i.ProyectoID == proyectoID {
			return i, nil
		}
	}
	return dominio.Integrante{}, dominio.ErrIntegranteNoExiste
}
