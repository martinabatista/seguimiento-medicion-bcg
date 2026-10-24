package memoria

import (
	"context"
	"slices"
	"sync"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// RepositorioEsfuerzo guarda los registros de esfuerzo de todas las historias en un slice, en orden de carga.
type RepositorioEsfuerzo struct {
	mu        sync.Mutex
	registros []dominio.RegistroEsfuerzo
	ultimoID  int64
}

func NuevoRepositorioEsfuerzo() *RepositorioEsfuerzo {
	return &RepositorioEsfuerzo{}
}

func (r *RepositorioEsfuerzo) Guardar(_ context.Context, e dominio.RegistroEsfuerzo) (dominio.RegistroEsfuerzo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ultimoID++
	e.ID = r.ultimoID
	e.IntegranteNombre = "" // el nombre no se guarda con el registro
	r.registros = append(r.registros, e)
	return e, nil
}

func (r *RepositorioEsfuerzo) Listar(_ context.Context, proyectoID, historiaID int64) ([]dominio.RegistroEsfuerzo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	lista := []dominio.RegistroEsfuerzo{}
	for _, e := range r.registros {
		if e.ProyectoID == proyectoID && e.HistoriaID == historiaID {
			lista = append(lista, e)
		}
	}
	return lista, nil
}

func (r *RepositorioEsfuerzo) Eliminar(_ context.Context, proyectoID, historiaID, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	i := slices.IndexFunc(r.registros, func(e dominio.RegistroEsfuerzo) bool {
		return e.ID == id && e.ProyectoID == proyectoID && e.HistoriaID == historiaID
	})
	if i < 0 {
		return dominio.ErrRegistroEsfuerzoNoExiste
	}
	r.registros = slices.Delete(r.registros, i, i+1)
	return nil
}
