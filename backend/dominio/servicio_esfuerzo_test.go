package dominio_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/memoria"
)

// esfuerzoDePrueba arma el servicio con repositorios en memoria: el proyecto "Tienda Online" (del
// 2026-11-01 al 2026-12-15) con la integrante Ana y la historia "Registrar usuario", otro proyecto
// "App Turnos" con la integrante Carla, y hoy = 2026-11-20.
type esfuerzoDePrueba struct {
	servicio  *dominio.ServicioEsfuerzo
	historias *memoria.RepositorioHistorias
	tienda    int64
	turnos    int64
	historia  int64 // "Registrar usuario", de la tienda
	ana       int64 // integrante de la tienda
	carla     int64 // integrante de App Turnos
}

func nuevoEsfuerzo(t *testing.T) esfuerzoDePrueba {
	t.Helper()
	ctx := context.Background()
	proyectos := memoria.NuevoRepositorioProyectos()
	historias := memoria.NuevoRepositorioHistorias()
	integrantes := memoria.NuevoRepositorioIntegrantes()
	tienda, err := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "Tienda Online", Descripcion: "Venta de ropa", FechaInicio: f("2026-11-01"), FechaFin: f("2026-12-15")})
	if err != nil {
		t.Fatalf("crear proyecto de partida: %v", err)
	}
	turnos, _ := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "App Turnos", Descripcion: "Turnos médicos", FechaInicio: f("2026-11-05"), FechaFin: f("2027-01-20")})
	h, _ := historias.Guardar(ctx, dominio.Historia{ProyectoID: tienda.ID, Titulo: "Registrar usuario", Descripcion: "Alta", Prioridad: 1, Estado: dominio.EstadoPendiente})
	ana, _ := integrantes.Guardar(ctx, dominio.Integrante{ProyectoID: tienda.ID, Nombre: "Ana Pérez", Email: "ana@mail.com", Rol: dominio.RolDeveloper})
	carla, _ := integrantes.Guardar(ctx, dominio.Integrante{ProyectoID: turnos.ID, Nombre: "Carla Gómez", Email: "carla@mail.com", Rol: dominio.RolDeveloper})
	hoy := f("2026-11-20")
	return esfuerzoDePrueba{
		servicio:  dominio.NuevoServicioEsfuerzo(proyectos, historias, integrantes, memoria.NuevoRepositorioEsfuerzo(), func() time.Time { return hoy }),
		historias: historias,
		tienda:    tienda.ID,
		turnos:    turnos.ID,
		historia:  h.ID,
		ana:       ana.ID,
		carla:     carla.ID,
	}
}

func esfuerzo(integrante int64, fecha, actividad, horas string) dominio.DatosEsfuerzo {
	return dominio.DatosEsfuerzo{IntegranteID: integrante, Fecha: fecha, Actividad: actividad, Horas: horas}
}

// spec 09 CA-1, CA-2, RN-8: registrar guarda el esfuerzo con un identificador y el nombre del integrante.
func TestEsfuerzoRegistrar(t *testing.T) {
	e := nuevoEsfuerzo(t)
	ctx := context.Background()

	r, err := e.servicio.Registrar(ctx, e.tienda, e.historia, esfuerzo(e.ana, "2026-11-18", "Diseño del formulario", "2.5"))
	if err != nil {
		t.Fatalf("Registrar: %v", err)
	}
	if r.ID == 0 || r.ProyectoID != e.tienda || r.HistoriaID != e.historia || r.IntegranteID != e.ana ||
		r.IntegranteNombre != "Ana Pérez" || r.Actividad != "Diseño del formulario" || r.Centesimas != 250 ||
		!r.Fecha.Equal(f("2026-11-18")) {
		t.Errorf("registro = %+v, no conserva los datos o no tiene identificador", r)
	}
}

