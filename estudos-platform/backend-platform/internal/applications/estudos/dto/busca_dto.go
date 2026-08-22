package dto

type BuscaResponse struct {
	Itens []ResultadoBusca `json:"itens"`
}

type ResultadoBusca struct {
	Slug       string  `json:"slug"`
	Titulo     string  `json:"titulo"`
	Similarity float64 `json:"similarity"`
}
