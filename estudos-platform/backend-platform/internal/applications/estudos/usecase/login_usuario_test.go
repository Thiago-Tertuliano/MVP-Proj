package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/port"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func usuarioTeste(t *testing.T) *entity.Usuario {
	t.Helper()
	email, _ := valueobject.NewEmail("t@ex.com")
	u, err := entity.NovoUsuario("Thiago", email, valueobject.NovoHashSenha("hash"))
	if err != nil {
		t.Fatal(err)
	}
	return u
}

func TestLogin_Sucesso(t *testing.T) {
	ctx := context.Background()
	repo := &MockUsuarioRepository{
		FindByEmailFn: func(ctx context.Context, email valueobject.Email) (*entity.Usuario, error) {
			return usuarioTeste(t), nil
		},
	}
	refresh := &MockRefreshTokenRepository{
		RevokeAllByUserFn: func(ctx context.Context, id uuid.UUID) error { return nil },
		SaveFn:            func(ctx context.Context, rt *repository.RefreshToken) error { return nil },
	}
	hasher := &MockSenhaHasher{CompararFn: func(hash, plain string) bool { return true }}
	tokens := &MockTokenGerador{
		GerarFn: func(c port.Claims, a, r time.Duration) (*port.TokenPar, error) { return tokenParFake(), nil },
	}

	uc := NewLoginUsuario(repo, refresh, hasher, tokens, TokenConfig{AccessTTLMin: 15, RefreshTTLHor: 168})
	resp, err := uc.Execute(ctx, dto.LoginRequest{Email: "t@ex.com", Senha: "senha123"})
	if err != nil {
		t.Fatalf("não deveria retornar erro: %v", err)
	}
	if resp.Tokens.RefreshToken != "refresh" {
		t.Error("refresh token ausente")
	}
}

func TestLogin_UsuarioNaoEncontrado(t *testing.T) {
	ctx := context.Background()
	repo := &MockUsuarioRepository{
		FindByEmailFn: func(ctx context.Context, email valueobject.Email) (*entity.Usuario, error) {
			return nil, domainErros.ErrNotFound("usuário não encontrado", "mock", nil)
		},
	}
	uc := NewLoginUsuario(repo, &MockRefreshTokenRepository{}, &MockSenhaHasher{}, &MockTokenGerador{}, TokenConfig{})
	_, err := uc.Execute(ctx, dto.LoginRequest{Email: "t@ex.com", Senha: "senha123"})
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.Unauthorized {
		t.Errorf("esperava Unauthorized (anti-enumeração), got %#v", err)
	}
}

func TestLogin_SenhaErrada(t *testing.T) {
	ctx := context.Background()
	repo := &MockUsuarioRepository{
		FindByEmailFn: func(ctx context.Context, email valueobject.Email) (*entity.Usuario, error) {
			return usuarioTeste(t), nil
		},
	}
	hasher := &MockSenhaHasher{CompararFn: func(hash, plain string) bool { return false }}
	uc := NewLoginUsuario(repo, &MockRefreshTokenRepository{}, hasher, &MockTokenGerador{}, TokenConfig{})
	_, err := uc.Execute(ctx, dto.LoginRequest{Email: "t@ex.com", Senha: "errada"})
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.Unauthorized {
		t.Errorf("esperava Unauthorized, got %#v", err)
	}
}
