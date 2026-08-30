package repository

import (
	"context"
	"encoding/json"
)

type Anotacao struct {
	UsuarioID string 
	ArtigoID  string
	Conteudo  json.RawMessage
}

type AnotacaoRepository interface {
	Upsert(ctx context.Context, a Anotacao) error
	FindByUsuarioEArtigo(ctx context.Context, usuarioID, artigoID string) (*Anotacao, error)
}