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

// maxCorpoRoadmap limita o payload do editor (200 nós com quiz cabem com folga).
const maxCorpoRoadmap = 2 << 20 // 2 MiB

type ListarRoadmapsUseCase interface {
	Execute(ctx context.Context) (*dto.ListarRoadmapsResponse, error)
}

type ObterRoadmapUseCase interface {
	Execute(ctx context.Context, slug string) (*dto.RoadmapResponse, error)
}

type RoadmapAlunoUseCase interface {
	ObterProgresso(ctx context.Context, usuarioID, slug string) (*dto.RoadmapProgressoResponse, error)
	Concluir(ctx context.Context, usuarioID, slug, noID string) (*dto.ConclusaoNoResponse, error)
	ResponderQuiz(ctx context.Context, usuarioID, slug, noID string, respostas map[string]string) (*dto.QuizResultadoResponse, error)
	ListarProgresso(ctx context.Context, usuarioID string) (*dto.ListarProgressoRoadmapsResponse, error)
	Gamificacao(ctx context.Context, usuarioID string) (*dto.GamificacaoResponse, error)
}

type RoadmapEditorUseCase interface {
	Meus(ctx context.Context, autorID string) (*dto.ListarRoadmapsResponse, error)
	Obter(ctx context.Context, autorID, slug string) (*dto.RoadmapResponse, error)
	Criar(ctx context.Context, autorID string, req dto.CriarRoadmapRequest) (*dto.RoadmapResponse, error)
	Salvar(ctx context.Context, autorID, slug string, req dto.SalvarRoadmapRequest) (*dto.RoadmapResponse, error)
	Publicar(ctx context.Context, autorID, slug string, publicar bool) (*dto.RoadmapResponse, error)
}

type RoadmapHandler struct {
	listar ListarRoadmapsUseCase
	obter  ObterRoadmapUseCase
	aluno  RoadmapAlunoUseCase
	editor RoadmapEditorUseCase
}

func NewRoadmapHandler(l ListarRoadmapsUseCase, o ObterRoadmapUseCase, a RoadmapAlunoUseCase, e RoadmapEditorUseCase) *RoadmapHandler {
	return &RoadmapHandler{listar: l, obter: o, aluno: a, editor: e}
}

func usuarioDe(r *http.Request) string {
	id, _ := r.Context().Value(middleware.CtxUsuarioID).(string)
	return id
}

// ---- Público ----

func (h *RoadmapHandler) Listar(w http.ResponseWriter, r *http.Request) {
	resp, err := h.listar.Execute(r.Context())
	h.responder(w, http.StatusOK, resp, err)
}

func (h *RoadmapHandler) Obter(w http.ResponseWriter, r *http.Request) {
	resp, err := h.obter.Execute(r.Context(), chi.URLParam(r, "slug"))
	h.responder(w, http.StatusOK, resp, err)
}

// ---- Aluno ----

func (h *RoadmapHandler) Progresso(w http.ResponseWriter, r *http.Request) {
	resp, err := h.aluno.ObterProgresso(r.Context(), usuarioDe(r), chi.URLParam(r, "slug"))
	h.responder(w, http.StatusOK, resp, err)
}

func (h *RoadmapHandler) ConcluirNo(w http.ResponseWriter, r *http.Request) {
	resp, err := h.aluno.Concluir(r.Context(), usuarioDe(r), chi.URLParam(r, "slug"), chi.URLParam(r, "noId"))
	h.responder(w, http.StatusOK, resp, err)
}

func (h *RoadmapHandler) ResponderQuiz(w http.ResponseWriter, r *http.Request) {
	var req dto.ResponderQuizRequest
	if !h.decodificar(w, r, &req) {
		return
	}
	resp, err := h.aluno.ResponderQuiz(r.Context(), usuarioDe(r), chi.URLParam(r, "slug"), chi.URLParam(r, "noId"), req.Respostas)
	h.responder(w, http.StatusOK, resp, err)
}

