import { api, ApiError } from "@/lib/api";
import type {
  ConclusaoNoResposta,
  Gamificacao,
  ListaProgressoRoadmaps,
  ListaRoadmaps,
  QuizResultado,
  Roadmap,
  RoadmapProgresso,
  RoadmapResumo,
} from "@/lib/roadmaps/types";

/**
 * Acesso à API de roadmaps. As leituras públicas rodam no servidor (RSC) sem cookie e só enxergam
 * roadmaps PUBLICADOS; progresso e gamificação exigem sessão e rodam no browser.
 *
 * Convenção do projeto: 404 vira `null` (a página chama `notFound()`); outras falhas propagam `ApiError`.
 */

const FRESH = { cache: "no-store" as const };

export async function listarRoadmaps(): Promise<RoadmapResumo[]> {
  const resp = await api<ListaRoadmaps>("/roadmaps", FRESH);
  return resp.itens ?? [];
}

export async function obterRoadmap(slug: string): Promise<Roadmap | null> {
  try {
    return await api<Roadmap>(`/roadmaps/${encodeURIComponent(slug)}`, FRESH);
  } catch (err) {
    if (err instanceof ApiError && err.isNotFound) return null;
    throw err;
  }
}

/* ------------------------------------------------------------------ aluno (sessão) */

export function obterProgressoRoadmap(slug: string, signal?: AbortSignal): Promise<RoadmapProgresso> {
  return api<RoadmapProgresso>(`/roadmaps/${encodeURIComponent(slug)}/progresso`, { signal, cache: "no-store" });
}

export async function listarProgressoRoadmaps(signal?: AbortSignal) {
  const resp = await api<ListaProgressoRoadmaps>("/gamificacao/roadmaps", { signal, cache: "no-store" });
  return resp.itens ?? [];
}

export function concluirNo(slug: string, noId: string): Promise<ConclusaoNoResposta> {
  return api<ConclusaoNoResposta>(`/roadmaps/${encodeURIComponent(slug)}/nos/${encodeURIComponent(noId)}/concluir`, {
    method: "POST",
  });
}

export function responderQuiz(slug: string, noId: string, respostas: Record<string, string>): Promise<QuizResultado> {
  return api<QuizResultado>(`/roadmaps/${encodeURIComponent(slug)}/nos/${encodeURIComponent(noId)}/quiz`, {
    method: "POST",
    body: { respostas },
  });
}

export function obterGamificacao(signal?: AbortSignal): Promise<Gamificacao> {
  return api<Gamificacao>("/gamificacao/me", { signal, cache: "no-store" });
}
