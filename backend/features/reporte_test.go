package features

// Pasos de godog para los escenarios del reporte (spec 14, reporte.feature).

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strings"

	"github.com/cucumber/godog"
	"github.com/go-pdf/fpdf"
)

// estadoReporte guarda el último reporte visto o descargado en un escenario.
type estadoReporte struct {
	visto     reporteVisto
	pdf       []byte
	nombrePDF string
}

// reporteVisto lee tanto el reporte del proyecto como el de un Sprint: el del proyecto trae
// "sprints" y "metricas"; el de un Sprint, "sprint".
type reporteVisto struct {
	Proyecto    proyecto `json:"proyecto"`
	Estado      string   `json:"estado"`
	GeneradoEl  string   `json:"generado_el"`
	Integrantes []struct {
		Nombre string `json:"nombre"`
		Rol    string `json:"rol"`
	} `json:"integrantes"`
	Historias []historiaDeReporte `json:"historias"`
	Sprints   []sprintDeReporte   `json:"sprints"`
	Sprint    *sprintDeReporte    `json:"sprint"`
	Metricas  struct {
		SprintsCerrados int      `json:"sprints_cerrados"`
		Velocidad       *float64 `json:"velocidad"`
		Totales         struct {
			StoryPointsPlanificados int `json:"story_points_planificados"`
			StoryPointsCompletados  int `json:"story_points_completados"`
		} `json:"totales"`
		Desviacion *desviacionDeReporte `json:"desviacion"`
	} `json:"metricas"`
	Defectos struct {
		Detectados int       `json:"detectados"`
		Resueltos  int       `json:"resueltos"`
		Lista      []defecto `json:"lista"`
	} `json:"defectos"`
}

type desviacionDeReporte struct {
	Horas      json.Number `json:"horas"`
	Porcentaje int         `json:"porcentaje"`
}

type historiaDeReporte struct {
	Titulo             string               `json:"titulo"`
	Estado             string               `json:"estado"`
	StoryPoints        *int                 `json:"story_points"`
	HorasEstimadas     *json.Number         `json:"horas_estimadas"`
	HorasReales        json.Number          `json:"horas_reales"`
	Desviacion         *desviacionDeReporte `json:"desviacion"`
	CompletadaEnSprint *bool                `json:"completada_en_sprint"`
}

type sprintDeReporte struct {
	Numero  int  `json:"numero"`
	Parcial bool `json:"parcial"`
	Avance  struct {
		HistoriasPlanificadas   int `json:"historias_planificadas"`
		HistoriasCompletadas    int `json:"historias_completadas"`
		PorcentajeCompletadas   int `json:"porcentaje_completadas"`
		StoryPointsPlanificados int `json:"story_points_planificados"`
		StoryPointsCompletados  int `json:"story_points_completados"`
	} `json:"avance"`
	Esfuerzo struct {
		HorasEstimadas json.Number          `json:"horas_estimadas"`
		HorasReales    json.Number          `json:"horas_reales"`
		Desviacion     *desviacionDeReporte `json:"desviacion"`
	} `json:"esfuerzo"`
}

