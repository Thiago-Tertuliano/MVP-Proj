package usecase

import (
	"context"
	stderrors "errors"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type PublicarTrilha struct {
	repo repository.TrilhaRepository
}

func NewPublicarTrilha(repo repository.TrilhaRepository) *PublicarTrilha {
	return &PublicarTrilha{repo: repo}
}

func (uc *PublicarTrilha) Execute(ctx context.Context, id string) (*dto.TrilhaResponse, error) {
	trilha, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) && de.Kind == errors.NotFound {
			return nil, errors.ErrNotFound("trilha não encontrada", "PublicarTrilha.Execute", nil)
		}
		return nil, errors.ErrInternal("falha ao buscar trilha", "PublicarTrilha.Execute", err)
	}
	if err := trilha.Publicar(); err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ctx, trilha); err != nil {
		return nil, errors.ErrInternal("falha ao publicar trilha", "PublicarTrilha.Execute", err)
	}
	return toTrilhaResponse(trilha), nil
}
