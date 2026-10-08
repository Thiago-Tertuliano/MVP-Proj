"use client";

import * as React from "react";

import { NoteEditor } from "@/components/relp/NoteEditor";
import { useAnotacao } from "@/lib/use-anotacao";
import { useSession } from "@/lib/session";
import { Skeleton } from "@/components/ui/skeleton";

/** Anotações do artigo: só para logado (visitante já vê o convite em `ArticleProgress`). */
export function ArticleNotes({ artigoId }: { artigoId: string }) {
  const { status } = useSession();
  const habilitado = status === "user";
  const { texto, alterar, status: salvamento, carregando, erroCarga, tentarNovamente } = useAnotacao(
    artigoId,
    habilitado,
  );

  if (status === "loading") return <Skeleton className="h-32 w-full" />;
  if (!habilitado) return null;

  return (
    <NoteEditor
      value={texto}
      onChange={alterar}
      status={salvamento}
      loading={carregando}
      loadError={erroCarga}
      onRetry={tentarNovamente}
    />
  );
}
