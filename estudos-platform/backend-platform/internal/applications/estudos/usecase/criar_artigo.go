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
	repo       repository.ArtigoRepository
	trilhaRepo repository.TrilhaRepository
}

func NewCriarArtigo(repo repository.ArtigoRepository, trilhaRepo repository.TrilhaRepository) *CriarArtigo {
	return &CriarArtigo{
		repo:       repo,
		trilhaRepo: trilhaRepo,
	}
}

func (uc *CriarArtigo) Execute(ctx context.Context, req dto.CriarArtigoRequest, autorID string) (*dto.ArtigoResponse, error) {
	autorUUID, err := uuid.Parse(autorID)
	if err != nil {
		return nil, errors.ErrInvalidArgument("autor_id inválido", "CriarArtigo.Execute", err)
	}

	var trilhaUUID *uuid.UUID
	if req.TrilhaID != nil {
		parsed, err := uuid.Parse(*req.TrilhaID)
		if err != nil {
			return nil, errors.ErrInvalidArgument("trilha_id inválido", "CriarArtigo.Execute", err)
		}
		trilhaUUID = &parsed
	}

	var moduloUUID *uuid.UUID
	if req.ModuloID != nil {
		parsed, err := uuid.Parse(*req.ModuloID)
		if err != nil {
			return nil, errors.ErrInvalidArgument("modulo_id inválido", "CriarArtigo.Execute", err)
		}
		moduloUUID = &parsed
	}

	if moduloUUID != nil {
		if trilhaUUID == nil {
			return nil, errors.ErrInvalidArgument("trilha_id é obrigatório quando modulo_id é informado", "CriarArtigo.Execute", nil)
		}

		trilha, err := uc.trilhaRepo.FindByID(ctx, trilhaUUID.String())
		if err != nil {
			return nil, errors.ErrNotFound("trilha especificada não encontrada", "CriarArtigo.Execute", err)
		}

		moduloExiste := false
		for _, mod := range trilha.Modulos() {
			if mod.ID() == *moduloUUID {
				moduloExiste = true
				break
			}
		}

		if !moduloExiste {
			return nil, errors.ErrInvalidArgument("o módulo informado não pertence à trilha selecionada", "CriarArtigo.Execute", nil)
		}
	} else if trilhaUUID != nil {
		_, err := uc.trilhaRepo.FindByID(ctx, trilhaUUID.String())
		if err != nil {
			return nil, errors.ErrNotFound("trilha especificada não encontrada", "CriarArtigo.Execute", err)
		}
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
		TrilhaID:  trilhaUUID,
		ModuloID:  moduloUUID,
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

	return toArtigoResponse(artigo), nil // Assumindo que toArtigoResponse existe no mesmo package
}