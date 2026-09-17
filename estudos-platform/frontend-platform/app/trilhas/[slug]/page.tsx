import Link from "next/link";
import { notFound } from "next/navigation";
import { RoadmapMap } from "@/components/relp/RoadmapMap";
import type { Trilha } from "@/lib/types";

type Props = { params: { slug: string } };

async function getTrilha(slug: string): Promise<Trilha | null> {
  try {
    const res = await fetch(`http://localhost:8080/api/v1/trilhas/${slug}`, {
      cache: "no-store", // Sempre busca atualizado para refletir o progresso
    });
    
    if (res.status === 404) return null;
    if (!res.ok) throw new Error("Falha ao buscar trilha detalhada");
    
    return await res.json();
  } catch (error) {
    console.error("Erro na integração:", error);
    return null;
  }
}

export default async function TrilhaPage({ params }: Props) {
  const trilha = await getTrilha(params.slug);
  if (!trilha) notFound();

  const progresso = trilha.progressoPct || 0;

  return (
    <div className="space-y-6">
      <nav className="text-sm text-muted-foreground">
        <Link href="/" className="hover:text-primary">
          Trilhas
        </Link>
        <span className="mx-2">/</span>
        <span className="text-fg">{trilha.titulo}</span>
      </nav>

      <header className="space-y-3">
        <div className="flex flex-wrap items-center gap-3">
          <h1 className="text-2xl font-bold text-fg">{trilha.titulo}</h1>
          <span className="rounded-full bg-done-muted px-2.5 py-0.5 text-xs font-medium text-done">
            {progresso}% concluído
          </span>
        </div>
        <p className="text-muted-foreground">{trilha.descricao}</p>
        <div className="h-2 max-w-md overflow-hidden rounded-full bg-progress-track">
          <div
            className="h-full rounded-full bg-done"
            style={{ width: `${progresso}%` }}
          />
        </div>
      </header>

      {/* Como RoadmapMap deve ler os módulos, certifique-se que trilha.modulos é passado */}
      <RoadmapMap trilha={{...trilha, modulos: trilha.modulos || []}} />
    </div>
  );
}