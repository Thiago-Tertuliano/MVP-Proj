package usecase

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

// --- MOCK INTELIGENTE E ISOLADO PARA TRILHAS ---
// Embutir a interface 'repository.TrilhaRepository' faz o Go aceitar essa struct
// sem precisarmos implementar os 10 métodos que ela tem. Sobrescrevemos só o que o teste usa.
type MockTrilhaRepoParaArtigo struct {
	repository.TrilhaRepository 
	FindByIDFn func(ctx context.Context, id string) (*entity.Trilha, error)
}

func (m *MockTrilhaRepoParaArtigo) FindByID(ctx context.Context, id string) (*entity.Trilha, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

// --- TESTES ---

func TestCriarArtigo_Sucesso(t *testing.T) {
	repo := &MockArtigoRepository{
		SlugExisteFn: func(ctx context.Context, slug valueobject.Slug) (bool, error) { return false, nil },
		SaveFn:       func(ctx context.Context, a *entity.Artigo) error { return nil },
	}
	
	// Usando o novo mock isolado
	trilhaRepo := &MockTrilhaRepoParaArtigo{} 

	uc := NewCriarArtigo(repo, trilhaRepo)
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
	
	// Usando o novo mock isolado
	trilhaRepo := &MockTrilhaRepoParaArtigo{} 

	uc := NewCriarArtigo(repo, trilhaRepo)
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
	
	uc := NewPublicarArtigo(repo, nil)
	resp, err := uc.Execute(context.Background(), artigo.ID().String(), autor.String())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != "publicado" {
		t.Fatalf("esperava publicado, got %s", resp.Status)
	}
}