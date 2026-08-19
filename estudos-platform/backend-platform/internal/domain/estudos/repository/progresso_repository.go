package repository

import "context"

type ProgressoArtigo struct {
	UsuarioID string
	ArtigoID string
	TrilhaID *string
	Concluido bool
}

type ProgressoRepository interface {
	UpsertArtigo(ctx context.Context, p ProgressoArtigo) error
	CountConcluidosNaTrilha(ctx context.Context, usuarioID, trilhaID string) (concluidos, total int, err error)
}