package features

// Pasos de godog para los escenarios de Planning Poker (spec 08, planning-poker.feature).
// Reutiliza los pasos de proyectos, integrantes, backlog y Story Points para armar el estado.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

// estadoPoker guarda lo que los pasos de Planning Poker recuerdan dentro de un escenario.
type estadoPoker struct {
	votacionID int64 // la última votación abierta
}

func (m *mundo) registrarPasosPlanningPoker(sc *godog.ScenarioContext) {
	// Acciones
	sc.Step(`^abro una votación para la historia "([^"]*)"$`, m.abroUnaVotacion)
	sc.Step(`^que abrí una votación para la historia "([^"]*)"$`, m.queAbriUnaVotacion)
	sc.Step(`^abro una votación para la historia "([^"]*)" del proyecto "([^"]*)"$`, m.abroUnaVotacionDelProyecto)
	sc.Step(`^abro una votación para la historia con identificador (\d+)$`, m.abroUnaVotacionParaLaHistoriaConID)
	sc.Step(`^"([^"]*)" vota (\d+)$`, m.vota)
	sc.Step(`^que "([^"]*)" votó (\d+)$`, m.queVoto)
	sc.Step(`^"([^"]*)" vota "([^"]*)"$`, m.votaConValorEscrito)
	sc.Step(`^veo la votación$`, m.veoLaVotacion)
	sc.Step(`^veo la votación con identificador (\d+)$`, m.veoLaVotacionConID)
	sc.Step(`^revelo la ronda$`, m.reveloLaRonda)
	sc.Step(`^que revelé la ronda$`, m.queReveleLaRonda)
	sc.Step(`^empiezo una nueva ronda$`, m.empiezoUnaNuevaRonda)
	sc.Step(`^registro la estimación acordada de (\d+)$`, m.registroLaEstimacionAcordada)
	sc.Step(`^que registré la estimación acordada de (\d+)$`, m.queRegistreLaEstimacionAcordada)

	// Verificaciones
	sc.Step(`^la votación está abierta en la ronda (\d+)$`, m.laVotacionEstaAbiertaEnLaRonda)
	sc.Step(`^la votación está cerrada con la estimación acordada de (\d+)$`, m.laVotacionEstaCerrada)
	sc.Step(`^faltan? votar ((?:"[^"]*"(?: y )?)+)$`, m.faltanVotar)
	sc.Step(`^ya votó "([^"]*)"$`, m.yaVoto)
	sc.Step(`^no se ve ningún valor de los votos$`, m.noSeVeNingunValor)
	sc.Step(`^veo los votos: "([^"]*)" (\d+) y "([^"]*)" (\d+)$`, m.veoLosVotos)
	sc.Step(`^en la ronda (\d+) se ven los votos: "([^"]*)" (\d+) y "([^"]*)" (\d+)$`, m.enLaRondaSeVenLosVotos)
	sc.Step(`^hay diferencias entre los votos$`, m.hayDiferencias)
	sc.Step(`^no hay diferencias entre los votos$`, m.noHayDiferencias)
	sc.Step(`^el voto más (bajo|alto) es (\d+) de ((?:"[^"]*"(?: y )?)+)$`, m.elVotoExtremoEs)
}

// ---------- helpers ----------

type integranteVotacion struct {
	ID     int64  `json:"id"`
	Nombre string `json:"nombre"`
}

type votoRevelado struct {
	IntegranteID int64  `json:"integrante_id"`
	Nombre       string `json:"nombre"`
	StoryPoints  int    `json:"story_points"`
}

type extremo struct {
	StoryPoints int                  `json:"story_points"`
	Integrantes []integranteVotacion `json:"integrantes"`
}

type rondaRevelada struct {
	Numero         int            `json:"numero"`
	Votos          []votoRevelado `json:"votos"`
	HayDiferencias bool           `json:"hay_diferencias"`
	MasBajo        extremo        `json:"mas_bajo"`
	MasAlto        extremo        `json:"mas_alto"`
}

