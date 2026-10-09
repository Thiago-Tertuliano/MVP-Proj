"use client";

import * as React from "react";

import { RoadmapCard } from "@/components/roadmap/RoadmapCard";
import { listarProgressoRoadmaps } from "@/lib/roadmaps/client";
import type { RoadmapProgressoItem, RoadmapResumo } from "@/lib/roadmaps/types";
import { useSession } from "@/lib/session";

/** Grade de roadmaps; sobrepõe o progresso do aluno logado (uma única chamada para todos os cards). */
export function RoadmapCatalogo({ roadmaps }: { roadmaps: RoadmapResumo[] }) {
  const { status } = useSession();
  const [progresso, setProgresso] = React.useState<Map<string, RoadmapProgressoItem> | "loading" | null>(null);

  React.useEffect(() => {
    if (status !== "user") {
      setProgresso(null);
      return;
    }
    const ctrl = new AbortController();
    setProgresso("loading");
    listarProgressoRoadmaps(ctrl.signal)
      .then((itens) => setProgresso(new Map(itens.map((i) => [i.slug, i]))))
      // Falha de progresso não esconde o catálogo: os cards ficam sem barra.
      .catch(() => !ctrl.signal.aborted && setProgresso(null));
    return () => ctrl.abort();
  }, [status]);

  return (
    <ul className="grid gap-4 md:grid-cols-2">
      {roadmaps.map((r) => {
        const item = progresso instanceof Map ? progresso.get(r.slug) : undefined;
        return (
          <li key={r.id}>
            <RoadmapCard
              slug={r.slug}
              titulo={r.titulo}
              descricao={r.descricao}
              icone={r.icone}
              totalNos={r.total_nos}
              totalChefes={r.total_chefes}
              xpTotal={r.xp_total}
              progresso={
                progresso === "loading"
                  ? { status: "loading" }
                  : progresso instanceof Map
                    ? {
                        status: "ready",
                        concluidos: item?.concluidos ?? 0,
                        total: item?.total ?? r.total_nos,
                        percentual: item?.percentual ?? 0,
                      }
                    : undefined
              }
            />
          </li>
        );
      })}
    </ul>
  );
}
