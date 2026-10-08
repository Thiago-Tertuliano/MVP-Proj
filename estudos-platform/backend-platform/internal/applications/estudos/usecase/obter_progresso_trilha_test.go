package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type mockProgressoCount struct {
	concluidos int
	total      int
	ids        []string
	ultimo     *repository.UltimoProgresso
	err        error
}

func (m *mockProgressoCount) UpsertArtigo(context.Context, repository.ProgressoArtigo) error {
	return nil
}

func (m *mockProgressoCount) CountConcluidosNaTrilha(context.Context, string, string) (int, int, error) {
	return m.concluidos, m.total, m.err
}

func (m *mockProgressoCount) ListarConcluidosNaTrilha(context.Context, string, string) ([]string, error) {
	return m.ids, m.err
}

func (m *mockProgressoCount) UltimoArtigo(context.Context, string) (*repository.UltimoProgresso, error) {
	return m.ultimo, m.err
}

func TestObterProgressoTrilha_Percentual(t *testing.T) {
	uc := NewObterProgressoTrilha(&mockProgressoCount{concluidos: 1, total: 4, ids: []string{"a1"}})
	resp, err := uc.Execute(context.Background(), uuid.New().String(), uuid.New().String())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Percentual != 25 || resp.Concluidos != 1 || resp.Total != 4 {
		t.Fatalf("got %+v", resp)
	}
	if len(resp.ArtigosConcluidos) != 1 || resp.ArtigosConcluidos[0] != "a1" {
		t.Fatalf("artigos concluídos: %+v", resp.ArtigosConcluidos)
	}
}

func TestObterProgressoTrilha_Vazia(t *testing.T) {
	uc := NewObterProgressoTrilha(&mockProgressoCount{concluidos: 0, total: 0})
	resp, err := uc.Execute(context.Background(), uuid.New().String(), uuid.New().String())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Percentual != 0 {
		t.Fatalf("percentual: %v", resp.Percentual)
	}
	if resp.ArtigosConcluidos == nil || len(resp.ArtigosConcluidos) != 0 {
		t.Fatalf("artigos_concluidos deve ser lista vazia (não nula): %#v", resp.ArtigosConcluidos)
	}
}

func TestObterProgressoTrilha_IDInvalido(t *testing.T) {
	uc := NewObterProgressoTrilha(&mockProgressoCount{})
	_, err := uc.Execute(context.Background(), "x", uuid.New().String())
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.InvalidArgument {
		t.Fatalf("esperava InvalidArgument, got %#v", err)
	}
}

func TestObterContinuar_SemHistorico(t *testing.T) {
	uc := NewObterContinuar(&mockProgressoCount{})
	resp, err := uc.Execute(context.Background(), uuid.New().String())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Item != nil {
		t.Fatalf("esperava item nulo, got %+v", resp.Item)
	}
}

func TestObterContinuar_ComHistorico(t *testing.T) {
	slug := "go-basico"
	quando := time.Unix(1700000000, 0)
	uc := NewObterContinuar(&mockProgressoCount{ultimo: &repository.UltimoProgresso{
		ArtigoID: "a1", ArtigoSlug: "pacotes-em-go", ArtigoTitulo: "Pacotes em Go",
		TrilhaSlug: &slug, Concluido: true, AtualizadoEm: quando,
	}})
	resp, err := uc.Execute(context.Background(), uuid.New().String())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Item == nil || resp.Item.ArtigoSlug != "pacotes-em-go" || !resp.Item.Concluido || resp.Item.AtualizadoEm != 1700000000 {
		t.Fatalf("got %+v", resp.Item)
	}
	if resp.Item.TrilhaSlug == nil || *resp.Item.TrilhaSlug != "go-basico" {
		t.Fatalf("trilha_slug: %+v", resp.Item.TrilhaSlug)
	}
}

func TestObterContinuar_UsuarioInvalido(t *testing.T) {
	uc := NewObterContinuar(&mockProgressoCount{})
	_, err := uc.Execute(context.Background(), "x")
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.InvalidArgument {
		t.Fatalf("esperava InvalidArgument, got %#v", err)
	}
}
