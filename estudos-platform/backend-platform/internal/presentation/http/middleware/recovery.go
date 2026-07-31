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
				http.Error(w, "erro interno do servidor", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}