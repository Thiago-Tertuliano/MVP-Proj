package content

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Block struct {
	Type  string   `json:"type"`
	Text  string   `json:"text,omitempty"`
	Level int      `json:"level,omitempty"`
	Lang  string   `json:"lang,omitempty"`
	Items []string `json:"items,omitempty"`
}

type Fonte struct {
	Titulo string `json:"titulo"`
	URL    string `json:"url"`
}

type QuizOpcao struct {
	ID    string `json:"id"`
	Texto string `json:"texto"`
}

type QuizQuestao struct {
	ID         string      `json:"id"`
	Enunciado  string      `json:"enunciado"`
	Opcoes     []QuizOpcao `json:"opcoes"`
	Correta    string      `json:"correta"`
	Explicacao string      `json:"explicacao"`
}

type MetadadosAula struct {
	TempoLeituraMin int          `json:"tempo_leitura_min"`
	Objetivo        string       `json:"objetivo"`
	Tags            []string     `json:"tags"`
	Origem          string       `json:"origem"`
	Fontes          []Fonte      `json:"fontes"`
	Quiz            QuizEnvelope `json:"quiz"`
}

type QuizEnvelope struct {
	Questoes []QuizQuestao `json:"questoes"`
}

func BlocksJSON(blocks []Block) (json.RawMessage, error) {
	type wrap struct {
		Blocks []Block `json:"blocks"`
	}
	b, err := json.Marshal(wrap{Blocks: blocks})
	if err != nil {
		return nil, err
	}
	return b, nil
}

func MetadadosJSON(m MetadadosAula) (json.RawMessage, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func StubCatalogo(titulo, url, tag string) ([]Block, MetadadosAula) {
	blocks := []Block{
		{Type: "h", Level: 2, Text: "Objetivo"},
		{Type: "p", Text: "Estudar o material indicado e conseguir explicar o tema em uma frase."},
		{Type: "h", Level: 2, Text: "Material"},
		{Type: "p", Text: "Fonte: " + titulo},
		{Type: "p", Text: url},
		{Type: "h", Level: 2, Text: "Checkpoint"},
		{Type: "p", Text: "1) Abri o link. 2) Anotei o assunto em uma frase. 3) Marquei como lido na plataforma."},
	}
	meta := MetadadosAula{
		TempoLeituraMin: 10,
		Objetivo:        fmt.Sprintf("Estudar %s a partir da fonte do catálogo", titulo),
		Tags:            []string{tag},
		Origem:          "courses-md",
		Fontes:          []Fonte{{Titulo: titulo, URL: url}},
		Quiz:            QuizEnvelope{Questoes: quizCatalogo(titulo)},
	}
	return blocks, meta
}

func quizCatalogo(titulo string) []QuizQuestao {
	return []QuizQuestao{
		{
			ID:        "q1",
			Enunciado: fmt.Sprintf("Você abriu o material %s?", titulo),
			Opcoes: []QuizOpcao{
				{ID: "a", Texto: "Sim"},
				{ID: "b", Texto: "Ainda não"},
				{ID: "c", Texto: "Link quebrado"},
			},
			Correta:    "a",
			Explicacao: "O checkpoint desta fase é acessar a fonte. Conteúdo de curso pago não é copiado.",
		},
		{
			ID:        "q2",
			Enunciado: "Você consegue dizer o assunto principal em uma frase?",
			Opcoes: []QuizOpcao{
				{ID: "a", Texto: "Sim"},
				{ID: "b", Texto: "Ainda não"},
				{ID: "c", Texto: "Não era o tema desta trilha"},
			},
			Correta:    "a",
			Explicacao: "O catálogo classifica pelo assunto principal.",
		},
		{
			ID:        "q3",
			Enunciado: "Este nó substitui o curso original?",
			Opcoes: []QuizOpcao{
				{ID: "a", Texto: "Sim, a aula está toda aqui"},
				{ID: "b", Texto: "Não — o link é a aula; aqui é o mapa e o progresso"},
				{ID: "c", Texto: "Só se pagar API de IA"},
			},
			Correta:    "b",
			Explicacao: "A plataforma organiza o caminho. O material continua no link.",
		},
	}
}

func MarkdownParaBlocks(md string) []Block {
	var blocks []Block
	lines := strings.Split(strings.ReplaceAll(md, "\r\n", "\n"), "\n")
	var para []string
	flushP := func() {
		t := strings.TrimSpace(strings.Join(para, " "))
		if t != "" {
			blocks = append(blocks, Block{Type: "p", Text: t})
		}
		para = nil
	}

	inCode := false
	var code []string
	lang := ""
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "```") {
			if !inCode {
				flushP()
				inCode = true
				lang = strings.TrimSpace(strings.TrimPrefix(trim, "```"))
				code = nil
				continue
			}
			blocks = append(blocks, Block{Type: "code", Lang: lang, Text: strings.Join(code, "\n")})
			inCode = false
			continue
		}
		if inCode {
			code = append(code, line)
			continue
		}
		if strings.HasPrefix(trim, "### ") {
			flushP()
			blocks = append(blocks, Block{Type: "h", Level: 3, Text: strings.TrimSpace(trim[4:])})
			continue
		}
		if strings.HasPrefix(trim, "## ") {
			flushP()
			blocks = append(blocks, Block{Type: "h", Level: 2, Text: strings.TrimSpace(trim[3:])})
			continue
		}
		if trim == "" || trim == "---" {
			flushP()
			continue
		}
		para = append(para, trim)
	}
	flushP()
	return blocks
}
