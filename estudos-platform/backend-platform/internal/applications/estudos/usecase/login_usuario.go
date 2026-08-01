package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	stderrors "errors"
	"time"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/port"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type TokenConfig struct {
	AccessTTLMin  int
	RefreshTTLHor int
}

type LoginUsuario struct {
	repo    repository.UsuarioRepository
	refresh repository.RefreshTokenRepository
	hasher  port.SenhaHasher
	tokens  port.TokenGerador
	cfg     TokenConfig
}

func NewLoginUsuario(repo repository.UsuarioRepository, refresh repository.RefreshTokenRepository, hasher port.SenhaHasher, tokens port.TokenGerador, cfg TokenConfig) *LoginUsuario {
	return &LoginUsuario{repo: repo, hasher: hasher, tokens: tokens, refresh: refresh, cfg: cfg}
}

func (uc *LoginUsuario) Execute(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// 1. busca usuário — se não existe, mesma mensagem de "credenciais inválidas"
	//    para não vazar quais e-mails estão cadastrados (anti enumeração).
	//    Só NotFound vira 401; falha de infraestrutura vira erro interno.
	usuario, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) && de.Kind == errors.NotFound {
			return nil, errors.ErrUnauthorized("credenciais inválidas", "LoginUsuario.Execute", nil)
		}
		return nil, errors.ErrInternal("falha ao buscar usuário", "LoginUsuario.Execute", err)
	}

	if !usuario.EstaAtiva() {
		return nil, errors.ErrUnauthorized("conta inativa", "LoginUsuario.Execute", nil)
	}

	// 2. compara o hash (via port)
	if !uc.hasher.Comparar(usuario.SenhaHash().Value(), req.Senha) {
		return nil, errors.ErrUnauthorized("credenciais inválidas", "LoginUsuario.Execute", nil)
	}

	// 3. gera par de tokens
	tokens, err := uc.tokens.Gerar(
		port.Claims{UsuarioID: usuario.ID().String(), Email: email.Value()},
		time.Duration(uc.cfg.AccessTTLMin)*time.Minute,
		time.Duration(uc.cfg.RefreshTTLHor)*time.Hour,
	)
	if err != nil {
		return nil, errors.ErrInternal("falha ao gerar tokens", "LoginUsuario.Execute", err)
	}

	// 4. rotação: revoga refreshs anteriores e persiste o novo (apenas o HASH)
	if err := uc.refresh.RevokeAllByUser(ctx, usuario.ID()); err != nil {
		return nil, errors.ErrInternal("falha ao revogar refresh tokens", "LoginUsuario.Execute", err)
	}
	hash := sha256.Sum256([]byte(tokens.RefreshToken))
	rt := &repository.RefreshToken{
		ID:        uuid.New().String(),
		UsuarioID: usuario.ID().String(),
		TokenHash: hex.EncodeToString(hash[:]),
		ExpiraEm:  tokens.RefreshExp,
		Revogado:  false,
		CriadoEm:  time.Now().UTC(),
	}
	if err := uc.refresh.Save(ctx, rt); err != nil {
		return nil, errors.ErrInternal("falha ao salvar refresh token", "LoginUsuario.Execute", err)
	}

	return &dto.AuthResponse{
		Tokens: dto.TokenResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			ExpiracaoEm:  tokens.AccessExp.Unix(),
		},
		Usuario: dto.UsuarioResponse{
			ID:    usuario.ID().String(),
			Nome:  usuario.Nome(),
			Email: usuario.Email().Value(),
		},
	}, nil
}
