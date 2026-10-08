package repository

import (
	"context"
	"time"
)

type ProgressoArtigo struct {
	UsuarioID string
	ArtigoID  string
	TrilhaID  *string
	Concluido bool
}

// UltimoProgresso descreve o artigo publicado com que o aluno interagiu por último
// (lido ou em curso). Alimenta o card "Continuar de onde parou" do front.
type UltimoProgresso struct {
	ArtigoID     string
	ArtigoSlug   string
	ArtigoTitulo string
	TrilhaID     *string
	TrilhaSlug   *string
	TrilhaTitulo *string
	Concluido    bool
	AtualizadoEm time.Time
}

type ProgressoRepository interface {
	UpsertArtigo(ctx context.Context, p ProgressoArtigo) error
	CountConcluidosNaTrilha(ctx context.Context, usuarioID, trilhaID string) (concluidos, total int, err error)
	// ListarConcluidosNaTrilha devolve os IDs dos artigos publicados da trilha que o aluno já concluiu.
	ListarConcluidosNaTrilha(ctx context.Context, usuarioID, trilhaID string) ([]string, error)
	// UltimoArtigo devolve a interação mais recente do aluno ou nil se ele ainda não estudou nada.
	UltimoArtigo(ctx context.Context, usuarioID string) (*UltimoProgresso, error)
}