func (m *mundo) registrarPasosReporte(sc *godog.ScenarioContext) {
	// Acciones
	sc.Step(`^genero el reporte del proyecto$`, m.generoElReporteDelProyecto)
	sc.Step(`^genero el reporte del proyecto "([^"]*)"$`, m.generoElReporteDelProyectoLlamado)
	sc.Step(`^genero el reporte del proyecto con identificador (\d+)$`, m.generoElReporteDelProyectoConID)
	sc.Step(`^genero el reporte del Sprint (\d+)$`, m.generoElReporteDelSprint)
	sc.Step(`^genero el reporte del Sprint con identificador (\d+)$`, m.generoElReporteDelSprintConID)
	sc.Step(`^genero el reporte del Sprint de "([^"]*)" desde el proyecto "([^"]*)"$`, m.generoElReporteDelSprintDeOtroProyecto)
	sc.Step(`^descargo el reporte del proyecto en PDF$`, m.descargoElReporteDelProyecto)
	sc.Step(`^descargo el reporte del proyecto "([^"]*)" en PDF$`, m.descargoElReporteDelProyectoLlamado)
	sc.Step(`^descargo el reporte del proyecto con identificador (\d+) en PDF$`, m.descargoElReporteDelProyectoConID)
	sc.Step(`^descargo el reporte del Sprint (\d+) en PDF$`, m.descargoElReporteDelSprint)

	// Verificaciones: encabezado e historias (CA-1, CA-2, CA-3)
	sc.Step(`^el reporte es del proyecto "([^"]*)", "([^"]*)", generado el "([^"]*)"$`, m.elReporteEsDelProyecto)
	sc.Step(`^el reporte tiene al integrante "([^"]*)" con rol "([^"]*)"$`, m.elReporteTieneAlIntegrante)
	sc.Step(`^el reporte muestra las historias "([^"]*)", "([^"]*)" y "([^"]*)" en ese orden$`, m.elReporteMuestraTresHistorias)
	sc.Step(`^el reporte del Sprint muestra las historias "([^"]*)" y "([^"]*)"$`, m.elReporteMuestraDosHistorias)
	sc.Step(`^la historia "([^"]*)" del reporte está "([^"]*)"$`, m.laHistoriaDelReporteEsta)
	sc.Step(`^en el reporte del Sprint la historia "([^"]*)" está completada y "([^"]*)" no$`, m.enElReporteDelSprintEstaCompletada)
	sc.Step(`^el reporte no incluye la historia "([^"]*)"$`, m.elReporteNoIncluyeLaHistoria)
	sc.Step(`^la historia "([^"]*)" del reporte tiene (\d+) Story Points y (\S+) horas estimadas$`, m.laHistoriaDelReporteTieneEstimaciones)
	sc.Step(`^la historia "([^"]*)" del reporte está sin estimar y sin horas estimadas$`, m.laHistoriaDelReporteEstaSinEstimar)
	sc.Step(`^la historia "([^"]*)" del reporte tiene (\d+) Story Points$`, m.laHistoriaDelReporteTieneStoryPoints)
	sc.Step(`^la historia "([^"]*)" del reporte tiene (\S+) horas reales y una desviación de (\S+) horas y (\S+) %$`, m.laHistoriaDelReporteTieneDesviacion)
	sc.Step(`^la historia "([^"]*)" del reporte tiene (\S+) horas reales y desviación "sin datos"$`, m.laHistoriaDelReporteTieneRealesSinDesviacion)

	// Verificaciones: métricas (CA-1, CA-3, CA-4)
	sc.Step(`^el Sprint (\d+) del reporte tiene (\d+) historias planificadas, (\d+) completadas? y (\d+) % completadas$`, m.elSprintDelReporteTieneHistorias)
	sc.Step(`^el Sprint (\d+) del reporte tiene (\d+) Story Points planificados y (\d+) completados$`, m.elSprintDelReporteTieneStoryPoints)
	sc.Step(`^el Sprint (\d+) del reporte es parcial$`, m.elSprintDelReporteEsParcial)
	sc.Step(`^el Sprint del reporte tiene (\S+) horas estimadas, (\S+) horas reales y una desviación de (\S+) horas y (\S+) %$`, m.elSprintDelReporteTieneEsfuerzo)
	sc.Step(`^el reporte del Sprint dice que es parcial$`, m.elReporteDelSprintEsParcial)
	sc.Step(`^la velocidad del reporte es ([\d.]+)$`, m.laVelocidadDelReporteEs)
	sc.Step(`^el reporte dice que todavía no hay velocidad$`, m.elReporteNoTieneVelocidad)
	sc.Step(`^el proyecto del reporte tiene (\d+) Sprints? cerrados? con (\d+) Story Points planificados y (\d+) completados$`, m.elProyectoDelReporteTieneStoryPoints)
	sc.Step(`^la desviación del proyecto en el reporte es "sin datos"$`, m.laDesviacionDelProyectoEnElReporteEsSinDatos)

	// Verificaciones: defectos (CA-5)
	sc.Step(`^el reporte tiene (\d+) defectos? detectados? y (\d+) resueltos?$`, m.elReporteTieneDefectos)
	sc.Step(`^el reporte muestra los defectos "([^"]*)" y "([^"]*)" en ese orden$`, m.elReporteMuestraLosDefectos)
	sc.Step(`^el defecto "([^"]*)" del reporte está "([^"]*)" en el Sprint (\d+)$`, m.elDefectoDelReporteEstaResueltoEn)
	sc.Step(`^el defecto "([^"]*)" del reporte está "([^"]*)" sin Sprint de resolución$`, m.elDefectoDelReporteEstaSinResolucion)

	// Verificaciones: PDF (CA-6) y sin datos (CA-7)
	sc.Step(`^recibo un archivo PDF llamado "([^"]*)"$`, m.reciboUnArchivoPDF)
	sc.Step(`^el PDF incluye el texto "([^"]*)"$`, m.elPDFIncluyeElTexto)
	sc.Step(`^el PDF incluye las secciones "([^"]*)", "([^"]*)", "([^"]*)" y "([^"]*)"$`, m.elPDFIncluyeLasSecciones)
	sc.Step(`^el reporte no tiene historias, Sprints ni defectos$`, m.elReporteNoTieneDatos)
	sc.Step(`^el reporte del Sprint no tiene historias y tiene (\d+) % completadas$`, m.elReporteDelSprintNoTieneHistorias)
}

