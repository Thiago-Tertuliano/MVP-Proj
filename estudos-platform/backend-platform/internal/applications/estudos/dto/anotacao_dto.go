package dto

import "encoding/json"

type SalvarAnotacaoRequest struct {
	// Vou ajustar isso depois, mas aceita o JSON livre no frontend, então não precisa de validação por enquanto - =) Thigas
	Conteudo json.RawMessage `json:"conteudo" validate:"required"`
}

type AnotacaoResponse struct {
	ArtigoID string          `json:"artigo_id"`
	Conteudo json.RawMessage `json:"conteudo"`
}