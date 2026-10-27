// Package pdf arma los reportes de la spec 14 como archivos PDF, con la librería go-pdf/fpdf (Go
// puro, sin programas externos). Solo presenta: recibe los reportes ya armados por el dominio y no
// calcula nada. El dominio no depende de este paquete.
package pdf

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// Medidas de la hoja (A4, en milímetros) y de las letras (en puntos).
const (
	margen         = 15.0
	altoLinea      = 5.5
	letraTitulo    = 18.0
	letraSeccion   = 13.0
	letraTexto     = 10.0
	letraPie       = 8.0
	separacionSecc = 4.0
)

// Textos para lo que falta (spec 14 RN-7, RN-8): los mismos que muestra la vista previa.
const (
	sinEstimar        = "sin estimar"
	sinHorasEstimadas = "sin horas estimadas"
	sinDatos          = "sin datos"
	sinVelocidad      = "todavía no hay velocidad (no hay Sprints cerrados)"
)

// documento envuelve el PDF con las ayudas para escribir el reporte. Todo el texto pasa por tr,
// que lo convierte a la codificación de las fuentes estándar (cp1252): así se ven las tildes, la ñ
// y los signos del español (RN-10).
type documento struct {
	pdf *fpdf.Fpdf
	tr  func(string) string
}

func nuevoDocumento(titulo string) *documento {
	p := fpdf.New("P", "mm", "A4", "")
	d := &documento{pdf: p, tr: p.UnicodeTranslatorFromDescriptor("")}
	p.SetMargins(margen, margen, margen)
	p.SetAutoPageBreak(true, margen) // si no entra en una hoja, sigue en otra (RN-10)
	// Sin compresión: el texto queda legible dentro del archivo y los tests pueden buscarlo. Los
	// reportes son chicos, así que el tamaño no cambia en la práctica.
	p.SetCompression(false)
	p.SetTitle(titulo, true)
	p.AliasNbPages("")
	p.SetFooterFunc(func() {
		p.SetY(-margen + 5)
		p.SetFont("Helvetica", "", letraPie)
		p.CellFormat(0, altoLinea, d.tr(fmt.Sprintf("%s · Página %d de {nb}", titulo, p.PageNo())), "", 0, "C", false, 0, "")
	})
	p.AddPage()
	return d
}

func (d *documento) titulo(texto string) {
	d.pdf.SetFont("Helvetica", "B", letraTitulo)
	d.pdf.MultiCell(0, altoLinea*1.6, d.tr(texto), "", "L", false)
	d.pdf.Ln(2)
}

func (d *documento) seccion(texto string) {
	d.pdf.Ln(separacionSecc)
	d.pdf.SetFont("Helvetica", "B", letraSeccion)
	d.pdf.MultiCell(0, altoLinea*1.3, d.tr(texto), "B", "L", false)
	d.pdf.Ln(1)
}

// linea escribe un renglón; si es largo, sigue en los renglones siguientes (RN-10).
func (d *documento) linea(texto string) {
	d.pdf.SetFont("Helvetica", "", letraTexto)
	d.pdf.MultiCell(0, altoLinea, d.tr(texto), "", "L", false)
}

// item escribe el renglón principal de un elemento de una lista (una historia, un Sprint, un defecto).
func (d *documento) item(texto string) {
	d.pdf.Ln(1)
	d.pdf.SetFont("Helvetica", "B", letraTexto)
	d.pdf.MultiCell(0, altoLinea, d.tr(texto), "", "L", false)
}

func (d *documento) bytes() ([]byte, error) {
	var b bytes.Buffer
	if err := d.pdf.Output(&b); err != nil {
		return nil, fmt.Errorf("generar el PDF: %w", err)
	}
	return b.Bytes(), nil
}

