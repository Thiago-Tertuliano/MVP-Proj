package middleware

import (
	"context"
	"net/http"
	"strconv"
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
		token := ""
		if strings.HasPrefix(header, "Bearer ") {
			token = strings.TrimPrefix(header, "Bearer ")
		}
		if token == "" {
			if c, err := r.Cookie("access_token"); err == nil {
				token = c.Value
			}
		}
		if token == "" {
			escreverErroJSON(w, http.StatusUnauthorized, "token ausente")
			return
		}

		claims, err := a.tokens.ValidarAccessToken(token)
		if err != nil {
			escreverErroJSON(w, http.StatusUnauthorized, "token inválido")
			return
		}

		ctx := context.WithValue(r.Context(), CtxUsuarioID, claims.UsuarioID)
		ctx = context.WithValue(ctx, CtxEmail, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// escreverErroJSON padroniza respostas de erro em JSON na camada de middleware.
func escreverErroJSON(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"erro":` + strconv.Quote(msg) + `}`))
}
