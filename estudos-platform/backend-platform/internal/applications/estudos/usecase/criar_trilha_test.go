package usecase

import (
	"context"
	"testing"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type MockTrilhaRepository struct {
	SaveFn           func(ctx context.Context, t *entity.Trilha) error
	FindByIDFn       func(ctx context.Context, id string) (*entity.Trilha, error)
	FindBySlugFn     func(ctx context.Context, slug valueobject.Slug) (*entity.Trilha, error)
	ListPublicadasFn func(ctx context.Context, limit, offset int) ([]*entity.Trilha, error)
	SlugExisteFn     func(ctx context.Context, slug valueobject.Slug) (bool, error)
}

func (m *MockTrilhaRepository) Save(ctx context.Context, t *entity.Trilha) error {
	return m.SaveFn(ctx, t)
}
func (m *MockTrilhaRepository) FindByID(ctx context.Context, id string) (*entity.Trilha, error) {
	return m.FindByIDFn(ctx, id)
}
func (m *MockTrilhaRepository) FindBySlug(ctx context.Context, slug valueobject.Slug) (*entity.Trilha, error) {
	return m.FindBySlugFn(ctx, slug)
}
func (m *MockTrilhaRepository) ListPublicadas(ctx context.Context, limit, offset int) ([]*entity.Trilha, error) {
	return m.ListPublicadasFn(ctx, limit, offset)
}
func (m *MockTrilhaRepository) SlugExiste(ctx context.Context, slug valueobject.Slug) (bool, error) {
	return m.SlugExisteFn(ctx, slug)
}

func TestCriarTrilha_Sucesso(t *testing.T) {
	repo := &MockTrilhaRepository{
		SlugExisteFn: func(ctx context.Context, slug valueobject.Slug) (bool, error) { return false, nil },
		SaveFn:       func(ctx context.Context, tr *entity.Trilha) error { return nil },
	}
	uc := NewCriarTrilha(repo)
	resp, err := uc.Execute(context.Background(), dto.CriarTrilhaRequest{Titulo: "Clean Architecture"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Slug != "clean-architecture" || resp.Publicada {
		t.Fatalf("resposta inesperada: %+v", resp)
	}
}

func TestAdicionarModuloEPublicar(t *testing.T) {
	slug, _ := valueobject.NewSlug("ddd")
	trilha, _ := entity.NovaTrilha(entity.NovaTrilhaInput{Slug: slug, Titulo: "DDD"})
	repo := &MockTrilhaRepository{
		FindByIDFn: func(ctx context.Context, id string) (*entity.Trilha, error) { return trilha, nil },
		SaveFn:     func(ctx context.Context, tr *entity.Trilha) error { return nil },
	}

	addUC := NewAdicionarModulo(repo)
	if _, err := addUC.Execute(context.Background(), trilha.ID().String(), dto.AdicionarModuloRequest{
		Titulo: "Aggregates",
	}); err != nil {
		t.Fatal(err)
	}

	pubUC := NewPublicarTrilha(repo)
	resp, err := pubUC.Execute(context.Background(), trilha.ID().String())
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Publicada || len(resp.Modulos) != 1 {
		t.Fatalf("esperava publicada com 1 módulo: %+v", resp)
	}
}

func TestPublicarTrilha_SemModulos(t *testing.T) {
	slug, _ := valueobject.NewSlug("vazia")
	trilha, _ := entity.NovaTrilha(entity.NovaTrilhaInput{Slug: slug, Titulo: "Vazia"})
	repo := &MockTrilhaRepository{
		FindByIDFn: func(ctx context.Context, id string) (*entity.Trilha, error) { return trilha, nil },
	}
	_, err := NewPublicarTrilha(repo).Execute(context.Background(), trilha.ID().String())
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.InvalidState {
		t.Fatalf("esperava InvalidState, got %#v", err)
	}
}
