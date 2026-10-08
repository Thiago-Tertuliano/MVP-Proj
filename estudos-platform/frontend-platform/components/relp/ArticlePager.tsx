import * as React from "react";
import Link from "next/link";
import { ChevronLeft, ChevronRight } from "lucide-react";

import { cn } from "@/lib/utils";

type Vizinho = { slug: string; titulo: string } | null;

export type ArticlePagerProps = {
  anterior: Vizinho;
  proximo: Vizinho;
  className?: string;
};

const BASE =
  "flex min-w-0 flex-1 flex-col gap-0.5 rounded-lg border border-border bg-card p-3 text-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background";

/** Anterior / Próximo dentro da trilha. Nas pontas o lado fica desabilitado (não some, para não mover o layout). */
export function ArticlePager({ anterior, proximo, className }: ArticlePagerProps) {
  return (
    <nav aria-label="Navegação entre artigos" className={cn("flex gap-3", className)}>
      {anterior ? (
        <Link href={`/artigos/${anterior.slug}`} rel="prev" className={cn(BASE, "hover:border-primary")}>
          <span className="flex items-center gap-1 text-xs text-muted-foreground">
            <ChevronLeft className="h-3.5 w-3.5" aria-hidden="true" />
            Anterior
          </span>
          <span className="truncate font-medium text-foreground">{anterior.titulo}</span>
        </Link>
      ) : (
        <span aria-disabled="true" className={cn(BASE, "cursor-not-allowed opacity-50")}>
          <span className="flex items-center gap-1 text-xs text-muted-foreground">
            <ChevronLeft className="h-3.5 w-3.5" aria-hidden="true" />
            Anterior
          </span>
          <span className="text-muted-foreground">Este é o primeiro artigo</span>
        </span>
      )}

      {proximo ? (
        <Link href={`/artigos/${proximo.slug}`} rel="next" className={cn(BASE, "items-end text-right hover:border-primary")}>
          <span className="flex items-center gap-1 text-xs text-muted-foreground">
            Próximo
            <ChevronRight className="h-3.5 w-3.5" aria-hidden="true" />
          </span>
          <span className="w-full truncate font-medium text-foreground">{proximo.titulo}</span>
        </Link>
      ) : (
        <span aria-disabled="true" className={cn(BASE, "cursor-not-allowed items-end text-right opacity-50")}>
          <span className="flex items-center gap-1 text-xs text-muted-foreground">
            Próximo
            <ChevronRight className="h-3.5 w-3.5" aria-hidden="true" />
          </span>
          <span className="text-muted-foreground">Este é o último artigo</span>
        </span>
      )}
    </nav>
  );
}
