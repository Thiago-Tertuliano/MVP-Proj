"use client";

import * as React from "react";

import { LevelUpDialog } from "@/components/roadmap/LevelUpDialog";
import { toast } from "@/components/ui/sonner";
import { onRecompensa } from "@/lib/roadmaps/gamificacao-store";
import { celebracaoDe } from "@/lib/roadmaps/recompensa";

/**
 * Único lugar que celebra: assina o store de gamificação e mostra toasts/level-up, não importa
 * se o XP veio de um nó, do quiz do chefe ou da leitura de um artigo. Montado uma vez no layout.
 */
export function RewardHost() {
  const [nivel, setNivel] = React.useState<number | null>(null);

  React.useEffect(
    () =>
      onRecompensa((r) => {
        const c = celebracaoDe(r);
        if (c.titulo) toast.success(c.titulo, { description: c.descricao });
        if (c.roadmapsCompletos.length > 0) {
          toast.success("Roadmap completo!", { description: "Você dominou todos os nós e ganhou o bônus de conclusão." });
        }
        c.conquistas.forEach((q) => toast.success(`Conquista: ${q.nome}`, { description: q.descricao }));
        if (c.subiuDeNivel) setNivel(c.subiuDeNivel);
      }),
    [],
  );

  return <LevelUpDialog open={nivel !== null} nivel={nivel ?? 1} onClose={() => setNivel(null)} />;
}
