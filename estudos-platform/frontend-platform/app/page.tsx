import { Map } from "lucide-react";

import { ContinueSection } from "@/components/relp/ContinueSection";
import { TrilhaCardConnected } from "@/components/relp/TrilhaCardConnected";
import { EmptyState } from "@/components/ui/empty-state";
import { listarTrilhas } from "@/lib/content";

// Dados vivos da API: sem cache estático (o layout já depende de sessão no cliente).
export const dynamic = "force-dynamic";

export default async function HomePage() {
  // Falha de rede/5xx propaga para app/error.tsx (estado de erro com "Tentar de novo").
  const trilhas = await listarTrilhas();

  return (
    <div className="space-y-8">
      <header className="space-y-2">
        <h1 className="text-3xl font-bold tracking-tight text-foreground">Trilhas de estudo</h1>
        <p className="max-w-2xl text-muted-foreground">
          Siga um mapa de módulos e artigos, marque o que já leu e continue de onde parou.
        </p>
      </header>

      <ContinueSection />

      <section aria-labelledby="trilhas-titulo" className="space-y-4">
        <h2 id="trilhas-titulo" className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">
          Trilhas disponíveis
        </h2>
        {trilhas.length === 0 ? (
          <EmptyState
            icon={<Map />}
            title="Nenhuma trilha publicada ainda"
            description="Assim que novas trilhas forem publicadas, elas aparecem aqui. Volte em breve!"
          />
        ) : (
          <ul className="grid gap-4 md:grid-cols-2">
            {trilhas.map((trilha) => (
              <li key={trilha.id}>
                <TrilhaCardConnected trilha={trilha} />
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
