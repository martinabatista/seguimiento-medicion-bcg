package features

// Pasos de godog para los escenarios de integrantes (spec 02, integrantes.feature).

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cucumber/godog"
)

// estadoIntegrantes guarda lo que los pasos de integrantes recuerdan dentro de un escenario.
// El proyecto de los Antecedentes ("el proyecto") es m.backlog.principal: lo anotan los pasos de proyectos.
type estadoIntegrantes struct {
	integranteID int64 // el último integrante registrado
	listaVista   int64 // proyecto de la última lista pedida (si no, el principal)
	altaPedida   bool  // en el escenario se pidió agregar un integrante: "registrado" se refiere a él
}

func (m *mundo) registrarPasosIntegrantes(sc *godog.ScenarioContext) {
	// Acciones
	sc.Step(`^agrego al integrante "([^"]*)" con email "([^"]*)" y rol "([^"]*)"$`, m.agregoAlIntegrante)
	sc.Step(`^agrego al proyecto "([^"]*)" al integrante "([^"]*)" con email "([^"]*)" y rol "([^"]*)"$`, m.agregoAlProyectoAlIntegrante)
	sc.Step(`^agrego un integrante con un nombre de (\d+) caracteres, email "([^"]*)" y rol "([^"]*)"$`, m.agregoUnIntegranteConNombreDeLargo)
	sc.Step(`^agrego al integrante "([^"]*)" con un email de (\d+) caracteres y rol "([^"]*)"$`, m.agregoAlIntegranteConEmailDeLargo)
	sc.Step(`^agrego un integrante al proyecto con identificador (\d+)$`, m.agregoUnIntegranteAlProyectoConID)
	sc.Step(`^que el proyecto tiene al integrante "([^"]*)" con email "([^"]*)" y rol "([^"]*)"$`, m.queElProyectoTieneAlIntegrante)
	sc.Step(`^listo los integrantes$`, m.listoLosIntegrantes)
	sc.Step(`^listo los integrantes del proyecto "([^"]*)"$`, m.listoLosIntegrantesDelProyecto)
	sc.Step(`^listo los integrantes del proyecto con identificador (\d+)$`, m.listoLosIntegrantesDelProyectoConID)
	sc.Step(`^veo ese integrante por su identificador$`, m.veoEseIntegrante)
	sc.Step(`^veo el integrante con identificador (\d+)$`, m.veoElIntegranteConID)
	sc.Step(`^veo ese integrante desde el proyecto "([^"]*)"$`, m.veoEseIntegranteDesdeElProyecto)

	// Verificaciones
	sc.Step(`^el integrante queda registrado con un identificador$`, m.elIntegranteQuedaRegistrado)
	sc.Step(`^el integrante tiene nombre "([^"]*)", email "([^"]*)" y rol "([^"]*)"$`, m.elIntegranteTieneTodo)
	sc.Step(`^el integrante tiene rol "([^"]*)"$`, m.elIntegranteTieneRol)
	sc.Step(`^el integrante tiene email "([^"]*)"$`, m.elIntegranteTieneEmail)
	sc.Step(`^veo los integrantes "([^"]*)" y "([^"]*)" en ese orden$`, m.veoLosIntegrantesEnOrden)
	sc.Step(`^la lista de integrantes está vacía$`, m.laListaDeIntegrantesEstaVacia)
	sc.Step(`^el proyecto tiene (\d+) integrantes?$`, m.elProyectoTieneIntegrantes)
}

// ---------- helpers ----------

type integrante struct {
	ID         int64  `json:"id"`
	ProyectoID int64  `json:"proyecto_id"`
	Nombre     string `json:"nombre"`
	Email      string `json:"email"`
	Rol        string `json:"rol"`
}

func rutaIntegrantes(proyectoID int64) string {
	return fmt.Sprintf("/api/proyectos/%d/integrantes", proyectoID)
}

func rutaIntegrante(proyectoID, integranteID int64) string {
	return fmt.Sprintf("/api/proyectos/%d/integrantes/%d", proyectoID, integranteID)
}

// agregarIntegrante pide agregar el integrante con los datos tal cual se escribieron en el
// escenario, para que los valide el sistema y no el paso.
func (m *mundo) agregarIntegrante(proyectoID int64, nombre, email, rol string) error {
	m.integrantes.altaPedida = true
	cuerpo := map[string]string{"nombre": texto(nombre), "email": texto(email), "rol": texto(rol)}
	if err := m.pedir(http.MethodPost, rutaIntegrantes(proyectoID), cuerpo); err != nil {
		return err
	}
	if m.respuesta.Code == http.StatusCreated {
		var creado integrante
		if err := json.Unmarshal(m.respuesta.Body.Bytes(), &creado); err != nil {
			return err
		}
		m.integrantes.integranteID = creado.ID
	}
	return nil
}

func (m *mundo) obtenerIntegrante(proyectoID, integranteID int64) (integrante, error) {
	rec := httptest.NewRecorder()
	m.api.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, rutaIntegrante(proyectoID, integranteID), nil))
	var i integrante
	if rec.Code != http.StatusOK {
		return i, fmt.Errorf("no se pudo ver el integrante %d: código %d", integranteID, rec.Code)
	}
	return i, json.Unmarshal(rec.Body.Bytes(), &i)
}

func (m *mundo) listarIntegrantes(proyectoID int64) ([]integrante, error) {
	rec := httptest.NewRecorder()
	m.api.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, rutaIntegrantes(proyectoID), nil))
	var lista []integrante
	if rec.Code != http.StatusOK {
		return nil, fmt.Errorf("no se pudo listar los integrantes: código %d", rec.Code)
	}
	return lista, json.Unmarshal(rec.Body.Bytes(), &lista)
}

