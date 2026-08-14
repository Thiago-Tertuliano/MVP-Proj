package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type CriarTrilhaUseCase interface {
	Execute(ctx context.Context, req dto.CriarTrilhaRequest) (*dto.TrilhaResponse, error)
}

type ObterTrilhaUseCase interface {
	Execute(ctx context.Context, slug string) (*dto.TrilhaResponse, error)
}

type ListarTrilhasUseCase interface {
	Execute(ctx context.Context, limit, offset int) (*dto.ListarTrilhasResponse, error)
}

type AdicionarModuloUseCase interface {
	Execute(ctx context.Context, trilhaID string, req dto.AdicionarModuloRequest) (*dto.TrilhaResponse, error)
}

type PublicarTrilhaUseCase interface {
	Execute(ctx context.Context, id string) (*dto.TrilhaResponse, error)
}

type TrilhaHandler struct {
	criar     CriarTrilhaUseCase
	obter     ObterTrilhaUseCase
	listar    ListarTrilhasUseCase
	adicionar AdicionarModuloUseCase
	publicar  PublicarTrilhaUseCase
	validate  *validator.Validate
}

func NewTrilhaHandler(
	criar CriarTrilhaUseCase,
	obter ObterTrilhaUseCase,
	listar ListarTrilhasUseCase,
	adicionar AdicionarModuloUseCase,
	publicar PublicarTrilhaUseCase,
) *TrilhaHandler {
	return &TrilhaHandler{
		criar: criar, obter: obter, listar: listar,
		adicionar: adicionar, publicar: publicar,
		validate: validator.New(),
	}
}

func (h *TrilhaHandler) Criar(w http.ResponseWriter, r *http.Request) {
	var req dto.CriarTrilhaRequest
	if !h.decodificarValidar(w, r, &req) {
		return
	}
	resp, err := h.criar.Execute(r.Context(), req)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusCreated, resp)
}

func (h *TrilhaHandler) Obter(w http.ResponseWriter, r *http.Request) {
	resp, err := h.obter.Execute(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, resp)
}

func (h *TrilhaHandler) Listar(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	resp, err := h.listar.Execute(r.Context(), limit, offset)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, resp)
}

func (h *TrilhaHandler) AdicionarModulo(w http.ResponseWriter, r *http.Request) {
	var req dto.AdicionarModuloRequest
	if !h.decodificarValidar(w, r, &req) {
		return
	}
	resp, err := h.adicionar.Execute(r.Context(), chi.URLParam(r, "id"), req)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusCreated, resp)
}

func (h *TrilhaHandler) Publicar(w http.ResponseWriter, r *http.Request) {
	resp, err := h.publicar.Execute(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, resp)
}

func (h *TrilhaHandler) decodificarValidar(w http.ResponseWriter, r *http.Request, dst any) bool {
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

func (h *TrilhaHandler) escreverJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *TrilhaHandler) escreverErro(w http.ResponseWriter, err error) {
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