type votacion struct {
	ID       int64                `json:"id"`
	Estado   string               `json:"estado"`
	Ronda    int                  `json:"ronda"`
	Acordada *int                 `json:"estimacion_acordada"`
	Votaron  []integranteVotacion `json:"votaron"`
	Faltan   []integranteVotacion `json:"faltan"`
	Rondas   []rondaRevelada      `json:"rondas"`
}

var nombresEntreComillas = regexp.MustCompile(`"([^"]*)"`)

// listaDeNombres convierte `"Ana" y "Beto"` en [Ana Beto].
func listaDeNombres(texto string) []string {
	var ns []string
	for _, m := range nombresEntreComillas.FindAllStringSubmatch(texto, -1) {
		ns = append(ns, m[1])
	}
	return ns
}

func nombresVotacion(is []integranteVotacion) []string {
	var ns []string
	for _, i := range is {
		ns = append(ns, i.Nombre)
	}
	return ns
}

func rutaVotaciones(proyectoID int64) string {
	return fmt.Sprintf("/api/proyectos/%d/votaciones", proyectoID)
}

func rutaVotacion(proyectoID, votacionID int64) string {
	return fmt.Sprintf("%s/%d", rutaVotaciones(proyectoID), votacionID)
}

// integranteLlamado busca al integrante por nombre en los proyectos creados en el escenario.
func (m *mundo) integranteLlamado(nombre string) (int64, error) {
	for _, proyectoID := range m.backlog.proyectos {
		rec := httptest.NewRecorder()
		m.api.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, rutaIntegrantes(proyectoID), nil))
		var lista []integrante
		if err := json.Unmarshal(rec.Body.Bytes(), &lista); err != nil {
			return 0, err
		}
		for _, i := range lista {
			if i.Nombre == nombre {
				return i.ID, nil
			}
		}
	}
	return 0, fmt.Errorf("el escenario no agregó al integrante %q", nombre)
}

func (m *mundo) abrirVotacion(proyectoID, historiaID int64) error {
	if err := m.pedir(http.MethodPost, rutaVotaciones(proyectoID), map[string]int64{"historia_id": historiaID}); err != nil {
		return err
	}
	if m.respuesta.Code == http.StatusCreated {
		var v votacion
		if err := json.Unmarshal(m.respuesta.Body.Bytes(), &v); err != nil {
			return err
		}
		m.poker.votacionID = v.ID
	}
	return nil
}

// votar manda el voto tal cual se escribió en el escenario (puede ser "2.5" o faltar).
func (m *mundo) votar(nombre, valor string) error {
	id, err := m.integranteLlamado(nombre)
	if err != nil {
		return err
	}
	cuerpo := map[string]any{"integrante_id": id}
	if valor != "" {
		cuerpo["story_points"] = json.Number(valor)
	}
	return m.pedir(http.MethodPost, rutaVotacion(m.backlog.principal, m.poker.votacionID)+"/votos", cuerpo)
}

// verVotacion lee la votación actual desde el sistema.
func (m *mundo) verVotacion() (votacion, error) {
	rec := httptest.NewRecorder()
	m.api.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, rutaVotacion(m.backlog.principal, m.poker.votacionID), nil))
	var v votacion
	if rec.Code != http.StatusOK {
		return v, fmt.Errorf("no se pudo ver la votación %d: código %d", m.poker.votacionID, rec.Code)
	}
	return v, json.Unmarshal(rec.Body.Bytes(), &v)
}

// ultimaRevelada devuelve la ronda revelada más reciente.
func (m *mundo) ultimaRevelada() (rondaRevelada, error) {
	v, err := m.verVotacion()
	if err != nil {
		return rondaRevelada{}, err
	}
	if len(v.Rondas) == 0 {
		return rondaRevelada{}, fmt.Errorf("la votación no tiene rondas reveladas: %+v", v)
	}
	return v.Rondas[len(v.Rondas)-1], nil
}