// ---------- helpers ----------

func rutaReporteProyecto(proyectoID int64) string {
	return fmt.Sprintf("/api/proyectos/%d/reporte", proyectoID)
}

func rutaReporteSprint(proyectoID, sprintID int64) string {
	return fmt.Sprintf("/api/proyectos/%d/sprints/%d/reporte", proyectoID, sprintID)
}

// generar pide la vista previa del reporte y, si salió bien, la recuerda.
func (m *mundo) generar(ruta string) error {
	if err := m.pedir(http.MethodGet, ruta, nil); err != nil {
		return err
	}
	m.reporte.visto = reporteVisto{}
	if m.respuesta.Code != http.StatusOK {
		return nil // lo verifica "se informa el error"
	}
	return json.Unmarshal(m.respuesta.Body.Bytes(), &m.reporte.visto)
}

// descargar pide el PDF y, si salió bien, recuerda el archivo y su nombre.
func (m *mundo) descargar(ruta string) error {
	if err := m.pedir(http.MethodGet, ruta+"/pdf", nil); err != nil {
		return err
	}
	m.reporte.pdf, m.reporte.nombrePDF = nil, ""
	if m.respuesta.Code != http.StatusOK {
		return nil
	}
	if tipo := m.respuesta.Header().Get("Content-Type"); tipo != "application/pdf" {
		return fmt.Errorf("el archivo es %q, se esperaba application/pdf", tipo)
	}
	_, params, err := mime.ParseMediaType(m.respuesta.Header().Get("Content-Disposition"))
	if err != nil {
		return fmt.Errorf("Content-Disposition inválido: %w", err)
	}
	m.reporte.pdf, m.reporte.nombrePDF = m.respuesta.Body.Bytes(), params["filename"]
	return nil
}

func (m *mundo) historiaDelReporte(titulo string) (historiaDeReporte, error) {
	for _, h := range m.reporte.visto.Historias {
		if h.Titulo == titulo {
			return h, nil
		}
	}
	return historiaDeReporte{}, fmt.Errorf("el reporte no tiene la historia %q: %+v", titulo, m.reporte.visto.Historias)
}

func (m *mundo) sprintDelReporte(numero int) (sprintDeReporte, error) {
	for _, s := range m.reporte.visto.Sprints {
		if s.Numero == numero {
			return s, nil
		}
	}
	return sprintDeReporte{}, fmt.Errorf("el reporte no tiene el Sprint %d: %+v", numero, m.reporte.visto.Sprints)
}

