import { describe, expect, it } from "vitest";

import { GAMIFICACAO_VAZIA, mereceRecompensa, mesclarResultado, nivelPorXP, progressoDoNivel, xpInicioNivel } from "@/lib/roadmaps/xp";
import type { Conquista, ResultadoGamificacao } from "@/lib/roadmaps/types";

const conquista = (codigo: string, conquistada = false): Conquista => ({
  codigo,
  nome: codigo,
  descricao: "",
  icone: "star",
  conquistada,
});

const resultado = (extra: Partial<ResultadoGamificacao> = {}): ResultadoGamificacao => ({
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
  ...extra,
});

describe("curva de nível (espelha o servidor)", () => {
  it.each([
    [0, 1],
    [49, 1],
    [50, 2],
    [199, 2],
    [200, 3],
    [450, 4],
    [4050, 10],
  ])("nivelPorXP(%i) = %i", (xp, nivel) => {
    expect(nivelPorXP(xp)).toBe(nivel);
  });

  it("xp negativo conta como zero", () => {
    expect(nivelPorXP(-10)).toBe(1);
  });

  it("início de nível é consistente com nivelPorXP", () => {
    for (let n = 1; n <= 25; n++) {
      expect(nivelPorXP(xpInicioNivel(n))).toBe(n);
      if (n > 1) expect(nivelPorXP(xpInicioNivel(n) - 1)).toBe(n - 1);
    }
  });

  it("progressoDoNivel devolve ganho, tamanho do nível e percentual", () => {
    expect(progressoDoNivel(120)).toEqual({ noNivel: 70, paraProximo: 150, percentual: (70 / 150) * 100 });
    expect(progressoDoNivel(0).percentual).toBe(0);
  });
});

describe("mesclarResultado", () => {
  it("atualiza XP, nível, barra e streak", () => {
    const g = mesclarResultado(GAMIFICACAO_VAZIA, resultado({ xp_ganho: 120, xp_total: 120, nivel: 2, nivel_anterior: 1, subiu_de_nivel: true, streak: 3 }));
    expect(g).toMatchObject({ xp_total: 120, nivel: 2, xp_no_nivel: 70, xp_para_proximo: 150, streak_atual: 3, streak_max: 3, ativo_hoje: true });
  });

  it("não reduz o recorde de sequência", () => {
    const g = mesclarResultado({ ...GAMIFICACAO_VAZIA, streak_max: 9 }, resultado({ streak: 1 }));
    expect(g.streak_max).toBe(9);
    expect(g.streak_atual).toBe(1);
  });

  it("marca como conquistadas as medalhas do catálogo já carregado", () => {
    const base = { ...GAMIFICACAO_VAZIA, conquistas: [conquista("a"), conquista("b")] };
    const g = mesclarResultado(base, resultado({ conquistas_novas: [{ ...conquista("b"), conquistada: true, em: 99 }] }));
    expect(g.conquistas.map((c) => c.conquistada)).toEqual([false, true]);
    expect(g.conquistas[1].em).toBe(99);
  });

  it("sem catálogo carregado, preserva as novas", () => {
    const g = mesclarResultado(GAMIFICACAO_VAZIA, resultado({ conquistas_novas: [conquista("x")] }));
    expect(g.conquistas).toHaveLength(1);
    expect(g.conquistas[0].conquistada).toBe(true);
  });

  it("não muta o resumo original", () => {
    const base = { ...GAMIFICACAO_VAZIA, conquistas: [conquista("a")] };
    mesclarResultado(base, resultado({ xp_total: 10, conquistas_novas: [conquista("a")] }));
    expect(base.conquistas[0].conquistada).toBe(false);
    expect(base.xp_total).toBe(0);
  });
});

describe("mereceRecompensa", () => {
  it("ação idempotente não celebra", () => {
    expect(mereceRecompensa(resultado())).toBe(false);
  });
  it.each([
    resultado({ xp_ganho: 10 }),
    resultado({ subiu_de_nivel: true }),
    resultado({ conquistas_novas: [conquista("a")] }),
    resultado({ roadmaps_completos: ["go"] }),
  ])("celebra %#", (r) => {
    expect(mereceRecompensa(r)).toBe(true);
  });
});
