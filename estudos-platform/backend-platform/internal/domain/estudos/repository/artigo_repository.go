package repository

import (
	"context"

	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
)

type ResultadoBusca struct {
	Slug       string
	Titulo     string
	Similarity float64
}

type ArtigoRepository interface {
	Save(ctx context.Context, a *entity.Artigo) error
	FindByID(ctx context.Context, id string) (*entity.Artigo, error)
	FindBySlug(ctx context.Context, slug valueobject.Slug) (*entity.Artigo, error)
	ListPublicados(ctx context.Context, limit, offset int) ([]*entity.Artigo, error)
	SlugExiste(ctx context.Context, slug valueobject.Slug) (bool, error)
	AtualizarEmbedding(ctx context.Context, id string, embedding []float32) error
	BuscarPublicados(ctx context.Context, q string, limit int) ([]ResultadoBusca, error)
}
