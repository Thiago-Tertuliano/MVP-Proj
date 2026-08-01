package usecase

import (
	"context"
	stderrors "errors"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

// ObterPerfil busca o usuário completo pelo ID extraído do token autenticado.
type ObterPerfil struct {
	repo repository.UsuarioRepository
}

func NewObterPerfil(repo repository.UsuarioRepository) *ObterPerfil {
	return &ObterPerfil{repo: repo}
}

func (uc *ObterPerfil) Execute(ctx context.Context, usuarioID string) (*dto.UsuarioResponse, error) {
	usuario, err := uc.repo.FindByID(ctx, usuarioID)
	if err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) && de.Kind == errors.NotFound {
			return nil, errors.ErrNotFound("usuário não encontrado", "ObterPerfil.Execute", nil)
		}
		return nil, errors.ErrInternal("falha ao buscar usuário", "ObterPerfil.Execute", err)
	}

	return &dto.UsuarioResponse{
		ID:    usuario.ID().String(),
		Nome:  usuario.Nome(),
		Email: usuario.Email().Value(),
	}, nil
}
