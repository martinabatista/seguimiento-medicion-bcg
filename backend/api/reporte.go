package api

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/pdf"
)

// integranteDeReporteJSON es un integrante en el encabezado del reporte del proyecto.
type integranteDeReporteJSON struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

// historiaDeReporteJSON es una historia con sus estimaciones y su esfuerzo (spec 14 CA-1 a CA-3).
// completada_en_sprint solo aparece en el reporte de un Sprint.
type historiaDeReporteJSON struct {
	HistoriaID         int64           `json:"historia_id"`
	Titulo             string          `json:"titulo"`
	Estado             string          `json:"estado"`
	StoryPoints        *int            `json:"story_points"`    // null = sin estimar
	HorasEstimadas     *json.Number    `json:"horas_estimadas"` // null = sin horas estimadas
	HorasReales        json.Number     `json:"horas_reales"`
	Desviacion         *desviacionJSON `json:"desviacion"` // null = "sin datos"
	CompletadaEnSprint *bool           `json:"completada_en_sprint,omitempty"`
}

// sprintDeReporteJSON es un Sprint con sus métricas de avance (spec 10) y de esfuerzo (spec 12),
// con el mismo JSON que esas rutas.
type sprintDeReporteJSON struct {
	SprintID    int64              `json:"sprint_id"`
	Numero      int                `json:"numero"`
	FechaInicio string             `json:"fecha_inicio"`
	FechaFin    string             `json:"fecha_fin"`
	SprintGoal  string             `json:"sprint_goal"`
	Estado      string             `json:"estado"`
	Parcial     bool               `json:"parcial"`
	Avance      metricasSprintJSON `json:"avance"`
	Esfuerzo    esfuerzoSprintJSON `json:"esfuerzo"`
}

// metricasDeReporteJSON son las métricas del proyecto (spec 10 RN-6, RN-7 y spec 12 RN-9, RN-10).
type metricasDeReporteJSON struct {
	SprintsCerrados              int             `json:"sprints_cerrados"`
	Velocidad                    *float64        `json:"velocidad"` // null = todavía no hay
	Totales                      totalesJSON     `json:"totales"`
	HorasEstimadas               json.Number     `json:"horas_estimadas"`
	HorasReales                  json.Number     `json:"horas_reales"`
	Desviacion                   *desviacionJSON `json:"desviacion"`
	CompletadasSinHorasEstimadas int             `json:"completadas_sin_horas_estimadas"`
}

type defectosDeReporteJSON struct {
	Detectados int           `json:"detectados"`
	Resueltos  int           `json:"resueltos"`
	Lista      []defectoJSON `json:"lista"`
}

// reporteProyectoJSON es la vista previa del reporte del proyecto (spec 14 Salidas esperadas).
type reporteProyectoJSON struct {
	Proyecto    proyectoJSON              `json:"proyecto"`
	Estado      string                    `json:"estado"`
	GeneradoEl  string                    `json:"generado_el"`
	Integrantes []integranteDeReporteJSON `json:"integrantes"`
	Historias   []historiaDeReporteJSON   `json:"historias"`
	Sprints     []sprintDeReporteJSON     `json:"sprints"`
	Metricas    metricasDeReporteJSON     `json:"metricas"`
	Defectos    defectosDeReporteJSON     `json:"defectos"`
}

// reporteSprintJSON es la vista previa del reporte de un Sprint.
type reporteSprintJSON struct {
	Proyecto   proyectoJSON            `json:"proyecto"`
	GeneradoEl string                  `json:"generado_el"`
	Sprint     sprintDeReporteJSON     `json:"sprint"`
	Historias  []historiaDeReporteJSON `json:"historias"`
	Defectos   defectosDeReporteJSON   `json:"defectos"`
}

type reporteHandler struct {
	servicio *dominio.ServicioReporte
}

// proyecto: GET /api/proyectos/{id}/reporte (CA-1 a CA-5, CA-7).
func (h reporteHandler) proyecto(w http.ResponseWriter, r *http.Request) {
	rp, ok := h.reporteProyecto(w, r)
	if ok {
		writeJSON(w, http.StatusOK, reporteProyectoAJSON(rp))
	}
}

// proyectoPDF: GET /api/proyectos/{id}/reporte/pdf (CA-6).
func (h reporteHandler) proyectoPDF(w http.ResponseWriter, r *http.Request) {
	rp, ok := h.reporteProyecto(w, r)
	if !ok {
		return
	}
	archivo, err := pdf.ReporteProyecto(rp)
	if err != nil {
		responderError(w, err)
		return
	}
	responderPDF(w, "Reporte "+rp.Proyecto.Nombre+".pdf", archivo)
}

// sprint: GET /api/proyectos/{id}/sprints/{sprintID}/reporte (CA-1 a CA-5, CA-7).
func (h reporteHandler) sprint(w http.ResponseWriter, r *http.Request) {
	rs, ok := h.reporteSprint(w, r)
	if ok {
		writeJSON(w, http.StatusOK, reporteSprintAJSON(rs))
	}
}

// sprintPDF: GET /api/proyectos/{id}/sprints/{sprintID}/reporte/pdf (CA-6).
func (h reporteHandler) sprintPDF(w http.ResponseWriter, r *http.Request) {
	rs, ok := h.reporteSprint(w, r)
	if !ok {
		return
	}
	archivo, err := pdf.ReporteSprint(rs)
	if err != nil {
		responderError(w, err)
		return
	}
	responderPDF(w, fmt.Sprintf("Reporte %s - Sprint %d.pdf", rs.Proyecto.Nombre, rs.Sprint.Sprint.Numero), archivo)
}

