package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/port"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func novoRegistrar(repo *MockUsuarioRepository, hasher *MockSenhaHasher, tokens *MockTokenGerador) *RegistrarUsuario {
	return NewRegistrarUsuario(repo, hasher, tokens, RegistrarConfig{AccessTTLMin: 15, RefreshTTLHor: 168})
}

func tokenParFake() *port.TokenPar {
	return &port.TokenPar{
		AccessToken: "access", RefreshToken: "refresh",
		AccessExp: time.Now().Add(15 * time.Minute), RefreshExp: time.Now().Add(168 * time.Hour),
	}
}

func TestRegistrar_Sucesso(t *testing.T) {
	ctx := context.Background()
	repo := &MockUsuarioRepository{
		EmailExisteFn: func(ctx context.Context, email valueobject.Email) (bool, error) { return false, nil },
		SaveFn:        func(ctx context.Context, u *entity.Usuario) error { return nil },
	}
	hasher := &MockSenhaHasher{
		HashFn: func(plain string) (string, error) { return "$2a$10$hash", nil },
	}
	tokens := &MockTokenGerador{
		GerarFn: func(c port.Claims, a, r time.Duration) (*port.TokenPar, error) { return tokenParFake(), nil },
	}

	uc := novoRegistrar(repo, hasher, tokens)
	resp, err := uc.Execute(ctx, dto.RegistrarRequest{Nome: "Thiago", Email: "t@ex.com", Senha: "senha123"})
	if err != nil {
		t.Fatalf("não deveria retornar erro: %v", err)
	}
	if resp.Tokens.AccessToken != "access" {
		t.Error("access token ausente na resposta")
	}
	if resp.Usuario.Email != "t@ex.com" {
		t.Errorf("email incorreto: %s", resp.Usuario.Email)
	}
}

func TestRegistrar_EmailJaCadastrado(t *testing.T) {
	ctx := context.Background()
	repo := &MockUsuarioRepository{
		EmailExisteFn: func(ctx context.Context, email valueobject.Email) (bool, error) { return true, nil },
	}
	hasher := &MockSenhaHasher{HashFn: func(plain string) (string, error) { return "h", nil }}
	tokens := &MockTokenGerador{}

	uc := novoRegistrar(repo, hasher, tokens)
	_, err := uc.Execute(ctx, dto.RegistrarRequest{Nome: "Thiago", Email: "t@ex.com", Senha: "senha123"})
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.AlreadyExists {
		t.Errorf("esperava AlreadyExists, got %#v", err)
	}
}

func TestRegistrar_EmailInvalido(t *testing.T) {
	repo := &MockUsuarioRepository{}
	hasher := &MockSenhaHasher{}
	tokens := &MockTokenGerador{}
	uc := novoRegistrar(repo, hasher, tokens)

	_, err := uc.Execute(context.Background(), dto.RegistrarRequest{Nome: "Thiago", Email: "invalido", Senha: "senha123"})
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.InvalidArgument {
		t.Errorf("esperava InvalidArgument, got %#v", err)
	}
}

func TestRegistrar_ErroNoRepo(t *testing.T) {
	ctx := context.Background()
	repo := &MockUsuarioRepository{
		EmailExisteFn: func(ctx context.Context, email valueobject.Email) (bool, error) {
			return false, context.Canceled
		},
	}
	uc := novoRegistrar(repo, &MockSenhaHasher{}, &MockTokenGerador{})
	_, err := uc.Execute(ctx, dto.RegistrarRequest{Nome: "Thiago", Email: "t@ex.com", Senha: "senha123"})
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.Internal {
		t.Errorf("esperava Internal, got %#v", err)
	}
}
