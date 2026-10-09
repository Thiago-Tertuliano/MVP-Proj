package usecase

import (
	"context"
	"errors"
	"testing"

	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type fakePapelRepo struct {
	papel, email string
	err          error
	definido     string
}

func (f *fakePapelRepo) Obter(context.Context, string) (string, string, error) {
	return f.papel, f.email, f.err
}
func (f *fakePapelRepo) Definir(_ context.Context, _ string, papel string) error {
	f.definido = papel
	f.papel = papel
	return nil
}

func TestPapelEditor(t *testing.T) {
	ctx := context.Background()

	t.Run("aluno comum", func(t *testing.T) {
		repo := &fakePapelRepo{papel: "aluno", email: "a@x.com"}
		ok, err := NewPapelEditor(repo, []string{"editor@x.com"}).EhEditor(ctx, "u")
		if err != nil || ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if repo.definido != "" {
			t.Fatal("não deveria promover")
		}
	})

	t.Run("já editor no banco", func(t *testing.T) {
		repo := &fakePapelRepo{papel: "editor", email: "a@x.com"}
		ok, err := NewPapelEditor(repo, nil).EhEditor(ctx, "u")
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
	})

	t.Run("e-mail da lista é promovido (case-insensitive) e persistido", func(t *testing.T) {
		repo := &fakePapelRepo{papel: "aluno", email: "Editor@X.com"}
		papel, err := NewPapelEditor(repo, []string{" editor@x.com "}).Papel(ctx, "u")
		if err != nil || papel != "editor" {
			t.Fatalf("papel=%s err=%v", papel, err)
		}
		if repo.definido != "editor" {
			t.Fatal("deveria persistir a promoção")
		}
	})

	t.Run("erro de banco vira interno", func(t *testing.T) {
		repo := &fakePapelRepo{err: errors.New("boom")}
		_, err := NewPapelEditor(repo, nil).EhEditor(ctx, "u")
		var de *domainErros.DomainError
		if !errors.As(err, &de) || de.Kind != domainErros.Internal {
			t.Fatalf("esperava Internal, got %v", err)
		}
	})

	t.Run("usuário inexistente mantém NotFound", func(t *testing.T) {
		repo := &fakePapelRepo{err: domainErros.ErrNotFound("x", "t", nil)}
		_, err := NewPapelEditor(repo, nil).EhEditor(ctx, "u")
		var de *domainErros.DomainError
		if !errors.As(err, &de) || de.Kind != domainErros.NotFound {
			t.Fatalf("esperava NotFound, got %v", err)
		}
	})
}
