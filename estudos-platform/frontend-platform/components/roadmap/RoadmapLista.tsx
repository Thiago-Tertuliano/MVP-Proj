import * as React from "react";

import { RoadmapNodeCard } from "@/components/roadmap/RoadmapNodeCard";
import { agruparPorEtapa } from "@/lib/roadmaps/graph";
import type { EstadoNo, RoadmapAresta, RoadmapNo } from "@/lib/roadmaps/types";
import { cn } from "@/lib/utils";

export type RoadmapListaProps = {
  nos: RoadmapNo[];
  arestas: RoadmapAresta[];
  estados: Record<string, EstadoNo> | null;
  selecionadoId?: string | null;
  onSelecionar: (id: string) => void;
  className?: string;
};

/**
 * Visão em lista por etapas: alternativa ao canvas para teclado, leitor de tela e telas pequenas.
 * Cada etapa é um grupo de nós que podem ser estudados em paralelo (mesma profundidade no grafo).
 */
export function RoadmapLista({ nos, arestas, estados, selecionadoId, onSelecionar, className }: RoadmapListaProps) {
  const etapas = React.useMemo(() => agruparPorEtapa(nos, arestas), [nos, arestas]);

  return (
    <ol className={cn("space-y-6", className)} aria-label="Etapas do roadmap">
      {etapas.map(({ etapa, nos: grupo }) => {
        const dominados = estados ? grupo.filter((n) => estados[n.id] === "dominado").length : 0;
        const titulo = `etapa-${etapa}`;
        return (
          <li key={etapa} aria-labelledby={titulo} className="space-y-3">
            <div className="flex flex-wrap items-baseline gap-x-3">
              <h3 id={titulo} className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">
                Etapa {etapa}
              </h3>
              {estados && (
                <span className="text-xs tabular-nums text-muted-foreground">
                  {dominados}/{grupo.length} dominados
                </span>
              )}
            </div>
            <ul className="grid gap-3 sm:grid-cols-2">
              {grupo.map((no) => (
                <li key={no.id}>
                  <RoadmapNodeCard
                    titulo={no.titulo}
                    tipo={no.tipo}
                    estado={estados ? (estados[no.id] ?? "bloqueado") : null}
                    xp={no.xp}
                    fixo={false}
                    selecionado={selecionadoId === no.id}
                    onSelecionar={() => onSelecionar(no.id)}
                  />
                </li>
              ))}
            </ul>
          </li>
        );
      })}
    </ol>
  );
}
