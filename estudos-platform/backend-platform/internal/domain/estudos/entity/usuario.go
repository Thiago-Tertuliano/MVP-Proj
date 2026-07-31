package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/kernel"
)

type StatusConta string

const (
	StatusContaAtiva    StatusConta = "ativa"
	StatusContaSuspensa StatusConta = "suspensa"
)

// Usuario é o Aggregate Root do contexto de autenticação.
type Usuario struct {
	kernel.BaseEntity
	nome      string
	email     valueobject.Email
	senhaHash valueobject.SenhaHash
	status    StatusConta
}

func NovoUsuario(nome string, email valueobject.Email, senhaHash valueobject.SenhaHash) (*Usuario, error) {
	if len(nome) < 2 {
		return nil, errors.ErrInvalidArgument("nome deve ter no mínimo 2 caracteres", "entity.NovoUsuario", nil)
	}
	return &Usuario{
		BaseEntity: kernel.NewBaseEntity(),
		nome:       nome,
		email:      email,
		senhaHash:  senhaHash,
		status:     StatusContaAtiva,
	}, nil
}

// ReconstruirUsuario é usado pela infra ao carregar do banco.
func ReconstruirUsuario(id uuid.UUID, nome string, email valueobject.Email, senhaHash valueobject.SenhaHash, status StatusConta, createdAt, updatedAt time.Time) *Usuario {
	return &Usuario{
		BaseEntity: kernel.NewBaseEntityWithID(id, createdAt, updatedAt),
		nome:       nome,
		email:      email,
		senhaHash:  senhaHash,
		status:     status,
	}
}

func (u *Usuario) Nome() string                { return u.nome }
func (u *Usuario) Email() valueobject.Email    { return u.email }
func (u *Usuario) SenhaHash() valueobject.SenhaHash { return u.senhaHash }
func (u *Usuario) Status() StatusConta         { return u.status }

func (u *Usuario) EstaAtiva() bool {
	return u.status == StatusContaAtiva
}

func (u *Usuario) AlterarSenha(novoHash valueobject.SenhaHash) {
	u.senhaHash = novoHash
	u.Touch()
}