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
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type RoadmapRepoPG struct {
	pool *pgxpool.Pool
}

func NewRoadmapRepoPG(pool *pgxpool.Pool) *RoadmapRepoPG { return &RoadmapRepoPG{pool: pool} }

func (r *RoadmapRepoPG) Save(ctx context.Context, rm *entity.Roadmap) error {
	const op = "roadmap_repo_pg.Save"
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `
		INSERT INTO roadmaps (id, slug, titulo, descricao, icone, publicado, autor_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (id) DO UPDATE SET
			titulo = EXCLUDED.titulo,
			descricao = EXCLUDED.descricao,
			icone = EXCLUDED.icone,
			publicado = EXCLUDED.publicado,
			updated_at = EXCLUDED.updated_at
	`, rm.ID, rm.Slug, rm.Titulo, nullStr(rm.Descricao), nullStr(rm.Icone), rm.Publicado, rm.AutorID, rm.CreatedAt, rm.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domainErros.ErrAlreadyExists("slug já utilizado", op, nil)
		}
		return MapPG(err, op)
	}

	// Arestas primeiro (dependem dos nós); depois remove nós que saíram do grafo.
	if _, err := tx.Exec(ctx, `DELETE FROM roadmap_arestas WHERE roadmap_id = $1`, rm.ID); err != nil {
		return err
	}
	manter := make([]string, 0, len(rm.Nos))
	for _, n := range rm.Nos {
		manter = append(manter, n.ID)
	}
	if _, err := tx.Exec(ctx,
		`DELETE FROM roadmap_nos WHERE roadmap_id = $1 AND id <> ALL($2::uuid[])`, rm.ID, manter,
	); err != nil {
		return err
	}

	for _, n := range rm.Nos {
		var quiz []byte
		if n.Quiz != nil {
			quiz, err = json.Marshal(n.Quiz)
			if err != nil {
				return err
			}
		}
		var artigo *string
		if n.ArtigoID != "" {
			artigo = &n.ArtigoID
		}
		ct, err := tx.Exec(ctx, `
			INSERT INTO roadmap_nos (id, roadmap_id, tipo, titulo, descricao, artigo_id, pos_x, pos_y, xp, conclusao, quiz)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
			ON CONFLICT (id) DO UPDATE SET
				tipo = EXCLUDED.tipo,
				titulo = EXCLUDED.titulo,
				descricao = EXCLUDED.descricao,
				artigo_id = EXCLUDED.artigo_id,
				pos_x = EXCLUDED.pos_x,
				pos_y = EXCLUDED.pos_y,
				xp = EXCLUDED.xp,
				conclusao = EXCLUDED.conclusao,
				quiz = EXCLUDED.quiz
			WHERE roadmap_nos.roadmap_id = EXCLUDED.roadmap_id
		`, n.ID, rm.ID, string(n.Tipo), n.Titulo, nullStr(n.Descricao), artigo, n.PosX, n.PosY, n.XP, string(n.Conclusao), quiz)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23503" {
				return domainErros.ErrInvalidArgument("artigo vinculado não existe", op, nil)
			}
			return err
		}
		if ct.RowsAffected() == 0 {
			return domainErros.ErrInvalidArgument("id de nó já pertence a outro roadmap", op, nil)
		}
	}
	for _, a := range rm.Arestas {
		if _, err := tx.Exec(ctx, `
			INSERT INTO roadmap_arestas (roadmap_id, origem_id, destino_id, tipo) VALUES ($1,$2,$3,$4)
		`, rm.ID, a.OrigemID, a.DestinoID, string(a.Tipo)); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *RoadmapRepoPG) FindBySlug(ctx context.Context, slug string) (*entity.Roadmap, error) {
	var (
		rm                   entity.Roadmap
		descricao, icone     *string
		createdAt, updatedAt time.Time
		id, autor            uuid.UUID
	)
	err := r.pool.QueryRow(ctx, `
		SELECT id, slug, titulo, descricao, icone, publicado, autor_id, created_at, updated_at
		FROM roadmaps WHERE slug = $1
	`, slug).Scan(&id, &rm.Slug, &rm.Titulo, &descricao, &icone, &rm.Publicado, &autor, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErros.ErrNotFound("roadmap não encontrado", "roadmap_repo_pg.FindBySlug", nil)
		}
		return nil, MapPG(err, "roadmap_repo_pg.FindBySlug")
	}
	rm.ID, rm.AutorID = id.String(), autor.String()
	rm.Descricao, rm.Icone = derefStr(descricao), derefStr(icone)
	rm.CreatedAt, rm.UpdatedAt = createdAt, updatedAt

	rm.Nos, err = r.carregarNos(ctx, rm.ID)
	if err != nil {
		return nil, err
	}
	rm.Arestas, err = r.carregarArestas(ctx, rm.ID)
	if err != nil {
		return nil, err
	}
	return &rm, nil
}

