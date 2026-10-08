import { describe, expect, it } from "vitest";

import { agruparPorModulo, artigosEmOrdem, montarRoadmap, proximoNaoLido, vizinhos } from "@/lib/roadmap";
import type { Artigo, Trilha } from "@/lib/types";

const trilha = {
  id: "t1",
  slug: "go-basico",
  titulo: "Go Básico",
  publicada: true,
  ordem: 0,
  created_at: 0,
  updated_at: 0,
  modulos: [
    { id: "m2", slug: "interfaces", titulo: "Interfaces", ordem: 2 },
    { id: "m1", slug: "sintaxe", titulo: "Sintaxe", ordem: 1 },
  ],
} satisfies Trilha;

function artigo(id: string, moduloId: string | null): Artigo {
  return {
    id,
    slug: `slug-${id}`,
    titulo: `Artigo ${id}`,
    trilha_id: "t1",
    modulo_id: moduloId,
    conteudo: {},
    metadados: {},
    autor_id: "u",
    status: "publicado",
    created_at: 0,
    updated_at: 0,
  };
}

// A API devolve fora de ordem de módulo de propósito.
const artigos = [artigo("a3", "m2"), artigo("a1", "m1"), artigo("a2", "m1")];

describe("roadmap", () => {
  it("ordena módulos por `ordem` e mantém a ordem dos artigos dentro do módulo", () => {
    const grupos = agruparPorModulo(trilha, artigos);
    expect(grupos.map((g) => g.slug)).toEqual(["sintaxe", "interfaces"]);
    expect(artigosEmOrdem(trilha, artigos).map((a) => a.id)).toEqual(["a1", "a2", "a3"]);
  });

  it("artigos sem módulo caem em 'Outros artigos'; módulos vazios somem", () => {
    const grupos = agruparPorModulo(trilha, [artigo("x", null)]);
    expect(grupos).toHaveLength(1);
    expect(grupos[0].titulo).toBe("Outros artigos");
  });

  it("visitante: todos 'futuro' (sem 'atual')", () => {
    const nos = montarRoadmap(trilha, artigos, null).flatMap((m) => m.nos);
    expect(nos.every((n) => n.estado === "futuro")).toBe(true);
  });

  it("logado: lido / atual (1º não lido) / futuro", () => {
    const nos = montarRoadmap(trilha, artigos, new Set(["a1"])).flatMap((m) => m.nos);
    expect(nos.map((n) => [n.id, n.estado])).toEqual([
      ["a1", "lido"],
      ["a2", "atual"],
      ["a3", "futuro"],
    ]);
  });

  it("tudo lido: nenhum 'atual'", () => {
    const nos = montarRoadmap(trilha, artigos, new Set(["a1", "a2", "a3"])).flatMap((m) => m.nos);
    expect(nos.every((n) => n.estado === "lido")).toBe(true);
  });

  it("índice reinicia em cada módulo", () => {
    const modulos = montarRoadmap(trilha, artigos, new Set());
    expect(modulos[0].nos.map((n) => n.indice)).toEqual([1, 2]);
    expect(modulos[1].nos.map((n) => n.indice)).toEqual([1]);
  });

  describe("vizinhos (próximo/anterior)", () => {
    const ordem = artigosEmOrdem(trilha, artigos);

    it("meio da trilha tem os dois", () => {
      const v = vizinhos(ordem, "slug-a2");
      expect(v.anterior?.slug).toBe("slug-a1");
      expect(v.proximo?.slug).toBe("slug-a3");
    });
    it("primeiro não tem anterior; último não tem próximo", () => {
      expect(vizinhos(ordem, "slug-a1").anterior).toBeNull();
      expect(vizinhos(ordem, "slug-a3").proximo).toBeNull();
    });
    it("slug fora da trilha → sem vizinhos", () => {
      expect(vizinhos(ordem, "nao-existe")).toEqual({ anterior: null, proximo: null });
    });
  });

  describe("proximoNaoLido", () => {
    const ordem = artigosEmOrdem(trilha, artigos);

    it("pula os já lidos depois do artigo atual", () => {
      expect(proximoNaoLido(ordem, new Set(["a2"]), "slug-a1")?.slug).toBe("slug-a3");
    });
    it("dá a volta no fim da trilha", () => {
      expect(proximoNaoLido(ordem, new Set(["a3"]), "slug-a2")?.slug).toBe("slug-a1");
    });
    it("tudo lido → null", () => {
      expect(proximoNaoLido(ordem, new Set(["a1", "a2", "a3"]), "slug-a1")).toBeNull();
    });
  });
});
