package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func novaTrilhaTeste(t *testing.T, publicada bool) *entity.Trilha {
	t.Helper()
	slug, _ := valueobject.NewSlug("go-basico")
	return entity.ReconstruirTrilha(
		uuid.New(), slug, "Go Básico", "", "", 0, publicada, nil, time.Now(), time.Now(),
	)
}

func novoArtigoTeste(slug string, status valueobject.ArtigoStatus) *entity.Artigo {
	return entity.ReconstruirArtigo(
		uuid.New(), valueobject.ReconstructSlug(slug), "Artigo "+slug, "", "",
		nil, nil, uuid.New(), status, nil, time.Now(), time.Now(), nil, nil,
	)
}

func TestListarArtigosPorTrilha_SoPublicados(t *testing.T) {
	trilha := novaTrilhaTeste(t, true)
	uc := NewListarArtigosPorTrilha(
		&MockArtigoRepository{
			ListarPorTrilhaFn: func(ctx context.Context, id uuid.UUID) ([]*entity.Artigo, error) {
				return []*entity.Artigo{
					novoArtigoTeste("publicado-um", valueobject.ArtigoStatusPublicado),
					novoArtigoTeste("rascunho", valueobject.ArtigoStatusRascunho),
					novoArtigoTeste("publicado-dois", valueobject.ArtigoStatusPublicado),
				}, nil
			},
		},
		&MockTrilhaRepository{
			FindBySlugFn: func(ctx context.Context, slug valueobject.Slug) (*entity.Trilha, error) { return trilha, nil },
		},
	)

	resp, err := uc.Execute(context.Background(), "go-basico")
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Itens) != 2 {
		t.Fatalf("esperava 2 artigos publicados, got %d", len(resp.Itens))
	}
	if resp.Itens[0].Slug != "publicado-um" || resp.Itens[1].Slug != "publicado-dois" {
		t.Fatalf("ordem do repositório deve ser preservada: %s, %s", resp.Itens[0].Slug, resp.Itens[1].Slug)
	}
}

func TestListarArtigosPorTrilha_ListaVaziaNaoENula(t *testing.T) {
	uc := NewListarArtigosPorTrilha(
		&MockArtigoRepository{},
		&MockTrilhaRepository{
			FindBySlugFn: func(ctx context.Context, slug valueobject.Slug) (*entity.Trilha, error) {
				return novaTrilhaTeste(t, true), nil
			},
		},
	)
	resp, err := uc.Execute(context.Background(), "go-basico")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Itens == nil {
		t.Fatal("itens deve ser lista vazia ([]), não nula")
	}
}

func TestListarArtigosPorTrilha_TrilhaRascunho404(t *testing.T) {
	uc := NewListarArtigosPorTrilha(
		&MockArtigoRepository{},
		&MockTrilhaRepository{
			FindBySlugFn: func(ctx context.Context, slug valueobject.Slug) (*entity.Trilha, error) {
				return novaTrilhaTeste(t, false), nil
			},
		},
	)
	_, err := uc.Execute(context.Background(), "go-basico")
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.NotFound {
		t.Fatalf("esperava NotFound, got %#v", err)
	}
}

func TestListarArtigosPorTrilha_TrilhaInexistente404(t *testing.T) {
	uc := NewListarArtigosPorTrilha(
		&MockArtigoRepository{},
		&MockTrilhaRepository{
			FindBySlugFn: func(ctx context.Context, slug valueobject.Slug) (*entity.Trilha, error) {
				return nil, domainErros.ErrNotFound("trilha não encontrada", "test", nil)
			},
		},
	)
	_, err := uc.Execute(context.Background(), "nao-existe")
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.NotFound {
		t.Fatalf("esperava NotFound, got %#v", err)
	}
}

func TestListarArtigosPorTrilha_FalhaDeBancoE500(t *testing.T) {
	uc := NewListarArtigosPorTrilha(
		&MockArtigoRepository{},
		&MockTrilhaRepository{
			FindBySlugFn: func(ctx context.Context, slug valueobject.Slug) (*entity.Trilha, error) {
				return nil, context.DeadlineExceeded
			},
		},
	)
	_, err := uc.Execute(context.Background(), "go-basico")
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.Internal {
		t.Fatalf("esperava Internal, got %#v", err)
	}
}
