package dominio

import (
	"context"
	"time"
)

// RepositorioEsfuerzo guarda y lee los registros de esfuerzo. Lo implementan postgres/ (base real)
// y memoria/ (tests). Las búsquedas reciben el proyecto y la historia: un registro de otra
// historia o de otro proyecto se trata como inexistente (spec 09 RN-8).
type RepositorioEsfuerzo interface {
	// Guardar registra un esfuerzo nuevo y lo devuelve con el ID asignado.
	Guardar(ctx context.Context, r RegistroEsfuerzo) (RegistroEsfuerzo, error)
	// Listar devuelve los registros de la historia (en cualquier orden: lo ordena el dominio).
	Listar(ctx context.Context, proyectoID, historiaID int64) ([]RegistroEsfuerzo, error)
	// Eliminar borra el registro, o devuelve ErrRegistroEsfuerzoNoExiste.
	Eliminar(ctx context.Context, proyectoID, historiaID, id int64) error
}

// ServicioEsfuerzo son los casos de uso de la spec 09. Lee proyectos, historias e integrantes para
// verificar que existan, y la fecha de hoy de un reloj que recibe (como la spec 03).
type ServicioEsfuerzo struct {
	proyectos   RepositorioProyectos
	historias   RepositorioHistorias
	integrantes RepositorioIntegrantes
	esfuerzo    RepositorioEsfuerzo
	hoy         func() time.Time
}

func NuevoServicioEsfuerzo(proyectos RepositorioProyectos, historias RepositorioHistorias, integrantes RepositorioIntegrantes, esfuerzo RepositorioEsfuerzo, hoy func() time.Time) *ServicioEsfuerzo {
	return &ServicioEsfuerzo{proyectos: proyectos, historias: historias, integrantes: integrantes, esfuerzo: esfuerzo, hoy: hoy}
}

// Registrar guarda el esfuerzo sobre la historia (CA-1, CA-2). Primero verifica el proyecto y la
// historia (RN-1, RN-9), después los datos (RN-3 a RN-5) y por último que el integrante sea del
// proyecto (RN-2). Si algo falla, no registra nada.
func (s *ServicioEsfuerzo) Registrar(ctx context.Context, proyectoID, historiaID int64, d DatosEsfuerzo) (RegistroEsfuerzo, error) {
	p, err := s.proyectoEHistoria(ctx, proyectoID, historiaID)
	if err != nil {
		return RegistroEsfuerzo{}, err
	}
	r, err := NuevoRegistroEsfuerzo(p, historiaID, d, s.hoy())
	if err != nil {
		return RegistroEsfuerzo{}, err
	}
	integrante, err := s.integrantes.Obtener(ctx, proyectoID, d.IntegranteID)
	if err != nil {
		return RegistroEsfuerzo{}, err
	}
	guardado, err := s.esfuerzo.Guardar(ctx, r)
	if err != nil {
		return RegistroEsfuerzo{}, err
	}
	guardado.IntegranteNombre = integrante.Nombre
	return guardado, nil
}

// Consultar devuelve el esfuerzo de la historia ordenado por fecha, con el nombre de cada
// integrante y el total de horas (CA-4, RN-6, RN-7).
func (s *ServicioEsfuerzo) Consultar(ctx context.Context, proyectoID, historiaID int64) (EsfuerzoDeHistoria, error) {
	if _, err := s.proyectoEHistoria(ctx, proyectoID, historiaID); err != nil {
		return EsfuerzoDeHistoria{}, err
	}
	registros, err := s.esfuerzo.Listar(ctx, proyectoID, historiaID)
	if err != nil {
		return EsfuerzoDeHistoria{}, err
	}
	integrantes, err := s.integrantes.Listar(ctx, proyectoID)
	if err != nil {
		return EsfuerzoDeHistoria{}, err
	}
	nombres := make(map[int64]string, len(integrantes))
	for _, i := range integrantes {
		nombres[i.ID] = i.Nombre
	}
	for i := range registros {
		registros[i].IntegranteNombre = nombres[registros[i].IntegranteID]
	}
	return ResumirEsfuerzo(registros), nil
}

// Eliminar borra un registro de la historia (CA-6). Si no existe o es de otra historia, devuelve
// ErrRegistroEsfuerzoNoExiste (RN-8).
func (s *ServicioEsfuerzo) Eliminar(ctx context.Context, proyectoID, historiaID, id int64) error {
	if _, err := s.proyectoEHistoria(ctx, proyectoID, historiaID); err != nil {
		return err
	}
	return s.esfuerzo.Eliminar(ctx, proyectoID, historiaID, id)
}

// proyectoEHistoria devuelve el proyecto, o ErrProyectoNoExiste / ErrHistoriaNoExiste (RN-1).
func (s *ServicioEsfuerzo) proyectoEHistoria(ctx context.Context, proyectoID, historiaID int64) (Proyecto, error) {
	p, err := s.proyectos.Obtener(ctx, proyectoID)
	if err != nil {
		return Proyecto{}, err
	}
	if _, err := s.historias.Obtener(ctx, proyectoID, historiaID); err != nil {
		return Proyecto{}, err
	}
	return p, nil
}
