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
      <nav className="text-sm text-muted-foreground">
        <Link href="/" className="hover:text-primary">
          Trilhas
        </Link>
        {trilha && (
          <>
            <span className="mx-2">/</span>
            <Link href={`/trilhas/${trilha.slug}`} className="hover:text-primary">
              {trilha.titulo}
            </Link>
          </>
        )}
        <span className="mx-2">/</span>
        <span className="text-fg">{artigo.titulo}</span>
      </nav>

      <div className="grid gap-8 lg:grid-cols-[1fr_320px]">
        <article className="space-y-6">
          <header className="space-y-2 border-b border-border pb-6">
            <div className="flex flex-wrap items-center gap-2">
              <span
                className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                  artigo.status === "publicado"
                    ? "bg-status-published/15 text-status-published"
                    : "bg-status-draft/15 text-status-draft"
                }`}
              >
                {artigo.status}
              </span>
              {artigo.metadados.origem && (
                <span className="text-xs text-muted-foreground">{artigo.metadados.origem}</span>
              )}
            </div>
            <h1 className="text-3xl font-bold text-fg">{artigo.titulo}</h1>
            {artigo.metadados.objetivo && (
              <p className="text-muted-foreground">{artigo.metadados.objetivo}</p>
            )}
            {artigo.metadados.tempo_leitura_min && (
              <p className="text-xs text-muted-foreground">
                ~{artigo.metadados.tempo_leitura_min} min de leitura
              </p>
            )}
          </header>

          <ArticleRenderer blocks={artigo.conteudo.blocks} />

          <footer className="flex flex-wrap gap-3 border-t border-border pt-6">
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
