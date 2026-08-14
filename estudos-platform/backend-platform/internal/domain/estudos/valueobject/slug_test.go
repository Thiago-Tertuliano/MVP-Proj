package valueobject

import "testing"

func TestNewSlug_Sucesso(t *testing.T) {
	s, err := NewSlug("  Arquitetura DDD  ")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if s.Value() != "arquitetura-ddd" {
		t.Errorf("got %q", s.Value())
	}
}

func TestNewSlug_Acentos(t *testing.T) {
	s, err := NewSlug("Introdução à Programação")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if s.Value() != "introducao-a-programacao" {
		t.Errorf("got %q", s.Value())
	}
}

func TestNewSlug_Invalido(t *testing.T) {
	if _, err := NewSlug("a"); err == nil {
		t.Fatal("esperava erro para slug curto")
	}
}
