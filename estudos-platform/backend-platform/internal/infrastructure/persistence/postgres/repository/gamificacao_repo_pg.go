package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
)

type GamificacaoRepoPG struct {
	pool *pgxpool.Pool
}

func NewGamificacaoRepoPG(pool *pgxpool.Pool) *GamificacaoRepoPG {
	return &GamificacaoRepoPG{pool: pool}
}

func scanGamificacao(row pgx.Row) (entity.Gamificacao, error) {
	var g entity.Gamificacao
	err := row.Scan(&g.XPTotal, &g.StreakAtual, &g.StreakMax, &g.UltimoDia)
	return g, err
}

func (r *GamificacaoRepoPG) Obter(ctx context.Context, usuarioID string) (entity.Gamificacao, error) {
	g, err := scanGamificacao(r.pool.QueryRow(ctx,
		`SELECT xp_total, streak_atual, streak_max, ultimo_dia FROM gamificacao WHERE usuario_id = $1`, usuarioID))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Gamificacao{}, nil
	}
	return g, MapPG(err, "gamificacao_repo_pg.Obter")
}

func (r *GamificacaoRepoPG) RegistrarAtividade(ctx context.Context, usuarioID string, agora time.Time) (entity.Gamificacao, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return entity.Gamificacao{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	g, err := atualizarTotais(ctx, tx, usuarioID, 0, agora)
	if err != nil {
		return entity.Gamificacao{}, err
	}
	return g, tx.Commit(ctx)
}

func (r *GamificacaoRepoPG) ConcederXP(ctx context.Context, usuarioID, origem, refID string, xp int, agora time.Time) (bool, entity.Gamificacao, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, entity.Gamificacao{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	ct, err := tx.Exec(ctx, `
		INSERT INTO xp_eventos (usuario_id, origem, ref_id, xp) VALUES ($1,$2,$3,$4)
		ON CONFLICT (usuario_id, origem, ref_id) DO NOTHING
	`, usuarioID, origem, refID, xp)
	if err != nil {
		return false, entity.Gamificacao{}, MapPG(err, "gamificacao_repo_pg.ConcederXP")
	}
	if ct.RowsAffected() == 0 {
		g, err := scanGamificacao(tx.QueryRow(ctx,
			`SELECT xp_total, streak_atual, streak_max, ultimo_dia FROM gamificacao WHERE usuario_id = $1`, usuarioID))
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return false, entity.Gamificacao{}, err
		}
		return false, g, nil
	}

	g, err := atualizarTotais(ctx, tx, usuarioID, xp, agora)
	if err != nil {
		return false, entity.Gamificacao{}, err
	}
	return true, g, tx.Commit(ctx)
}

// atualizarTotais garante a linha, trava-a e aplica XP + regra de sequência diária.
func atualizarTotais(ctx context.Context, tx pgx.Tx, usuarioID string, xp int, agora time.Time) (entity.Gamificacao, error) {
	if _, err := tx.Exec(ctx,
		`INSERT INTO gamificacao (usuario_id) VALUES ($1) ON CONFLICT (usuario_id) DO NOTHING`, usuarioID,
	); err != nil {
		return entity.Gamificacao{}, MapPG(err, "gamificacao_repo_pg.atualizarTotais")
	}
	g, err := scanGamificacao(tx.QueryRow(ctx,
		`SELECT xp_total, streak_atual, streak_max, ultimo_dia FROM gamificacao WHERE usuario_id = $1 FOR UPDATE`, usuarioID))
	if err != nil {
		return entity.Gamificacao{}, err
	}

	g.XPTotal += xp
	g.StreakAtual, g.StreakMax = entity.AtualizarStreak(g.UltimoDia, agora, g.StreakAtual, g.StreakMax)
	dia := entity.DiaCalendario(agora)
	if g.UltimoDia == nil || dia.After(*g.UltimoDia) {
		g.UltimoDia = &dia
	}

	if _, err := tx.Exec(ctx, `
		UPDATE gamificacao
		SET xp_total = $2, streak_atual = $3, streak_max = $4, ultimo_dia = $5, updated_at = now()
		WHERE usuario_id = $1
	`, usuarioID, g.XPTotal, g.StreakAtual, g.StreakMax, g.UltimoDia); err != nil {
		return entity.Gamificacao{}, err
	}
	return g, nil
}

func (r *GamificacaoRepoPG) ListarConquistas(ctx context.Context, usuarioID string) ([]repository.ConquistaObtida, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT codigo, conquistada_em FROM usuario_conquistas WHERE usuario_id = $1 ORDER BY conquistada_em`, usuarioID)
	if err != nil {
		return nil, MapPG(err, "gamificacao_repo_pg.ListarConquistas")
	}
	defer rows.Close()

	out := []repository.ConquistaObtida{}
	for rows.Next() {
		var c repository.ConquistaObtida
		if err := rows.Scan(&c.Codigo, &c.Em); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *GamificacaoRepoPG) ConcederConquista(ctx context.Context, usuarioID, codigo string) (bool, error) {
	ct, err := r.pool.Exec(ctx, `
		INSERT INTO usuario_conquistas (usuario_id, codigo) VALUES ($1,$2)
		ON CONFLICT (usuario_id, codigo) DO NOTHING
	`, usuarioID, codigo)
	if err != nil {
		return false, MapPG(err, "gamificacao_repo_pg.ConcederConquista")
	}
	return ct.RowsAffected() == 1, nil
}
