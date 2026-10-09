import { describe, expect, it } from "vitest";

import {
  agruparPorEtapa,
  arestaPercorrida,
  autoLayout,
  descricaoEstado,
  emOrdemDeEstudo,
  prerequisitosPendentes,
  profundidades,
  proximoRecomendado,
  semPosicoes,
} from "@/lib/roadmaps/graph";
import type { EstadoNo, RoadmapAresta, RoadmapNo } from "@/lib/roadmaps/types";

function no(id: string, extra: Partial<RoadmapNo> = {}): RoadmapNo {
  return { id, tipo: "topico", titulo: `Nó ${id}`, pos_x: 0, pos_y: 0, xp: 10, conclusao: "manual", ...extra };
}
const req = (origem: string, destino: string): RoadmapAresta => ({ origem, destino, tipo: "requer" });
const opc = (origem: string, destino: string): RoadmapAresta => ({ origem, destino, tipo: "opcional" });

// a → b → d ; a → c → d ; e solto
const nos = [no("a"), no("b"), no("c"), no("d", { tipo: "chefe", conclusao: "quiz" }), no("e")];
const arestas = [req("a", "b"), req("a", "c"), req("b", "d"), req("c", "d")];

describe("profundidades", () => {
  it("usa o maior caminho a partir da raiz", () => {
    const p = profundidades(nos, arestas);
    expect([p.get("a"), p.get("b"), p.get("c"), p.get("d"), p.get("e")]).toEqual([0, 1, 1, 2, 0]);
  });

  it("arestas opcionais também posicionam o nó depois da origem (só não bloqueiam)", () => {
    const p = profundidades([no("x"), no("y")], [opc("x", "y")]);
    expect(p.get("y")).toBe(1);
  });

  it("ignora arestas para nós inexistentes e não trava em ciclo", () => {
    expect(profundidades([no("x")], [req("x", "fantasma")]).get("x")).toBe(0);
    const ciclo = profundidades([no("x"), no("y")], [req("x", "y"), req("y", "x")]);
    expect(ciclo.size).toBe(2);
  });
});

describe("autoLayout", () => {
  it("coloca cada profundidade em uma linha e centraliza as camadas menores", () => {
    const pos = autoLayout(nos, arestas);
    expect(pos.get("a")!.y).toBe(0);
    expect(pos.get("b")!.y).toBe(pos.get("c")!.y);
    expect(pos.get("d")!.y).toBeGreaterThan(pos.get("b")!.y);
    expect(pos.get("b")!.x).not.toBe(pos.get("c")!.x);
  });

  it("nunca sobrepõe dois nós", () => {
    const pos = autoLayout(nos, arestas);
    const chaves = new Set(Array.from(pos.values()).map((p) => `${p.x},${p.y}`));
    expect(chaves.size).toBe(nos.length);
  });
});

describe("semPosicoes", () => {
  it("detecta nós ainda não posicionados", () => {
    expect(semPosicoes(nos)).toBe(true);
    expect(semPosicoes([no("a", { pos_x: 10 }), no("b")])).toBe(false);
    expect(semPosicoes([no("a")])).toBe(false);
  });
});

describe("ordem e etapas", () => {
  it("ordena por profundidade", () => {
    const ordem = emOrdemDeEstudo(nos, arestas).map((n) => n.id);
    expect(ordem.indexOf("a")).toBeLessThan(ordem.indexOf("b"));
    expect(ordem.indexOf("b")).toBeLessThan(ordem.indexOf("d"));
  });

  it("agrupa em etapas numeradas a partir de 1", () => {
    const g = agruparPorEtapa(nos, arestas);
    expect(g.map((x) => x.etapa)).toEqual([1, 2, 3]);
    expect(g[0].nos.map((n) => n.id).sort()).toEqual(["a", "e"]);
    expect(g[2].nos.map((n) => n.id)).toEqual(["d"]);
  });
});

describe("proximoRecomendado", () => {
  const est = (o: Record<string, EstadoNo>) => o;

  it("visitante (sem estados) não tem recomendação", () => {
    expect(proximoRecomendado(nos, arestas, null)).toBeNull();
  });

  it("prefere nó em curso ao primeiro disponível", () => {
    const e = est({ a: "dominado", b: "disponivel", c: "em_curso", d: "bloqueado", e: "disponivel" });
    expect(proximoRecomendado(nos, arestas, e)?.id).toBe("c");
  });

  it("sem nada em curso, devolve o primeiro disponível na ordem de estudo", () => {
    const e = est({ a: "dominado", b: "disponivel", c: "disponivel", d: "bloqueado", e: "dominado" });
    expect(["b", "c"]).toContain(proximoRecomendado(nos, arestas, e)?.id);
  });

  it("tudo dominado → null", () => {
    const e = est({ a: "dominado", b: "dominado", c: "dominado", d: "dominado", e: "dominado" });
    expect(proximoRecomendado(nos, arestas, e)).toBeNull();
  });
});

describe("prerequisitosPendentes / arestaPercorrida", () => {
  const e: Record<string, EstadoNo> = { a: "dominado", b: "disponivel", c: "dominado", d: "bloqueado", e: "disponivel" };

  it("lista só os pré-requisitos requer ainda não dominados", () => {
    expect(prerequisitosPendentes("d", nos, arestas, e).map((n) => n.id)).toEqual(["b"]);
    expect(prerequisitosPendentes("a", nos, arestas, e)).toEqual([]);
    expect(prerequisitosPendentes("d", nos, arestas, null)).toEqual([]);
  });

  it("aresta é percorrida quando a origem está dominada", () => {
    expect(arestaPercorrida(req("a", "b"), e)).toBe(true);
    expect(arestaPercorrida(req("b", "d"), e)).toBe(false);
    expect(arestaPercorrida(req("a", "b"), null)).toBe(false);
  });
});

describe("descricaoEstado", () => {
  it("diferencia chefe de tópico", () => {
    expect(descricaoEstado("dominado", "chefe")).toMatch(/derrotado/i);
    expect(descricaoEstado("dominado", "topico")).toBe("Dominado");
    expect(descricaoEstado("bloqueado", "topico")).toMatch(/pré-requisitos/);
    expect(descricaoEstado("disponivel", "chefe")).toMatch(/chefe/i);
  });
});
