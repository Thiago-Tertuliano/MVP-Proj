package dto

type MarcarArtigoLidoRequest struct {
	Concluido bool `json:"concluido"`
}

type ProgressoArtigoResponse struct {
	ArtigoID  string `json:"artigo_id"`
	Concluido bool   `json:"concluido"`
	// Gamificacao traz XP/conquistas gerados por esta leitura (nós de roadmap concluídos).
	Gamificacao *ResultadoGamificacao `json:"gamificacao,omitempty"`
}

type ProgressoTrilhaResponse struct {
	TrilhaID   string  `json:"trilha_id"`
	Concluidos int     `json:"concluidos"`
	Total      int     `json:"total"`
	Percentual float64 `json:"percentual"`
	// ArtigosConcluidos lista os IDs dos artigos lidos: o front pinta cada nó do mapa com isso.
	ArtigosConcluidos []string `json:"artigos_concluidos"`
}

// ContinuarItem é o último artigo com que o aluno interagiu.
type ContinuarItem struct {
	ArtigoID     string  `json:"artigo_id"`
	ArtigoSlug   string  `json:"artigo_slug"`
	ArtigoTitulo string  `json:"artigo_titulo"`
	TrilhaID     *string `json:"trilha_id"`
	TrilhaSlug   *string `json:"trilha_slug"`
	TrilhaTitulo *string `json:"trilha_titulo"`
	Concluido    bool    `json:"concluido"`
	AtualizadoEm int64   `json:"atualizado_em"`
}

// ContinuarResponse traz Item nulo quando o aluno ainda não estudou nenhum artigo.
type ContinuarResponse struct {
	Item *ContinuarItem `json:"item"`
}
