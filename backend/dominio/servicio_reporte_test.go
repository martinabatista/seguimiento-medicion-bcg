package dominio_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/memoria"
)

// reporteDePrueba arma el servicio del reporte sobre repositorios en memoria, con los servicios de
// las specs 03, 10, 11 y 12 que reutiliza, y hoy = 2026-11-20.
type reporteDePrueba struct {
	servicio    *dominio.ServicioReporte
	metricas    *dominio.ServicioMetricas
	esfuerzo    *dominio.ServicioMetricasEsfuerzo
	defectos    *dominio.ServicioDefectos
	sprints     *dominio.ServicioSprints
	proyectos   *memoria.RepositorioProyectos
	historias   *memoria.RepositorioHistorias
	integrantes *memoria.RepositorioIntegrantes
	registros   *memoria.RepositorioEsfuerzo
}

func nuevoReporte() reporteDePrueba {
	hoy := func() time.Time { return f("2026-11-20") }
	proyectos := memoria.NuevoRepositorioProyectos()
	historias := memoria.NuevoRepositorioHistorias()
	sprints := memoria.NuevoRepositorioSprints()
	integrantes := memoria.NuevoRepositorioIntegrantes()
	registros := memoria.NuevoRepositorioEsfuerzo()
	estado := dominio.NuevoServicioEstado(proyectos, historias, sprints, hoy)
	metricas := dominio.NuevoServicioMetricas(proyectos, historias, sprints)
	esfuerzo := dominio.NuevoServicioMetricasEsfuerzo(proyectos, historias, sprints, registros)
	defectos := dominio.NuevoServicioDefectos(proyectos, historias, sprints, memoria.NuevoRepositorioDefectos())
	return reporteDePrueba{
		servicio:    dominio.NuevoServicioReporte(estado, metricas, esfuerzo, defectos, integrantes, historias, sprints, hoy),
		metricas:    metricas,
		esfuerzo:    esfuerzo,
		defectos:    defectos,
		sprints:     dominio.NuevoServicioSprints(proyectos, historias, sprints),
		proyectos:   proyectos,
		historias:   historias,
		integrantes: integrantes,
		registros:   registros,
	}
}