// ReporteProyecto arma el PDF del reporte de un proyecto (spec 14 CA-6), con las secciones en el
// orden de la vista previa (RN-8).
func ReporteProyecto(r dominio.ReporteProyecto) ([]byte, error) {
	d := nuevoDocumento("Reporte del proyecto " + r.Proyecto.Nombre)
	d.titulo("Reporte del proyecto " + r.Proyecto.Nombre)
	if r.Proyecto.Descripcion != "" {
		d.linea(r.Proyecto.Descripcion)
	}
	d.linea(fmt.Sprintf("Fechas: %s al %s · Estado: %s", fecha(r.Proyecto.FechaInicio), fecha(r.Proyecto.FechaFin), r.Estado))
	d.linea("Generado el " + fecha(r.GeneradoEl))

	d.seccion("Integrantes")
	if len(r.Integrantes) == 0 {
		d.linea("sin integrantes")
	}
	for _, i := range r.Integrantes {
		d.linea(fmt.Sprintf("%s · %s", i.Nombre, i.Rol))
	}

	d.seccion("Historias")
	if len(r.Historias) == 0 {
		d.linea("sin historias")
	}
	for _, h := range r.Historias {
		d.item(fmt.Sprintf("%s — %s", h.Titulo, h.Estado))
		escribirEstimacionYEsfuerzo(d, h)
	}

	d.seccion("Sprints")
	if len(r.Sprints) == 0 {
		d.linea("sin Sprints")
	}
	for _, s := range r.Sprints {
		d.item(encabezadoSprint(s))
		escribirMetricasSprint(d, s)
	}

	d.seccion("Métricas del proyecto")
	if r.HayVelocidad {
		d.linea("Velocidad del equipo: " + velocidad(r.VelocidadDecimas) + " Story Points por Sprint")
	} else {
		d.linea("Velocidad del equipo: " + sinVelocidad)
	}
	t := r.Totales
	d.linea(fmt.Sprintf("Sprints cerrados: %d · Historias: %d de %d completadas (%d %%) · Story Points: %d de %d completados",
		r.SprintsCerrados, t.HistoriasCompletadas, t.HistoriasPlanificadas, t.PorcentajeCompletadas, t.StoryPointsCompletados, t.StoryPointsPlanificados))
	d.linea(fmt.Sprintf("Horas estimadas: %s · Horas reales: %s", horas(r.HorasEstimadas), horas(r.HorasReales)))
	d.linea("Desviación: " + desviacion(r.Desviacion))
	d.linea(fmt.Sprintf("Historias completadas sin horas estimadas: %d", r.CompletadasSinEstimar))

	escribirDefectos(d, r.Defectos)
	return d.bytes()
}

// ReporteSprint arma el PDF del reporte de un Sprint (spec 14 CA-6).
func ReporteSprint(r dominio.ReporteSprint) ([]byte, error) {
	s := r.Sprint
	titulo := fmt.Sprintf("Reporte del Sprint %d", s.Sprint.Numero)
	d := nuevoDocumento(titulo + " · " + r.Proyecto.Nombre)
	d.titulo(titulo)
	d.linea("Proyecto: " + r.Proyecto.Nombre)
	d.linea(fmt.Sprintf("Fechas: %s al %s", fecha(s.Sprint.FechaInicio), fecha(s.Sprint.FechaFin)))
	d.linea("Sprint Goal: " + s.Sprint.Goal)
	d.linea("Estado: " + estadoSprint(s))
	d.linea("Generado el " + fecha(r.GeneradoEl))

	d.seccion("Historias planificadas")
	if len(r.Historias) == 0 {
		d.linea("sin historias")
	}
	for _, h := range r.Historias {
		completada := "no completada en el Sprint"
		if h.CompletadaEnSprint {
			completada = "completada en el Sprint"
		}
		d.item(fmt.Sprintf("%s — %s", h.Titulo, completada))
		escribirEstimacionYEsfuerzo(d, h)
	}

	d.seccion("Métricas del Sprint")
	escribirMetricasSprint(d, s)

	escribirDefectos(d, r.Defectos)
	return d.bytes()
}

// escribirEstimacionYEsfuerzo escribe las estimaciones y el esfuerzo de una historia (CA-2, CA-3),
// en dos renglones cortos.
func escribirEstimacionYEsfuerzo(d *documento, h dominio.HistoriaDeReporte) {
	d.linea(fmt.Sprintf("Story Points: %s · Horas estimadas: %s", storyPoints(h.StoryPoints), horasOpcionales(h.Estimadas)))
	d.linea(fmt.Sprintf("Horas reales: %s · Desviación: %s", horas(h.Reales), desviacion(h.Desviacion)))
}

