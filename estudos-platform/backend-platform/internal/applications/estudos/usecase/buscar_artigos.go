package usecase

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type BuscarArtigos struct {
	repo repository.ArtigoRepository
}

func NewBuscarArtigos(repo repository.ArtigoRepository) *BuscarArtigos {
	return &BuscarArtigos{repo: repo}
}

func (uc *BuscarArtigos) Execute(ctx context.Context, q string, limit int) (*dto.BuscaResponse, error) {
	q = strings.TrimSpace(q)
	if len(q) < 2 {
		return nil, errors.ErrInvalidArgument("q deve ter pelo menos 2 caracteres", "BuscarArtigos.Execute", nil)
	}

	itens, err := uc.repo.BuscarPublicados(ctx, q, limit)
	if err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) {
			return nil, err
		}
		return nil, errors.ErrInternal("falha na busca", "BuscarArtigos.Execute", err)
	}

	out := make([]dto.ResultadoBusca, 0, len(itens))
	for _, it := range itens {
		out = append(out, dto.ResultadoBusca{Slug: it.Slug, Titulo: it.Titulo, Similarity: it.Similarity})
	}
	return &dto.BuscaResponse{Itens: out}, nil
}
