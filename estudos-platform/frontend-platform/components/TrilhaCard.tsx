import Link from "next/link";
import type { Trilha } from "@/lib/types";

type Props = { trilha: Trilha };

export function TrilhaCard({ trilha }: Props) {
  const totalArtigos = trilha.modulos.reduce((n, m) => n + m.artigos.length, 0);
  const concluidos = trilha.modulos.reduce(
    (n, m) => n + m.artigos.filter((a) => a.concluido).length,
    0,
  );

  return (
    <article className="card flex flex-col gap-4 transition hover:border-accent/40">
      <div className="flex items-start justify-between gap-4">
        <div>
          <h2 className="text-lg font-semibold">{trilha.titulo}</h2>
          <p className="mt-1 text-sm text-zinc-400">{trilha.descricao}</p>
        </div>
        <span className="shrink-0 rounded-full bg-accent-muted px-2.5 py-0.5 text-xs font-medium text-accent">
          {trilha.progressoPct}%
        </span>
      </div>

      <div className="h-1.5 overflow-hidden rounded-full bg-zinc-800">
        <div
          className="h-full rounded-full bg-accent transition-all"
          style={{ width: `${trilha.progressoPct}%` }}
        />
      </div>

      <p className="text-xs text-zinc-500">
        {concluidos} / {totalArtigos} nós lidos · {trilha.modulos.length} módulos
      </p>

      <Link href={`/trilhas/${trilha.slug}`} className="btn-primary w-fit">
        Abrir trilha
      </Link>
    </article>
  );
}
