package dominio_test

import (
	"context"
	"errors"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// horasLeidas devuelve las horas estimadas de la historia en centésimas, o "sin horas estimadas".
func horasLeidas(h dominio.Historia) any {
	if h.HorasEstimadas == nil {
		return "sin horas estimadas"
	}
	return *h.HorasEstimadas
}

func tieneHoras(h dominio.Historia, centesimas int) bool {
	return h.HorasEstimadas != nil && *h.HorasEstimadas == centesimas
}

// spec 12 RN-2: mayor a 0, hasta 2 decimales y sin máximo.
func TestParsearHorasEstimadas(t *testing.T) {
	casos := []struct {
		texto      string
		centesimas int
		err        error
	}{
		{"8", 800, nil},
		{"6.5", 650, nil},
		{"0.01", 1, nil},    // el mínimo con 2 decimales
		{"500", 50000, nil}, // sin máximo
		{"12.25", 1225, nil},
		{"", 0, dominio.ErrHorasEstimadasObligatorias},
		{"0", 0, dominio.ErrHorasEstimadasInvalidas},
		{"0.00", 0, dominio.ErrHorasEstimadasInvalidas},
		{"-3", 0, dominio.ErrHorasEstimadasInvalidas},
		{"1.125", 0, dominio.ErrHorasEstimadasInvalidas}, // 3 decimales
		{"ocho", 0, dominio.ErrHorasEstimadasInvalidas},
		{"99999999999999999999", 0, dominio.ErrHorasEstimadasInvalidas}, // no entra en un entero
	}
	for _, c := range casos {
		centesimas, err := dominio.ParsearHorasEstimadas(c.texto)
		if !errors.Is(err, c.err) || centesimas != c.centesimas {
			t.Errorf("ParsearHorasEstimadas(%q) = %d, %v; se esperaba %d, %v", c.texto, centesimas, err, c.centesimas, c.err)
		}
	}
}

// spec 12 CA-1, RN-1, RN-4: cargar, cambiar y quitar las horas estimadas; el resto de la historia no cambia.
func TestEstimarHoras(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()
	h, _ := b.servicio.Crear(ctx, b.tienda, datos("Registrar usuario", 1))
	b.servicio.Estimar(ctx, b.tienda, h.ID, 5)
	h, _ = b.servicio.Obtener(ctx, b.tienda, h.ID)

	estimada, err := b.servicio.EstimarHoras(ctx, b.tienda, h.ID, "8")
	if err != nil {
		t.Fatalf("EstimarHoras: %v", err)
	}
	visto, _ := b.servicio.Obtener(ctx, b.tienda, h.ID)
	if !tieneHoras(estimada, 800) || !tieneHoras(visto, 800) {
		t.Fatalf("horas estimadas = %v / %v, se esperaban 800 centésimas", horasLeidas(estimada), horasLeidas(visto))
	}
	// RN-4: no cambia ningún otro dato, tampoco los Story Points.
	sinHoras := visto
	sinHoras.HorasEstimadas = nil
	if sinHoras.Titulo != h.Titulo || sinHoras.Estado != h.Estado || !tieneSP(sinHoras, 5) {
		t.Errorf("historia = %+v, se esperaba igual a %+v salvo las horas estimadas", visto, h)
	}

	if _, err := b.servicio.EstimarHoras(ctx, b.tienda, h.ID, "6.5"); err != nil {
		t.Fatalf("cambiar las horas estimadas: %v", err)
	}
	if visto, _ := b.servicio.Obtener(ctx, b.tienda, h.ID); !tieneHoras(visto, 650) {
		t.Errorf("horas estimadas = %v, se esperaban 650 centésimas", horasLeidas(visto))
	}

	sinEstimar, err := b.servicio.QuitarHorasEstimadas(ctx, b.tienda, h.ID)
	if err != nil {
		t.Fatalf("QuitarHorasEstimadas: %v", err)
	}
	if visto, _ := b.servicio.Obtener(ctx, b.tienda, h.ID); sinEstimar.HorasEstimadas != nil || visto.HorasEstimadas != nil {
		t.Errorf("horas estimadas = %v, se esperaba sin horas estimadas", horasLeidas(visto))
	}
}

// spec 12 RN-4: modificar la historia desde el backlog o estimar sus Story Points no cambia sus horas.
func TestHorasEstimadasSeConservan(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()
	h, _ := b.servicio.Crear(ctx, b.tienda, datos("Registrar usuario", 1))
	b.servicio.EstimarHoras(ctx, b.tienda, h.ID, "8")

	b.servicio.Modificar(ctx, b.tienda, h.ID, datos("Registrar usuario v2", 2), "En progreso")
	b.servicio.Estimar(ctx, b.tienda, h.ID, 3)
	if visto, _ := b.servicio.Obtener(ctx, b.tienda, h.ID); !tieneHoras(visto, 800) {
		t.Errorf("horas estimadas = %v, se esperaban 800 centésimas", horasLeidas(visto))
	}
}

// spec 12 CA-1, CA-2, CA-7, RN-3, RN-11: errores; si algo falla, las horas no cambian.
func TestEstimarHorasErrores(t *testing.T) {
	b := nuevoBacklog(t)
	ctx := context.Background()
	h, _ := b.servicio.Crear(ctx, b.tienda, datos("Registrar usuario", 1))
	b.servicio.EstimarHoras(ctx, b.tienda, h.ID, "8")
	completada, _ := b.historias.Guardar(ctx, dominio.Historia{ProyectoID: b.tienda, Titulo: "Ver catálogo", Descripcion: "d", Prioridad: 2, Estado: dominio.EstadoCompletada})

	casos := []struct {
		nombre     string
		proyectoID int64
		historiaID int64
		texto      string
		err        error
	}{
		{"faltan las horas (CA-1)", b.tienda, h.ID, "", dominio.ErrHorasEstimadasObligatorias},
		{"horas inválidas (CA-1)", b.tienda, h.ID, "0", dominio.ErrHorasEstimadasInvalidas},
		{"historia completada (CA-2, RN-3)", b.tienda, completada.ID, "8", dominio.ErrHistoriaCompletada},
		{"proyecto inexistente (CA-7)", 999, h.ID, "8", dominio.ErrProyectoNoExiste},
		{"historia inexistente (CA-7)", b.tienda, 999, "8", dominio.ErrHistoriaNoExiste},
		{"historia de otro proyecto (RN-11)", b.turnos, h.ID, "8", dominio.ErrHistoriaNoExiste},
	}
	for _, c := range casos {
		if _, err := b.servicio.EstimarHoras(ctx, c.proyectoID, c.historiaID, c.texto); !errors.Is(err, c.err) {
			t.Errorf("%s: err = %v, se esperaba %v", c.nombre, err, c.err)
		}
	}
	if visto, _ := b.servicio.Obtener(ctx, b.tienda, h.ID); !tieneHoras(visto, 800) {
		t.Errorf("horas estimadas = %v, se esperaban 800 centésimas sin cambios", horasLeidas(visto))
	}

	// RN-3: tampoco se quitan las de una historia completada.
	if _, err := b.servicio.QuitarHorasEstimadas(ctx, b.tienda, completada.ID); !errors.Is(err, dominio.ErrHistoriaCompletada) {
		t.Errorf("quitar en una completada: err = %v, se esperaba %v", err, dominio.ErrHistoriaCompletada)
	}
	if _, err := b.servicio.QuitarHorasEstimadas(ctx, b.tienda, 999); !errors.Is(err, dominio.ErrHistoriaNoExiste) {
		t.Errorf("quitar en una inexistente: err = %v, se esperaba %v", err, dominio.ErrHistoriaNoExiste)
	}
}
