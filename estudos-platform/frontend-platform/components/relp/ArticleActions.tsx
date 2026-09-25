"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

type Props = {
  artigoId: string;
  trilhaSlug?: string;
};

export function ArticleActions({ artigoId, trilhaSlug }: Props) {
  const [isPending, setIsPending] = useState(false);
  const [nota, setNota] = useState("");
  const [salvandoNota, setSalvandoNota] = useState(false);
  const router = useRouter();

  // Busca a anotação existente assim que o aluno abre o artigo
  useEffect(() => {
    async function carregarAnotacao() {
      try {
        const res = await fetch(`http://localhost:8080/api/v1/artigos/${artigoId}/anotacoes`);
        if (res.ok) {
          const data = await res.json();
          // O backend devolve "{}", tentamos extrair a propriedade "texto" se ela existir
          const parsed = JSON.parse(data.conteudo || "{}");
          if (parsed.texto) setNota(parsed.texto);
        }
      } catch (e) {
        console.error("Erro ao carregar anotação:", e);
      }
    }
    carregarAnotacao();
  }, [artigoId]);

  async function salvarAnotacao() {
    setSalvandoNota(true);
    try {
      const payload = { conteudo: { texto: nota } };
      const res = await fetch(`http://localhost:8080/api/v1/artigos/${artigoId}/anotacoes`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      
      if (!res.ok) throw new Error("Falha ao salvar a anotação");
      // Pequeno feedback visual sem interromper a leitura
      alert("Anotação guardada com sucesso!"); 
    } catch (error) {
      console.error(error);
      alert("Erro ao guardar a anotação.");
    } finally {
      setSalvandoNota(false);
    }
  }

  async function marcarComoLido() {
    setIsPending(true);
    try {
      const res = await fetch(`http://localhost:8080/api/v1/progresso/artigos/${artigoId}`, {
        method: "PUT",
      });

      if (!res.ok) throw new Error("Falha ao registar o progresso.");

      router.refresh(); 
      
      if (trilhaSlug) {
        router.push(`/trilhas/${trilhaSlug}`);
      } else {
        alert("Artigo marcado como lido!");
      }
    } catch (error) {
      console.error(error);
      alert("Ocorreu um erro ao guardar o progresso.");
    } finally {
      setIsPending(false);
    }
  }

  return (
    <div className="space-y-8 border-t border-border pt-8">
      
      {/* Bloco de Anotações */}
      <div className="space-y-3 rounded-relp-md border border-border bg-bg-muted p-5">
        <div className="flex items-center justify-between">
          <h3 className="text-sm font-semibold uppercase tracking-wide text-fg">As tuas Anotações</h3>
          <button
            type="button"
            className="text-xs font-medium text-primary hover:underline disabled:opacity-50"
            onClick={salvarAnotacao}
            disabled={salvandoNota}
          >
            {salvandoNota ? "A guardar..." : "Guardar anotação"}
          </button>
        </div>
        <textarea
          className="min-h-[120px] w-full rounded-md border border-border bg-bg p-3 text-sm text-fg focus:border-primary focus:outline-none"
          placeholder="Escreve aqui os teus resumos, grifos ou pontos-chave da leitura..."
          value={nota}
          onChange={(e) => setNota(e.target.value)}
        />
      </div>

      {/* Rodapé de Navegação / Progresso */}
      <footer className="flex flex-wrap gap-3">
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
          {isPending ? "A registar..." : "Marcar como lido"}
        </button>
      </footer>
    </div>
  );
}