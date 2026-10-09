package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRateLimit_BloqueiaNa11a(t *testing.T) {
	rl := NewRateLimit(10, time.Minute)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := rl.Handler(next)

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		req.RemoteAddr = "1.2.3.4:9999"
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("req %d deveria passar, got %d", i+1, w.Code)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	req.RemoteAddr = "1.2.3.4:9999"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("esperava 429, got %d", w.Code)
	}
	if w.Header().Get("Retry-After") != "60" {
		t.Fatalf("Retry-After ausente: %q", w.Header().Get("Retry-After"))
	}
	if !strings.Contains(w.Body.String(), "muitas tentativas") {
		t.Fatalf("corpo 429 inesperado: %s", w.Body.String())
	}
}

func TestRateLimit_IPDiferenteNaoBloqueia(t *testing.T) {
	rl := NewRateLimit(1, time.Minute)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := rl.Handler(next)

	reqA := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	reqA.RemoteAddr = "1.1.1.1:1"
	wA := httptest.NewRecorder()
	h.ServeHTTP(wA, reqA)

	reqB := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	reqB.RemoteAddr = "2.2.2.2:1"
	wB := httptest.NewRecorder()
	h.ServeHTTP(wB, reqB)

	if wA.Code != http.StatusOK || wB.Code != http.StatusOK {
		t.Fatalf("IPs distintos devem passar: %d %d", wA.Code, wB.Code)
	}
}

func TestRateLimit_IPCliente(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	req.RemoteAddr = "10.0.0.1:443"
	req.Header.Set("X-Forwarded-For", " 203.0.113.7 , 76.76.21.21")

	if ip := NewRateLimit(1, time.Minute).ipCliente(req); ip != "10.0.0.1" {
		t.Fatalf("sem confiar no proxy deve usar RemoteAddr, got %q", ip)
	}
	if ip := NewRateLimit(1, time.Minute).ConfiandoNoProxy(true).ipCliente(req); ip != "203.0.113.7" {
		t.Fatalf("confiando no proxy deve usar o 1º X-Forwarded-For, got %q", ip)
	}

	req.Header.Del("X-Forwarded-For")
	if ip := NewRateLimit(1, time.Minute).ConfiandoNoProxy(true).ipCliente(req); ip != "10.0.0.1" {
		t.Fatalf("sem cabeçalho cai no RemoteAddr, got %q", ip)
	}
}

func TestRateLimit_AtrasDoProxyUsuariosNaoDividemLimite(t *testing.T) {
	h := NewRateLimit(1, time.Minute).ConfiandoNoProxy(true).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	for _, cliente := range []string{"203.0.113.1", "203.0.113.2"} {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		req.RemoteAddr = "10.0.0.1:443"
		req.Header.Set("X-Forwarded-For", cliente)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("cliente %s não deveria ser limitado pelo outro: %d", cliente, w.Code)
		}
	}
}
