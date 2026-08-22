package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type mockProgressoCount struct {
	concluidos int
	total      int
	err        error
}

func (m *mockProgressoCount) UpsertArtigo(context.Context, repository.ProgressoArtigo) error {
	return nil
}

func (m *mockProgressoCount) CountConcluidosNaTrilha(context.Context, string, string) (int, int, error) {
	return m.concluidos, m.total, m.err
}

func TestObterProgressoTrilha_Percentual(t *testing.T) {
	uc := NewObterProgressoTrilha(&mockProgressoCount{concluidos: 1, total: 4})
	resp, err := uc.Execute(context.Background(), uuid.New().String(), uuid.New().String())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Percentual != 25 || resp.Concluidos != 1 || resp.Total != 4 {
		t.Fatalf("got %+v", resp)
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
}

func TestObterProgressoTrilha_IDInvalido(t *testing.T) {
	uc := NewObterProgressoTrilha(&mockProgressoCount{})
	_, err := uc.Execute(context.Background(), "x", uuid.New().String())
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.InvalidArgument {
		t.Fatalf("esperava InvalidArgument, got %#v", err)
	}
}
