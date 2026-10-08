"use client";

import * as React from "react";

import { Skeleton } from "@/components/ui/skeleton";
import { ContinueCard } from "@/components/relp/ContinueCard";
import { api } from "@/lib/api";
import { useSession } from "@/lib/session";
import type { Continuar } from "@/lib/types";

/** Home: só aparece para quem está logado e já começou algum artigo. Falhas são silenciosas (não é crítico). */
export function ContinueSection() {
  const { status } = useSession();
  const [estado, setEstado] = React.useState<{ loading: boolean; dados: Continuar["item"] }>({
    loading: false,
    dados: null,
  });

  React.useEffect(() => {
    if (status !== "user") {
      setEstado({ loading: false, dados: null });
      return;
    }
    const ac = new AbortController();
    setEstado({ loading: true, dados: null });
    api<Continuar>("/progresso/continuar", { signal: ac.signal })
      .then((r) => setEstado({ loading: false, dados: r.item }))
      .catch(() => {
        if (!ac.signal.aborted) setEstado({ loading: false, dados: null });
      });
    return () => ac.abort();
  }, [status]);

  if (status !== "user") return null;
  if (estado.loading) return <Skeleton className="h-24 w-full rounded-xl" />;
  if (!estado.dados) return null;

  return (
    <ContinueCard
      artigoSlug={estado.dados.artigo_slug}
      artigoTitulo={estado.dados.artigo_titulo}
      trilhaTitulo={estado.dados.trilha_titulo}
      concluido={estado.dados.concluido}
    />
  );
}