func (m *mundo) sprintDelReporteDeSprint() (sprintDeReporte, error) {
	if m.reporte.visto.Sprint == nil {
		return sprintDeReporte{}, fmt.Errorf("no se generó el reporte de un Sprint: %s", m.respuesta.Body.String())
	}
	return *m.reporte.visto.Sprint, nil
}

func (m *mundo) defectoDelReporte(descripcion string) (defecto, error) {
	for _, d := range m.reporte.visto.Defectos.Lista {
		if d.Descripcion == descripcion {
			return d, nil
		}
	}
	return defecto{}, fmt.Errorf("el reporte no tiene el defecto %q", descripcion)
}

// desviacionEs compara una desviación con la del escenario ("2" horas y "25" %).
func desviacionEs(d *desviacionDeReporte, horas, porcentaje string) error {
	if d == nil {
		return fmt.Errorf("la desviación es \"sin datos\", se esperaba %s horas y %s %%", horas, porcentaje)
	}
	if d.Horas.String() != horas || fmt.Sprint(d.Porcentaje) != porcentaje {
		return fmt.Errorf("desviación = %s horas y %d %%, se esperaba %s horas y %s %%", d.Horas, d.Porcentaje, horas, porcentaje)
	}
	return nil
}

// ---------- acciones ----------

func (m *mundo) generoElReporteDelProyecto() error {
	return m.generar(rutaReporteProyecto(m.backlog.principal))
}

func (m *mundo) generoElReporteDelProyectoLlamado(nombre string) error {
	id, err := m.proyectoLlamado(nombre)
	if err != nil {
		return err
	}
	return m.generar(rutaReporteProyecto(id))
}

func (m *mundo) generoElReporteDelProyectoConID(id int64) error {
	return m.generar(rutaReporteProyecto(id))
}

func (m *mundo) generoElReporteDelSprint(numero int) error {
	id, err := m.sprintNumero(numero)
	if err != nil {
		return err
	}
	return m.generar(rutaReporteSprint(m.backlog.principal, id))
}

func (m *mundo) generoElReporteDelSprintConID(id int64) error {
	return m.generar(rutaReporteSprint(m.backlog.principal, id))
}

func (m *mundo) generoElReporteDelSprintDeOtroProyecto(duenio, desde string) error {
	sprintID, ok := m.defectos.sprintsAjenos[duenio]
	if !ok {
		return fmt.Errorf("el escenario no creó un Sprint en %q", duenio)
	}
	proyectoID, err := m.proyectoLlamado(desde)
	if err != nil {
		return err
	}
	return m.generar(rutaReporteSprint(proyectoID, sprintID))
}

func (m *mundo) descargoElReporteDelProyecto() error {
	return m.descargar(rutaReporteProyecto(m.backlog.principal))
}

func (m *mundo) descargoElReporteDelProyectoLlamado(nombre string) error {
	id, err := m.proyectoLlamado(nombre)
	if err != nil {
		return err
	}
	return m.descargar(rutaReporteProyecto(id))
}

func (m *mundo) descargoElReporteDelProyectoConID(id int64) error {
	return m.descargar(rutaReporteProyecto(id))
}

func (m *mundo) descargoElReporteDelSprint(numero int) error {
	id, err := m.sprintNumero(numero)
	if err != nil {
		return err
	}
	return m.descargar(rutaReporteSprint(m.backlog.principal, id))
}

// ---------- verificaciones: encabezado e historias ----------

func (m *mundo) elReporteEsDelProyecto(nombre, estado, generado string) error {
	r := m.reporte.visto
	if r.Proyecto.Nombre != nombre || r.Estado != estado || r.GeneradoEl != generado {
		return fmt.Errorf("reporte de %q, %q, generado el %q; se esperaba %q, %q, %q", r.Proyecto.Nombre, r.Estado, r.GeneradoEl, nombre, estado, generado)
	}
	return nil
}

func (m *mundo) elReporteTieneAlIntegrante(nombre, rol string) error {
	for _, i := range m.reporte.visto.Integrantes {
		if i.Nombre == nombre && i.Rol == rol {
			return nil
		}
	}
	return fmt.Errorf("el reporte no tiene a %q con rol %q: %+v", nombre, rol, m.reporte.visto.Integrantes)
}

