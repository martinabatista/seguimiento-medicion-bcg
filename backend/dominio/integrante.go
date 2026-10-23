package dominio

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Largos máximos de los datos de un integrante (spec 02 RN-2, RN-4).
const (
	MaxNombreIntegrante = 100
	MaxEmail            = 254
)

// RolIntegrante es el rol Scrum de un integrante dentro del proyecto (spec 02 RN-5).
type RolIntegrante string

const (
	RolProductOwner RolIntegrante = "Product Owner"
	RolScrumMaster  RolIntegrante = "Scrum Master"
	RolDeveloper    RolIntegrante = "Developer"
)

// DatosIntegrante son los datos que carga la usuaria al agregar un integrante.
type DatosIntegrante struct {
	Nombre string
	Email  string
	Rol    string
}

// Integrante es una persona que trabaja en un proyecto (spec 02).
type Integrante struct {
	ID         int64
	ProyectoID int64
	Nombre     string
	Email      string // normalizado: sin espacios al borde y en minúsculas (RN-3)
	Rol        RolIntegrante
}

// NuevoIntegrante arma un integrante validando las reglas de la spec 02, en el orden de RN-9:
// nombre, email y rol. El ID queda en cero: lo asigna el sistema al guardarlo (RN-8).
// El email repetido (RN-6) no se controla acá: depende de los otros integrantes del proyecto.
func NuevoIntegrante(proyectoID int64, d DatosIntegrante) (Integrante, error) {
	// RN-2: nombre obligatorio, hasta 100 caracteres (no bytes: cuenta bien las tildes y la ñ).
	if strings.TrimSpace(d.Nombre) == "" {
		return Integrante{}, ErrNombreObligatorio
	}
	if utf8.RuneCountInString(d.Nombre) > MaxNombreIntegrante {
		return Integrante{}, ErrNombreIntegranteMuyLargo
	}
	// RN-3: el email se normaliza antes de validarlo y se guarda así.
	email := strings.ToLower(strings.TrimSpace(d.Email))
	if email == "" {
		return Integrante{}, ErrEmailObligatorio
	}
	// RN-4: hasta 254 caracteres y con formato válido.
	if utf8.RuneCountInString(email) > MaxEmail {
		return Integrante{}, ErrEmailMuyLargo
	}
	if !emailValido(email) {
		return Integrante{}, ErrEmailInvalido
	}
	// RN-5: rol obligatorio, uno de los tres.
	if d.Rol == "" {
		return Integrante{}, ErrRolObligatorio
	}
	rol, err := ParsearRol(d.Rol)
	if err != nil {
		return Integrante{}, err
	}
	return Integrante{ProyectoID: proyectoID, Nombre: d.Nombre, Email: email, Rol: rol}, nil
}

// ParsearRol convierte un texto en rol. Se escribe exactamente como en la spec (con mayúsculas):
// cualquier otro texto es inválido (RN-5).
func ParsearRol(texto string) (RolIntegrante, error) {
	switch r := RolIntegrante(texto); r {
	case RolProductOwner, RolScrumMaster, RolDeveloper:
		return r, nil
	}
	return "", ErrRolInvalido
}

// emailValido aplica el formato de RN-4: sin espacios, exactamente una @, algo antes de la @ y,
// después, un dominio con al menos un punto que no esté ni al principio ni al final.
func emailValido(email string) bool {
	if strings.IndexFunc(email, unicode.IsSpace) >= 0 {
		return false
	}
	local, dominio, ok := strings.Cut(email, "@")
	if !ok || local == "" || strings.Contains(dominio, "@") {
		return false
	}
	return strings.Contains(dominio, ".") &&
		!strings.HasPrefix(dominio, ".") && !strings.HasSuffix(dominio, ".")
}
