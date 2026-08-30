package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseCourses_ExtraiTrilhasENaoDuplicaPython(t *testing.T) {
	md := `# Biblioteca

## Categorias
* [Dados](#dados)

## Dados

| Material | Link |
| --- | --- |
| SPARK | [Acessar](https://example.com/spark) |
| Python | [Acessar](https://example.com/py) |

## Desenvolvimento

| Material | Link |
| -- | -- |
| Python | [Acessar](https://example.com/py) |
| API REST | [Acessar](https://example.com/api) |
`
	plano, err := ParseCourses(md)
	if err != nil {
		t.Fatal(err)
	}
	if len(plano.Trilhas) != 2 {
		t.Fatalf("trilhas=%d avisos=%v", len(plano.Trilhas), plano.Avisos)
	}
	if plano.Trilhas[0].Slug != "dados" || len(plano.Trilhas[0].Modulos[0].Aulas) != 2 {
		t.Fatalf("dados: %+v", plano.Trilhas[0])
	}
	if plano.Trilhas[1].Slug != "desenvolvimento" {
		t.Fatalf("segunda trilha: %s", plano.Trilhas[1].Slug)
	}
	if len(plano.Trilhas[1].Modulos[0].Aulas) != 1 {
		t.Fatalf("python duplicado deveria ser pulado, aulas=%d avisos=%v", len(plano.Trilhas[1].Modulos[0].Aulas), plano.Avisos)
	}
	if !strings.Contains(strings.Join(plano.Avisos, " "), "python") {
		t.Fatalf("esperava aviso de slug reusado: %v", plano.Avisos)
	}
}

func TestParseCourses_URLVazia(t *testing.T) {
	md := `## Dados
| Material | Link |
| --- | --- |
| SPARK | [Acessar]() |
`
	_, err := ParseCourses(md)
	if err == nil {
		t.Fatal("esperava erro de URL vazia")
	}
}

func TestParseCourses_SemTabela(t *testing.T) {
	_, err := ParseCourses("# só texto")
	if err == nil {
		t.Fatal("esperava erro")
	}
}

func TestParseAulaGo_AssociaPorKeyword(t *testing.T) {
	md := `
## 3.1 Pacotes e arquivos

package entity na pasta.

## 3.2 Struct, campos e métodos

Receiver por ponteiro.

## Interfaces implícitas

Satisfaz sem implements.
`
	trilha, err := ParseAulaGo(md)
	if err != nil {
		t.Fatal(err)
	}
	if trilha.Slug != "go-basico" || len(trilha.Modulos) != 2 {
		t.Fatalf("%+v", trilha)
	}
	sintaxe := trilha.Modulos[0].Aulas
	if len(sintaxe) != 2 || sintaxe[0].Slug != "pacotes-em-go" {
		t.Fatalf("sintaxe: %+v", sintaxe)
	}
	if trilha.Modulos[1].Aulas[0].Slug != "interfaces-implicitas" {
		t.Fatalf("interfaces: %+v", trilha.Modulos[1].Aulas)
	}
}

func TestParseAulaGo_FalhaSeNaoCasar(t *testing.T) {
	_, err := ParseAulaGo("## Só título sem keyword útil\n\ntexto")
	if err == nil {
		t.Fatal("esperava erro")
	}
}

func TestParseCourses_LinhaSemPipeInicialETituloCurto(t *testing.T) {
	md := `## SRE / Infraestrutura
| Material | Link |
| --- | --- |
Linux  |[Acessar](https://example.com/linux) |

## SAP
| Módulo | Link |
| --- | --- |
| RE | [Acessar](https://example.com/re) |
`
	plano, err := ParseCourses(md)
	if err != nil {
		t.Fatal(err)
	}
	if len(plano.Trilhas) != 2 {
		t.Fatalf("trilhas=%d", len(plano.Trilhas))
	}
	linux := plano.Trilhas[0].Modulos[0].Aulas
	if len(linux) != 1 || linux[0].Slug != "linux" {
		t.Fatalf("linux: %+v", linux)
	}
	sap := plano.Trilhas[1].Modulos[0].Aulas
	if len(sap) != 1 || sap[0].Slug != "sap-re" {
		t.Fatalf("título curto deveria virar SAP RE: %+v", sap)
	}
}

func TestParseCourses_ArquivoReal(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "fontes", "Courses.md"))
	if err != nil {
		t.Fatal(err)
	}
	plano, err := ParseCourses(string(raw))
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]int{}
	for _, tr := range plano.Trilhas {
		n := 0
		for _, m := range tr.Modulos {
			n += len(m.Aulas)
		}
		got[tr.Slug] = n
	}
	want := []string{"dados", "sre-infra", "desenvolvimento", "sap", "ferramentas"}
	if len(plano.Trilhas) != len(want) {
		t.Fatalf("trilhas=%v", got)
	}
	for i, slug := range want {
		if plano.Trilhas[i].Slug != slug {
			t.Fatalf("ordem: want %s got %s", slug, plano.Trilhas[i].Slug)
		}
		if got[slug] == 0 {
			t.Fatalf("trilha %s sem aulas", slug)
		}
	}
	if got["dados"] < 10 || got["sap"] < 10 {
		t.Fatalf("catálogo curto demais: %v", got)
	}
}

func TestParseAulaGo_ArquivoReal(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "fontes", "AULA-CODE-REVIEW-GO-SPRINT.md"))
	if err != nil {
		t.Fatal(err)
	}
	trilha, err := ParseAulaGo(string(raw))
	if err != nil {
		t.Fatal(err)
	}
	if trilha.Slug != "go-basico" || len(trilha.Modulos) != 2 {
		t.Fatalf("%+v", trilha)
	}
	if len(trilha.Modulos[0].Aulas) != 2 || trilha.Modulos[1].Aulas[0].Slug != "interfaces-implicitas" {
		t.Fatalf("módulos: %+v", trilha.Modulos)
	}
	if !strings.Contains(string(trilha.Modulos[0].Aulas[0].Conteudo), "Pacotes") {
		t.Fatal("pacotes-em-go deveria trazer a seção da aula")
	}
}

func TestMarkdownParaBlocks_CodeFence(t *testing.T) {
	blocks := MarkdownParaBlocks("## Objetivo\n\nOlá\n\n```go\npackage x\n```\n")
	if len(blocks) < 3 {
		t.Fatalf("%+v", blocks)
	}
	if blocks[0].Type != "h" || blocks[1].Type != "p" || blocks[2].Type != "code" {
		t.Fatalf("%+v", blocks)
	}
}
