package valueobject

import (
	"testing"

	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func TestNewEmail_Valido(t *testing.T) {
	e, err := NewEmail("Thiago@Exemplo.com")
	if err != nil {
		t.Fatalf("não deveria retornar erro: %v", err)
	}
	if e.Value() != "thiago@exemplo.com" {
		t.Errorf("esperava email normalizado, got %q", e.Value())
	}
}

func TestNewEmail_NormalizaEspacos(t *testing.T) {
	e, err := NewEmail("  thiago@exemplo.com  ")
	if err != nil {
		t.Fatalf("não deveria retornar erro: %v", err)
	}
	if e.Value() != "thiago@exemplo.com" {
		t.Errorf("esperava trim, got %q", e.Value())
	}
}

func TestNewEmail_Invalido(t *testing.T) {
	invalidos := []string{"", "sem-arroba", "nome@", "@dominio.com", "nome dominio.com"}
	for _, raw := range invalidos {
		_, err := NewEmail(raw)
		if err == nil {
			t.Errorf("esperava erro para %q", raw)
			continue
		}
		de, ok := err.(*domainErros.DomainError)
		if !ok || de.Kind != domainErros.InvalidArgument {
			t.Errorf("esperava DomainError InvalidArgument para %q, got %#v", raw, err)
		}
	}
}
