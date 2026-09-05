import Link from "next/link";
import type { Trilha } from "@/lib/types";

type Props = { trilha: Trilha };

export function RoadmapMap({ trilha }: Props) {
  return (
    <div className="grid gap-8 md:grid-cols-2">
      {trilha.modulos.map((modulo) => (
        <section key={modulo.slug} className="card">
          <div className="mb-4 border-b border-surface-border pb-3">
            <h2 className="font-semibold">{modulo.titulo}</h2>
            <p className="text-sm text-zinc-400">{modulo.descricao}</p>
          </div>
          <ul className="flex flex-col gap-2">
            {modulo.artigos.map((artigo, idx) => (
              <li key={artigo.slug}>
                <Link
                  href={`/artigos/${artigo.slug}`}
                  className={`flex items-center gap-3 rounded-lg border px-4 py-3 text-sm transition ${
                    artigo.concluido
                      ? "border-success/30 bg-success/5 text-zinc-200"
                      : "border-surface-border hover:border-accent/50 hover:bg-accent-muted/30"
                  }`}
                >
                  <span
                    className={`flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-xs font-bold ${
                      artigo.concluido ? "bg-success/20 text-success" : "bg-zinc-800 text-zinc-400"
                    }`}
                  >
                    {artigo.concluido ? "L" : idx + 1}
                  </span>
                  <span>{artigo.titulo}</span>
                </Link>
              </li>
            ))}
          </ul>
        </section>
      ))}
    </div>
  );
}
