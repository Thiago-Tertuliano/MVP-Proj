package external

import "context"

// StubEmbedding não chama OpenAI: embedding fica nulo e a busca usa título (ILIKE).
type StubEmbedding struct{}

func (StubEmbedding) Gerar(context.Context, string) ([]float32, error) {
	return nil, nil
}
