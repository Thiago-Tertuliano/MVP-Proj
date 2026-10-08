package usecase

import (
	"context"
	stderrors "errors"

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

// Execute devolve os artigos PUBLICADOS de uma trilha PUBLICADA, na ordem do
// roadmap (ordem do módulo, depois criação). Trilha em rascunho responde 404,
// igual a GET /trilhas/{slug} — o aluno nunca enxerga rascunho.
func (uc *ListarArtigosPorTrilha) Execute(ctx context.Context, slug string) (*dto.ListarArtigosResponse, error) {
	slugVO, err := valueobject.NewSlug(slug)
	if err != nil {
		return nil, errors.ErrInvalidArgument("slug inválido", "ListarArtigosPorTrilha.Execute", err)
	}

	trilha, err := uc.trilhaRepo.FindBySlug(ctx, slugVO)
	if err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) && de.Kind == errors.NotFound {
			return nil, errors.ErrNotFound("trilha não encontrada", "ListarArtigosPorTrilha.Execute", nil)
		}
		return nil, errors.ErrInternal("falha ao buscar trilha", "ListarArtigosPorTrilha.Execute", err)
	}
	if !trilha.Publicada() {
		return nil, errors.ErrNotFound("trilha não encontrada", "ListarArtigosPorTrilha.Execute", nil)
	}

	artigos, err := uc.artigoRepo.ListarPorTrilha(ctx, trilha.ID())
	if err != nil {
		return nil, errors.ErrInternal("erro ao buscar artigos da trilha", "ListarArtigosPorTrilha.Execute", err)
	}

	itens := make([]*dto.ArtigoResponse, 0, len(artigos))
	for _, a := range artigos {
		if a.Status() != valueobject.ArtigoStatusPublicado {
			continue
		}
		itens = append(itens, toArtigoResponse(a))
	}

	return &dto.ListarArtigosResponse{Itens: itens}, nil
}
