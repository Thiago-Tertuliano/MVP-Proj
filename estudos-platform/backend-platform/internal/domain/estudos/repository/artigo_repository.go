package repository

import (
	"context"

	"github.com/google/uuid"
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
	ListarPorTrilha(ctx context.Context, trilhaID uuid.UUID) ([]*entity.Artigo, error) // Adicionado Sprint A2
	SlugExiste(ctx context.Context, slug valueobject.Slug) (bool, error)
	AtualizarEmbedding(ctx context.Context, id string, embedding []float32) error
	BuscarPublicados(ctx context.Context, q string, limit int) ([]ResultadoBusca, error)
}