package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
	"github.com/thiago-tertuliano/estudos-platform/internal/presentation/http/middleware"
)

type fakeListarRM struct{}

func (fakeListarRM) Execute(context.Context) (*dto.ListarRoadmapsResponse, error) {
	return &dto.ListarRoadmapsResponse{Itens: []dto.RoadmapResumoResponse{{Slug: "go"}}}, nil
}

type fakeObterRM struct {
	err error
}

func (f fakeObterRM) Execute(_ context.Context, slug string) (*dto.RoadmapResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &dto.RoadmapResponse{Slug: slug}, nil
}

type fakeAlunoRM struct {
	err       error
	respostas map[string]string
	usuario   string
}

func (f *fakeAlunoRM) ObterProgresso(context.Context, string, string) (*dto.RoadmapProgressoResponse, error) {
	return &dto.RoadmapProgressoResponse{}, f.err
}
func (f *fakeAlunoRM) Concluir(_ context.Context, u, _, noID string) (*dto.ConclusaoNoResponse, error) {
	f.usuario = u
	if f.err != nil {
		return nil, f.err
	}
	return &dto.ConclusaoNoResponse{NoID: noID}, nil
}
func (f *fakeAlunoRM) ResponderQuiz(_ context.Context, _, _, _ string, r map[string]string) (*dto.QuizResultadoResponse, error) {
	f.respostas = r
	return &dto.QuizResultadoResponse{Aprovado: true}, f.err
}
func (f *fakeAlunoRM) ListarProgresso(context.Context, string) (*dto.ListarProgressoRoadmapsResponse, error) {
	return &dto.ListarProgressoRoadmapsResponse{}, nil
}
func (f *fakeAlunoRM) Gamificacao(context.Context, string) (*dto.GamificacaoResponse, error) {
	return &dto.GamificacaoResponse{Nivel: 3}, nil
}

type fakeEditorRM struct {
	err error
}

func (f fakeEditorRM) Meus(context.Context, string) (*dto.ListarRoadmapsResponse, error) {
	return &dto.ListarRoadmapsResponse{}, f.err
}
func (f fakeEditorRM) Obter(context.Context, string, string) (*dto.RoadmapResponse, error) {
	return &dto.RoadmapResponse{}, f.err
}
func (f fakeEditorRM) Criar(_ context.Context, _ string, r dto.CriarRoadmapRequest) (*dto.RoadmapResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &dto.RoadmapResponse{Titulo: r.Titulo, Slug: "novo"}, nil
}
func (f fakeEditorRM) Salvar(context.Context, string, string, dto.SalvarRoadmapRequest) (*dto.RoadmapResponse, error) {
	return &dto.RoadmapResponse{}, f.err
}
func (f fakeEditorRM) Publicar(context.Context, string, string, bool) (*dto.RoadmapResponse, error) {
	return &dto.RoadmapResponse{}, f.err
}

func roteador(h *RoadmapHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/roadmaps", h.Listar)
	r.Get("/roadmaps/{slug}", h.Obter)
	r.Post("/roadmaps", h.Criar)
	r.Put("/roadmaps/{slug}", h.Salvar)
	r.Post("/roadmaps/{slug}/nos/{noId}/concluir", h.ConcluirNo)
	r.Post("/roadmaps/{slug}/nos/{noId}/quiz", h.ResponderQuiz)
	r.Get("/gamificacao/me", h.Gamificacao)
	return r
}

func chamarRM(h http.Handler, metodo, caminho, corpo string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(metodo, caminho, strings.NewReader(corpo))
	req = req.WithContext(context.WithValue(req.Context(), middleware.CtxUsuarioID, "u-1"))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestRoadmapHandler_Sucessos(t *testing.T) {
	aluno := &fakeAlunoRM{}
	h := roteador(NewRoadmapHandler(fakeListarRM{}, fakeObterRM{}, aluno, fakeEditorRM{}))

	if rec := chamarRM(h, "GET", "/roadmaps", ""); rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"slug":"go"`) {
		t.Errorf("listar: %d %s", rec.Code, rec.Body)
	}
	if rec := chamarRM(h, "GET", "/roadmaps/go-basico", ""); rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "go-basico") {
		t.Errorf("obter: %d %s", rec.Code, rec.Body)
	}
	if rec := chamarRM(h, "POST", "/roadmaps", `{"titulo":"Novo"}`); rec.Code != http.StatusCreated {
		t.Errorf("criar: %d %s", rec.Code, rec.Body)
	}
	if rec := chamarRM(h, "PUT", "/roadmaps/x", `{"titulo":"Novo"}`); rec.Code != http.StatusOK {
		t.Errorf("salvar: %d", rec.Code)
	}
	rec := chamarRM(h, "POST", "/roadmaps/go/nos/n1/concluir", "")
	if rec.Code != http.StatusOK || aluno.usuario != "u-1" {
		t.Errorf("concluir: %d usuario=%q", rec.Code, aluno.usuario)
	}
	rec = chamarRM(h, "POST", "/roadmaps/go/nos/n1/quiz", `{"respostas":{"q1":"a"}}`)
	if rec.Code != http.StatusOK || aluno.respostas["q1"] != "a" {
		t.Errorf("quiz: %d %v", rec.Code, aluno.respostas)
	}
	rec = chamarRM(h, "GET", "/gamificacao/me", "")
	var g dto.GamificacaoResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &g)
	if rec.Code != http.StatusOK || g.Nivel != 3 {
		t.Errorf("gamificacao: %d %+v", rec.Code, g)
	}
}

func TestRoadmapHandler_MapeamentoDeErros(t *testing.T) {
	casos := []struct {
		kind   domainErros.Kind
		status int
	}{
		{domainErros.InvalidArgument, http.StatusBadRequest},
		{domainErros.NotFound, http.StatusNotFound},
		{domainErros.InvalidState, http.StatusConflict},
		{domainErros.AlreadyExists, http.StatusConflict},
		{domainErros.Forbidden, http.StatusForbidden},
		{domainErros.Unauthorized, http.StatusUnauthorized},
		{domainErros.Internal, http.StatusInternalServerError},
	}
	for _, c := range casos {
		erro := domainErros.New(c.kind, "mensagem", "op", nil)
		h := roteador(NewRoadmapHandler(fakeListarRM{}, fakeObterRM{err: erro}, &fakeAlunoRM{}, fakeEditorRM{}))
		rec := chamarRM(h, "GET", "/roadmaps/x", "")
		if rec.Code != c.status {
			t.Errorf("%s → %d, esperado %d", c.kind, rec.Code, c.status)
		}
		if c.kind != domainErros.Internal && !strings.Contains(rec.Body.String(), `"erro":"mensagem"`) {
			t.Errorf("%s: corpo %s", c.kind, rec.Body)
		}
		if c.kind == domainErros.Internal && strings.Contains(rec.Body.String(), "mensagem") {
			t.Errorf("erro interno não deve expor a mensagem: %s", rec.Body)
		}
	}
}

func TestRoadmapHandler_CorpoInvalido(t *testing.T) {
	h := roteador(NewRoadmapHandler(fakeListarRM{}, fakeObterRM{}, &fakeAlunoRM{}, fakeEditorRM{}))
	for _, c := range [][2]string{{"POST", "/roadmaps"}, {"PUT", "/roadmaps/x"}, {"POST", "/roadmaps/x/nos/n/quiz"}} {
		if rec := chamarRM(h, c[0], c[1], "{nao-e-json"); rec.Code != http.StatusBadRequest {
			t.Errorf("%s %s → %d", c[0], c[1], rec.Code)
		}
	}
}
