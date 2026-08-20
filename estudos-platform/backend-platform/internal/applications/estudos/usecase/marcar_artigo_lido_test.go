package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type mockProgressoRepo struct {
	UpsertFn func(ctx context.Context, p repository.ProgressoArtigo) error
}

func (m *mockProgressoRepo) UpsertArtigo(ctx context.Context, p repository.ProgressoArtigo) error {
	return m.UpsertFn(ctx, p)
}

func (m *mockProgressoRepo) CountConcluidosNaTrilha(context.Context, string, string) (int, int, error) {
	return 0, 0, nil
}

func TestMarcarArtigoLido_Sucesso(t *testing.T) {
	artigoID := uuid.New().String()
	userID := uuid.New().String()
	var saved repository.ProgressoArtigo

	uc := NewMarcarArtigoLido(
		&MockArtigoRepository{
			FindByIDFn: func(ctx context.Context, id string) (*entity.Artigo, error) {
				return &entity.Artigo{}, nil
			},
		},
		&mockProgressoRepo{
			UpsertFn: func(ctx context.Context, p repository.ProgressoArtigo) error {
				saved = p
				return nil
			},
		},
	)

	resp, err := uc.Execute(context.Background(), userID, artigoID, true)
	if err != nil {
		t.Fatalf("inesperado: %v", err)
	}
	if !resp.Concluido || resp.ArtigoID != artigoID {
		t.Fatalf("resposta: %+v", resp)
	}
	if saved.ArtigoID != artigoID || saved.UsuarioID != userID || !saved.Concluido {
		t.Fatalf("upsert: %+v", saved)
	}
}

func TestMarcarArtigoLido_ArtigoInexistente(t *testing.T) {
	uc := NewMarcarArtigoLido(
		&MockArtigoRepository{
			FindByIDFn: func(ctx context.Context, id string) (*entity.Artigo, error) {
				return nil, domainErros.ErrNotFound("artigo não encontrado", "test", nil)
			},
		},
		&mockProgressoRepo{
			UpsertFn: func(ctx context.Context, p repository.ProgressoArtigo) error {
				t.Fatal("não deveria upsertar")
				return nil
			},
		},
	)
	_, err := uc.Execute(context.Background(), uuid.New().String(), uuid.New().String(), true)
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.NotFound {
		t.Errorf("esperava NotFound, got %#v", err)
	}
}

func TestMarcarArtigoLido_IDInvalido(t *testing.T) {
	uc := NewMarcarArtigoLido(&MockArtigoRepository{}, &mockProgressoRepo{})
	_, err := uc.Execute(context.Background(), "nao-uuid", uuid.New().String(), true)
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.InvalidArgument {
		t.Errorf("esperava InvalidArgument, got %#v", err)
	}
}
