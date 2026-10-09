package usecase

import (
	"context"
	stderrors "errors"
	"log"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

// ArtigoConcluidoHook reage a uma leitura concluída (gamificação de roadmaps).
type ArtigoConcluidoHook interface {
	AoConcluirArtigo(ctx context.Context, usuarioID, artigoID string) (*dto.ResultadoGamificacao, error)
}

type MarcarArtigoLido struct {
	artigos   repository.ArtigoRepository
	progresso repository.ProgressoRepository
	hook      ArtigoConcluidoHook
}

func NewMarcarArtigoLido(artigos repository.ArtigoRepository, progresso repository.ProgressoRepository) *MarcarArtigoLido {
	return &MarcarArtigoLido{artigos: artigos, progresso: progresso}
}

// ComGamificacao liga o efeito colateral de XP/roadmaps. É best-effort: uma falha aqui
// nunca desfaz nem bloqueia o registro da leitura.
func (uc *MarcarArtigoLido) ComGamificacao(h ArtigoConcluidoHook) *MarcarArtigoLido {
	uc.hook = h
	return uc
}

func (uc *MarcarArtigoLido) Execute(ctx context.Context, usuarioID, artigoID string, concluido bool) (*dto.ProgressoArtigoResponse, error) {
	if _, err := uuid.Parse(usuarioID); err != nil {
		return nil, errors.ErrInvalidArgument("usuario_id inválido", "MarcarArtigoLido.Execute", err)
	}
	if _, err := uuid.Parse(artigoID); err != nil {
		return nil, errors.ErrInvalidArgument("artigo_id inválido", "MarcarArtigoLido.Execute", err)
	}

	artigo, err := uc.artigos.FindByID(ctx, artigoID)
	if err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) {
			return nil, err
		}
		return nil, errors.ErrInternal("falha ao buscar artigo", "MarcarArtigoLido.Execute", err)
	}

	p := repository.ProgressoArtigo{
		UsuarioID: usuarioID,
		ArtigoID:  artigoID,
		TrilhaID:  uuidPtrString(artigo.TrilhaID()),
		Concluido: concluido,
	}
	if err := uc.progresso.UpsertArtigo(ctx, p); err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) {
			return nil, err
		}
		return nil, errors.ErrInternal("falha ao salvar progresso", "MarcarArtigoLido.Execute", err)
	}

	resp := &dto.ProgressoArtigoResponse{ArtigoID: artigoID, Concluido: concluido}
	if concluido && uc.hook != nil {
		if res, err := uc.hook.AoConcluirArtigo(ctx, usuarioID, artigoID); err != nil {
			log.Printf("gamificação: falha ao processar artigo %s: %v", artigoID, err)
		} else {
			resp.Gamificacao = res
		}
	}
	return resp, nil
}
