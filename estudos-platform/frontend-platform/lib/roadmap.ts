import type { Artigo, Trilha } from "@/lib/types";

/**
 * Modelo de tela do mapa da trilha, montado a partir de:
 *   GET /trilhas/{slug}           → módulos (ordem)
 *   GET /trilhas/{slug}/artigos   → nós
 *   GET /progresso/trilhas/{id}   → artigos concluídos (só logado)
 */

export type NoEstado = "lido" | "atual" | "futuro";

export type RoadmapNoData = {
  id: string;
  slug: string;
  titulo: string;
  /** Posição dentro do módulo (1-based). */
  indice: number;
  estado: NoEstado;
};

export type RoadmapModuloData = {
  slug: string;
  titulo: string;
  descricao?: string;
  nos: RoadmapNoData[];
};

type ModuloComArtigos = { slug: string; titulo: string; descricao?: string; artigos: Artigo[] };

const MODULO_AVULSO = { slug: "outros", titulo: "Outros artigos" };

/** Módulos por `ordem`; artigos mantêm a ordem da API. Artigos sem módulo caem em "Outros artigos". */
export function agruparPorModulo(trilha: Trilha, artigos: Artigo[]): ModuloComArtigos[] {
  const modulos = [...trilha.modulos].sort((a, b) => a.ordem - b.ordem);
  const porModulo = new Map<string, Artigo[]>(modulos.map((m) => [m.id, []]));
  const avulsos: Artigo[] = [];

  for (const artigo of artigos) {
    const lista = artigo.modulo_id ? porModulo.get(artigo.modulo_id) : undefined;
    if (lista) lista.push(artigo);
    else avulsos.push(artigo);
  }

  const grupos: ModuloComArtigos[] = modulos
    .map((m) => ({ slug: m.slug, titulo: m.titulo, descricao: m.descricao, artigos: porModulo.get(m.id) ?? [] }))
    .filter((g) => g.artigos.length > 0);

  if (avulsos.length > 0) grupos.push({ ...MODULO_AVULSO, artigos: avulsos });
  return grupos;
}

/** Todos os artigos na ordem em que o aluno deve ler (módulos em ordem, depois artigos). */
export function artigosEmOrdem(trilha: Trilha, artigos: Artigo[]): Artigo[] {
  return agruparPorModulo(trilha, artigos).flatMap((g) => g.artigos);
}

/**
 * Estados dos nós. `concluidos = null` = visitante (sem progresso): todos "futuro" e sem "atual".
 * Logado: lido = concluído; atual = primeiro ainda não lido da trilha; futuro = o restante.
 */
export function montarRoadmap(
  trilha: Trilha,
  artigos: Artigo[],
  concluidos: ReadonlySet<string> | null,
): RoadmapModuloData[] {
  const grupos = agruparPorModulo(trilha, artigos);
  let atualMarcado = concluidos === null;

  return grupos.map((g) => ({
    slug: g.slug,
    titulo: g.titulo,
    descricao: g.descricao,
    nos: g.artigos.map((a, i) => {
      let estado: NoEstado = "futuro";
      if (concluidos?.has(a.id)) {
        estado = "lido";
      } else if (!atualMarcado) {
        estado = "atual";
        atualMarcado = true;
      }
      return { id: a.id, slug: a.slug, titulo: a.titulo, indice: i + 1, estado };
    }),
  }));
}

export type ArtigoVizinho = { slug: string; titulo: string };

/** Anterior/próximo dentro da MESMA trilha: nulos nas pontas (o botão fica desabilitado). */
export function vizinhos(ordenados: { slug: string; titulo: string }[], slug: string): {
  anterior: ArtigoVizinho | null;
  proximo: ArtigoVizinho | null;
} {
  const i = ordenados.findIndex((a) => a.slug === slug);
  if (i === -1) return { anterior: null, proximo: null };
  const pick = (a?: { slug: string; titulo: string }) => (a ? { slug: a.slug, titulo: a.titulo } : null);
  return { anterior: pick(ordenados[i - 1]), proximo: pick(ordenados[i + 1]) };
}

/** Próximo artigo ainda não lido depois de `depoisDe` (dá a volta no fim da trilha). */
export function proximoNaoLido(
  ordenados: { id: string; slug: string; titulo: string }[],
  concluidos: ReadonlySet<string>,
  depoisDe?: string,
): { slug: string; titulo: string } | null {
  const inicio = depoisDe ? ordenados.findIndex((a) => a.slug === depoisDe) + 1 : 0;
  const rotacionados = [...ordenados.slice(inicio), ...ordenados.slice(0, inicio)];
  const alvo = rotacionados.find((a) => !concluidos.has(a.id));
  return alvo ? { slug: alvo.slug, titulo: alvo.titulo } : null;
}
