package dominio_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

func puntos(n int) *int { return &n }

// datosDeReporte arma los datos de un proyecto con tres historias y dos Sprints, como los
// devuelven las specs 03, 10, 11 y 12:
//   - Sprint 1 (cerrado, ID 10): "Registrar usuario" (completada) y "Ver catálogo". "Ver catálogo"
//     tenía 5 Story Points al cerrarlo y ahora tiene 8.
//   - Sprint 2 (abierto, ID 20): "Pagar pedido".
//   - Defectos: A detectado en el 1 y resuelto en el 2, B detectado en el 2, C detectado en el 1.
func datosDeReporte() dominio.DatosReporte {
	p := dominio.Proyecto{ID: 1, Nombre: "Tienda Online", FechaInicio: f("2026-11-01"), FechaFin: f("2026-12-31")}
	// En desorden, para verificar que el reporte usa el orden del backlog.
	historias := []dominio.Historia{
		{ID: 3, Titulo: "Pagar pedido", Prioridad: 3, Estado: dominio.EstadoEnProgreso},
		{ID: 1, Titulo: "Registrar usuario", Prioridad: 1, Estado: dominio.EstadoCompletada, StoryPoints: puntos(5), HorasEstimadas: horas(800)},
		{ID: 2, Titulo: "Ver catálogo", Prioridad: 2, Estado: dominio.EstadoPendiente, StoryPoints: puntos(8)},
	}
	s1 := dominio.Sprint{ID: 10, Numero: 1, FechaInicio: f("2026-11-01"), FechaFin: f("2026-11-14"), Goal: "Alta", Estado: dominio.SprintCerrado,
		Historias: []int64{1, 2}, Completadas: []int64{1}, StoryPointsAlCerrar: map[int64]*int{1: puntos(5), 2: puntos(5)}}
	s2 := dominio.Sprint{ID: 20, Numero: 2, FechaInicio: f("2026-11-16"), FechaFin: f("2026-11-28"), Goal: "Pago", Estado: dominio.SprintAbierto,
		Historias: []int64{3}, Completadas: []int64{}}
	desv := &dominio.Desviacion{Centesimas: 200, Porcentaje: 25}
	return dominio.DatosReporte{
		Estado:      dominio.EstadoDeProyecto{Proyecto: p, Estado: dominio.EstadoEnCurso},
		Integrantes: []dominio.Integrante{{ID: 1, Nombre: "Ana", Rol: dominio.RolDeveloper}},
		Historias:   historias,
		Sprints:     []dominio.Sprint{s2, s1},
		Avance: dominio.MetricasProyecto{
			Sprints: []dominio.MetricasSprint{
				{SprintID: 10, Numero: 1, Estado: dominio.SprintCerrado, HistoriasPlanificadas: 2, HistoriasCompletadas: 1, PorcentajeCompletadas: 50, StoryPointsPlanificados: 10, StoryPointsCompletados: 5},
				{SprintID: 20, Numero: 2, Estado: dominio.SprintAbierto, Parcial: true, HistoriasPlanificadas: 1, SinEstimar: 1},
			},
			SprintsCerrados: 1, HayVelocidad: true, VelocidadDecimas: 50,
			Totales: dominio.TotalesSprints{HistoriasPlanificadas: 2, HistoriasCompletadas: 1, PorcentajeCompletadas: 50, StoryPointsPlanificados: 10, StoryPointsCompletados: 5},
		},
		Esfuerzo: dominio.MetricasEsfuerzo{
			Historias: []dominio.EsfuerzoHistoria{
				{HistoriaID: 1, Titulo: "Registrar usuario", Estado: dominio.EstadoCompletada, Estimadas: horas(800), Reales: 1000, Desviacion: desv},
				{HistoriaID: 2, Titulo: "Ver catálogo", Estado: dominio.EstadoPendiente, Reales: 200},
				{HistoriaID: 3, Titulo: "Pagar pedido", Estado: dominio.EstadoEnProgreso},
			},
			Sprints: []dominio.EsfuerzoSprint{
				{SprintID: 10, Numero: 1, Estado: dominio.SprintCerrado, Estimadas: 800, Reales: 1000, Desviacion: desv},
				{SprintID: 20, Numero: 2, Estado: dominio.SprintAbierto, Parcial: true},
			},
			Estimadas: 800, Reales: 1200, Desviacion: desv,
		},
		Defectos: []dominio.Defecto{
			{ID: 1, Descripcion: "A", Estado: dominio.DefectoResuelto, SprintDeteccionID: 10, SprintResolucionID: 20},
			{ID: 2, Descripcion: "B", Estado: dominio.DefectoAbierto, SprintDeteccionID: 20},
			{ID: 3, Descripcion: "C", Estado: dominio.DefectoAbierto, SprintDeteccionID: 10},
		},
		MetricasDefectos: dominio.MetricasDefectos{Detectados: 3, Resueltos: 1, Sprints: []dominio.DefectosDeSprint{
			{SprintID: 10, Numero: 1, Detectados: 2, Resueltos: 0},
			{SprintID: 20, Numero: 2, Detectados: 1, Resueltos: 1},
		}},
		Hoy: f("2026-11-20"),
	}
}

