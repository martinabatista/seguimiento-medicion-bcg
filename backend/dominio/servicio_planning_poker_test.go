package dominio_test

import (
	"context"
	"errors"
	"testing"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/memoria"
)

// pokerDePrueba arma el servicio con repositorios en memoria: el proyecto "Tienda Online" con Ana y
// Beto (Developers), Caro (Product Owner) y la historia "Registrar usuario"; y "App Turnos" sin
// Developers: solo tiene a Eli (Scrum Master) y la historia "Pedir turno".
type pokerDePrueba struct {
	servicio          *dominio.ServicioPlanningPoker
	historias         *memoria.RepositorioHistorias
	tienda, turnos    int64
	ana, beto, caro   int64
	registrar, pedir  int64
	eliDeOtroProyecto int64
}

func nuevoPoker(t *testing.T) pokerDePrueba {
	t.Helper()
	ctx := context.Background()
	proyectos := memoria.NuevoRepositorioProyectos()
	historias := memoria.NuevoRepositorioHistorias()
	integrantes := memoria.NuevoRepositorioIntegrantes()
	backlog := dominio.NuevoServicioBacklog(proyectos, historias, memoria.NuevoRepositorioSprints(), memoria.NuevoRepositorioEsfuerzo())

	tienda, _ := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "Tienda Online", Descripcion: "d", FechaInicio: f("2026-11-01"), FechaFin: f("2026-12-15")})
	turnos, _ := proyectos.Guardar(ctx, dominio.Proyecto{Nombre: "App Turnos", Descripcion: "d", FechaInicio: f("2026-11-05"), FechaFin: f("2027-01-20")})
	agregar := func(proyectoID int64, nombre, email string, rol dominio.RolIntegrante) int64 {
		i, err := integrantes.Guardar(ctx, dominio.Integrante{ProyectoID: proyectoID, Nombre: nombre, Email: email, Rol: rol})
		if err != nil {
			t.Fatalf("agregar integrante de partida: %v", err)
		}
		return i.ID
	}
	historia := func(proyectoID int64, titulo string) int64 {
		h, _ := historias.Guardar(ctx, dominio.Historia{ProyectoID: proyectoID, Titulo: titulo, Descripcion: "d", Prioridad: 1, Estado: dominio.EstadoPendiente})
		return h.ID
	}
	p := pokerDePrueba{
		servicio:  dominio.NuevoServicioPlanningPoker(proyectos, historias, integrantes, memoria.NuevoRepositorioVotaciones(), backlog),
		historias: historias,
		tienda:    tienda.ID,
		turnos:    turnos.ID,
	}
	p.ana = agregar(tienda.ID, "Ana", "ana@mail.com", dominio.RolDeveloper)
	p.beto = agregar(tienda.ID, "Beto", "beto@mail.com", dominio.RolDeveloper)
	p.caro = agregar(tienda.ID, "Caro", "caro@mail.com", dominio.RolProductOwner)
	p.eliDeOtroProyecto = agregar(turnos.ID, "Eli", "eli@mail.com", dominio.RolScrumMaster)
	p.registrar = historia(tienda.ID, "Registrar usuario")
	p.pedir = historia(turnos.ID, "Pedir turno")
	return p
}

func nombresDe(is []dominio.Integrante) []string {
	var ns []string
	for _, i := range is {
		ns = append(ns, i.Nombre)
	}
	return ns
}

