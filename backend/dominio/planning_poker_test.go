package dominio

import (
	"errors"
	"testing"
)

// Equipo de prueba del proyecto 7, en el orden en que se agregaron: Ana y Beto votan, Caro no.
var (
	ana  = Integrante{ID: 1, ProyectoID: 7, Nombre: "Ana", Rol: RolDeveloper}
	beto = Integrante{ID: 2, ProyectoID: 7, Nombre: "Beto", Rol: RolDeveloper}
	caro = Integrante{ID: 3, ProyectoID: 7, Nombre: "Caro", Rol: RolProductOwner}
	dani = Integrante{ID: 4, ProyectoID: 7, Nombre: "Dani", Rol: RolDeveloper}

	equipo            = []Integrante{ana, beto, caro}
	historiaPendiente = Historia{ID: 10, ProyectoID: 7, Titulo: "Registrar usuario", Prioridad: 1, Estado: EstadoPendiente}
)

// votacionDePrueba abre una votación para la historia pendiente con el equipo de prueba.
func votacionDePrueba(t *testing.T) Votacion {
	t.Helper()
	v, err := NuevaVotacion(historiaPendiente, equipo, nil)
	if err != nil {
		t.Fatalf("NuevaVotacion: %v", err)
	}
	return v
}

// votar registra votos que tienen que ser válidos.
func votar(t *testing.T, v *Votacion, votos map[Integrante]int) {
	t.Helper()
	for i, valor := range votos {
		if err := v.Votar(i, valor); err != nil {
			t.Fatalf("%s vota %d: %v", i.Nombre, valor, err)
		}
	}
}

func nombres(is []Integrante) []string {
	var ns []string
	for _, i := range is {
		ns = append(ns, i.Nombre)
	}
	return ns
}

func mismosNombres(is []Integrante, esperados ...string) bool {
	ns := nombres(is)
	if len(ns) != len(esperados) {
		return false
	}
	for k := range ns {
		if ns[k] != esperados[k] {
			return false
		}
	}
	return true
}

// spec 08 RN-1, RN-2 (CA-1): cuándo se puede abrir una votación.
func TestNuevaVotacion(t *testing.T) {
	completada := historiaPendiente
	completada.Estado = EstadoCompletada
	enProgreso := historiaPendiente
	enProgreso.Estado = EstadoEnProgreso
	abiertaDeLaMisma := Votacion{ID: 1, ProyectoID: 7, HistoriaID: 10, Estado: VotacionAbierta}
	cerradaDeLaMisma := Votacion{ID: 2, ProyectoID: 7, HistoriaID: 10, Estado: VotacionCerrada}
	abiertaDeOtra := Votacion{ID: 3, ProyectoID: 7, HistoriaID: 11, Estado: VotacionAbierta}

	casos := []struct {
		nombre      string
		historia    Historia
		integrantes []Integrante
		existentes  []Votacion
		err         error
	}{
		{"historia pendiente", historiaPendiente, equipo, nil, nil},
		{"historia en progreso", enProgreso, equipo, nil, nil},
		{"un solo Developer", historiaPendiente, []Integrante{ana, caro}, nil, nil},
		{"la misma historia ya tuvo una votación cerrada", historiaPendiente, equipo, []Votacion{cerradaDeLaMisma}, nil},
		{"otra historia tiene una votación abierta", historiaPendiente, equipo, []Votacion{abiertaDeOtra}, nil},
		{"historia completada (RN-1)", completada, equipo, nil, ErrHistoriaCompletada},
		{"la historia ya tiene una votación abierta (RN-1)", historiaPendiente, equipo, []Votacion{abiertaDeLaMisma}, ErrVotacionAbierta},
		{"proyecto sin Developers (RN-2)", historiaPendiente, []Integrante{caro}, nil, ErrSinDevelopers},
		{"proyecto sin integrantes (RN-2)", historiaPendiente, nil, nil, ErrSinDevelopers},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			v, err := NuevaVotacion(c.historia, c.integrantes, c.existentes)
			if !errors.Is(err, c.err) {
				t.Fatalf("error = %v, se esperaba %v", err, c.err)
			}
			if err != nil {
				return
			}
			if v.ID != 0 || v.ProyectoID != 7 || v.HistoriaID != 10 || v.Estado != VotacionAbierta || v.Acordada != nil {
				t.Errorf("votación = %+v, se esperaba abierta, de la historia 10 y sin ID", v)
			}
			if v.RondaActual().Numero != 1 || len(v.RondaActual().Votos) != 0 || v.RondaActual().Revelada {
				t.Errorf("ronda actual = %+v, se esperaba la 1 sin votos ni revelar", v.RondaActual())
			}
		})
	}
}

