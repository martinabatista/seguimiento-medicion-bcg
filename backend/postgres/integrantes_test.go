package postgres

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// integrantesDePrueba arma el repositorio de integrantes sobre la base vacía, con dos proyectos
// creados. Como repoDePrueba, se saltea si no está DATABASE_URL.
func integrantesDePrueba(t *testing.T) (r *RepositorioIntegrantes, tienda, turnos int64) {
	t.Helper()
	proyectos := repoDePrueba(t)
	ctx := context.Background()
	a, err := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "Tienda Online", Descripcion: "Venta de ropa", FechaInicio: fecha("2026-11-01"), FechaFin: fecha("2026-12-15")})
	if err != nil {
		t.Fatalf("crear proyecto de partida: %v", err)
	}
	b, _ := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "App Turnos", Descripcion: "Turnos médicos", FechaInicio: fecha("2026-11-05"), FechaFin: fecha("2027-01-20")})
	return NuevoRepositorioIntegrantes(proyectos.db), a.ID, b.ID
}

// spec 02 CA-1 a CA-3, CA-6 contra PostgreSQL.
func TestRepositorioIntegrantes(t *testing.T) {
	r, tienda, turnos := integrantesDePrueba(t)
	ctx := context.Background()

	// Guardar: el ID lo genera PostgreSQL.
	ana, err := r.Guardar(ctx, dominio.Integrante{ProyectoID: tienda, Nombre: "Ana Pérez", Email: "ana@mail.com", Rol: dominio.RolDeveloper})
	if err != nil || ana.ID == 0 {
		t.Fatalf("Guardar = %+v, %v; se esperaba un ID asignado", ana, err)
	}
	// Datos en el límite, con eñes: VARCHAR cuenta caracteres, no bytes.
	limite, err := r.Guardar(ctx, dominio.Integrante{ProyectoID: tienda, Nombre: strings.Repeat("ñ", 100), Email: strings.Repeat("ñ", 245) + "@mail.com", Rol: dominio.RolScrumMaster})
	if err != nil {
		t.Fatalf("Guardar con datos en el límite: %v", err)
	}
	// El mismo email en otro proyecto se acepta (RN-6).
	carla, err := r.Guardar(ctx, dominio.Integrante{ProyectoID: turnos, Nombre: "Ana Pérez", Email: "ana@mail.com", Rol: dominio.RolProductOwner})
	if err != nil {
		t.Fatalf("Guardar el mismo email en otro proyecto: %v", err)
	}

	visto, err := r.Obtener(ctx, tienda, ana.ID)
	if err != nil || visto != ana {
		t.Errorf("Obtener = %+v, %v; se esperaba %+v", visto, err, ana)
	}

	// Listar: solo los del proyecto, en orden de alta (RN-7).
	lista, err := r.Listar(ctx, tienda)
	if err != nil || len(lista) != 2 || lista[0] != ana || lista[1] != limite {
		t.Errorf("Listar = %+v, %v; se esperaba [%+v %+v]", lista, err, ana, limite)
	}

	// Un integrante de otro proyecto o inexistente no existe (RN-8).
	if _, err := r.Obtener(ctx, tienda, carla.ID); !errors.Is(err, dominio.ErrIntegranteNoExiste) {
		t.Errorf("Obtener de otro proyecto: error = %v", err)
	}
	if _, err := r.Obtener(ctx, tienda, 999); !errors.Is(err, dominio.ErrIntegranteNoExiste) {
		t.Errorf("Obtener inexistente: error = %v", err)
	}
}

// spec 02 RN-6: la base rechaza el email repetido en el proyecto aunque no pase por el servicio
// (por ejemplo, dos altas al mismo tiempo).
func TestGuardarIntegranteRepetido(t *testing.T) {
	r, tienda, _ := integrantesDePrueba(t)
	ctx := context.Background()
	if _, err := r.Guardar(ctx, dominio.Integrante{ProyectoID: tienda, Nombre: "Ana Pérez", Email: "ana@mail.com", Rol: dominio.RolDeveloper}); err != nil {
		t.Fatalf("Guardar de partida: %v", err)
	}

	_, err := r.Guardar(ctx, dominio.Integrante{ProyectoID: tienda, Nombre: "Ana P.", Email: "ana@mail.com", Rol: dominio.RolScrumMaster})
	if !errors.Is(err, dominio.ErrIntegranteYaRegistrado) {
		t.Fatalf("error = %v, se esperaba %v", err, dominio.ErrIntegranteYaRegistrado)
	}
	if lista, _ := r.Listar(ctx, tienda); len(lista) != 1 {
		t.Errorf("hay %d integrantes, se esperaba 1", len(lista))
	}
}

// spec 02 CA-2: un proyecto sin integrantes tiene la lista vacía (no nil ni error).
func TestListarIntegrantesVacio(t *testing.T) {
	r, tienda, _ := integrantesDePrueba(t)
	lista, err := r.Listar(context.Background(), tienda)
	if err != nil || lista == nil || len(lista) != 0 {
		t.Errorf("Listar = %#v, %v; se esperaba una lista vacía", lista, err)
	}
}
