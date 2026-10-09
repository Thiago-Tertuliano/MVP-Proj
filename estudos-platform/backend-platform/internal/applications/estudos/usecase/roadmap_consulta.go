package usecase

import (
	"context"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

// ListarRoadmaps devolve o catálogo público (apenas publicados).
type ListarRoadmaps struct {
	roadmaps repository.RoadmapRepository
}

func NewListarRoadmaps(r repository.RoadmapRepository) *ListarRoadmaps {
	return &ListarRoadmaps{roadmaps: r}
}

func (uc *ListarRoadmaps) Execute(ctx context.Context) (*dto.ListarRoadmapsResponse, error) {
	itens, err := uc.roadmaps.ListarPublicados(ctx)
	if err != nil {
		return nil, traduzir(err, "falha ao listar roadmaps", "ListarRoadmaps.Execute")
	}
	out := &dto.ListarRoadmapsResponse{Itens: make([]dto.RoadmapResumoResponse, 0, len(itens))}
	for _, it := range itens {
		out.Itens = append(out.Itens, resumoParaDTO(it))
	}
	return out, nil
}

// ObterRoadmap devolve o grafo público de um roadmap publicado (sem gabarito do quiz).
type ObterRoadmap struct {
	roadmaps repository.RoadmapRepository
}

func NewObterRoadmap(r repository.RoadmapRepository) *ObterRoadmap { return &ObterRoadmap{roadmaps: r} }

func (uc *ObterRoadmap) Execute(ctx context.Context, slug string) (*dto.RoadmapResponse, error) {
	rm, err := uc.roadmaps.FindBySlug(ctx, slug)
	if err != nil {
		return nil, traduzir(err, "falha ao buscar roadmap", "ObterRoadmap.Execute")
	}
	// Rascunho não é visível fora do editor: responde como inexistente.
	if !rm.Publicado {
		return nil, errors.ErrNotFound("roadmap não encontrado", "ObterRoadmap.Execute", nil)
	}
	return roadmapParaDTO(rm, false), nil
}
