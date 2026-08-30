package content

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type alvoGo struct {
	Slug     string
	Modulo   string
	Titulo   string
	Keywords []string
}

var alvosGo = []alvoGo{
	{Slug: "pacotes-em-go", Modulo: "sintaxe", Titulo: "Pacotes em Go", Keywords: []string{"pacote", "package", "internal/"}},
	{Slug: "structs-e-metodos", Modulo: "sintaxe", Titulo: "Structs e métodos", Keywords: []string{"struct", "método", "metodo", "receiver"}},
	{Slug: "interfaces-implicitas", Modulo: "interfaces", Titulo: "Interfaces implícitas", Keywords: []string{"interface"}},
}

var secaoRE = regexp.MustCompile(`(?m)^(#{2,3})\s+(.+)$`)

func ParseAulaGo(md string) (*dto.TrilhaImportacao, error) {
	md = strings.ReplaceAll(md, "\r\n", "\n")
	if strings.TrimSpace(md) == "" {
		return nil, errors.ErrInvalidArgument("aula Go vazia", "content.ParseAulaGo", nil)
	}

	type secao struct {
		titulo string
		corpo  string
	}
	locs := secaoRE.FindAllStringSubmatchIndex(md, -1)
	var secoes []secao
	for i, loc := range locs {
		titulo := strings.TrimSpace(md[loc[4]:loc[5]])
		start := loc[1]
		end := len(md)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		secoes = append(secoes, secao{titulo: titulo, corpo: strings.TrimSpace(md[start:end])})
	}

	porSlug := map[string][]string{}
	for _, s := range secoes {
		t := strings.ToLower(s.titulo)
		corpo := strings.TrimSpace("## " + s.titulo + "\n\n" + s.corpo)
		for _, a := range alvosGo {
			if casaKeywords(t, a.Keywords) {
				porSlug[a.Slug] = append(porSlug[a.Slug], corpo)
			}
		}
	}

	modulos := map[string]*dto.ModuloImportacao{
		"sintaxe":    {Slug: "sintaxe", Titulo: "Sintaxe", Descricao: "Pacotes, nomes e tipos"},
		"interfaces": {Slug: "interfaces", Titulo: "Interfaces", Descricao: "Contratos implícitos"},
	}

	for _, a := range alvosGo {
		partes := porSlug[a.Slug]
		if len(partes) == 0 {
			return nil, errors.ErrInvalidArgument(
				fmt.Sprintf("nenhuma seção da aula Go casou com %s", a.Slug),
				"content.ParseAulaGo",
				nil,
			)
		}
		mdAula := strings.Join(partes, "\n\n")
		blocks := MarkdownParaBlocks(mdAula)
		if len(blocks) == 0 {
			return nil, errors.ErrInvalidArgument("seção vazia para "+a.Slug, "content.ParseAulaGo", nil)
		}
		conteudo, err := BlocksJSON(blocks)
		if err != nil {
			return nil, err
		}
		meta := MetadadosAula{
			TempoLeituraMin: 10,
			Objetivo:        "Revisar " + a.Titulo + " a partir da aula do repositório.",
			Tags:            []string{"go"},
			Origem:          "aula-code-review",
			Fontes:          []Fonte{{Titulo: "Aula Code Review Go Sprint", URL: "fontes/AULA-CODE-REVIEW-GO-SPRINT.md"}},
			Quiz: QuizEnvelope{Questoes: []QuizQuestao{
				{
					ID:        "q1",
					Enunciado: "Este texto veio de um curso pago externo?",
					Opcoes: []QuizOpcao{
						{ID: "a", Texto: "Sim"},
						{ID: "b", Texto: "Não — é o doc da aula no nosso Git"},
						{ID: "c", Texto: "Só o LinkedIn"},
					},
					Correta:    "b",
					Explicacao: "A fonte é fontes/AULA-CODE-REVIEW-GO-SPRINT.md.",
				},
				{
					ID:        "q2",
					Enunciado: "Onde está o código de exemplo desta trilha?",
					Opcoes: []QuizOpcao{
						{ID: "a", Texto: "hello world genérico"},
						{ID: "b", Texto: "Neste repositório (internal/, cmd/api)"},
						{ID: "c", Texto: "Somente no Courses.md"},
					},
					Correta:    "b",
					Explicacao: "A aula usa o backend da plataforma.",
				},
				{
					ID:        "q3",
					Enunciado: "Depois de ler, o próximo passo na plataforma é?",
					Opcoes: []QuizOpcao{
						{ID: "a", Texto: "Pedir para o job publicar sozinho"},
						{ID: "b", Texto: "Marcar o artigo como lido"},
						{ID: "c", Texto: "Apagar o seed"},
					},
					Correta:    "b",
					Explicacao: "Publicar é pela API; progresso é PUT /progresso/artigos/{id}.",
				},
			}},
		}
		metadados, err := MetadadosJSON(meta)
		if err != nil {
			return nil, err
		}
		mod := modulos[a.Modulo]
		mod.Aulas = append(mod.Aulas, dto.AulaImportacao{
			Slug:      a.Slug,
			Titulo:    a.Titulo,
			Conteudo:  conteudo,
			Metadados: metadados,
		})
	}

	return &dto.TrilhaImportacao{
		Slug:      "go-basico",
		Titulo:    "Go Básico",
		Descricao: "Fundamentos de Go para ler o backend da plataforma.",
		Ordem:     0,
		Modulos:   []dto.ModuloImportacao{*modulos["sintaxe"], *modulos["interfaces"]},
	}, nil
}

func casaKeywords(titulo string, keys []string) bool {
	for _, k := range keys {
		if strings.Contains(titulo, strings.ToLower(k)) {
			return true
		}
	}
	return false
}
