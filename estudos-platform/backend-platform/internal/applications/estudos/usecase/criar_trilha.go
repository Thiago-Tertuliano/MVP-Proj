package usecase

import (
	"context"
	stderrors "errors"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type CriarTrilha struct {
	repo repository.TrilhaRepository
}

func NewCriarTrilha(repo repository.TrilhaRepository) *CriarTrilha {
	return &CriarTrilha{repo: repo}
}

func (uc *CriarTrilha) Execute(ctx context.Context, req dto.CriarTrilhaRequest) (*dto.TrilhaResponse, error) {
	slugRaw := req.Slug
	if slugRaw == "" {
		slugRaw = req.Titulo
	}
	slug, err := valueobject.NewSlug(slugRaw)
	if err != nil {
		return nil, err
	}
	existe, err := uc.repo.SlugExiste(ctx, slug)
	if err != nil {
		return nil, errors.ErrInternal("falha ao verificar slug", "CriarTrilha.Execute", err)
	}
	if existe {
		return nil, errors.ErrAlreadyExists("slug já utilizado", "CriarTrilha.Execute", nil)
	}

	trilha, err := entity.NovaTrilha(entity.NovaTrilhaInput{
		Slug: slug, Titulo: req.Titulo, Descricao: req.Descricao, CapaURL: req.CapaURL, Ordem: req.Ordem,
	})
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ctx, trilha); err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) {
			return nil, err
		}
		return nil, errors.ErrInternal("falha ao salvar trilha", "CriarTrilha.Execute", err)
	}
	return toTrilhaResponse(trilha), nil
}
