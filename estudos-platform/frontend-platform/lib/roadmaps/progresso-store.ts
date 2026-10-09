import { obterProgressoRoadmap } from "@/lib/roadmaps/client";
import type { ConclusaoNoResposta, RoadmapProgresso } from "@/lib/roadmaps/types";

/**
 * Cache em memória do progresso por roadmap (estados dos nós). Mesmo padrão do `progress-store`:
 * fora do React, assinado via useSyncExternalStore, limpo no logout.
 */

export type RoadmapProgressoEntry = {
  status: "loading" | "ready" | "error";
  data: RoadmapProgresso | null;
  /** Algo mudou fora desta tela (ex.: leitura de artigo ligado a um nó): refaz o GET no próximo uso. */
  obsoleto?: boolean;
};

const entries = new Map<string, RoadmapProgressoEntry>();
const listeners = new Set<() => void>();

function set(slug: string, entry: RoadmapProgressoEntry) {
  entries.set(slug, entry);
  listeners.forEach((fn) => fn());
}

export function subscribeProgressoRoadmap(listener: () => void): () => void {
  listeners.add(listener);
  return () => listeners.delete(listener);
}

export function getProgressoRoadmapEntry(slug: string): RoadmapProgressoEntry | undefined {
  return entries.get(slug);
}

export function limparProgressoRoadmaps() {
  entries.clear();
  listeners.forEach((fn) => fn());
}

/** Marca tudo como desatualizado sem notificar: a próxima montagem refaz o GET, sem piscar a tela atual. */
export function invalidarProgressoRoadmaps() {
  entries.forEach((entry, slug) => entries.set(slug, { ...entry, obsoleto: true }));
}

export async function carregarProgressoRoadmap(slug: string, opts: { force?: boolean } = {}): Promise<void> {
  const atual = entries.get(slug);
  if (!opts.force && atual && atual.status !== "error" && !atual.obsoleto) return;

  set(slug, { status: "loading", data: atual?.data ?? null });
  try {
    const data = await obterProgressoRoadmap(slug);
    set(slug, { status: "ready", data });
  } catch {
    set(slug, { status: "error", data: atual?.data ?? null });
  }
}

/** Funde a resposta de "concluir nó"/"quiz" no cache (puro: facilita testar). */
export function aplicarConclusaoNoProgresso(atual: RoadmapProgresso, r: ConclusaoNoResposta): RoadmapProgresso {
  const total = r.total || atual.total;
  return {
    ...atual,
    estados: r.estados,
    concluidos: r.concluidos,
    total,
    percentual: total > 0 ? Math.min(100, (r.concluidos / total) * 100) : 0,
    xp_ganho: atual.xp_ganho + r.resultado.xp_ganho,
    completo: total > 0 && r.concluidos >= total,
  };
}

export function aplicarConclusao(slug: string, r: ConclusaoNoResposta) {
  const atual = entries.get(slug);
  if (!atual?.data) {
    // Sem base para fundir: pede o estado completo ao servidor.
    void carregarProgressoRoadmap(slug, { force: true });
    return;
  }
  set(slug, { status: "ready", data: aplicarConclusaoNoProgresso(atual.data, r) });
}
