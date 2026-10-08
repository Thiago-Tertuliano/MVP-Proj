"use client";

import { TrilhaCard } from "@/components/relp/TrilhaCard";
import { useProgressoTrilha } from "@/lib/use-progresso";
import type { Trilha } from "@/lib/types";

/** Liga o card puro ao progresso do aluno (só logado; visitante vê o card sem barra). */
export function TrilhaCardConnected({ trilha }: { trilha: Trilha }) {
  const entry = useProgressoTrilha(trilha.id);

  const progresso = entry
    ? entry.data
      ? {
          status: "ready" as const,
          concluidos: entry.data.concluidos,
          total: entry.data.total,
          percentual: entry.data.percentual,
        }
      : { status: entry.status }
    : undefined;

  return (
    <TrilhaCard
      slug={trilha.slug}
      titulo={trilha.titulo}
      descricao={trilha.descricao}
      totalModulos={trilha.modulos.length}
      progresso={progresso?.status === "error" ? undefined : progresso}
    />
  );
}
