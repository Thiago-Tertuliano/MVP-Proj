package usecase

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type MockArtigoRepository struct {
	SaveFn           func(ctx context.Context, a *entity.Artigo) error
	FindByIDFn       func(ctx context.Context, id string) (*entity.Artigo, error)
	FindBySlugFn     func(ctx context.Context, slug valueobject.Slug) (*entity.Artigo, error)
	ListPublicadosFn func(ctx context.Context, limit, offset int) ([]*entity.Artigo, error)
	SlugExisteFn     func(ctx context.Context, slug valueobject.Slug) (bool, error)
}

func (m *MockArtigoRepository) Save(ctx context.Context, a *entity.Artigo) error {
	return m.SaveFn(ctx, a)
}
func (m *MockArtigoRepository) FindByID(ctx context.Context, id string) (*entity.Artigo, error) {
	return m.FindByIDFn(ctx, id)
}
func (m *MockArtigoRepository) FindBySlug(ctx context.Context, slug valueobject.Slug) (*entity.Artigo, error) {
	return m.FindBySlugFn(ctx, slug)
}
func (m *MockArtigoRepository) ListPublicados(ctx context.Context, limit, offset int) ([]*entity.Artigo, error) {
	return m.ListPublicadosFn(ctx, limit, offset)
}
func (m *MockArtigoRepository) SlugExiste(ctx context.Context, slug valueobject.Slug) (bool, error) {
	return m.SlugExisteFn(ctx, slug)
}

func TestCriarArtigo_Sucesso(t *testing.T) {
	repo := &MockArtigoRepository{
		SlugExisteFn: func(ctx context.Context, slug valueobject.Slug) (bool, error) { return false, nil },
		SaveFn:       func(ctx context.Context, a *entity.Artigo) error { return nil },
	}
	uc := NewCriarArtigo(repo)
	resp, err := uc.Execute(context.Background(), dto.CriarArtigoRequest{
		Titulo:   "Arquitetura Hexagonal",
		Conteudo: json.RawMessage(`{"blocks":[]}`),
	}, uuid.New().String())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Slug != "arquitetura-hexagonal" {
		t.Errorf("slug inesperado: %s", resp.Slug)
	}
	if resp.Status != "rascunho" {
		t.Errorf("status inesperado: %s", resp.Status)
	}
}

func TestCriarArtigo_SlugDuplicado(t *testing.T) {
	repo := &MockArtigoRepository{
		SlugExisteFn: func(ctx context.Context, slug valueobject.Slug) (bool, error) { return true, nil },
	}
	uc := NewCriarArtigo(repo)
	_, err := uc.Execute(context.Background(), dto.CriarArtigoRequest{Titulo: "Teste"}, uuid.New().String())
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.AlreadyExists {
		t.Fatalf("esperava AlreadyExists, got %#v", err)
	}
}

func TestPublicarArtigo_Sucesso(t *testing.T) {
	autor := uuid.New()
	slug, _ := valueobject.NewSlug("pub")
	artigo, _ := entity.NovoArtigo(entity.NovoArtigoInput{
		Titulo:   "Publicável",
		Conteudo: json.RawMessage(`{"ok":true}`),
		Slug:     slug,
		AutorID:  autor,
	})
	repo := &MockArtigoRepository{
		FindByIDFn: func(ctx context.Context, id string) (*entity.Artigo, error) { return artigo, nil },
		SaveFn:     func(ctx context.Context, a *entity.Artigo) error { return nil },
	}
	uc := NewPublicarArtigo(repo)
	resp, err := uc.Execute(context.Background(), artigo.ID().String(), autor.String())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != "publicado" {
		t.Fatalf("esperava publicado, got %s", resp.Status)
	}
}
