package middleware

import (
	"context"
	"net/http"
)

// EditorVerificador diz se o usuário autenticado tem o papel de editor.
type EditorVerificador interface {
	EhEditor(ctx context.Context, usuarioID string) (bool, error)
}

// RequerEditor bloqueia (403) quem não é editor. Deve rodar depois do Autenticador.
func RequerEditor(v EditorVerificador) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			usuarioID, _ := r.Context().Value(CtxUsuarioID).(string)
			if usuarioID == "" {
				escreverErroJSON(w, http.StatusUnauthorized, "token ausente")
				return
			}
			ok, err := v.EhEditor(r.Context(), usuarioID)
			if err != nil {
				escreverErroJSON(w, http.StatusInternalServerError, "falha ao verificar permissão")
				return
			}
			if !ok {
				escreverErroJSON(w, http.StatusForbidden, "apenas editores podem realizar esta ação")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
