import * as React from "react";
import Link from "next/link";
import { ArrowRight, Layers } from "lucide-react";

import { Skeleton } from "@/components/ui/skeleton";
import { ProgressSummary } from "@/components/relp/ProgressSummary";
import { cn } from "@/lib/utils";

export type TrilhaCardProps = {
  slug: string;
  titulo: string;
  descricao?: string;
  totalModulos?: number;
  /** Só para logado. `loading` mostra skeleton da barra; ausente = visitante (sem barra). */
  progresso?: { status: "loading" | "ready" | "error"; concluidos?: number; total?: number; percentual?: number };
  className?: string;
};

/** Card de trilha da home. Todo o card é um link (um único foco de teclado). */
export function TrilhaCard({ slug, titulo, descricao, totalModulos, progresso, className }: TrilhaCardProps) {
  return (
    <Link
      href={`/trilhas/${slug}`}
      className={cn(
        "group flex h-full flex-col gap-4 rounded-xl border border-border bg-card p-5 text-card-foreground transition-colors hover:border-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background",
        className,
      )}
    >
      <div className="flex-1 space-y-2">
        <h3 className="text-lg font-semibold leading-tight tracking-tight group-hover:text-primary">{titulo}</h3>
        {descricao && <p className="line-clamp-3 text-sm text-muted-foreground">{descricao}</p>}
      </div>

      {progresso?.status === "loading" && <Skeleton className="h-8 w-full" />}
      {progresso?.status === "ready" && progresso.total !== undefined && (
        <ProgressSummary
          concluidos={progresso.concluidos ?? 0}
          total={progresso.total}
          percentual={progresso.percentual}
        />
      )}

      <div className="flex items-center justify-between text-sm">
        <span className="inline-flex items-center gap-1.5 text-muted-foreground">
          <Layers className="h-4 w-4" aria-hidden="true" />
          {totalModulos !== undefined
            ? `${totalModulos} ${totalModulos === 1 ? "módulo" : "módulos"}`
            : "Trilha"}
        </span>
        <span className="inline-flex items-center gap-1 font-medium text-primary">
          Ver trilha
          <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-0.5" aria-hidden="true" />
        </span>
      </div>
    </Link>
  );
}
