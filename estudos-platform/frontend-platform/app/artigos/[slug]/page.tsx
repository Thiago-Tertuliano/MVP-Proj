import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import { Clock } from "lucide-react";

import { ArticleNotes } from "@/components/relp/ArticleNotes";
import { ArticlePager } from "@/components/relp/ArticlePager";
import { ArticleProgress } from "@/components/relp/ArticleProgress";
import { ArticleRenderer } from "@/components/relp/ArticleRenderer";
import { QuizPanel } from "@/components/relp/QuizPanel";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { listarArtigosDaTrilha, listarTrilhas, obterArtigo } from "@/lib/content";
import { artigosEmOrdem, vizinhos } from "@/lib/roadmap";
import type { Artigo, Trilha } from "@/lib/types";

export const dynamic = "force-dynamic";

type Props = { params: { slug: string } };

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const artigo = await obterArtigo(params.slug).catch(() => null);
  return { title: artigo?.titulo ?? "Artigo" };
}

/** Contexto da trilha (breadcrumb + anterior/próximo). Falha aqui não derruba a leitura do artigo. */
async function contextoDaTrilha(artigo: Artigo): Promise<{ trilha: Trilha; ordenados: Artigo[] } | null> {
  if (!artigo.trilha_id) return null;
  try {
    const trilhas = await listarTrilhas();
    const trilha = trilhas.find((t) => t.id === artigo.trilha_id);
    if (!trilha) return null;
    const artigos = await listarArtigosDaTrilha(trilha.slug);
    return { trilha, ordenados: artigosEmOrdem(trilha, artigos) };
  } catch {
    return null;
  }
}

export default async function ArtigoPage({ params }: Props) {
  const artigo = await obterArtigo(params.slug);
  if (!artigo) notFound();

  const contexto = await contextoDaTrilha(artigo);
  const { anterior, proximo } = contexto
    ? vizinhos(contexto.ordenados, artigo.slug)
    : { anterior: null, proximo: null };
  const questoes = artigo.metadados?.quiz?.questoes ?? [];
  const blocks = artigo.conteudo?.blocks ?? [];

  return (
    <div className="space-y-8">
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link href="/">Trilhas</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          {contexto && (
            <>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbLink asChild>
                  <Link href={`/trilhas/${contexto.trilha.slug}`}>{contexto.trilha.titulo}</Link>
                </BreadcrumbLink>
              </BreadcrumbItem>
            </>
          )}
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{artigo.titulo}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <div className="grid gap-8 lg:grid-cols-[minmax(0,1fr)_20rem]">
        <article className="min-w-0 space-y-8">
          <header className="space-y-3 border-b border-border pb-6">
            <h1 className="text-3xl font-bold tracking-tight text-foreground">{artigo.titulo}</h1>
            {artigo.subtitulo && <p className="text-lg text-muted-foreground">{artigo.subtitulo}</p>}
            {artigo.metadados?.objetivo && (
              <p className="text-muted-foreground">
                <span className="font-medium text-foreground">Objetivo: </span>
                {artigo.metadados.objetivo}
              </p>
            )}
            {artigo.metadados?.tempo_leitura_min ? (
              <p className="flex items-center gap-1.5 text-sm text-muted-foreground">
                <Clock className="h-4 w-4" aria-hidden="true" />
                {artigo.metadados.tempo_leitura_min} min de leitura
              </p>
            ) : null}
          </header>

          <ArticleRenderer blocks={blocks} />

          <section aria-label="Seu progresso neste artigo" className="space-y-4 border-t border-border pt-6">
            <ArticleProgress
              key={artigo.id}
              artigoId={artigo.id}
              artigoSlug={artigo.slug}
              trilhaId={artigo.trilha_id}
              ordenados={contexto?.ordenados ?? []}
            />
            {contexto && <ArticlePager anterior={anterior} proximo={proximo} />}
          </section>
        </article>

        <aside aria-label="Estudo" className="min-w-0 space-y-6 lg:sticky lg:top-24 lg:self-start">
          <ArticleNotes key={artigo.id} artigoId={artigo.id} />
          <QuizPanel questoes={questoes} />
        </aside>
      </div>
    </div>
  );
}
