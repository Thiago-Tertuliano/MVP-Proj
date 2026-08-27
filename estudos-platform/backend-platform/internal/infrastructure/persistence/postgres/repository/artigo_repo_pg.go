package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	domrepo "github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
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
			autor_id, status, publicado_em, created_at, updated_at, trilha_id, modulo_id
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		ON CONFLICT (id) DO UPDATE SET
			titulo = EXCLUDED.titulo,
			subtitulo = EXCLUDED.subtitulo,
			capa_url = EXCLUDED.capa_url,
			conteudo = EXCLUDED.conteudo,
			metadados = EXCLUDED.metadados,
			status = EXCLUDED.status,
			publicado_em = EXCLUDED.publicado_em,
			updated_at = EXCLUDED.updated_at,
			trilha_id = EXCLUDED.trilha_id,
			modulo_id = EXCLUDED.modulo_id
	`, a.ID(), a.Slug().Value(), a.Titulo(), nullIfEmpty(a.Subtitulo()), nullIfEmpty(a.CapaURL()),
		[]byte(a.Conteudo()), []byte(a.Metadados()), a.AutorID(), a.Status().String(),
		a.PublicadoEm(), a.CreatedAt(), a.UpdatedAt(), a.TrilhaID(), a.ModuloID())
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
		       autor_id, status, publicado_em, created_at, updated_at, trilha_id, modulo_id
		FROM artigos WHERE id = $1
	`, id)
	return scanArtigo(row)
}

func (r *ArtigoRepoPG) FindBySlug(ctx context.Context, slug valueobject.Slug) (*entity.Artigo, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, slug, titulo, subtitulo, capa_url, conteudo, metadados,
		       autor_id, status, publicado_em, created_at, updated_at, trilha_id, modulo_id
		FROM artigos WHERE slug = $1
	`, slug.Value())
	return scanArtigo(row)
}

func (r *ArtigoRepoPG) ListPublicados(ctx context.Context, limit, offset int) ([]*entity.Artigo, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, slug, titulo, subtitulo, capa_url, conteudo, metadados,
		       autor_id, status, publicado_em, created_at, updated_at, trilha_id, modulo_id
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

func (r *ArtigoRepoPG) ListarPorTrilha(ctx context.Context, trilhaID uuid.UUID) ([]*entity.Artigo, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, slug, titulo, subtitulo, capa_url, conteudo, metadados,
		       autor_id, status, publicado_em, created_at, updated_at, trilha_id, modulo_id
		FROM artigos 
		WHERE trilha_id = $1
		ORDER BY created_at ASC
	`, trilhaID)
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

func (r *ArtigoRepoPG) AtualizarEmbedding(ctx context.Context, id string, embedding []float32) error {
	if len(embedding) == 0 {
		return nil
	}
	_, err := r.pool.Exec(ctx, `
		UPDATE artigos SET embedding = $2::vector, updated_at = now() WHERE id = $1
	`, id, float32VectorLiteral(embedding))
	return MapPG(err, "artigo_repo_pg.AtualizarEmbedding") // Assumindo que MapPG está definido neste ou em outro arquivo do pacote
}

func (r *ArtigoRepoPG) BuscarPublicados(ctx context.Context, q string, limit int) ([]domrepo.ResultadoBusca, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	rows, err := r.pool.Query(ctx, `
		SELECT slug, titulo, 0::float8 AS similarity
		FROM artigos
		WHERE status = 'publicado'
		  AND (titulo ILIKE '%' || $1 || '%' OR slug ILIKE '%' || $1 || '%')
		ORDER BY publicado_em DESC NULLS LAST
		LIMIT $2
	`, q, limit)
	if err != nil {
		return nil, MapPG(err, "artigo_repo_pg.BuscarPublicados")
	}
	defer rows.Close()

	var out []domrepo.ResultadoBusca
	for rows.Next() {
		var item domrepo.ResultadoBusca
		if err := rows.Scan(&item.Slug, &item.Titulo, &item.Similarity); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func float32VectorLiteral(v []float32) string {
	if len(v) == 0 {
		return "[]"
	}
	b := make([]byte, 0, 2+len(v)*8)
	b = append(b, '[')
	for i, n := range v {
		if i > 0 {
			b = append(b, ',')
		}
		b = append(b, strconv.FormatFloat(float64(n), 'f', 6, 32)...)
	}
	b = append(b, ']')
	return string(b)
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
		trilhaID    *uuid.UUID
		moduloID    *uuid.UUID
	)
	if err := row.Scan(&id, &slug, &titulo, &subtitulo, &capaURL, &conteudo, &metadados,
		&autorID, &status, &publicadoEm, &createdAt, &updatedAt, &trilhaID, &moduloID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErros.ErrNotFound("artigo não encontrado", "artigo_repo_pg.scanArtigo", nil)
		}
		return nil, err
	}
	st, err := valueobject.ParseArtigoStatus(status)
	if err != nil {
		return nil, err
	}
	
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