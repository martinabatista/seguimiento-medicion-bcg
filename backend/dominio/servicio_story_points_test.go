package dominio_test

import (
	"context"
	"errors"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// spLeidos devuelve los Story Points de la historia como texto, para los mensajes de los tests.
func spLeidos(h dominio.Historia) any {
	if h.StoryPoints == nil {
		return "sin estimar"
	}
	return *h.StoryPoints
}

func tieneSP(h dominio.Historia, sp int) bool {
	return h.StoryPoints != nil && *h.StoryPoints == sp
}

// spec 07 CA-1, CA-2, RN-3, RN-5: asignar, cambiar y repetir el valor; el resto de la historia no cambia.
func TestEstimar(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()
	h, _ := b.servicio.Crear(ctx, b.tienda, datos("Registrar usuario", 1))

	estimada, err := b.servicio.Estimar(ctx, b.tienda, h.ID, 5)
	if err != nil {
		t.Fatalf("Estimar: %v", err)
	}
	visto, _ := b.servicio.Obtener(ctx, b.tienda, h.ID)
	if !tieneSP(estimada, 5) || !tieneSP(visto, 5) {
		t.Fatalf("Story Points = %v / %v, se esperaba 5", spLeidos(estimada), spLeidos(visto))
	}
	// RN-5: estimar no cambia ningún otro dato.
	sinSP := visto
	sinSP.StoryPoints = nil
	if sinSP != h {
		t.Errorf("historia = %+v, se esperaba igual a %+v salvo los Story Points", visto, h)
	}

	for _, sp := range []int{13, 13, 1} { // cambiar, repetir el mismo valor y volver al menor
		if _, err := b.servicio.Estimar(ctx, b.tienda, h.ID, sp); err != nil {
			t.Fatalf("Estimar %d: %v", sp, err)
		}
		if visto, _ := b.servicio.Obtener(ctx, b.tienda, h.ID); !tieneSP(visto, sp) {
			t.Errorf("Story Points = %v, se esperaba %d", spLeidos(visto), sp)
		}
	}
}

// spec 07 CA-1: una historia en progreso también se estima.
func TestEstimarEnProgreso(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()
	h, _ := b.historias.Guardar(ctx, dominio.Historia{ProyectoID: b.tienda, Titulo: "Registrar usuario", Descripcion: "d", Prioridad: 1, Estado: dominio.EstadoEnProgreso})

	if _, err := b.servicio.Estimar(ctx, b.tienda, h.ID, 8); err != nil {
		t.Fatalf("Estimar: %v", err)
	}
	if visto, _ := b.servicio.Obtener(ctx, b.tienda, h.ID); !tieneSP(visto, 8) || visto.Estado != dominio.EstadoEnProgreso {
		t.Errorf("historia = %+v, se esperaba En progreso con 8 Story Points", visto)
	}
}

// spec 07 CA-2, RN-3: quitar la estimación la deja sin estimar; si ya lo estaba, no es un error.
func TestQuitarEstimacion(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()
	h, _ := b.servicio.Crear(ctx, b.tienda, datos("Registrar usuario", 1))
	b.servicio.Estimar(ctx, b.tienda, h.ID, 3)

	for i := 0; i < 2; i++ { // la segunda vez ya está sin estimar
		sinEstimar, err := b.servicio.QuitarEstimacion(ctx, b.tienda, h.ID)
		if err != nil {
			t.Fatalf("QuitarEstimacion (vez %d): %v", i+1, err)
		}
		visto, _ := b.servicio.Obtener(ctx, b.tienda, h.ID)
		if sinEstimar.StoryPoints != nil || visto.StoryPoints != nil {
			t.Errorf("Story Points = %v, se esperaba sin estimar", spLeidos(visto))
		}
	}
}

// spec 07 CA-3, RN-2: un valor fuera de la escala se rechaza y la estimación no cambia.
func TestEstimarFueraDeEscalaNoCambia(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()
	h, _ := b.servicio.Crear(ctx, b.tienda, datos("Registrar usuario", 1))
	b.servicio.Estimar(ctx, b.tienda, h.ID, 3)

	for _, sp := range []int{0, 4, 22, -3} {
		if _, err := b.servicio.Estimar(ctx, b.tienda, h.ID, sp); !errors.Is(err, dominio.ErrStoryPointsFueraDeEscala) {
			t.Errorf("Estimar %d: error = %v, se esperaba %v", sp, err, dominio.ErrStoryPointsFueraDeEscala)
		}
	}
	if visto, _ := b.servicio.Obtener(ctx, b.tienda, h.ID); !tieneSP(visto, 3) {
		t.Errorf("Story Points = %v, se esperaba 3 sin cambios", spLeidos(visto))
	}
}

// spec 07 CA-4, RN-4: una historia completada no se estima ni se le quita la estimación.
func TestEstimarCompletada(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()
	sp := 3
	h, _ := b.historias.Guardar(ctx, dominio.Historia{ProyectoID: b.tienda, Titulo: "Registrar usuario", Descripcion: "d", Prioridad: 1, Estado: dominio.EstadoCompletada, StoryPoints: &sp})

	if _, err := b.servicio.Estimar(ctx, b.tienda, h.ID, 5); !errors.Is(err, dominio.ErrHistoriaCompletada) {
		t.Errorf("Estimar: error = %v, se esperaba %v", err, dominio.ErrHistoriaCompletada)
	}
	if _, err := b.servicio.QuitarEstimacion(ctx, b.tienda, h.ID); !errors.Is(err, dominio.ErrHistoriaCompletada) {
		t.Errorf("QuitarEstimacion: error = %v, se esperaba %v", err, dominio.ErrHistoriaCompletada)
	}
	if visto, _ := b.servicio.Obtener(ctx, b.tienda, h.ID); !tieneSP(visto, 3) {
		t.Errorf("Story Points = %v, se esperaba 3 sin cambios", spLeidos(visto))
	}
}

// spec 07 CA-6, RN-6: proyecto o historia que no existe; una historia de otro proyecto no existe.
func TestEstimarInexistente(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()
	h, _ := b.servicio.Crear(ctx, b.tienda, datos("Registrar usuario", 1))

	casos := []struct {
		nombre               string
		proyectoID, historia int64
		err                  error
	}{
		{"proyecto que no existe", 999, h.ID, dominio.ErrProyectoNoExiste},
		{"historia que no existe", b.tienda, 999, dominio.ErrHistoriaNoExiste},
		{"historia de otro proyecto", b.turnos, h.ID, dominio.ErrHistoriaNoExiste},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if _, err := b.servicio.Estimar(ctx, c.proyectoID, c.historia, 3); !errors.Is(err, c.err) {
				t.Errorf("Estimar: error = %v, se esperaba %v", err, c.err)
			}
			if _, err := b.servicio.QuitarEstimacion(ctx, c.proyectoID, c.historia); !errors.Is(err, c.err) {
				t.Errorf("QuitarEstimacion: error = %v, se esperaba %v", err, c.err)
			}
		})
	}
	if visto, _ := b.servicio.Obtener(ctx, b.tienda, h.ID); visto.StoryPoints != nil {
		t.Errorf("Story Points = %v, no se debía estimar", spLeidos(visto))
	}
}
