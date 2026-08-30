package usecase

import (
	"context"
	stderrors "errors"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type ImportarConteudo struct {
	usuarios repository.UsuarioRepository
	trilhas  repository.TrilhaRepository
	artigos  repository.ArtigoRepository
}

func NewImportarConteudo(
	usuarios repository.UsuarioRepository,
	trilhas repository.TrilhaRepository,
	artigos repository.ArtigoRepository,
) *ImportarConteudo {
	return &ImportarConteudo{usuarios: usuarios, trilhas: trilhas, artigos: artigos}
}

func (uc *ImportarConteudo) Execute(ctx context.Context, autorEmail string, plano *dto.PlanoImportacao, dryRun bool) (*dto.RelatorioImportacao, error) {
	if plano == nil || len(plano.Trilhas) == 0 {
		return nil, errors.ErrInvalidArgument("plano de importação vazio", "ImportarConteudo.Execute", nil)
	}

	rel := &dto.RelatorioImportacao{DryRun: dryRun, Avisos: append([]string{}, plano.Avisos...)}
	if dryRun {
		for _, t := range plano.Trilhas {
			rel.TrilhasOK++
			for _, m := range t.Modulos {
				rel.ArtigosOK += len(m.Aulas)
			}
		}
		return rel, nil
	}

	email, err := valueobject.NewEmail(autorEmail)
	if err != nil {
		return nil, err
	}
	autor, err := uc.usuarios.FindByEmail(ctx, email)
	if err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) {
			return nil, err
		}
		return nil, errors.ErrInternal("falha ao buscar autor", "ImportarConteudo.Execute", err)
	}

	for _, tp := range plano.Trilhas {
		criou, err := uc.upsertTrilha(ctx, autor.ID(), tp, rel)
		if err != nil {
			return rel, err
		}
		if criou {
			rel.TrilhasCriadas++
		}
		rel.TrilhasOK++
	}
	return rel, nil
}

func (uc *ImportarConteudo) upsertTrilha(ctx context.Context, autorID uuid.UUID, tp dto.TrilhaImportacao, rel *dto.RelatorioImportacao) (bool, error) {
	slug, err := valueobject.NewSlug(tp.Slug)
	if err != nil {
		return false, err
	}
	criou := false
	trilha, err := uc.trilhas.FindBySlug(ctx, slug)
	if err != nil {
		var de *errors.DomainError
		if !stderrors.As(err, &de) || de.Kind != errors.NotFound {
			return false, errors.ErrInternal("falha ao buscar trilha", "ImportarConteudo.upsertTrilha", err)
		}
		trilha, err = entity.NovaTrilha(entity.NovaTrilhaInput{
			Slug: slug, Titulo: tp.Titulo, Descricao: tp.Descricao, Ordem: tp.Ordem,
		})
		if err != nil {
			return false, err
		}
		criou = true
	} else if err := trilha.AtualizarCatalogo(tp.Titulo, tp.Descricao, tp.Ordem); err != nil {
		return false, err
	}

	for _, mp := range tp.Modulos {
		mslug, err := valueobject.NewSlug(mp.Slug)
		if err != nil {
			return criou, err
		}
		if trilha.ModuloPorSlug(mslug) == nil {
			if _, err := trilha.AdicionarModulo(mslug, mp.Titulo, mp.Descricao); err != nil {
				return criou, err
			}
		}
	}
	if err := uc.trilhas.Save(ctx, trilha); err != nil {
		return criou, errors.ErrInternal("falha ao salvar trilha", "ImportarConteudo.upsertTrilha", err)
	}
	trilha, err = uc.trilhas.FindBySlug(ctx, slug)
	if err != nil {
		return criou, err
	}

	tid := trilha.ID()
	for _, mp := range tp.Modulos {
		mslug, err := valueobject.NewSlug(mp.Slug)
		if err != nil {
			return criou, err
		}
		mod := trilha.ModuloPorSlug(mslug)
		if mod == nil {
			return criou, errors.ErrInternal("módulo sumiu após save: "+mp.Slug, "ImportarConteudo.upsertTrilha", nil)
		}
		mid := mod.ID()
		for _, aula := range mp.Aulas {
			n, err := uc.upsertArtigo(ctx, autorID, &tid, &mid, aula)
			if err != nil {
				return criou, err
			}
			if n {
				rel.ArtigosCriados++
			}
			rel.ArtigosOK++
		}
	}
	return criou, nil
}

func (uc *ImportarConteudo) upsertArtigo(ctx context.Context, autorID uuid.UUID, trilhaID, moduloID *uuid.UUID, aula dto.AulaImportacao) (bool, error) {
	slug, err := valueobject.NewSlug(aula.Slug)
	if err != nil {
		return false, err
	}
	artigo, err := uc.artigos.FindBySlug(ctx, slug)
	if err != nil {
		var de *errors.DomainError
		if !stderrors.As(err, &de) || de.Kind != errors.NotFound {
			return false, errors.ErrInternal("falha ao buscar artigo", "ImportarConteudo.upsertArtigo", err)
		}
		artigo, err = entity.NovoArtigo(entity.NovoArtigoInput{
			Titulo:    aula.Titulo,
			Subtitulo: aula.Subtitulo,
			Conteudo:  aula.Conteudo,
			Metadados: aula.Metadados,
			Slug:      slug,
			AutorID:   autorID,
			TrilhaID:  trilhaID,
			ModuloID:  moduloID,
		})
		if err != nil {
			return false, err
		}
		if err := uc.artigos.Save(ctx, artigo); err != nil {
			return false, errors.ErrInternal("falha ao criar artigo "+aula.Slug, "ImportarConteudo.upsertArtigo", err)
		}
		return true, nil
	}

	if err := artigo.AtualizarConteudo(aula.Titulo, aula.Subtitulo, artigo.CapaURL(), aula.Conteudo, aula.Metadados); err != nil {
		return false, err
	}
	if err := artigo.VincularTrilhaEModulo(trilhaID, moduloID); err != nil {
		return false, err
	}
	if err := uc.artigos.Save(ctx, artigo); err != nil {
		return false, errors.ErrInternal("falha ao atualizar artigo "+aula.Slug, "ImportarConteudo.upsertArtigo", err)
	}
	return false, nil
}
