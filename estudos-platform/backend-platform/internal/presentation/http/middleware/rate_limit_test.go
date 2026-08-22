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
