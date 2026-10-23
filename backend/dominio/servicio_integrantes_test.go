package dominio_test

import (
	"context"
	"errors"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/memoria"
)

// equipoDePrueba arma el servicio con repositorios en memoria y dos proyectos ya creados.
type equipoDePrueba struct {
	servicio *dominio.ServicioIntegrantes
	tienda   int64 // proyecto "Tienda Online"
	turnos   int64 // proyecto "App Turnos"
}

func nuevoEquipo(t *testing.T) equipoDePrueba {
	t.Helper()
	ctx := context.Background()
	proyectos := memoria.NuevoRepositorioProyectos()
	tienda, err := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "Tienda Online", Descripcion: "Venta de ropa", FechaInicio: f("2026-11-01"), FechaFin: f("2026-12-15")})
	if err != nil {
		t.Fatalf("crear proyecto de partida: %v", err)
	}
	turnos, _ := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "App Turnos", Descripcion: "Turnos médicos", FechaInicio: f("2026-11-05"), FechaFin: f("2027-01-20")})
	return equipoDePrueba{
		servicio: dominio.NuevoServicioIntegrantes(proyectos, memoria.NuevoRepositorioIntegrantes()),
		tienda:   tienda.ID,
		turnos:   turnos.ID,
	}
}

func integrante(nombre, email, rol string) dominio.DatosIntegrante {
	return dominio.DatosIntegrante{Nombre: nombre, Email: email, Rol: rol}
}

// cantidadDeIntegrantes lista el proyecto y devuelve cuántos integrantes tiene.
func (e equipoDePrueba) cantidadDeIntegrantes(t *testing.T, proyectoID int64) int {
	t.Helper()
	lista, err := e.servicio.Listar(context.Background(), proyectoID)
	if err != nil {
		t.Fatalf("Listar: %v", err)
	}
	return len(lista)
}

// spec 02 CA-1, RN-3, RN-8: al agregar, el sistema asigna el identificador y normaliza el email.
func TestIntegrantesAgregar(t *testing.T) {
	e := nuevoEquipo(t)
	ctx := context.Background()

	ana, err := e.servicio.Agregar(ctx, e.tienda, integrante("Ana Pérez", " Ana@Mail.COM ", "Developer"))
	if err != nil {
		t.Fatalf("Agregar: %v", err)
	}
	bruno, err := e.servicio.Agregar(ctx, e.tienda, integrante("Bruno Díaz", "bruno@mail.com", "Scrum Master"))
	if err != nil {
		t.Fatalf("Agregar el segundo: %v", err)
	}
	if ana.ID == 0 || bruno.ID == 0 || ana.ID == bruno.ID {
		t.Errorf("IDs = %d y %d, se esperaban distintos de cero y entre sí", ana.ID, bruno.ID)
	}
	esperado := dominio.Integrante{ID: ana.ID, ProyectoID: e.tienda, Nombre: "Ana Pérez", Email: "ana@mail.com", Rol: dominio.RolDeveloper}
	if ana != esperado {
		t.Errorf("integrante = %+v, se esperaba %+v", ana, esperado)
	}
}

// spec 02 CA-4, CA-5: si un dato es inválido, no se registra nada.
func TestIntegrantesAgregarInvalidoNoRegistra(t *testing.T) {
	e := nuevoEquipo(t)
	ctx := context.Background()

	casos := []struct {
		nombre string
		datos  dominio.DatosIntegrante
		err    error
	}{
		{"nombre vacío (CA-4)", integrante("", "ana@mail.com", "Developer"), dominio.ErrNombreObligatorio},
		{"email inválido (CA-4)", integrante("Ana Pérez", "ana.mail.com", "Developer"), dominio.ErrEmailInvalido},
		{"rol inválido (CA-5)", integrante("Ana Pérez", "ana@mail.com", "Tester"), dominio.ErrRolInvalido},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if _, err := e.servicio.Agregar(ctx, e.tienda, c.datos); !errors.Is(err, c.err) {
				t.Fatalf("error = %v, se esperaba %v", err, c.err)
			}
			if n := e.cantidadDeIntegrantes(t, e.tienda); n != 0 {
				t.Errorf("hay %d integrantes, no se debía registrar ninguno", n)
			}
		})
	}
}

// spec 02 CA-3, RN-6: el mismo email (normalizado) no se registra dos veces en un proyecto.
func TestIntegrantesEmailRepetido(t *testing.T) {
	for _, email := range []string{"ana@mail.com", "ANA@Mail.com", " ana@mail.com ", "\tana@MAIL.com"} {
		t.Run(email, func(t *testing.T) {
			e := nuevoEquipo(t)
			ctx := context.Background()
			if _, err := e.servicio.Agregar(ctx, e.tienda, integrante("Ana Pérez", "ana@mail.com", "Developer")); err != nil {
				t.Fatalf("Agregar de partida: %v", err)
			}

			_, err := e.servicio.Agregar(ctx, e.tienda, integrante("Ana P.", email, "Scrum Master"))
			if !errors.Is(err, dominio.ErrIntegranteYaRegistrado) {
				t.Fatalf("error = %v, se esperaba %v", err, dominio.ErrIntegranteYaRegistrado)
			}
			if n := e.cantidadDeIntegrantes(t, e.tienda); n != 1 {
				t.Errorf("hay %d integrantes, se esperaba 1", n)
			}
		})
	}
}