// reporteProyecto pide el reporte; si falla, ya respondió el error y devuelve false.
func (h reporteHandler) reporteProyecto(w http.ResponseWriter, r *http.Request) (dominio.ReporteProyecto, bool) {
	proyectoID, err := leerID(r)
	if err != nil {
		responderError(w, err)
		return dominio.ReporteProyecto{}, false
	}
	rp, err := h.servicio.ReporteProyecto(r.Context(), proyectoID)
	if err != nil {
		responderError(w, err)
		return dominio.ReporteProyecto{}, false
	}
	return rp, true
}

func (h reporteHandler) reporteSprint(w http.ResponseWriter, r *http.Request) (dominio.ReporteSprint, bool) {
	proyectoID, sprintID, err := leerIDsSprint(r)
	if err != nil {
		responderError(w, err)
		return dominio.ReporteSprint{}, false
	}
	rs, err := h.servicio.ReporteSprint(r.Context(), proyectoID, sprintID)
	if err != nil {
		responderError(w, err)
		return dominio.ReporteSprint{}, false
	}
	return rs, true
}

// responderPDF manda el archivo para descargar. mime.FormatMediaType codifica el nombre si tiene
// tildes (filename*=utf-8”…), como pide HTTP, y el navegador lo muestra con sus tildes.
func responderPDF(w http.ResponseWriter, nombre string, archivo []byte) {
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": nombre}))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(archivo)
}

func reporteProyectoAJSON(r dominio.ReporteProyecto) reporteProyectoJSON {
	out := reporteProyectoJSON{
		Proyecto:    aJSON(r.Proyecto),
		Estado:      string(r.Estado),
		GeneradoEl:  r.GeneradoEl.Format(dominio.FormatoFecha),
		Integrantes: make([]integranteDeReporteJSON, 0, len(r.Integrantes)),
		Historias:   historiasDeReporteAJSON(r.Historias, false),
		Sprints:     make([]sprintDeReporteJSON, 0, len(r.Sprints)),
		Metricas: metricasDeReporteJSON{
			SprintsCerrados:              r.SprintsCerrados,
			Velocidad:                    velocidadAJSON(r.HayVelocidad, r.VelocidadDecimas),
			Totales:                      totalesJSON(r.Totales),
			HorasEstimadas:               horasConSigno(r.HorasEstimadas),
			HorasReales:                  horasConSigno(r.HorasReales),
			Desviacion:                   desviacionAJSON(r.Desviacion),
			CompletadasSinHorasEstimadas: r.CompletadasSinEstimar,
		},
		Defectos: defectosDeReporteAJSON(r.Defectos),
	}
	for _, i := range r.Integrantes {
		out.Integrantes = append(out.Integrantes, integranteDeReporteJSON{Nombre: i.Nombre, Rol: string(i.Rol)})
	}
	for _, s := range r.Sprints {
		out.Sprints = append(out.Sprints, sprintDeReporteAJSON(s))
	}
	return out
}

func reporteSprintAJSON(r dominio.ReporteSprint) reporteSprintJSON {
	return reporteSprintJSON{
		Proyecto:   aJSON(r.Proyecto),
		GeneradoEl: r.GeneradoEl.Format(dominio.FormatoFecha),
		Sprint:     sprintDeReporteAJSON(r.Sprint),
		Historias:  historiasDeReporteAJSON(r.Historias, true),
		Defectos:   defectosDeReporteAJSON(r.Defectos),
	}
}

// historiasDeReporteAJSON pasa las historias al JSON; deSprint agrega completada_en_sprint.
func historiasDeReporteAJSON(historias []dominio.HistoriaDeReporte, deSprint bool) []historiaDeReporteJSON {
	out := make([]historiaDeReporteJSON, 0, len(historias))
	for _, h := range historias {
		hj := historiaDeReporteJSON{
			HistoriaID:     h.HistoriaID,
			Titulo:         h.Titulo,
			Estado:         string(h.Estado),
			StoryPoints:    h.StoryPoints,
			HorasEstimadas: horasOpcionales(h.Estimadas),
			HorasReales:    horasConSigno(h.Reales),
			Desviacion:     desviacionAJSON(h.Desviacion),
		}
		if deSprint {
			completada := h.CompletadaEnSprint
			hj.CompletadaEnSprint = &completada
		}
		out = append(out, hj)
	}
	return out
}

func sprintDeReporteAJSON(s dominio.SprintDeReporte) sprintDeReporteJSON {
	return sprintDeReporteJSON{
		SprintID:    s.Sprint.ID,
		Numero:      s.Sprint.Numero,
		FechaInicio: s.Sprint.FechaInicio.Format(dominio.FormatoFecha),
		FechaFin:    s.Sprint.FechaFin.Format(dominio.FormatoFecha),
		SprintGoal:  s.Sprint.Goal,
		Estado:      string(s.Sprint.Estado),
		Parcial:     s.Avance.Parcial,
		Avance:      metricasSprintAJSON(s.Avance),
		Esfuerzo:    esfuerzoSprintAJSON(s.Esfuerzo),
	}
}

func defectosDeReporteAJSON(d dominio.DefectosDeReporte) defectosDeReporteJSON {
	out := defectosDeReporteJSON{Detectados: d.Detectados, Resueltos: d.Resueltos, Lista: make([]defectoJSON, 0, len(d.Lista))}
	for _, def := range d.Lista {
		out.Lista = append(out.Lista, defectoAJSON(def))
	}
	return out
}
