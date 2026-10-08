import * as React from "react";
import { Loader2 } from "lucide-react";

import { cn } from "@/lib/utils";

/** Bloco de carregamento: ocupa o espaço do conteúdo final para evitar "flicker" de tela branca. */
function Skeleton({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return <div aria-hidden="true" className={cn("animate-pulse rounded-md bg-muted", className)} {...props} />;
}

/** Spinner para ações pontuais (botão, salvar). Para páginas, prefira `Skeleton`. */
function Spinner({ className, label = "Carregando…" }: { className?: string; label?: string }) {
  return (
    <span role="status" className="inline-flex items-center">
      <Loader2 className={cn("h-4 w-4 animate-spin text-muted-foreground", className)} aria-hidden="true" />
      <span className="sr-only">{label}</span>
    </span>
  );
}

export { Skeleton, Spinner };