// ultimoIntegrante devuelve el último integrante registrado en el proyecto principal.
func (m *mundo) ultimoIntegrante() (integrante, error) {
	return m.obtenerIntegrante(m.backlog.principal, m.integrantes.integranteID)
}

// ---------- acciones ----------

func (m *mundo) agregoAlIntegrante(nombre, email, rol string) error {
	return m.agregarIntegrante(m.backlog.principal, nombre, email, rol)
}

func (m *mundo) agregoAlProyectoAlIntegrante(proyecto, nombre, email, rol string) error {
	id, err := m.proyectoLlamado(proyecto)
	if err != nil {
		return err
	}
	return m.agregarIntegrante(id, nombre, email, rol)
}

func (m *mundo) agregoUnIntegranteConNombreDeLargo(largo int, email, rol string) error {
	return m.agregarIntegrante(m.backlog.principal, strings.Repeat("n", largo), email, rol)
}

// agregoAlIntegranteConEmailDeLargo arma un email válido de exactamente ese largo.
func (m *mundo) agregoAlIntegranteConEmailDeLargo(nombre string, largo int, rol string) error {
	const dominio = "@mail.com"
	return m.agregarIntegrante(m.backlog.principal, nombre, strings.Repeat("a", largo-len(dominio))+dominio, rol)
}

func (m *mundo) agregoUnIntegranteAlProyectoConID(proyectoID int64) error {
	return m.agregarIntegrante(proyectoID, "Ana Pérez", "ana@mail.com", "Developer")
}

func (m *mundo) queElProyectoTieneAlIntegrante(nombre, email, rol string) error {
	if err := m.agregarIntegrante(m.backlog.principal, nombre, email, rol); err != nil {
		return err
	}
	if m.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("no se pudo agregar el integrante de partida: %s", m.respuesta.Body.String())
	}
	return nil
}

func (m *mundo) listoLosIntegrantes() error {
	return m.listoLosIntegrantesDelProyectoConID(m.backlog.principal)
}

func (m *mundo) listoLosIntegrantesDelProyecto(nombre string) error {
	id, err := m.proyectoLlamado(nombre)
	if err != nil {
		return err
	}
	return m.listoLosIntegrantesDelProyectoConID(id)
}

func (m *mundo) listoLosIntegrantesDelProyectoConID(proyectoID int64) error {
	m.integrantes.listaVista = proyectoID
	return m.pedir(http.MethodGet, rutaIntegrantes(proyectoID), nil)
}

func (m *mundo) veoEseIntegrante() error {
	return m.pedir(http.MethodGet, rutaIntegrante(m.backlog.principal, m.integrantes.integranteID), nil)
}

func (m *mundo) veoElIntegranteConID(id int64) error {
	return m.pedir(http.MethodGet, rutaIntegrante(m.backlog.principal, id), nil)
}

func (m *mundo) veoEseIntegranteDesdeElProyecto(nombre string) error {
	id, err := m.proyectoLlamado(nombre)
	if err != nil {
		return err
	}
	return m.pedir(http.MethodGet, rutaIntegrante(id, m.integrantes.integranteID), nil)
}

// ---------- verificaciones ----------

func (m *mundo) elIntegranteQuedaRegistrado() error {
	if m.respuesta.Code != http.StatusCreated || m.integrantes.integranteID == 0 {
		return fmt.Errorf("no se registró: código %d, %s", m.respuesta.Code, m.respuesta.Body.String())
	}
	return nil
}

func (m *mundo) elIntegranteTieneTodo(nombre, email, rol string) error {
	i, err := m.ultimoIntegrante()
	if err != nil {
		return err
	}
	if i.Nombre != nombre || i.Email != email || i.Rol != rol {
		return fmt.Errorf("integrante = %+v, se esperaba nombre %q, email %q y rol %q", i, nombre, email, rol)
	}
	return nil
}

// elIntegranteTieneRol verifica la respuesta del alta, que puede ser de cualquier proyecto.
func (m *mundo) elIntegranteTieneRol(rol string) error {
	var i integrante
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &i); err != nil {
		return err
	}
	if i.Rol != rol {
		return fmt.Errorf("rol = %q, se esperaba %q", i.Rol, rol)
	}
	return nil
}

func (m *mundo) elIntegranteTieneEmail(email string) error {
	i, err := m.ultimoIntegrante()
	if err != nil {
		return err
	}
	if i.Email != email {
		return fmt.Errorf("email = %q, se esperaba %q", i.Email, email)
	}
	return nil
}

func (m *mundo) veoLosIntegrantesEnOrden(primero, segundo string) error {
	var lista []integrante
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &lista); err != nil {
		return err
	}
	if len(lista) != 2 || lista[0].Nombre != primero || lista[1].Nombre != segundo {
		return fmt.Errorf("integrantes = %+v, se esperaba %q y %q en ese orden", lista, primero, segundo)
	}
	return nil
}

// laListaDeIntegrantesEstaVacia vuelve a listar el último proyecto pedido (o el principal) y espera [].
func (m *mundo) laListaDeIntegrantesEstaVacia() error {
	id := m.integrantes.listaVista
	if id == 0 {
		id = m.backlog.principal
	}
	return m.proyectoConIntegrantes(id, 0)
}

func (m *mundo) elProyectoTieneIntegrantes(cantidad int) error {
	return m.proyectoConIntegrantes(m.backlog.principal, cantidad)
}

func (m *mundo) proyectoConIntegrantes(proyectoID int64, cantidad int) error {
	lista, err := m.listarIntegrantes(proyectoID)
	if err != nil {
		return err
	}
	if len(lista) != cantidad {
		return fmt.Errorf("integrantes = %+v, se esperaban %d", lista, cantidad)
	}
	return nil
}
