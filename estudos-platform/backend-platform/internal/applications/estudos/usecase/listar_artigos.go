package usecase

import (
	"context"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type ListarArtigos struct {
	repo repository.ArtigoRepository
}

func NewListarArtigos(repo repository.ArtigoRepository) *ListarArtigos {
	return &ListarArtigos{repo: repo}
}

func (uc *ListarArtigos) Execute(ctx context.Context, limit, offset int) (*dto.ListarArtigosResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	itens, err := uc.repo.ListPublicados(ctx, limit, offset)
	if err != nil {
		return nil, errors.ErrInternal("falha ao listar artigos", "ListarArtigos.Execute", err)
	}
	out := make([]*dto.ArtigoResponse, 0, len(itens))
	for _, a := range itens {
		out = append(out, toArtigoResponse(a))
	}
	return &dto.ListarArtigosResponse{Itens: out}, nil
}