func compararVotos(r rondaRevelada, nombre1 string, valor1 int, nombre2 string, valor2 int) error {
	esperado := []votoRevelado{{Nombre: nombre1, StoryPoints: valor1}, {Nombre: nombre2, StoryPoints: valor2}}
	if len(r.Votos) != 2 || r.Votos[0].Nombre != esperado[0].Nombre || r.Votos[0].StoryPoints != valor1 ||
		r.Votos[1].Nombre != esperado[1].Nombre || r.Votos[1].StoryPoints != valor2 {
		return fmt.Errorf("votos de la ronda %d = %+v, se esperaba %s %d y %s %d", r.Numero, r.Votos, nombre1, valor1, nombre2, valor2)
	}
	return nil
}

func (m *mundo) exigirOK(paso string) error {
	if m.respuesta.Code != http.StatusOK && m.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("no se pudo %s: %s", paso, m.respuesta.Body.String())
	}
	return nil
}

// ---------- acciones ----------

func (m *mundo) abroUnaVotacion(titulo string) error {
	id, err := m.historiaLlamada(titulo)
	if err != nil {
		return err
	}
	return m.abrirVotacion(m.backlog.principal, id)
}

func (m *mundo) queAbriUnaVotacion(titulo string) error {
	if err := m.abroUnaVotacion(titulo); err != nil {
		return err
	}
	return m.exigirOK("abrir la votación de partida")
}

func (m *mundo) abroUnaVotacionDelProyecto(titulo, proyecto string) error {
	historiaID, err := m.historiaLlamada(titulo)
	if err != nil {
		return err
	}
	proyectoID, err := m.proyectoLlamado(proyecto)
	if err != nil {
		return err
	}
	return m.abrirVotacion(proyectoID, historiaID)
}

func (m *mundo) abroUnaVotacionParaLaHistoriaConID(historiaID int64) error {
	return m.abrirVotacion(m.backlog.principal, historiaID)
}

func (m *mundo) vota(nombre string, valor int) error {
	return m.votar(nombre, strconv.Itoa(valor))
}

func (m *mundo) queVoto(nombre string, valor int) error {
	if err := m.vota(nombre, valor); err != nil {
		return err
	}
	return m.exigirOK("registrar el voto de partida")
}

func (m *mundo) votaConValorEscrito(nombre, valor string) error {
	return m.votar(nombre, texto(valor))
}

func (m *mundo) veoLaVotacion() error {
	return m.pedir(http.MethodGet, rutaVotacion(m.backlog.principal, m.poker.votacionID), nil)
}

func (m *mundo) veoLaVotacionConID(id int64) error {
	return m.pedir(http.MethodGet, rutaVotacion(m.backlog.principal, id), nil)
}

func (m *mundo) reveloLaRonda() error {
	return m.pedir(http.MethodPost, rutaVotacion(m.backlog.principal, m.poker.votacionID)+"/revelar", nil)
}

func (m *mundo) queReveleLaRonda() error {
	if err := m.reveloLaRonda(); err != nil {
		return err
	}
	return m.exigirOK("revelar la ronda de partida")
}

func (m *mundo) empiezoUnaNuevaRonda() error {
	return m.pedir(http.MethodPost, rutaVotacion(m.backlog.principal, m.poker.votacionID)+"/rondas", nil)
}

func (m *mundo) registroLaEstimacionAcordada(valor int) error {
	return m.pedir(http.MethodPost, rutaVotacion(m.backlog.principal, m.poker.votacionID)+"/acordar",
		map[string]int{"story_points": valor})
}

func (m *mundo) queRegistreLaEstimacionAcordada(valor int) error {
	if err := m.registroLaEstimacionAcordada(valor); err != nil {
		return err
	}
	return m.exigirOK("registrar la estimación acordada de partida")
}

// ---------- verificaciones ----------

