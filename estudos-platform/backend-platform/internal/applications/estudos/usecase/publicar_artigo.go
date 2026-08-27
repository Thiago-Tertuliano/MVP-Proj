package usecase

import (
	"context"
	stderrors "errors"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/port"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type PublicarArtigo struct {
	repo  repository.ArtigoRepository
	embed port.EmbeddingGerador
}

func NewPublicarArtigo(repo repository.ArtigoRepository, embed port.EmbeddingGerador) *PublicarArtigo {
	return &PublicarArtigo{repo: repo, embed: embed}
}

func (uc *PublicarArtigo) Execute(ctx context.Context, id, autorID string) (*dto.ArtigoResponse, error) {
	autorUUID, err := uuid.Parse(autorID)
	if err != nil {
		return nil, errors.ErrInvalidArgument("autor_id inválido", "PublicarArtigo.Execute", err)
	}

	artigo, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) && de.Kind == errors.NotFound {
			return nil, errors.ErrNotFound("artigo não encontrado", "PublicarArtigo.Execute", nil)
		}
		return nil, errors.ErrInternal("falha ao buscar artigo", "PublicarArtigo.Execute", err)
	}
	if !artigo.EhAutor(autorUUID) {
		return nil, errors.ErrForbidden("apenas o autor pode publicar o artigo", "PublicarArtigo.Execute", nil)
	}
	if err := artigo.Publicar(); err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ctx, artigo); err != nil {
		return nil, errors.ErrInternal("falha ao publicar artigo", "PublicarArtigo.Execute", err)
	}
	if uc.embed != nil {
		vec, embErr := uc.embed.Gerar(ctx, artigo.Titulo()+" "+string(artigo.Conteudo()))
		if embErr == nil && len(vec) > 0 {
			_ = uc.repo.AtualizarEmbedding(ctx, artigo.ID().String(), vec) // embedding é best-effort
		}
	}
	return toArtigoResponse(artigo), nil
}
