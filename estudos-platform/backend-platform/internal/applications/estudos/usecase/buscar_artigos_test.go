package usecase

import (
	"context"
	"testing"

	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func TestBuscarArtigos_QueryCurta(t *testing.T) {
	uc := NewBuscarArtigos(&MockArtigoRepository{})
	_, err := uc.Execute(context.Background(), "a", 10)
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.InvalidArgument {
		t.Fatalf("esperava InvalidArgument, got %#v", err)
	}
}

func TestBuscarArtigos_Sucesso(t *testing.T) {
	uc := NewBuscarArtigos(&MockArtigoRepository{
		BuscarPublicadosFn: func(ctx context.Context, q string, limit int) ([]repository.ResultadoBusca, error) {
			return []repository.ResultadoBusca{{Slug: "pacotes-em-go", Titulo: "Pacotes em Go", Similarity: 0}}, nil
		},
	})
	resp, err := uc.Execute(context.Background(), "pacotes", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Itens) != 1 || resp.Itens[0].Slug != "pacotes-em-go" {
		t.Fatalf("itens: %+v", resp.Itens)
	}
}
