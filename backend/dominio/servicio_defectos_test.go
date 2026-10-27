package dominio_test

import (
	"context"
	"errors"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/memoria"
)

// defectosDePrueba arma el servicio con repositorios en memoria: el proyecto "Tienda Online" con la
// historia "Registrar usuario" y dos Sprints (del 2026-11-01 al 14 y del 16 al 28), y otro proyecto
// "App Turnos" con su historia y su Sprint.
type defectosDePrueba struct {
	servicio      *dominio.ServicioDefectos
	tienda        int64
	turnos        int64
	historia      int64 // "Registrar usuario", de la tienda
	sprint1       int64
	sprint2       int64
	historiaTurno int64 // de App Turnos
	sprintTurno   int64 // de App Turnos
}

func nuevoDefectos(t *testing.T) defectosDePrueba {
	t.Helper()
	ctx := context.Background()
	proyectos := memoria.NuevoRepositorioProyectos()
	historias := memoria.NuevoRepositorioHistorias()
	sprints := memoria.NuevoRepositorioSprints()
	tienda, err := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "Tienda Online", Descripcion: "Venta de ropa", FechaInicio: f("2026-11-01"), FechaFin: f("2026-12-15")})
	if err != nil {
		t.Fatalf("crear proyecto de partida: %v", err)
	}
	turnos, _ := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "App Turnos", Descripcion: "Turnos médicos", FechaInicio: f("2026-11-05"), FechaFin: f("2027-01-20")})
	h, _ := historias.Guardar(ctx, dominio.Historia{ProyectoID: tienda.ID, Titulo: "Registrar usuario", Descripcion: "Alta", Prioridad: 1, Estado: dominio.EstadoPendiente})
	ht, _ := historias.Guardar(ctx, dominio.Historia{ProyectoID: turnos.ID, Titulo: "Pedir turno", Descripcion: "Turno", Prioridad: 1, Estado: dominio.EstadoPendiente})
	s1, _ := sprints.Guardar(ctx, dominio.Sprint{ProyectoID: tienda.ID, Numero: 1, FechaInicio: f("2026-11-01"), FechaFin: f("2026-11-14"), Goal: "Alta", Estado: dominio.SprintAbierto})
	s2, _ := sprints.Guardar(ctx, dominio.Sprint{ProyectoID: tienda.ID, Numero: 2, FechaInicio: f("2026-11-16"), FechaFin: f("2026-11-28"), Goal: "Compra", Estado: dominio.SprintAbierto})
	st, _ := sprints.Guardar(ctx, dominio.Sprint{ProyectoID: turnos.ID, Numero: 1, FechaInicio: f("2026-11-05"), FechaFin: f("2026-11-20"), Goal: "Turnos", Estado: dominio.SprintAbierto})
	return defectosDePrueba{
		servicio:      dominio.NuevoServicioDefectos(proyectos, historias, sprints, memoria.NuevoRepositorioDefectos()),
		tienda:        tienda.ID,
		turnos:        turnos.ID,
		historia:      h.ID,
		sprint1:       s1.ID,
		sprint2:       s2.ID,
		historiaTurno: ht.ID,
		sprintTurno:   st.ID,
	}
}

// spec 11 CA-1, RN-6, RN-12: registrar guarda el defecto con un identificador, Abierto y sin
// resolución, con el título de la historia y el número del Sprint.
func TestDefectosRegistrar(t *testing.T) {
	d := nuevoDefectos(t)
	ctx := context.Background()

	got, err := d.servicio.Registrar(ctx, d.tienda, defecto("El alta acepta emails sin @", "Alta", d.historia, d.sprint1))
	if err != nil {
		t.Fatalf("Registrar: %v", err)
	}
	if got.ID == 0 || got.ProyectoID != d.tienda || got.HistoriaID != d.historia || got.HistoriaTitulo != "Registrar usuario" ||
		got.SprintDeteccionID != d.sprint1 || got.SprintDeteccionNumero != 1 || got.Estado != dominio.DefectoAbierto ||
		got.SprintResolucionID != 0 || got.SprintResolucionNumero != 0 {
		t.Errorf("defecto = %+v, no conserva los datos o no arranca Abierto y sin resolución", got)
	}
}

