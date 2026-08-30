package usecase

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type mockAnotacaoRepo struct {
	UpsertFn func(ctx context.Context, a repository.Anotacao) error
	FindFn   func(ctx context.Context, usuarioID, artigoID string) (*repository.Anotacao, error)
}

func (m *mockAnotacaoRepo) Upsert(ctx context.Context, a repository.Anotacao) error {
	return m.UpsertFn(ctx, a)
}

func (m *mockAnotacaoRepo) FindByUsuarioEArtigo(ctx context.Context, usuarioID, artigoID string) (*repository.Anotacao, error) {
	if m.FindFn != nil {
		return m.FindFn(ctx, usuarioID, artigoID)
	}
	return &repository.Anotacao{UsuarioID: usuarioID, ArtigoID: artigoID, Conteudo: json.RawMessage(`{}`)}, nil
}

func TestSalvarAnotacao_Sucesso(t *testing.T) {
	artigoID := uuid.New().String()
	userID := uuid.New().String()
	var saved repository.Anotacao

	uc := NewSalvarAnotacao(
		&MockArtigoRepository{
			FindByIDFn: func(ctx context.Context, id string) (*entity.Artigo, error) {
				return &entity.Artigo{}, nil
			},
		},
		&mockAnotacaoRepo{
			UpsertFn: func(ctx context.Context, a repository.Anotacao) error {
				saved = a
				return nil
			},
		},
	)

	conteudo := json.RawMessage(`{"notes":[{"texto":"lembrar do aggregate"}]}`)
	resp, err := uc.Execute(context.Background(), userID, artigoID, dto.SalvarAnotacaoRequest{Conteudo: conteudo})
	if err != nil {
		t.Fatalf("inesperado: %v", err)
	}
	if resp.ArtigoID != artigoID {
		t.Fatalf("resp: %+v", resp)
	}
	if saved.UsuarioID != userID || saved.ArtigoID != artigoID {
		t.Fatalf("upsert: %+v", saved)
	}
}

func TestSalvarAnotacao_ArtigoInexistente(t *testing.T) {
	uc := NewSalvarAnotacao(
		&MockArtigoRepository{
			FindByIDFn: func(ctx context.Context, id string) (*entity.Artigo, error) {
				return nil, domainErros.ErrNotFound("artigo não encontrado", "test", nil)
			},
		},
		&mockAnotacaoRepo{
			UpsertFn: func(ctx context.Context, a repository.Anotacao) error {
				t.Fatal("não deveria upsertar")
				return nil
			},
		},
	)
	_, err := uc.Execute(context.Background(), uuid.New().String(), uuid.New().String(), dto.SalvarAnotacaoRequest{
		Conteudo: json.RawMessage(`{"notes":[]}`),
	})
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.NotFound {
		t.Fatalf("esperava NotFound, got %#v", err)
	}
}

func TestSalvarAnotacao_JSONInvalido(t *testing.T) {
	uc := NewSalvarAnotacao(&MockArtigoRepository{}, &mockAnotacaoRepo{})
	_, err := uc.Execute(context.Background(), uuid.New().String(), uuid.New().String(), dto.SalvarAnotacaoRequest{
		Conteudo: json.RawMessage(`{quebrado`),
	})
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.InvalidArgument {
		t.Fatalf("esperava InvalidArgument, got %#v", err)
	}
}