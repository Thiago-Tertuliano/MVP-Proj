package usecase

import (
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
)

func toArtigoResponse(a *entity.Artigo) *dto.ArtigoResponse {
	resp := &dto.ArtigoResponse{
		ID:        a.ID().String(),
		Slug:      a.Slug().Value(),
		Titulo:    a.Titulo(),
		Subtitulo: a.Subtitulo(),
		CapaURL:   a.CapaURL(),
		Conteudo:  a.Conteudo(),
		Metadados: a.Metadados(),
		AutorID:   a.AutorID().String(),
		Status:    a.Status().String(),
		CreatedAt: a.CreatedAt().Unix(),
		UpdatedAt: a.UpdatedAt().Unix(),
	}
	if pe := a.PublicadoEm(); pe != nil {
		ts := pe.Unix()
		resp.PublicadoEm = &ts
	}
	resp.TrilhaID = uuidPtrString(a.TrilhaID())
	resp.ModuloID = uuidPtrString(a.ModuloID())
	return resp
}
