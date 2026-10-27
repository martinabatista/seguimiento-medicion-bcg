package memoria

import (
	"context"
	"slices"
	"sync"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// RepositorioDefectos guarda los defectos de todos los proyectos en un slice, en orden de registro.
type RepositorioDefectos struct {
	mu       sync.Mutex
	defectos []dominio.Defecto
	ultimoID int64
}

func NuevoRepositorioDefectos() *RepositorioDefectos {
	return &RepositorioDefectos{}
}

func (r *RepositorioDefectos) Guardar(_ context.Context, d dominio.Defecto) (dominio.Defecto, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ultimoID++
	d.ID = r.ultimoID
	d = sinDatosDeConsulta(d)
	r.defectos = append(r.defectos, d)
	return d, nil
}

func (r *RepositorioDefectos) Actualizar(_ context.Context, d dominio.Defecto) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	i := slices.IndexFunc(r.defectos, func(e dominio.Defecto) bool { return e.ID == d.ID && e.ProyectoID == d.ProyectoID })
	if i < 0 {
		return dominio.ErrDefectoNoExiste
	}
	// Solo cambian el estado y el Sprint de resolución (spec 11: lo demás no se modifica).
	r.defectos[i].Estado = d.Estado
	r.defectos[i].SprintResolucionID = d.SprintResolucionID
	return nil
}

func (r *RepositorioDefectos) Listar(_ context.Context, proyectoID int64) ([]dominio.Defecto, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	lista := []dominio.Defecto{}
	for _, d := range r.defectos {
		if d.ProyectoID == proyectoID {
			lista = append(lista, d)
		}
	}
	return lista, nil
}

func (r *RepositorioDefectos) Obtener(_ context.Context, proyectoID, id int64) (dominio.Defecto, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, d := range r.defectos {
		if d.ID == id && d.ProyectoID == proyectoID {
			return d, nil
		}
	}
	return dominio.Defecto{}, dominio.ErrDefectoNoExiste
}

// sinDatosDeConsulta borra lo que se completa al consultar: no se guarda con el defecto.
func sinDatosDeConsulta(d dominio.Defecto) dominio.Defecto {
	d.HistoriaTitulo = ""
	d.SprintDeteccionNumero = 0
	d.SprintResolucionNumero = 0
	return d
}
