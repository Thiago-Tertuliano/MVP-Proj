package dto

import "encoding/json"

type CriarArtigoRequest struct {
	Titulo    string          `json:"titulo" validate:"required,min=3,max=300"`
	Subtitulo string          `json:"subtitulo" validate:"omitempty,max=500"`
	CapaURL   string          `json:"capa_url" validate:"omitempty,url"`
	Slug      string          `json:"slug" validate:"omitempty,min=2,max=200"`
	TrilhaID  *string         `json:"trilha_id"`
	ModuloID  *string         `json:"modulo_id"`
	Conteudo  json.RawMessage `json:"conteudo"`
	Metadados json.RawMessage `json:"metadados"`
}

type AtualizarArtigoRequest struct {
	Titulo    string          `json:"titulo" validate:"required,min=3,max=300"`
	Subtitulo string          `json:"subtitulo" validate:"omitempty,max=500"`
	CapaURL   string          `json:"capa_url" validate:"omitempty,url"`
	TrilhaID  *string         `json:"trilha_id"`
	ModuloID  *string         `json:"modulo_id"`
	Conteudo  json.RawMessage `json:"conteudo"`
	Metadados json.RawMessage `json:"metadados"`
}

type ArtigoResponse struct {
	ID          string          `json:"id"`
	Slug        string          `json:"slug"`
	Titulo      string          `json:"titulo"`
	Subtitulo   string          `json:"subtitulo,omitempty"`
	CapaURL     string          `json:"capa_url,omitempty"`
	TrilhaID    *string         `json:"trilha_id"`
	ModuloID    *string         `json:"modulo_id"`
	Conteudo    json.RawMessage `json:"conteudo"`
	Metadados   json.RawMessage `json:"metadados"`
	AutorID     string          `json:"autor_id"`
	Status      string          `json:"status"`
	PublicadoEm *int64          `json:"publicado_em,omitempty"`
	CreatedAt   int64           `json:"created_at"`
	UpdatedAt   int64           `json:"updated_at"`
}

type ListarArtigosResponse struct {
	Itens []*ArtigoResponse `json:"itens"`
}