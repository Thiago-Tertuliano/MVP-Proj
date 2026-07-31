package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/port"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type RefreshTokenUC struct {
	refresh repository.RefreshTokenRepository
	repo    repository.UsuarioRepository
	tokens  port.TokenGerador
	cfg     LoginConfig
}

func NewRefreshTokenUC(refresh repository.RefreshTokenRepository, repo repository.UsuarioRepository, tokens port.TokenGerador, cfg LoginConfig) *RefreshTokenUC {
	return &RefreshTokenUC{refresh: refresh, repo: repo, tokens: tokens, cfg: cfg}
}

func (uc *RefreshTokenUC) Execute(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	// 1. localiza pelo hash (nunca guardamos o token em texto puro)
	hash := sha256.Sum256([]byte(refreshToken))
	rt, err := uc.refresh.FindByHash(ctx, hex.EncodeToString(hash[:]))
	if err != nil {
		return nil, errors.ErrUnauthorized("refresh token inválido", "RefreshTokenUC.Execute", nil)
	}

	if rt.Revogado || rt.ExpiraEm.Before(time.Now()) {
		return nil, errors.ErrUnauthorized("refresh token expirado ou revogado", "RefreshTokenUC.Execute", nil)
	}

	// 2. carrega o usuário dono do token
	usuario, err := uc.repo.FindByID(ctx, rt.UsuarioID)
	if err != nil {
		return nil, errors.ErrUnauthorized("usuário não encontrado", "RefreshTokenUC.Execute", nil)
	}
	if !usuario.EstaAtiva() {
		return nil, errors.ErrUnauthorized("conta inativa", "RefreshTokenUC.Execute", nil)
	}

	// 3. emite novo par (rotação) e persiste novo hash
	tokens, err := uc.tokens.Gerar(
		port.Claims{UsuarioID: usuario.ID().String(), Email: usuario.Email().Value()},
		time.Duration(uc.cfg.AccessTTLMin)*time.Minute,
		time.Duration(uc.cfg.RefreshTTLHor)*time.Hour,
	)
	if err != nil {
		return nil, errors.ErrInternal("falha ao gerar tokens", "RefreshTokenUC.Execute", err)
	}

	if err := uc.refresh.Revoke(ctx, rt.ID); err != nil {
		return nil, errors.ErrInternal("falha ao revogar token antigo", "RefreshTokenUC.Execute", err)
	}

	newHash := sha256.Sum256([]byte(tokens.RefreshToken))
	novo := &repository.RefreshToken{
		ID:        uuid.New().String(),
		UsuarioID: usuario.ID().String(),
		TokenHash: hex.EncodeToString(newHash[:]),
		ExpiraEm:  tokens.RefreshExp,
		Revogado:  false,
		CriadoEm:  time.Now().UTC(),
	}
	if err := uc.refresh.Save(ctx, novo); err != nil {
		return nil, errors.ErrInternal("falha ao salvar refresh token", "RefreshTokenUC.Execute", err)
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