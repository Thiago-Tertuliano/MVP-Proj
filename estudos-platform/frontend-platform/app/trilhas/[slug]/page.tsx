import Link from "next/link";
import { notFound } from "next/navigation";
import { RoadmapMap } from "@/components/relp/RoadmapMap";
import { getTrilha } from "@/lib/mock-data";

type Props = { params: { slug: string } };

export default function TrilhaPage({ params }: Props) {
  const trilha = getTrilha(params.slug);
  if (!trilha) notFound();

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
            {trilha.progressoPct}% concluído
          </span>
        </div>
        <p className="text-muted-foreground">{trilha.descricao}</p>
        <div className="h-2 max-w-md overflow-hidden rounded-full bg-progress-track">
          <div
            className="h-full rounded-full bg-done"
            style={{ width: `${trilha.progressoPct}%` }}
          />
        </div>
      </header>

      <RoadmapMap trilha={trilha} />
    </div>
  );
}
