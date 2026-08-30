package dto

import "encoding/json"

// PlanoImportacao é o resultado parseado das fontes (Courses.md + aula Go).
type PlanoImportacao struct {
	Trilhas []TrilhaImportacao `json:"trilhas"`
	Avisos  []string           `json:"avisos,omitempty"`
}

type TrilhaImportacao struct {
	Slug      string             `json:"slug"`
	Titulo    string             `json:"titulo"`
	Descricao string             `json:"descricao"`
	Ordem     int                `json:"ordem"`
	Modulos   []ModuloImportacao `json:"modulos"`
}

type ModuloImportacao struct {
	Slug      string           `json:"slug"`
	Titulo    string           `json:"titulo"`
	Descricao string           `json:"descricao"`
	Aulas     []AulaImportacao `json:"aulas"`
}

type AulaImportacao struct {
	Slug      string          `json:"slug"`
	Titulo    string          `json:"titulo"`
	Subtitulo string          `json:"subtitulo,omitempty"`
	Conteudo  json.RawMessage `json:"conteudo"`
	Metadados json.RawMessage `json:"metadados"`
}

type RelatorioImportacao struct {
	DryRun         bool     `json:"dry_run"`
	TrilhasCriadas int      `json:"trilhas_criadas"`
	TrilhasOK      int      `json:"trilhas_ok"`
	ArtigosCriados int      `json:"artigos_criados"`
	ArtigosOK      int      `json:"artigos_ok"`
	Avisos         []string `json:"avisos,omitempty"`
}