// spec 08 RN-2 a RN-4 (CA-2): quién vota, con qué valores, y que el voto nuevo reemplaza al anterior.
func TestVotar(t *testing.T) {
	t.Run("un Developer vota", func(t *testing.T) {
		v := votacionDePrueba(t)
		if err := v.Votar(ana, 5); err != nil {
			t.Fatalf("Votar: %v", err)
		}
		if votos := v.RondaActual().Votos; len(votos) != 1 || votos[0] != (Voto{IntegranteID: ana.ID, Valor: 5}) {
			t.Errorf("votos = %+v, se esperaba el de Ana con 5", votos)
		}
	})
	t.Run("los extremos de la escala", func(t *testing.T) {
		v := votacionDePrueba(t)
		for _, valor := range []int{1, 21} {
			if err := v.Votar(ana, valor); err != nil {
				t.Errorf("Votar %d: %v", valor, err)
			}
		}
	})
	t.Run("cambiar el voto antes de revelar: cuenta el último (RN-4)", func(t *testing.T) {
		v := votacionDePrueba(t)
		votar(t, &v, map[Integrante]int{ana: 3})
		votar(t, &v, map[Integrante]int{ana: 8})
		if votos := v.RondaActual().Votos; len(votos) != 1 || votos[0].Valor != 8 {
			t.Errorf("votos = %+v, se esperaba un solo voto de Ana con 8", votos)
		}
	})

	errores := []struct {
		nombre     string
		preparar   func(*Votacion)
		integrante Integrante
		valor      int
		err        error
	}{
		{"valor fuera de la escala (RN-3)", nil, ana, 4, ErrStoryPointsFueraDeEscala},
		{"cero (RN-3)", nil, ana, 0, ErrStoryPointsFueraDeEscala},
		{"un integrante que no es Developer (RN-2)", nil, caro, 5, ErrSoloDevelopers},
		{"la ronda ya fue revelada (RN-6)", func(v *Votacion) {
			votar(t, v, map[Integrante]int{ana: 5, beto: 5})
			v.Revelar(equipo)
		}, ana, 8, ErrRondaRevelada},
		{"la votación está cerrada (RN-10)", func(v *Votacion) {
			votar(t, v, map[Integrante]int{ana: 5, beto: 5})
			v.Revelar(equipo)
			v.Acordar(5)
		}, ana, 8, ErrVotacionCerrada},
		{"primero el dato inválido, después el conflicto", func(v *Votacion) {
			votar(t, v, map[Integrante]int{ana: 5, beto: 5})
			v.Revelar(equipo)
		}, caro, 4, ErrStoryPointsFueraDeEscala},
	}
	for _, c := range errores {
		t.Run(c.nombre, func(t *testing.T) {
			v := votacionDePrueba(t)
			if c.preparar != nil {
				c.preparar(&v)
			}
			antes := len(v.RondaActual().Votos)
			if err := v.Votar(c.integrante, c.valor); !errors.Is(err, c.err) {
				t.Fatalf("error = %v, se esperaba %v", err, c.err)
			}
			if len(v.RondaActual().Votos) != antes {
				t.Errorf("votos = %+v, un voto rechazado no debe registrarse", v.RondaActual().Votos)
			}
		})
	}
}

// spec 08 RN-6 (CA-4): se revela solo cuando votaron todos los Developers.
func TestRevelar(t *testing.T) {
	t.Run("votaron todos", func(t *testing.T) {
		v := votacionDePrueba(t)
		votar(t, &v, map[Integrante]int{ana: 3, beto: 13})
		if err := v.Revelar(equipo); err != nil {
			t.Fatalf("Revelar: %v", err)
		}
		if !v.RondaActual().Revelada {
			t.Error("la ronda debería quedar revelada")
		}
	})
	t.Run("falta un voto", func(t *testing.T) {
		v := votacionDePrueba(t)
		votar(t, &v, map[Integrante]int{ana: 5})
		if err := v.Revelar(equipo); !errors.Is(err, ErrFaltanVotos) {
			t.Fatalf("error = %v, se esperaba %v", err, ErrFaltanVotos)
		}
		if v.RondaActual().Revelada {
			t.Error("la ronda no debería quedar revelada")
		}
	})
	t.Run("un Developer que se agregó después también tiene que votar", func(t *testing.T) {
		v := votacionDePrueba(t)
		votar(t, &v, map[Integrante]int{ana: 5, beto: 5})
		if err := v.Revelar(append(equipo, dani)); !errors.Is(err, ErrFaltanVotos) {
			t.Fatalf("error = %v, se esperaba %v", err, ErrFaltanVotos)
		}
	})
	t.Run("revelar dos veces", func(t *testing.T) {
		v := votacionDePrueba(t)
		votar(t, &v, map[Integrante]int{ana: 5, beto: 5})
		v.Revelar(equipo)
		if err := v.Revelar(equipo); !errors.Is(err, ErrRondaRevelada) {
			t.Fatalf("error = %v, se esperaba %v", err, ErrRondaRevelada)
		}
	})
	t.Run("votación cerrada", func(t *testing.T) {
		v := votacionDePrueba(t)
		votar(t, &v, map[Integrante]int{ana: 5, beto: 5})
		v.Revelar(equipo)
		v.Acordar(5)
		if err := v.Revelar(equipo); !errors.Is(err, ErrVotacionCerrada) {
			t.Fatalf("error = %v, se esperaba %v", err, ErrVotacionCerrada)
		}
	})
}

