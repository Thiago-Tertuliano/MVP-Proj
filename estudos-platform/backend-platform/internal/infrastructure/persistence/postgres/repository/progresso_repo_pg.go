package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	domrepo "github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type ProgressoRepoPG struct {
	pool *pgxpool.Pool
}

func NewProgressoRepoPG(pool *pgxpool.Pool) *ProgressoRepoPG {
	return &ProgressoRepoPG{pool: pool}
}

func (r *ProgressoRepoPG) UpsertArtigo(ctx context.Context, p domrepo.ProgressoArtigo) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO progresso_estudo (id, usuario_id, artigo_id, trilha_id, concluido, percentual, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, now())
		ON CONFLICT (usuario_id, artigo_id) DO UPDATE SET
			concluido  = EXCLUDED.concluido,
			trilha_id  = EXCLUDED.trilha_id,
			percentual = EXCLUDED.percentual,
			updated_at = now()
	`, uuid.New(), p.UsuarioID, p.ArtigoID, p.TrilhaID, p.Concluido, boolToPercent(p.Concluido))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgForeignKey {
			return domainErros.ErrNotFound("artigo ou usuário não encontrado", "progresso_repo_pg.UpsertArtigo", nil)
		}
		return MapPG(err, "progresso_repo_pg.UpsertArtigo")
	}
	return nil
}

func (r *ProgressoRepoPG) CountConcluidosNaTrilha(ctx context.Context, usuarioID, trilhaID string) (int, int, error) {
	var concluidos, total int
	err := r.pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE p.concluido),
			COUNT(*)
		FROM artigos a
		LEFT JOIN progresso_estudo p
		  ON p.artigo_id = a.id AND p.usuario_id = $1
		WHERE a.trilha_id = $2
		  AND a.status = 'publicado'
	`, usuarioID, trilhaID).Scan(&concluidos, &total)
	if err != nil {
		return 0, 0, MapPG(err, "progresso_repo_pg.CountConcluidosNaTrilha")
	}
	return concluidos, total, nil
}

func boolToPercent(ok bool) int {
	if ok {
		return 100
	}
	return 0
}
