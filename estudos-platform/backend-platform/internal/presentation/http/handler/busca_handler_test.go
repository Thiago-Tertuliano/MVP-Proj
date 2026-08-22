package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type fakeBusca struct {
	resp *dto.BuscaResponse
	err  error
}

func (f *fakeBusca) Execute(ctx context.Context, q string, limit int) (*dto.BuscaResponse, error) {
	return f.resp, f.err
}

func TestBuscaHandler_OK(t *testing.T) {
	h := NewBuscaHandler(&fakeBusca{resp: &dto.BuscaResponse{Itens: []dto.ResultadoBusca{{Slug: "pacotes-em-go", Titulo: "Pacotes em Go"}}}})
	req := httptest.NewRequest(http.MethodGet, "/busca?q=pacotes", nil)
	w := httptest.NewRecorder()
	h.Buscar(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", w.Code)
	}
}

func TestBuscaHandler_QueryInvalida(t *testing.T) {
	h := NewBuscaHandler(&fakeBusca{err: domainErros.ErrInvalidArgument("q deve ter pelo menos 2 caracteres", "test", nil)})
	req := httptest.NewRequest(http.MethodGet, "/busca?q=a", nil)
	w := httptest.NewRecorder()
	h.Buscar(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400, got %d", w.Code)
	}
}
