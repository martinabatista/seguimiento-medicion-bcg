package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// integranteJSON es un integrante tal como sale por la API (spec 02).
type integranteJSON struct {
	ID         int64  `json:"id"`
	ProyectoID int64  `json:"proyecto_id"`
	Nombre     string `json:"nombre"`
	Email      string `json:"email"`
	Rol        string `json:"rol"`
}

// integranteEntrada es lo que llega al agregar un integrante.
type integranteEntrada struct {
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	Rol    string `json:"rol"`
}

type integrantesHandler struct {
	servicio *dominio.ServicioIntegrantes
}

// agregar: POST /api/proyectos/{id}/integrantes (CA-1). 409 si el email ya está en el proyecto (CA-3).
func (h integrantesHandler) agregar(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	var entrada integranteEntrada
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, errCuerpoInvalido)
		return
	}
	i, err := h.servicio.Agregar(r.Context(), proyectoID, dominio.DatosIntegrante{
		Nombre: entrada.Nombre,
		Email:  entrada.Email,
		Rol:    entrada.Rol,
	})
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, integranteAJSON(i))
}

// listar: GET /api/proyectos/{id}/integrantes (CA-2), en orden de alta.
func (h integrantesHandler) listar(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	integrantes, err := h.servicio.Listar(r.Context(), proyectoID)
	if err != nil {
		responderError(w, err)
		return
	}
	lista := make([]integranteJSON, 0, len(integrantes))
	for _, i := range integrantes {
		lista = append(lista, integranteAJSON(i))
	}
	writeJSON(w, http.StatusOK, lista)
}

// ver: GET /api/proyectos/{id}/integrantes/{integranteID} (CA-2, CA-6).
func (h integrantesHandler) ver(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	// Un id que no es un número no puede existir (CA-6), igual que en el backlog.
	integranteID, err := strconv.ParseInt(r.PathValue("integranteID"), 10, 64)
	if err != nil {
		responderError(w, dominio.ErrIntegranteNoExiste)
		return
	}
	i, err := h.servicio.Obtener(r.Context(), proyectoID, integranteID)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, integranteAJSON(i))
}

func integranteAJSON(i dominio.Integrante) integranteJSON {
	return integranteJSON{
		ID:         i.ID,
		ProyectoID: i.ProyectoID,
		Nombre:     i.Nombre,
		Email:      i.Email,
		Rol:        string(i.Rol),
	}
}