// spec 09 CA-4, RN-6, RN-7: consultar devuelve los registros por fecha, con nombre, y el total exacto.
func TestEsfuerzoConsultar(t *testing.T) {
	e := nuevoEsfuerzo(t)
	ctx := context.Background()

	vacio, err := e.servicio.Consultar(ctx, e.tienda, e.historia)
	if err != nil || len(vacio.Registros) != 0 || vacio.TotalCentesimas != 0 {
		t.Fatalf("Consultar sin registros = %+v, %v; se esperaba vacío y total 0", vacio, err)
	}
	e.servicio.Registrar(ctx, e.tienda, e.historia, esfuerzo(e.ana, "2026-11-19", "Backend", "2"))
	e.servicio.Registrar(ctx, e.tienda, e.historia, esfuerzo(e.ana, "2026-11-18", "Diseño", "1.5"))
	e.servicio.Registrar(ctx, e.tienda, e.historia, esfuerzo(e.ana, "2026-11-19", "Revisión", "0.25"))

	c, err := e.servicio.Consultar(ctx, e.tienda, e.historia)
	if err != nil {
		t.Fatalf("Consultar: %v", err)
	}
	var actividades []string
	for _, r := range c.Registros {
		actividades = append(actividades, r.Actividad)
		if r.IntegranteNombre != "Ana Pérez" {
			t.Errorf("integrante = %q, se esperaba Ana Pérez", r.IntegranteNombre)
		}
	}
	if len(actividades) != 3 || actividades[0] != "Diseño" || actividades[1] != "Backend" || actividades[2] != "Revisión" {
		t.Errorf("actividades = %v, se esperaba Diseño, Backend y Revisión", actividades)
	}
	if c.TotalCentesimas != 375 {
		t.Errorf("total = %d centésimas, se esperaba 375", c.TotalCentesimas)
	}
}

// spec 09 CA-4: cada historia tiene su propio esfuerzo.
func TestEsfuerzoPorHistoria(t *testing.T) {
	e := nuevoEsfuerzo(t)
	ctx := context.Background()
	otra, _ := e.historias.Guardar(ctx, dominio.Historia{ProyectoID: e.tienda, Titulo: "Ver catálogo", Descripcion: "Listado", Prioridad: 2, Estado: dominio.EstadoPendiente})
	e.servicio.Registrar(ctx, e.tienda, e.historia, esfuerzo(e.ana, "2026-11-18", "Backend", "2"))
	e.servicio.Registrar(ctx, e.tienda, otra.ID, esfuerzo(e.ana, "2026-11-18", "Catálogo", "0.5"))

	c, err := e.servicio.Consultar(ctx, e.tienda, otra.ID)
	if err != nil || len(c.Registros) != 1 || c.Registros[0].Actividad != "Catálogo" || c.TotalCentesimas != 50 {
		t.Errorf("esfuerzo de la otra historia = %+v, %v; se esperaba solo Catálogo con 0,5 horas", c, err)
	}
}

// spec 09 CA-3, CA-5: si un dato es inválido, no se registra nada.
func TestEsfuerzoInvalidoNoRegistra(t *testing.T) {
	e := nuevoEsfuerzo(t)
	ctx := context.Background()
	casos := []struct {
		nombre string
		datos  dominio.DatosEsfuerzo
		err    error
	}{
		{"horas en cero (CA-3)", esfuerzo(e.ana, "2026-11-18", "Backend", "0"), dominio.ErrHorasInvalidas},
		{"fecha futura (CA-5)", esfuerzo(e.ana, "2026-11-21", "Backend", "1"), dominio.ErrFechaFutura},
		{"actividad vacía (CA-5)", esfuerzo(e.ana, "2026-11-18", "", "1"), dominio.ErrActividadObligatoria},
		{"sin integrante (CA-8)", esfuerzo(0, "2026-11-18", "Backend", "1"), dominio.ErrIntegranteObligatorio},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if _, err := e.servicio.Registrar(ctx, e.tienda, e.historia, c.datos); !errors.Is(err, c.err) {
				t.Fatalf("error = %v, se esperaba %v", err, c.err)
			}
			got, _ := e.servicio.Consultar(ctx, e.tienda, e.historia)
			if len(got.Registros) != 0 {
				t.Errorf("hay %d registros, no se debía registrar ninguno", len(got.Registros))
			}
		})
	}
}

