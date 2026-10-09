import type { Metadata } from "next";
import { Swords } from "lucide-react";

import { RoadmapCatalogo } from "@/components/roadmap/RoadmapCatalogo";
import { EmptyState } from "@/components/ui/empty-state";
import { listarRoadmaps } from "@/lib/roadmaps/client";

export const metadata: Metadata = { title: "Roadmaps" };

// Dados vivos da API: sem cache estático.
export const dynamic = "force-dynamic";

export default async function RoadmapsPage() {
  const roadmaps = await listarRoadmaps();

  return (
    <div className="space-y-8">
      <header className="space-y-2">
        <h1 className="text-3xl font-bold tracking-tight text-foreground">Roadmaps</h1>
        <p className="max-w-2xl text-muted-foreground">
          Mapas de estudo gamificados: domine nós, derrote chefes, ganhe XP e mantenha sua sequência diária.
        </p>
      </header>

      {roadmaps.length === 0 ? (
        <EmptyState
          icon={<Swords />}
          title="Nenhum roadmap publicado ainda"
          description="Assim que um roadmap for publicado, ele aparece aqui. Volte em breve!"
        />
      ) : (
        <RoadmapCatalogo roadmaps={roadmaps} />
      )}
    </div>
  );
}
