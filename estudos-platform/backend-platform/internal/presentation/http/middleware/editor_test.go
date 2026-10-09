package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeVerif struct {
	ok  bool
	err error
}

func (f fakeVerif) EhEditor(context.Context, string) (bool, error) { return f.ok, f.err }

func chamar(v EditorVerificador, usuarioID string) int {
	chamado := false
	h := RequerEditor(v)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chamado = true
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if usuarioID != "" {
		req = req.WithContext(context.WithValue(req.Context(), CtxUsuarioID, usuarioID))
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK && chamado {
		panic("handler não deveria ter sido chamado")
	}
	return rec.Code
}

func TestRequerEditor(t *testing.T) {
	if c := chamar(fakeVerif{ok: true}, "u-1"); c != http.StatusOK {
		t.Errorf("editor: %d", c)
	}
	if c := chamar(fakeVerif{ok: false}, "u-1"); c != http.StatusForbidden {
		t.Errorf("aluno: %d", c)
	}
	if c := chamar(fakeVerif{ok: true}, ""); c != http.StatusUnauthorized {
		t.Errorf("anônimo: %d", c)
	}
	if c := chamar(fakeVerif{err: errors.New("db")}, "u-1"); c != http.StatusInternalServerError {
		t.Errorf("erro: %d", c)
	}
}
