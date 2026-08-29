package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	domrepo "github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type AnotacaoRepoPG struct {
	pool *pgxpool.Pool
}

func NewAnotacaoRepoPG(pool *pgxpool.Pool) *AnotacaoRepoPG {
	return &AnotacaoRepoPG{pool: pool}
}

func (r *AnotacaoRepoPG) Upsert(ctx context.Context, a domrepo.Anotacao) error {
	conteudo := a.Conteudo
	if len(conteudo) == 0 {
		conteudo = json.RawMessage("{}")
	}

	_, err := r.pool.Exec(ctx, `
		INSERT INTO anotacoes (id, usuario_id, artigo_id, conteudo, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (usuario_id, artigo_id) DO UPDATE SET
			conteudo = EXCLUDED.conteudo,
			updated_at = NOW()
	`, uuid.New().String(), a.UsuarioID, a.ArtigoID, conteudo)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgForeignKey {
			return domainErros.ErrNotFound("artigo ou usuário nao encontrado", "anotacao_repo_pg.Upsert", nil)
		}
		return MapPG(err, "anotacao_repo_pg.Upsert")
	}
	return nil
}

func (r *AnotacaoRepoPG) FindByUsuarioEArtigo(ctx context.Context, usuarioID, artigoID string) (*domrepo.Anotacao, error) {
	var conteudo []byte
	err := r.pool.QueryRow(ctx, `
		SELECT conteudo
		FROM anotacoes
		WHERE usuario_id = $1 AND artigo_id = $2
	`, usuarioID, artigoID).Scan(&conteudo)
	if err !=  nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Front recebe vazio em vez de 404 na primeira visita
			return &domrepo.Anotacao{
				UsuarioID: usuarioID,
				ArtigoID:  artigoID,
				Conteudo:  json.RawMessage("{}"),
			}, nil
		}
		return nil, MapPG(err, "anotacao_repo_pg.FindByUsuarioEArtigo")
	}
	return &domrepo.Anotacao{
		UsuarioID: usuarioID,
		ArtigoID:  artigoID,
		Conteudo:  json.RawMessage(conteudo),
	}, nil
}