// spec 11 CA-5, CA-7, RN-1, RN-13: el proyecto se verifica antes que los datos, los datos antes que
// la existencia de la historia y del Sprint, y una historia o un Sprint de otro proyecto no existen.
// Si algo falla, no se registra nada.
func TestDefectosRegistrarErrores(t *testing.T) {
	d := nuevoDefectos(t)
	ctx := context.Background()

	casos := []struct {
		nombre   string
		proyecto int64
		datos    dominio.DatosDefecto
		err      error
	}{
		{"proyecto inexistente, con datos inválidos", 999, defecto("", "", 0, 0), dominio.ErrProyectoNoExiste},
		{"dato inválido antes que la historia inexistente", d.tienda, defecto("Falla", "x", 999, d.sprint1), dominio.ErrSeveridadInvalida},
		{"sin Sprint de detección antes que la historia inexistente", d.tienda, defecto("Falla", "Alta", 999, 0), dominio.ErrSprintDeteccionObligatorio},
		{"historia inexistente", d.tienda, defecto("Falla", "Alta", 999, d.sprint1), dominio.ErrHistoriaNoExiste},
		{"historia de otro proyecto", d.tienda, defecto("Falla", "Alta", d.historiaTurno, d.sprint1), dominio.ErrHistoriaNoExiste},
		{"historia inexistente antes que el Sprint inexistente", d.tienda, defecto("Falla", "Alta", 999, 999), dominio.ErrHistoriaNoExiste},
		{"Sprint inexistente", d.tienda, defecto("Falla", "Alta", d.historia, 999), dominio.ErrSprintNoExiste},
		{"Sprint de otro proyecto", d.tienda, defecto("Falla", "Alta", d.historia, d.sprintTurno), dominio.ErrSprintNoExiste},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if _, err := d.servicio.Registrar(ctx, c.proyecto, c.datos); !errors.Is(err, c.err) {
				t.Errorf("error = %v, se esperaba %v", err, c.err)
			}
		})
	}
	lista, _ := d.servicio.Listar(ctx, d.tienda)
	if len(lista) != 0 {
		t.Errorf("hay %d defectos, no se debía registrar ninguno", len(lista))
	}
}

// spec 11 CA-5, RN-12: se listan en el orden en que se registraron, cada proyecto los suyos, y se
// puede ver uno; uno de otro proyecto no existe.
func TestDefectosListarYObtener(t *testing.T) {
	d := nuevoDefectos(t)
	ctx := context.Background()

	vacio, err := d.servicio.Listar(ctx, d.tienda)
	if err != nil || vacio == nil || len(vacio) != 0 {
		t.Fatalf("Listar sin defectos = %v, %v; se esperaba una lista vacía", vacio, err)
	}
	primero, _ := d.servicio.Registrar(ctx, d.tienda, defecto("Primero", "Baja", d.historia, d.sprint2))
	segundo, _ := d.servicio.Registrar(ctx, d.tienda, defecto("Segundo", "Crítica", d.historia, d.sprint1))
	ajeno, _ := d.servicio.Registrar(ctx, d.turnos, defecto("De turnos", "Media", d.historiaTurno, d.sprintTurno))

	lista, err := d.servicio.Listar(ctx, d.tienda)
	if err != nil || len(lista) != 2 || lista[0].ID != primero.ID || lista[1].ID != segundo.ID {
		t.Fatalf("Listar = %+v, %v; se esperaban %q y %q en ese orden", lista, err, "Primero", "Segundo")
	}
	if lista[0].HistoriaTitulo != "Registrar usuario" || lista[0].SprintDeteccionNumero != 2 || lista[1].SprintDeteccionNumero != 1 {
		t.Errorf("Listar no completa la historia y los Sprints: %+v", lista)
	}

	got, err := d.servicio.Obtener(ctx, d.tienda, segundo.ID)
	if err != nil || got.Descripcion != "Segundo" || got.Severidad != dominio.SeveridadCritica || got.SprintDeteccionNumero != 1 {
		t.Errorf("Obtener = %+v, %v; se esperaba el defecto %q", got, err, "Segundo")
	}
	if _, err := d.servicio.Obtener(ctx, d.tienda, ajeno.ID); !errors.Is(err, dominio.ErrDefectoNoExiste) {
		t.Errorf("Obtener de otro proyecto: error = %v, se esperaba %v", err, dominio.ErrDefectoNoExiste)
	}
	if _, err := d.servicio.Obtener(ctx, d.tienda, 999); !errors.Is(err, dominio.ErrDefectoNoExiste) {
		t.Errorf("Obtener inexistente: error = %v, se esperaba %v", err, dominio.ErrDefectoNoExiste)
	}
	if _, err := d.servicio.Listar(ctx, 999); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Listar de proyecto inexistente: error = %v, se esperaba %v", err, dominio.ErrProyectoNoExiste)
	}
}

