package usecase

import (
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
)

func toTrilhaResponse(t *entity.Trilha) *dto.TrilhaResponse {
	mods := make([]dto.ModuloResponse, 0, len(t.Modulos()))
	for _, m := range t.Modulos() {
		mods = append(mods, dto.ModuloResponse{
			ID:        m.ID().String(),
			Slug:      m.Slug().Value(),
			Titulo:    m.Titulo(),
			Descricao: m.Descricao(),
			Ordem:     m.Ordem(),
		})
	}
	return &dto.TrilhaResponse{
		ID:        t.ID().String(),
		Slug:      t.Slug().Value(),
		Titulo:    t.Titulo(),
		Descricao: t.Descricao(),
		CapaURL:   t.CapaURL(),
		Ordem:     t.Ordem(),
		Publicada: t.Publicada(),
		Modulos:   mods,
		CreatedAt: t.CreatedAt().Unix(),
		UpdatedAt: t.UpdatedAt().Unix(),
	}
}
