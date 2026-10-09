import * as React from "react";
import Link from "next/link";
import { Flame, Lock } from "lucide-react";

import { Icone } from "@/components/roadmap/icones";
import { ProgressSummary } from "@/components/relp/ProgressSummary";
import { Progress } from "@/components/ui/progress";
import type { Gamificacao } from "@/lib/roadmaps/types";
import { cn } from "@/lib/utils";

export type JornadaRoadmap = { slug: string; titulo: string; concluidos: number; total: number; percentual: number };

export type JornadaViewProps = {
  gamificacao: Gamificacao;
  /** Roadmaps em que o aluno já avançou (ou todos, com 0%). */
  roadmaps?: JornadaRoadmap[];
  className?: string;
};

function formatarData(em: number): string {
  return new Date(em * 1000).toLocaleDateString("pt-BR", { day: "2-digit", month: "short", year: "numeric" });
}

/** Painel "Minha jornada": nível, sequência, conquistas e roadmaps em andamento. Puro. */
export function JornadaView({ gamificacao: g, roadmaps = [], className }: JornadaViewProps) {
  const pct = g.xp_para_proximo > 0 ? (g.xp_no_nivel / g.xp_para_proximo) * 100 : 0;
  const conquistadas = g.conquistas.filter((c) => c.conquistada).length;

  return (
    <div className={cn("space-y-8", className)}>
      <div className="grid gap-4 md:grid-cols-2">
        <section aria-labelledby="jornada-nivel" className="space-y-4 rounded-xl border border-border bg-card p-5">
          <h2 id="jornada-nivel" className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">
            Nível
          </h2>
          <div className="flex items-center gap-4">
            <span
              aria-hidden="true"
              className="flex h-16 w-16 shrink-0 items-center justify-center rounded-full bg-notify text-3xl font-bold tabular-nums text-primary-foreground"
            >
              {g.nivel}
            </span>
            <div className="min-w-0 flex-1 space-y-1.5">
              <p className="text-lg font-semibold">
                Nível {g.nivel} <span className="text-sm font-normal text-muted-foreground">· {g.xp_total} XP no total</span>
              </p>
              <Progress
                value={pct}
                aria-label={`${g.xp_no_nivel} de ${g.xp_para_proximo} XP para o nível ${g.nivel + 1}`}
                className="h-2.5"
                indicatorClassName="bg-notify motion-reduce:transition-none"
              />
              <p className="text-xs text-muted-foreground">
                {g.xp_para_proximo - g.xp_no_nivel} XP para o nível {g.nivel + 1}
              </p>
            </div>
          </div>
        </section>

        <section aria-labelledby="jornada-streak" className="space-y-4 rounded-xl border border-border bg-card p-5">
          <h2 id="jornada-streak" className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">
            Sequência
          </h2>
          <div className="flex items-center gap-4">
            <span
              aria-hidden="true"
              className={cn(
                "flex h-16 w-16 shrink-0 items-center justify-center rounded-full",
                g.ativo_hoje ? "bg-warning-muted text-warning-fg" : "bg-muted text-muted-foreground",
              )}
            >
              <Flame className={cn("h-8 w-8", g.ativo_hoje && "fill-warning/40")} />
            </span>
            <div className="space-y-1">
              <p className="text-lg font-semibold tabular-nums">
                {g.streak_atual} {g.streak_atual === 1 ? "dia" : "dias"} seguidos
              </p>
              <p className="text-xs text-muted-foreground">
                Recorde: {g.streak_max} {g.streak_max === 1 ? "dia" : "dias"} ·{" "}
                {g.ativo_hoje ? "você já estudou hoje" : "estude hoje para manter a sequência"}
              </p>
            </div>
          </div>
        </section>
      </div>

      <section aria-labelledby="jornada-roadmaps" className="space-y-3">
        <h2 id="jornada-roadmaps" className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">
          Meus roadmaps
        </h2>
        {roadmaps.length === 0 ? (
          <p className="rounded-xl border border-dashed border-border bg-card p-5 text-sm text-muted-foreground">
            Você ainda não começou nenhum roadmap.{" "}
            <Link href="/roadmaps" className="font-medium text-primary underline-offset-4 hover:underline">
              Escolha um para começar
            </Link>
            .
          </p>
        ) : (
          <ul className="grid gap-3 md:grid-cols-2">
            {roadmaps.map((r) => (
              <li key={r.slug}>
                <Link
                  href={`/roadmaps/${r.slug}`}
                  className="block space-y-2 rounded-xl border border-border bg-card p-4 transition-colors hover:border-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                >
                  <span className="block font-medium">{r.titulo}</span>
                  <ProgressSummary
                    concluidos={r.concluidos}
                    total={r.total}
                    percentual={r.percentual}
                    label={`${r.concluidos} de ${r.total} nós dominados`}
                  />
                </Link>
              </li>
            ))}
          </ul>
        )}
      </section>

      <section aria-labelledby="jornada-conquistas" className="space-y-3">
        <div className="flex items-baseline justify-between gap-2">
          <h2 id="jornada-conquistas" className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">
            Conquistas
          </h2>
          <span className="text-xs tabular-nums text-muted-foreground">
            {conquistadas} de {g.conquistas.length}
          </span>
        </div>
        <ul className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {g.conquistas.map((c) => (
            <li
              key={c.codigo}
              data-conquistada={c.conquistada}
              className={cn(
                "flex items-start gap-3 rounded-xl border p-4",
                c.conquistada ? "border-notify/40 bg-card" : "border-dashed border-border bg-muted",
              )}
            >
              <span
                aria-hidden="true"
                className={cn(
                  "flex h-10 w-10 shrink-0 items-center justify-center rounded-full",
                  c.conquistada ? "bg-notify text-primary-foreground" : "bg-card text-muted-foreground",
                )}
              >
                {c.conquistada ? <Icone nome={c.icone} className="h-5 w-5" /> : <Lock className="h-4 w-4" />}
              </span>
              <div className="min-w-0 space-y-0.5">
                <p className={cn("text-sm font-semibold", !c.conquistada && "text-muted-foreground")}>
                  {c.nome}
                  <span className="sr-only">{c.conquistada ? " (conquistada)" : " (bloqueada)"}</span>
                </p>
                <p className="text-xs text-muted-foreground">{c.descricao}</p>
                {c.conquistada && c.em && <p className="text-xs text-notify">{formatarData(c.em)}</p>}
              </div>
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}
