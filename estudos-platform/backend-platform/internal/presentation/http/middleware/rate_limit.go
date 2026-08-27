package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// RateLimit é um sliding window em memória (1 instância). 10 req/min por IP no auth público.
type RateLimit struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	limite int
	janela time.Duration
}

func NewRateLimit(limite int, janela time.Duration) *RateLimit {
	return &RateLimit{
		hits:   make(map[string][]time.Time),
		limite: limite,
		janela: janela,
	}
}

func (rl *RateLimit) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		if !rl.permite(ip) {
			w.Header().Set("Retry-After", "60")
			escreverErroJSON(w, http.StatusTooManyRequests, "muitas tentativas — aguarde 1 minuto")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimit) permite(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	agora := time.Now()
	corte := agora.Add(-rl.janela)

	vivos := rl.hits[ip][:0]
	for _, t := range rl.hits[ip] {
		if t.After(corte) {
			vivos = append(vivos, t)
		}
	}

	if len(vivos) >= rl.limite {
		rl.hits[ip] = vivos
		return false
	}

	rl.hits[ip] = append(vivos, agora)
	return true
}
