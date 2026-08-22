package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type BuscarArtigosUseCase interface {
	Execute(ctx context.Context, q string, limit int) (*dto.BuscaResponse, error)
}

type BuscaHandler struct {
	buscar BuscarArtigosUseCase
}

func NewBuscaHandler(buscar BuscarArtigosUseCase) *BuscaHandler {
	return &BuscaHandler{buscar: buscar}
}

func (h *BuscaHandler) Buscar(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	resp, err := h.buscar.Execute(r.Context(), q, limit)
	if err != nil {
		h.escreverErro(w, err)
		return
	}
	h.escreverJSON(w, http.StatusOK, resp)
}

func (h *BuscaHandler) escreverJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *BuscaHandler) escreverErro(w http.ResponseWriter, err error) {
	var de *domainErros.DomainError
	if !errors.As(err, &de) {
		h.escreverJSON(w, http.StatusInternalServerError, map[string]string{"erro": "erro interno"})
		return
	}
	status := http.StatusInternalServerError
	if de.Kind == domainErros.InvalidArgument {
		status = http.StatusBadRequest
	}
	h.escreverJSON(w, status, map[string]string{"erro": de.Message})
}