// spec 02 CA-3, RN-2, RN-6: el mismo email en otro proyecto y el mismo nombre con otro email se aceptan.
func TestIntegrantesRepetidosPermitidos(t *testing.T) {
	e := nuevoEquipo(t)
	ctx := context.Background()
	e.servicio.Agregar(ctx, e.tienda, integrante("Ana Pérez", "ana@mail.com", "Developer"))

	if _, err := e.servicio.Agregar(ctx, e.turnos, integrante("Ana Pérez", "ana@mail.com", "Product Owner")); err != nil {
		t.Errorf("mismo email en otro proyecto: %v", err)
	}
	if _, err := e.servicio.Agregar(ctx, e.tienda, integrante("Ana Pérez", "ana.perez@otro.com", "Developer")); err != nil {
		t.Errorf("mismo nombre con otro email: %v", err)
	}
	if n := e.cantidadDeIntegrantes(t, e.tienda); n != 2 {
		t.Errorf("la tienda tiene %d integrantes, se esperaban 2", n)
	}
}

// spec 02 RN-9: los datos se validan antes que el email repetido.
func TestIntegrantesDatoInvalidoAntesQueRepetido(t *testing.T) {
	e := nuevoEquipo(t)
	ctx := context.Background()
	e.servicio.Agregar(ctx, e.tienda, integrante("Ana Pérez", "ana@mail.com", "Developer"))

	_, err := e.servicio.Agregar(ctx, e.tienda, integrante("Ana P.", "ana@mail.com", "Tester"))
	if !errors.Is(err, dominio.ErrRolInvalido) {
		t.Errorf("error = %v, se esperaba %v", err, dominio.ErrRolInvalido)
	}
}

// spec 02 CA-2, RN-7: se listan en el orden en que se agregaron y cada proyecto ve solo los suyos.
func TestIntegrantesListar(t *testing.T) {
	e := nuevoEquipo(t)
	ctx := context.Background()

	vacio, err := e.servicio.Listar(ctx, e.tienda)
	if err != nil || len(vacio) != 0 {
		t.Fatalf("Listar sin integrantes = %v, %v; se esperaba lista vacía", vacio, err)
	}
	e.servicio.Agregar(ctx, e.tienda, integrante("Ana Pérez", "ana@mail.com", "Developer"))
	e.servicio.Agregar(ctx, e.turnos, integrante("Carla Gómez", "carla@mail.com", "Product Owner"))
	e.servicio.Agregar(ctx, e.tienda, integrante("Bruno Díaz", "bruno@mail.com", "Scrum Master"))

	lista, _ := e.servicio.Listar(ctx, e.tienda)
	if len(lista) != 2 || lista[0].Nombre != "Ana Pérez" || lista[1].Nombre != "Bruno Díaz" {
		t.Errorf("integrantes = %+v, se esperaban Ana Pérez y Bruno Díaz en ese orden", lista)
	}
}

// spec 02 CA-2, CA-6, RN-8: ver un integrante; si es de otro proyecto, no existe.
func TestIntegrantesObtener(t *testing.T) {
	e := nuevoEquipo(t)
	ctx := context.Background()
	ana, _ := e.servicio.Agregar(ctx, e.tienda, integrante("Ana Pérez", "ana@mail.com", "Developer"))

	visto, err := e.servicio.Obtener(ctx, e.tienda, ana.ID)
	if err != nil || visto != ana {
		t.Errorf("Obtener = %+v, %v; se esperaba %+v", visto, err, ana)
	}
	if _, err := e.servicio.Obtener(ctx, e.turnos, ana.ID); !errors.Is(err, dominio.ErrIntegranteNoExiste) {
		t.Errorf("Obtener desde otro proyecto: error = %v, se esperaba %v", err, dominio.ErrIntegranteNoExiste)
	}
	if _, err := e.servicio.Obtener(ctx, e.tienda, 999); !errors.Is(err, dominio.ErrIntegranteNoExiste) {
		t.Errorf("Obtener inexistente: error = %v, se esperaba %v", err, dominio.ErrIntegranteNoExiste)
	}
}

// spec 02 RN-1, CA-6: todas las operaciones sobre un proyecto que no existe informan que no existe.
func TestIntegrantesProyectoInexistente(t *testing.T) {
	e := nuevoEquipo(t)
	ctx := context.Background()
	const fantasma = 999

	if _, err := e.servicio.Agregar(ctx, fantasma, integrante("Ana Pérez", "ana@mail.com", "Developer")); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Agregar: error = %v, se esperaba %v", err, dominio.ErrProyectoNoExiste)
	}
	if _, err := e.servicio.Listar(ctx, fantasma); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Listar: error = %v, se esperaba %v", err, dominio.ErrProyectoNoExiste)
	}
	if _, err := e.servicio.Obtener(ctx, fantasma, 1); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Obtener: error = %v, se esperaba %v", err, dominio.ErrProyectoNoExiste)
	}
	// Un proyecto inexistente con datos inválidos informa primero el proyecto (RN-1).
	if _, err := e.servicio.Agregar(ctx, fantasma, integrante("", "", "")); !errors.Is(err, dominio.ErrProyectoNoExiste) {
		t.Errorf("Agregar con datos inválidos: error = %v, se esperaba %v", err, dominio.ErrProyectoNoExiste)
	}
}
