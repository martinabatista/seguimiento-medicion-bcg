package memoria

import (
	"context"
	"slices"
	"sync"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// RepositorioVotaciones guarda las votaciones de Planning Poker en un slice. Guarda y devuelve
// copias: así un cambio en una votación leída no modifica la guardada sin Actualizar.
type RepositorioVotaciones struct {
	mu         sync.Mutex
	votaciones []dominio.Votacion
	ultimoID   int64
}

func NuevoRepositorioVotaciones() *RepositorioVotaciones {
	return &RepositorioVotaciones{}
}

func (r *RepositorioVotaciones) Guardar(_ context.Context, v dominio.Votacion) (dominio.Votacion, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ultimoID++
	v.ID = r.ultimoID
	r.votaciones = append(r.votaciones, copiarVotacion(v))
	return copiarVotacion(v), nil
}

func (r *RepositorioVotaciones) Actualizar(_ context.Context, v dominio.Votacion) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	i := r.buscar(v.ProyectoID, v.ID)
	if i < 0 {
		return dominio.ErrVotacionNoExiste
	}
	r.votaciones[i] = copiarVotacion(v)
	return nil
}

func (r *RepositorioVotaciones) Listar(_ context.Context, proyectoID int64) ([]dominio.Votacion, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	lista := []dominio.Votacion{}
	for _, v := range r.votaciones {
		if v.ProyectoID == proyectoID {
			lista = append(lista, copiarVotacion(v))
		}
	}
	return lista, nil
}

func (r *RepositorioVotaciones) Obtener(_ context.Context, proyectoID, id int64) (dominio.Votacion, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	i := r.buscar(proyectoID, id)
	if i < 0 {
		return dominio.Votacion{}, dominio.ErrVotacionNoExiste
	}
	return copiarVotacion(r.votaciones[i]), nil
}

// buscar devuelve la posición de la votación en el slice, o -1 si no está en ese proyecto.
func (r *RepositorioVotaciones) buscar(proyectoID, id int64) int {
	return slices.IndexFunc(r.votaciones, func(v dominio.Votacion) bool {
		return v.ID == id && v.ProyectoID == proyectoID
	})
}

// copiarVotacion copia también las rondas, los votos y la acordada (son slices y punteros).
func copiarVotacion(v dominio.Votacion) dominio.Votacion {
	rondas := make([]dominio.Ronda, len(v.Rondas))
	for k, ronda := range v.Rondas {
		ronda.Votos = slices.Clone(ronda.Votos)
		rondas[k] = ronda
	}
	v.Rondas = rondas
	if v.Acordada != nil {
		acordada := *v.Acordada
		v.Acordada = &acordada
	}
	return v
}
