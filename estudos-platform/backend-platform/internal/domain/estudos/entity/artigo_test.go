package entity

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func TestNovoArtigo_EPublicar(t *testing.T) {
	slug, _ := valueobject.NewSlug("ddd-basico")
	a, err := NovoArtigo(NovoArtigoInput{
		Titulo:   "DDD Básico",
		Conteudo: json.RawMessage(`{"blocks":[{"type":"p","text":"olá"}]}`),
		Slug:     slug,
		AutorID:  uuid.New(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if a.Status() != valueobject.ArtigoStatusRascunho {
		t.Fatal("deveria nascer como rascunho")
	}
	if err := a.Publicar(); err != nil {
		t.Fatal(err)
	}
	if a.Status() != valueobject.ArtigoStatusPublicado || a.PublicadoEm() == nil {
		t.Fatal("deveria estar publicado")
	}
}

func TestArtigo_PublicarConteudoVazio(t *testing.T) {
	slug, _ := valueobject.NewSlug("vazio")
	a, err := NovoArtigo(NovoArtigoInput{
		Titulo:  "Vazio",
		Slug:    slug,
		AutorID: uuid.New(),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = a.Publicar()
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.InvalidArgument {
		t.Fatalf("esperava InvalidArgument, got %#v", err)
	}
}
