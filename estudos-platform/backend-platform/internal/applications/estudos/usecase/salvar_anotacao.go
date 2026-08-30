package usecase

import (
	"context"
	"encoding/json"
	stderrors "errors"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type SalvarAnotacaoUseCase struct {
	artigos repository.ArtigoRepository
	anotacoes repository.AnotacaoRepository
}

func NewSalvarAnotacao(artigos repository.ArtigoRepository, anotacoes repository.AnotacaoRepository) *SalvarAnotacaoUseCase {
	return &SalvarAnotacaoUseCase{
		artigos: artigos,
		anotacoes: anotacoes,
	}
}

func (uc *SalvarAnotacaoUseCase) Execute(ctx context.Context, usuarioID, artigoID string, request dto.SalvarAnotacaoRequest) (*dto.AnotacaoResponse, error) {
	if _, err := uuid.Parse(usuarioID); err != nil {
		return nil, errors.ErrInvalidArgument("usuario_id inválido", "SalvarAnotacaoUseCase.Execute", err)
	}
	if _, err := uuid.Parse(artigoID); err != nil {
		return nil, errors.ErrInvalidArgument("artigo_id inválido", "SalvarAnotacaoUseCase.Execute", err)
	}
	if len(request.Conteudo) == 0 || !json.Valid(request.Conteudo) {
		return nil, errors.ErrInvalidArgument("conteudo JSON inválido", "SalvarAnotacaoUseCase.Execute", nil)
	}

	if _, err := uc.artigos.FindByID(ctx, artigoID); err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) {
			return nil, err
		}
		return nil, errors.ErrInternal("falha ao buscar", "SalvarAnotacaoUseCase.Execute", err)
	}

	a := repository.Anotacao{
		UsuarioID: usuarioID,
		ArtigoID:  artigoID,
		Conteudo:  request.Conteudo,
	}
	if err := uc.anotacoes.Upsert(ctx, a); err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) {
			return nil, err
		}
		return nil, errors.ErrInternal("falha ao salvar anotacao", "SalvarAnotacaoUseCase.Execute", err)
	}

	return &dto.AnotacaoResponse{
		ArtigoID: artigoID,
		Conteudo: request.Conteudo,
	}, nil
}