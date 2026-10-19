// Package features ejecuta los escenarios BDD (*.feature) con godog como parte de go test.
package features

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cucumber/godog"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/api"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: inicializarEscenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"."},
			Strict:   true,
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("fallaron escenarios BDD")
	}
}

// mundo guarda el estado de un escenario entre sus pasos.
type mundo struct {
	respuesta *httptest.ResponseRecorder
}

func inicializarEscenario(sc *godog.ScenarioContext) {
	m := &mundo{}
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		*m = mundo{}
		return ctx, nil
	})

	sc.Step(`^consulto el estado del servicio$`, m.consultoElEstadoDelServicio)
	sc.Step(`^la respuesta tiene código (\d+)$`, m.laRespuestaTieneCodigo)
	sc.Step(`^el estado informado es "([^"]*)"$`, m.elEstadoInformadoEs)
}

func (m *mundo) consultoElEstadoDelServicio() error {
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	m.respuesta = httptest.NewRecorder()
	api.NewRouter().ServeHTTP(m.respuesta, req)
	return nil
}

func (m *mundo) laRespuestaTieneCodigo(codigo int) error {
	if m.respuesta.Code != codigo {
		return fmt.Errorf("código = %d, se esperaba %d", m.respuesta.Code, codigo)
	}
	return nil
}

func (m *mundo) elEstadoInformadoEs(esperado string) error {
	var body map[string]string
	if err := json.Unmarshal(m.respuesta.Body.Bytes(), &body); err != nil {
		return fmt.Errorf("la respuesta no es JSON válido: %w", err)
	}
	if body["status"] != esperado {
		return fmt.Errorf("status = %q, se esperaba %q", body["status"], esperado)
	}
	return nil
}
