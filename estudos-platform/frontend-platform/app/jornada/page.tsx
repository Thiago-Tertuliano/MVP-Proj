import type { Metadata } from "next";

import { JornadaConnected } from "@/components/roadmap/JornadaConnected";

export const metadata: Metadata = { title: "Minha jornada" };

export default function JornadaPage() {
  return (
    <div className="space-y-8">
      <header className="space-y-2">
        <h1 className="text-3xl font-bold tracking-tight text-foreground">Minha jornada</h1>
        <p className="max-w-2xl text-muted-foreground">
          Seu nível, sua sequência de estudos e as conquistas desbloqueadas nos roadmaps.
        </p>
      </header>
      <JornadaConnected />
    </div>
  );
}
