package usecase

import (
	"context"
	stderrors "errors"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type MarcarArtigoLido struct {
	artigos   repository.ArtigoRepository
	progresso repository.ProgressoRepository
}

func NewMarcarArtigoLido(artigos repository.ArtigoRepository, progresso repository.ProgressoRepository) *MarcarArtigoLido {
	return &MarcarArtigoLido{artigos: artigos, progresso: progresso}
}

func (uc *MarcarArtigoLido) Execute(ctx context.Context, usuarioID, artigoID string, concluido bool) (*dto.ProgressoArtigoResponse, error) {
	if _, err := uuid.Parse(usuarioID); err != nil {
		return nil, errors.ErrInvalidArgument("usuario_id inválido", "MarcarArtigoLido.Execute", err)
	}
	if _, err := uuid.Parse(artigoID); err != nil {
		return nil, errors.ErrInvalidArgument("artigo_id inválido", "MarcarArtigoLido.Execute", err)
	}

	if _, err := uc.artigos.FindByID(ctx, artigoID); err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) {
			return nil, err
		}
		return nil, errors.ErrInternal("falha ao buscar artigo", "MarcarArtigoLido.Execute", err)
	}

	p := repository.ProgressoArtigo{
		UsuarioID: usuarioID,
		ArtigoID:  artigoID,
		TrilhaID:  nil, // B2: preencher quando A1 ligar artigo à trilha
		Concluido: concluido,
	}
	if err := uc.progresso.UpsertArtigo(ctx, p); err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) {
			return nil, err
		}
		return nil, errors.ErrInternal("falha ao salvar progresso", "MarcarArtigoLido.Execute", err)
	}

	return &dto.ProgressoArtigoResponse{ArtigoID: artigoID, Concluido: concluido}, nil
}
