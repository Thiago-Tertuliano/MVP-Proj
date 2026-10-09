package repository

import "context"

const (
	PapelAluno  = "aluno"
	PapelEditor = "editor"
)

// PapelRepository lê e grava o papel do usuário. A autorização consulta o banco a cada
// requisição (não confia em claim do JWT), então revogar um editor vale imediatamente.
type PapelRepository interface {
	// Obter devolve o papel e o e-mail do usuário.
	Obter(ctx context.Context, usuarioID string) (papel, email string, err error)
	Definir(ctx context.Context, usuarioID, papel string) error
}
