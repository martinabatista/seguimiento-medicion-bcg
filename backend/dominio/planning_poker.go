package dominio

import (
	"cmp"
	"slices"
)

// EstadoVotacion indica si una votación de Planning Poker admite cambios (spec 08 RN-10).
type EstadoVotacion string

const (
	VotacionAbierta EstadoVotacion = "Abierta"
	VotacionCerrada EstadoVotacion = "Cerrada"
)

// Voto es la carta que eligió un Developer en una ronda (spec 08 RN-3, RN-4).
type Voto struct {
	IntegranteID int64
	Valor        int
}

// Ronda es una vuelta de votación. Mientras no está revelada, sus valores no se muestran (RN-5).
type Ronda struct {
	Numero   int
	Revelada bool
	Votos    []Voto // un voto por Developer (RN-4)
}

// Votacion es una sesión de Planning Poker para una historia del backlog (spec 08). La última
// ronda es la actual; las anteriores se conservan con sus votos (RN-8).
type Votacion struct {
	ID         int64
	ProyectoID int64
	HistoriaID int64
	Estado     EstadoVotacion
	Rondas     []Ronda
	Acordada   *int // la estimación acordada al cerrar (RN-9); nil mientras está abierta
}

// NuevaVotacion abre una votación para la historia, en la ronda 1 sin votos (spec 08 CA-1). La
// historia no puede estar completada ni tener otra votación abierta (RN-1), y el proyecto tiene
// que tener al menos un Developer (RN-2). El ID queda en cero: lo asigna el sistema al guardarla.
func NuevaVotacion(h Historia, integrantes []Integrante, existentes []Votacion) (Votacion, error) {
	if h.Estado == EstadoCompletada {
		return Votacion{}, ErrHistoriaCompletada
	}
	for _, v := range existentes {
		if v.HistoriaID == h.ID && v.Estado == VotacionAbierta {
			return Votacion{}, ErrVotacionAbierta
		}
	}
	if !slices.ContainsFunc(integrantes, esDeveloper) {
		return Votacion{}, ErrSinDevelopers
	}
	return Votacion{
		ProyectoID: h.ProyectoID,
		HistoriaID: h.ID,
		Estado:     VotacionAbierta,
		Rondas:     []Ronda{{Numero: 1}},
	}, nil
}

// RondaActual devuelve la última ronda de la votación.
func (v Votacion) RondaActual() Ronda {
	return v.Rondas[len(v.Rondas)-1]
}

// Votar registra el voto de un Developer en la ronda actual (spec 08 CA-2). Primero se valida el
// valor (RN-3, dato inválido) y después los conflictos: votación cerrada (RN-10), ronda revelada
// (RN-6) y rol (RN-2). Si el Developer ya había votado, el voto nuevo reemplaza al anterior (RN-4).
func (v *Votacion) Votar(i Integrante, valor int) error {
	if !EnEscala(valor) {
		return ErrStoryPointsFueraDeEscala
	}
	if err := v.rondaAbierta(); err != nil {
		return err
	}
	if i.ProyectoID != v.ProyectoID { // RN-11: un integrante de otro proyecto no existe acá
		return ErrIntegranteNoExiste
	}
	if !esDeveloper(i) {
		return ErrSoloDevelopers
	}
	r := &v.Rondas[len(v.Rondas)-1]
	if k := slices.IndexFunc(r.Votos, func(vt Voto) bool { return vt.IntegranteID == i.ID }); k >= 0 {
		r.Votos[k].Valor = valor
		return nil
	}
	r.Votos = append(r.Votos, Voto{IntegranteID: i.ID, Valor: valor})
	return nil
}

// Revelar muestra los votos de la ronda actual, solo si votaron todos los Developers del proyecto
// (spec 08 CA-4, RN-6). Un Developer que se agregó durante la ronda también tiene que votar.
func (v *Votacion) Revelar(integrantes []Integrante) error {
	if err := v.rondaAbierta(); err != nil {
		return err
	}
	r := &v.Rondas[len(v.Rondas)-1]
	for _, i := range integrantes {
		if esDeveloper(i) && !voto(*r, i.ID) {
			return ErrFaltanVotos
		}
	}
	r.Revelada = true
	return nil
}

// NuevaRonda empieza la ronda siguiente sin votos, después de revelar la actual (spec 08 CA-6, RN-8).
func (v *Votacion) NuevaRonda() error {
	if err := v.rondaRevelada(); err != nil {
		return err
	}
	v.Rondas = append(v.Rondas, Ronda{Numero: v.RondaActual().Numero + 1})
	return nil
}

// Acordar registra la estimación acordada y cierra la votación (spec 08 CA-7, RN-9). Puede ser
// cualquier valor de la escala, aunque nadie lo haya votado. Guardarla en la historia (spec 07) lo
// hace el servicio.
func (v *Votacion) Acordar(valor int) error {
	if !EnEscala(valor) {
		return ErrStoryPointsFueraDeEscala
	}
	if err := v.rondaRevelada(); err != nil {
		return err
	}
	v.Acordada = &valor
	v.Estado = VotacionCerrada
	return nil
}

