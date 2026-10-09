import { api } from "@/lib/api";
import { aplicarResultado } from "@/lib/roadmaps/gamificacao-store";
import type { ResultadoGamificacao } from "@/lib/roadmaps/types";
import type { ProgressoTrilha } from "@/lib/types";

type ProgressoArtigoResposta = {
  artigo_id: string;
  concluido: boolean;
  gamificacao?: ResultadoGamificacao | null;
};

/**
 * Cache em memória do progresso por trilha, compartilhado entre card, mapa e artigo.
 * Navegar do artigo para o mapa mostra o nó verde na hora, sem esperar a rede.
 * Fica fora do React (useSyncExternalStore) para sobreviver a trocas de página.
 */

export type ProgressoEntry = {
  status: "loading" | "ready" | "error";
  data: ProgressoTrilha | null;
};

const entries = new Map<string, ProgressoEntry>();
const listeners = new Set<() => void>();

function set(trilhaId: string, entry: ProgressoEntry) {
  entries.set(trilhaId, entry);
  listeners.forEach((fn) => fn());
}

export function subscribeProgresso(listener: () => void): () => void {
  listeners.add(listener);
  return () => listeners.delete(listener);
}

export function getProgressoEntry(trilhaId: string): ProgressoEntry | undefined {
  return entries.get(trilhaId);
}

/** Logout / sessão expirada: nunca mostrar progresso do usuário anterior. */
export function limparProgresso() {
  entries.clear();
  listeners.forEach((fn) => fn());
}

export async function carregarProgresso(trilhaId: string, opts: { force?: boolean } = {}): Promise<void> {
  const atual = entries.get(trilhaId);
  if (!opts.force && atual && atual.status !== "error") return;

  set(trilhaId, { status: "loading", data: atual?.data ?? null });
  try {
    const data = await api<ProgressoTrilha>(`/progresso/trilhas/${trilhaId}`);
    set(trilhaId, { status: "ready", data });
  } catch {
    set(trilhaId, { status: "error", data: atual?.data ?? null });
  }
}

/** Recalcula contadores após marcar/desmarcar um artigo (puro: facilita testar). */
export function aplicarMarcacao(progresso: ProgressoTrilha, artigoId: string, concluido: boolean): ProgressoTrilha {
  const ids = new Set(progresso.artigos_concluidos);
  if (concluido) ids.add(artigoId);
  else ids.delete(artigoId);
  const concluidos = ids.size;
  const total = progresso.total;
  return {
    ...progresso,
    artigos_concluidos: Array.from(ids),
    concluidos,
    percentual: total > 0 ? Math.min(100, (concluidos / total) * 100) : 0,
  };
}

/**
 * Optimistic UI: atualiza o cache antes da API responder; se o PUT falhar, desfaz e propaga o erro
 * para quem chamou mostrar o toast.
 */
export async function marcarArtigo(params: {
  artigoId: string;
  trilhaId: string | null;
  concluido: boolean;
}): Promise<void> {
  const { artigoId, trilhaId, concluido } = params;
  const antes = trilhaId ? entries.get(trilhaId) : undefined;

  if (trilhaId && antes?.data) {
    set(trilhaId, { status: "ready", data: aplicarMarcacao(antes.data, artigoId, concluido) });
  }

  try {
    const resp = await api<ProgressoArtigoResposta | undefined>(`/progresso/artigos/${artigoId}`, {
      method: "PUT",
      body: { concluido },
    });
    // Ler um artigo ligado a um nó de roadmap pode render XP/conquistas: o HUD e os toasts reagem.
    aplicarResultado(resp?.gamificacao);
  } catch (err) {
    if (trilhaId && antes) set(trilhaId, antes);
    throw err;
  }
}

/** Registra a visita ("em curso") sem concluir. Silencioso: é só para o card "Continuar". */
export async function registrarVisita(artigoId: string): Promise<void> {
  try {
    await api(`/progresso/artigos/${artigoId}`, { method: "PUT", body: { concluido: false } });
  } catch {
    /* melhor esforço */
  }
}
