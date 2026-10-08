import * as React from "react";
import Link from "next/link";
import { ArrowRight, BookOpen } from "lucide-react";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export type ContinueCardProps = {
  artigoTitulo: string;
  artigoSlug: string;
  trilhaTitulo?: string | null;
  /** Se o último artigo já foi concluído, o texto muda para "Último artigo lido". */
  concluido?: boolean;
  className?: string;
};

/** Card "Continue de onde parou" — último artigo tocado pelo aluno. */
export function ContinueCard({ artigoTitulo, artigoSlug, trilhaTitulo, concluido, className }: ContinueCardProps) {
  return (
    <section
      aria-labelledby="continuar-titulo"
      className={cn(
        "flex flex-col gap-4 rounded-xl border border-primary/30 bg-primary-muted p-5 sm:flex-row sm:items-center sm:justify-between",
        className,
      )}
    >
      <div className="flex min-w-0 items-start gap-3">
        <span
          aria-hidden="true"
          className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground"
        >
          <BookOpen className="h-5 w-5" />
        </span>
        <div className="min-w-0">
          <h2 id="continuar-titulo" className="text-xs font-semibold uppercase tracking-wide text-primary">
            {concluido ? "Último artigo lido" : "Continue de onde parou"}
          </h2>
          <p className="truncate text-base font-semibold text-foreground">{artigoTitulo}</p>
          {trilhaTitulo && <p className="truncate text-sm text-muted-foreground">{trilhaTitulo}</p>}
        </div>
      </div>
      <Button asChild className="shrink-0">
        <Link href={`/artigos/${artigoSlug}`}>
          {concluido ? "Reler" : "Continuar"}
          <ArrowRight />
        </Link>
      </Button>
    </section>
  );
}
