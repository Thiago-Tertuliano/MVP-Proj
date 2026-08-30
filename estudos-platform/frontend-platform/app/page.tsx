import { TrilhaCard } from "@/components/TrilhaCard";
import { trilhas, trilhasPublicadas } from "@/lib/mock-data";

export default function HomePage() {
  const publicadas = trilhasPublicadas();
  const rascunho = trilhas.filter((t) => !t.publicada);

  return (
    <div className="space-y-8">
      <header className="space-y-2">
        <p className="text-sm font-medium uppercase tracking-wider text-accent">MVP · leitura</p>
        <h1 className="text-3xl font-bold tracking-tight">Suas trilhas</h1>
        <p className="max-w-2xl text-zinc-400">
          Mapa estilo roadmap: módulos como regiões, artigos como nós. Dados mockados para
          validar fidelidade visual antes de ligar na API.
        </p>
      </header>

      <section className="space-y-4">
        <h2 className="text-sm font-semibold uppercase tracking-wide text-zinc-500">
          Publicadas
        </h2>
        <div className="grid gap-4 md:grid-cols-2">
          {publicadas.map((trilha) => (
            <TrilhaCard key={trilha.slug} trilha={trilha} />
          ))}
        </div>
      </section>

      {rascunho.length > 0 && (
        <section className="space-y-4">
          <h2 className="text-sm font-semibold uppercase tracking-wide text-zinc-500">
            No banco, ainda não públicas
          </h2>
          <p className="text-sm text-zinc-500">
            Após o content-job, trilhas como Dados nascem rascunho até o Bruno publicar pela API.
          </p>
          <div className="grid gap-4 opacity-60 md:grid-cols-2">
            {rascunho.map((trilha) => (
              <TrilhaCard key={trilha.slug} trilha={trilha} />
            ))}
          </div>
        </section>
      )}
    </div>
  );
}
