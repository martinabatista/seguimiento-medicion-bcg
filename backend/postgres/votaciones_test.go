package postgres

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// spec 08: votaciones con sus rondas, votos y estimación acordada contra PostgreSQL.
func TestRepositorioVotaciones(t *testing.T) {
	historias, tienda, turnos := historiasDePrueba(t)
	ctx := context.Background()
	integrantes := NuevoRepositorioIntegrantes(historias.db)
	ana, err := integrantes.Guardar(ctx, dominio.Integrante{ProyectoID: tienda, Nombre: "Ana", Email: "ana@mail.com", Rol: dominio.RolDeveloper})
	if err != nil {
		t.Fatalf("agregar integrante de partida: %v", err)
	}
	beto, _ := integrantes.Guardar(ctx, dominio.Integrante{ProyectoID: tienda, Nombre: "Beto", Email: "beto@mail.com", Rol: dominio.RolDeveloper})
	h, _ := historias.Guardar(ctx, dominio.Historia{ProyectoID: tienda, Titulo: "Registrar usuario", Descripcion: "d", Prioridad: 1, Estado: dominio.EstadoPendiente})
	r := NuevoRepositorioVotaciones(historias.db)

	// Guardar: el ID lo genera PostgreSQL; arranca en la ronda 1 sin votos.
	v, err := r.Guardar(ctx, dominio.Votacion{ProyectoID: tienda, HistoriaID: h.ID, Estado: dominio.VotacionAbierta, Rondas: []dominio.Ronda{{Numero: 1}}})
	if err != nil || v.ID == 0 {
		t.Fatalf("Guardar = %+v, %v; se esperaba un ID asignado", v, err)
	}
	visto, err := r.Obtener(ctx, tienda, v.ID)
	if err != nil || !reflect.DeepEqual(visto, v) {
		t.Fatalf("Obtener = %+v, %v; se esperaba %+v", visto, err, v)
	}

	// Actualizar: una ronda revelada con dos votos y otra en curso con uno.
	v.Rondas = []dominio.Ronda{
		{Numero: 1, Revelada: true, Votos: []dominio.Voto{{IntegranteID: ana.ID, Valor: 3}, {IntegranteID: beto.ID, Valor: 13}}},
		{Numero: 2, Votos: []dominio.Voto{{IntegranteID: ana.ID, Valor: 5}}},
	}
	if err := r.Actualizar(ctx, v); err != nil {
		t.Fatalf("Actualizar con rondas: %v", err)
	}
	if visto, _ = r.Obtener(ctx, tienda, v.ID); !reflect.DeepEqual(visto, v) {
		t.Errorf("después de votar = %+v, se esperaba %+v", visto, v)
	}

	// Cerrar con la estimación acordada.
	acordada := 5
	v.Rondas[1].Votos = append(v.Rondas[1].Votos, dominio.Voto{IntegranteID: beto.ID, Valor: 5})
	v.Rondas[1].Revelada = true
	v.Estado, v.Acordada = dominio.VotacionCerrada, &acordada
	if err := r.Actualizar(ctx, v); err != nil {
		t.Fatalf("Actualizar al cerrar: %v", err)
	}
	if visto, _ = r.Obtener(ctx, tienda, v.ID); !reflect.DeepEqual(visto, v) {
		t.Errorf("después de cerrar = %+v, se esperaba %+v", visto, v)
	}

	// Listar por proyecto; una votación de otro proyecto no existe (RN-11).
	if lista, err := r.Listar(ctx, tienda); err != nil || len(lista) != 1 || !reflect.DeepEqual(lista[0], v) {
		t.Errorf("Listar = %+v, %v; se esperaba solo la votación de la tienda", lista, err)
	}
	if lista, err := r.Listar(ctx, turnos); err != nil || lista == nil || len(lista) != 0 {
		t.Errorf("Listar otro proyecto = %#v, %v; se esperaba una lista vacía", lista, err)
	}
	if _, err := r.Obtener(ctx, turnos, v.ID); !errors.Is(err, dominio.ErrVotacionNoExiste) {
		t.Errorf("Obtener desde otro proyecto: error = %v", err)
	}
	otra := v
	otra.ProyectoID = turnos
	if err := r.Actualizar(ctx, otra); !errors.Is(err, dominio.ErrVotacionNoExiste) {
		t.Errorf("Actualizar desde otro proyecto: error = %v", err)
	}

	// Si se elimina la historia del backlog, sus votaciones se eliminan con ella.
	if err := historias.Eliminar(ctx, tienda, h.ID); err != nil {
		t.Fatalf("eliminar la historia votada: %v", err)
	}
	if _, err := r.Obtener(ctx, tienda, v.ID); !errors.Is(err, dominio.ErrVotacionNoExiste) {
		t.Errorf("votación de una historia eliminada: error = %v, se esperaba %v", err, dominio.ErrVotacionNoExiste)
	}
}
