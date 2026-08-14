package entity

import (
	"testing"

	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func TestTrilha_AdicionarModuloEPublicar(t *testing.T) {
	slug, _ := valueobject.NewSlug("go-avancado")
	trilha, err := NovaTrilha(NovaTrilhaInput{Slug: slug, Titulo: "Go Avançado", Descricao: "trilha"})
	if err != nil {
		t.Fatal(err)
	}
	modSlug, _ := valueobject.NewSlug("concurrency")
	if _, err := trilha.AdicionarModulo(modSlug, "Concorrência", "goroutines"); err != nil {
		t.Fatal(err)
	}
	if err := trilha.Publicar(); err != nil {
		t.Fatal(err)
	}
	if !trilha.Publicada() || len(trilha.Modulos()) != 1 {
		t.Fatal("trilha deveria estar publicada com 1 módulo")
	}
}

func TestTrilha_PublicarSemModulos(t *testing.T) {
	slug, _ := valueobject.NewSlug("vazia")
	trilha, _ := NovaTrilha(NovaTrilhaInput{Slug: slug, Titulo: "Vazia"})
	err := trilha.Publicar()
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.InvalidState {
		t.Fatalf("esperava InvalidState, got %#v", err)
	}
}
