import type { Artigo, Trilha } from "./types";

export const trilhas: Trilha[] = [
  {
    slug: "go-basico",
    titulo: "Go Básico",
    descricao: "Fundamentos de Go para ler o backend da plataforma.",
    publicada: true,
    progressoPct: 33,
    modulos: [
      {
        slug: "sintaxe",
        titulo: "Sintaxe",
        descricao: "Pacotes, nomes e tipos",
        artigos: [
          { slug: "pacotes-em-go", titulo: "Pacotes em Go", concluido: true },
          { slug: "structs-e-metodos", titulo: "Structs e métodos", concluido: false },
        ],
      },
      {
        slug: "interfaces",
        titulo: "Interfaces",
        descricao: "Contratos implícitos",
        artigos: [
          { slug: "interfaces-implicitas", titulo: "Interfaces implícitas", concluido: false },
        ],
      },
    ],
  },
  {
    slug: "dados",
    titulo: "Dados",
    descricao: "Engenharia de dados, cloud e plataformas.",
    publicada: false,
    progressoPct: 0,
    modulos: [
      {
        slug: "catalogo",
        titulo: "Catálogo",
        descricao: "Materiais do Courses.md",
        artigos: [
          { slug: "spark", titulo: "SPARK", concluido: false },
          { slug: "python", titulo: "Python", concluido: false },
        ],
      },
    ],
  },
];

export const artigos: Record<string, Artigo> = {
  "pacotes-em-go": {
    id: "55555555-5555-5555-5555-555555555555",
    slug: "pacotes-em-go",
    titulo: "Pacotes em Go",
    status: "publicado",
    trilhaSlug: "go-basico",
    moduloSlug: "sintaxe",
    conteudo: {
      blocks: [
        { type: "h", level: 2, text: "3.1 Pacotes e arquivos" },
        { type: "p", text: "Todo arquivo .go na pasta declara o mesmo package. Nome exportado começa com maiúscula." },
        {
          type: "code",
          lang: "go",
          text: `package entity

import (
    "fmt"
    "github.com/google/uuid"
    ".../internal/domain/shared/errors"
)`,
        },
        { type: "p", text: "Um diretório = um pacote (em geral). internal/ não é importável de fora do módulo." },
      ],
    },
    metadados: {
      tempo_leitura_min: 10,
      objetivo: "Revisar Pacotes em Go a partir da aula do repositório.",
      tags: ["go"],
      origem: "aula-code-review",
      fontes: [{ titulo: "Aula Code Review Go Sprint", url: "fontes/AULA-CODE-REVIEW-GO-SPRINT.md" }],
      quiz: {
        questoes: [
          {
            id: "q1",
            enunciado: "Este texto veio de um curso pago externo?",
            opcoes: [
              { id: "a", texto: "Sim" },
              { id: "b", texto: "Não — é o doc da aula no nosso Git" },
              { id: "c", texto: "Só o LinkedIn" },
            ],
            correta: "b",
            explicacao: "A fonte é fontes/AULA-CODE-REVIEW-GO-SPRINT.md.",
          },
          {
            id: "q2",
            enunciado: "Onde está o código de exemplo desta trilha?",
            opcoes: [
              { id: "a", texto: "hello world genérico" },
              { id: "b", texto: "Neste repositório (internal/, cmd/api)" },
              { id: "c", texto: "Somente no Courses.md" },
            ],
            correta: "b",
            explicacao: "A aula usa o backend da plataforma.",
          },
        ],
      },
    },
  },
  spark: {
    id: "a0000000-0000-0000-0000-000000000001",
    slug: "spark",
    titulo: "SPARK",
    status: "rascunho",
    trilhaSlug: "dados",
    moduloSlug: "catalogo",
    conteudo: {
      blocks: [
        { type: "h", level: 2, text: "Objetivo" },
        { type: "p", text: "Estudar o material indicado e conseguir explicar o tema em uma frase." },
        { type: "h", level: 2, text: "Material" },
        { type: "p", text: "Fonte: SPARK" },
        { type: "p", text: "https://lnkd.in/ddhHyvPH" },
        { type: "h", level: 2, text: "Checkpoint" },
        {
          type: "p",
          text: "1) Abri o link. 2) Anotei o assunto em uma frase. 3) Marquei como lido na plataforma.",
        },
      ],
    },
    metadados: {
      tempo_leitura_min: 10,
      objetivo: "Estudar SPARK a partir da fonte do catálogo",
      tags: ["dados"],
      origem: "courses-md",
      fontes: [{ titulo: "SPARK", url: "https://lnkd.in/ddhHyvPH" }],
      quiz: {
        questoes: [
          {
            id: "q1",
            enunciado: "Você abriu o material SPARK?",
            opcoes: [
              { id: "a", texto: "Sim" },
              { id: "b", texto: "Ainda não" },
              { id: "c", texto: "Link quebrado" },
            ],
            correta: "a",
            explicacao: "O checkpoint desta fase é acessar a fonte.",
          },
        ],
      },
    },
  },
};

export function getTrilha(slug: string): Trilha | undefined {
  return trilhas.find((t) => t.slug === slug);
}

export function getArtigo(slug: string): Artigo | undefined {
  return artigos[slug];
}

export function trilhasPublicadas(): Trilha[] {
  return trilhas.filter((t) => t.publicada);
}
