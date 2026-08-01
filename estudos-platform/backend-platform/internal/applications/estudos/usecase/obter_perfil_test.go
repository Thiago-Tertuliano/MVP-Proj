package usecase

import (
	"context"
	"testing"

	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func TestObterPerfil_Sucesso(t *testing.T) {
	ctx := context.Background()
	u := usuarioTeste(t)
	repo := &MockUsuarioRepository{
		FindByIDFn: func(ctx context.Context, id string) (*entity.Usuario, error) { return u, nil },
	}
	uc := NewObterPerfil(repo)

	resp, err := uc.Execute(ctx, u.ID().String())
	if err != nil {
		t.Fatalf("não deveria retornar erro: %v", err)
	}
	if resp.ID != u.ID().String() || resp.Nome != "Thiago" || resp.Email != "t@ex.com" {
		t.Errorf("perfil incorreto: %+v", resp)
	}
}

func TestObterPerfil_NaoEncontrado(t *testing.T) {
	ctx := context.Background()
	repo := &MockUsuarioRepository{
		FindByIDFn: func(ctx context.Context, id string) (*entity.Usuario, error) {
			return nil, domainErros.ErrNotFound("usuário não encontrado", "mock", nil)
		},
	}
	uc := NewObterPerfil(repo)

	_, err := uc.Execute(ctx, "u-inexistente")
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.NotFound {
		t.Errorf("esperava NotFound, got %#v", err)
	}
}
