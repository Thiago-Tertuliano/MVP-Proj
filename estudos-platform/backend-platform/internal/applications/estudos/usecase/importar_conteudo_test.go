package usecase

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func planoStub() *dto.PlanoImportacao {
	return &dto.PlanoImportacao{
		Trilhas: []dto.TrilhaImportacao{{
			Slug: "dados", Titulo: "Dados", Descricao: "Catálogo", Ordem: 0,
			Modulos: []dto.ModuloImportacao{{
				Slug: "catalogo", Titulo: "Catálogo",
				Aulas: []dto.AulaImportacao{{
					Slug:      "spark",
					Titulo:    "SPARK",
					Conteudo:  json.RawMessage(`{"blocks":[{"type":"p","text":"ok"}]}`),
					Metadados: json.RawMessage(`{"origem":"courses-md"}`),
				}},
			}},
		}},
	}
}

func TestImportarConteudo_DryRunNaoGrava(t *testing.T) {
	saved := 0
	uc := NewImportarConteudo(
		&MockUsuarioRepository{},
		&MockTrilhaRepository{SaveFn: func(ctx context.Context, tr *entity.Trilha) error {
			saved++
			return nil
		}},
		&MockArtigoRepository{SaveFn: func(ctx context.Context, a *entity.Artigo) error {
			saved++
			return nil
		}},
	)
	rel, err := uc.Execute(context.Background(), "autor.seed@estudos.local", planoStub(), true)
	if err != nil {
		t.Fatal(err)
	}
	if !rel.DryRun || rel.ArtigosOK != 1 || rel.TrilhasOK != 1 || saved != 0 {
		t.Fatalf("rel=%+v saved=%d", rel, saved)
	}
}

func TestImportarConteudo_CriaTrilhaEArtigo(t *testing.T) {
	email, _ := valueobject.NewEmail("autor.seed@estudos.local")
	autor, err := entity.NovoUsuario("Autor Seed", email, valueobject.NovoHashSenha("hash"))
	if err != nil {
		t.Fatal(err)
	}

	trilhas := map[string]*entity.Trilha{}
	artigos := map[string]*entity.Artigo{}

	uc := NewImportarConteudo(
		&MockUsuarioRepository{
			FindByEmailFn: func(ctx context.Context, e valueobject.Email) (*entity.Usuario, error) {
				return autor, nil
			},
		},
		&MockTrilhaRepository{
			FindBySlugFn: func(ctx context.Context, slug valueobject.Slug) (*entity.Trilha, error) {
				if t, ok := trilhas[slug.Value()]; ok {
					return t, nil
				}
				return nil, domainErros.ErrNotFound("trilha não encontrada", "test", nil)
			},
			SaveFn: func(ctx context.Context, tr *entity.Trilha) error {
				trilhas[tr.Slug().Value()] = tr
				return nil
			},
		},
		&MockArtigoRepository{
			FindBySlugFn: func(ctx context.Context, slug valueobject.Slug) (*entity.Artigo, error) {
				if a, ok := artigos[slug.Value()]; ok {
					return a, nil
				}
				return nil, domainErros.ErrNotFound("artigo não encontrado", "test", nil)
			},
			SaveFn: func(ctx context.Context, a *entity.Artigo) error {
				artigos[a.Slug().Value()] = a
				return nil
			},
		},
	)

	rel, err := uc.Execute(context.Background(), "autor.seed@estudos.local", planoStub(), false)
	if err != nil {
		t.Fatal(err)
	}
	if rel.TrilhasCriadas != 1 || rel.ArtigosCriados != 1 {
		t.Fatalf("%+v artigos=%d", rel, len(artigos))
	}
	rel2, err := uc.Execute(context.Background(), "autor.seed@estudos.local", planoStub(), false)
	if err != nil {
		t.Fatal(err)
	}
	if rel2.ArtigosCriados != 0 || len(artigos) != 1 {
		t.Fatalf("segunda vez duplicou: %+v n=%d", rel2, len(artigos))
	}
	if artigos["spark"].Status() != valueobject.ArtigoStatusRascunho {
		t.Fatal("novo artigo deve nascer rascunho")
	}
	if artigos["spark"].TrilhaID() == nil {
		t.Fatal("trilha_id obrigatório")
	}
}

func TestImportarConteudo_PlanoVazio(t *testing.T) {
	uc := NewImportarConteudo(&MockUsuarioRepository{}, &MockTrilhaRepository{}, &MockArtigoRepository{})
	_, err := uc.Execute(context.Background(), "a@b.com", &dto.PlanoImportacao{}, false)
	if err == nil {
		t.Fatal("esperava erro")
	}
}

func TestImportarConteudo_NaoRebaixaPublicado(t *testing.T) {
	email, _ := valueobject.NewEmail("autor.seed@estudos.local")
	autor, err := entity.NovoUsuario("Autor Seed", email, valueobject.NovoHashSenha("hash"))
	if err != nil {
		t.Fatal(err)
	}
	slug, _ := valueobject.NewSlug("spark")
	tid, mid := uuid.New(), uuid.New()
	now := time.Now().UTC()
	pub := now
	existente := entity.ReconstruirArtigo(
		uuid.New(), slug, "SPARK", "", "",
		json.RawMessage(`{"blocks":[]}`), json.RawMessage(`{}`),
		autor.ID(), valueobject.ArtigoStatusPublicado, &pub, now, now, &tid, &mid,
	)
	artigos := map[string]*entity.Artigo{"spark": existente}
	trilhas := map[string]*entity.Trilha{}

	uc := NewImportarConteudo(
		&MockUsuarioRepository{
			FindByEmailFn: func(ctx context.Context, e valueobject.Email) (*entity.Usuario, error) {
				return autor, nil
			},
		},
		&MockTrilhaRepository{
			FindBySlugFn: func(ctx context.Context, s valueobject.Slug) (*entity.Trilha, error) {
				if t, ok := trilhas[s.Value()]; ok {
					return t, nil
				}
				return nil, domainErros.ErrNotFound("trilha não encontrada", "test", nil)
			},
			SaveFn: func(ctx context.Context, tr *entity.Trilha) error {
				trilhas[tr.Slug().Value()] = tr
				return nil
			},
		},
		&MockArtigoRepository{
			FindBySlugFn: func(ctx context.Context, s valueobject.Slug) (*entity.Artigo, error) {
				if a, ok := artigos[s.Value()]; ok {
					return a, nil
				}
				return nil, domainErros.ErrNotFound("artigo não encontrado", "test", nil)
			},
			SaveFn: func(ctx context.Context, a *entity.Artigo) error {
				artigos[a.Slug().Value()] = a
				return nil
			},
		},
	)

	rel, err := uc.Execute(context.Background(), "autor.seed@estudos.local", planoStub(), false)
	if err != nil {
		t.Fatal(err)
	}
	if rel.ArtigosCriados != 0 {
		t.Fatalf("não deveria criar: %+v", rel)
	}
	if artigos["spark"].Status() != valueobject.ArtigoStatusPublicado {
		t.Fatal("job não pode rebaixar publicado")
	}
}
