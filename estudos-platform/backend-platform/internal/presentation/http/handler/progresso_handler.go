package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
	"github.com/thiago-tertuliano/estudos-platform/internal/presentation/http/middleware"
)

type MarcarArtigoLidoUseCase interface {
	Execute(ctx context.Context, usuarioID, artigoID string, concluido bool) (*dto.ProgressoArtigoResponse, error)
}

type ObterProgressoTrilhaUseCase interface {
	Execute(ctx context.Context, usuarioID, trilhaID string) (*dto.ProgressoTrilhaResponse, error)
}

type ProgressoHandler struct {
	marcar MarcarArtigoLidoUseCase
	obter  ObterProgressoTrilhaUseCase
}

func NewProgressoHandler(marcar MarcarArtigoLidoUseCase, obter ObterProgressoTrilhaUseCase) *ProgressoHandler {
	return &ProgressoHandler{marcar: marcar, obter: obter}
}

func (h *ProgressoHandler) MarcarArtigo(w http.ResponseWriter, r *http.Request) {
	var req dto.MarcarArtigoLidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.escreverJSON(w, http.StatusBadRequest, map[string]string{"erro": "corpo da requisição inválido"})
		return
	}

	usuarioID, _ := r.Context().Value(middleware.CtxUsuarioID).(string)
	artigoID := chi.URLParam(r, "id")

	resp, err := h.marcar.Execute(r.Context(), usuarioID, artigoID, req.Concluido)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, resp)
}

func (h *ProgressoHandler) ObterTrilha(w http.ResponseWriter, r *http.Request) {
	usuarioID, _ := r.Context().Value(middleware.CtxUsuarioID).(string)
	trilhaID := chi.URLParam(r, "id")
	resp, err := h.obter.Execute(r.Context(), usuarioID, trilhaID)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, resp)
}

func (h *ProgressoHandler) escreverJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *ProgressoHandler) escreverErro(w http.ResponseWriter, err error) {
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
