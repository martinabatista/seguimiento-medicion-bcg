package dominio_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

func defecto(descripcion, severidad string, historiaID, sprintID int64) dominio.DatosDefecto {
	return dominio.DatosDefecto{Descripcion: descripcion, Severidad: severidad, HistoriaID: historiaID, SprintDeteccionID: sprintID}
}

// spec 11 CA-1, RN-6: un defecto nuevo arranca Abierto y sin Sprint de resolución.
func TestNuevoDefecto(t *testing.T) {
	d, err := dominio.NuevoDefecto(7, defecto("El alta acepta emails sin @", "Alta", 3, 5))
	if err != nil {
		t.Fatalf("NuevoDefecto: %v", err)
	}
	if d.ID != 0 || d.ProyectoID != 7 || d.Descripcion != "El alta acepta emails sin @" || d.Severidad != dominio.SeveridadAlta ||
		d.HistoriaID != 3 || d.SprintDeteccionID != 5 || d.Estado != dominio.DefectoAbierto || d.SprintResolucionID != 0 {
		t.Errorf("defecto = %+v, no conserva los datos o no arranca Abierto y sin resolución", d)
	}
}

// spec 11 CA-5, RN-2 a RN-5, RN-13: validaciones, informando el primer dato inválido en orden.
func TestNuevoDefectoValidaciones(t *testing.T) {
	casos := []struct {
		nombre string
		datos  dominio.DatosDefecto
		err    error
	}{
		{"cada severidad: Baja", defecto("Falla", "Baja", 1, 1), nil},
		{"cada severidad: Media", defecto("Falla", "Media", 1, 1), nil},
		{"cada severidad: Crítica", defecto("Falla", "Crítica", 1, 1), nil},
		{"descripción de 1000", defecto(strings.Repeat("ñ", 1000), "Baja", 1, 1), nil},
		{"descripción vacía", defecto("", "Alta", 1, 1), dominio.ErrDescripcionObligatoria},
		{"descripción con solo espacios", defecto("   ", "Alta", 1, 1), dominio.ErrDescripcionObligatoria},
		{"descripción de 1001", defecto(strings.Repeat("ñ", 1001), "Baja", 1, 1), dominio.ErrDescripcionHistoriaMuyLarga},
		{"sin severidad", defecto("Falla", "", 1, 1), dominio.ErrSeveridadObligatoria},
		{"severidad en minúscula", defecto("Falla", "alta", 1, 1), dominio.ErrSeveridadInvalida},
		{"severidad sin tilde", defecto("Falla", "Critica", 1, 1), dominio.ErrSeveridadInvalida},
		{"sin historia", defecto("Falla", "Alta", 0, 1), dominio.ErrHistoriaDefectoObligatoria},
		{"sin Sprint de detección", defecto("Falla", "Alta", 1, 0), dominio.ErrSprintDeteccionObligatorio},
		{"todo mal: primero la descripción", defecto("", "x", 0, 0), dominio.ErrDescripcionObligatoria},
		{"severidad, historia y Sprint mal: primero la severidad", defecto("Falla", "x", 0, 0), dominio.ErrSeveridadInvalida},
		{"historia y Sprint mal: primero la historia", defecto("Falla", "Alta", 0, 0), dominio.ErrHistoriaDefectoObligatoria},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if _, err := dominio.NuevoDefecto(1, c.datos); !errors.Is(err, c.err) {
				t.Errorf("error = %v, se esperaba %v", err, c.err)
			}
		})
	}
}

// spec 11 RN-7: solo los tres estados, escritos exactamente así.
func TestParsearEstadoDefecto(t *testing.T) {
	validos := map[string]dominio.EstadoDefecto{
		"Abierto":       dominio.DefectoAbierto,
		"En corrección": dominio.DefectoEnCorreccion,
		"Resuelto":      dominio.DefectoResuelto,
	}
	for texto, esperado := range validos {
		if got, err := dominio.ParsearEstadoDefecto(texto); err != nil || got != esperado {
			t.Errorf("ParsearEstadoDefecto(%q) = %q, %v; se esperaba %q", texto, got, err, esperado)
		}
	}
	for _, texto := range []string{"", "Cerrado", "resuelto", "En correccion"} {
		if _, err := dominio.ParsearEstadoDefecto(texto); !errors.Is(err, dominio.ErrEstadoInvalido) {
			t.Errorf("ParsearEstadoDefecto(%q): error = %v, se esperaba %v", texto, err, dominio.ErrEstadoInvalido)
		}
	}
}

