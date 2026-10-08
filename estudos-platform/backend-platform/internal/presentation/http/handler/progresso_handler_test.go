package handler

import (
	"context"
	"encoding/json"
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

type fakeProgressoTrilha struct {
	resp *dto.ProgressoTrilhaResponse
	err  error
}

func (f *fakeProgressoTrilha) Execute(ctx context.Context, usuarioID, trilhaID string) (*dto.ProgressoTrilhaResponse, error) {
	return f.resp, f.err
}

type fakeContinuar struct {
	resp *dto.ContinuarResponse
	err  error
}

func (f *fakeContinuar) Execute(ctx context.Context, usuarioID string) (*dto.ContinuarResponse, error) {
	return f.resp, f.err
}

func comUsuario(req *http.Request) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), middleware.CtxUsuarioID, uuid.New().String()))
}

func TestProgressoHandler_MarcarArtigo_OK(t *testing.T) {
	id := uuid.New().String()
	h := NewProgressoHandler(&fakeMarcarLido{
		resp: &dto.ProgressoArtigoResponse{ArtigoID: id, Concluido: true},
	}, nil, nil)

	r := chi.NewRouter()
	r.Put("/progresso/artigos/{id}", h.MarcarArtigo)

	req := comUsuario(httptest.NewRequest(http.MethodPut, "/progresso/artigos/"+id, strings.NewReader(`{"concluido":true}`)))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestProgressoHandler_MarcarArtigo_NotFound(t *testing.T) {
	h := NewProgressoHandler(&fakeMarcarLido{
		err: domainErros.ErrNotFound("artigo não encontrado", "test", nil),
	}, nil, nil)
	r := chi.NewRouter()
	r.Put("/progresso/artigos/{id}", h.MarcarArtigo)

	req := comUsuario(httptest.NewRequest(http.MethodPut, "/progresso/artigos/"+uuid.New().String(), strings.NewReader(`{"concluido":true}`)))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("esperava 404, got %d", w.Code)
	}
}

func TestProgressoHandler_MarcarArtigo_BodyInvalido(t *testing.T) {
	h := NewProgressoHandler(&fakeMarcarLido{}, nil, nil)
	r := chi.NewRouter()
	r.Put("/progresso/artigos/{id}", h.MarcarArtigo)

	req := httptest.NewRequest(http.MethodPut, "/progresso/artigos/"+uuid.New().String(), strings.NewReader(`{`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400, got %d", w.Code)
	}
}

func TestProgressoHandler_ObterTrilha_OK(t *testing.T) {
	id := uuid.New().String()
	h := NewProgressoHandler(&fakeMarcarLido{}, &fakeProgressoTrilha{
		resp: &dto.ProgressoTrilhaResponse{
			TrilhaID: id, Concluidos: 1, Total: 3, Percentual: 33.333,
			ArtigosConcluidos: []string{"a1"},
		},
	}, nil)
	r := chi.NewRouter()
	r.Get("/progresso/trilhas/{id}", h.ObterTrilha)

	req := comUsuario(httptest.NewRequest(http.MethodGet, "/progresso/trilhas/"+id, nil))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body dto.ProgressoTrilhaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.ArtigosConcluidos) != 1 || body.ArtigosConcluidos[0] != "a1" {
		t.Fatalf("artigos_concluidos: %+v", body.ArtigosConcluidos)
	}
}

func TestProgressoHandler_Continuar_OK(t *testing.T) {
	trilha := "go-basico"
	h := NewProgressoHandler(&fakeMarcarLido{}, &fakeProgressoTrilha{}, &fakeContinuar{
		resp: &dto.ContinuarResponse{Item: &dto.ContinuarItem{
			ArtigoID: uuid.New().String(), ArtigoSlug: "pacotes-em-go", ArtigoTitulo: "Pacotes em Go",
			TrilhaSlug: &trilha,
		}},
	})
	r := chi.NewRouter()
	r.Get("/progresso/continuar", h.Continuar)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, comUsuario(httptest.NewRequest(http.MethodGet, "/progresso/continuar", nil)))

	if w.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"artigo_slug":"pacotes-em-go"`) {
		t.Fatalf("body inesperado: %s", w.Body.String())
	}
}

func TestProgressoHandler_Continuar_SemHistorico(t *testing.T) {
	h := NewProgressoHandler(&fakeMarcarLido{}, &fakeProgressoTrilha{}, &fakeContinuar{resp: &dto.ContinuarResponse{}})
	r := chi.NewRouter()
	r.Get("/progresso/continuar", h.Continuar)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, comUsuario(httptest.NewRequest(http.MethodGet, "/progresso/continuar", nil)))

	if w.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", w.Code)
	}
	if strings.TrimSpace(w.Body.String()) != `{"item":null}` {
		t.Fatalf("esperava item nulo, got %s", w.Body.String())
	}
}

func TestProgressoHandler_Continuar_Erro(t *testing.T) {
	h := NewProgressoHandler(&fakeMarcarLido{}, &fakeProgressoTrilha{}, &fakeContinuar{
		err: domainErros.ErrInvalidArgument("usuario_id inválido", "test", nil),
	})
	r := chi.NewRouter()
	r.Get("/progresso/continuar", h.Continuar)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/progresso/continuar", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400, got %d", w.Code)
	}
}
