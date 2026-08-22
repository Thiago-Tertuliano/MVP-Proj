package valueobject

import "testing"

func TestSenhaHash_Value(t *testing.T) {
	h := NovoHashSenha("$2a$10$abcdef")
	if h.Value() != "$2a$10$abcdef" {
		t.Error("Value() deveria retornar o hash")
	}
}

func TestSenhaHash_IsZero(t *testing.T) {
	if !NovoHashSenha("").IsZero() {
		t.Error("hash vazio deveria ser zero")
	}
	if NovoHashSenha("abc").IsZero() {
		t.Error("hash não vazio não deveria ser zero")
	}
}
