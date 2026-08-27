package port

import "context"

// EmbeddingGerador gera vetor do texto do artigo. Stub local devolve nil (busca cai no ILIKE).
type EmbeddingGerador interface {
	Gerar(ctx context.Context, texto string) ([]float32, error)
}
