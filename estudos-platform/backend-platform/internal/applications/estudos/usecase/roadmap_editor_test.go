package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func TestRoadmapEditor_CriarGeraSlugUnico(t *testing.T) {
	repo := newFakeRoadmapRepo()
	uc := NewRoadmapEditor(repo)
	autor := uuid.NewString()
	ctx := context.Background()

	a, err := uc.Criar(ctx, autor, dto.CriarRoadmapRequest{Titulo: "Go do Zero"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := uc.Criar(ctx, autor, dto.CriarRoadmapRequest{Titulo: "Go do Zero"})
	if err != nil {
		t.Fatal(err)
	}
	if a.Slug != "go-do-zero" || b.Slug != "go-do-zero-2" {
		t.Fatalf("slugs: %s / %s", a.Slug, b.Slug)
	}
	if a.Publicado {
		t.Fatal("deveria nascer como rascunho")
	}
	meus, _ := uc.Meus(ctx, autor)
	if len(meus.Itens) != 2 {
		t.Fatalf("meus = %d", len(meus.Itens))
	}
}

func TestRoadmapEditor_SlugReservado(t *testing.T) {
	uc := NewRoadmapEditor(newFakeRoadmapRepo())
	r, err := uc.Criar(context.Background(), uuid.NewString(), dto.CriarRoadmapRequest{Titulo: "Meus"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Slug == "meus" {
		t.Fatal("slug reservado não pode ser usado")
	}
}

func TestRoadmapEditor_CriarTituloInvalido(t *testing.T) {
	uc := NewRoadmapEditor(newFakeRoadmapRepo())
	_, err := uc.Criar(context.Background(), uuid.NewString(), dto.CriarRoadmapRequest{Titulo: "ab"})
	if kind(t, err) != domainErros.InvalidArgument {
		t.Fatalf("got %v", err)
	}
}

func TestRoadmapEditor_SoOAutorEdita(t *testing.T) {
	repo := newFakeRoadmapRepo()
	uc := NewRoadmapEditor(repo)
	dono, outro := uuid.NewString(), uuid.NewString()
	ctx := context.Background()
	r, _ := uc.Criar(ctx, dono, dto.CriarRoadmapRequest{Titulo: "Meu mapa"})

	if _, err := uc.Obter(ctx, outro, r.Slug); kind(t, err) != domainErros.Forbidden {
		t.Fatalf("obter: %v", err)
	}
	if _, err := uc.Salvar(ctx, outro, r.Slug, dto.SalvarRoadmapRequest{Titulo: "Hack"}); kind(t, err) != domainErros.Forbidden {
		t.Fatalf("salvar: %v", err)
	}
	if _, err := uc.Publicar(ctx, outro, r.Slug, true); kind(t, err) != domainErros.Forbidden {
		t.Fatalf("publicar: %v", err)
	}
	if _, err := uc.Obter(ctx, dono, "nao-existe"); kind(t, err) != domainErros.NotFound {
		t.Fatalf("inexistente: %v", err)
	}
}

func TestRoadmapEditor_RascunhoAceitaIncompletoMasPublicarExige(t *testing.T) {
	uc := NewRoadmapEditor(newFakeRoadmapRepo())
	autor := uuid.NewString()
	ctx := context.Background()
	r, _ := uc.Criar(ctx, autor, dto.CriarRoadmapRequest{Titulo: "Com chefe"})

	chefeID := uuid.NewString()
	req := dto.SalvarRoadmapRequest{
		Titulo: "Com chefe",
		Nos:    []dto.RoadmapNoDTO{{ID: chefeID, Tipo: "chefe", Titulo: "Chefe final", XP: 100}}, // sem quiz ainda
	}
	salvo, err := uc.Salvar(ctx, autor, r.Slug, req)
	if err != nil {
		t.Fatalf("rascunho incompleto deveria salvar: %v", err)
	}
	if salvo.Nos[0].Conclusao != "quiz" {
		t.Fatalf("chefe deveria assumir conclusão por quiz: %+v", salvo.Nos[0])
	}

	if _, err := uc.Publicar(ctx, autor, r.Slug, true); kind(t, err) != domainErros.InvalidArgument {
		t.Fatalf("publicar sem quiz deveria falhar: %v", err)
	}

	req.Nos[0].Quiz = &dto.QuizDTO{Questoes: []dto.QuestaoQuizDTO{{
		ID: "q1", Enunciado: "?", Correta: "a", Opcoes: []dto.OpcaoQuizDTO{{ID: "a", Texto: "A"}, {ID: "b", Texto: "B"}},
	}}}
	if _, err := uc.Salvar(ctx, autor, r.Slug, req); err != nil {
		t.Fatal(err)
	}
	pub, err := uc.Publicar(ctx, autor, r.Slug, true)
	if err != nil || !pub.Publicado {
		t.Fatalf("publicar: %v %+v", err, pub)
	}

	// Publicado: salvar de novo exige grafo íntegro.
	req.Nos[0].Quiz = nil
	if _, err := uc.Salvar(ctx, autor, r.Slug, req); kind(t, err) != domainErros.InvalidArgument {
		t.Fatalf("salvar publicado inválido deveria falhar: %v", err)
	}

	desp, err := uc.Publicar(ctx, autor, r.Slug, false)
	if err != nil || desp.Publicado {
		t.Fatalf("despublicar: %v %+v", err, desp)
	}
}

func TestRoadmapEditor_SalvarRejeitaCiclo(t *testing.T) {
	uc := NewRoadmapEditor(newFakeRoadmapRepo())
	autor := uuid.NewString()
	ctx := context.Background()
	r, _ := uc.Criar(ctx, autor, dto.CriarRoadmapRequest{Titulo: "Ciclo"})
	a, b := uuid.NewString(), uuid.NewString()
	_, err := uc.Salvar(ctx, autor, r.Slug, dto.SalvarRoadmapRequest{
		Titulo: "Ciclo",
		Nos: []dto.RoadmapNoDTO{
			{ID: a, Titulo: "A", XP: 10}, {ID: b, Titulo: "B", XP: 10},
		},
		Arestas: []dto.RoadmapArestaDTO{{Origem: a, Destino: b}, {Origem: b, Destino: a}},
	})
	if kind(t, err) != domainErros.InvalidArgument {
		t.Fatalf("got %v", err)
	}
}

func TestRoadmapEditor_ObterIncluiGabaritoMasAlunoNao(t *testing.T) {
	uc := NewRoadmapEditor(newFakeRoadmapRepo())
	autor := uuid.NewString()
	ctx := context.Background()
	r, _ := uc.Criar(ctx, autor, dto.CriarRoadmapRequest{Titulo: "Gabarito"})
	_, err := uc.Salvar(ctx, autor, r.Slug, dto.SalvarRoadmapRequest{
		Titulo: "Gabarito",
		Nos: []dto.RoadmapNoDTO{{ID: uuid.NewString(), Tipo: "chefe", Titulo: "Chefe", XP: 50, Quiz: &dto.QuizDTO{Questoes: []dto.QuestaoQuizDTO{{
			ID: "q", Enunciado: "?", Correta: "a", Explicacao: "e", Opcoes: []dto.OpcaoQuizDTO{{ID: "a", Texto: "A"}, {ID: "b", Texto: "B"}},
		}}}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	ed, err := uc.Obter(ctx, autor, r.Slug)
	if err != nil {
		t.Fatal(err)
	}
	if ed.Nos[0].Quiz.Questoes[0].Correta != "a" {
		t.Fatal("editor deveria ver o gabarito")
	}
}