// spec 09 CA-6, RN-8: eliminar saca el registro del total; eliminarlo otra vez informa que no existe.
func TestEsfuerzoEliminar(t *testing.T) {
	e := nuevoEsfuerzo(t)
	ctx := context.Background()
	backend, _ := e.servicio.Registrar(ctx, e.tienda, e.historia, esfuerzo(e.ana, "2026-11-18", "Backend", "2"))
	e.servicio.Registrar(ctx, e.tienda, e.historia, esfuerzo(e.ana, "2026-11-19", "Pruebas", "1"))

	if err := e.servicio.Eliminar(ctx, e.tienda, e.historia, backend.ID); err != nil {
		t.Fatalf("Eliminar: %v", err)
	}
	c, _ := e.servicio.Consultar(ctx, e.tienda, e.historia)
	if len(c.Registros) != 1 || c.Registros[0].Actividad != "Pruebas" || c.TotalCentesimas != 100 {
		t.Errorf("después de eliminar = %+v, se esperaba solo Pruebas con 1 hora", c)
	}
	if err := e.servicio.Eliminar(ctx, e.tienda, e.historia, backend.ID); !errors.Is(err, dominio.ErrRegistroEsfuerzoNoExiste) {
		t.Errorf("Eliminar otra vez: error = %v, se esperaba %v", err, dominio.ErrRegistroEsfuerzoNoExiste)
	}
}

// spec 09 CA-8, RN-1, RN-2, RN-9: proyecto, historia o integrante inexistentes (o de otro proyecto).
func TestEsfuerzoInexistentes(t *testing.T) {
	e := nuevoEsfuerzo(t)
	ctx := context.Background()
	valido := func(integrante int64) dominio.DatosEsfuerzo {
		return esfuerzo(integrante, "2026-11-18", "Backend", "1")
	}

	casos := []struct {
		nombre string
		hacer  func() error
		err    error
	}{
		{"registrar en proyecto inexistente", func() error {
			_, err := e.servicio.Registrar(ctx, 999, e.historia, valido(e.ana))
			return err
		}, dominio.ErrProyectoNoExiste},
		{"registrar en historia inexistente", func() error {
			_, err := e.servicio.Registrar(ctx, e.tienda, 999, valido(e.ana))
			return err
		}, dominio.ErrHistoriaNoExiste},
		{"registrar en historia de otro proyecto", func() error {
			_, err := e.servicio.Registrar(ctx, e.turnos, e.historia, valido(e.carla))
			return err
		}, dominio.ErrHistoriaNoExiste},
		{"integrante inexistente", func() error {
			_, err := e.servicio.Registrar(ctx, e.tienda, e.historia, valido(999))
			return err
		}, dominio.ErrIntegranteNoExiste},
		{"integrante de otro proyecto", func() error {
			_, err := e.servicio.Registrar(ctx, e.tienda, e.historia, valido(e.carla))
			return err
		}, dominio.ErrIntegranteNoExiste},
		{"historia inexistente antes que datos inválidos (RN-9)", func() error {
			_, err := e.servicio.Registrar(ctx, e.tienda, 999, esfuerzo(0, "", "", ""))
			return err
		}, dominio.ErrHistoriaNoExiste},
		{"consultar en proyecto inexistente", func() error {
			_, err := e.servicio.Consultar(ctx, 999, e.historia)
			return err
		}, dominio.ErrProyectoNoExiste},
		{"eliminar en historia inexistente", func() error {
			return e.servicio.Eliminar(ctx, e.tienda, 999, 1)
		}, dominio.ErrHistoriaNoExiste},
		{"eliminar registro inexistente", func() error {
			return e.servicio.Eliminar(ctx, e.tienda, e.historia, 999)
		}, dominio.ErrRegistroEsfuerzoNoExiste},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if err := c.hacer(); !errors.Is(err, c.err) {
				t.Errorf("error = %v, se esperaba %v", err, c.err)
			}
		})
	}
	if c, _ := e.servicio.Consultar(ctx, e.tienda, e.historia); len(c.Registros) != 0 {
		t.Errorf("hay %d registros, no se debía registrar ninguno", len(c.Registros))
	}
}
