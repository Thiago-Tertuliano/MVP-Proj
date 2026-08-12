package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
	"github.com/thiago-tertuliano/estudos-platform/internal/presentation/http/middleware"
)

// Interfaces locais: dependemos de contratos, não de implementações concretas.
// As structs de use case satisfazem essas interfaces sem nenhuma mudança nelas.
type RegistrarUseCase interface {
	Execute(ctx context.Context, req dto.RegistrarRequest) (*dto.AuthResponse, error)
}

type LoginUseCase interface {
	Execute(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
}

type RefreshUseCase interface {
	Execute(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
}

type ObterPerfilUseCase interface {
	Execute(ctx context.Context, usuarioID string) (*dto.UsuarioResponse, error)
}

type LogoutUseCase interface {
	Execute(ctx context.Context, usuarioID string) error
}

// AuthHandler é LEVE: só faz binding, validação de borda e delega para o use case.
type AuthHandler struct {
	registrar RegistrarUseCase
	login     LoginUseCase
	refresh   RefreshUseCase
	perfil    ObterPerfilUseCase
	logout    LogoutUseCase
	validate  *validator.Validate
}

func NewAuthHandler(
	registrar RegistrarUseCase,
	login LoginUseCase,
	refresh RefreshUseCase,
	perfil ObterPerfilUseCase,
	logout LogoutUseCase,
) *AuthHandler {
	return &AuthHandler{
		registrar: registrar,
		login:     login,
		refresh:   refresh,
		perfil:    perfil,
		logout:    logout,
		validate:  validator.New(),
	}
}

func (h *AuthHandler) Registrar(w http.ResponseWriter, r *http.Request) {
	var req dto.RegistrarRequest
	if !h.decodificarValidar(w, r, &req) {
		return
	}
	resp, err := h.registrar.Execute(r.Context(), req)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if !h.decodificarValidar(w, r, &req) {
		return
	}
	resp, err := h.login.Execute(r.Context(), req)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if !h.decodificarValidar(w, r, &req) {
		return
	}
	resp, err := h.refresh.Execute(r.Context(), req.RefreshToken)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, resp)
}

// Me retorna o perfil completo do usuário logado (ID vem do token validado no
// middleware; nome/email são buscados no banco pelo use case).
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	usuarioID, _ := r.Context().Value(middleware.CtxUsuarioID).(string)
	resp, err := h.perfil.Execute(r.Context(), usuarioID)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, resp)
}

// Logout revoga todos os refresh tokens do usuário autenticado.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	usuarioID, _ := r.Context().Value(middleware.CtxUsuarioID).(string)
	if err := h.logout.Execute(r.Context(), usuarioID); err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, map[string]string{"mensagem": "logout efetuado"})
}

// ---- helpers ----

func (h *AuthHandler) decodificarValidar(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		h.escreverJSON(w, http.StatusBadRequest, map[string]string{"erro": "corpo da requisição inválido"})
		return false
	}
	if err := h.validate.Struct(dst); err != nil {
		h.escreverJSON(w, http.StatusUnprocessableEntity, map[string]string{"erro": err.Error()})
		return false
	}
	return true
}

func (h *AuthHandler) escreverJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// escreverErro mapeia DomainError → HTTP status.
func (h *AuthHandler) escreverErro(w http.ResponseWriter, err error) {
	var de *domainErros.DomainError
	if !errors.As(err, &de) {
		h.escreverJSON(w, http.StatusInternalServerError, map[string]string{"erro": "erro interno"})
		return
	}

	status := http.StatusInternalServerError
	switch de.Kind {
	case domainErros.InvalidArgument:
		status = http.StatusBadRequest
	case domainErros.NotFound:
		status = http.StatusNotFound
	case domainErros.AlreadyExists:
		status = http.StatusConflict
	case domainErros.Unauthorized:
		status = http.StatusUnauthorized
	case domainErros.Forbidden:
		status = http.StatusForbidden
	case domainErros.InvalidState:
		status = http.StatusConflict
	}

	h.escreverJSON(w, status, map[string]string{"erro": de.Message})
}
