package usecase

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

// ---- fakeRoadmapRepo ----

type fakeRoadmapRepo struct {
	porSlug map[string]*entity.Roadmap
	saves   int
}

func newFakeRoadmapRepo(rms ...*entity.Roadmap) *fakeRoadmapRepo {
	f := &fakeRoadmapRepo{porSlug: map[string]*entity.Roadmap{}}
	for _, r := range rms {
		f.porSlug[r.Slug] = r
	}
	return f
}

func (f *fakeRoadmapRepo) Save(_ context.Context, r *entity.Roadmap) error {
	f.saves++
	cp := *r
	f.porSlug[r.Slug] = &cp
	return nil
}

func (f *fakeRoadmapRepo) FindBySlug(_ context.Context, slug string) (*entity.Roadmap, error) {
	r, ok := f.porSlug[slug]
	if !ok {
		return nil, domainErros.ErrNotFound("roadmap não encontrado", "fake", nil)
	}
	cp := *r
	return &cp, nil
}

func (f *fakeRoadmapRepo) SlugExiste(_ context.Context, slug string) (bool, error) {
	_, ok := f.porSlug[slug]
	return ok, nil
}

func (f *fakeRoadmapRepo) resumos(filtro func(*entity.Roadmap) bool) []repository.RoadmapResumo {
	var out []repository.RoadmapResumo
	for _, r := range f.porSlug {
		if !filtro(r) {
			continue
		}
		chefes := 0
		for _, n := range r.Nos {
			if n.Tipo == entity.NoChefe {
				chefes++
			}
		}
		out = append(out, repository.RoadmapResumo{
			ID: r.ID, Slug: r.Slug, Titulo: r.Titulo, Publicado: r.Publicado, AutorID: r.AutorID,
			TotalNos: len(r.Nos), TotalChefes: chefes, XPTotal: r.XPTotal(), AtualizadoEm: r.UpdatedAt,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Slug < out[j].Slug })
	return out
}

func (f *fakeRoadmapRepo) ListarPublicados(context.Context) ([]repository.RoadmapResumo, error) {
	return f.resumos(func(r *entity.Roadmap) bool { return r.Publicado }), nil
}

func (f *fakeRoadmapRepo) ListarPorAutor(_ context.Context, autorID string) ([]repository.RoadmapResumo, error) {
	return f.resumos(func(r *entity.Roadmap) bool { return r.AutorID == autorID }), nil
}

func (f *fakeRoadmapRepo) ListarNosPorArtigo(_ context.Context, artigoID string) ([]repository.NoVinculado, error) {
	var out []repository.NoVinculado
	for _, r := range f.porSlug {
		if !r.Publicado {
			continue
		}
		for _, n := range r.Nos {
			if n.ArtigoID == artigoID && n.Conclusao == entity.ConclusaoArtigo {
				out = append(out, repository.NoVinculado{RoadmapSlug: r.Slug, NoID: n.ID})
			}
		}
	}
	return out, nil
}

// ---- fakeRoadmapProgRepo ----

type fakeRoadmapProgRepo struct {
	concluidos map[string]map[string]bool // usuario -> noID
	roadmapDe  map[string]string          // noID -> roadmapID
	tipoDe     map[string]entity.TipoNo
	artigos    map[string]bool // artigoID -> concluido (qualquer usuário)
}

func newFakeProg(rms ...*entity.Roadmap) *fakeRoadmapProgRepo {
	f := &fakeRoadmapProgRepo{
		concluidos: map[string]map[string]bool{}, roadmapDe: map[string]string{},
		tipoDe: map[string]entity.TipoNo{}, artigos: map[string]bool{},
	}
	for _, r := range rms {
		for _, n := range r.Nos {
			f.roadmapDe[n.ID] = r.ID
			f.tipoDe[n.ID] = n.Tipo
		}
	}
	return f
}

func (f *fakeRoadmapProgRepo) ListarNosConcluidos(_ context.Context, usuarioID, roadmapID string) ([]string, error) {
	var out []string
	for no := range f.concluidos[usuarioID] {
		if f.roadmapDe[no] == roadmapID {
			out = append(out, no)
		}
	}
	return out, nil
}

func (f *fakeRoadmapProgRepo) MarcarNoConcluido(_ context.Context, usuarioID, _, noID string) (bool, error) {
	if f.concluidos[usuarioID] == nil {
		f.concluidos[usuarioID] = map[string]bool{}
	}
	if f.concluidos[usuarioID][noID] {
		return false, nil
	}
	f.concluidos[usuarioID][noID] = true
	return true, nil
}

func (f *fakeRoadmapProgRepo) EstadoArtigos(_ context.Context, _ string, ids []string) (map[string]bool, error) {
	out := map[string]bool{}
	for _, id := range ids {
		if c, ok := f.artigos[id]; ok {
			out[id] = c
		}
	}
	return out, nil
}

func (f *fakeRoadmapProgRepo) ContarConcluidos(_ context.Context, usuarioID string) (int, int, error) {
	nos, chefes := 0, 0
	for no := range f.concluidos[usuarioID] {
		nos++
		if f.tipoDe[no] == entity.NoChefe {
			chefes++
		}
	}
	return nos, chefes, nil
}

func (f *fakeRoadmapProgRepo) ListarProgressoPorRoadmap(context.Context, string) ([]repository.ProgressoRoadmapItem, error) {
	return nil, nil
}

// ---- fakeGamRepo ----

type fakeGamRepo struct {
	g          map[string]*entity.Gamificacao
	eventos    map[string]bool
	conquistas map[string]map[string]time.Time
}

func newFakeGam() *fakeGamRepo {
	return &fakeGamRepo{g: map[string]*entity.Gamificacao{}, eventos: map[string]bool{}, conquistas: map[string]map[string]time.Time{}}
}

func (f *fakeGamRepo) get(u string) *entity.Gamificacao {
	if f.g[u] == nil {
		f.g[u] = &entity.Gamificacao{}
	}
	return f.g[u]
}

func (f *fakeGamRepo) Obter(_ context.Context, u string) (entity.Gamificacao, error) {
	return *f.get(u), nil
}

func (f *fakeGamRepo) tocar(g *entity.Gamificacao, xp int, agora time.Time) {
	g.XPTotal += xp
	g.StreakAtual, g.StreakMax = entity.AtualizarStreak(g.UltimoDia, agora, g.StreakAtual, g.StreakMax)
	d := entity.DiaCalendario(agora)
	g.UltimoDia = &d
}

func (f *fakeGamRepo) RegistrarAtividade(_ context.Context, u string, agora time.Time) (entity.Gamificacao, error) {
	g := f.get(u)
	f.tocar(g, 0, agora)
	return *g, nil
}

func (f *fakeGamRepo) ConcederXP(_ context.Context, u, origem, ref string, xp int, agora time.Time) (bool, entity.Gamificacao, error) {
	chave := u + "|" + origem + "|" + ref
	g := f.get(u)
	if f.eventos[chave] {
		return false, *g, nil
	}
	f.eventos[chave] = true
	f.tocar(g, xp, agora)
	return true, *g, nil
}

func (f *fakeGamRepo) ListarConquistas(_ context.Context, u string) ([]repository.ConquistaObtida, error) {
	var out []repository.ConquistaObtida
	for c, t := range f.conquistas[u] {
		out = append(out, repository.ConquistaObtida{Codigo: c, Em: t})
	}
	return out, nil
}

func (f *fakeGamRepo) ConcederConquista(_ context.Context, u, codigo string) (bool, error) {
	if f.conquistas[u] == nil {
		f.conquistas[u] = map[string]time.Time{}
	}
	if _, ok := f.conquistas[u][codigo]; ok {
		return false, nil
	}
	f.conquistas[u][codigo] = time.Now()
	return true, nil
}

// ---- builders de cenário ----

type cenario struct {
	rm    *entity.Roadmap
	repo  *fakeRoadmapRepo
	prog  *fakeRoadmapProgRepo
	gam   *fakeGamRepo
	svc   *GamificacaoRoadmap
	aluno *RoadmapAluno
	user  string
	nos   map[string]string // apelido -> id
	agora time.Time
}

// novoCenario monta: a (manual, 10xp) → b (manual, 10xp) → chefe (quiz, 100xp); c solta (artigo).
func novoCenario(publicado bool) *cenario {
	ids := map[string]string{}
	for _, k := range []string{"a", "b", "chefe", "c"} {
		ids[k] = uuid.NewString()
	}
	artigo := uuid.NewString()
	rm, _ := entity.NovoRoadmap("go-basico", "Go básico", "", "", uuid.NewString())
	rm.Nos = []entity.RoadmapNo{
		{ID: ids["a"], Tipo: entity.NoTopico, Titulo: "Variáveis", XP: 10, Conclusao: entity.ConclusaoManual},
		{ID: ids["b"], Tipo: entity.NoTopico, Titulo: "Funções", XP: 10, Conclusao: entity.ConclusaoManual},
		{ID: ids["chefe"], Tipo: entity.NoChefe, Titulo: "Chefe", XP: 100, Conclusao: entity.ConclusaoQuiz, Quiz: &entity.Quiz{Questoes: []entity.QuestaoQuiz{
			{ID: "q1", Enunciado: "1?", Opcoes: []entity.OpcaoQuiz{{ID: "x", Texto: "x"}, {ID: "y", Texto: "y"}}, Correta: "x", Explicacao: "porque x"},
			{ID: "q2", Enunciado: "2?", Opcoes: []entity.OpcaoQuiz{{ID: "x", Texto: "x"}, {ID: "y", Texto: "y"}}, Correta: "y"},
		}}},
		{ID: ids["c"], Tipo: entity.NoTopico, Titulo: "Pacotes", XP: 20, Conclusao: entity.ConclusaoArtigo, ArtigoID: artigo},
	}
	rm.Arestas = []entity.RoadmapAresta{
		{OrigemID: ids["a"], DestinoID: ids["b"], Tipo: entity.ArestaRequer},
		{OrigemID: ids["b"], DestinoID: ids["chefe"], Tipo: entity.ArestaRequer},
	}
	rm.Publicado = publicado

	repo := newFakeRoadmapRepo(rm)
	prog := newFakeProg(rm)
	gam := newFakeGam()
	agora := time.Date(2026, 10, 8, 15, 0, 0, 0, time.UTC)
	svc := NewGamificacaoRoadmap(repo, prog, gam).ComRelogio(func() time.Time { return agora })
	ids["artigo"] = artigo
	return &cenario{
		rm: rm, repo: repo, prog: prog, gam: gam, svc: svc,
		aluno: NewRoadmapAluno(repo, prog, svc), user: uuid.NewString(), nos: ids, agora: agora,
	}
}
