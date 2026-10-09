package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimit é um sliding window em memória (1 instância). 10 req/min por IP no auth público.
type RateLimit struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	limite int
	janela time.Duration
	// confiarProxy: atrás de proxy/CDN o RemoteAddr é o do proxy (todos os usuários dividiriam o limite).
	confiarProxy bool
}

func NewRateLimit(limite int, janela time.Duration) *RateLimit {
	return &RateLimit{
		hits:   make(map[string][]time.Time),
		limite: limite,
		janela: janela,
	}
}

// ConfiandoNoProxy passa a identificar o cliente pelo 1º IP de X-Forwarded-For.
func (rl *RateLimit) ConfiandoNoProxy(confiar bool) *RateLimit {
	rl.confiarProxy = confiar
	return rl
}

func (rl *RateLimit) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.permite(rl.ipCliente(r)) {
			w.Header().Set("Retry-After", "60")
			escreverErroJSON(w, http.StatusTooManyRequests, "muitas tentativas — aguarde 1 minuto")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimit) ipCliente(r *http.Request) string {
	if rl.confiarProxy {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			if primeiro := strings.TrimSpace(strings.Split(xff, ",")[0]); primeiro != "" {
				return primeiro
			}
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
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
