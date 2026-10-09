package usecase

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

// PapelEditor resolve o papel efetivo do usuário. E-mails listados em EDITOR_EMAILS são
// promovidos a editor no primeiro acesso (persistido), sem precisar de passo manual no banco.
type PapelEditor struct {
	repo   repository.PapelRepository
	emails map[string]struct{}
}

func NewPapelEditor(repo repository.PapelRepository, emailsEditores []string) *PapelEditor {
	set := make(map[string]struct{}, len(emailsEditores))
	for _, e := range emailsEditores {
		e = strings.ToLower(strings.TrimSpace(e))
		if e != "" {
			set[e] = struct{}{}
		}
	}
	return &PapelEditor{repo: repo, emails: set}
}

// Papel devolve "aluno" ou "editor".
func (p *PapelEditor) Papel(ctx context.Context, usuarioID string) (string, error) {
	papel, email, err := p.repo.Obter(ctx, usuarioID)
	if err != nil {
		var de *errors.DomainError
		if stderrors.As(err, &de) {
			return "", err
		}
		return "", errors.ErrInternal("falha ao consultar papel", "PapelEditor.Papel", err)
	}
	if papel == repository.PapelEditor {
		return papel, nil
	}
	if _, ok := p.emails[strings.ToLower(email)]; ok {
		if err := p.repo.Definir(ctx, usuarioID, repository.PapelEditor); err != nil {
			return "", errors.ErrInternal("falha ao promover editor", "PapelEditor.Papel", err)
		}
		return repository.PapelEditor, nil
	}
	return repository.PapelAluno, nil
}

// EhEditor implementa middleware.EditorVerificador.
func (p *PapelEditor) EhEditor(ctx context.Context, usuarioID string) (bool, error) {
	papel, err := p.Papel(ctx, usuarioID)
	if err != nil {
		return false, err
	}
	return papel == repository.PapelEditor, nil
}