// spec 11 CA-2, CA-3, RN-7 a RN-9: cambiar el estado y registrar o borrar el Sprint de resolución.
func TestCambiarEstadoDefecto(t *testing.T) {
	s1 := dominio.Sprint{ID: 10, FechaInicio: f("2026-11-01"), FechaFin: f("2026-11-14")}
	s2 := dominio.Sprint{ID: 20, FechaInicio: f("2026-11-16"), FechaFin: f("2026-11-28")}
	abierto := dominio.Defecto{ID: 1, SprintDeteccionID: s1.ID, Estado: dominio.DefectoAbierto}
	resuelto := dominio.Defecto{ID: 1, SprintDeteccionID: s1.ID, Estado: dominio.DefectoResuelto, SprintResolucionID: s2.ID}

	casos := []struct {
		nombre       string
		defecto      dominio.Defecto
		estado       dominio.EstadoDefecto
		deteccion    dominio.Sprint
		resolucion   *dominio.Sprint
		esperado     dominio.EstadoDefecto
		resolucionID int64
		err          error
	}{
		{"resolver en un Sprint posterior", abierto, dominio.DefectoResuelto, s1, &s2, dominio.DefectoResuelto, s2.ID, nil},
		{"resolver en el mismo Sprint", abierto, dominio.DefectoResuelto, s1, &s1, dominio.DefectoResuelto, s1.ID, nil},
		{"volver a resolver con otro Sprint", resuelto, dominio.DefectoResuelto, s1, &s1, dominio.DefectoResuelto, s1.ID, nil},
		{"pasar a En corrección", abierto, dominio.DefectoEnCorreccion, s1, nil, dominio.DefectoEnCorreccion, 0, nil},
		{"reabrir borra la resolución", resuelto, dominio.DefectoAbierto, s1, nil, dominio.DefectoAbierto, 0, nil},
		{"En corrección ignora el Sprint indicado", resuelto, dominio.DefectoEnCorreccion, s1, &s2, dominio.DefectoEnCorreccion, 0, nil},
		{"resolver sin Sprint", abierto, dominio.DefectoResuelto, s1, nil, "", 0, dominio.ErrSprintResolucionObligatorio},
		{"resolver antes de la detección", dominio.Defecto{ID: 1, SprintDeteccionID: s2.ID, Estado: dominio.DefectoAbierto}, dominio.DefectoResuelto, s2, &s1, "", 0, dominio.ErrResolucionAnteriorADeteccion},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got, err := dominio.CambiarEstadoDefecto(c.defecto, c.estado, c.deteccion, c.resolucion)
			if !errors.Is(err, c.err) {
				t.Fatalf("error = %v, se esperaba %v", err, c.err)
			}
			if err != nil {
				return
			}
			if got.Estado != c.esperado || got.SprintResolucionID != c.resolucionID {
				t.Errorf("defecto = %+v; se esperaba estado %q y Sprint de resolución %d", got, c.esperado, c.resolucionID)
			}
		})
	}
}

// spec 11 CA-4, RN-10, RN-11: detectados por Sprint de detección, resueltos solo los Resueltos por
// Sprint de resolución; todos los Sprints ordenados por fecha de inicio, con 0 si no tienen defectos.
func TestContarDefectos(t *testing.T) {
	// En desorden y con números que no siguen las fechas, para verificar que se ordena por inicio.
	sprints := []dominio.Sprint{
		{ID: 30, Numero: 2, FechaInicio: f("2026-11-30")},
		{ID: 10, Numero: 3, FechaInicio: f("2026-11-01")},
		{ID: 20, Numero: 1, FechaInicio: f("2026-11-16")},
	}
	defectos := []dominio.Defecto{
		{SprintDeteccionID: 10, Estado: dominio.DefectoResuelto, SprintResolucionID: 20},
		{SprintDeteccionID: 10, Estado: dominio.DefectoAbierto},
		{SprintDeteccionID: 20, Estado: dominio.DefectoResuelto, SprintResolucionID: 20},
		{SprintDeteccionID: 20, Estado: dominio.DefectoEnCorreccion}, // reabierto: no cuenta como resuelto
	}

	m := dominio.ContarDefectos(sprints, defectos)

	if m.Detectados != 4 || m.Resueltos != 2 {
		t.Errorf("proyecto: %d detectados y %d resueltos; se esperaban 4 y 2", m.Detectados, m.Resueltos)
	}
	esperados := []dominio.DefectosDeSprint{
		{SprintID: 10, Numero: 3, Detectados: 2, Resueltos: 0},
		{SprintID: 20, Numero: 1, Detectados: 2, Resueltos: 2},
		{SprintID: 30, Numero: 2, Detectados: 0, Resueltos: 0},
	}
	if len(m.Sprints) != len(esperados) {
		t.Fatalf("Sprints = %+v; se esperaban %+v", m.Sprints, esperados)
	}
	for i, e := range esperados {
		if m.Sprints[i] != e {
			t.Errorf("Sprint %d = %+v; se esperaba %+v", i, m.Sprints[i], e)
		}
	}
}

// spec 11 Casos límite: sin Sprints ni defectos, todo en 0 y la lista vacía (no nil, para el JSON).
func TestContarDefectosSinDatos(t *testing.T) {
	m := dominio.ContarDefectos(nil, nil)
	if m.Detectados != 0 || m.Resueltos != 0 || m.Sprints == nil || len(m.Sprints) != 0 {
		t.Errorf("métricas = %+v; se esperaban 0, 0 y una lista vacía", m)
	}
}
