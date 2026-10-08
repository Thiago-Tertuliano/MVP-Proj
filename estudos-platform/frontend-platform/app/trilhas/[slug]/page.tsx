import type { Metadata } from "next";
import { notFound } from "next/navigation";

import { TrilhaRoadmap } from "@/components/relp/TrilhaRoadmap";
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbPage, BreadcrumbSeparator } from "@/components/ui/breadcrumb";
import { EmptyState } from "@/components/ui/empty-state";
import { listarArtigosDaTrilha, obterTrilha } from "@/lib/content";
import Link from "next/link";

export const dynamic = "force-dynamic";

type Props = { params: { slug: string } };

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const trilha = await obterTrilha(params.slug).catch(() => null);
  return { title: trilha?.titulo ?? "Trilha" };
}

export default async function TrilhaPage({ params }: Props) {
  const trilha = await obterTrilha(params.slug);
  if (!trilha) notFound();

  const artigos = await listarArtigosDaTrilha(params.slug);

  return (
    <div className="space-y-6">
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link href="/">Trilhas</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{trilha.titulo}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <header className="space-y-2">
        <h1 className="text-3xl font-bold tracking-tight text-foreground">{trilha.titulo}</h1>
        {trilha.descricao && <p className="max-w-2xl text-muted-foreground">{trilha.descricao}</p>}
      </header>

      {artigos.length === 0 ? (
        <EmptyState
          title="Esta trilha ainda não tem artigos publicados"
          description="Os artigos aparecem aqui assim que forem publicados."
          action={
            <Link href="/" className="text-sm font-medium text-primary underline-offset-4 hover:underline">
              Voltar para as trilhas
            </Link>
          }
        />
      ) : (
        <TrilhaRoadmap trilha={trilha} artigos={artigos} />
      )}
    </div>
  );
}
