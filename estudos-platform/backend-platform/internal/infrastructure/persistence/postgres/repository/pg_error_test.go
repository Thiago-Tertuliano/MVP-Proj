package repository

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func TestMapPG_Nil(t *testing.T) {
	if err := MapPG(nil, "op"); err != nil {
		t.Fatalf("nil deve passar, got %v", err)
	}
}

func TestMapPG_ErroComumNaoTraduz(t *testing.T) {
	base := errors.New("rede caiu")
	got := MapPG(base, "op")
	if got != base {
		t.Fatalf("erro sem PgError deve voltar igual, got %v", got)
	}
}

func TestMapPG_TabelaAusente(t *testing.T) {
	pgErr := &pgconn.PgError{Code: pgUndefinedTable}
	got := MapPG(pgErr, "usuario_repo_pg.EmailExiste")

	var de *domainErros.DomainError
	if !errors.As(got, &de) {
		t.Fatalf("esperava DomainError, got %#v", got)
	}
	if de.Kind != domainErros.Internal {
		t.Errorf("Kind=%s, queria INTERNAL", de.Kind)
	}
	if de.Message == "" {
		t.Fatal("mensagem vazia — o Bruno precisa do texto de migrate")
	}
}

func TestMapPG_OutroCodigoPostgres(t *testing.T) {
	pgErr := &pgconn.PgError{Code: pgUniqueViolation}
	got := MapPG(pgErr, "op")
	if got != pgErr {
		t.Fatalf("23505 não é papel do MapPG, got %v", got)
	}

	fk := &pgconn.PgError{Code: pgForeignKey}
	if MapPG(fk, "op") != fk {
		t.Fatal("23503 deve passar cru")
	}
}
