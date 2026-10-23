package dominio

import (
	"errors"
	"strings"
	"testing"
)

// emailDeLargo arma un email válido de exactamente n caracteres (para los límites de RN-4).
func emailDeLargo(n int) string {
	const dominio = "@mail.com"
	return strings.Repeat("a", n-len(dominio)) + dominio
}

// spec 02 RN-2 a RN-5 (CA-1, CA-4, CA-5)
func TestNuevoIntegrante(t *testing.T) {
	valido := DatosIntegrante{Nombre: "Ana Pérez", Email: "ana@mail.com", Rol: "Developer"}
	con := func(cambiar func(*DatosIntegrante)) DatosIntegrante {
		d := valido
		cambiar(&d)
		return d
	}
	casos := []struct {
		nombre string
		datos  DatosIntegrante
		err    error
	}{
		{"datos válidos (CA-1)", valido, nil},
		{"nombre vacío (spec 02 RN-2, CA-4)", con(func(d *DatosIntegrante) { d.Nombre = "" }), ErrNombreObligatorio},
		{"nombre con solo espacios (spec 02 RN-2, CA-4)", con(func(d *DatosIntegrante) { d.Nombre = "   " }), ErrNombreObligatorio},
		{"nombre de 100 caracteres (spec 02 RN-2, CA-4)", con(func(d *DatosIntegrante) { d.Nombre = strings.Repeat("n", 100) }), nil},
		{"nombre de 101 caracteres (spec 02 RN-2, CA-4)", con(func(d *DatosIntegrante) { d.Nombre = strings.Repeat("n", 101) }), ErrNombreIntegranteMuyLargo},
		{"nombre de 100 caracteres con eñes (spec 02 RN-2, CA-4)", con(func(d *DatosIntegrante) { d.Nombre = strings.Repeat("ñ", 100) }), nil},
		{"email vacío (spec 02 RN-3, CA-4)", con(func(d *DatosIntegrante) { d.Email = "" }), ErrEmailObligatorio},
		{"email con solo espacios (spec 02 RN-3, CA-4)", con(func(d *DatosIntegrante) { d.Email = "   " }), ErrEmailObligatorio},
		{"email con subdominios (spec 02 RN-4, CA-4)", con(func(d *DatosIntegrante) { d.Email = "ana.p@mail.com.ar" }), nil},
		{"email sin @ (spec 02 RN-4, CA-4)", con(func(d *DatosIntegrante) { d.Email = "ana.mail.com" }), ErrEmailInvalido},
		{"email con dos @ (spec 02 RN-4, CA-4)", con(func(d *DatosIntegrante) { d.Email = "ana@@mail.com" }), ErrEmailInvalido},
		{"email sin texto antes de la @ (spec 02 RN-4, CA-4)", con(func(d *DatosIntegrante) { d.Email = "@mail.com" }), ErrEmailInvalido},
		{"email sin punto en el dominio (spec 02 RN-4, CA-4)", con(func(d *DatosIntegrante) { d.Email = "ana@mail" }), ErrEmailInvalido},
		{"email con el punto al principio del dominio (spec 02 RN-4, CA-4)", con(func(d *DatosIntegrante) { d.Email = "ana@.com" }), ErrEmailInvalido},
		{"email con el punto al final del dominio (spec 02 RN-4, CA-4)", con(func(d *DatosIntegrante) { d.Email = "ana@mail." }), ErrEmailInvalido},
		{"email con un espacio en el medio (spec 02 RN-4, CA-4)", con(func(d *DatosIntegrante) { d.Email = "ana pérez@mail.com" }), ErrEmailInvalido},
		{"email de 254 caracteres (spec 02 RN-4, CA-4)", con(func(d *DatosIntegrante) { d.Email = emailDeLargo(254) }), nil},
		{"email de 255 caracteres (spec 02 RN-4, CA-4)", con(func(d *DatosIntegrante) { d.Email = emailDeLargo(255) }), ErrEmailMuyLargo},
		{"rol Product Owner (spec 02 RN-5, CA-5)", con(func(d *DatosIntegrante) { d.Rol = "Product Owner" }), nil},
		{"rol Scrum Master (spec 02 RN-5, CA-5)", con(func(d *DatosIntegrante) { d.Rol = "Scrum Master" }), nil},
		{"rol vacío (spec 02 RN-5, CA-5)", con(func(d *DatosIntegrante) { d.Rol = "" }), ErrRolObligatorio},
		{"rol que no existe (spec 02 RN-5, CA-5)", con(func(d *DatosIntegrante) { d.Rol = "Tester" }), ErrRolInvalido},
		{"rol en minúscula (spec 02 RN-5, CA-5)", con(func(d *DatosIntegrante) { d.Rol = "developer" }), ErrRolInvalido},
		{"varios datos inválidos: primero el nombre (spec 02 RN-9, CA-4)", DatosIntegrante{Nombre: "", Email: "ana.mail.com", Rol: "Tester"}, ErrNombreObligatorio},
		{"email y rol inválidos: primero el email (spec 02 RN-9, CA-4)", con(func(d *DatosIntegrante) { d.Email, d.Rol = "ana.mail.com", "Tester" }), ErrEmailInvalido},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			i, err := NuevoIntegrante(7, c.datos)
			if !errors.Is(err, c.err) {
				t.Fatalf("error = %v, se esperaba %v", err, c.err)
			}
			if err != nil {
				return
			}
			if i.ProyectoID != 7 || i.Nombre != c.datos.Nombre || string(i.Rol) != c.datos.Rol {
				t.Errorf("integrante = %+v, no conserva los datos recibidos", i)
			}
			if i.ID != 0 {
				t.Errorf("ID = %d, el identificador lo asigna el sistema al guardar (spec 02 RN-8)", i.ID)
			}
		})
	}
}

// spec 02 RN-3 (CA-1, CA-3): el email se guarda sin espacios al borde y en minúsculas.
func TestNuevoIntegranteNormalizaEmail(t *testing.T) {
	casos := []struct{ email, esperado string }{
		{"ana@mail.com", "ana@mail.com"},
		{" Ana@Mail.COM ", "ana@mail.com"},
		{"ANA@MAIL.COM", "ana@mail.com"},
		{"\tana@mail.com\n", "ana@mail.com"},
	}
	for _, c := range casos {
		t.Run(c.email, func(t *testing.T) {
			i, err := NuevoIntegrante(7, DatosIntegrante{Nombre: "Ana Pérez", Email: c.email, Rol: "Developer"})
			if err != nil {
				t.Fatalf("NuevoIntegrante: %v", err)
			}
			if i.Email != c.esperado {
				t.Errorf("email = %q, se esperaba %q", i.Email, c.esperado)
			}
		})
	}
}

// Los errores de validación del integrante son datos inválidos (400); "no existe" (404) y
// "ya registrado" (409) no.
func TestErroresDeIntegrante(t *testing.T) {
	invalidos := []error{
		ErrNombreIntegranteMuyLargo, ErrEmailObligatorio, ErrEmailInvalido, ErrEmailMuyLargo,
		ErrRolObligatorio, ErrRolInvalido,
	}
	for _, e := range invalidos {
		var d ErrDatoInvalido
		if !errors.As(e, &d) {
			t.Errorf("%q debería ser un ErrDatoInvalido", e)
		}
	}
	for _, e := range []error{ErrIntegranteNoExiste, ErrIntegranteYaRegistrado} {
		var d ErrDatoInvalido
		if errors.As(e, &d) {
			t.Errorf("%q no es un dato inválido", e)
		}
	}
}
