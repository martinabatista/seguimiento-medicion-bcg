package pdf

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/go-pdf/fpdf"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// aWinAnsi pasa el texto a como lo escribe el PDF (las fuentes estándar usan cp1252), para buscarlo
// en el archivo.
var aWinAnsi = fpdf.New("P", "mm", "A4", "").UnicodeTranslatorFromDescriptor("")

func contiene(t *testing.T, archivo []byte, textos ...string) {
	t.Helper()
	for _, texto := range textos {
		if !bytes.Contains(archivo, []byte(aWinAnsi(texto))) {
			t.Errorf("el PDF no incluye el texto %q", texto)
		}
	}
}

func esPDF(t *testing.T, archivo []byte, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("generar el PDF: %v", err)
	}
	if !bytes.HasPrefix(archivo, []byte("%PDF-")) {
		t.Fatalf("el archivo no es un PDF: empieza con %q", archivo[:min(len(archivo), 10)])
	}
}

func dia(s string) time.Time {
	d, _ := time.Parse(dominio.FormatoFecha, s)
	return d
}

func cinco() *int { n := 5; return &n }

func reporteProyecto() dominio.ReporteProyecto {
	desv := &dominio.Desviacion{Centesimas: 200, Porcentaje: 25}
	ochoHoras := 800
	return dominio.ReporteProyecto{
		Proyecto:    dominio.Proyecto{ID: 1, Nombre: "Gestión de Turnos", Descripcion: "Reservas en línea", FechaInicio: dia("2026-11-01"), FechaFin: dia("2026-12-31")},
		Estado:      dominio.EstadoEnCurso,
		Integrantes: []dominio.Integrante{{Nombre: "Ana Pérez", Rol: dominio.RolDeveloper}},
		GeneradoEl:  dia("2026-11-20"),
		Historias: []dominio.HistoriaDeReporte{
			{Titulo: "Registrar usuario", Estado: dominio.EstadoCompletada, StoryPoints: cinco(), Estimadas: &ochoHoras, Reales: 1000, Desviacion: desv},
			{Titulo: "Ver catálogo", Estado: dominio.EstadoPendiente},
		},
		Sprints: []dominio.SprintDeReporte{{
			Sprint:   dominio.Sprint{Numero: 1, FechaInicio: dia("2026-11-01"), FechaFin: dia("2026-11-14"), Goal: "Alta de clientes", Estado: dominio.SprintAbierto},
			Avance:   dominio.MetricasSprint{Numero: 1, Parcial: true, HistoriasPlanificadas: 2, HistoriasCompletadas: 1, PorcentajeCompletadas: 50},
			Esfuerzo: dominio.EsfuerzoSprint{Numero: 1, Parcial: true, Estimadas: 800, Reales: 1000, Desviacion: desv},
		}},
		HorasEstimadas: 800, HorasReales: 1000, Desviacion: desv,
		Defectos: dominio.DefectosDeReporte{Detectados: 1, Lista: []dominio.Defecto{{
			Descripcion: "El alta acepta emails sin @", Severidad: dominio.SeveridadCritica, Estado: dominio.DefectoAbierto,
			HistoriaTitulo: "Registrar usuario", SprintDeteccionNumero: 1,
		}}},
	}
}

// spec 14 CA-6, RN-8, RN-10: el PDF del proyecto tiene el título, las secciones y los datos, con
// tildes y ñ, y los textos que faltan se muestran como en la vista previa.
func TestReporteProyectoPDF(t *testing.T) {
	archivo, err := ReporteProyecto(reporteProyecto())
	esPDF(t, archivo, err)
	contiene(t, archivo,
		"Reporte del proyecto", "Gestión de Turnos", "En curso", "Generado el 2026-11-20", "Ana Pérez",
		"Historias", "Sprints", "Métricas", "Defectos",
		"Registrar usuario", "Ver catálogo", "sin estimar", "sin horas estimadas", "sin datos",
		"Alta de clientes", "parcial", "todavía no hay velocidad",
		"El alta acepta emails sin @", "Crítica",
	)
}

// spec 14 CA-6: el PDF del Sprint dice de qué Sprint es y marca las historias completadas en él.
func TestReporteSprintPDF(t *testing.T) {
	rp := reporteProyecto()
	r := dominio.ReporteSprint{
		Proyecto:   rp.Proyecto,
		GeneradoEl: rp.GeneradoEl,
		Sprint:     rp.Sprints[0],
		Historias: []dominio.HistoriaDeReporte{
			{Titulo: "Registrar usuario", Estado: dominio.EstadoCompletada, StoryPoints: cinco(), CompletadaEnSprint: true},
			{Titulo: "Ver catálogo", Estado: dominio.EstadoPendiente},
		},
		Defectos: rp.Defectos,
	}
	archivo, err := ReporteSprint(r)
	esPDF(t, archivo, err)
	contiene(t, archivo, "Reporte del Sprint 1", "Gestión de Turnos", "Alta de clientes", "parcial",
		"Historias planificadas", "Registrar usuario", "Métricas del Sprint", "Defectos")
}

// spec 14 RN-7: sin datos igual hay PDF, y cada sección dice que no hay.
func TestReporteSinDatosPDF(t *testing.T) {
	archivo, err := ReporteProyecto(dominio.ReporteProyecto{Proyecto: dominio.Proyecto{Nombre: "App Turnos"}, Estado: dominio.EstadoNoIniciado})
	esPDF(t, archivo, err)
	contiene(t, archivo, "App Turnos", "sin integrantes", "sin historias", "sin Sprints", "sin defectos", "todavía no hay velocidad")
}

// spec 14 RN-10: los textos largos siguen en varias líneas y, si no entran, en otra hoja.
func TestReporteLargoPDF(t *testing.T) {
	r := reporteProyecto()
	r.Defectos.Lista[0].Descripcion = strings.Repeat("palabra ", 125) // 1000 caracteres
	for i := 0; i < 80; i++ {
		r.Historias = append(r.Historias, dominio.HistoriaDeReporte{Titulo: "Historia de relleno", Estado: dominio.EstadoPendiente})
	}
	archivo, err := ReporteProyecto(r)
	esPDF(t, archivo, err)
	if paginas := bytes.Count(archivo, []byte("/Type /Page\n")); paginas < 2 {
		t.Errorf("el PDF tiene %d páginas, se esperaban varias", paginas)
	}
	if bytes.Contains(archivo, []byte(strings.Repeat("palabra ", 125))) {
		t.Error("la descripción larga quedó en una sola línea")
	}
}

// spec 14 RN-8: los números se muestran como en la pantalla de Métricas (coma decimal, signo).
func TestFormatoDeNumeros(t *testing.T) {
	casos := map[string]string{
		desviacion(&dominio.Desviacion{Centesimas: 200, Porcentaje: 25}): "+2 h (+25 %)",
		desviacion(&dominio.Desviacion{Centesimas: -50, Porcentaje: -6}): "-0,5 h (-6 %)",
		desviacion(&dominio.Desviacion{}):                                "0 h (0 %)",
		desviacion(nil):                                                  "sin datos",
		horas(250):                                                       "2,5",
		velocidad(65):                                                    "6,5",
		velocidad(50):                                                    "5",
	}
	for got, esperado := range casos {
		if got != esperado {
			t.Errorf("se mostró %q, se esperaba %q", got, esperado)
		}
	}
}
