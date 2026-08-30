package usecase

import (
	"context"
	stderrors "errors"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type ObterAnotacao struct {
	anotacoes repository.AnotacaoRepository
}

func NewObterAnotacao(anotacoes repository.AnotacaoRepository) *ObterAnotacao {
	return &ObterAnotacao{
		anotacoes: anotacoes,
	}
}

func (uc *ObterAnotacao) Execute(ctx context.Context, usuarioID, artigoID string) (*dto.AnotacaoResponse, error) {
	if _, err := uuid.Parse(usuarioID); err != nil {
		return nil, errors.ErrInvalidArgument("usuario_id inválido", "ObterAnotacao.Execute", err)
	}
	if _, err := uuid.Parse(artigoID); err != nil {
		return nil, errors.ErrInvalidArgument("artigo_id inválido", "ObterAnotacao.Execute", err)
	}

	a, err := uc.anotacoes.FindByUsuarioEArtigo(ctx, usuarioID, artigoID)
	if err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) {
			return nil, err
		}
		return nil, errors.ErrInternal("falha ao buscar anotacao", "ObterAnotacao.Execute", err)
	}

	return &dto.AnotacaoResponse{
		ArtigoID: a.ArtigoID,
		Conteudo: a.Conteudo,
	}, nil
}