package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
)

type RefreshTokenRepoPG struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepoPG(pool *pgxpool.Pool) *RefreshTokenRepoPG {
	return &RefreshTokenRepoPG{pool: pool}
}

func (r *RefreshTokenRepoPG) Save(ctx context.Context, t *repository.RefreshToken) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO refresh_tokens (id, usuario_id, token_hash, expira_em, revogado, criado_em)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, t.ID, t.UsuarioID, t.TokenHash, t.ExpiraEm, t.Revogado, t.CriadoEm)
	return err
}

func (r *RefreshTokenRepoPG) FindByHash(ctx context.Context, hash string) (*repository.RefreshToken, error) {
	var t repository.RefreshToken
	err := r.pool.QueryRow(ctx, `
		SELECT id, usuario_id, token_hash, expira_em, revogado, criado_em
		FROM refresh_tokens WHERE token_hash = $1
	`, hash).Scan(&t.ID, &t.UsuarioID, &t.TokenHash, &t.ExpiraEm, &t.Revogado, &t.CriadoEm)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *RefreshTokenRepoPG) Revoke(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE refresh_tokens SET revogado = true WHERE id = $1`, id)
	return err
}

func (r *RefreshTokenRepoPG) RevokeAllByUser(ctx context.Context, usuarioID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE refresh_tokens SET revogado = true WHERE usuario_id = $1`, usuarioID)
	return err
}