// spec 08 CA-1 a CA-7: el recorrido completo de una votación, de abrir a la estimación acordada.
func TestPlanningPokerRecorridoCompleto(t *testing.T) {
	p := nuevoPoker(t)
	ctx := context.Background()

	v, err := p.servicio.Abrir(ctx, p.tienda, p.registrar)
	if err != nil {
		t.Fatalf("Abrir: %v", err)
	}
	if v.ID == 0 || v.Estado != dominio.VotacionAbierta || v.Ronda != 1 || len(nombresDe(v.Faltan)) != 2 {
		t.Fatalf("votación = %+v, se esperaba abierta en la ronda 1 con Ana y Beto faltando", v)
	}

	// Ronda 1: votan distinto y el voto de Ana no se ve hasta revelar (CA-2, CA-3).
	if v, err = p.servicio.Votar(ctx, p.tienda, v.ID, p.ana, 3); err != nil {
		t.Fatalf("Votar Ana: %v", err)
	}
	if len(v.Reveladas) != 0 || len(v.Votaron) != 1 {
		t.Fatalf("vista = %+v, el voto de Ana tiene que estar oculto", v)
	}
	p.servicio.Votar(ctx, p.tienda, v.ID, p.beto, 13)
	if v, err = p.servicio.Revelar(ctx, p.tienda, v.ID); err != nil {
		t.Fatalf("Revelar: %v", err)
	}
	if r := v.Reveladas[0].Resultado; !r.HayDiferencias || r.Minimo.Valor != 3 || r.Maximo.Valor != 13 {
		t.Fatalf("resultado = %+v, se esperaban diferencias entre 3 y 13 (CA-5)", r)
	}

	// Ronda 2: se ponen de acuerdo y se registra la estimación (CA-6, CA-7).
	if v, err = p.servicio.NuevaRonda(ctx, p.tienda, v.ID); err != nil || v.Ronda != 2 {
		t.Fatalf("NuevaRonda = ronda %d, %v; se esperaba la 2", v.Ronda, err)
	}
	p.servicio.Votar(ctx, p.tienda, v.ID, p.ana, 5)
	p.servicio.Votar(ctx, p.tienda, v.ID, p.beto, 5)
	p.servicio.Revelar(ctx, p.tienda, v.ID)
	if v, err = p.servicio.Acordar(ctx, p.tienda, v.ID, 5); err != nil {
		t.Fatalf("Acordar: %v", err)
	}
	if v.Estado != dominio.VotacionCerrada || v.Acordada == nil || *v.Acordada != 5 {
		t.Errorf("votación = %+v, se esperaba cerrada con la acordada 5", v)
	}
	h, _ := p.historias.Obtener(ctx, p.tienda, p.registrar)
	if h.StoryPoints == nil || *h.StoryPoints != 5 {
		t.Errorf("Story Points de la historia = %v, se esperaba 5 (spec 07)", h.StoryPoints)
	}

	// Se vuelve a consultar con las dos rondas, y se puede abrir otra votación para la historia.
	visto, err := p.servicio.Ver(ctx, p.tienda, v.ID)
	if err != nil || len(visto.Reveladas) != 2 {
		t.Errorf("Ver = %d rondas reveladas, %v; se esperaban 2", len(visto.Reveladas), err)
	}
	if _, err := p.servicio.Abrir(ctx, p.tienda, p.registrar); err != nil {
		t.Errorf("abrir otra votación después de cerrar: %v", err)
	}
}

// spec 08 CA-1, CA-8: cuándo no se puede abrir una votación.
func TestPlanningPokerAbrirErrores(t *testing.T) {
	p := nuevoPoker(t)
	ctx := context.Background()
	completada, _ := p.historias.Guardar(ctx, dominio.Historia{ProyectoID: p.tienda, Titulo: "Pagar pedido", Descripcion: "d", Prioridad: 2, Estado: dominio.EstadoCompletada})
	p.servicio.Abrir(ctx, p.tienda, p.registrar)

	casos := []struct {
		nombre               string
		proyectoID, historia int64
		err                  error
	}{
		{"proyecto que no existe", 999, p.registrar, dominio.ErrProyectoNoExiste},
		{"historia que no existe", p.tienda, 999, dominio.ErrHistoriaNoExiste},
		{"historia de otro proyecto", p.tienda, p.pedir, dominio.ErrHistoriaNoExiste},
		{"historia completada", p.tienda, completada.ID, dominio.ErrHistoriaCompletada},
		{"la historia ya tiene una votación abierta", p.tienda, p.registrar, dominio.ErrVotacionAbierta},
		{"proyecto sin Developers", p.turnos, p.pedir, dominio.ErrSinDevelopers},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if _, err := p.servicio.Abrir(ctx, c.proyectoID, c.historia); !errors.Is(err, c.err) {
				t.Errorf("error = %v, se esperaba %v", err, c.err)
			}
		})
	}
}

