import * as React from "react";
import Link from "next/link";
import { FileText } from "lucide-react";

import type { ResultadoBusca } from "@/lib/types";
import { cn } from "@/lib/utils";

export type SearchResultsProps = {
  itens: ResultadoBusca[];
  termo: string;
  className?: string;
};

/** Lista de resultados da busca. Estados vazio/erro/curto ficam na página (EmptyState). */
export function SearchResults({ itens, termo, className }: SearchResultsProps) {
  return (
    <section aria-label={`Resultados para ${termo}`} className={cn("space-y-3", className)}>
      <p className="text-sm text-muted-foreground" role="status">
        {itens.length} {itens.length === 1 ? "resultado" : "resultados"} para{" "}
        <strong className="font-semibold text-foreground">“{termo}”</strong>
      </p>
      <ul className="space-y-2">
        {itens.map((item) => (
          <li key={item.slug}>
            <Link
              href={`/artigos/${item.slug}`}
              className="group flex items-center gap-3 rounded-xl border border-border bg-card p-4 transition-colors hover:border-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background"
            >
              <FileText className="h-5 w-5 shrink-0 text-muted-foreground group-hover:text-primary" aria-hidden="true" />
              <span className="min-w-0 flex-1 truncate font-medium text-foreground group-hover:text-primary">
                {item.titulo}
              </span>
            </Link>
          </li>
        ))}
      </ul>
    </section>
  );
}
