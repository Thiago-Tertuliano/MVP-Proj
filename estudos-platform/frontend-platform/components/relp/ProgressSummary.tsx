import * as React from "react";

import { Progress } from "@/components/ui/progress";
import { cn } from "@/lib/utils";

export type ProgressSummaryProps = {
  concluidos: number;
  total: number;
  /** Se ausente, calcula a partir de concluidos/total. */
  percentual?: number;
  /** Texto à esquerda; padrão "X de Y artigos lidos". */
  label?: string;
  className?: string;
};

/** Barra verde + "3 de 12 artigos lidos" + %. A cor é sempre `done` (educação), nunca o azul de marca. */
export function ProgressSummary({ concluidos, total, percentual, label, className }: ProgressSummaryProps) {
  const pct = Math.round(percentual ?? (total > 0 ? (concluidos / total) * 100 : 0));
  const texto = label ?? `${concluidos} de ${total} ${total === 1 ? "artigo lido" : "artigos lidos"}`;

  return (
    <div className={cn("space-y-1.5", className)}>
      <div className="flex items-baseline justify-between gap-2 text-xs">
        <span className="text-muted-foreground">{texto}</span>
        <span className="font-semibold tabular-nums text-done-fg">{pct}%</span>
      </div>
      <Progress value={pct} aria-label={`Progresso: ${texto}`} />
    </div>
  );
}