func (h *RoadmapHandler) ProgressoCatalogo(w http.ResponseWriter, r *http.Request) {
	resp, err := h.aluno.ListarProgresso(r.Context(), usuarioDe(r))
	h.responder(w, http.StatusOK, resp, err)
}

func (h *RoadmapHandler) Gamificacao(w http.ResponseWriter, r *http.Request) {
	resp, err := h.aluno.Gamificacao(r.Context(), usuarioDe(r))
	h.responder(w, http.StatusOK, resp, err)
}

// ---- Editor ----

func (h *RoadmapHandler) Meus(w http.ResponseWriter, r *http.Request) {
	resp, err := h.editor.Meus(r.Context(), usuarioDe(r))
	h.responder(w, http.StatusOK, resp, err)
}

func (h *RoadmapHandler) ObterParaEdicao(w http.ResponseWriter, r *http.Request) {
	resp, err := h.editor.Obter(r.Context(), usuarioDe(r), chi.URLParam(r, "slug"))
	h.responder(w, http.StatusOK, resp, err)
}

func (h *RoadmapHandler) Criar(w http.ResponseWriter, r *http.Request) {
	var req dto.CriarRoadmapRequest
	if !h.decodificar(w, r, &req) {
		return
	}
	resp, err := h.editor.Criar(r.Context(), usuarioDe(r), req)
	h.responder(w, http.StatusCreated, resp, err)
}

func (h *RoadmapHandler) Salvar(w http.ResponseWriter, r *http.Request) {
	var req dto.SalvarRoadmapRequest
	if !h.decodificar(w, r, &req) {
		return
	}
	resp, err := h.editor.Salvar(r.Context(), usuarioDe(r), chi.URLParam(r, "slug"), req)
	h.responder(w, http.StatusOK, resp, err)
}

func (h *RoadmapHandler) Publicar(w http.ResponseWriter, r *http.Request) {
	resp, err := h.editor.Publicar(r.Context(), usuarioDe(r), chi.URLParam(r, "slug"), true)
	h.responder(w, http.StatusOK, resp, err)
}

func (h *RoadmapHandler) Despublicar(w http.ResponseWriter, r *http.Request) {
	resp, err := h.editor.Publicar(r.Context(), usuarioDe(r), chi.URLParam(r, "slug"), false)
	h.responder(w, http.StatusOK, resp, err)
}

// ---- Infra ----

func (h *RoadmapHandler) decodificar(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxCorpoRoadmap)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		escreverJSONRoadmap(w, http.StatusBadRequest, map[string]string{"erro": "corpo da requisição inválido"})
		return false
	}
	return true
}

func (h *RoadmapHandler) responder(w http.ResponseWriter, status int, payload any, err error) {
	if err != nil {
		escreverErroRoadmap(w, err)
		return
	}
	escreverJSONRoadmap(w, status, payload)
}

func escreverJSONRoadmap(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func escreverErroRoadmap(w http.ResponseWriter, err error) {
	var de *domainErros.DomainError
	if !errors.As(err, &de) {
		escreverJSONRoadmap(w, http.StatusInternalServerError, map[string]string{"erro": "erro interno"})
		return
	}
	status := http.StatusInternalServerError
	switch de.Kind {
	case domainErros.InvalidArgument:
		status = http.StatusBadRequest
	case domainErros.NotFound:
		status = http.StatusNotFound
	case domainErros.InvalidState, domainErros.AlreadyExists:
		status = http.StatusConflict
	case domainErros.Unauthorized:
		status = http.StatusUnauthorized
	case domainErros.Forbidden:
		status = http.StatusForbidden
	}
	msg := de.Message
	if status == http.StatusInternalServerError {
		msg = "erro interno"
	}
	escreverJSONRoadmap(w, status, map[string]string{"erro": msg})
}
