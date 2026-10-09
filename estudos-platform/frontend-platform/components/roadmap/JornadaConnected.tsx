"use client";

import * as React from "react";

import { JornadaView, type JornadaRoadmap } from "@/components/roadmap/JornadaView";
import { GuestPrompt } from "@/components/relp/GuestPrompt";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { carregarGamificacao } from "@/lib/roadmaps/gamificacao-store";
import { listarProgressoRoadmaps, listarRoadmaps } from "@/lib/roadmaps/client";
import { useGamificacao } from "@/lib/roadmaps/hooks";
import { useSession } from "@/lib/session";

/** /jornada: exige sessão. Visitante vê convite; logado vê XP, sequência, roadmaps e conquistas. */
export function JornadaConnected() {
  const { status } = useSession();
  const entry = useGamificacao();
  const [roadmaps, setRoadmaps] = React.useState<JornadaRoadmap[]>([]);

  React.useEffect(() => {
    if (status !== "user") return;
    const ctrl = new AbortController();
    Promise.all([listarRoadmaps(), listarProgressoRoadmaps(ctrl.signal)])
      .then(([lista, progresso]) => {
        const porSlug = new Map(progresso.map((p) => [p.slug, p]));
        setRoadmaps(
          lista
            .map((r) => {
              const p = porSlug.get(r.slug);
              return { slug: r.slug, titulo: r.titulo, concluidos: p?.concluidos ?? 0, total: p?.total ?? r.total_nos, percentual: p?.percentual ?? 0 };
            })
            // Só os que o aluno já tocou; o catálogo completo está em /roadmaps.
            .filter((r) => r.concluidos > 0),
        );
      })
      .catch(() => undefined);
    return () => ctrl.abort();
  }, [status]);

  if (status === "guest") {
    return (
      <GuestPrompt
        titulo="Entre para acompanhar sua jornada"
        descricao="XP, nível, sequência e conquistas ficam salvos na sua conta."
        next="/jornada"
      />
    );
  }

  if (!entry || (!entry.data && entry.status !== "error")) {
    return (
      <div className="space-y-4" role="status" aria-label="Carregando sua jornada">
        <div className="grid gap-4 md:grid-cols-2">
          <Skeleton className="h-36 rounded-xl" />
          <Skeleton className="h-36 rounded-xl" />
        </div>
        <Skeleton className="h-48 rounded-xl" />
        <span className="sr-only">Carregando…</span>
      </div>
    );
  }

  if (!entry.data) {
    return (
      <Alert variant="danger">
        <AlertTitle>Não foi possível carregar sua jornada</AlertTitle>
        <AlertDescription>
          <Button type="button" size="sm" variant="outline" onClick={() => carregarGamificacao({ force: true })}>
            Tentar de novo
          </Button>
        </AlertDescription>
      </Alert>
    );
  }

  return <JornadaView gamificacao={entry.data} roadmaps={roadmaps} />;
}