func (m *mundo) titulosDelReporte() []string {
	titulos := []string{}
	for _, h := range m.reporte.visto.Historias {
		titulos = append(titulos, h.Titulo)
	}
	return titulos
}

func (m *mundo) elReporteMuestraTresHistorias(a, b, c string) error {
	if got := strings.Join(m.titulosDelReporte(), " | "); got != strings.Join([]string{a, b, c}, " | ") {
		return fmt.Errorf("historias = %s; se esperaban %q, %q y %q en ese orden", got, a, b, c)
	}
	return nil
}

func (m *mundo) elReporteMuestraDosHistorias(a, b string) error {
	if got := strings.Join(m.titulosDelReporte(), " | "); got != a+" | "+b {
		return fmt.Errorf("historias = %s; se esperaban %q y %q", got, a, b)
	}
	return nil
}

func (m *mundo) laHistoriaDelReporteEsta(titulo, estado string) error {
	h, err := m.historiaDelReporte(titulo)
	if err != nil {
		return err
	}
	if h.Estado != estado {
		return fmt.Errorf("estado = %q, se esperaba %q", h.Estado, estado)
	}
	return nil
}

func (m *mundo) enElReporteDelSprintEstaCompletada(completada, otra string) error {
	for titulo, esperado := range map[string]bool{completada: true, otra: false} {
		h, err := m.historiaDelReporte(titulo)
		if err != nil {
			return err
		}
		if h.CompletadaEnSprint == nil || *h.CompletadaEnSprint != esperado {
			return fmt.Errorf("historia %q: completada en el Sprint = %v, se esperaba %v", titulo, h.CompletadaEnSprint, esperado)
		}
	}
	return nil
}

func (m *mundo) elReporteNoIncluyeLaHistoria(titulo string) error {
	if _, err := m.historiaDelReporte(titulo); err == nil {
		return fmt.Errorf("el reporte incluye la historia %q", titulo)
	}
	return nil
}

func (m *mundo) laHistoriaDelReporteTieneEstimaciones(titulo string, storyPoints int, horas string) error {
	h, err := m.historiaDelReporte(titulo)
	if err != nil {
		return err
	}
	if h.StoryPoints == nil || *h.StoryPoints != storyPoints || h.HorasEstimadas == nil || h.HorasEstimadas.String() != horas {
		return fmt.Errorf("historia = %+v; se esperaban %d Story Points y %s horas estimadas", h, storyPoints, horas)
	}
	return nil
}

func (m *mundo) laHistoriaDelReporteEstaSinEstimar(titulo string) error {
	h, err := m.historiaDelReporte(titulo)
	if err != nil {
		return err
	}
	if h.StoryPoints != nil || h.HorasEstimadas != nil {
		return fmt.Errorf("historia = %+v; se esperaba sin estimar y sin horas estimadas", h)
	}
	return nil
}

func (m *mundo) laHistoriaDelReporteTieneStoryPoints(titulo string, storyPoints int) error {
	h, err := m.historiaDelReporte(titulo)
	if err != nil {
		return err
	}
	if h.StoryPoints == nil || *h.StoryPoints != storyPoints {
		return fmt.Errorf("Story Points = %v, se esperaban %d", h.StoryPoints, storyPoints)
	}
	return nil
}

func (m *mundo) laHistoriaDelReporteTieneDesviacion(titulo, reales, horas, porcentaje string) error {
	h, err := m.historiaDelReporte(titulo)
	if err != nil {
		return err
	}
	if h.HorasReales.String() != reales {
		return fmt.Errorf("horas reales = %s, se esperaban %s", h.HorasReales, reales)
	}
	return desviacionEs(h.Desviacion, horas, porcentaje)
}

