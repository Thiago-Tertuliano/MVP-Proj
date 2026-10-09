import * as React from "react";
import Link from "next/link";
import { ArrowRight, Swords } from "lucide-react";

import { Icone } from "@/components/roadmap/icones";
import { ProgressSummary } from "@/components/relp/ProgressSummary";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";

export type RoadmapCardProps = {
  slug: string;
  titulo: string;
  descricao?: string;
  icone?: string;
  totalNos: number;
  totalChefes: number;
  xpTotal: number;
  /** Só para logado. `loading` mostra skeleton; ausente = visitante (sem barra). */
  progresso?: { status: "loading" | "ready"; concluidos?: number; total?: number; percentual?: number };
  className?: string;
};

/** Card do catálogo de roadmaps. Todo o card é um link (um único foco de teclado). */
export function RoadmapCard({
  slug,
  titulo,
  descricao,
  icone,
  totalNos,
  totalChefes,
  xpTotal,
  progresso,
  className,
}: RoadmapCardProps) {
  return (
    <Link
      href={`/roadmaps/${slug}`}
      className={cn(
        "group flex h-full flex-col gap-4 rounded-xl border border-border bg-card p-5 text-card-foreground transition-colors hover:border-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background",
        className,
      )}
    >
      <div className="flex flex-1 gap-3">
        <span
          aria-hidden="true"
          className="flex h-11 w-11 shrink-0 items-center justify-center rounded-lg bg-primary-muted text-primary"
        >
          <Icone nome={icone} className="h-6 w-6" />
        </span>
        <div className="min-w-0 flex-1 space-y-1.5">
          <h3 className="text-lg font-semibold leading-tight tracking-tight group-hover:text-primary">{titulo}</h3>
          {descricao && <p className="line-clamp-3 text-sm text-muted-foreground">{descricao}</p>}
        </div>
      </div>

      {progresso?.status === "loading" && <Skeleton className="h-8 w-full" />}
      {progresso?.status === "ready" && progresso.total !== undefined && (
        <ProgressSummary
          concluidos={progresso.concluidos ?? 0}
          total={progresso.total}
          percentual={progresso.percentual}
          label={`${progresso.concluidos ?? 0} de ${progresso.total} nós dominados`}
        />
      )}

      <div className="flex items-center justify-between gap-2 text-sm">
        <span className="inline-flex flex-wrap items-center gap-x-3 gap-y-1 text-muted-foreground">
          <span>
            {totalNos} {totalNos === 1 ? "nó" : "nós"}
          </span>
          {totalChefes > 0 && (
            <span className="inline-flex items-center gap-1">
              <Swords className="h-4 w-4" aria-hidden="true" />
              {totalChefes} {totalChefes === 1 ? "chefe" : "chefes"}
            </span>
          )}
          <span className="font-medium text-notify">{xpTotal} XP</span>
        </span>
        <span className="inline-flex shrink-0 items-center gap-1 font-medium text-primary">
          Abrir
          <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-0.5" aria-hidden="true" />
        </span>
      </div>
    </Link>
  );
}
