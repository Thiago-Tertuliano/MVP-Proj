package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func TestLogout_Sucesso(t *testing.T) {
	called := false
	uid := uuid.New()
	refresh := &MockRefreshTokenRepository{
		RevokeAllByUserFn: func(ctx context.Context, usuarioID uuid.UUID) error {
			called = true
			if usuarioID != uid {
				t.Fatalf("usuario_id inesperado: %s", usuarioID)
			}
			return nil
		},
	}
	uc := NewLogoutUsuario(refresh)
	if err := uc.Execute(context.Background(), uid.String()); err != nil {
		t.Fatalf("não deveria falhar: %v", err)
	}
	if !called {
		t.Fatal("deveria ter revogado tokens")
	}
}

func TestLogout_IDInvalido(t *testing.T) {
	uc := NewLogoutUsuario(&MockRefreshTokenRepository{})
	err := uc.Execute(context.Background(), "nao-uuid")
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.InvalidArgument {
		t.Errorf("esperava InvalidArgument, got %#v", err)
	}
}

func TestLogout_FalhaRepo(t *testing.T) {
	uid := uuid.New()
	refresh := &MockRefreshTokenRepository{
		RevokeAllByUserFn: func(ctx context.Context, usuarioID uuid.UUID) error {
			return context.Canceled
		},
	}
	uc := NewLogoutUsuario(refresh)
	err := uc.Execute(context.Background(), uid.String())
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.Internal {
		t.Errorf("esperava Internal, got %#v", err)
	}
}
