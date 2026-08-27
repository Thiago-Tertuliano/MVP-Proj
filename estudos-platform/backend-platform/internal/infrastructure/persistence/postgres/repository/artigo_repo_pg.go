package repository

import (
	"context"
	"encoding/json"
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

type ArtigoRepoPG struct {
	pool *pgxpool.Pool
}

func NewArtigoRepoPG(pool *pgxpool.Pool) *ArtigoRepoPG {
	return &ArtigoRepoPG{pool: pool}
}

func (r *ArtigoRepoPG) Save(ctx context.Context, a *entity.Artigo) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO artigos (
			id, slug, titulo, subtitulo, capa_url, conteudo, metadados,
			autor_id, status, publicado_em, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		ON CONFLICT (id) DO UPDATE SET
			titulo = EXCLUDED.titulo,
			subtitulo = EXCLUDED.subtitulo,
			capa_url = EXCLUDED.capa_url,
			conteudo = EXCLUDED.conteudo,
			metadados = EXCLUDED.metadados,
			status = EXCLUDED.status,
			publicado_em = EXCLUDED.publicado_em,
			updated_at = EXCLUDED.updated_at
	`, a.ID(), a.Slug().Value(), a.Titulo(), nullIfEmpty(a.Subtitulo()), nullIfEmpty(a.CapaURL()),
		[]byte(a.Conteudo()), []byte(a.Metadados()), a.AutorID(), a.Status().String(),
		a.PublicadoEm(), a.CreatedAt(), a.UpdatedAt())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domainErros.ErrAlreadyExists("slug já utilizado", "artigo_repo_pg.Save", nil)
		}
		return err
	}
	return nil
}

func (r *ArtigoRepoPG) FindByID(ctx context.Context, id string) (*entity.Artigo, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, slug, titulo, subtitulo, capa_url, conteudo, metadados,
		       autor_id, status, publicado_em, created_at, updated_at
		FROM artigos WHERE id = $1
	`, id)
	return scanArtigo(row)
}

func (r *ArtigoRepoPG) FindBySlug(ctx context.Context, slug valueobject.Slug) (*entity.Artigo, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, slug, titulo, subtitulo, capa_url, conteudo, metadados,
		       autor_id, status, publicado_em, created_at, updated_at
		FROM artigos WHERE slug = $1
	`, slug.Value())
	return scanArtigo(row)
}

func (r *ArtigoRepoPG) ListPublicados(ctx context.Context, limit, offset int) ([]*entity.Artigo, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, slug, titulo, subtitulo, capa_url, conteudo, metadados,
		       autor_id, status, publicado_em, created_at, updated_at
		FROM artigos
		WHERE status = 'publicado'
		ORDER BY publicado_em DESC NULLS LAST
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*entity.Artigo
	for rows.Next() {
		a, err := scanArtigo(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *ArtigoRepoPG) SlugExiste(ctx context.Context, slug valueobject.Slug) (bool, error) {
	var existe bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM artigos WHERE slug = $1)`, slug.Value()).Scan(&existe)
	return existe, err
}

type scannable interface {
	Scan(dest ...any) error
}

func scanArtigo(row scannable) (*entity.Artigo, error) {
	var (
		id          uuid.UUID
		slug        string
		titulo      string
		subtitulo   *string
		capaURL     *string
		conteudo    []byte
		metadados   []byte
		autorID     uuid.UUID
		status      string
		publicadoEm *time.Time
		createdAt   time.Time
		updatedAt   time.Time
	)
	if err := row.Scan(&id, &slug, &titulo, &subtitulo, &capaURL, &conteudo, &metadados,
		&autorID, &status, &publicadoEm, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErros.ErrNotFound("artigo não encontrado", "artigo_repo_pg.scanArtigo", nil)
		}
		return nil, err
	}
	st, err := valueobject.ParseArtigoStatus(status)
	if err != nil {
		return nil, err
	}
	var trilhaID, moduloID *uuid.UUID // Adicionado (Sprint A1)

	return entity.ReconstruirArtigo(
		id,
		valueobject.ReconstructSlug(slug),
		titulo,
		deref(subtitulo),
		deref(capaURL),
		json.RawMessage(conteudo),
		json.RawMessage(metadados),
		autorID,
		st,
		publicadoEm,
		createdAt,
		updatedAt,
		trilhaID,
		moduloID,
	), nil
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