func (m *mundo) laHistoriaDelReporteTieneRealesSinDesviacion(titulo, reales string) error {
	h, err := m.historiaDelReporte(titulo)
	if err != nil {
		return err
	}
	if h.HorasReales.String() != reales || h.Desviacion != nil {
		return fmt.Errorf("historia = %+v; se esperaban %s horas reales y desviación \"sin datos\"", h, reales)
	}
	return nil
}

// ---------- verificaciones: métricas ----------

func (m *mundo) elSprintDelReporteTieneHistorias(numero, planificadas, completadas, porcentaje int) error {
	s, err := m.sprintDelReporte(numero)
	if err != nil {
		return err
	}
	a := s.Avance
	if a.HistoriasPlanificadas != planificadas || a.HistoriasCompletadas != completadas || a.PorcentajeCompletadas != porcentaje {
		return fmt.Errorf("Sprint %d: %d planificadas, %d completadas, %d %%; se esperaban %d, %d, %d %%", numero,
			a.HistoriasPlanificadas, a.HistoriasCompletadas, a.PorcentajeCompletadas, planificadas, completadas, porcentaje)
	}
	return nil
}

func (m *mundo) elSprintDelReporteTieneStoryPoints(numero, planificados, completados int) error {
	s, err := m.sprintDelReporte(numero)
	if err != nil {
		return err
	}
	if s.Avance.StoryPointsPlanificados != planificados || s.Avance.StoryPointsCompletados != completados {
		return fmt.Errorf("Sprint %d: %d/%d Story Points, se esperaban %d/%d", numero, s.Avance.StoryPointsPlanificados, s.Avance.StoryPointsCompletados, planificados, completados)
	}
	return nil
}

func (m *mundo) elSprintDelReporteEsParcial(numero int) error {
	s, err := m.sprintDelReporte(numero)
	if err != nil {
		return err
	}
	if !s.Parcial {
		return fmt.Errorf("el Sprint %d del reporte no es parcial", numero)
	}
	return nil
}

func (m *mundo) elSprintDelReporteTieneEsfuerzo(estimadas, reales, horas, porcentaje string) error {
	s, err := m.sprintDelReporteDeSprint()
	if err != nil {
		return err
	}
	if s.Esfuerzo.HorasEstimadas.String() != estimadas || s.Esfuerzo.HorasReales.String() != reales {
		return fmt.Errorf("Sprint: %s horas estimadas y %s reales, se esperaban %s y %s", s.Esfuerzo.HorasEstimadas, s.Esfuerzo.HorasReales, estimadas, reales)
	}
	return desviacionEs(s.Esfuerzo.Desviacion, horas, porcentaje)
}

func (m *mundo) elReporteDelSprintEsParcial() error {
	s, err := m.sprintDelReporteDeSprint()
	if err != nil {
		return err
	}
	if !s.Parcial {
		return fmt.Errorf("el reporte del Sprint no es parcial")
	}
	return nil
}

func (m *mundo) laVelocidadDelReporteEs(esperada float64) error {
	v := m.reporte.visto.Metricas.Velocidad
	if v == nil || *v != esperada {
		return fmt.Errorf("velocidad = %v, se esperaba %v", v, esperada)
	}
	return nil
}

func (m *mundo) elReporteNoTieneVelocidad() error {
	if v := m.reporte.visto.Metricas.Velocidad; v != nil {
		return fmt.Errorf("velocidad = %v, se esperaba que todavía no haya", *v)
	}
	return nil
}

func (m *mundo) elProyectoDelReporteTieneStoryPoints(cerrados, planificados, completados int) error {
	me := m.reporte.visto.Metricas
	if me.SprintsCerrados != cerrados || me.Totales.StoryPointsPlanificados != planificados || me.Totales.StoryPointsCompletados != completados {
		return fmt.Errorf("proyecto: %d cerrados con %d/%d Story Points, se esperaban %d con %d/%d", me.SprintsCerrados,
			me.Totales.StoryPointsPlanificados, me.Totales.StoryPointsCompletados, cerrados, planificados, completados)
	}
	return nil
}

func (m *mundo) laDesviacionDelProyectoEnElReporteEsSinDatos() error {
	if d := m.reporte.visto.Metricas.Desviacion; d != nil {
		return fmt.Errorf("desviación = %+v, se esperaba \"sin datos\"", *d)
	}
	return nil
}

