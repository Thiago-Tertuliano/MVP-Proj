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
	"github.com/thiago-tertuliano/estudos-platform/internal/presentation/http/middleware"
)

type CriarArtigoUseCase interface {
	Execute(ctx context.Context, req dto.CriarArtigoRequest, autorID string) (*dto.ArtigoResponse, error)
}

type ObterArtigoUseCase interface {
	Execute(ctx context.Context, slug string) (*dto.ArtigoResponse, error)
}

type ListarArtigosUseCase interface {
	Execute(ctx context.Context, limit, offset int) (*dto.ListarArtigosResponse, error)
}

type AtualizarArtigoUseCase interface {
	Execute(ctx context.Context, id, autorID string, req dto.AtualizarArtigoRequest) (*dto.ArtigoResponse, error)
}

type PublicarArtigoUseCase interface {
	Execute(ctx context.Context, id, autorID string) (*dto.ArtigoResponse, error)
}

type ArtigoHandler struct {
	criar     CriarArtigoUseCase
	obter     ObterArtigoUseCase
	listar    ListarArtigosUseCase
	atualizar AtualizarArtigoUseCase
	publicar  PublicarArtigoUseCase
	validate  *validator.Validate
}

func NewArtigoHandler(
	criar CriarArtigoUseCase,
	obter ObterArtigoUseCase,
	listar ListarArtigosUseCase,
	atualizar AtualizarArtigoUseCase,
	publicar PublicarArtigoUseCase,
) *ArtigoHandler {
	return &ArtigoHandler{
		criar: criar, obter: obter, listar: listar,
		atualizar: atualizar, publicar: publicar,
		validate: validator.New(),
	}
}

func (h *ArtigoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	var req dto.CriarArtigoRequest
	if !h.decodificarValidar(w, r, &req) {
		return
	}
	autorID, _ := r.Context().Value(middleware.CtxUsuarioID).(string)
	resp, err := h.criar.Execute(r.Context(), req, autorID)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusCreated, resp)
}

func (h *ArtigoHandler) Obter(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	resp, err := h.obter.Execute(r.Context(), slug)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, resp)
}

func (h *ArtigoHandler) Listar(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	resp, err := h.listar.Execute(r.Context(), limit, offset)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, resp)
}

func (h *ArtigoHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	var req dto.AtualizarArtigoRequest
	if !h.decodificarValidar(w, r, &req) {
		return
	}
	autorID, _ := r.Context().Value(middleware.CtxUsuarioID).(string)
	id := chi.URLParam(r, "id")
	resp, err := h.atualizar.Execute(r.Context(), id, autorID, req)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, resp)
}

func (h *ArtigoHandler) Publicar(w http.ResponseWriter, r *http.Request) {
	autorID, _ := r.Context().Value(middleware.CtxUsuarioID).(string)
	id := chi.URLParam(r, "id")
	resp, err := h.publicar.Execute(r.Context(), id, autorID)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, resp)
}

func (h *ArtigoHandler) decodificarValidar(w http.ResponseWriter, r *http.Request, dst any) bool {
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

func (h *ArtigoHandler) escreverJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *ArtigoHandler) escreverErro(w http.ResponseWriter, err error) {
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
