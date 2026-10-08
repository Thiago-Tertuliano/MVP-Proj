"use client";

import * as React from "react";

import {
  carregarProgresso,
  getProgressoEntry,
  subscribeProgresso,
  type ProgressoEntry,
} from "@/lib/progress-store";
import { useSession } from "@/lib/session";

/**
 * Progresso do aluno logado em uma trilha. Visitante recebe `undefined` e nenhuma chamada é feita.
 * O 1º componente que pedir dispara o GET /progresso/trilhas/{id}; os demais reaproveitam o cache.
 */
export function useProgressoTrilha(trilhaId: string | null | undefined): ProgressoEntry | undefined {
  const { status } = useSession();
  const habilitado = status === "user" && !!trilhaId;

  const entry = React.useSyncExternalStore(
    subscribeProgresso,
    () => (trilhaId ? getProgressoEntry(trilhaId) : undefined),
    () => undefined,
  );

  React.useEffect(() => {
    if (habilitado && trilhaId) void carregarProgresso(trilhaId);
  }, [habilitado, trilhaId]);

  return habilitado ? entry : undefined;
}
