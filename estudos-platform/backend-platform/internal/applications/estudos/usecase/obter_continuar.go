package usecase

import (
	"context"
	stderrors "errors"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

// ObterContinuar devolve o último artigo com que o aluno interagiu (card "Continuar").
type ObterContinuar struct {
	progresso repository.ProgressoRepository
}

func NewObterContinuar(progresso repository.ProgressoRepository) *ObterContinuar {
	return &ObterContinuar{progresso: progresso}
}

func (uc *ObterContinuar) Execute(ctx context.Context, usuarioID string) (*dto.ContinuarResponse, error) {
	if _, err := uuid.Parse(usuarioID); err != nil {
		return nil, errors.ErrInvalidArgument("usuario_id inválido", "ObterContinuar.Execute", err)
	}

	ultimo, err := uc.progresso.UltimoArtigo(ctx, usuarioID)
	if err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) {
			return nil, err
		}
		return nil, errors.ErrInternal("falha ao buscar último artigo", "ObterContinuar.Execute", err)
	}
	if ultimo == nil {
		return &dto.ContinuarResponse{}, nil
	}

	return &dto.ContinuarResponse{Item: &dto.ContinuarItem{
		ArtigoID:     ultimo.ArtigoID,
		ArtigoSlug:   ultimo.ArtigoSlug,
		ArtigoTitulo: ultimo.ArtigoTitulo,
		TrilhaID:     ultimo.TrilhaID,
		TrilhaSlug:   ultimo.TrilhaSlug,
		TrilhaTitulo: ultimo.TrilhaTitulo,
		Concluido:    ultimo.Concluido,
		AtualizadoEm: ultimo.AtualizadoEm.Unix(),
	}}, nil
}
