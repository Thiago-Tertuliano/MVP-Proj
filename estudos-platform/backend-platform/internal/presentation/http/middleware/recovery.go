package middleware

import (
	"log/slog"
	"net/http"
)

// Recovery converte panic em 500 para o servidor nunca cair.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recuperado", "erro", rec)
				escreverErroJSON(w, http.StatusInternalServerError, "erro interno do servidor")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
