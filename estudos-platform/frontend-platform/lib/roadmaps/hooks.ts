"use client";

import * as React from "react";

import { carregarGamificacao, getGamificacaoEntry, subscribeGamificacao, type GamificacaoEntry } from "@/lib/roadmaps/gamificacao-store";
import {
  carregarProgressoRoadmap,
  getProgressoRoadmapEntry,
  subscribeProgressoRoadmap,
  type RoadmapProgressoEntry,
} from "@/lib/roadmaps/progresso-store";
import { useSession } from "@/lib/session";

const ENTRY_VISITANTE: GamificacaoEntry = { status: "idle", data: null };

/** Resumo de XP/nível/streak do aluno logado. Visitante recebe `undefined` e nenhuma chamada é feita. */
export function useGamificacao(): GamificacaoEntry | undefined {
  const { status } = useSession();
  const habilitado = status === "user";

  const entry = React.useSyncExternalStore(subscribeGamificacao, getGamificacaoEntry, () => ENTRY_VISITANTE);

  React.useEffect(() => {
    if (habilitado) void carregarGamificacao();
  }, [habilitado]);

  return habilitado ? entry : undefined;
}

/** Estados dos nós de um roadmap para o aluno logado. */
export function useProgressoRoadmap(slug: string | null | undefined): RoadmapProgressoEntry | undefined {
  const { status } = useSession();
  const habilitado = status === "user" && !!slug;

  const entry = React.useSyncExternalStore(
    subscribeProgressoRoadmap,
    () => (slug ? getProgressoRoadmapEntry(slug) : undefined),
    () => undefined,
  );

  React.useEffect(() => {
    if (habilitado && slug) void carregarProgressoRoadmap(slug);
  }, [habilitado, slug]);

  return habilitado ? entry : undefined;
}
