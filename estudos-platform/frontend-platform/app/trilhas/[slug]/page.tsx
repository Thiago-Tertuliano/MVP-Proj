import Link from "next/link";
import { notFound } from "next/navigation";
import { RoadmapMap } from "@/components/RoadmapMap";
import { getTrilha } from "@/lib/mock-data";

type Props = { params: { slug: string } };

export default function TrilhaPage({ params }: Props) {
  const trilha = getTrilha(params.slug);
  if (!trilha) notFound();

  return (
    <div className="space-y-6">
      <nav className="text-sm text-zinc-500">
        <Link href="/" className="hover:text-accent">
          Trilhas
        </Link>
        <span className="mx-2">/</span>
        <span className="text-zinc-300">{trilha.titulo}</span>
      </nav>

      <header className="space-y-3">
        <div className="flex flex-wrap items-center gap-3">
          <h1 className="text-2xl font-bold">{trilha.titulo}</h1>
          <span className="rounded-full bg-accent-muted px-2.5 py-0.5 text-xs font-medium text-accent">
            {trilha.progressoPct}% concluído
          </span>
        </div>
        <p className="text-zinc-400">{trilha.descricao}</p>
        <div className="h-2 max-w-md overflow-hidden rounded-full bg-zinc-800">
          <div
            className="h-full rounded-full bg-accent"
            style={{ width: `${trilha.progressoPct}%` }}
          />
        </div>
      </header>

      <RoadmapMap trilha={trilha} />
    </div>
  );
}
