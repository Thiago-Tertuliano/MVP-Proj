"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { Check, CircleCheck } from "lucide-react";

import { Button } from "@/components/ui/button";
import { toast } from "@/components/ui/sonner";
import { GuestPrompt } from "@/components/relp/GuestPrompt";
import { ApiError } from "@/lib/api";
import { marcarArtigo, registrarVisita } from "@/lib/progress-store";
import { proximoNaoLido } from "@/lib/roadmap";
import { useSession } from "@/lib/session";
import { useProgressoTrilha } from "@/lib/use-progresso";

export type ArticleProgressProps = {
  artigoId: string;
  artigoSlug: string;
  trilhaId: string | null;
  /** Artigos da trilha na ordem de leitura (para sugerir o próximo não lido). */
  ordenados: { id: string; slug: string; titulo: string }[];
};

/**
 * Botão "Marcar como lido" com optimistic UI: o estado muda na hora; se a API falhar, volta e avisa.
 * Visitante vê um convite para entrar em vez do botão.
 */
export function ArticleProgress({ artigoId, artigoSlug, trilhaId, ordenados }: ArticleProgressProps) {
  const { status } = useSession();
  const router = useRouter();
  const entry = useProgressoTrilha(trilhaId);
  const [otimista, setOtimista] = React.useState<boolean | null>(null);
  const [pendente, setPendente] = React.useState(false);
  const visitaRegistrada = React.useRef(false);

  const doStore = entry?.data ? entry.data.artigos_concluidos.includes(artigoId) : false;
  const concluido = otimista ?? doStore;

  // Quando o store chega com a confirmação, o estado local não é mais necessário.
  React.useEffect(() => {
    if (otimista !== null && !pendente) setOtimista(null);
  }, [doStore, otimista, pendente]);

  // Registra a visita ("em curso") para alimentar o "Continuar de onde parou".
  // Só depois de saber o progresso, para nunca desmarcar um artigo já lido.
  React.useEffect(() => {
    if (status !== "user" || !trilhaId || visitaRegistrada.current) return;
    if (entry?.status !== "ready" || !entry.data) return;
    visitaRegistrada.current = true;
    if (!entry.data.artigos_concluidos.includes(artigoId)) void registrarVisita(artigoId);
  }, [status, trilhaId, entry, artigoId]);

  if (status === "loading") return <div className="h-10" aria-hidden="true" />;

  if (status === "guest") {
    return (
      <GuestPrompt
        titulo="Salve seu progresso"
        descricao="Entre para marcar artigos como lidos, acompanhar sua trilha e guardar anotações."
        next={`/artigos/${artigoSlug}`}
      />
    );
  }

  async function alternar() {
    const novo = !concluido;
    setOtimista(novo);
    setPendente(true);
    try {
      await marcarArtigo({ artigoId, trilhaId, concluido: novo });
      if (novo) {
        const concluidos = new Set([...(entry?.data?.artigos_concluidos ?? []), artigoId]);
        const proximo = proximoNaoLido(ordenados, concluidos, artigoSlug);
        toast.success("Artigo marcado como lido.", {
          action: proximo
            ? { label: "Próximo", onClick: () => router.push(`/artigos/${proximo.slug}`) }
            : undefined,
          description: proximo ? undefined : "Você leu todos os artigos desta trilha. Parabéns!",
        });
      }
    } catch (err) {
      setOtimista(null);
      // 401 já derrubou a sessão (toast próprio); aqui só os demais erros.
      if (!(err instanceof ApiError) || !err.isUnauthorized) {
        toast.error("Não foi possível salvar seu progresso.", {
          description: err instanceof ApiError ? err.message : "Tente novamente.",
        });
      }
    } finally {
      setPendente(false);
    }
  }

  return (
    <Button
      type="button"
      variant={concluido ? "success" : "default"}
      onClick={alternar}
      aria-pressed={concluido}
      loading={false}
      className="min-w-44"
    >
      {concluido ? <CircleCheck /> : <Check />}
      {concluido ? "Lido — desfazer" : "Marcar como lido"}
    </Button>
  );
}
