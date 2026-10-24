package dominio_test

import (
	"context"
	"errors"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/memoria"
)

// backlogDePrueba arma el servicio con repositorios en memoria y dos proyectos ya creados.
type backlogDePrueba struct {
	servicio  *dominio.ServicioBacklog
	historias *memoria.RepositorioHistorias
	esfuerzo  *memoria.RepositorioEsfuerzo
	tienda    int64 // proyecto "Tienda Online"
	turnos    int64 // proyecto "App Turnos"
}

func nuevoBacklog(t *testing.T) backlogDePrueba {
	t.Helper()
	ctx := context.Background()
	proyectos := memoria.NuevoRepositorioProyectos()
	historias := memoria.NuevoRepositorioHistorias()
	tienda, err := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "Tienda Online", Descripcion: "Venta de ropa", FechaInicio: f("2026-11-01"), FechaFin: f("2026-12-15")})
	if err != nil {
		t.Fatalf("crear proyecto de partida: %v", err)
	}
	turnos, _ := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "App Turnos", Descripcion: "Turnos médicos", FechaInicio: f("2026-11-05"), FechaFin: f("2027-01-20")})
	esfuerzo := memoria.NuevoRepositorioEsfuerzo()
	return backlogDePrueba{
		servicio:  dominio.NuevoServicioBacklog(proyectos, historias, memoria.NuevoRepositorioSprints(), esfuerzo),
		historias: historias,
		esfuerzo:  esfuerzo,
		tienda:    tienda.ID,
		turnos:    turnos.ID,
	}
}

func datos(titulo string, prioridad int) dominio.DatosHistoria {
	return dominio.DatosHistoria{Titulo: titulo, Descripcion: "Alta de clientes", Prioridad: prioridad, Criterios: "Valida el email"}
}

// spec 04 CA-1, RN-6, RN-10: al crear, el sistema asigna el identificador; arranca Pendiente y sin estimar.
func TestBacklogCrear(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()

	h1, err := b.servicio.Crear(ctx, b.tienda, datos("Registrar usuario", 1))
	if err != nil {
		t.Fatalf("Crear: %v", err)
	}
	h2, err := b.servicio.Crear(ctx, b.tienda, datos("Registrar usuario", 1))
	if err != nil {
		t.Fatalf("Crear con título repetido: %v", err)
	}
	if h1.ID == 0 || h2.ID == 0 || h1.ID == h2.ID {
		t.Errorf("IDs = %d y %d, se esperaban distintos de cero y entre sí", h1.ID, h2.ID)
	}
	if h1.ProyectoID != b.tienda || h1.Estado != dominio.EstadoPendiente || h1.StoryPoints != nil {
		t.Errorf("historia = %+v, se esperaba del proyecto %d, Pendiente y sin estimar", h1, b.tienda)
	}
}

// spec 04 CA-1, CA-5: si un dato es inválido, no se registra nada.
func TestBacklogCrearInvalidoNoRegistra(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()

	_, err := b.servicio.Crear(ctx, b.tienda, datos("", 1))
	if !errors.Is(err, dominio.ErrTituloObligatorio) {
		t.Fatalf("error = %v, se esperaba %v", err, dominio.ErrTituloObligatorio)
	}
	lista, _ := b.servicio.Listar(ctx, b.tienda)
	if len(lista) != 0 {
		t.Errorf("hay %d historias, no se debía registrar ninguna", len(lista))
	}
}

