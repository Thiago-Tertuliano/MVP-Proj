import Link from "next/link";
import { notFound } from "next/navigation";
import { ArticleRenderer } from "@/components/ArticleRenderer";
import { QuizPanel } from "@/components/QuizPanel";
import { getArtigo, getTrilha } from "@/lib/mock-data";

type Props = { params: { slug: string } };

export default function ArtigoPage({ params }: Props) {
  const artigo = getArtigo(params.slug);
  if (!artigo) notFound();

  const trilha = artigo.trilhaSlug ? getTrilha(artigo.trilhaSlug) : undefined;
  const questoes = artigo.metadados.quiz?.questoes ?? [];

  return (
    <div className="space-y-8">
      <nav className="text-sm text-zinc-500">
        <Link href="/" className="hover:text-accent">
          Trilhas
        </Link>
        {trilha && (
          <>
            <span className="mx-2">/</span>
            <Link href={`/trilhas/${trilha.slug}`} className="hover:text-accent">
              {trilha.titulo}
            </Link>
          </>
        )}
        <span className="mx-2">/</span>
        <span className="text-zinc-300">{artigo.titulo}</span>
      </nav>

      <div className="grid gap-8 lg:grid-cols-[1fr_320px]">
        <article className="space-y-6">
          <header className="space-y-2 border-b border-surface-border pb-6">
            <div className="flex flex-wrap items-center gap-2">
              <span
                className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                  artigo.status === "publicado"
                    ? "bg-success/15 text-success"
                    : "bg-warning/15 text-warning"
                }`}
              >
                {artigo.status}
              </span>
              {artigo.metadados.origem && (
                <span className="text-xs text-zinc-500">{artigo.metadados.origem}</span>
              )}
            </div>
            <h1 className="text-3xl font-bold">{artigo.titulo}</h1>
            {artigo.metadados.objetivo && (
              <p className="text-zinc-400">{artigo.metadados.objetivo}</p>
            )}
            {artigo.metadados.tempo_leitura_min && (
              <p className="text-xs text-zinc-500">
                ~{artigo.metadados.tempo_leitura_min} min de leitura
              </p>
            )}
          </header>

          <ArticleRenderer blocks={artigo.conteudo.blocks} />

          <footer className="flex flex-wrap gap-3 border-t border-surface-border pt-6">
            {trilha && (
              <Link href={`/trilhas/${trilha.slug}`} className="btn-ghost">
                Voltar ao mapa
              </Link>
            )}
            <button type="button" className="btn-primary">
              Marcar como lido
            </button>
          </footer>
        </article>

        {questoes.length > 0 && <QuizPanel questoes={questoes} />}
      </div>
    </div>
  );
}
