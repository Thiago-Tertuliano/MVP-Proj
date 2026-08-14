package usecase

import (
	"context"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type ListarTrilhas struct {
	repo repository.TrilhaRepository
}

func NewListarTrilhas(repo repository.TrilhaRepository) *ListarTrilhas {
	return &ListarTrilhas{repo: repo}
}

func (uc *ListarTrilhas) Execute(ctx context.Context, limit, offset int) (*dto.ListarTrilhasResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	itens, err := uc.repo.ListPublicadas(ctx, limit, offset)
	if err != nil {
		return nil, errors.ErrInternal("falha ao listar trilhas", "ListarTrilhas.Execute", err)
	}
	out := make([]*dto.TrilhaResponse, 0, len(itens))
	for _, t := range itens {
		out = append(out, toTrilhaResponse(t))
	}
	return &dto.ListarTrilhasResponse{Itens: out}, nil
}
