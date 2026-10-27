package api

import (
	"encoding/json"
	"net/http"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// horasEstimadasEntrada es lo que llega al cargar las horas estimadas. Se leen como json.Number para
// que el dominio distinga "faltan" de "no son válidas" (spec 12 RN-2).
type horasEstimadasEntrada struct {
	HorasEstimadas json.Number `json:"horas_estimadas"`
}

// desviacionJSON es la desviación de la spec 12 RN-6: horas con signo (-0.5) y porcentaje entero.
type desviacionJSON struct {
	Horas      json.Number `json:"horas"`
	Porcentaje int         `json:"porcentaje"`
}

type esfuerzoHistoriaJSON struct {
	HistoriaID     int64           `json:"historia_id"`
	Titulo         string          `json:"titulo"`
	Estado         string          `json:"estado"`
	HorasEstimadas *json.Number    `json:"horas_estimadas"` // null = sin horas estimadas
	HorasReales    json.Number     `json:"horas_reales"`
	Desviacion     *desviacionJSON `json:"desviacion"` // null = no entra en la desviación (RN-7)
}

type esfuerzoSprintJSON struct {
	SprintID                     int64           `json:"sprint_id"`
	Numero                       int             `json:"numero"`
	Estado                       string          `json:"estado"`
	Parcial                      bool            `json:"parcial"`
	HorasEstimadas               json.Number     `json:"horas_estimadas"`
	HorasReales                  json.Number     `json:"horas_reales"`
	Desviacion                   *desviacionJSON `json:"desviacion"` // null = "sin datos" (RN-10)
	CompletadasSinHorasEstimadas int             `json:"completadas_sin_horas_estimadas"`
}

// metricasEsfuerzoJSON son las métricas de esfuerzo de un proyecto (spec 12).
type metricasEsfuerzoJSON struct {
	Historias                    []esfuerzoHistoriaJSON `json:"historias"`
	Sprints                      []esfuerzoSprintJSON   `json:"sprints"`
	HorasEstimadas               json.Number            `json:"horas_estimadas"`
	HorasReales                  json.Number            `json:"horas_reales"`
	Desviacion                   *desviacionJSON        `json:"desviacion"`
	CompletadasSinHorasEstimadas int                    `json:"completadas_sin_horas_estimadas"`
}

// estimarHoras: PUT /api/proyectos/{id}/backlog/{historiaID}/horas-estimadas (spec 12 CA-1, CA-2).
func (h backlogHandler) estimarHoras(w http.ResponseWriter, r *http.Request) {
	proyectoID, historiaID, err := leerIDs(r)
	if err != nil {
		responderError(w, err)
		return
	}
	var entrada horasEstimadasEntrada
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, errCuerpoInvalido)
		return
	}
	hist, err := h.servicio.EstimarHoras(r.Context(), proyectoID, historiaID, entrada.HorasEstimadas.String())
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, historiaAJSON(hist))
}

// quitarHorasEstimadas: DELETE /api/proyectos/{id}/backlog/{historiaID}/horas-estimadas (spec 12 CA-1).
func (h backlogHandler) quitarHorasEstimadas(w http.ResponseWriter, r *http.Request) {
	proyectoID, historiaID, err := leerIDs(r)
	if err != nil {
		responderError(w, err)
		return
	}
	hist, err := h.servicio.QuitarHorasEstimadas(r.Context(), proyectoID, historiaID)
	if err != nil {
		responderError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, historiaAJSON(hist))
}

type metricasEsfuerzoHandler struct {
	servicio *dominio.ServicioMetricasEsfuerzo
}

// consultar: GET /api/proyectos/{id}/metricas/esfuerzo (spec 12 CA-3 a CA-5, CA-7).
func (h metricasEsfuerzoHandler) consultar(w http.ResponseWriter, r *http.Request) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return
	}
	m, err := h.servicio.Consultar(r.Context(), proyectoID)
	if err != nil {
		responderError(w, err)
		return
	}
	respuesta := metricasEsfuerzoJSON{
		Historias:                    make([]esfuerzoHistoriaJSON, 0, len(m.Historias)),
		Sprints:                      make([]esfuerzoSprintJSON, 0, len(m.Sprints)),
		HorasEstimadas:               horasConSigno(m.Estimadas),
		HorasReales:                  horasConSigno(m.Reales),
		Desviacion:                   desviacionAJSON(m.Desviacion),
		CompletadasSinHorasEstimadas: m.CompletadasSinEstimar,
	}
	for _, e := range m.Historias {
		respuesta.Historias = append(respuesta.Historias, esfuerzoHistoriaJSON{
			HistoriaID:     e.HistoriaID,
			Titulo:         e.Titulo,
			Estado:         string(e.Estado),
			HorasEstimadas: horasOpcionales(e.Estimadas),
			HorasReales:    horasConSigno(e.Reales),
			Desviacion:     desviacionAJSON(e.Desviacion),
		})
	}
	for _, s := range m.Sprints {
		respuesta.Sprints = append(respuesta.Sprints, esfuerzoSprintJSON{
			SprintID:                     s.SprintID,
			Numero:                       s.Numero,
			Estado:                       string(s.Estado),
			Parcial:                      s.Parcial,
			HorasEstimadas:               horasConSigno(s.Estimadas),
			HorasReales:                  horasConSigno(s.Reales),
			Desviacion:                   desviacionAJSON(s.Desviacion),
			CompletadasSinHorasEstimadas: s.CompletadasSinEstimar,
		})
	}
	writeJSON(w, http.StatusOK, respuesta)
}

// horasConSigno pasa centésimas a horas para el JSON (250 → 2.5; −50 → −0.5). El dominio suma con
// enteros; recién acá se convierte, como en la spec 09.
func horasConSigno(centesimas int) json.Number {
	if centesimas < 0 {
		return json.Number("-" + dominio.HorasATexto(-centesimas))
	}
	return json.Number(dominio.HorasATexto(centesimas))
}

// horasOpcionales es horasConSigno para horas que pueden faltar: nil sale como null.
func horasOpcionales(centesimas *int) *json.Number {
	if centesimas == nil {
		return nil
	}
	horas := horasConSigno(*centesimas)
	return &horas
}

func desviacionAJSON(d *dominio.Desviacion) *desviacionJSON {
	if d == nil {
		return nil
	}
	return &desviacionJSON{Horas: horasConSigno(d.Centesimas), Porcentaje: d.Porcentaje}
}
