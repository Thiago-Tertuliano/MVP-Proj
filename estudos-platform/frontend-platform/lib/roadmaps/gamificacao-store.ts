import { obterGamificacao } from "@/lib/roadmaps/client";
import { invalidarProgressoRoadmaps } from "@/lib/roadmaps/progresso-store";
import type { Gamificacao, ResultadoGamificacao } from "@/lib/roadmaps/types";
import { GAMIFICACAO_VAZIA, mereceRecompensa, mesclarResultado } from "@/lib/roadmaps/xp";

/**
 * Cache em memória do resumo de gamificação (XP, nível, streak, conquistas), compartilhado entre
 * o HUD do header, o mapa e a Jornada. Fica fora do React (useSyncExternalStore) para sobreviver
 * à troca de página — mesmo padrão do `progress-store`.
 *
 * Toda ação que pode render XP (concluir nó, quiz, ler artigo) chama `aplicarResultado`: o HUD
 * muda na hora (optimistic) e uma releitura em segundo plano corrige qualquer divergência.
 */

export type GamificacaoEntry = {
  status: "idle" | "loading" | "ready" | "error";
  data: Gamificacao | null;
};

let entry: GamificacaoEntry = { status: "idle", data: null };
let emAndamento: Promise<void> | null = null;
let geracao = 0; // invalida respostas de requisições antigas após logout

const listeners = new Set<() => void>();
const recompensaListeners = new Set<(r: ResultadoGamificacao) => void>();

function set(next: GamificacaoEntry) {
  entry = next;
  listeners.forEach((fn) => fn());
}

export function subscribeGamificacao(listener: () => void): () => void {
  listeners.add(listener);
  return () => listeners.delete(listener);
}

export function getGamificacaoEntry(): GamificacaoEntry {
  return entry;
}

/** Logout / sessão expirada: nunca mostrar XP do usuário anterior. */
export function limparGamificacao() {
  geracao++;
  emAndamento = null;
  set({ status: "idle", data: null });
}

export async function carregarGamificacao(opts: { force?: boolean } = {}): Promise<void> {
  if (!opts.force && (entry.status === "ready" || entry.status === "loading")) return;
  if (emAndamento) return emAndamento;

  const minhaGeracao = geracao;
  set({ status: "loading", data: entry.data });
  emAndamento = obterGamificacao()
    .then((data) => {
      if (minhaGeracao !== geracao) return;
      set({ status: "ready", data });
    })
    .catch(() => {
      if (minhaGeracao !== geracao) return;
      set({ status: "error", data: entry.data });
    })
    .finally(() => {
      if (minhaGeracao === geracao) emAndamento = null;
    });
  return emAndamento;
}

/** Quem exibe celebrações (toast/level-up) assina aqui; a fonte da ação só reporta o resultado. */
export function onRecompensa(listener: (r: ResultadoGamificacao) => void): () => void {
  recompensaListeners.add(listener);
  return () => recompensaListeners.delete(listener);
}

/** Aplica o resultado de uma ação ao cache, anuncia a recompensa e confirma com o servidor. */
export function aplicarResultado(r: ResultadoGamificacao | null | undefined) {
  if (!r) return;
  // Nós concluídos fora da tela do roadmap (ex.: leitura de artigo) deixam o mapa em cache defasado.
  if (r.nos_concluidos.length > 0) invalidarProgressoRoadmaps();
  set({ status: "ready", data: mesclarResultado(entry.data ?? GAMIFICACAO_VAZIA, r) });
  if (mereceRecompensa(r)) recompensaListeners.forEach((fn) => fn(r));
  void carregarGamificacao({ force: true });
}