// spec 14 RN-2 (CA-1 a CA-5): el reporte da los mismos números que las consultas de métricas
// (spec 10), de esfuerzo (spec 12) y de defectos (spec 11), porque las reutiliza.
func TestReporteReutilizaLasMetricas(t *testing.T) {
	r := nuevoReporte()
	ctx := context.Background()
	p, _ := r.proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "Tienda Online", Descripcion: "d", FechaInicio: f("2026-11-01"), FechaFin: f("2026-12-31")})
	ana, _ := r.integrantes.Guardar(ctx, dominio.Integrante{ProyectoID: p.ID, Nombre: "Ana", Email: "ana@mail.com", Rol: dominio.RolDeveloper})
	registrar, _ := r.historias.Guardar(ctx, dominio.Historia{ProyectoID: p.ID, Titulo: "Registrar usuario", Descripcion: "d", Prioridad: 1, Estado: dominio.EstadoPendiente, StoryPoints: puntos(5), HorasEstimadas: horas(800)})
	catalogo, _ := r.historias.Guardar(ctx, dominio.Historia{ProyectoID: p.ID, Titulo: "Ver catálogo", Descripcion: "d", Prioridad: 2, Estado: dominio.EstadoPendiente, StoryPoints: puntos(3)})
	r.registros.Guardar(ctx, dominio.RegistroEsfuerzo{ProyectoID: p.ID, HistoriaID: registrar.ID, IntegranteID: ana.ID, Fecha: f("2026-11-05"), Actividad: "Backend", Centesimas: 1000})

	s1, err := r.sprints.Crear(ctx, p.ID, f("2026-11-01"), f("2026-11-14"), "Alta")
	if err != nil {
		t.Fatalf("crear Sprint: %v", err)
	}
	r.sprints.AsignarHistoria(ctx, p.ID, s1.ID, registrar.ID)
	r.sprints.AsignarHistoria(ctx, p.ID, s1.ID, catalogo.ID)
	r.sprints.CompletarHistoria(ctx, p.ID, s1.ID, registrar.ID)
	r.sprints.Cerrar(ctx, p.ID, s1.ID)
	r.defectos.Registrar(ctx, p.ID, dominio.DatosDefecto{Descripcion: "Falla", Severidad: "Alta", HistoriaID: registrar.ID, SprintDeteccionID: s1.ID})

	rp, err := r.servicio.ReporteProyecto(ctx, p.ID)
	if err != nil {
		t.Fatalf("ReporteProyecto: %v", err)
	}
	m, _ := r.metricas.Consultar(ctx, p.ID)
	me, _ := r.esfuerzo.Consultar(ctx, p.ID)
	md, _ := r.defectos.Metricas(ctx, p.ID)

	if rp.Proyecto.ID != p.ID || rp.Estado != dominio.EstadoEnCurso || !rp.GeneradoEl.Equal(f("2026-11-20")) ||
		len(rp.Integrantes) != 1 || rp.Integrantes[0].Nombre != "Ana" {
		t.Errorf("encabezado = %+v %q %v %+v", rp.Proyecto, rp.Estado, rp.GeneradoEl, rp.Integrantes)
	}
	if len(rp.Sprints) != 1 || rp.Sprints[0].Avance != m.Sprints[0] || rp.VelocidadDecimas != m.VelocidadDecimas || rp.Totales != m.Totales {
		t.Errorf("avance del reporte = %+v, se esperaba el de la spec 10 %+v", rp.Sprints, m)
	}
	if rp.HorasReales != me.Reales || rp.HorasEstimadas != me.Estimadas || rp.Desviacion == nil || *rp.Desviacion != *me.Desviacion {
		t.Errorf("esfuerzo del reporte = %d/%d/%v, se esperaba el de la spec 12 %+v", rp.HorasEstimadas, rp.HorasReales, rp.Desviacion, me)
	}
	if rp.Defectos.Detectados != md.Detectados || len(rp.Defectos.Lista) != 1 || rp.Defectos.Lista[0].HistoriaTitulo != "Registrar usuario" {
		t.Errorf("defectos del reporte = %+v, se esperaba el de la spec 11 con la historia completada", rp.Defectos)
	}

	rs, err := r.servicio.ReporteSprint(ctx, p.ID, s1.ID)
	if err != nil {
		t.Fatalf("ReporteSprint: %v", err)
	}
	if rs.Sprint.Avance != m.Sprints[0] || len(rs.Historias) != 2 || !rs.Historias[0].CompletadaEnSprint || rs.Defectos.Detectados != 1 {
		t.Errorf("reporte del Sprint = %+v", rs)
	}
}

// spec 14 RN-1 (CA-7): el proyecto tiene que existir y el Sprint ser de ese proyecto.
func TestReporteInexistente(t *testing.T) {
	r := nuevoReporte()
	ctx := context.Background()
	tienda, _ := r.proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "Tienda Online", Descripcion: "d", FechaInicio: f("2026-11-01"), FechaFin: f("2026-12-31")})
	turnos, _ := r.proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "App Turnos", Descripcion: "d", FechaInicio: f("2026-11-05"), FechaFin: f("2027-01-20")})
	ajeno, _ := r.sprints.Crear(ctx, turnos.ID, f("2026-11-05"), f("2026-11-20"), "Turnos")

	if _, err := r.servicio.ReporteProyecto(ctx, 999); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("ReporteProyecto inexistente: error = %v", err)
	}
	if _, err := r.servicio.ReporteSprint(ctx, 999, ajeno.ID); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("ReporteSprint de proyecto inexistente: error = %v", err)
	}
	if _, err := r.servicio.ReporteSprint(ctx, tienda.ID, 999); !errors.Is(err, dominio.ErrSprintNoExiste) {
		t.Errorf("ReporteSprint inexistente: error = %v", err)
	}
	if _, err := r.servicio.ReporteSprint(ctx, tienda.ID, ajeno.ID); !errors.Is(err, dominio.ErrSprintNoExiste) {
		t.Errorf("ReporteSprint de otro proyecto: error = %v", err)
	}
}
