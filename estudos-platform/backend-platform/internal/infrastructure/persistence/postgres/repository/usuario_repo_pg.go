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

// UsuarioRepoPG implementa repository.UsuarioRepository com PostgreSQL.
type UsuarioRepoPG struct {
	pool *pgxpool.Pool
}

func NewUsuarioRepoPG(pool *pgxpool.Pool) *UsuarioRepoPG {
	return &UsuarioRepoPG{pool: pool}
}

func (r *UsuarioRepoPG) Save(ctx context.Context, u *entity.Usuario) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO usuarios (id, nome, email, senha_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			nome = EXCLUDED.nome,
			senha_hash = EXCLUDED.senha_hash,
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at
	`, u.ID(), u.Nome(), u.Email().Value(), u.SenhaHash().Value(), u.Status(), u.CreatedAt(), u.UpdatedAt())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domainErros.ErrAlreadyExists("e-mail já cadastrado", "usuario_repo_pg.Save", nil)
		}
		return err
	}
	return nil
}

func (r *UsuarioRepoPG) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.Usuario, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, nome, email, senha_hash, status, created_at, updated_at
		FROM usuarios WHERE email = $1
	`, email.Value())
	return scanUsuario(row)
}

func (r *UsuarioRepoPG) FindByID(ctx context.Context, id string) (*entity.Usuario, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, nome, email, senha_hash, status, created_at, updated_at
		FROM usuarios WHERE id = $1
	`, id)
	return scanUsuario(row)
}

func (r *UsuarioRepoPG) EmailExiste(ctx context.Context, email valueobject.Email) (bool, error) {
	var existe bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM usuarios WHERE email = $1)`, email.Value(),
	).Scan(&existe)
	if err != nil {
		return false, MapPG(err, "usuario_pg.EmailExiste")
	}
	return existe, nil
}

func scanUsuario(row pgx.Row) (*entity.Usuario, error) {
	var (
		id        uuid.UUID
		nome      string
		email     string
		senhaHash string
		status    string
		createdAt time.Time
		updatedAt time.Time
	)
	if err := row.Scan(&id, &nome, &email, &senhaHash, &status, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErros.ErrNotFound("usuário não encontrado", "usuario_repo_pg.scanUsuario", nil)
		}
		return nil, err
	}
	return entity.ReconstruirUsuario(
		id, nome,
		valueobject.ReconstructEmail(email),
		valueobject.NovoHashSenha(senhaHash),
		entity.StatusConta(status),
		createdAt, updatedAt,
	), nil
}
