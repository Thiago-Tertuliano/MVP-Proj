package usecase

import (
	"context"
	stderrors "errors"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type CriarArtigo struct {
	repo repository.ArtigoRepository
}

func NewCriarArtigo(repo repository.ArtigoRepository) *CriarArtigo {
	return &CriarArtigo{repo: repo}
}

func (uc *CriarArtigo) Execute(ctx context.Context, req dto.CriarArtigoRequest, autorID string) (*dto.ArtigoResponse, error) {
	autorUUID, err := uuid.Parse(autorID)
	if err != nil {
		return nil, errors.ErrInvalidArgument("autor_id inválido", "CriarArtigo.Execute", err)
	}

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
		return nil, errors.ErrInternal("falha ao verificar slug", "CriarArtigo.Execute", err)
	}
	if existe {
		return nil, errors.ErrAlreadyExists("slug já utilizado", "CriarArtigo.Execute", nil)
	}

	artigo, err := entity.NovoArtigo(entity.NovoArtigoInput{
		Titulo:    req.Titulo,
		Subtitulo: req.Subtitulo,
		CapaURL:   req.CapaURL,
		Conteudo:  req.Conteudo,
		Metadados: req.Metadados,
		Slug:      slug,
		AutorID:   autorUUID,
	})
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, artigo); err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) {
			return nil, err
		}
		return nil, errors.ErrInternal("falha ao salvar artigo", "CriarArtigo.Execute", err)
	}

	return toArtigoResponse(artigo), nil
}
