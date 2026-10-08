import type { Metadata } from "next";
import Link from "next/link";
import { Search, SearchX } from "lucide-react";

import { BuscaForm } from "@/components/relp/BuscaForm";
import { SearchResults } from "@/components/relp/SearchResults";
import { EmptyState } from "@/components/ui/empty-state";
import { BUSCA_MIN_CHARS, buscarArtigos } from "@/lib/content";

export const dynamic = "force-dynamic";
export const metadata: Metadata = { title: "Busca" };

type Props = { searchParams: { q?: string | string[] } };

export default async function BuscaPage({ searchParams }: Props) {
  const bruto = Array.isArray(searchParams.q) ? searchParams.q[0] : searchParams.q;
  const termo = (bruto ?? "").trim();
  const curto = termo.length > 0 && termo.length < BUSCA_MIN_CHARS;

  // Falha da API não vira "nenhum resultado": mostra erro com saída (não mascara o problema).
  let itens: Awaited<ReturnType<typeof buscarArtigos>> = [];
  let falhou = false;
  if (termo.length >= BUSCA_MIN_CHARS) {
    try {
      itens = await buscarArtigos(termo);
    } catch {
      falhou = true;
    }
  }

  return (
    <div className="space-y-6">
      <header className="space-y-3">
        <h1 className="text-3xl font-bold tracking-tight text-foreground">Busca</h1>
        <div className="max-w-xl">
          <BuscaForm termo={termo} />
        </div>
      </header>

      {!termo && (
        <EmptyState
          icon={<Search />}
          title="Busque por um assunto"
          description={`Digite pelo menos ${BUSCA_MIN_CHARS} letras para encontrar artigos publicados.`}
        />
      )}

      {curto && (
        <EmptyState
          icon={<Search />}
          title="Digite um pouco mais"
          description={`Use pelo menos ${BUSCA_MIN_CHARS} caracteres para buscar.`}
        />
      )}

      {falhou && (
        <EmptyState
          tone="danger"
          title="Não foi possível buscar agora"
          description="Tivemos um problema ao consultar os artigos. Tente novamente em instantes."
          action={
            <Link
              href={`/busca?q=${encodeURIComponent(termo)}`}
              className="text-sm font-medium text-primary underline-offset-4 hover:underline"
            >
              Tentar de novo
            </Link>
          }
        />
      )}

      {!falhou && termo.length >= BUSCA_MIN_CHARS &&
        (itens.length === 0 ? (
          <EmptyState
            icon={<SearchX />}
            title="Nenhum artigo encontrado"
            description={`Não achamos nada para “${termo}”. Tente outras palavras ou explore as trilhas.`}
            action={
              <Link href="/" className="text-sm font-medium text-primary underline-offset-4 hover:underline">
                Ver trilhas
              </Link>
            }
          />
        ) : (
          <SearchResults itens={itens} termo={termo} />
        ))}
    </div>
  );
}
