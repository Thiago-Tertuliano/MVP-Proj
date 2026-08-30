package content

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type trilhaSpec struct {
	Slug      string
	Titulo    string
	Descricao string
	Match     []string
}

var trilhasCatalogo = []trilhaSpec{
	{Slug: "dados", Titulo: "Dados", Descricao: "Engenharia de dados, cloud e plataformas.", Match: []string{"dados"}},
	{Slug: "sre-infra", Titulo: "SRE / Infraestrutura", Descricao: "Linux, DevOps, containers e GitOps.", Match: []string{"sre", "infraestrutura"}},
	{Slug: "desenvolvimento", Titulo: "Desenvolvimento", Descricao: "Linguagens, APIs e práticas de engenharia.", Match: []string{"desenvolvimento"}},
	{Slug: "sap", Titulo: "SAP", Descricao: "Módulos SAP.", Match: []string{"sap"}},
	{Slug: "ferramentas", Titulo: "Ferramentas", Descricao: "Ferramentas gerais de trabalho técnico.", Match: []string{"ferramentas"}},
}

var (
	headingRE = regexp.MustCompile(`(?m)^##\s+(.+)$`)
	linkRE    = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	rowRE     = regexp.MustCompile(`^\|?\s*(.*?)\s*\|\s*(.*?)\s*\|?\s*$`)
)

func ParseCourses(md string) (*dto.PlanoImportacao, error) {
	md = strings.ReplaceAll(md, "\r\n", "\n")
	if !strings.Contains(md, "|") || !linkRE.MatchString(md) {
		return nil, errors.ErrInvalidArgument("Courses.md sem tabela parseável", "content.ParseCourses", nil)
	}

	type secao struct {
		titulo string
		corpo  string
	}
	var secoes []secao
	idxs := headingRE.FindAllStringSubmatchIndex(md, -1)
	if len(idxs) == 0 {
		return nil, errors.ErrInvalidArgument("Courses.md sem headings ##", "content.ParseCourses", nil)
	}
	for i, loc := range idxs {
		titulo := strings.TrimSpace(md[loc[2]:loc[3]])
		start := loc[1]
		end := len(md)
		if i+1 < len(idxs) {
			end = idxs[i+1][0]
		}
		secoes = append(secoes, secao{titulo: titulo, corpo: md[start:end]})
	}

	plano := &dto.PlanoImportacao{}
	vistos := map[string]string{} // slug artigo -> trilha
	ordem := 0

	for _, spec := range trilhasCatalogo {
		var corpo string
		for _, s := range secoes {
			if headingCasa(s.titulo, spec.Match) {
				corpo = s.corpo
				break
			}
		}
		if corpo == "" {
			continue
		}
		itens, avisos, err := parseTabela(corpo, spec.Slug, vistos)
		if err != nil {
			return nil, err
		}
		if len(itens) == 0 {
			continue
		}
		plano.Avisos = append(plano.Avisos, avisos...)
		plano.Trilhas = append(plano.Trilhas, dto.TrilhaImportacao{
			Slug:      spec.Slug,
			Titulo:    spec.Titulo,
			Descricao: spec.Descricao,
			Ordem:     ordem,
			Modulos: []dto.ModuloImportacao{{
				Slug:      "catalogo",
				Titulo:    "Catálogo",
				Descricao: "Materiais do Courses.md",
				Aulas:     itens,
			}},
		})
		ordem++
	}

	if len(plano.Trilhas) == 0 {
		return nil, errors.ErrInvalidArgument("nenhuma trilha extraída do Courses.md", "content.ParseCourses", nil)
	}
	return plano, nil
}

func isCabecalhoTabela(col1 string) bool {
	if col1 == "" || strings.Contains(col1, "---") {
		return true
	}
	for _, h := range []string{"Material", "Módulo", "Modulo", "Ferramenta", "Link"} {
		if strings.EqualFold(col1, h) {
			return true
		}
	}
	return false
}

func headingCasa(titulo string, keys []string) bool {
	n := normaliza(titulo)
	if strings.Contains(n, "categorias") {
		return false
	}
	for _, k := range keys {
		if strings.Contains(n, k) {
			return true
		}
	}
	return false
}

func normaliza(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == ' ' || r == '/' {
			b.WriteRune(r)
		}
	}
	return strings.TrimSpace(b.String())
}

func parseTabela(corpo, tag string, vistos map[string]string) ([]dto.AulaImportacao, []string, error) {
	var aulas []dto.AulaImportacao
	var avisos []string
	for _, line := range strings.Split(corpo, "\n") {
		line = strings.TrimSpace(line)
		m := rowRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		col1, col2 := strings.TrimSpace(m[1]), strings.TrimSpace(m[2])
		if isCabecalhoTabela(col1) {
			continue
		}
		lm := linkRE.FindStringSubmatch(col2)
		if lm == nil {
			lm = linkRE.FindStringSubmatch(col1)
		}
		if lm == nil {
			if strings.Contains(col2, "http") || strings.Contains(col1, "http") {
				return nil, nil, errors.ErrInvalidArgument("linha de tabela sem URL markdown: "+col1, "content.ParseCourses", nil)
			}
			continue
		}
		titulo := strings.TrimSpace(col1)
		titulo = linkRE.ReplaceAllString(titulo, "$1")
		titulo = strings.TrimSpace(titulo)
		if titulo == "" || strings.EqualFold(titulo, "Acessar") {
			titulo = strings.TrimSpace(lm[1])
		}
		if len([]rune(titulo)) < 3 {
			titulo = strings.ToUpper(tag) + " " + titulo
		}
		url := strings.TrimSpace(lm[2])
		if url == "" {
			return nil, nil, errors.ErrInvalidArgument("URL vazia em "+titulo, "content.ParseCourses", nil)
		}
		slugVO, err := valueobject.NewSlug(titulo)
		if err != nil {
			return nil, nil, errors.ErrInvalidArgument(fmt.Sprintf("slug inválido para %q", titulo), "content.ParseCourses", err)
		}
		slug := slugVO.Value()
		if outra, ok := vistos[slug]; ok {
			avisos = append(avisos, fmt.Sprintf("slug %s já usado na trilha %s — reusa, não duplica", slug, outra))
			continue
		}
		blocks, meta := StubCatalogo(titulo, url, tag)
		conteudo, err := BlocksJSON(blocks)
		if err != nil {
			return nil, nil, err
		}
		metadados, err := MetadadosJSON(meta)
		if err != nil {
			return nil, nil, err
		}
		aulas = append(aulas, dto.AulaImportacao{
			Slug:      slug,
			Titulo:    titulo,
			Conteudo:  conteudo,
			Metadados: metadados,
		})
		vistos[slug] = tag
	}
	return aulas, avisos, nil
}
