package usecase

import (
	"context"
	stderrors "errors"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type ObterTrilha struct {
	repo repository.TrilhaRepository
}

func NewObterTrilha(repo repository.TrilhaRepository) *ObterTrilha {
	return &ObterTrilha{repo: repo}
}

func (uc *ObterTrilha) Execute(ctx context.Context, slugRaw string) (*dto.TrilhaResponse, error) {
	slug, err := valueobject.NewSlug(slugRaw)
	if err != nil {
		return nil, err
	}
	trilha, err := uc.repo.FindBySlug(ctx, slug)
	if err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) && de.Kind == errors.NotFound {
			return nil, errors.ErrNotFound("trilha não encontrada", "ObterTrilha.Execute", nil)
		}
		return nil, errors.ErrInternal("falha ao buscar trilha", "ObterTrilha.Execute", err)
	}
	if !trilha.Publicada() {
		return nil, errors.ErrNotFound("trilha não encontrada", "ObterTrilha.Execute", nil)
	}
	return toTrilhaResponse(trilha), nil
}
