package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func uuidFromString(t *testing.T, s string) uuid.UUID {
	t.Helper()
	id, err := uuid.Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func novoEmail(t *testing.T, raw string) valueobject.Email {
	t.Helper()
	e, err := valueobject.NewEmail(raw)
	if err != nil {
		t.Fatalf("email inválido de teste: %v", err)
	}
	return e
}

func TestNovoUsuario_Sucesso(t *testing.T) {
	u, err := NovoUsuario("Thiago", novoEmail(t, "t@ex.com"), valueobject.NovoHashSenha("hash"))
	if err != nil {
		t.Fatalf("não deveria retornar erro: %v", err)
	}
	if u.ID().String() == "" {
		t.Error("ID não deveria ser vazio")
	}
	if !u.EstaAtiva() {
		t.Error("novo usuário deveria estar ativo")
	}
	if u.Nome() != "Thiago" {
		t.Errorf("nome incorreto: %q", u.Nome())
	}
}

func TestNovoUsuario_NomeCurto(t *testing.T) {
	_, err := NovoUsuario("A", novoEmail(t, "t@ex.com"), valueobject.NovoHashSenha("hash"))
	if err == nil {
		t.Fatal("esperava erro para nome com 1 caractere")
	}
	if de, ok := err.(*domainErros.DomainError); !ok || de.Kind != domainErros.InvalidArgument {
		t.Errorf("esperava InvalidArgument, got %#v", err)
	}
}

func TestUsuario_AlterarSenha_AtualizaTimestamp(t *testing.T) {
	u, _ := NovoUsuario("Thiago", novoEmail(t, "t@ex.com"), valueobject.NovoHashSenha("hash1"))
	antes := u.UpdatedAt()
	time.Sleep(10 * time.Millisecond)

	u.AlterarSenha(valueobject.NovoHashSenha("hash2"))

	if u.SenhaHash().Value() != "hash2" {
		t.Error("senha não atualizada")
	}
	if !u.UpdatedAt().After(antes) {
		t.Error("updated_at deveria avançar após AlterarSenha")
	}
}

func TestReconstruirUsuario(t *testing.T) {
	id := "550e8400-e29b-41d4-a716-446655440000"
	now := time.Now().UTC()
	u := ReconstruirUsuario(
		uuidFromString(t, id), "Thiago",
		novoEmail(t, "t@ex.com"), valueobject.NovoHashSenha("hash"),
		StatusContaSuspensa, now, now,
	)
	if u.ID().String() != id {
		t.Error("ID deveria ser preservado na reconstrução")
	}
	if u.EstaAtiva() {
		t.Error("status suspensa não deveria ser ativa")
	}
}
