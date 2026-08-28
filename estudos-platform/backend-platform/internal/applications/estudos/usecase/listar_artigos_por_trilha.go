package usecase

import (
	"context"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type ListarArtigosPorTrilha struct {
	artigoRepo repository.ArtigoRepository
	trilhaRepo repository.TrilhaRepository
}

func NewListarArtigosPorTrilha(artigoRepo repository.ArtigoRepository, trilhaRepo repository.TrilhaRepository) *ListarArtigosPorTrilha {
	return &ListarArtigosPorTrilha{
		artigoRepo: artigoRepo,
		trilhaRepo: trilhaRepo,
	}
}

func (uc *ListarArtigosPorTrilha) Execute(ctx context.Context, slug string) (*dto.ListarArtigosResponse, error) {
	slugVO, err := valueobject.NewSlug(slug)
	if err != nil {
		return nil, errors.ErrInvalidArgument("slug inválido", "ListarArtigosPorTrilha.Execute", err)
	}

	trilha, err := uc.trilhaRepo.FindBySlug(ctx, slugVO)
	if err != nil {
		return nil, errors.ErrNotFound("trilha não encontrada", "ListarArtigosPorTrilha.Execute", err)
	}

	artigos, err := uc.artigoRepo.ListarPorTrilha(ctx, trilha.ID())
	if err != nil {
		return nil, errors.ErrInternal("erro ao buscar artigos da trilha", "ListarArtigosPorTrilha.Execute", err)
	}

	var itens []*dto.ArtigoResponse
	for _, a := range artigos {
		itens = append(itens, toArtigoResponse(a)) // Assumindo que toArtigoResponse existe no mesmo package
	}

	return &dto.ListarArtigosResponse{Itens: itens}, nil
}