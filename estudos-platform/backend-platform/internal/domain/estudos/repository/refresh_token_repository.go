package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// RefreshToken é um registro de infraestrutura (não é Aggregate).
type RefreshToken struct {
	ID        string
	UsuarioID string
	TokenHash string
	ExpiraEm  time.Time
	Revogado  bool
	CriadoEm  time.Time
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, t *RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*RefreshToken, error)
	Revoke(ctx context.Context, id string) error
	RevokeAllByUser(ctx context.Context, usuarioID uuid.UUID) error
}