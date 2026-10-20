package memoria

import (
	"context"
	"slices"
	"sync"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// RepositorioHistorias guarda las historias de todos los backlogs en un slice, en orden de creación.
type RepositorioHistorias struct {
	mu        sync.Mutex
	historias []dominio.Historia
	ultimoID  int64
}

func NuevoRepositorioHistorias() *RepositorioHistorias {
	return &RepositorioHistorias{}
}

func (r *RepositorioHistorias) Guardar(_ context.Context, h dominio.Historia) (dominio.Historia, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ultimoID++
	h.ID = r.ultimoID
	r.historias = append(r.historias, h)
	return h, nil
}

func (r *RepositorioHistorias) Actualizar(_ context.Context, h dominio.Historia) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	i := r.buscar(h.ProyectoID, h.ID)
	if i < 0 {
		return dominio.ErrHistoriaNoExiste
	}
	r.historias[i] = h
	return nil
}

func (r *RepositorioHistorias) Eliminar(_ context.Context, proyectoID, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	i := r.buscar(proyectoID, id)
	if i < 0 {
		return dominio.ErrHistoriaNoExiste
	}
	r.historias = slices.Delete(r.historias, i, i+1)
	return nil
}

func (r *RepositorioHistorias) Listar(_ context.Context, proyectoID int64) ([]dominio.Historia, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	lista := []dominio.Historia{}
	for _, h := range r.historias {
		if h.ProyectoID == proyectoID {
			lista = append(lista, h)
		}
	}
	return lista, nil
}

func (r *RepositorioHistorias) Obtener(_ context.Context, proyectoID, id int64) (dominio.Historia, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	i := r.buscar(proyectoID, id)
	if i < 0 {
		return dominio.Historia{}, dominio.ErrHistoriaNoExiste
	}
	return r.historias[i], nil
}

// buscar devuelve la posición de la historia en el slice, o -1 si no está en ese proyecto.
func (r *RepositorioHistorias) buscar(proyectoID, id int64) int {
	return slices.IndexFunc(r.historias, func(h dominio.Historia) bool {
		return h.ID == id && h.ProyectoID == proyectoID
	})
}
