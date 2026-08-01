package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/port"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func refreshTokenValido() *repository.RefreshToken {
	return &repository.RefreshToken{
		ID: "rt-1", UsuarioID: "u-1", TokenHash: "hash",
		ExpiraEm: time.Now().Add(24 * time.Hour), Revogado: false,
	}
}

func TestRefresh_Sucesso(t *testing.T) {
	ctx := context.Background()
	refresh := &MockRefreshTokenRepository{
		FindByHashFn: func(ctx context.Context, hash string) (*repository.RefreshToken, error) {
			return refreshTokenValido(), nil
		},
		RevokeFn: func(ctx context.Context, id string) error { return nil },
		SaveFn:   func(ctx context.Context, rt *repository.RefreshToken) error { return nil },
	}
	repo := &MockUsuarioRepository{
		FindByIDFn: func(ctx context.Context, id string) (*entity.Usuario, error) { return usuarioTeste(t), nil },
	}
	tokens := &MockTokenGerador{
		GerarFn: func(c port.Claims, a, r time.Duration) (*port.TokenPar, error) { return tokenParFake(), nil },
	}

	uc := NewRefreshTokenUC(refresh, repo, tokens, TokenConfig{AccessTTLMin: 15, RefreshTTLHor: 168})
	resp, err := uc.Execute(ctx, "token-bruto")
	if err != nil {
		t.Fatalf("não deveria retornar erro: %v", err)
	}
	if resp.Tokens.AccessToken != "access" {
		t.Error("access token ausente")
	}
}

func TestRefresh_TokenNaoEncontrado(t *testing.T) {
	ctx := context.Background()
	refresh := &MockRefreshTokenRepository{
		FindByHashFn: func(ctx context.Context, hash string) (*repository.RefreshToken, error) {
			return nil, domainErros.ErrNotFound("não encontrado", "mock", nil)
		},
	}
	uc := NewRefreshTokenUC(refresh, &MockUsuarioRepository{}, &MockTokenGerador{}, TokenConfig{})
	_, err := uc.Execute(ctx, "token")
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.Unauthorized {
		t.Errorf("esperava Unauthorized, got %#v", err)
	}
}

func TestRefresh_TokenNaoEncontradoNil(t *testing.T) {
	ctx := context.Background()
	refresh := &MockRefreshTokenRepository{
		FindByHashFn: func(ctx context.Context, hash string) (*repository.RefreshToken, error) {
			return nil, nil
		},
	}
	uc := NewRefreshTokenUC(refresh, &MockUsuarioRepository{}, &MockTokenGerador{}, TokenConfig{})
	_, err := uc.Execute(ctx, "token")
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.Unauthorized {
		t.Errorf("esperava Unauthorized sem panic, got %#v", err)
	}
}

func TestRefresh_Revogado(t *testing.T) {
	ctx := context.Background()
	rt := refreshTokenValido()
	rt.Revogado = true
	refresh := &MockRefreshTokenRepository{
		FindByHashFn: func(ctx context.Context, hash string) (*repository.RefreshToken, error) { return rt, nil },
	}
	uc := NewRefreshTokenUC(refresh, &MockUsuarioRepository{}, &MockTokenGerador{}, TokenConfig{})
	_, err := uc.Execute(ctx, "token")
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.Unauthorized {
		t.Errorf("esperava Unauthorized, got %#v", err)
	}
}

func TestRefresh_Expirado(t *testing.T) {
	ctx := context.Background()
	rt := refreshTokenValido()
	rt.ExpiraEm = time.Now().Add(-time.Minute)
	refresh := &MockRefreshTokenRepository{
		FindByHashFn: func(ctx context.Context, hash string) (*repository.RefreshToken, error) { return rt, nil },
	}
	uc := NewRefreshTokenUC(refresh, &MockUsuarioRepository{}, &MockTokenGerador{}, TokenConfig{})
	_, err := uc.Execute(ctx, "token")
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.Unauthorized {
		t.Errorf("esperava Unauthorized, got %#v", err)
	}
}
