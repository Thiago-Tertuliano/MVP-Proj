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
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

// RegistrarConfig guarda os TTLs de token injetados pelo compositor (router).
type RegistrarConfig struct {
	AccessTTLMin  int
	RefreshTTLHor int
}

type RegistrarUsuario struct {
	repo    repository.UsuarioRepository
	refresh repository.RefreshTokenRepository
	hasher  port.SenhaHasher
	tokens  port.TokenGerador
	cfg     RegistrarConfig
}

func NewRegistrarUsuario(
	repo repository.UsuarioRepository,
	refresh repository.RefreshTokenRepository,
	hasher port.SenhaHasher,
	tokens port.TokenGerador,
	cfg RegistrarConfig,
) *RegistrarUsuario {
	return &RegistrarUsuario{repo: repo, refresh: refresh, hasher: hasher, tokens: tokens, cfg: cfg}
}

func (uc *RegistrarUsuario) Execute(ctx context.Context, req dto.RegistrarRequest) (*dto.AuthResponse, error) {
	// 1. valida e cria o VO de email (regra de domínio)
	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// 2. garante unicidade (invariante de negócio)
	existe, err := uc.repo.EmailExiste(ctx, email)
	if err != nil {
		return nil, errors.ErrInternal("falha ao verificar e-mail", "RegistrarUsuario.Execute", err)
	}
	if existe {
		return nil, errors.ErrAlreadyExists("e-mail já cadastrado", "RegistrarUsuario.Execute", nil)
	}

	// 3. gera hash da senha via port (infra decide o algoritmo)
	hash, err := uc.hasher.Hash(req.Senha)
	if err != nil {
		return nil, errors.ErrInternal("falha ao gerar hash de senha", "RegistrarUsuario.Execute", err)
	}

	// 4. constrói a entidade (validações de domínio rodam aqui)
	usuario, err := entity.NovoUsuario(req.Nome, email, valueobject.NovoHashSenha(hash))
	if err != nil {
		return nil, err
	}

	// 5. persiste via porta do domínio. Erro de domínio já mapeado (ex.: e-mail
	//    duplicado por corrida) passa direto; demais viram erro interno.
	if err := uc.repo.Save(ctx, usuario); err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) {
			return nil, err
		}
		return nil, errors.ErrInternal("falha ao salvar usuário", "RegistrarUsuario.Execute", err)
	}

	// 6. gera tokens e persiste o refresh (hash) — paridade com login
	tokens, err := uc.tokens.Gerar(
		port.Claims{UsuarioID: usuario.ID().String(), Email: email.Value()},
		time.Duration(uc.cfg.AccessTTLMin)*time.Minute,
		time.Duration(uc.cfg.RefreshTTLHor)*time.Hour,
	)
	if err != nil {
		return nil, errors.ErrInternal("falha ao gerar tokens", "RegistrarUsuario.Execute", err)
	}

	tokenHash := sha256.Sum256([]byte(tokens.RefreshToken))
	rt := &repository.RefreshToken{
		ID:        uuid.New().String(),
		UsuarioID: usuario.ID().String(),
		TokenHash: hex.EncodeToString(tokenHash[:]),
		ExpiraEm:  tokens.RefreshExp,
		Revogado:  false,
		CriadoEm:  time.Now().UTC(),
	}
	if err := uc.refresh.Save(ctx, rt); err != nil {
		return nil, errors.ErrInternal("falha ao salvar refresh token", "RegistrarUsuario.Execute", err)
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
