package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *ProgressoRepoPG) ListarConcluidosNaTrilha(ctx context.Context, usuarioID, trilhaID string) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT a.id::text
		FROM artigos a
		JOIN progresso_estudo p
		  ON p.artigo_id = a.id AND p.usuario_id = $1 AND p.concluido
		WHERE a.trilha_id = $2
		  AND a.status = 'publicado'
		ORDER BY a.created_at ASC
	`, usuarioID, trilhaID)
	if err != nil {
		return nil, MapPG(err, "progresso_repo_pg.ListarConcluidosNaTrilha")
	}
	defer rows.Close()

	out := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, MapPG(err, "progresso_repo_pg.ListarConcluidosNaTrilha")
		}
		out = append(out, id)
	}
	if err := rows.Err(); err != nil {
		return nil, MapPG(err, "progresso_repo_pg.ListarConcluidosNaTrilha")
	}
	return out, nil
}

// UltimoArtigo devolve a interação mais recente com um artigo publicado.
// Artigos de trilha em rascunho ficam de fora: o aluno não consegue abri-los.
func (r *ProgressoRepoPG) UltimoArtigo(ctx context.Context, usuarioID string) (*domrepo.UltimoProgresso, error) {
	var out domrepo.UltimoProgresso
	err := r.pool.QueryRow(ctx, `
		SELECT a.id::text, a.slug, a.titulo,
		       t.id::text, t.slug, t.titulo,
		       p.concluido, p.updated_at
		FROM progresso_estudo p
		JOIN artigos a ON a.id = p.artigo_id AND a.status = 'publicado'
		LEFT JOIN trilhas t ON t.id = a.trilha_id
		WHERE p.usuario_id = $1
		  AND (t.id IS NULL OR t.publicada = true)
		ORDER BY p.updated_at DESC
		LIMIT 1
	`, usuarioID).Scan(
		&out.ArtigoID, &out.ArtigoSlug, &out.ArtigoTitulo,
		&out.TrilhaID, &out.TrilhaSlug, &out.TrilhaTitulo,
		&out.Concluido, &out.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, MapPG(err, "progresso_repo_pg.UltimoArtigo")
	}
	return &out, nil
}

func boolToPercent(ok bool) int {
	if ok {
		return 100
	}
	return 0
}
