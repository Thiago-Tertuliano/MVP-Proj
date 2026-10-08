"use client";

import * as React from "react";
import Link from "next/link";
import { ArrowRight } from "lucide-react";

import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { RoadmapMap } from "@/components/relp/RoadmapMap";
import { carregarProgresso } from "@/lib/progress-store";
import { artigosEmOrdem, montarRoadmap, proximoNaoLido } from "@/lib/roadmap";
import { useSession } from "@/lib/session";
import type { Artigo, Trilha } from "@/lib/types";
import { useProgressoTrilha } from "@/lib/use-progresso";

/**
 * Mapa da trilha. O conteúdo vem do servidor; o progresso (lido/atual) é aplicado no cliente
 * quando há sessão. Visitante vê o mapa sem estados.
 */
export function TrilhaRoadmap({ trilha, artigos }: { trilha: Trilha; artigos: Artigo[] }) {
  const { status } = useSession();
  const entry = useProgressoTrilha(trilha.id);

  const concluidos = React.useMemo(
    () => (entry?.data ? new Set(entry.data.artigos_concluidos) : null),
    [entry?.data],
  );
  const modulos = React.useMemo(() => montarRoadmap(trilha, artigos, concluidos), [trilha, artigos, concluidos]);
  const ordenados = React.useMemo(() => artigosEmOrdem(trilha, artigos), [trilha, artigos]);
  const proximo = concluidos ? proximoNaoLido(ordenados, concluidos) : null;

  const aguardando = status === "loading" || (status === "user" && (!entry || (entry.status === "loading" && !entry.data)));
  const falhou = status === "user" && entry?.status === "error" && !entry.data;

  return (
    <div className="space-y-6">
      {aguardando && <Skeleton className="h-16 w-full rounded-xl" />}

      {falhou && (
        <Alert variant="warning">
          <AlertDescription>
            <p>Não conseguimos carregar seu progresso agora. O mapa continua disponível.</p>
            <Button
              type="button"
              size="sm"
              variant="outline"
              className="mt-3"
              onClick={() => void carregarProgresso(trilha.id, { force: true })}
            >
              Tentar de novo
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {proximo && (
        <div className="flex flex-col gap-3 rounded-xl border border-primary/30 bg-primary-muted p-4 sm:flex-row sm:items-center sm:justify-between">
          <div className="min-w-0">
            <p className="text-xs font-semibold uppercase tracking-wide text-primary">Próximo a ler</p>
            <p className="truncate font-semibold text-foreground">{proximo.titulo}</p>
          </div>
          <Button asChild className="shrink-0">
            <Link href={`/artigos/${proximo.slug}`}>
              Continuar <ArrowRight />
            </Link>
          </Button>
        </div>
      )}

      {concluidos && ordenados.length > 0 && concluidos.size >= ordenados.length && (
        <Alert variant="success">
          <AlertDescription>Você concluiu esta trilha. Parabéns! Pode revisitar qualquer artigo quando quiser.</AlertDescription>
        </Alert>
      )}

      <RoadmapMap
        modulos={modulos}
        mostrarEstado={concluidos !== null}
        progresso={
          entry?.data
            ? { concluidos: entry.data.concluidos, total: entry.data.total, percentual: entry.data.percentual }
            : null
        }
      />
    </div>
  );
}
