package middleware

import (
	"net/http"
	"strings"
)

// CORS libera origens configuradas e responde preflight OPTIONS.
type CORS struct {
	allowedOrigins map[string]struct{}
	allowAll       bool
}

func NewCORS(origins []string) *CORS {
	c := &CORS{allowedOrigins: make(map[string]struct{}, len(origins))}
	for _, o := range origins {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		if o == "*" {
			c.allowAll = true
			continue
		}
		c.allowedOrigins[o] = struct{}{}
	}
	return c
}

func (c *CORS) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && c.origemPermitida(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Max-Age", "600")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (c *CORS) origemPermitida(origin string) bool {
	if c.allowAll {
		return true
	}
	_, ok := c.allowedOrigins[origin]
	return ok
}