// escribirMetricasSprint escribe las métricas de avance (spec 10) y de esfuerzo (spec 12) de un Sprint.
func escribirMetricasSprint(d *documento, s dominio.SprintDeReporte) {
	a, e := s.Avance, s.Esfuerzo
	d.linea(fmt.Sprintf("Historias: %d de %d completadas (%d %%) · Sin estimar: %d",
		a.HistoriasCompletadas, a.HistoriasPlanificadas, a.PorcentajeCompletadas, a.SinEstimar))
	d.linea(fmt.Sprintf("Story Points: %d de %d completados", a.StoryPointsCompletados, a.StoryPointsPlanificados))
	d.linea(fmt.Sprintf("Horas estimadas: %s · Horas reales: %s · Desviación: %s",
		horas(e.Estimadas), horas(e.Reales), desviacion(e.Desviacion)))
	d.linea(fmt.Sprintf("Historias completadas sin horas estimadas: %d", e.CompletadasSinEstimar))
}

// escribirDefectos escribe las cantidades y la lista de defectos (CA-5).
func escribirDefectos(d *documento, defectos dominio.DefectosDeReporte) {
	d.seccion("Defectos")
	d.linea(fmt.Sprintf("Detectados: %d · Resueltos: %d", defectos.Detectados, defectos.Resueltos))
	if len(defectos.Lista) == 0 {
		d.linea("sin defectos")
	}
	for _, def := range defectos.Lista {
		d.item(def.Descripcion)
		d.linea(fmt.Sprintf("Severidad: %s · Estado: %s · Historia: %s", def.Severidad, def.Estado, def.HistoriaTitulo))
		resolucion := "sin Sprint de resolución"
		if def.SprintResolucionNumero != 0 {
			resolucion = fmt.Sprintf("resuelto en el Sprint %d", def.SprintResolucionNumero)
		}
		d.linea(fmt.Sprintf("Detectado en el Sprint %d · %s", def.SprintDeteccionNumero, resolucion))
	}
}

func encabezadoSprint(s dominio.SprintDeReporte) string {
	return fmt.Sprintf("Sprint %d (%s al %s) — %s · Sprint Goal: %s",
		s.Sprint.Numero, fecha(s.Sprint.FechaInicio), fecha(s.Sprint.FechaFin), estadoSprint(s), s.Sprint.Goal)
}

// estadoSprint marca el Sprint abierto como parcial (spec 14 RN-6).
func estadoSprint(s dominio.SprintDeReporte) string {
	if s.Avance.Parcial {
		return string(s.Sprint.Estado) + ", parcial"
	}
	return string(s.Sprint.Estado)
}

func fecha(t time.Time) string {
	return t.Format(dominio.FormatoFecha)
}

func storyPoints(p *int) string {
	if p == nil {
		return sinEstimar
	}
	return strconv.Itoa(*p)
}

func horasOpcionales(centesimas *int) string {
	if centesimas == nil {
		return sinHorasEstimadas
	}
	return horas(*centesimas)
}

// desviacion muestra la desviación como la pantalla de Métricas: con "+" si es positiva ("+2 h
// (+25 %)", "-0,5 h (-6 %)") o "sin datos" (spec 12 RN-10).
func desviacion(d *dominio.Desviacion) string {
	if d == nil {
		return sinDatos
	}
	signo, centesimas := "", d.Centesimas
	if centesimas > 0 {
		signo = "+"
	} else if centesimas < 0 {
		signo, centesimas = "-", -centesimas
	}
	porcentaje := strconv.Itoa(d.Porcentaje)
	if d.Porcentaje > 0 {
		porcentaje = "+" + porcentaje
	}
	return fmt.Sprintf("%s%s h (%s %%)", signo, horas(centesimas), porcentaje)
}

// horas muestra centésimas como las pantallas: con coma y hasta 2 decimales (250 → "2,5").
func horas(centesimas int) string {
	return strings.Replace(dominio.HorasATexto(centesimas), ".", ",", 1)
}

// velocidad muestra las décimas como en las pantallas: 50 → "5", 65 → "6,5" (spec 10 RN-6).
func velocidad(decimas int) string {
	return strings.Replace(strconv.FormatFloat(float64(decimas)/10, 'f', -1, 64), ".", ",", 1)
}
