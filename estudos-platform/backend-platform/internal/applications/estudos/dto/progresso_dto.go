package dto

type MarcarArtigoLidoRequest struct {
	Concluido bool `json:"concluido"`
}

type ProgressoArtigoResponse struct {
	ArtigoID  string `json:"artigo_id"`
	Concluido bool   `json:"concluido"`
}

type ProgressoTrilhaResponse struct {
	TrilhaID   string  `json:"trilha_id"`
	Concluidos int     `json:"concluidos"`
	Total      int     `json:"total"`
	Percentual float64 `json:"percentual"`
}
