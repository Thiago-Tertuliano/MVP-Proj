package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/port"
)

const (
	CtxUsuarioID ctxKey = "usuario_id"
	CtxEmail     ctxKey = "email"
)

// Autenticador protege rotas que exigem JWT válido.
type Autenticador struct {
	tokens port.TokenGerador
}

func NewAutenticador(tokens port.TokenGerador) *Autenticador {
	return &Autenticador{tokens: tokens}
}

func (a *Autenticador) Proteger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, `{"erro":"token ausente"}`, http.StatusUnauthorized)
			return
		}

		claims, err := a.tokens.ValidarAccessToken(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			http.Error(w, `{"erro":"token inválido"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), CtxUsuarioID, claims.UsuarioID)
		ctx = context.WithValue(ctx, CtxEmail, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}