func (m *mundo) laVotacionEstaAbiertaEnLaRonda(ronda int) error {
	v, err := m.verVotacion()
	if err != nil {
		return err
	}
	if v.Estado != "Abierta" || v.Ronda != ronda {
		return fmt.Errorf("votación %s en la ronda %d, se esperaba Abierta en la ronda %d", v.Estado, v.Ronda, ronda)
	}
	return nil
}

func (m *mundo) laVotacionEstaCerrada(acordada int) error {
	v, err := m.verVotacion()
	if err != nil {
		return err
	}
	if v.Estado != "Cerrada" || v.Acordada == nil || *v.Acordada != acordada {
		return fmt.Errorf("votación %s con acordada %v, se esperaba Cerrada con %d", v.Estado, v.Acordada, acordada)
	}
	return nil
}

func (m *mundo) faltanVotar(lista string) error {
	v, err := m.verVotacion()
	if err != nil {
		return err
	}
	if faltan := nombresVotacion(v.Faltan); !slices.Equal(faltan, listaDeNombres(lista)) {
		return fmt.Errorf("faltan votar %v, se esperaba %v", faltan, listaDeNombres(lista))
	}
	return nil
}

func (m *mundo) yaVoto(nombre string) error {
	v, err := m.verVotacion()
	if err != nil {
		return err
	}
	if !slices.Contains(nombresVotacion(v.Votaron), nombre) {
		return fmt.Errorf("votaron %v, se esperaba que haya votado %s", nombresVotacion(v.Votaron), nombre)
	}
	return nil
}

// noSeVeNingunValor revisa la respuesta de "veo la votación": ningún valor de voto ni ronda revelada.
func (m *mundo) noSeVeNingunValor() error {
	cuerpo := m.respuesta.Body.String()
	var v votacion
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &v); err != nil {
		return fmt.Errorf("la respuesta no es JSON válido (%s): %w", cuerpo, err)
	}
	if strings.Contains(cuerpo, `"story_points"`) || len(v.Rondas) != 0 {
		return fmt.Errorf("la votación muestra valores de votos sin revelar: %s", cuerpo)
	}
	return nil
}

func (m *mundo) veoLosVotos(nombre1 string, valor1 int, nombre2 string, valor2 int) error {
	r, err := m.ultimaRevelada()
	if err != nil {
		return err
	}
	return compararVotos(r, nombre1, valor1, nombre2, valor2)
}

func (m *mundo) enLaRondaSeVenLosVotos(numero int, nombre1 string, valor1 int, nombre2 string, valor2 int) error {
	v, err := m.verVotacion()
	if err != nil {
		return err
	}
	for _, r := range v.Rondas {
		if r.Numero == numero {
			return compararVotos(r, nombre1, valor1, nombre2, valor2)
		}
	}
	return fmt.Errorf("la ronda %d no está entre las reveladas: %+v", numero, v.Rondas)
}

func (m *mundo) hayDiferencias() error {
	r, err := m.ultimaRevelada()
	if err != nil {
		return err
	}
	if !r.HayDiferencias {
		return fmt.Errorf("ronda %d sin diferencias, se esperaban: %+v", r.Numero, r.Votos)
	}
	return nil
}

func (m *mundo) noHayDiferencias() error {
	r, err := m.ultimaRevelada()
	if err != nil {
		return err
	}
	if r.HayDiferencias {
		return fmt.Errorf("ronda %d con diferencias, no se esperaban: %+v", r.Numero, r.Votos)
	}
	return nil
}

func (m *mundo) elVotoExtremoEs(cual string, valor int, lista string) error {
	r, err := m.ultimaRevelada()
	if err != nil {
		return err
	}
	e := r.MasAlto
	if cual == "bajo" {
		e = r.MasBajo
	}
	if quienes := nombresVotacion(e.Integrantes); e.StoryPoints != valor || !slices.Equal(quienes, listaDeNombres(lista)) {
		return fmt.Errorf("voto más %s = %d de %v, se esperaba %d de %v", cual, e.StoryPoints, quienes, valor, listaDeNombres(lista))
	}
	return nil
}
