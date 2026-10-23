package dominio

import "context"

// RepositorioIntegrantes guarda y lee los integrantes de los proyectos. Lo implementan postgres/
// (base real) y memoria/ (tests). Las búsquedas por integrante reciben también el proyecto: un
// integrante de otro proyecto se trata como inexistente (spec 02 RN-8).
type RepositorioIntegrantes interface {
	// Guardar registra un integrante nuevo y lo devuelve con el ID asignado. Si el email ya está
	// en el proyecto devuelve ErrIntegranteYaRegistrado: así la regla RN-6 se cumple aunque dos
	// altas lleguen a la vez.
	Guardar(ctx context.Context, i Integrante) (Integrante, error)
	// Listar devuelve los integrantes del proyecto en el orden en que se agregaron (RN-7).
	Listar(ctx context.Context, proyectoID int64) ([]Integrante, error)
	// Obtener devuelve un integrante del proyecto, o ErrIntegranteNoExiste.
	Obtener(ctx context.Context, proyectoID, id int64) (Integrante, error)
}

// ServicioIntegrantes son los casos de uso de la spec 02. Usa el repositorio de proyectos solo
// para verificar que el proyecto exista (RN-1).
type ServicioIntegrantes struct {
	proyectos   RepositorioProyectos
	integrantes RepositorioIntegrantes
}

func NuevoServicioIntegrantes(proyectos RepositorioProyectos, integrantes RepositorioIntegrantes) *ServicioIntegrantes {
	return &ServicioIntegrantes{proyectos: proyectos, integrantes: integrantes}
}

// Agregar valida los datos y registra al integrante en el proyecto (CA-1). Primero verifica el
// proyecto (RN-1), después los datos (RN-2 a RN-5) y por último el email repetido (RN-6, RN-9).
// Si algo falla, no registra nada.
func (s *ServicioIntegrantes) Agregar(ctx context.Context, proyectoID int64, d DatosIntegrante) (Integrante, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return Integrante{}, err
	}
	nuevo, err := NuevoIntegrante(proyectoID, d)
	if err != nil {
		return Integrante{}, err
	}
	actuales, err := s.integrantes.Listar(ctx, proyectoID)
	if err != nil {
		return Integrante{}, err
	}
	for _, i := range actuales {
		if i.Email == nuevo.Email { // los dos están normalizados (RN-3)
			return Integrante{}, ErrIntegranteYaRegistrado
		}
	}
	return s.integrantes.Guardar(ctx, nuevo)
}

// Listar devuelve los integrantes del proyecto en el orden en que se agregaron (CA-2, RN-7).
func (s *ServicioIntegrantes) Listar(ctx context.Context, proyectoID int64) ([]Integrante, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return nil, err
	}
	return s.integrantes.Listar(ctx, proyectoID)
}

// Obtener devuelve un integrante del proyecto (CA-2) o ErrIntegranteNoExiste (CA-6, RN-8).
func (s *ServicioIntegrantes) Obtener(ctx context.Context, proyectoID, id int64) (Integrante, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return Integrante{}, err
	}
	return s.integrantes.Obtener(ctx, proyectoID, id)
}