func titulos(historias []dominio.HistoriaDeReporte) []string {
	lista := []string{}
	for _, h := range historias {
		lista = append(lista, h.Titulo)
	}
	return lista
}

func descripciones(defectos []dominio.Defecto) []string {
	lista := []string{}
	for _, d := range defectos {
		lista = append(lista, d.Descripcion)
	}
	return lista
}

// spec 14 CA-1 a CA-5, RN-2 a RN-4, RN-9: el reporte del proyecto junta los datos sin recalcularlos.
func TestArmarReporteProyecto(t *testing.T) {
	d := datosDeReporte()
	r := dominio.ArmarReporteProyecto(d)

	if r.Proyecto.Nombre != "Tienda Online" || r.Estado != dominio.EstadoEnCurso || !r.GeneradoEl.Equal(f("2026-11-20")) ||
		len(r.Integrantes) != 1 || r.Integrantes[0].Nombre != "Ana" {
		t.Errorf("encabezado = %+v %q %v %+v", r.Proyecto, r.Estado, r.GeneradoEl, r.Integrantes)
	}

	// CA-1, CA-2, CA-3: todo el backlog en su orden, con los Story Points actuales y el esfuerzo de la spec 12.
	if got := titulos(r.Historias); !slices.Equal(got, []string{"Registrar usuario", "Ver catálogo", "Pagar pedido"}) {
		t.Fatalf("historias = %v; se esperaba el orden del backlog", got)
	}
	registrar, catalogo := r.Historias[0], r.Historias[1]
	if registrar.Estado != dominio.EstadoCompletada || *registrar.StoryPoints != 5 || *registrar.Estimadas != 800 ||
		registrar.Reales != 1000 || registrar.Desviacion == nil || registrar.Desviacion.Porcentaje != 25 {
		t.Errorf("Registrar usuario = %+v", registrar)
	}
	if *catalogo.StoryPoints != 8 || catalogo.Estimadas != nil || catalogo.Reales != 200 || catalogo.Desviacion != nil {
		t.Errorf("Ver catálogo = %+v; se esperaban los Story Points actuales (8), sin horas estimadas y sin desviación", catalogo)
	}
	if r.Historias[2].StoryPoints != nil {
		t.Errorf("Pagar pedido = %+v; se esperaba sin estimar", r.Historias[2])
	}

	// CA-1, CA-4: los Sprints por fecha de inicio, con sus métricas de las specs 10 y 12.
	if len(r.Sprints) != 2 || r.Sprints[0].Sprint.ID != 10 || r.Sprints[1].Sprint.ID != 20 {
		t.Fatalf("Sprints = %+v; se esperaban el 1 y el 2 en ese orden", r.Sprints)
	}
	if r.Sprints[0].Avance.HistoriasCompletadas != 1 || r.Sprints[0].Avance.StoryPointsPlanificados != 10 ||
		r.Sprints[0].Esfuerzo.Reales != 1000 || r.Sprints[0].Sprint.Goal != "Alta" {
		t.Errorf("Sprint 1 = %+v", r.Sprints[0])
	}
	if !r.Sprints[1].Avance.Parcial || !r.Sprints[1].Esfuerzo.Parcial {
		t.Errorf("Sprint 2 = %+v; se esperaba parcial", r.Sprints[1])
	}

	// CA-4: métricas del proyecto.
	if !r.HayVelocidad || r.VelocidadDecimas != 50 || r.SprintsCerrados != 1 || r.Totales.StoryPointsCompletados != 5 ||
		r.HorasEstimadas != 800 || r.HorasReales != 1200 || r.Desviacion == nil || r.CompletadasSinEstimar != 0 {
		t.Errorf("métricas del proyecto = %+v", r)
	}

	// CA-5: todos los defectos, en orden de registro, con las cantidades de la spec 11.
	if r.Defectos.Detectados != 3 || r.Defectos.Resueltos != 1 || !slices.Equal(descripciones(r.Defectos.Lista), []string{"A", "B", "C"}) {
		t.Errorf("defectos = %+v", r.Defectos)
	}
}

