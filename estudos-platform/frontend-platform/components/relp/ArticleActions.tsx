"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

type Props = {
  artigoId: string;
  trilhaSlug?: string;
};

export function ArticleActions({ artigoId, trilhaSlug }: Props) {
  const [isPending, setIsPending] = useState(false);
  const router = useRouter();

  async function marcarComoLido() {
    setIsPending(true);
    try {
      const res = await fetch(`http://localhost:8080/api/v1/progresso/artigos/${artigoId}`, {
        method: "PUT",
        // No futuro, passaremos aqui o cabeçalho de Autorização com o token JWT
      });

      if (!res.ok) throw new Error("Falha ao registar o progresso.");

      // Força a atualização dos dados da página no Next.js (para a barra verde atualizar)
      router.refresh(); 
      
      // Se estivermos dentro de uma trilha, redireciona de volta para o mapa
      if (trilhaSlug) {
        router.push(`/trilhas/${trilhaSlug}`);
      } else {
        alert("Artigo marcado como lido com sucesso!");
      }
    } catch (error) {
      console.error(error);
      alert("Ocorreu um erro ao guardar o progresso.");
    } finally {
      setIsPending(false);
    }
  }

  return (
    <footer className="flex flex-wrap gap-3 border-t border-border pt-6">
      {trilhaSlug && (
        <Link href={`/trilhas/${trilhaSlug}`} className="btn-ghost">
          Voltar ao mapa
        </Link>
      )}
      <button
        type="button"
        className="btn-primary disabled:opacity-50"
        onClick={marcarComoLido}
        disabled={isPending}
      >
        {isPending ? "A guardar..." : "Marcar como lido"}
      </button>
    </footer>
  );
}