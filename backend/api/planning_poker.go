package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// votacionJSON es una votación de Planning Poker tal como sale por la API (spec 08). Mientras la
// ronda actual no se revela, solo dice quién votó y quién falta: los valores aparecen únicamente
// en "rondas", que trae las rondas ya reveladas (RN-5).
type votacionJSON struct {
	ID         int64                    `json:"id"`
	ProyectoID int64                    `json:"proyecto_id"`
	HistoriaID int64                    `json:"historia_id"`
	Estado     string                   `json:"estado"`
	Ronda      int                      `json:"ronda"`
	Acordada   *int                     `json:"estimacion_acordada"`
	Votaron    []integranteVotacionJSON `json:"votaron"`
	Faltan     []integranteVotacionJSON `json:"faltan"`
	Rondas     []rondaReveladaJSON      `json:"rondas"`
}

type integranteVotacionJSON struct {
	ID     int64  `json:"id"`
	Nombre string `json:"nombre"`
}

type rondaReveladaJSON struct {
	Numero         int                `json:"numero"`
	Votos          []votoReveladoJSON `json:"votos"`
	HayDiferencias bool               `json:"hay_diferencias"`
	MasBajo        extremoRondaJSON   `json:"mas_bajo"`
	MasAlto        extremoRondaJSON   `json:"mas_alto"`
}

type votoReveladoJSON struct {
	IntegranteID int64  `json:"integrante_id"`
	Nombre       string `json:"nombre"`
	StoryPoints  int    `json:"story_points"`
}

type extremoRondaJSON struct {
	StoryPoints int                      `json:"story_points"`
	Integrantes []integranteVotacionJSON `json:"integrantes"`
}

// votoEntrada es lo que llega al votar. Los Story Points se leen como json.Number, como en la
// spec 07, para distinguir "falta el valor" de "no es de la escala".
type votoEntrada struct {
	IntegranteID int64       `json:"integrante_id"`
	StoryPoints  json.Number `json:"story_points"`
}

type pokerHandler struct {
	servicio *dominio.ServicioPlanningPoker
}

// abrir: POST /api/proyectos/{id}/votaciones {"historia_id"} (spec 08 CA-1).
func (h pokerHandler) abrir(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	var entrada struct {
		HistoriaID int64 `json:"historia_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, errCuerpoInvalido)
		return
	}
	v, err := h.servicio.Abrir(r.Context(), proyectoID, entrada.HistoriaID)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, votacionAJSON(v))
}

// ver: GET /api/proyectos/{id}/votaciones/{votacionID} (spec 08 CA-3).
func (h pokerHandler) ver(w http.ResponseWriter, r *http.Request) {
	h.operar(w, r, h.servicio.Ver)
}

// votar: POST /api/proyectos/{id}/votaciones/{votacionID}/votos {"integrante_id", "story_points"} (CA-2).
func (h pokerHandler) votar(w http.ResponseWriter, r *http.Request) {
	proyectoID, votacionID, err := leerIDsVotacion(r)
	if err != nil {
		responderError(w, err)
		return
	}
	var entrada votoEntrada
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, errCuerpoInvalido)
		return
	}
	valor, err := dominio.ParsearStoryPoints(entrada.StoryPoints.String())
	if err != nil {
		responderError(w, err)
		return
	}
	v, err := h.servicio.Votar(r.Context(), proyectoID, votacionID, entrada.IntegranteID, valor)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, votacionAJSON(v))
}

// revelar: POST /api/proyectos/{id}/votaciones/{votacionID}/revelar (CA-4, CA-5).
func (h pokerHandler) revelar(w http.ResponseWriter, r *http.Request) {
	h.operar(w, r, h.servicio.Revelar)
}

// nuevaRonda: POST /api/proyectos/{id}/votaciones/{votacionID}/rondas (CA-6).
func (h pokerHandler) nuevaRonda(w http.ResponseWriter, r *http.Request) {
	h.operar(w, r, h.servicio.NuevaRonda)
}

// acordar: POST /api/proyectos/{id}/votaciones/{votacionID}/acordar {"story_points"} (CA-7).
func (h pokerHandler) acordar(w http.ResponseWriter, r *http.Request) {
	proyectoID, votacionID, err := leerIDsVotacion(r)
	if err != nil {
		responderError(w, err)
		return
	}
	var entrada estimacionEntrada
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, errCuerpoInvalido)
		return
	}
	valor, err := dominio.ParsearStoryPoints(entrada.StoryPoints.String())
	if err != nil {
		responderError(w, err)
		return
	}
	v, err := h.servicio.Acordar(r.Context(), proyectoID, votacionID, valor)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, votacionAJSON(v))
}

// operar resuelve las rutas que solo necesitan el proyecto y la votación.
func (h pokerHandler) operar(w http.ResponseWriter, r *http.Request, op func(ctx context.Context, proyectoID, votacionID int64) (dominio.VistaVotacion, error)) {
	proyectoID, votacionID, err := leerIDsVotacion(r)
	if err != nil {
		responderError(w, err)
		return
	}
	v, err := op(r.Context(), proyectoID, votacionID)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, votacionAJSON(v))
}

// leerIDsVotacion lee el {id} del proyecto y el {votacionID} de la ruta. Un id que no es un número
// no puede existir (CA-8).
func leerIDsVotacion(r *http.Request) (proyectoID, votacionID int64, err error) {
	proyectoID, err = leerID(r)
	if err != nil {
		return 0, 0, err
	}
	votacionID, err = strconv.ParseInt(r.PathValue("votacionID"), 10, 64)
	if err != nil {
		return 0, 0, dominio.ErrVotacionNoExiste
	}
	return proyectoID, votacionID, nil
}

func votacionAJSON(v dominio.VistaVotacion) votacionJSON {
	out := votacionJSON{
		ID:         v.ID,
		ProyectoID: v.ProyectoID,
		HistoriaID: v.HistoriaID,
		Estado:     string(v.Estado),
		Ronda:      v.Ronda,
		Acordada:   v.Acordada,
		Votaron:    integrantesVotacionAJSON(v.Votaron),
		Faltan:     integrantesVotacionAJSON(v.Faltan),
		Rondas:     make([]rondaReveladaJSON, 0, len(v.Reveladas)),
	}
	for _, r := range v.Reveladas {
		ronda := rondaReveladaJSON{
			Numero:         r.Numero,
			Votos:          make([]votoReveladoJSON, 0, len(r.Votos)),
			HayDiferencias: r.Resultado.HayDiferencias,
			MasBajo:        extremoAJSON(r.Resultado.Minimo),
			MasAlto:        extremoAJSON(r.Resultado.Maximo),
		}
		for _, voto := range r.Votos {
			ronda.Votos = append(ronda.Votos, votoReveladoJSON{IntegranteID: voto.Integrante.ID, Nombre: voto.Integrante.Nombre, StoryPoints: voto.Valor})
		}
		out.Rondas = append(out.Rondas, ronda)
	}
	return out
}

func extremoAJSON(e dominio.ExtremoRonda) extremoRondaJSON {
	return extremoRondaJSON{StoryPoints: e.Valor, Integrantes: integrantesVotacionAJSON(e.Integrantes)}
}

// integrantesVotacionAJSON devuelve [] (no null) cuando no hay nadie.
func integrantesVotacionAJSON(is []dominio.Integrante) []integranteVotacionJSON {
	out := make([]integranteVotacionJSON, 0, len(is))
	for _, i := range is {
		out = append(out, integranteVotacionJSON{ID: i.ID, Nombre: i.Nombre})
	}
	return out
}
