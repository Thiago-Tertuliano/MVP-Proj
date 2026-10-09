package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type PapelRepoPG struct {
	pool *pgxpool.Pool
}

func NewPapelRepoPG(pool *pgxpool.Pool) *PapelRepoPG { return &PapelRepoPG{pool: pool} }

func (r *PapelRepoPG) Obter(ctx context.Context, usuarioID string) (string, string, error) {
	if _, err := uuid.Parse(usuarioID); err != nil {
		return "", "", domainErros.ErrNotFound("usuário não encontrado", "papel_repo_pg.Obter", nil)
	}
	var papel, email string
	err := r.pool.QueryRow(ctx, `SELECT papel, email FROM usuarios WHERE id = $1`, usuarioID).Scan(&papel, &email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", domainErros.ErrNotFound("usuário não encontrado", "papel_repo_pg.Obter", nil)
		}
		return "", "", MapPG(err, "papel_repo_pg.Obter")
	}
	return papel, email, nil
}

func (r *PapelRepoPG) Definir(ctx context.Context, usuarioID, papel string) error {
	_, err := r.pool.Exec(ctx, `UPDATE usuarios SET papel = $2, updated_at = now() WHERE id = $1`, usuarioID, papel)
	return MapPG(err, "papel_repo_pg.Definir")
}
