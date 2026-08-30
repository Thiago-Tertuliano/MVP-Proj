package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
	"github.com/thiago-tertuliano/estudos-platform/internal/presentation/http/middleware"
)

type SalvarAnotacaoUseCase interface {
	Execute(ctx context.Context, usuarioID, artigoID string, req dto.SalvarAnotacaoRequest) (*dto.AnotacaoResponse, error)
}

type ObterAnotacaoUseCase interface {
	Execute(ctx context.Context, usuarioID, artigoID string) (*dto.AnotacaoResponse, error)
}

type AnotacaoHandler struct {
	salvar SalvarAnotacaoUseCase
	obter ObterAnotacaoUseCase
	validate *validator.Validate
}

func NewAnotacaoHandler(salvar SalvarAnotacaoUseCase, obter ObterAnotacaoUseCase) *AnotacaoHandler {
	return &AnotacaoHandler{
		salvar: salvar,
		obter: obter,
		validate: validator.New(),
	}
}

func (h *AnotacaoHandler) Salvar(w http.ResponseWriter, r *http.Request) {
	var req dto.SalvarAnotacaoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.escreverJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		h.escreverJSON(w, http.StatusBadRequest, map[string]string{"error": "validation failed"})
		return
	}

	usuarioID, _ := r.Context().Value(middleware.CtxUsuarioID).(string)
	artigoID := chi.URLParam(r, "id")

	resp, err := h.salvar.Execute(r.Context(), usuarioID, artigoID, req)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, resp)
}

func (h *AnotacaoHandler) Obter(w http.ResponseWriter, r *http.Request) {
	usuarioID, _ := r.Context().Value(middleware.CtxUsuarioID).(string)
	artigoID := chi.URLParam(r, "id")
	resp, err := h.obter.Execute(r.Context(), usuarioID, artigoID)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, resp)
}

func (h *AnotacaoHandler) escreverJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
func (h *AnotacaoHandler) escreverErro(w http.ResponseWriter, err error) {
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
	case domainErros.InvalidState:
		status = http.StatusConflict
	}
	h.escreverJSON(w, status, map[string]string{"erro": de.Message})
}