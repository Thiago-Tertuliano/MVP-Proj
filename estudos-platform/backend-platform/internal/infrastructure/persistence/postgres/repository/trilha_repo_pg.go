package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type TrilhaRepoPG struct {
	pool *pgxpool.Pool
}

func NewTrilhaRepoPG(pool *pgxpool.Pool) *TrilhaRepoPG {
	return &TrilhaRepoPG{pool: pool}
}

func (r *TrilhaRepoPG) Save(ctx context.Context, t *entity.Trilha) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO trilhas (id, slug, titulo, descricao, capa_url, ordem, publicada, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (id) DO UPDATE SET
			titulo = EXCLUDED.titulo,
			descricao = EXCLUDED.descricao,
			capa_url = EXCLUDED.capa_url,
			ordem = EXCLUDED.ordem,
			publicada = EXCLUDED.publicada,
			updated_at = EXCLUDED.updated_at
	`, t.ID(), t.Slug().Value(), t.Titulo(), nullStr(t.Descricao()), nullStr(t.CapaURL()),
		t.Ordem(), t.Publicada(), t.CreatedAt(), t.UpdatedAt())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domainErros.ErrAlreadyExists("slug já utilizado", "trilha_repo_pg.Save", nil)
		}
		return err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM modulos WHERE trilha_id = $1`, t.ID()); err != nil {
		return err
	}
	for _, m := range t.Modulos() {
		if _, err := tx.Exec(ctx, `
			INSERT INTO modulos (id, trilha_id, slug, titulo, descricao, ordem, created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
		`, m.ID(), t.ID(), m.Slug().Value(), m.Titulo(), nullStr(m.Descricao()), m.Ordem(), m.CreatedAt()); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *TrilhaRepoPG) FindByID(ctx context.Context, id string) (*entity.Trilha, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, slug, titulo, descricao, capa_url, ordem, publicada, created_at, updated_at
		FROM trilhas WHERE id = $1
	`, id)
	return r.scanTrilhaComModulos(ctx, row)
}

func (r *TrilhaRepoPG) FindBySlug(ctx context.Context, slug valueobject.Slug) (*entity.Trilha, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, slug, titulo, descricao, capa_url, ordem, publicada, created_at, updated_at
		FROM trilhas WHERE slug = $1
	`, slug.Value())
	return r.scanTrilhaComModulos(ctx, row)
}

func (r *TrilhaRepoPG) ListPublicadas(ctx context.Context, limit, offset int) ([]*entity.Trilha, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, slug, titulo, descricao, capa_url, ordem, publicada, created_at, updated_at
		FROM trilhas
		WHERE publicada = true
		ORDER BY ordem ASC, created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*entity.Trilha
	for rows.Next() {
		t, err := r.scanTrilhaComModulos(ctx, rows)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *TrilhaRepoPG) SlugExiste(ctx context.Context, slug valueobject.Slug) (bool, error) {
	var existe bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM trilhas WHERE slug = $1)`, slug.Value()).Scan(&existe)
	return existe, err
}

type trilhaRow interface {
	Scan(dest ...any) error
}

func (r *TrilhaRepoPG) scanTrilhaComModulos(ctx context.Context, row trilhaRow) (*entity.Trilha, error) {
	var (
		id         uuid.UUID
		slug       string
		titulo     string
		descricao  *string
		capaURL    *string
		ordem      int
		publicada  bool
		createdAt  time.Time
		updatedAt  time.Time
	)
	if err := row.Scan(&id, &slug, &titulo, &descricao, &capaURL, &ordem, &publicada, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErros.ErrNotFound("trilha não encontrada", "trilha_repo_pg.scan", nil)
		}
		return nil, err
	}

	modulos, err := r.loadModulos(ctx, id)
	if err != nil {
		return nil, err
	}

	return entity.ReconstruirTrilha(
		id,
		valueobject.ReconstructSlug(slug),
		titulo,
		derefStr(descricao),
		derefStr(capaURL),
		ordem,
		publicada,
		modulos,
		createdAt,
		updatedAt,
	), nil
}

func (r *TrilhaRepoPG) loadModulos(ctx context.Context, trilhaID uuid.UUID) ([]*entity.Modulo, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, slug, titulo, descricao, ordem, created_at
		FROM modulos WHERE trilha_id = $1 ORDER BY ordem ASC
	`, trilhaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*entity.Modulo
	for rows.Next() {
		var (
			id        uuid.UUID
			slug      string
			titulo    string
			descricao *string
			ordem     int
			createdAt time.Time
		)
		if err := rows.Scan(&id, &slug, &titulo, &descricao, &ordem, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, entity.ReconstruirModulo(
			id, valueobject.ReconstructSlug(slug), titulo, derefStr(descricao), ordem, createdAt,
		))
	}
	return out, rows.Err()
}

func nullStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
