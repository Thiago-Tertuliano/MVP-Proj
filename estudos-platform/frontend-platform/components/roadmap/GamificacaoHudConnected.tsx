"use client";

import { GamificacaoHud, GamificacaoHudSkeleton } from "@/components/roadmap/GamificacaoHud";
import { useGamificacao } from "@/lib/roadmaps/hooks";

/** HUD ligado ao store de gamificação. Falha de rede some em silêncio: o header nunca quebra por isso. */
export function GamificacaoHudConnected() {
  const entry = useGamificacao();
  if (!entry) return null;
  if (!entry.data) return entry.status === "error" ? null : <GamificacaoHudSkeleton />;

  const g = entry.data;
  return (
    <GamificacaoHud
      nivel={g.nivel}
      xpNoNivel={g.xp_no_nivel}
      xpParaProximo={g.xp_para_proximo}
      streak={g.streak_atual}
      ativoHoje={g.ativo_hoje}
    />
  );
}