// spec 08 RN-7 (CA-5): diferencias, voto más bajo y más alto, con quiénes los dieron.
func TestResultadoDeRonda(t *testing.T) {
	casos := []struct {
		nombre         string
		integrantes    []Integrante
		votos          map[Integrante]int
		hayDiferencias bool
		minimo, maximo int
		quienesMin     []string
		quienesMax     []string
	}{
		{"votos distintos", equipo, map[Integrante]int{ana: 3, beto: 13}, true, 3, 13, []string{"Ana"}, []string{"Beto"}},
		{"todos iguales", equipo, map[Integrante]int{ana: 5, beto: 5}, false, 5, 5, []string{"Ana", "Beto"}, []string{"Ana", "Beto"}},
		{"empate en el más bajo", append(equipo, dani), map[Integrante]int{ana: 3, beto: 8, dani: 3}, true, 3, 8, []string{"Ana", "Dani"}, []string{"Beto"}},
		{"un solo Developer", []Integrante{ana, caro}, map[Integrante]int{ana: 8}, false, 8, 8, []string{"Ana"}, []string{"Ana"}},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			v, _ := NuevaVotacion(historiaPendiente, c.integrantes, nil)
			votar(t, &v, c.votos)
			if err := v.Revelar(c.integrantes); err != nil {
				t.Fatalf("Revelar: %v", err)
			}
			vista := v.Vista(c.integrantes)
			if len(vista.Reveladas) != 1 {
				t.Fatalf("rondas reveladas = %d, se esperaba 1", len(vista.Reveladas))
			}
			r := vista.Reveladas[0].Resultado
			if r.HayDiferencias != c.hayDiferencias || r.Minimo.Valor != c.minimo || r.Maximo.Valor != c.maximo ||
				!mismosNombres(r.Minimo.Integrantes, c.quienesMin...) || !mismosNombres(r.Maximo.Integrantes, c.quienesMax...) {
				t.Errorf("resultado = diferencias %v, más bajo %d de %v, más alto %d de %v; se esperaba %v, %d de %v, %d de %v",
					r.HayDiferencias, r.Minimo.Valor, nombres(r.Minimo.Integrantes), r.Maximo.Valor, nombres(r.Maximo.Integrantes),
					c.hayDiferencias, c.minimo, c.quienesMin, c.maximo, c.quienesMax)
			}
		})
	}
}

// spec 08 RN-8 (CA-6): nueva ronda después de revelar; las anteriores se conservan.
func TestNuevaRonda(t *testing.T) {
	v := votacionDePrueba(t)
	votar(t, &v, map[Integrante]int{ana: 3})
	if err := v.NuevaRonda(); !errors.Is(err, ErrRondaSinRevelar) {
		t.Fatalf("sin revelar: error = %v, se esperaba %v", err, ErrRondaSinRevelar)
	}
	votar(t, &v, map[Integrante]int{beto: 13})
	v.Revelar(equipo)
	if err := v.NuevaRonda(); err != nil {
		t.Fatalf("NuevaRonda: %v", err)
	}
	actual := v.RondaActual()
	if actual.Numero != 2 || len(actual.Votos) != 0 || actual.Revelada {
		t.Errorf("ronda actual = %+v, se esperaba la 2 sin votos", actual)
	}
	if len(v.Rondas) != 2 || len(v.Rondas[0].Votos) != 2 || !v.Rondas[0].Revelada {
		t.Errorf("rondas = %+v, la 1 se debía conservar revelada con sus 2 votos", v.Rondas)
	}

	votar(t, &v, map[Integrante]int{ana: 5, beto: 5})
	v.Revelar(equipo)
	v.Acordar(5)
	if err := v.NuevaRonda(); !errors.Is(err, ErrVotacionCerrada) {
		t.Errorf("votación cerrada: error = %v, se esperaba %v", err, ErrVotacionCerrada)
	}
}

