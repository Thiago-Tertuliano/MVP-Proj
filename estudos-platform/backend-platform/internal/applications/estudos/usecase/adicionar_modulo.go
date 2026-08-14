package usecase

import (
	"context"
	stderrors "errors"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type AdicionarModulo struct {
	repo repository.TrilhaRepository
}

func NewAdicionarModulo(repo repository.TrilhaRepository) *AdicionarModulo {
	return &AdicionarModulo{repo: repo}
}

func (uc *AdicionarModulo) Execute(ctx context.Context, trilhaID string, req dto.AdicionarModuloRequest) (*dto.TrilhaResponse, error) {
	trilha, err := uc.repo.FindByID(ctx, trilhaID)
	if err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) && de.Kind == errors.NotFound {
			return nil, errors.ErrNotFound("trilha não encontrada", "AdicionarModulo.Execute", nil)
		}
		return nil, errors.ErrInternal("falha ao buscar trilha", "AdicionarModulo.Execute", err)
	}

	slugRaw := req.Slug
	if slugRaw == "" {
		slugRaw = req.Titulo
	}
	slug, err := valueobject.NewSlug(slugRaw)
	if err != nil {
		return nil, err
	}
	if _, err := trilha.AdicionarModulo(slug, req.Titulo, req.Descricao); err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ctx, trilha); err != nil {
		return nil, errors.ErrInternal("falha ao salvar módulo", "AdicionarModulo.Execute", err)
	}
	return toTrilhaResponse(trilha), nil
}