// rondaAbierta controla que se pueda votar o revelar: votación abierta y ronda actual sin revelar.
func (v Votacion) rondaAbierta() error {
	if v.Estado == VotacionCerrada {
		return ErrVotacionCerrada
	}
	if v.RondaActual().Revelada {
		return ErrRondaRevelada
	}
	return nil
}

// rondaRevelada controla que se pueda pasar de ronda o acordar: votación abierta y ronda revelada.
func (v Votacion) rondaRevelada() error {
	if v.Estado == VotacionCerrada {
		return ErrVotacionCerrada
	}
	if !v.RondaActual().Revelada {
		return ErrRondaSinRevelar
	}
	return nil
}

func esDeveloper(i Integrante) bool {
	return i.Rol == RolDeveloper
}

func voto(r Ronda, integranteID int64) bool {
	return slices.ContainsFunc(r.Votos, func(vt Voto) bool { return vt.IntegranteID == integranteID })
}

// ---------- lo que se muestra ----------

// VistaVotacion es la votación tal como se puede mostrar (spec 08 RN-5): de la ronda actual sin
// revelar solo dice quién votó y quién falta, nunca los valores. Los valores aparecen únicamente en
// las rondas reveladas.
type VistaVotacion struct {
	ID         int64
	ProyectoID int64
	HistoriaID int64
	Estado     EstadoVotacion
	Ronda      int // número de la ronda actual
	Acordada   *int
	Votaron    []Integrante // en la ronda actual
	Faltan     []Integrante // Developers que todavía no votaron (vacío si la ronda está revelada)
	Reveladas  []RondaRevelada
}

// RondaRevelada es una ronda con sus votos a la vista y el resultado (spec 08 CA-4, CA-5).
type RondaRevelada struct {
	Numero    int
	Votos     []VotoRevelado // en el orden en que se agregaron los integrantes al proyecto
	Resultado ResultadoRonda
}

type VotoRevelado struct {
	Integrante Integrante
	Valor      int
}

// ResultadoRonda dice si hay diferencias entre los votos y cuáles son los extremos (spec 08 RN-7).
type ResultadoRonda struct {
	HayDiferencias bool
	Minimo, Maximo ExtremoRonda
}

// ExtremoRonda es el voto más bajo o el más alto, con todos los integrantes que lo dieron.
type ExtremoRonda struct {
	Valor       int
	Integrantes []Integrante
}

// Vista arma lo que se puede mostrar de la votación. integrantes son los del proyecto, en el
// orden en que se agregaron (spec 02 RN-7): define el orden de votaron, faltan y los votos.
func (v Votacion) Vista(integrantes []Integrante) VistaVotacion {
	actual := v.RondaActual()
	vista := VistaVotacion{
		ID:         v.ID,
		ProyectoID: v.ProyectoID,
		HistoriaID: v.HistoriaID,
		Estado:     v.Estado,
		Ronda:      actual.Numero,
		Acordada:   v.Acordada,
	}
	for _, i := range integrantes {
		switch {
		case voto(actual, i.ID):
			vista.Votaron = append(vista.Votaron, i)
		case esDeveloper(i) && !actual.Revelada:
			vista.Faltan = append(vista.Faltan, i)
		}
	}
	for _, r := range v.Rondas {
		if r.Revelada {
			vista.Reveladas = append(vista.Reveladas, revelar(r, integrantes))
		}
	}
	return vista
}

// revelar arma los votos de una ronda revelada y calcula sus diferencias (RN-7).
func revelar(r Ronda, integrantes []Integrante) RondaRevelada {
	rev := RondaRevelada{Numero: r.Numero}
	for _, i := range integrantes {
		k := slices.IndexFunc(r.Votos, func(vt Voto) bool { return vt.IntegranteID == i.ID })
		if k >= 0 {
			rev.Votos = append(rev.Votos, VotoRevelado{Integrante: i, Valor: r.Votos[k].Valor})
		}
	}
	if len(rev.Votos) == 0 {
		return rev
	}
	minimo, maximo := rev.Votos[0].Valor, rev.Votos[0].Valor
	for _, vt := range rev.Votos {
		minimo, maximo = min(minimo, vt.Valor), max(maximo, vt.Valor)
	}
	rev.Resultado = ResultadoRonda{
		HayDiferencias: minimo != maximo,
		Minimo:         ExtremoRonda{Valor: minimo, Integrantes: quienesVotaron(rev.Votos, minimo)},
		Maximo:         ExtremoRonda{Valor: maximo, Integrantes: quienesVotaron(rev.Votos, maximo)},
	}
	return rev
}

func quienesVotaron(votos []VotoRevelado, valor int) []Integrante {
	var is []Integrante
	for _, vt := range votos {
		if vt.Valor == valor {
			is = append(is, vt.Integrante)
		}
	}
	return is
}

// VotacionesAbiertas devuelve las votaciones abiertas en el orden en que se abrieron (ID creciente),
// para retomarlas sin saber su identificador (spec 08 RN-12). Las cerradas no aparecen.
func VotacionesAbiertas(votaciones []Votacion) []Votacion {
	abiertas := []Votacion{}
	for _, v := range votaciones {
		if v.Estado == VotacionAbierta {
			abiertas = append(abiertas, v)
		}
	}
	slices.SortFunc(abiertas, func(a, b Votacion) int { return cmp.Compare(a.ID, b.ID) })
	return abiertas
}
