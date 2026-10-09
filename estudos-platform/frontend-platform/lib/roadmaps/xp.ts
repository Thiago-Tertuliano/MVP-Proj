import type { Gamificacao, ResultadoGamificacao } from "@/lib/roadmaps/types";

/**
 * Espelho da curva do servidor (usado só para exibição imediata; o servidor continua sendo a verdade
 * e a próxima leitura de `/gamificacao/me` corrige qualquer divergência).
 * nível = floor(sqrt(xp / 50)) + 1  →  XP para alcançar o nível n = 50 · (n − 1)².
 */
export const XP_BASE_NIVEL = 50;

export function nivelPorXP(xp: number): number {
  return Math.floor(Math.sqrt(Math.max(0, xp) / XP_BASE_NIVEL)) + 1;
}

export function xpInicioNivel(nivel: number): number {
  return nivel <= 1 ? 0 : XP_BASE_NIVEL * (nivel - 1) ** 2;
}

/** XP ganho dentro do nível atual e quanto o nível inteiro vale. */
export function progressoDoNivel(xp: number): { noNivel: number; paraProximo: number; percentual: number } {
  const nivel = nivelPorXP(xp);
  const inicio = xpInicioNivel(nivel);
  const paraProximo = xpInicioNivel(nivel + 1) - inicio;
  const noNivel = Math.max(0, xp - inicio);
  return { noNivel, paraProximo, percentual: paraProximo > 0 ? Math.min(100, (noNivel / paraProximo) * 100) : 0 };
}

/** Resumo inicial para quem ainda não tem registro de gamificação. */
export const GAMIFICACAO_VAZIA: Gamificacao = {
  xp_total: 0,
  nivel: 1,
  xp_no_nivel: 0,
  xp_para_proximo: XP_BASE_NIVEL,
  streak_atual: 0,
  streak_max: 0,
  ativo_hoje: false,
  conquistas: [],
};

/** Aplica o resultado de uma ação ao resumo em cache (optimistic: o servidor confirma depois). */
export function mesclarResultado(atual: Gamificacao, r: ResultadoGamificacao): Gamificacao {
  const { noNivel, paraProximo } = progressoDoNivel(r.xp_total);
  const novas = new Map(r.conquistas_novas.map((c) => [c.codigo, c]));
  const conquistas = atual.conquistas.map((c) => (novas.has(c.codigo) ? { ...c, ...novas.get(c.codigo)!, conquistada: true } : c));
  // Catálogo ainda não carregado: guarda só as novas para não perder a informação.
  if (atual.conquistas.length === 0) conquistas.push(...r.conquistas_novas.map((c) => ({ ...c, conquistada: true })));

  return {
    ...atual,
    xp_total: r.xp_total,
    nivel: r.nivel,
    xp_no_nivel: noNivel,
    xp_para_proximo: paraProximo,
    streak_atual: r.streak,
    streak_max: Math.max(atual.streak_max, r.streak),
    ativo_hoje: atual.ativo_hoje || r.xp_ganho > 0 || r.streak > 0,
    conquistas,
  };
}

/** Vale a pena celebrar? (evita toast vazio em ações idempotentes). */
export function mereceRecompensa(r: ResultadoGamificacao): boolean {
  return r.xp_ganho > 0 || r.subiu_de_nivel || r.conquistas_novas.length > 0 || r.roadmaps_completos.length > 0;
}
