package dto

type MarcarArtigoLidoRequest struct {
	Concluido bool `json:"concluido"`
}

type ProgressoArtigoResponse struct {
	ArtigoID  string `json:"artigo_id"`
	Concluido bool   `json:"concluido"`
}
