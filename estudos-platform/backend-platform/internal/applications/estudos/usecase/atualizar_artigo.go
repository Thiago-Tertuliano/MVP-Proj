package usecase

import (
	"context"
	stderrors "errors"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type AtualizarArtigo struct {
	repo repository.ArtigoRepository
}

func NewAtualizarArtigo(repo repository.ArtigoRepository) *AtualizarArtigo {
	return &AtualizarArtigo{repo: repo}
}

func (uc *AtualizarArtigo) Execute(ctx context.Context, id string, autorID string, req dto.AtualizarArtigoRequest) (*dto.ArtigoResponse, error) {
	autorUUID, err := uuid.Parse(autorID)
	if err != nil {
		return nil, errors.ErrInvalidArgument("autor_id inválido", "AtualizarArtigo.Execute", err)
	}

	artigo, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) && de.Kind == errors.NotFound {
			return nil, errors.ErrNotFound("artigo não encontrado", "AtualizarArtigo.Execute", nil)
		}
		return nil, errors.ErrInternal("falha ao buscar artigo", "AtualizarArtigo.Execute", err)
	}
	if !artigo.EhAutor(autorUUID) {
		return nil, errors.ErrForbidden("apenas o autor pode editar o artigo", "AtualizarArtigo.Execute", nil)
	}

	if err := artigo.AtualizarConteudo(req.Titulo, req.Subtitulo, req.CapaURL, req.Conteudo, req.Metadados); err != nil {
		return nil, err
	}
	trilhaID, err := parseOptionalUUID(req.TrilhaID, "AtualizarArtigo.Execute")
	if err != nil {
		return nil, err
	}
	moduloID, err := parseOptionalUUID(req.ModuloID, "AtualizarArtigo.Execute")
	if err != nil {
		return nil, err
	}
	if req.TrilhaID != nil || req.ModuloID != nil {
		if err := artigo.VincularTrilhaEModulo(trilhaID, moduloID); err != nil {
			return nil, err
		}
	}
	if err := uc.repo.Save(ctx, artigo); err != nil {
		return nil, errors.ErrInternal("falha ao atualizar artigo", "AtualizarArtigo.Execute", err)
	}
	return toArtigoResponse(artigo), nil
}