// spec 08 RN-9, RN-10 (CA-7): la estimación acordada cierra la votación.
func TestAcordar(t *testing.T) {
	v := votacionDePrueba(t)
	votar(t, &v, map[Integrante]int{ana: 3, beto: 8})
	if err := v.Acordar(5); !errors.Is(err, ErrRondaSinRevelar) {
		t.Fatalf("sin revelar: error = %v, se esperaba %v", err, ErrRondaSinRevelar)
	}
	v.Revelar(equipo)
	if err := v.Acordar(4); !errors.Is(err, ErrStoryPointsFueraDeEscala) {
		t.Fatalf("fuera de la escala: error = %v, se esperaba %v", err, ErrStoryPointsFueraDeEscala)
	}
	if v.Estado != VotacionAbierta || v.Acordada != nil {
		t.Fatalf("votación = %+v, un acuerdo rechazado no debe cerrarla", v)
	}
	if err := v.Acordar(5); err != nil { // nadie votó 5: igual se acepta
		t.Fatalf("Acordar: %v", err)
	}
	if v.Estado != VotacionCerrada || v.Acordada == nil || *v.Acordada != 5 {
		t.Errorf("votación = %+v, se esperaba cerrada con la acordada 5", v)
	}
	if err := v.Acordar(8); !errors.Is(err, ErrVotacionCerrada) {
		t.Errorf("acordar otra vez: error = %v, se esperaba %v", err, ErrVotacionCerrada)
	}
}

// spec 08 RN-5 (CA-3, CA-4): la vista nunca muestra valores de una ronda sin revelar.
func TestVistaDeVotacion(t *testing.T) {
	v := votacionDePrueba(t)
	votar(t, &v, map[Integrante]int{ana: 5})

	vista := v.Vista(equipo)
	if vista.Ronda != 1 || !mismosNombres(vista.Votaron, "Ana") || !mismosNombres(vista.Faltan, "Beto") {
		t.Errorf("ronda %d, votaron %v, faltan %v; se esperaba la 1, votó Ana y falta Beto", vista.Ronda, nombres(vista.Votaron), nombres(vista.Faltan))
	}
	if len(vista.Reveladas) != 0 {
		t.Fatalf("reveladas = %+v, una ronda sin revelar no muestra votos", vista.Reveladas)
	}

	votar(t, &v, map[Integrante]int{beto: 13})
	v.Revelar(equipo)
	vista = v.Vista(equipo)
	if len(vista.Reveladas) != 1 || len(vista.Faltan) != 0 {
		t.Fatalf("reveladas = %d, faltan %v; se esperaba 1 revelada y nadie faltando", len(vista.Reveladas), nombres(vista.Faltan))
	}
	votos := vista.Reveladas[0].Votos // en el orden en que se agregaron los integrantes
	if len(votos) != 2 || votos[0].Integrante.Nombre != "Ana" || votos[0].Valor != 5 || votos[1].Integrante.Nombre != "Beto" || votos[1].Valor != 13 {
		t.Errorf("votos revelados = %+v, se esperaba Ana 5 y Beto 13", votos)
	}

	v.NuevaRonda()
	votar(t, &v, map[Integrante]int{beto: 8})
	vista = v.Vista(equipo)
	if vista.Ronda != 2 || len(vista.Reveladas) != 1 || vista.Reveladas[0].Numero != 1 ||
		!mismosNombres(vista.Votaron, "Beto") || !mismosNombres(vista.Faltan, "Ana") {
		t.Errorf("vista = %+v, se esperaba la ronda 2 sin valores y solo la 1 revelada", vista)
	}
}

// Los conflictos del Planning Poker son 409; "no existe" es 404.
func TestErroresDePlanningPoker(t *testing.T) {
	conflictos := []error{
		ErrVotacionAbierta, ErrSinDevelopers, ErrSoloDevelopers, ErrRondaRevelada,
		ErrFaltanVotos, ErrRondaSinRevelar, ErrVotacionCerrada,
	}
	for _, e := range conflictos {
		var c ErrConflicto
		if !errors.As(e, &c) {
			t.Errorf("%q debería ser un ErrConflicto", e)
		}
	}
	var c ErrConflicto
	if errors.As(ErrVotacionNoExiste, &c) {
		t.Error("ErrVotacionNoExiste no es un conflicto")
	}
}
