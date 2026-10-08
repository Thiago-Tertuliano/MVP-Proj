import { api, ApiError } from "@/lib/api";
import type { Artigo, Busca, ListaArtigos, ListaTrilhas, ResultadoBusca, Trilha } from "@/lib/types";

/**
 * Leitura de conteúdo público (trilhas, artigos, busca). Roda no servidor (RSC) sem cookie:
 * a API só devolve trilhas/artigos PUBLICADOS nesses endpoints, então rascunho nunca chega ao aluno.
 *
 * Convenção: 404 vira `null` (a página chama `notFound()`); qualquer outra falha propaga `ApiError`
 * para o `error.tsx` mostrar o estado de erro — nunca uma tela vazia que parece "sem conteúdo".
 */

const FRESH = { cache: "no-store" as const };

async function ou404<T>(promise: Promise<T>): Promise<T | null> {
  try {
    return await promise;
  } catch (err) {
    if (err instanceof ApiError && err.isNotFound) return null;
    throw err;
  }
}

export async function listarTrilhas(): Promise<Trilha[]> {
  const resp = await api<ListaTrilhas>("/trilhas", { query: { limit: 100 }, ...FRESH });
  return resp.itens ?? [];
}

export function obterTrilha(slug: string): Promise<Trilha | null> {
  return ou404(api<Trilha>(`/trilhas/${encodeURIComponent(slug)}`, FRESH));
}

export async function listarArtigosDaTrilha(slug: string): Promise<Artigo[]> {
  const resp = await api<ListaArtigos>(`/trilhas/${encodeURIComponent(slug)}/artigos`, FRESH);
  return resp.itens ?? [];
}

/** 404 também quando o artigo existe mas não está publicado (a API não distingue de propósito). */
export function obterArtigo(slug: string): Promise<Artigo | null> {
  return ou404(api<Artigo>(`/artigos/${encodeURIComponent(slug)}`, FRESH));
}

export const BUSCA_MIN_CHARS = 2;

export async function buscarArtigos(q: string): Promise<ResultadoBusca[]> {
  const termo = q.trim();
  if (termo.length < BUSCA_MIN_CHARS) return [];
  const resp = await api<Busca>("/busca", { query: { q: termo }, ...FRESH });
  return resp.itens ?? [];
}
