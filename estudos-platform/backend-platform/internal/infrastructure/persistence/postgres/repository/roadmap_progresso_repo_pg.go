package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
)

type RoadmapProgressoRepoPG struct {
	pool *pgxpool.Pool
}

func NewRoadmapProgressoRepoPG(pool *pgxpool.Pool) *RoadmapProgressoRepoPG {
	return &RoadmapProgressoRepoPG{pool: pool}
}

func (r *RoadmapProgressoRepoPG) ListarNosConcluidos(ctx context.Context, usuarioID, roadmapID string) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT no_id FROM roadmap_progresso WHERE usuario_id = $1 AND roadmap_id = $2`, usuarioID, roadmapID)
	if err != nil {
		return nil, MapPG(err, "roadmap_progresso_repo_pg.ListarNosConcluidos")
	}
	defer rows.Close()

	out := []string{}
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id.String())
	}
	return out, rows.Err()
}

func (r *RoadmapProgressoRepoPG) MarcarNoConcluido(ctx context.Context, usuarioID, roadmapID, noID string) (bool, error) {
	ct, err := r.pool.Exec(ctx, `
		INSERT INTO roadmap_progresso (usuario_id, no_id, roadmap_id) VALUES ($1,$2,$3)
		ON CONFLICT (usuario_id, no_id) DO NOTHING
	`, usuarioID, noID, roadmapID)
	if err != nil {
		return false, MapPG(err, "roadmap_progresso_repo_pg.MarcarNoConcluido")
	}
	return ct.RowsAffected() == 1, nil
}

func (r *RoadmapProgressoRepoPG) EstadoArtigos(ctx context.Context, usuarioID string, artigoIDs []string) (map[string]bool, error) {
	out := map[string]bool{}
	if len(artigoIDs) == 0 {
		return out, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT artigo_id, concluido FROM progresso_estudo
		WHERE usuario_id = $1 AND artigo_id = ANY($2::uuid[])
	`, usuarioID, artigoIDs)
	if err != nil {
		return nil, MapPG(err, "roadmap_progresso_repo_pg.EstadoArtigos")
	}
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		var concluido bool
		if err := rows.Scan(&id, &concluido); err != nil {
			return nil, err
		}
		out[id.String()] = concluido
	}
	return out, rows.Err()
}

func (r *RoadmapProgressoRepoPG) ContarConcluidos(ctx context.Context, usuarioID string) (int, int, error) {
	var nos, chefes int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*), COUNT(*) FILTER (WHERE n.tipo = 'chefe')
		FROM roadmap_progresso p
		JOIN roadmap_nos n ON n.id = p.no_id
		WHERE p.usuario_id = $1
	`, usuarioID).Scan(&nos, &chefes)
	return nos, chefes, MapPG(err, "roadmap_progresso_repo_pg.ContarConcluidos")
}

func (r *RoadmapProgressoRepoPG) ListarProgressoPorRoadmap(ctx context.Context, usuarioID string) ([]repository.ProgressoRoadmapItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT r.slug,
		       COUNT(n.id),
		       COUNT(p.no_id)
		FROM roadmaps r
		JOIN roadmap_nos n ON n.roadmap_id = r.id
		LEFT JOIN roadmap_progresso p ON p.no_id = n.id AND p.usuario_id = $1
		WHERE r.publicado = true
		GROUP BY r.id, r.slug
		ORDER BY r.slug
	`, usuarioID)
	if err != nil {
		return nil, MapPG(err, "roadmap_progresso_repo_pg.ListarProgressoPorRoadmap")
	}
	defer rows.Close()

	out := []repository.ProgressoRoadmapItem{}
	for rows.Next() {
		var it repository.ProgressoRoadmapItem
		if err := rows.Scan(&it.Slug, &it.Total, &it.Concluidos); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}
