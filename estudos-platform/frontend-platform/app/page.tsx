import { TrilhaCard } from "@/components/relp/TrilhaCard";
import type { ListarTrilhasResponse, Trilha } from "@/lib/types";

// Função para buscar os dados reais da API
async function getTrilhas(): Promise<Trilha[]> {
  try {
    const res = await fetch("http://localhost:8080/api/v1/trilhas", {
      cache: "no-store", // Garante que a tela atualize logo após rodar o content-job
    });
    
    if (!res.ok) throw new Error("Falha ao buscar trilhas da API");
    
    const data: ListarTrilhasResponse = await res.json();
    return data.itens || [];
  } catch (error) {
    console.error("Erro na integração:", error);
    return [];
  }
}

export default async function HomePage() {
  const todasTrilhas = await getTrilhas();
  
  // Filtra as trilhas baseadas no status que vem do banco
  const publicadas = todasTrilhas.filter((t) => t.publicada);
  const rascunhos = todasTrilhas.filter((t) => !t.publicada);

  return (
    <div className="space-y-8">
      <header className="space-y-2">
        <p className="text-sm font-medium uppercase tracking-wider text-primary">MVP · leitura</p>
        <h1 className="text-3xl font-bold tracking-tight text-fg">Suas trilhas</h1>
        <p className="max-w-2xl text-muted-foreground">
          Mapa estilo roadmap: módulos como regiões, artigos como nós. Dados integrados 
          diretamente com a API em Go.
        </p>
      </header>

      <section className="space-y-4">
        <h2 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">
          Publicadas
        </h2>
        {publicadas.length === 0 ? (
          <p className="text-sm text-muted-foreground">
            Nenhuma trilha publicada no momento.
          </p>
        ) : (
          <div className="grid gap-4 md:grid-cols-2">
            {publicadas.map((trilha) => (
              <TrilhaCard key={trilha.slug} trilha={trilha} />
            ))}
          </div>
        )}
      </section>

      {rascunhos.length > 0 && (
        <section className="space-y-4">
          <h2 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">
            No banco, ainda não públicas
          </h2>
          <p className="text-sm text-muted-foreground">
            Após o content-job, trilhas como Dados nascem rascunho até o Bruno publicar pela API.
          </p>
          <div className="grid gap-4 opacity-60 md:grid-cols-2">
            {rascunhos.map((trilha) => (
              <TrilhaCard key={trilha.slug} trilha={trilha} />
            ))}
          </div>
        </section>
      )}
    </div>
  );
}