package external

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcryptHasher_HashEComparar(t *testing.T) {
	h := NewBcryptHasher(bcrypt.MinCost) // custo baixo p/ teste rápido
	hash, err := h.Hash("senha123")
	if err != nil {
		t.Fatalf("erro ao gerar hash: %v", err)
	}
	if !h.Comparar(hash, "senha123") {
		t.Error("senha correta deveria validar")
	}
	if h.Comparar(hash, "senha-errada") {
		t.Error("senha errada não deveria validar")
	}
}