// ---------- verificaciones: defectos ----------

func (m *mundo) elReporteTieneDefectos(detectados, resueltos int) error {
	d := m.reporte.visto.Defectos
	if d.Detectados != detectados || d.Resueltos != resueltos {
		return fmt.Errorf("reporte: %d detectados y %d resueltos, se esperaban %d y %d", d.Detectados, d.Resueltos, detectados, resueltos)
	}
	return nil
}

func (m *mundo) elReporteMuestraLosDefectos(a, b string) error {
	lista := m.reporte.visto.Defectos.Lista
	if len(lista) != 2 || lista[0].Descripcion != a || lista[1].Descripcion != b {
		return fmt.Errorf("defectos = %+v, se esperaban %q y %q en ese orden", lista, a, b)
	}
	return nil
}

func (m *mundo) elDefectoDelReporteEstaResueltoEn(descripcion, estado string, sprint int) error {
	d, err := m.defectoDelReporte(descripcion)
	if err != nil {
		return err
	}
	if d.Estado != estado || d.SprintResolucion == nil || d.SprintResolucion.Numero != sprint {
		return fmt.Errorf("defecto = %+v, se esperaba %q en el Sprint %d", d, estado, sprint)
	}
	return nil
}

func (m *mundo) elDefectoDelReporteEstaSinResolucion(descripcion, estado string) error {
	d, err := m.defectoDelReporte(descripcion)
	if err != nil {
		return err
	}
	if d.Estado != estado || d.SprintResolucion != nil {
		return fmt.Errorf("defecto = %+v, se esperaba %q sin Sprint de resolución", d, estado)
	}
	return nil
}

// ---------- verificaciones: PDF y sin datos ----------

func (m *mundo) reciboUnArchivoPDF(nombre string) error {
	if m.reporte.pdf == nil {
		return fmt.Errorf("no se descargó el PDF: código %d, %s", m.respuesta.Code, m.respuesta.Body.String())
	}
	if !bytes.HasPrefix(m.reporte.pdf, []byte("%PDF-")) {
		return fmt.Errorf("el archivo descargado no es un PDF")
	}
	if m.reporte.nombrePDF != nombre {
		return fmt.Errorf("el archivo se llama %q, se esperaba %q", m.reporte.nombrePDF, nombre)
	}
	return nil
}

// aWinAnsi pasa el texto a como lo escribe el PDF (cp1252), para buscarlo en el archivo.
var aWinAnsi = fpdf.New("P", "mm", "A4", "").UnicodeTranslatorFromDescriptor("")

func (m *mundo) elPDFIncluyeElTexto(texto string) error {
	if !bytes.Contains(m.reporte.pdf, []byte(aWinAnsi(texto))) {
		return fmt.Errorf("el PDF no incluye el texto %q", texto)
	}
	return nil
}

func (m *mundo) elPDFIncluyeLasSecciones(a, b, c, d string) error {
	for _, s := range []string{a, b, c, d} {
		if err := m.elPDFIncluyeElTexto(s); err != nil {
			return err
		}
	}
	return nil
}

func (m *mundo) elReporteNoTieneDatos() error {
	r := m.reporte.visto
	if len(r.Historias) != 0 || len(r.Sprints) != 0 || len(r.Defectos.Lista) != 0 || r.Proyecto.ID == 0 {
		return fmt.Errorf("reporte = %+v, se esperaba sin historias, Sprints ni defectos", r)
	}
	return nil
}

func (m *mundo) elReporteDelSprintNoTieneHistorias(porcentaje int) error {
	s, err := m.sprintDelReporteDeSprint()
	if err != nil {
		return err
	}
	if len(m.reporte.visto.Historias) != 0 || s.Avance.PorcentajeCompletadas != porcentaje {
		return fmt.Errorf("reporte del Sprint = %+v, se esperaba sin historias y %d %%", m.reporte.visto, porcentaje)
	}
	return nil
}
