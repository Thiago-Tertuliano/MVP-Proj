import Link from "next/link";
import { notFound } from "next/navigation";
import { ArticleRenderer } from "@/components/relp/ArticleRenderer";
import { QuizPanel } from "@/components/relp/QuizPanel";
import type { Artigo, Trilha } from "@/lib/types";

type Props = { params: { slug: string } };

// Busca o artigo real da API Go
async function getArtigo(slug: string): Promise<Artigo | null> {
  try {
    const res = await fetch(`http://localhost:8080/api/v1/artigos/${slug}`, {
      cache: "no-store", 
    });
    if (res.status === 404) return null;
    if (!res.ok) throw new Error("Falha ao buscar artigo");
    return await res.json();
  } catch (error) {
    console.error("Erro na integração do artigo:", error);
    return null;
  }
}

// Busca a trilha real para montar o Breadcrumb (navegação de retorno)
async function getTrilhaContext(slug?: string): Promise<Trilha | null> {
  if (!slug) return null;
  try {
    const res = await fetch(`http://localhost:8080/api/v1/trilhas/${slug}`, {
      cache: "no-store",
    });
    if (!res.ok) return null;
    return await res.json();
  } catch (error) {
    return null;
  }
}

export default async function ArtigoPage({ params }: Props) {
  const artigo = await getArtigo(params.slug);
  if (!artigo) notFound();

  const trilha = await getTrilhaContext(artigo.trilhaSlug);
  const questoes = artigo.metadados?.quiz?.questoes ?? [];

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
              {artigo.metadados?.origem && (
                <span className="text-xs text-muted-foreground">{artigo.metadados.origem}</span>
              )}
            </div>
            <h1 className="text-3xl font-bold text-fg">{artigo.titulo}</h1>
            {artigo.metadados?.objetivo && (
              <p className="text-muted-foreground">{artigo.metadados.objetivo}</p>
            )}
            {artigo.metadados?.tempo_leitura_min && (
              <p className="text-xs text-muted-foreground">
                ~{artigo.metadados.tempo_leitura_min} min de leitura
              </p>
            )}
          </header>

          <ArticleRenderer blocks={artigo.conteudo?.blocks || []} />

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