// spec 04 CA-2: modificar todos los datos y el estado; el ID, el proyecto y los Story Points no cambian.
func TestBacklogModificar(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()
	sp := 5 // como si ya la hubieran estimado (spec 07): modificar no los pierde
	h, _ := b.historias.Guardar(ctx, dominio.Historia{ProyectoID: b.tienda, Titulo: "Registrar usuario", Descripcion: "Alta de clientes", Prioridad: 3, Estado: dominio.EstadoPendiente, StoryPoints: &sp})

	nuevos := dominio.DatosHistoria{Titulo: "Registrar cliente", Descripcion: "Alta con email", Prioridad: 1, Criterios: "Envía un email"}
	m, err := b.servicio.Modificar(ctx, b.tienda, h.ID, nuevos, "En progreso")
	if err != nil {
		t.Fatalf("Modificar: %v", err)
	}
	visto, _ := b.servicio.Obtener(ctx, b.tienda, h.ID)
	if visto.ID != h.ID || visto.ProyectoID != b.tienda || visto.Titulo != "Registrar cliente" ||
		visto.Descripcion != "Alta con email" || visto.Prioridad != 1 || visto.Criterios != "Envía un email" ||
		visto.Estado != dominio.EstadoEnProgreso {
		t.Errorf("historia = %+v, no quedó modificada", visto)
	}
	if visto.StoryPoints == nil || *visto.StoryPoints != 5 || m.StoryPoints == nil || *m.StoryPoints != 5 {
		t.Errorf("Story Points = %v, modificar no debe perder la estimación", visto.StoryPoints)
	}
}

// spec 04 RN-9, CA-2, CA-6: una modificación inválida deja la historia como estaba.
func TestBacklogModificarInvalidoNoCambia(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()
	h, _ := b.servicio.Crear(ctx, b.tienda, datos("Registrar usuario", 3))

	casos := []struct {
		nombre string
		datos  dominio.DatosHistoria
		estado string
		err    error
	}{
		{"prioridad 0", datos("Registrar usuario", 0), "Pendiente", dominio.ErrPrioridadInvalida},
		{"estado que no existe", datos("Registrar usuario", 1), "Terminada", dominio.ErrEstadoInvalido},
		{"estado en minúscula", datos("Registrar usuario", 1), "pendiente", dominio.ErrEstadoInvalido},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			_, err := b.servicio.Modificar(ctx, b.tienda, h.ID, c.datos, c.estado)
			if !errors.Is(err, c.err) {
				t.Fatalf("error = %v, se esperaba %v", err, c.err)
			}
			visto, _ := b.servicio.Obtener(ctx, b.tienda, h.ID)
			if visto.Prioridad != 3 || visto.Estado != dominio.EstadoPendiente {
				t.Errorf("historia = %+v, se esperaba sin cambios", visto)
			}
		})
	}
}

// spec 04 CA-3: eliminar saca la historia del backlog; eliminarla otra vez informa que no existe.
func TestBacklogEliminar(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()
	registrar, _ := b.servicio.Crear(ctx, b.tienda, datos("Registrar usuario", 1))
	catalogo, _ := b.servicio.Crear(ctx, b.tienda, datos("Ver catálogo", 2))

	if err := b.servicio.Eliminar(ctx, b.tienda, registrar.ID); err != nil {
		t.Fatalf("Eliminar: %v", err)
	}
	lista, _ := b.servicio.Listar(ctx, b.tienda)
	if len(lista) != 1 || lista[0].ID != catalogo.ID {
		t.Errorf("backlog = %+v, se esperaba solo %q", lista, catalogo.Titulo)
	}
	if err := b.servicio.Eliminar(ctx, b.tienda, registrar.ID); !errors.Is(err, dominio.ErrHistoriaNoExiste) {
		t.Errorf("Eliminar otra vez: error = %v, se esperaba %v", err, dominio.ErrHistoriaNoExiste)
	}
}

// spec 04 CA-4, RN-8: el backlog se lista ordenado y cada proyecto ve solo el suyo.
func TestBacklogListar(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()

	vacio, err := b.servicio.Listar(ctx, b.tienda)
	if err != nil || len(vacio) != 0 {
		t.Fatalf("Listar sin historias = %v, %v; se esperaba lista vacía", vacio, err)
	}
	b.servicio.Crear(ctx, b.tienda, datos("Ver catálogo", 2))
	b.servicio.Crear(ctx, b.tienda, datos("Buscar producto", 2))
	b.servicio.Crear(ctx, b.tienda, datos("Registrar usuario", 1))
	b.servicio.Crear(ctx, b.turnos, datos("Pedir turno", 1))

	lista, _ := b.servicio.Listar(ctx, b.tienda)
	var titulos []string
	for _, h := range lista {
		titulos = append(titulos, h.Titulo)
	}
	esperado := []string{"Registrar usuario", "Ver catálogo", "Buscar producto"}
	if len(titulos) != len(esperado) || titulos[0] != esperado[0] || titulos[1] != esperado[1] || titulos[2] != esperado[2] {
		t.Errorf("backlog = %v, se esperaba %v", titulos, esperado)
	}
}

