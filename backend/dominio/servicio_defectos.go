package dominio

import "context"

// RepositorioDefectos guarda y lee los defectos. Lo implementan postgres/ (base real) y memoria/
// (tests). Las búsquedas reciben el proyecto: un defecto de otro proyecto se trata como inexistente
// (spec 11 RN-12). Los datos que se completan al consultar (título y números) no se guardan.
type RepositorioDefectos interface {
	// Guardar registra un defecto nuevo y lo devuelve con el ID asignado.
	Guardar(ctx context.Context, d Defecto) (Defecto, error)
	// Actualizar guarda el estado y el Sprint de resolución, o devuelve ErrDefectoNoExiste.
	Actualizar(ctx context.Context, d Defecto) error
	// Listar devuelve los defectos del proyecto en el orden en que se registraron.
	Listar(ctx context.Context, proyectoID int64) ([]Defecto, error)
	// Obtener devuelve el defecto, o ErrDefectoNoExiste.
	Obtener(ctx context.Context, proyectoID, id int64) (Defecto, error)
}

// ServicioDefectos son los casos de uso de la spec 11. Lee proyectos, historias y Sprints para
// verificar que existan y para completar el título de la historia y el número de cada Sprint.
type ServicioDefectos struct {
	proyectos RepositorioProyectos
	historias RepositorioHistorias
	sprints   RepositorioSprints
	defectos  RepositorioDefectos
}

func NuevoServicioDefectos(proyectos RepositorioProyectos, historias RepositorioHistorias, sprints RepositorioSprints, defectos RepositorioDefectos) *ServicioDefectos {
	return &ServicioDefectos{proyectos: proyectos, historias: historias, sprints: sprints, defectos: defectos}
}

// Registrar guarda el defecto (CA-1, CA-5). Verifica el proyecto, después los datos y por último
// que la historia y el Sprint de detección sean del proyecto (RN-1, RN-13). Si algo falla, no
// registra nada.
func (s *ServicioDefectos) Registrar(ctx context.Context, proyectoID int64, datos DatosDefecto) (Defecto, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return Defecto{}, err
	}
	d, err := NuevoDefecto(proyectoID, datos)
	if err != nil {
		return Defecto{}, err
	}
	h, err := s.historias.Obtener(ctx, proyectoID, d.HistoriaID)
	if err != nil {
		return Defecto{}, err
	}
	deteccion, err := s.sprints.Obtener(ctx, proyectoID, d.SprintDeteccionID)
	if err != nil {
		return Defecto{}, err
	}
	guardado, err := s.defectos.Guardar(ctx, d)
	if err != nil {
		return Defecto{}, err
	}
	guardado.HistoriaTitulo = h.Titulo
	guardado.SprintDeteccionNumero = deteccion.Numero
	return guardado, nil
}

// Listar devuelve los defectos del proyecto en el orden en que se registraron, con su historia y
// sus Sprints (CA-5).
func (s *ServicioDefectos) Listar(ctx context.Context, proyectoID int64) ([]Defecto, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return nil, err
	}
	defectos, err := s.defectos.Listar(ctx, proyectoID)
	if err != nil {
		return nil, err
	}
	if err := s.completar(ctx, proyectoID, defectos); err != nil {
		return nil, err
	}
	return defectos, nil
}

// Obtener devuelve un defecto del proyecto con todos sus datos (CA-5, RN-12).
func (s *ServicioDefectos) Obtener(ctx context.Context, proyectoID, id int64) (Defecto, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return Defecto{}, err
	}
	d, err := s.defectos.Obtener(ctx, proyectoID, id)
	if err != nil {
		return Defecto{}, err
	}
	return s.completarUno(ctx, d)
}

// CambiarEstado pasa el defecto al estado indicado (CA-2, CA-3). El Sprint de resolución solo se
// busca al resolver (RN-8); para los otros estados no se tiene en cuenta (RN-9). Si algo falla, el
// defecto no cambia.
func (s *ServicioDefectos) CambiarEstado(ctx context.Context, proyectoID, id int64, estadoTexto string, sprintResolucionID int64) (Defecto, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return Defecto{}, err
	}
	d, err := s.defectos.Obtener(ctx, proyectoID, id)
	if err != nil {
		return Defecto{}, err
	}
	estado, err := ParsearEstadoDefecto(estadoTexto)
	if err != nil {
		return Defecto{}, err
	}
	deteccion, err := s.sprints.Obtener(ctx, proyectoID, d.SprintDeteccionID)
	if err != nil {
		return Defecto{}, err
	}
	var resolucion *Sprint
	if estado == DefectoResuelto && sprintResolucionID != 0 {
		sr, err := s.sprints.Obtener(ctx, proyectoID, sprintResolucionID)
		if err != nil {
			return Defecto{}, err
		}
		resolucion = &sr
	}
	cambiado, err := CambiarEstadoDefecto(d, estado, deteccion, resolucion)
	if err != nil {
		return Defecto{}, err
	}
	if err := s.defectos.Actualizar(ctx, cambiado); err != nil {
		return Defecto{}, err
	}
	return s.completarUno(ctx, cambiado)
}

// Metricas devuelve los defectos detectados y resueltos del proyecto y de cada Sprint (CA-4).
func (s *ServicioDefectos) Metricas(ctx context.Context, proyectoID int64) (MetricasDefectos, error) {
	if err := verificarProyecto(ctx, s.proyectos, proyectoID); err != nil {
		return MetricasDefectos{}, err
	}
	sprints, err := s.sprints.Listar(ctx, proyectoID)
	if err != nil {
		return MetricasDefectos{}, err
	}
	defectos, err := s.defectos.Listar(ctx, proyectoID)
	if err != nil {
		return MetricasDefectos{}, err
	}
	return ContarDefectos(sprints, defectos), nil
}

func (s *ServicioDefectos) completarUno(ctx context.Context, d Defecto) (Defecto, error) {
	defectos := []Defecto{d}
	if err := s.completar(ctx, d.ProyectoID, defectos); err != nil {
		return Defecto{}, err
	}
	return defectos[0], nil
}

// completar agrega a cada defecto el título de su historia y el número de sus Sprints.
func (s *ServicioDefectos) completar(ctx context.Context, proyectoID int64, defectos []Defecto) error {
	historias, err := s.historias.Listar(ctx, proyectoID)
	if err != nil {
		return err
	}
	sprints, err := s.sprints.Listar(ctx, proyectoID)
	if err != nil {
		return err
	}
	titulos := make(map[int64]string, len(historias))
	for _, h := range historias {
		titulos[h.ID] = h.Titulo
	}
	numeros := make(map[int64]int, len(sprints))
	for _, sp := range sprints {
		numeros[sp.ID] = sp.Numero
	}
	for i := range defectos {
		defectos[i].HistoriaTitulo = titulos[defectos[i].HistoriaID]
		defectos[i].SprintDeteccionNumero = numeros[defectos[i].SprintDeteccionID]
		defectos[i].SprintResolucionNumero = numeros[defectos[i].SprintResolucionID] // 0 si no está resuelto
	}
	return nil
}
