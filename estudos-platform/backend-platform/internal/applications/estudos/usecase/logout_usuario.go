package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

// LogoutUsuario revoga todos os refresh tokens do usuário autenticado.
type LogoutUsuario struct {
	refresh repository.RefreshTokenRepository
}

func NewLogoutUsuario(refresh repository.RefreshTokenRepository) *LogoutUsuario {
	return &LogoutUsuario{refresh: refresh}
}

func (uc *LogoutUsuario) Execute(ctx context.Context, usuarioID string) error {
	id, err := uuid.Parse(usuarioID)
	if err != nil {
		return errors.ErrInvalidArgument("usuario_id inválido", "LogoutUsuario.Execute", err)
	}
	if err := uc.refresh.RevokeAllByUser(ctx, id); err != nil {
		return errors.ErrInternal("falha ao revogar refresh tokens", "LogoutUsuario.Execute", err)
	}
	return nil
}
