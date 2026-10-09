package usecase

import (
	"context"
	stderrors "errors"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

// PapelResolver informa o papel efetivo do usuário (aluno/editor).
type PapelResolver interface {
	Papel(ctx context.Context, usuarioID string) (string, error)
}

// ObterPerfil busca o usuário completo pelo ID extraído do token autenticado.
type ObterPerfil struct {
	repo   repository.UsuarioRepository
	papeis PapelResolver
}

func NewObterPerfil(repo repository.UsuarioRepository) *ObterPerfil {
	return &ObterPerfil{repo: repo}
}

// ComPapel habilita o campo "papel" na resposta.
func (uc *ObterPerfil) ComPapel(p PapelResolver) *ObterPerfil {
	uc.papeis = p
	return uc
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

	resp := &dto.UsuarioResponse{
		ID:    usuario.ID().String(),
		Nome:  usuario.Nome(),
		Email: usuario.Email().Value(),
	}
	if uc.papeis != nil {
		papel, err := uc.papeis.Papel(ctx, usuarioID)
		if err != nil {
			return nil, errors.ErrInternal("falha ao obter papel", "ObterPerfil.Execute", err)
		}
		resp.Papel = papel
	}
	return resp, nil
}
