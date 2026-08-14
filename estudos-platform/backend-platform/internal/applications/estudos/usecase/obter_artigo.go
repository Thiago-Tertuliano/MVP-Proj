package usecase

import (
	"context"
	stderrors "errors"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type ObterArtigo struct {
	repo repository.ArtigoRepository
}

func NewObterArtigo(repo repository.ArtigoRepository) *ObterArtigo {
	return &ObterArtigo{repo: repo}
}

func (uc *ObterArtigo) Execute(ctx context.Context, slugRaw string) (*dto.ArtigoResponse, error) {
	slug, err := valueobject.NewSlug(slugRaw)
	if err != nil {
		return nil, err
	}
	artigo, err := uc.repo.FindBySlug(ctx, slug)
	if err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) && de.Kind == errors.NotFound {
			return nil, errors.ErrNotFound("artigo não encontrado", "ObterArtigo.Execute", nil)
		}
		return nil, errors.ErrInternal("falha ao buscar artigo", "ObterArtigo.Execute", err)
	}
	if artigo.Status() != valueobject.ArtigoStatusPublicado {
		return nil, errors.ErrNotFound("artigo não encontrado", "ObterArtigo.Execute", nil)
	}
	return toArtigoResponse(artigo), nil
}
