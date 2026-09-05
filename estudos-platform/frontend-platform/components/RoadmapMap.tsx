import Link from "next/link";
import type { Trilha } from "@/lib/types";

type Props = { trilha: Trilha };

export function RoadmapMap({ trilha }: Props) {
  return (
    <div className="grid gap-8 md:grid-cols-2">
      {trilha.modulos.map((modulo) => (
        <section key={modulo.slug} className="card">
          <div className="mb-4 border-b border-border pb-3">
            <h2 className="font-semibold text-fg">{modulo.titulo}</h2>
            <p className="text-sm text-muted">{modulo.descricao}</p>
          </div>
          <ul className="flex flex-col gap-2">
            {modulo.artigos.map((artigo, idx) => (
              <li key={artigo.slug}>
                <Link
                  href={`/artigos/${artigo.slug}`}
                  className={`flex items-center gap-3 rounded-relp-md border px-4 py-3 text-sm transition ${
                    artigo.concluido
                      ? "border-done/30 bg-done-muted/50 text-fg"
                      : "border-border text-fg hover:border-primary/50 hover:bg-primary-muted/50"
                  }`}
                >
                  <span
                    className={`flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-xs font-bold ${
                      artigo.concluido
                        ? "bg-done/20 text-done"
                        : "bg-progress-track text-muted"
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
