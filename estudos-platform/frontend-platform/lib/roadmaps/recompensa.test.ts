import { describe, expect, it } from "vitest";

import { celebracaoDe } from "@/lib/roadmaps/recompensa";
import type { ResultadoGamificacao } from "@/lib/roadmaps/types";

const base: ResultadoGamificacao = {
  xp_ganho: 0,
  xp_total: 0,
  nivel: 1,
  nivel_anterior: 1,
  subiu_de_nivel: false,
  streak: 0,
  nos_concluidos: [],
  nos_desbloqueados: [],
  roadmaps_completos: [],
  conquistas_novas: [],
};

describe("celebracaoDe", () => {
  it("XP simples vira só o título", () => {
    expect(celebracaoDe({ ...base, xp_ganho: 30, nos_concluidos: ["a"] })).toMatchObject({
      titulo: "+30 XP",
      descricao: undefined,
      subiuDeNivel: null,
    });
  });

  it("descreve desbloqueios, vários nós e sequência", () => {
    const c = celebracaoDe({
      ...base,
      xp_ganho: 80,
      nos_concluidos: ["a", "b"],
      nos_desbloqueados: ["c"],
      streak: 4,
    });
    expect(c.descricao).toBe("2 nós dominados · 1 nó desbloqueado · sequência de 4 dias");
  });

  it("pluraliza desbloqueios", () => {
    expect(celebracaoDe({ ...base, nos_desbloqueados: ["a", "b"] }).descricao).toBe("2 nós desbloqueados");
  });

  it("sem XP não há título, mas level-up e conquistas continuam", () => {
    const c = celebracaoDe({
      ...base,
      subiu_de_nivel: true,
      nivel: 3,
      conquistas_novas: [{ codigo: "x", nome: "X", descricao: "d", icone: "star", conquistada: true }],
      roadmaps_completos: ["go"],
    });
    expect(c.titulo).toBeNull();
    expect(c.subiuDeNivel).toBe(3);
    expect(c.conquistas).toEqual([{ codigo: "x", nome: "X", descricao: "d" }]);
    expect(c.roadmapsCompletos).toEqual(["go"]);
  });
});