// spec 11 CA-2, CA-3, RN-7 a RN-9: resolver registra el Sprint de resolución y reabrir lo borra; el
// cambio queda guardado.
func TestDefectosCambiarEstado(t *testing.T) {
	d := nuevoDefectos(t)
	ctx := context.Background()
	falla, _ := d.servicio.Registrar(ctx, d.tienda, defecto("Falla", "Alta", d.historia, d.sprint1))

	resuelto, err := d.servicio.CambiarEstado(ctx, d.tienda, falla.ID, "Resuelto", d.sprint2)
	if err != nil {
		t.Fatalf("resolver: %v", err)
	}
	if resuelto.Estado != dominio.DefectoResuelto || resuelto.SprintResolucionID != d.sprint2 || resuelto.SprintResolucionNumero != 2 {
		t.Errorf("resuelto = %+v; se esperaba Resuelto en el Sprint 2", resuelto)
	}
	guardado, _ := d.servicio.Obtener(ctx, d.tienda, falla.ID)
	if guardado.Estado != dominio.DefectoResuelto || guardado.SprintResolucionID != d.sprint2 {
		t.Errorf("el cambio no quedó guardado: %+v", guardado)
	}

	reabierto, err := d.servicio.CambiarEstado(ctx, d.tienda, falla.ID, "Abierto", 0)
	if err != nil || reabierto.Estado != dominio.DefectoAbierto || reabierto.SprintResolucionID != 0 || reabierto.SprintResolucionNumero != 0 {
		t.Errorf("reabrir = %+v, %v; se esperaba Abierto y sin resolución", reabierto, err)
	}

	// RN-9: a En corrección con un Sprint indicado, el Sprint no se tiene en cuenta (ni se busca).
	enCorreccion, err := d.servicio.CambiarEstado(ctx, d.tienda, falla.ID, "En corrección", 999)
	if err != nil || enCorreccion.Estado != dominio.DefectoEnCorreccion || enCorreccion.SprintResolucionID != 0 {
		t.Errorf("En corrección = %+v, %v; se esperaba En corrección y sin resolución", enCorreccion, err)
	}
}

// spec 11 CA-2, CA-3, CA-7: errores al cambiar el estado; si algo falla, el defecto no cambia.
func TestDefectosCambiarEstadoErrores(t *testing.T) {
	d := nuevoDefectos(t)
	ctx := context.Background()
	falla, _ := d.servicio.Registrar(ctx, d.tienda, defecto("Falla", "Alta", d.historia, d.sprint2))

	casos := []struct {
		nombre   string
		proyecto int64
		defecto  int64
		estado   string
		sprint   int64
		err      error
	}{
		{"proyecto inexistente", 999, falla.ID, "Resuelto", d.sprint2, dominio.ErrProyectoNoExiste},
		{"defecto inexistente", d.tienda, 999, "Resuelto", d.sprint2, dominio.ErrDefectoNoExiste},
		{"estado inválido", d.tienda, falla.ID, "Cerrado", 0, dominio.ErrEstadoInvalido},
		{"resolver sin Sprint", d.tienda, falla.ID, "Resuelto", 0, dominio.ErrSprintResolucionObligatorio},
		{"resolver en un Sprint inexistente", d.tienda, falla.ID, "Resuelto", 999, dominio.ErrSprintNoExiste},
		{"resolver en un Sprint de otro proyecto", d.tienda, falla.ID, "Resuelto", d.sprintTurno, dominio.ErrSprintNoExiste},
		{"resolver antes de la detección", d.tienda, falla.ID, "Resuelto", d.sprint1, dominio.ErrResolucionAnteriorADeteccion},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if _, err := d.servicio.CambiarEstado(ctx, c.proyecto, c.defecto, c.estado, c.sprint); !errors.Is(err, c.err) {
				t.Errorf("error = %v, se esperaba %v", err, c.err)
			}
		})
	}
	got, _ := d.servicio.Obtener(ctx, d.tienda, falla.ID)
	if got.Estado != dominio.DefectoAbierto || got.SprintResolucionID != 0 {
		t.Errorf("el defecto cambió: %+v", got)
	}
}

// spec 11 CA-4, RN-10, RN-11: las métricas cuentan los defectos del proyecto con todos sus Sprints.
func TestDefectosMetricas(t *testing.T) {
	d := nuevoDefectos(t)
	ctx := context.Background()
	a, _ := d.servicio.Registrar(ctx, d.tienda, defecto("A", "Alta", d.historia, d.sprint1))
	d.servicio.Registrar(ctx, d.tienda, defecto("B", "Baja", d.historia, d.sprint1))
	c, _ := d.servicio.Registrar(ctx, d.tienda, defecto("C", "Media", d.historia, d.sprint2))
	d.servicio.Registrar(ctx, d.turnos, defecto("De turnos", "Media", d.historiaTurno, d.sprintTurno))
	d.servicio.CambiarEstado(ctx, d.tienda, a.ID, "Resuelto", d.sprint2)
	d.servicio.CambiarEstado(ctx, d.tienda, c.ID, "Resuelto", d.sprint2)

	m, err := d.servicio.Metricas(ctx, d.tienda)
	if err != nil {
		t.Fatalf("Metricas: %v", err)
	}
	if m.Detectados != 3 || m.Resueltos != 2 || len(m.Sprints) != 2 ||
		m.Sprints[0] != (dominio.DefectosDeSprint{SprintID: d.sprint1, Numero: 1, Detectados: 2, Resueltos: 0}) ||
		m.Sprints[1] != (dominio.DefectosDeSprint{SprintID: d.sprint2, Numero: 2, Detectados: 1, Resueltos: 2}) {
		t.Errorf("métricas = %+v; se esperaban 3 detectados y 2 resueltos (Sprint 1: 2 y 0, Sprint 2: 1 y 2)", m)
	}
	if _, err := d.servicio.Metricas(ctx, 999); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Metricas de proyecto inexistente: error = %v, se esperaba %v", err, dominio.ErrProyectoNoExiste)
	}
}
