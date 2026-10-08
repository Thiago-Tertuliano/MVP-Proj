/**
 * Dados de exemplo SÓ para Storybook/testes visuais. Espelham o contrato da API (snake_case).
 * O app (`app/`, `components/` fora de `*.stories.tsx`) nunca importa daqui — ver regra de lint.
 */
import type { Artigo, Block, QuizQuestao, Trilha } from "@/lib/types";
import { montarRoadmap } from "@/lib/roadmap";

const AGORA = 1_760_000_000;

export const trilhaGo: Trilha = {
  id: "22222222-2222-2222-2222-222222222222",
  slug: "go-basico",
  titulo: "Go Básico",
  descricao: "Fundamentos de Go para ler o backend da plataforma.",
  ordem: 1,
  publicada: true,
  created_at: AGORA,
  updated_at: AGORA,
  modulos: [
    { id: "m1", slug: "sintaxe", titulo: "Sintaxe", descricao: "Pacotes, nomes e tipos", ordem: 1 },
    { id: "m2", slug: "interfaces", titulo: "Interfaces", descricao: "Contratos implícitos", ordem: 2 },
  ],
};

export const blocosExemplo: Block[] = [
  { type: "h", level: 2, text: "Pacotes e arquivos" },
  {
    type: "p",
    text: "Todo arquivo `.go` na pasta declara o mesmo **package**. Nome exportado começa com maiúscula.",
  },
  {
    type: "code",
    lang: "go",
    text: `package entity

import (
    "fmt"
)

// Saudar imprime uma mensagem.
func Saudar(nome string) {
    fmt.Println("Olá,", nome, 42)
}`,
  },
  { type: "p", text: "- Um diretório = um pacote (em geral)" },
  { type: "p", text: "- `internal/` não é importável de fora do módulo" },
  { type: "callout", variant: "tip", text: "Rode `go vet ./...` antes de abrir o PR." },
  { type: "callout", variant: "warning", title: "Cuidado", text: "Nomes minúsculos não são exportados." },
  { type: "callout", variant: "info", text: "Pacotes de teste terminam em `_test`." },
  { type: "list", ordered: true, items: ["Crie o módulo", "Escreva o `main`", "Rode `go run .`"] },
  { type: "link", url: "https://go.dev/doc/", text: "Documentação oficial" },
  { type: "p", text: "https://go.dev/tour/" },
];

export const questoesExemplo: QuizQuestao[] = [
  {
    id: "q1",
    enunciado: "Onde está o código de exemplo desta trilha?",
    opcoes: [
      { id: "a", texto: "Hello world genérico" },
      { id: "b", texto: "Neste repositório (internal/, cmd/api)" },
      { id: "c", texto: "Somente no Courses.md" },
    ],
    correta: "b",
    explicacao: "A aula usa o backend da plataforma.",
  },
];

function artigo(i: number, slug: string, titulo: string, modulo_id: string): Artigo {
  return {
    id: `a000000${i}-0000-0000-0000-000000000000`,
    slug,
    titulo,
    trilha_id: trilhaGo.id,
    modulo_id,
    conteudo: { blocks: blocosExemplo },
    metadados: { tempo_leitura_min: 10, objetivo: `Revisar ${titulo}.`, quiz: { questoes: questoesExemplo } },
    autor_id: "u1",
    status: "publicado",
    created_at: AGORA,
    updated_at: AGORA,
  };
}

export const artigosGo: Artigo[] = [
  artigo(1, "pacotes-em-go", "Pacotes em Go", "m1"),
  artigo(2, "structs-e-metodos", "Structs e métodos", "m1"),
  artigo(3, "interfaces-implicitas", "Interfaces implícitas", "m2"),
];

/** Mapa de um aluno que leu o 1º artigo (2º = atual, demais = futuros). */
export const roadmapLogado = montarRoadmap(trilhaGo, artigosGo, new Set([artigosGo[0].id]));
/** Mapa de visitante (sem estados). */
export const roadmapVisitante = montarRoadmap(trilhaGo, artigosGo, null);
