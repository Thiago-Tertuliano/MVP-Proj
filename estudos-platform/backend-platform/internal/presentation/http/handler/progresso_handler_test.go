package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
	"github.com/thiago-tertuliano/estudos-platform/internal/presentation/http/middleware"
)

type fakeMarcarLido struct {
	resp *dto.ProgressoArtigoResponse
	err  error
}

func (f *fakeMarcarLido) Execute(ctx context.Context, usuarioID, artigoID string, concluido bool) (*dto.ProgressoArtigoResponse, error) {
	return f.resp, f.err
}

func TestProgressoHandler_MarcarArtigo_OK(t *testing.T) {
	id := uuid.New().String()
	h := NewProgressoHandler(&fakeMarcarLido{
		resp: &dto.ProgressoArtigoResponse{ArtigoID: id, Concluido: true},
	})

	r := chi.NewRouter()
	r.Put("/progresso/artigos/{id}", h.MarcarArtigo)

	req := httptest.NewRequest(http.MethodPut, "/progresso/artigos/"+id, strings.NewReader(`{"concluido":true}`))
	req = req.WithContext(context.WithValue(req.Context(), middleware.CtxUsuarioID, uuid.New().String()))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestProgressoHandler_MarcarArtigo_NotFound(t *testing.T) {
	h := NewProgressoHandler(&fakeMarcarLido{
		err: domainErros.ErrNotFound("artigo não encontrado", "test", nil),
	})
	r := chi.NewRouter()
	r.Put("/progresso/artigos/{id}", h.MarcarArtigo)

	req := httptest.NewRequest(http.MethodPut, "/progresso/artigos/"+uuid.New().String(), strings.NewReader(`{"concluido":true}`))
	req = req.WithContext(context.WithValue(req.Context(), middleware.CtxUsuarioID, uuid.New().String()))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("esperava 404, got %d", w.Code)
	}
}

func TestProgressoHandler_MarcarArtigo_BodyInvalido(t *testing.T) {
	h := NewProgressoHandler(&fakeMarcarLido{})
	r := chi.NewRouter()
	r.Put("/progresso/artigos/{id}", h.MarcarArtigo)

	req := httptest.NewRequest(http.MethodPut, "/progresso/artigos/"+uuid.New().String(), strings.NewReader(`{`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400, got %d", w.Code)
	}
}
