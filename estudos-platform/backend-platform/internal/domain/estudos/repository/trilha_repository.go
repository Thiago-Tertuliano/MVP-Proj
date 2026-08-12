package repository

import (
	"context"

	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
)

type TrilhaRepository interface {
	Save(ctx context.Context, t *entity.Trilha) error
	FindByID(ctx context.Context, id string) (*entity.Trilha, error)
	FindBySlug(ctx context.Context, slug valueobject.Slug) (*entity.Trilha, error)
	ListPublicadas(ctx context.Context, limit, offset int) ([]*entity.Trilha, error)
	SlugExiste(ctx context.Context, slug valueobject.Slug) (bool, error)
}