// spec 14 CA-1 a CA-5, RN-3, RN-5: el reporte de un Sprint cerrado usa los Story Points congelados,
// marca las completadas en él y trae solo los defectos detectados o resueltos en él.
func TestArmarReporteSprintCerrado(t *testing.T) {
	r, err := dominio.ArmarReporteSprint(datosDeReporte(), 10)
	if err != nil {
		t.Fatalf("ArmarReporteSprint: %v", err)
	}
	if r.Proyecto.Nombre != "Tienda Online" || r.Sprint.Sprint.Numero != 1 || r.Sprint.Avance.Parcial || !r.GeneradoEl.Equal(f("2026-11-20")) {
		t.Errorf("encabezado = %+v %+v", r.Proyecto, r.Sprint)
	}
	if got := titulos(r.Historias); !slices.Equal(got, []string{"Registrar usuario", "Ver catálogo"}) {
		t.Fatalf("historias = %v; se esperaban las planificadas en el orden del backlog", got)
	}
	registrar, catalogo := r.Historias[0], r.Historias[1]
	if !registrar.CompletadaEnSprint || registrar.Desviacion == nil || registrar.Reales != 1000 {
		t.Errorf("Registrar usuario = %+v; se esperaba completada en el Sprint y con desviación", registrar)
	}
	if catalogo.CompletadaEnSprint || catalogo.StoryPoints == nil || *catalogo.StoryPoints != 5 || catalogo.Desviacion != nil {
		t.Errorf("Ver catálogo = %+v; se esperaba no completada, con los 5 Story Points congelados y sin desviación", catalogo)
	}
	if r.Sprint.Avance.PorcentajeCompletadas != 50 || r.Sprint.Esfuerzo.Estimadas != 800 {
		t.Errorf("métricas del Sprint = %+v", r.Sprint)
	}
	// Detectados en el 1: A y C; resuelto en el 1: ninguno.
	if r.Defectos.Detectados != 2 || r.Defectos.Resueltos != 0 || !slices.Equal(descripciones(r.Defectos.Lista), []string{"A", "C"}) {
		t.Errorf("defectos = %+v; se esperaban A y C, 2 detectados y 0 resueltos", r.Defectos)
	}
}

// spec 14 RN-3, RN-5, RN-6: el reporte de un Sprint abierto es parcial y usa los Story Points
// actuales; un defecto resuelto en él aparece aunque se haya detectado en otro.
func TestArmarReporteSprintAbierto(t *testing.T) {
	d := datosDeReporte()
	d.Historias[0].StoryPoints = puntos(3) // "Pagar pedido" se estima mientras el Sprint está abierto

	r, err := dominio.ArmarReporteSprint(d, 20)
	if err != nil {
		t.Fatalf("ArmarReporteSprint: %v", err)
	}
	if !r.Sprint.Avance.Parcial {
		t.Errorf("Sprint = %+v; se esperaba parcial", r.Sprint)
	}
	if len(r.Historias) != 1 || r.Historias[0].StoryPoints == nil || *r.Historias[0].StoryPoints != 3 || r.Historias[0].CompletadaEnSprint {
		t.Errorf("historias = %+v; se esperaba Pagar pedido con sus 3 Story Points actuales", r.Historias)
	}
	if r.Defectos.Detectados != 1 || r.Defectos.Resueltos != 1 || !slices.Equal(descripciones(r.Defectos.Lista), []string{"A", "B"}) {
		t.Errorf("defectos = %+v; se esperaban A (resuelto acá) y B (detectado acá)", r.Defectos)
	}
}

// spec 14 RN-1: un Sprint que no es del proyecto no existe.
func TestArmarReporteSprintInexistente(t *testing.T) {
	if _, err := dominio.ArmarReporteSprint(datosDeReporte(), 999); !errors.Is(err, dominio.ErrSprintNoExiste) {
		t.Errorf("error = %v, se esperaba %v", err, dominio.ErrSprintNoExiste)
	}
}

// spec 14 RN-7 (CA-7): sin datos igual hay reporte, con listas vacías (no nil, para el JSON).
func TestArmarReporteSinDatos(t *testing.T) {
	p := dominio.Proyecto{ID: 1, Nombre: "App Turnos"}
	d := dominio.DatosReporte{Estado: dominio.EstadoDeProyecto{Proyecto: p, Estado: dominio.EstadoNoIniciado}, Hoy: f("2026-11-20")}

	r := dominio.ArmarReporteProyecto(d)
	if r.Integrantes == nil || r.Historias == nil || r.Sprints == nil || r.Defectos.Lista == nil ||
		len(r.Historias)+len(r.Sprints)+len(r.Defectos.Lista) != 0 || r.HayVelocidad || r.Desviacion != nil {
		t.Errorf("reporte = %+v; se esperaban listas vacías, sin velocidad y desviación sin datos", r)
	}

	d.Sprints = []dominio.Sprint{{ID: 10, Numero: 1, Estado: dominio.SprintAbierto}}
	d.Avance.Sprints = []dominio.MetricasSprint{{SprintID: 10, Numero: 1, Parcial: true}}
	d.Esfuerzo.Sprints = []dominio.EsfuerzoSprint{{SprintID: 10, Numero: 1, Parcial: true}}
	d.MetricasDefectos.Sprints = []dominio.DefectosDeSprint{{SprintID: 10, Numero: 1}}
	rs, err := dominio.ArmarReporteSprint(d, 10)
	if err != nil || rs.Historias == nil || len(rs.Historias) != 0 || rs.Defectos.Lista == nil || rs.Sprint.Avance.PorcentajeCompletadas != 0 {
		t.Errorf("reporte del Sprint = %+v, %v; se esperaban listas vacías y 0 %%", rs, err)
	}
}
