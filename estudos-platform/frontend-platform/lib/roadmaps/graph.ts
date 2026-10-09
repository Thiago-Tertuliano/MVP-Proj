import type { EstadoNo, RoadmapAresta, RoadmapNo, TipoNo } from "@/lib/roadmaps/types";

/**
 * Regras de apresentação do grafo (puras, sem rede). A fonte de verdade dos estados é o servidor;
 * aqui só ordenamos, posicionamos e rotulamos.
 */

export const ESTADO_ROTULO: Record<EstadoNo, string> = {
  bloqueado: "Bloqueado",
  disponivel: "Disponível",
  em_curso: "Em curso",
  dominado: "Dominado",
};

export const TIPO_ROTULO: Record<TipoNo, string> = {
  topico: "Tópico",
  chefe: "Chefe",
  marco: "Marco",
};

/** Texto lido por leitores de tela para o estado de um nó. */
export function descricaoEstado(estado: EstadoNo, tipo: TipoNo): string {
  if (estado === "bloqueado") return "Bloqueado: conclua os pré-requisitos para liberar";
  if (estado === "dominado") return tipo === "chefe" ? "Chefe derrotado" : "Dominado";
  if (estado === "em_curso") return "Em curso";
  return tipo === "chefe" ? "Chefe disponível para enfrentar" : "Disponível para estudar";
}

/* ------------------------------------------------------------------ profundidade / layout */

/**
 * Profundidade (maior caminho a partir de uma raiz). Arestas `opcional` contam: elas não bloqueiam
 * o aluno, mas dizem onde o nó fica no mapa e em que ordem faz sentido estudá-lo.
 */
export function profundidades(nos: RoadmapNo[], arestas: RoadmapAresta[]): Map<string, number> {
  const ids = new Set(nos.map((n) => n.id));
  const pais = new Map<string, string[]>();
  for (const a of arestas) {
    if (!ids.has(a.origem) || !ids.has(a.destino)) continue;
    pais.set(a.destino, [...(pais.get(a.destino) ?? []), a.origem]);
  }

  const memo = new Map<string, number>();
  const visitando = new Set<string>();
  const calc = (id: string): number => {
    const cached = memo.get(id);
    if (cached !== undefined) return cached;
    if (visitando.has(id)) return 0; // defesa contra ciclo (o servidor já rejeita)
    visitando.add(id);
    const ps = pais.get(id) ?? [];
    const d = ps.length === 0 ? 0 : 1 + Math.max(...ps.map(calc));
    visitando.delete(id);
    memo.set(id, d);
    return d;
  };
  for (const n of nos) calc(n.id);
  return memo;
}

export const LAYOUT = { larguraNo: 220, alturaNo: 84, espacoX: 40, espacoY: 70 } as const;

/** Posiciona em camadas (topo → base). Usado quando o editor ainda não posicionou os nós. */
export function autoLayout(nos: RoadmapNo[], arestas: RoadmapAresta[]): Map<string, { x: number; y: number }> {
  const prof = profundidades(nos, arestas);
  const camadas = new Map<number, RoadmapNo[]>();
  for (const n of nos) {
    const d = prof.get(n.id) ?? 0;
    camadas.set(d, [...(camadas.get(d) ?? []), n]);
  }
  const pos = new Map<string, { x: number; y: number }>();
  const passoX = LAYOUT.larguraNo + LAYOUT.espacoX;
  const passoY = LAYOUT.alturaNo + LAYOUT.espacoY;
  const maiorCamada = Math.max(1, ...Array.from(camadas.values()).map((c) => c.length));
  Array.from(camadas.entries()).forEach(([d, lista]) => {
    const offset = ((maiorCamada - lista.length) * passoX) / 2;
    lista.forEach((n, i) => pos.set(n.id, { x: offset + i * passoX, y: d * passoY }));
  });
  return pos;
}

/** Verdadeiro quando nenhum nó foi posicionado (todos em 0,0) — cai no layout automático. */
export function semPosicoes(nos: RoadmapNo[]): boolean {
  return nos.length > 1 && nos.every((n) => n.pos_x === 0 && n.pos_y === 0);
}

export function posicaoDoNo(n: RoadmapNo, auto: Map<string, { x: number; y: number }> | null) {
  if (auto) return auto.get(n.id) ?? { x: 0, y: 0 };
  return { x: n.pos_x, y: n.pos_y };
}

/* ------------------------------------------------------------------ navegação / recomendação */

/** Ordem de estudo estável: profundidade, depois posição, depois título. */
export function emOrdemDeEstudo(nos: RoadmapNo[], arestas: RoadmapAresta[]): RoadmapNo[] {
  const prof = profundidades(nos, arestas);
  return [...nos].sort(
    (a, b) =>
      (prof.get(a.id) ?? 0) - (prof.get(b.id) ?? 0) ||
      a.pos_x - b.pos_x ||
      a.pos_y - b.pos_y ||
      a.titulo.localeCompare(b.titulo, "pt-BR"),
  );
}

/** Agrupa por profundidade para a visão em lista (alternativa acessível ao canvas). */
export function agruparPorEtapa(nos: RoadmapNo[], arestas: RoadmapAresta[]): { etapa: number; nos: RoadmapNo[] }[] {
  const prof = profundidades(nos, arestas);
  const grupos = new Map<number, RoadmapNo[]>();
  for (const n of emOrdemDeEstudo(nos, arestas)) {
    const d = prof.get(n.id) ?? 0;
    grupos.set(d, [...(grupos.get(d) ?? []), n]);
  }
  return Array.from(grupos.entries())
    .sort((a, b) => a[0] - b[0])
    .map(([etapa, lista]) => ({ etapa: etapa + 1, nos: lista }));
}

/** Próximo nó sugerido: prioriza "em curso", depois o primeiro "disponível" na ordem de estudo. */
export function proximoRecomendado(
  nos: RoadmapNo[],
  arestas: RoadmapAresta[],
  estados: Record<string, EstadoNo> | null,
): RoadmapNo | null {
  if (!estados) return null;
  const ordem = emOrdemDeEstudo(nos, arestas);
  return (
    ordem.find((n) => estados[n.id] === "em_curso") ?? ordem.find((n) => estados[n.id] === "disponivel") ?? null
  );
}

/** Pré-requisitos (`requer`) ainda não dominados — explica por que um nó está bloqueado. */
export function prerequisitosPendentes(
  noId: string,
  nos: RoadmapNo[],
  arestas: RoadmapAresta[],
  estados: Record<string, EstadoNo> | null,
): RoadmapNo[] {
  if (!estados) return [];
  const porId = new Map(nos.map((n) => [n.id, n]));
  return arestas
    .filter((a) => a.tipo === "requer" && a.destino === noId && estados[a.origem] !== "dominado")
    .map((a) => porId.get(a.origem))
    .filter((n): n is RoadmapNo => !!n);
}

/** Aresta "acesa": o caminho já foi percorrido (origem dominada). */
export function arestaPercorrida(a: RoadmapAresta, estados: Record<string, EstadoNo> | null): boolean {
  return !!estados && estados[a.origem] === "dominado";
}
