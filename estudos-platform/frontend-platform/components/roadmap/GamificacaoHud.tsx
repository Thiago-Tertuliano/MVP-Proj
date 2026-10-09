import * as React from "react";
import Link from "next/link";
import { Flame } from "lucide-react";

import { Progress } from "@/components/ui/progress";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";

export type GamificacaoHudProps = {
  nivel: number;
  xpNoNivel: number;
  xpParaProximo: number;
  streak: number;
  /** Já estudou hoje? Sem isso a chama fica apagada (sequência em risco, mas não perdida). */
  ativoHoje?: boolean;
  className?: string;
};

/** Indicador compacto do header: nível, barra de XP e sequência. Leva para a Jornada. */
export function GamificacaoHud({ nivel, xpNoNivel, xpParaProximo, streak, ativoHoje, className }: GamificacaoHudProps) {
  const pct = xpParaProximo > 0 ? (xpNoNivel / xpParaProximo) * 100 : 0;
  const resumo = `Nível ${nivel}, ${xpNoNivel} de ${xpParaProximo} XP para o próximo nível. Sequência de ${streak} ${streak === 1 ? "dia" : "dias"}${ativoHoje ? ", ativa hoje" : ""}.`;

  return (
    <Link
      href="/jornada"
      aria-label={`Minha jornada. ${resumo}`}
      className={cn(
        "flex items-center gap-3 rounded-md px-2 py-1 text-xs transition-colors hover:bg-accent focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
        className,
      )}
    >
      <span className="flex flex-col gap-1">
        <span className="flex items-baseline gap-1.5 font-semibold tabular-nums text-foreground">
          Nv {nivel}
          <span className="hidden font-normal text-muted-foreground md:inline">
            {xpNoNivel}/{xpParaProximo} XP
          </span>
        </span>
        <Progress
          value={pct}
          aria-hidden="true"
          tabIndex={-1}
          className="h-1.5 w-16 md:w-24"
          indicatorClassName="bg-notify motion-reduce:transition-none"
        />
      </span>
      <span
        className={cn(
          "inline-flex items-center gap-1 font-semibold tabular-nums",
          ativoHoje ? "text-warning-fg" : "text-muted-foreground",
        )}
        title={ativoHoje ? "Sequência ativa hoje" : "Estude hoje para manter a sequência"}
      >
        <Flame className={cn("h-4 w-4", ativoHoje && "fill-warning/40")} aria-hidden="true" />
        {streak}
      </span>
    </Link>
  );
}

export function GamificacaoHudSkeleton() {
  return (
    <div role="status" aria-label="Carregando seu nível" className="flex items-center gap-3 px-2">
      <Skeleton className="h-8 w-20" />
      <Skeleton className="h-5 w-8" />
    </div>
  );
}
