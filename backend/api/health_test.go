package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthRespondeOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	NewRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, se esperaba %q", ct, "application/json")
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("la respuesta no es JSON válido: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %q, se esperaba %q", body["status"], "ok")
	}
}

func TestHealthRechazaMetodoNoPermitido(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/health", nil)
	rec := httptest.NewRecorder()

	NewRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("código = %d, se esperaba %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