// spec 08 CA-2, CA-8: votos rechazados por quién vota o por la votación.
func TestPlanningPokerVotarErrores(t *testing.T) {
	p := nuevoPoker(t)
	ctx := context.Background()
	v, _ := p.servicio.Abrir(ctx, p.tienda, p.registrar)

	casos := []struct {
		nombre                             string
		proyectoID, votacionID, integrante int64
		err                                error
	}{
		{"votación que no existe", p.tienda, 999, p.ana, dominio.ErrVotacionNoExiste},
		{"votación de otro proyecto", p.turnos, v.ID, p.ana, dominio.ErrVotacionNoExiste},
		{"integrante que no existe", p.tienda, v.ID, 999, dominio.ErrIntegranteNoExiste},
		{"integrante de otro proyecto", p.tienda, v.ID, p.eliDeOtroProyecto, dominio.ErrIntegranteNoExiste},
		{"integrante que no es Developer", p.tienda, v.ID, p.caro, dominio.ErrSoloDevelopers},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if _, err := p.servicio.Votar(ctx, c.proyectoID, c.votacionID, c.integrante, 5); !errors.Is(err, c.err) {
				t.Errorf("error = %v, se esperaba %v", err, c.err)
			}
		})
	}
	if visto, _ := p.servicio.Ver(ctx, p.tienda, v.ID); len(visto.Votaron) != 0 {
		t.Errorf("votaron = %v, ningún voto rechazado debía registrarse", nombresDe(visto.Votaron))
	}
	for _, op := range []func() error{
		func() error { _, err := p.servicio.Ver(ctx, p.tienda, 999); return err },
		func() error { _, err := p.servicio.Revelar(ctx, p.tienda, 999); return err },
		func() error { _, err := p.servicio.NuevaRonda(ctx, p.tienda, 999); return err },
		func() error { _, err := p.servicio.Acordar(ctx, p.tienda, 999, 5); return err },
	} {
		if err := op(); !errors.Is(err, dominio.ErrVotacionNoExiste) {
			t.Errorf("operación sobre votación inexistente: error = %v, se esperaba %v", err, dominio.ErrVotacionNoExiste)
		}
	}
}

// spec 08 RN-9: si la historia se completó durante la votación, la acordada se rechaza y la
// votación sigue abierta.
func TestPlanningPokerAcordarHistoriaCompletada(t *testing.T) {
	p := nuevoPoker(t)
	ctx := context.Background()
	v, _ := p.servicio.Abrir(ctx, p.tienda, p.registrar)
	p.servicio.Votar(ctx, p.tienda, v.ID, p.ana, 5)
	p.servicio.Votar(ctx, p.tienda, v.ID, p.beto, 5)
	p.servicio.Revelar(ctx, p.tienda, v.ID)

	h, _ := p.historias.Obtener(ctx, p.tienda, p.registrar)
	h.Estado = dominio.EstadoCompletada
	p.historias.Actualizar(ctx, h)

	if _, err := p.servicio.Acordar(ctx, p.tienda, v.ID, 5); !errors.Is(err, dominio.ErrHistoriaCompletada) {
		t.Fatalf("error = %v, se esperaba %v", err, dominio.ErrHistoriaCompletada)
	}
	if visto, _ := p.servicio.Ver(ctx, p.tienda, v.ID); visto.Estado != dominio.VotacionAbierta || visto.Acordada != nil {
		t.Errorf("votación = %+v, tenía que seguir abierta y sin acordada", visto)
	}
}
