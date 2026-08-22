package repository

import (
	"context"

	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
)

// UsuarioRepository é a PORTA de persistência. O domínio define o contrato;
// a implementação concreta vive na infraestrutura (usuario_repo_pg.go).
type UsuarioRepository interface {
	Save(ctx context.Context, u *entity.Usuario) error
	FindByEmail(ctx context.Context, email valueobject.Email) (*entity.Usuario, error)
	FindByID(ctx context.Context, id string) (*entity.Usuario, error)
	EmailExiste(ctx context.Context, email valueobject.Email) (bool, error)
}
