package usecase

import (
	"context"
	stderrors "errors"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type ObterProgressoTrilha struct {
	progresso repository.ProgressoRepository
}

func NewObterProgressoTrilha(progresso repository.ProgressoRepository) *ObterProgressoTrilha {
	return &ObterProgressoTrilha{progresso: progresso}
}

func (uc *ObterProgressoTrilha) Execute(ctx context.Context, usuarioID, trilhaID string) (*dto.ProgressoTrilhaResponse, error) {
	if _, err := uuid.Parse(usuarioID); err != nil {
		return nil, errors.ErrInvalidArgument("usuario_id inválido", "ObterProgressoTrilha.Execute", err)
	}
	if _, err := uuid.Parse(trilhaID); err != nil {
		return nil, errors.ErrInvalidArgument("trilha_id inválido", "ObterProgressoTrilha.Execute", err)
	}

	concluidos, total, err := uc.progresso.CountConcluidosNaTrilha(ctx, usuarioID, trilhaID)
	if err != nil {
		return nil, uc.mapErr(err, "falha ao calcular progresso")
	}

	ids, err := uc.progresso.ListarConcluidosNaTrilha(ctx, usuarioID, trilhaID)
	if err != nil {
		return nil, uc.mapErr(err, "falha ao listar artigos concluídos")
	}
	if ids == nil {
		ids = []string{}
	}

	pct := 0.0
	if total > 0 {
		pct = float64(concluidos) / float64(total) * 100
	}
	return &dto.ProgressoTrilhaResponse{
		TrilhaID:          trilhaID,
		Concluidos:        concluidos,
		Total:             total,
		Percentual:        pct,
		ArtigosConcluidos: ids,
	}, nil
}

func (uc *ObterProgressoTrilha) mapErr(err error, msg string) error {
	var de *errors.DomainError
	if stderrors.As(err, &de) {
		return err
	}
	return errors.ErrInternal(msg, "ObterProgressoTrilha.Execute", err)
}