func (r *RoadmapRepoPG) carregarNos(ctx context.Context, roadmapID string) ([]entity.RoadmapNo, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT n.id, n.tipo, n.titulo, n.descricao, n.artigo_id, a.slug, a.titulo,
		       n.pos_x, n.pos_y, n.xp, n.conclusao, n.quiz
		FROM roadmap_nos n
		LEFT JOIN artigos a ON a.id = n.artigo_id
		WHERE n.roadmap_id = $1
		ORDER BY n.pos_y, n.pos_x, n.id
	`, roadmapID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []entity.RoadmapNo{}
	for rows.Next() {
		var (
			n                    entity.RoadmapNo
			id                   uuid.UUID
			tipo, conclusao      string
			descricao, aSlug, aT *string
			artigoID             *uuid.UUID
			quiz                 []byte
		)
		if err := rows.Scan(&id, &tipo, &n.Titulo, &descricao, &artigoID, &aSlug, &aT,
			&n.PosX, &n.PosY, &n.XP, &conclusao, &quiz); err != nil {
			return nil, err
		}
		n.ID, n.Tipo, n.Conclusao = id.String(), entity.TipoNo(tipo), entity.ConclusaoNo(conclusao)
		n.Descricao = derefStr(descricao)
		if artigoID != nil {
			n.ArtigoID = artigoID.String()
		}
		n.ArtigoSlug, n.ArtigoTitulo = derefStr(aSlug), derefStr(aT)
		if len(quiz) > 0 {
			var q entity.Quiz
			if err := json.Unmarshal(quiz, &q); err != nil {
				return nil, err
			}
			n.Quiz = &q
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func (r *RoadmapRepoPG) carregarArestas(ctx context.Context, roadmapID string) ([]entity.RoadmapAresta, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, origem_id, destino_id, tipo FROM roadmap_arestas WHERE roadmap_id = $1 ORDER BY id
	`, roadmapID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []entity.RoadmapAresta{}
	for rows.Next() {
		var id, o, d uuid.UUID
		var tipo string
		if err := rows.Scan(&id, &o, &d, &tipo); err != nil {
			return nil, err
		}
		out = append(out, entity.RoadmapAresta{ID: id.String(), OrigemID: o.String(), DestinoID: d.String(), Tipo: entity.TipoAresta(tipo)})
	}
	return out, rows.Err()
}

func (r *RoadmapRepoPG) SlugExiste(ctx context.Context, slug string) (bool, error) {
	var existe bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM roadmaps WHERE slug = $1)`, slug).Scan(&existe)
	return existe, MapPG(err, "roadmap_repo_pg.SlugExiste")
}

const resumoSQL = `
	SELECT r.id, r.slug, r.titulo, r.descricao, r.icone, r.publicado, r.autor_id, r.updated_at,
	       COUNT(n.id), COUNT(n.id) FILTER (WHERE n.tipo = 'chefe'), COALESCE(SUM(n.xp), 0)
	FROM roadmaps r
	LEFT JOIN roadmap_nos n ON n.roadmap_id = r.id
`

func (r *RoadmapRepoPG) ListarPublicados(ctx context.Context) ([]repository.RoadmapResumo, error) {
	return r.listarResumos(ctx, resumoSQL+`WHERE r.publicado = true GROUP BY r.id ORDER BY r.updated_at DESC`)
}

func (r *RoadmapRepoPG) ListarPorAutor(ctx context.Context, autorID string) ([]repository.RoadmapResumo, error) {
	if _, err := uuid.Parse(autorID); err != nil {
		return []repository.RoadmapResumo{}, nil
	}
	return r.listarResumos(ctx, resumoSQL+`WHERE r.autor_id = $1 GROUP BY r.id ORDER BY r.updated_at DESC`, autorID)
}

func (r *RoadmapRepoPG) listarResumos(ctx context.Context, sql string, args ...any) ([]repository.RoadmapResumo, error) {
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, MapPG(err, "roadmap_repo_pg.listarResumos")
	}
	defer rows.Close()

	out := []repository.RoadmapResumo{}
	for rows.Next() {
		var (
			x                repository.RoadmapResumo
			id, autor        uuid.UUID
			descricao, icone *string
		)
		if err := rows.Scan(&id, &x.Slug, &x.Titulo, &descricao, &icone, &x.Publicado, &autor, &x.AtualizadoEm,
			&x.TotalNos, &x.TotalChefes, &x.XPTotal); err != nil {
			return nil, err
		}
		x.ID, x.AutorID = id.String(), autor.String()
		x.Descricao, x.Icone = derefStr(descricao), derefStr(icone)
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *RoadmapRepoPG) ListarNosPorArtigo(ctx context.Context, artigoID string) ([]repository.NoVinculado, error) {
	if _, err := uuid.Parse(artigoID); err != nil {
		return nil, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT r.slug, n.id
		FROM roadmap_nos n
		JOIN roadmaps r ON r.id = n.roadmap_id AND r.publicado = true
		WHERE n.artigo_id = $1 AND n.conclusao = 'artigo'
	`, artigoID)
	if err != nil {
		return nil, MapPG(err, "roadmap_repo_pg.ListarNosPorArtigo")
	}
	defer rows.Close()

	var out []repository.NoVinculado
	for rows.Next() {
		var v repository.NoVinculado
		var id uuid.UUID
		if err := rows.Scan(&v.RoadmapSlug, &id); err != nil {
			return nil, err
		}
		v.NoID = id.String()
		out = append(out, v)
	}
	return out, rows.Err()
}