// spec 04 CA-4, CA-7, RN-10: ver una historia; si es de otro proyecto, no existe.
func TestBacklogObtener(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()
	h, _ := b.servicio.Crear(ctx, b.tienda, datos("Registrar usuario", 1))

	visto, err := b.servicio.Obtener(ctx, b.tienda, h.ID)
	if err != nil || visto.ID != h.ID || visto.Titulo != "Registrar usuario" {
		t.Errorf("Obtener = %+v, %v; se esperaba %+v", visto, err, h)
	}
	if _, err := b.servicio.Obtener(ctx, b.turnos, h.ID); !errors.Is(err, dominio.ErrHistoriaNoExiste) {
		t.Errorf("Obtener desde otro proyecto: error = %v, se esperaba %v", err, dominio.ErrHistoriaNoExiste)
	}
	if _, err := b.servicio.Obtener(ctx, b.tienda, 999); !errors.Is(err, dominio.ErrHistoriaNoExiste) {
		t.Errorf("Obtener inexistente: error = %v, se esperaba %v", err, dominio.ErrHistoriaNoExiste)
	}
}

// spec 04 RN-1, CA-7: todas las operaciones sobre un proyecto que no existe informan que no existe.
func TestBacklogProyectoInexistente(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()
	h, _ := b.servicio.Crear(ctx, b.tienda, datos("Registrar usuario", 1))
	const fantasma = 999

	if _, err := b.servicio.Crear(ctx, fantasma, datos("Registrar usuario", 1)); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Crear: error = %v", err)
	}
	if _, err := b.servicio.Listar(ctx, fantasma); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Listar: error = %v", err)
	}
	if _, err := b.servicio.Obtener(ctx, fantasma, h.ID); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Obtener: error = %v", err)
	}
	if _, err := b.servicio.Modificar(ctx, fantasma, h.ID, datos("x", 1), "Pendiente"); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Modificar: error = %v", err)
	}
	if err := b.servicio.Eliminar(ctx, fantasma, h.ID); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Eliminar: error = %v", err)
	}
	if _, err := b.servicio.Modificar(ctx, b.tienda, 999, datos("x", 1), "Pendiente"); !errors.Is(err, dominio.ErrHistoriaNoExiste) {
		t.Errorf("Modificar historia inexistente: error = %v", err)
	}
}

// spec 09 RN-10 (CA-7): no se elimina del backlog una historia con esfuerzo registrado; sin
// registros, sí.
func TestBacklogNoEliminaHistoriaConEsfuerzo(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()
	h, _ := b.servicio.Crear(ctx, b.tienda, datos("Registrar usuario", 1))
	r, _ := b.esfuerzo.Guardar(ctx, dominio.RegistroEsfuerzo{ProyectoID: b.tienda, HistoriaID: h.ID, IntegranteID: 1, Fecha: f("2026-11-18"), Actividad: "Backend", Centesimas: 200})

	if err := b.servicio.Eliminar(ctx, b.tienda, h.ID); !errors.Is(err, dominio.ErrHistoriaConEsfuerzo) {
		t.Fatalf("Eliminar con esfuerzo: error = %v, se esperaba %v", err, dominio.ErrHistoriaConEsfuerzo)
	}
	if _, err := b.servicio.Obtener(ctx, b.tienda, h.ID); err != nil {
		t.Fatalf("la historia ya no existe: %v", err)
	}

	b.esfuerzo.Eliminar(ctx, b.tienda, h.ID, r.ID)
	if err := b.servicio.Eliminar(ctx, b.tienda, h.ID); err != nil {
		t.Errorf("Eliminar sin esfuerzo: %v", err)
	}
}
