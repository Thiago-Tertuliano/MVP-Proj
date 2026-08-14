package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS_OrigemPermitida(t *testing.T) {
	cors := NewCORS([]string{"http://localhost:3000"})
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := cors.Handler(next)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Fatalf("header CORS ausente: %v", w.Header())
	}
	if w.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", w.Code)
	}
}

func TestCORS_OrigemNegada(t *testing.T) {
	cors := NewCORS([]string{"http://localhost:3000"})
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := cors.Handler(next)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "https://evil.example")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("não deveria liberar origem não listada")
	}
}

func TestCORS_Preflight(t *testing.T) {
	cors := NewCORS([]string{"http://localhost:3000"})
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	h := cors.Handler(next)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/login", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("esperava 204, got %d", w.Code)
	}
	if called {
		t.Fatal("preflight não deve chamar o próximo handler")
	}
}

func TestCORS_AllowAll(t *testing.T) {
	cors := NewCORS([]string{"*"})
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := cors.Handler(next)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "https://qualquer.app")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "https://qualquer.app" {
		t.Fatalf("esperava ecoar origem com *, got %q", w.Header().Get("Access-Control-Allow-Origin"))
	}
}
