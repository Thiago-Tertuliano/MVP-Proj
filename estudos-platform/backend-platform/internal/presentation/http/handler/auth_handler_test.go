package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
	"github.com/thiago-tertuliano/estudos-platform/internal/presentation/http/middleware"
)

// ---- fakes das interfaces do AuthHandler ----

type fakeRegistrar struct {
	resp *dto.AuthResponse
	err  error
}

func (f *fakeRegistrar) Execute(ctx context.Context, req dto.RegistrarRequest) (*dto.AuthResponse, error) {
	return f.resp, f.err
}

type fakeLogin struct {
	resp *dto.AuthResponse
	err  error
}

func (f *fakeLogin) Execute(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	return f.resp, f.err
}

type fakeRefresh struct {
	resp *dto.AuthResponse
	err  error
}

func (f *fakeRefresh) Execute(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	return f.resp, f.err
}

type fakePerfil struct {
	resp *dto.UsuarioResponse
	err  error
}

func (f *fakePerfil) Execute(ctx context.Context, usuarioID string) (*dto.UsuarioResponse, error) {
	return f.resp, f.err
}

func authResponseFake() *dto.AuthResponse {
	return &dto.AuthResponse{
		Tokens:  dto.TokenResponse{AccessToken: "access", RefreshToken: "refresh", ExpiracaoEm: 123},
		Usuario: dto.UsuarioResponse{ID: "u-1", Nome: "Thiago", Email: "t@ex.com"},
	}
}

func executarHandler(fn http.HandlerFunc, method, target, corpo string, ctx context.Context) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, strings.NewReader(corpo))
	req.Header.Set("Content-Type", "application/json")
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	w := httptest.NewRecorder()
	fn(w, req)
	return w
}

// ---- testes ----

func TestRegistrar_Sucesso(t *testing.T) {
	h := NewAuthHandler(&fakeRegistrar{resp: authResponseFake()}, &fakeLogin{}, &fakeRefresh{}, &fakePerfil{})
	w := executarHandler(h.Registrar, http.MethodPost, "/auth/registrar", `{"nome":"Thiago","email":"t@ex.com","senha":"senha123"}`, nil)
	if w.Code != http.StatusCreated {
		t.Fatalf("esperava 201, got %d", w.Code)
	}
}

func TestRegistrar_CorpoInvalido(t *testing.T) {
	h := NewAuthHandler(&fakeRegistrar{}, &fakeLogin{}, &fakeRefresh{}, &fakePerfil{})
	w := executarHandler(h.Registrar, http.MethodPost, "/auth/registrar", `{nao-json`, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400, got %d", w.Code)
	}
}

func TestRegistrar_Validacao(t *testing.T) {
	h := NewAuthHandler(&fakeRegistrar{}, &fakeLogin{}, &fakeRefresh{}, &fakePerfil{})
	w := executarHandler(h.Registrar, http.MethodPost, "/auth/registrar", `{"nome":"T","email":"x","senha":"1"}`, nil)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("esperava 422, got %d", w.Code)
	}
}

func TestRegistrar_EmailJaCadastrado(t *testing.T) {
	erro := domainErros.ErrAlreadyExists("e-mail já cadastrado", "test", nil)
	h := NewAuthHandler(&fakeRegistrar{err: erro}, &fakeLogin{}, &fakeRefresh{}, &fakePerfil{})
	w := executarHandler(h.Registrar, http.MethodPost, "/auth/registrar", `{"nome":"Thiago","email":"t@ex.com","senha":"senha123"}`, nil)
	if w.Code != http.StatusConflict {
		t.Fatalf("esperava 409, got %d", w.Code)
	}
}

func TestLogin_Sucesso(t *testing.T) {
	h := NewAuthHandler(&fakeRegistrar{}, &fakeLogin{resp: authResponseFake()}, &fakeRefresh{}, &fakePerfil{})
	w := executarHandler(h.Login, http.MethodPost, "/auth/login", `{"email":"t@ex.com","senha":"senha123"}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", w.Code)
	}
	var resp dto.AuthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("resposta não é JSON: %v", err)
	}
	if resp.Tokens.AccessToken != "access" {
		t.Error("access token ausente na resposta")
	}
}

func TestLogin_CredenciaisInvalidas(t *testing.T) {
	erro := domainErros.ErrUnauthorized("credenciais inválidas", "test", nil)
	h := NewAuthHandler(&fakeRegistrar{}, &fakeLogin{err: erro}, &fakeRefresh{}, &fakePerfil{})
	w := executarHandler(h.Login, http.MethodPost, "/auth/login", `{"email":"t@ex.com","senha":"errada"}`, nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("esperava 401, got %d", w.Code)
	}
}

func TestRefresh_Sucesso(t *testing.T) {
	h := NewAuthHandler(&fakeRegistrar{}, &fakeLogin{}, &fakeRefresh{resp: authResponseFake()}, &fakePerfil{})
	w := executarHandler(h.Refresh, http.MethodPost, "/auth/refresh", `{"refresh_token":"token"}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", w.Code)
	}
}

func TestRefresh_TokenInvalido(t *testing.T) {
	erro := domainErros.ErrUnauthorized("refresh token inválido", "test", nil)
	h := NewAuthHandler(&fakeRegistrar{}, &fakeLogin{}, &fakeRefresh{err: erro}, &fakePerfil{})
	w := executarHandler(h.Refresh, http.MethodPost, "/auth/refresh", `{"refresh_token":"token"}`, nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("esperava 401, got %d", w.Code)
	}
}

func TestMe(t *testing.T) {
	perfil := &fakePerfil{resp: &dto.UsuarioResponse{ID: "u-1", Nome: "Thiago", Email: "t@ex.com"}}
	h := NewAuthHandler(&fakeRegistrar{}, &fakeLogin{}, &fakeRefresh{}, perfil)
	ctx := context.WithValue(context.Background(), middleware.CtxUsuarioID, "u-1")
	w := executarHandler(h.Me, http.MethodGet, "/auth/me", "", ctx)
	if w.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", w.Code)
	}
	var resp dto.UsuarioResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("resposta não é JSON: %v", err)
	}
	if resp.Nome != "Thiago" || resp.Email != "t@ex.com" {
		t.Errorf("perfil incorreto: %+v", resp)
	}
}

func TestMe_Erro(t *testing.T) {
	erro := domainErros.ErrNotFound("usuário não encontrado", "test", nil)
	h := NewAuthHandler(&fakeRegistrar{}, &fakeLogin{}, &fakeRefresh{}, &fakePerfil{err: erro})
	ctx := context.WithValue(context.Background(), middleware.CtxUsuarioID, "u-inexistente")
	w := executarHandler(h.Me, http.MethodGet, "/auth/me", "", ctx)
	if w.Code != http.StatusNotFound {
		t.Fatalf("esperava 404, got %d", w.Code)
	}
}
