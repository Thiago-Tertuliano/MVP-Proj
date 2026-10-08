import * as React from "react";
import Link from "next/link";
import { Check } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { ProgressSummary } from "@/components/relp/ProgressSummary";
import type { NoEstado, RoadmapModuloData, RoadmapNoData } from "@/lib/roadmap";
import { cn } from "@/lib/utils";

const ESTADO_TEXTO: Record<NoEstado, string> = {
  lido: "Lido",
  atual: "Próximo a ler",
  futuro: "Ainda não lido",
};

export type RoadmapNodeProps = {
  no: RoadmapNoData;
  href?: string;
  /** Último nó do módulo: não desenha a linha de conexão abaixo. */
  ultimo?: boolean;
  /** Visitante não tem estados (tudo "futuro"): esconde o rótulo para não parecer "pendente". */
  mostrarEstado?: boolean;
};

export function RoadmapNode({ no, href, ultimo, mostrarEstado = true }: RoadmapNodeProps) {
  const { estado } = no;
  return (
    <li className="relative flex gap-4 pb-6 last:pb-0" data-estado={estado}>
      {!ultimo && (
        <span
          aria-hidden="true"
          className={cn(
            "absolute left-[15px] top-8 h-[calc(100%-2rem)] w-0.5",
            estado === "lido" ? "bg-done" : "bg-border",
          )}
        />
      )}

      <span
        aria-hidden="true"
        className={cn(
          "relative z-10 flex h-8 w-8 shrink-0 items-center justify-center rounded-full border-2 text-sm font-semibold",
          estado === "lido" && "border-done-fg bg-done-fg text-primary-foreground",
          estado === "atual" && "border-primary bg-primary-muted text-primary ring-4 ring-primary/20",
          estado === "futuro" && "border-border bg-card text-muted-foreground",
        )}
      >
        {estado === "lido" ? <Check className="h-4 w-4" /> : no.indice}
      </span>

      <Link
        href={href ?? `/artigos/${no.slug}`}
        aria-current={estado === "atual" ? "step" : undefined}
        className={cn(
          "flex min-h-8 flex-1 flex-wrap items-center gap-x-3 gap-y-1 rounded-md py-1 text-sm font-medium transition-colors hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background",
          estado === "futuro" ? "text-foreground" : estado === "lido" ? "text-foreground" : "text-primary",
        )}
      >
        <span>{no.titulo}</span>
        {mostrarEstado && estado === "atual" && <Badge variant="default">Continue aqui</Badge>}
        {mostrarEstado && estado === "lido" && <Badge variant="done">Lido</Badge>}
        {mostrarEstado && <span className="sr-only">{`— ${ESTADO_TEXTO[estado]}`}</span>}
      </Link>
    </li>
  );
}

export type RoadmapMapProps = {
  modulos: RoadmapModuloData[];
  /** Cabeçalho de progresso total (só logado). */
  progresso?: { concluidos: number; total: number; percentual?: number } | null;
  /** `false` para visitante: sem rótulos de estado. */
  mostrarEstado?: boolean;
  className?: string;
};

/** Mapa da trilha: módulos em sequência, cada artigo é um nó (lido / atual / futuro). */
export function RoadmapMap({ modulos, progresso, mostrarEstado = true, className }: RoadmapMapProps) {
  return (
    <div className={cn("space-y-8", className)}>
      {progresso && (
        <div className="rounded-xl border border-border bg-card p-4">
          <ProgressSummary
            concluidos={progresso.concluidos}
            total={progresso.total}
            percentual={progresso.percentual}
          />
        </div>
      )}

      {modulos.map((modulo, i) => {
        const lidos = modulo.nos.filter((n) => n.estado === "lido").length;
        const headingId = `modulo-${modulo.slug}-${i}`;
        return (
          <section key={`${modulo.slug}-${i}`} aria-labelledby={headingId} className="space-y-4">
            <header className="space-y-1">
              <div className="flex flex-wrap items-baseline gap-x-3">
                <h2 id={headingId} className="text-xl font-semibold tracking-tight">
                  <span className="mr-2 text-muted-foreground">{i + 1}.</span>
                  {modulo.titulo}
                </h2>
                {mostrarEstado && progresso && (
                  <span className="text-xs tabular-nums text-muted-foreground">
                    {lidos}/{modulo.nos.length} lidos
                  </span>
                )}
              </div>
              {modulo.descricao && <p className="text-sm text-muted-foreground">{modulo.descricao}</p>}
            </header>
            <ol className="rounded-xl border border-border bg-card p-5">
              {modulo.nos.map((no, idx) => (
                <RoadmapNode
                  key={no.id}
                  no={no}
                  ultimo={idx === modulo.nos.length - 1}
                  mostrarEstado={mostrarEstado}
                />
              ))}
            </ol>
          </section>
        );
      })}
    </div>
  );
}
