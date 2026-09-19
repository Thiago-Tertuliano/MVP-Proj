package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type ObterAnotacaoUseCase struct {
	anotacoes repository.AnotacaoRepository
}

func NewObterAnotacao(anotacoes repository.AnotacaoRepository) *ObterAnotacaoUseCase {
	return &ObterAnotacaoUseCase{anotacoes: anotacoes}
}

func (uc *ObterAnotacaoUseCase) Execute(ctx context.Context, usuarioID, artigoID string) (*dto.AnotacaoResponse, error) {
	if _, err := uuid.Parse(usuarioID); err != nil {
		return nil, errors.ErrInvalidArgument("usuario_id inválido", "ObterAnotacaoUseCase.Execute", err)
	}
	if _, err := uuid.Parse(artigoID); err != nil {
		return nil, errors.ErrInvalidArgument("artigo_id inválido", "ObterAnotacaoUseCase.Execute", err)
	}

	anotacao, err := uc.anotacoes.FindByUsuarioEArtigo(ctx, usuarioID, artigoID)
	if err != nil {
		return nil, err // O repositório já retorna o DomainError adequado
	}

	return &dto.AnotacaoResponse{
		ArtigoID: anotacao.ArtigoID,
		Conteudo: anotacao.Conteudo,
	}, nil
}