"use client";

import * as React from "react";
import dynamic from "next/dynamic";
import { usePathname } from "next/navigation";
import { LayoutList, Map as MapIcon, Sparkles } from "lucide-react";

import { ChefeQuiz } from "@/components/roadmap/ChefeQuiz";
import { RoadmapDetalhe } from "@/components/roadmap/RoadmapDetalhe";
import { RoadmapLista } from "@/components/roadmap/RoadmapLista";
import { ProgressSummary } from "@/components/relp/ProgressSummary";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { toast } from "@/components/ui/sonner";
import { ApiError } from "@/lib/api";
import { concluirNo, responderQuiz } from "@/lib/roadmaps/client";
import { aplicarResultado } from "@/lib/roadmaps/gamificacao-store";
import { prerequisitosPendentes, proximoRecomendado } from "@/lib/roadmaps/graph";
import { useProgressoRoadmap } from "@/lib/roadmaps/hooks";
import { aplicarConclusao, carregarProgressoRoadmap } from "@/lib/roadmaps/progresso-store";
import type { ConclusaoNoResposta, QuizResultado, Roadmap } from "@/lib/roadmaps/types";
import { useSession } from "@/lib/session";

// React Flow depende de ResizeObserver/medidas do DOM: só no cliente.
const RoadmapCanvas = dynamic(() => import("@/components/roadmap/RoadmapCanvas").then((m) => m.RoadmapCanvas), {
  ssr: false,
  loading: () => <Skeleton className="h-[min(72vh,680px)] min-h-[420px] w-full rounded-xl" />,
});

type Vista = "mapa" | "lista";

function mensagem(err: unknown, padrao: string): string {
  return err instanceof ApiError ? err.message : padrao;
}

/** Conecta o roadmap (conteúdo público) ao progresso do aluno e às ações que rendem XP. */
export function RoadmapPlayer({ roadmap }: { roadmap: Roadmap }) {
  const { status } = useSession();
  const pathname = usePathname();
  const progresso = useProgressoRoadmap(roadmap.slug);

  const [vista, setVista] = React.useState<Vista>("mapa");
  const [selecionadoId, setSelecionadoId] = React.useState<string | null>(null);
  const [concluindo, setConcluindo] = React.useState(false);
  const [quizAberto, setQuizAberto] = React.useState(false);

  // Telas estreitas começam na lista (canvas pede mouse/gestos e área grande).
  React.useEffect(() => {
    if (window.matchMedia("(max-width: 639px)").matches) setVista("lista");
  }, []);

  const logado = status === "user";
  const aguardando = status === "loading" || (logado && !progresso?.data && progresso?.status !== "error");
  const estados = logado ? (progresso?.data?.estados ?? null) : null;

  const selecionado = roadmap.nos.find((n) => n.id === selecionadoId) ?? null;
  const proximo = proximoRecomendado(roadmap.nos, roadmap.arestas, estados);
  const estadoSelecionado = selecionado ? (estados ? (estados[selecionado.id] ?? "bloqueado") : null) : null;

  const selecionar = React.useCallback((id: string | null) => setSelecionadoId(id), []);

  function registrar(resp: ConclusaoNoResposta) {
    aplicarConclusao(roadmap.slug, resp);
    aplicarResultado(resp.resultado);
  }

  async function concluir() {
    if (!selecionado) return;
    setConcluindo(true);
    try {
      registrar(await concluirNo(roadmap.slug, selecionado.id));
    } catch (err) {
      toast.error(mensagem(err, "Não foi possível concluir este nó."));
      // Estado pode ter mudado em outra aba: realinha com o servidor.
      void carregarProgressoRoadmap(roadmap.slug, { force: true });
    } finally {
      setConcluindo(false);
    }
  }

  async function enviarQuiz(respostas: Record<string, string>): Promise<QuizResultado> {
    const resultado = await responderQuiz(roadmap.slug, selecionado!.id, respostas);
    if (resultado.aprovado && resultado.conclusao) registrar(resultado.conclusao);
    return resultado;
  }

  const loginHref = `/login?next=${encodeURIComponent(pathname ?? `/roadmaps/${roadmap.slug}`)}`;

  return (
    <div className="space-y-6">
      {logado && progresso?.data && (
        <div className="space-y-3 rounded-xl border border-border bg-card p-4">
          <ProgressSummary
            concluidos={progresso.data.concluidos}
            total={progresso.data.total}
            percentual={progresso.data.percentual}
            label={`${progresso.data.concluidos} de ${progresso.data.total} nós dominados`}
          />
          <div className="flex flex-wrap items-center justify-between gap-2">
            <p className="text-xs text-muted-foreground">
              {progresso.data.xp_ganho} de {progresso.data.xp_total} XP conquistados neste roadmap
            </p>
            {proximo && !progresso.data.completo && (
              <Button type="button" size="sm" onClick={() => setSelecionadoId(proximo.id)}>
                <Sparkles />
                {`${progresso.data.concluidos === 0 ? "Começar" : "Continuar"}: ${proximo.titulo}`}
              </Button>
            )}
            {progresso.data.completo && (
              <span className="text-sm font-semibold text-done-fg">Roadmap completo! 🎉</span>
            )}
          </div>
        </div>
      )}

      {logado && progresso?.status === "error" && !progresso.data && (
        <Alert variant="warning">
          <AlertTitle>Não foi possível carregar seu progresso</AlertTitle>
          <AlertDescription>
            <Button type="button" size="sm" variant="outline" onClick={() => carregarProgressoRoadmap(roadmap.slug, { force: true })}>
              Tentar de novo
            </Button>
          </AlertDescription>
        </Alert>
      )}

      <Tabs value={vista} onValueChange={(v) => setVista(v as Vista)}>
        <TabsList aria-label="Modo de visualização">
          <TabsTrigger value="mapa">
            <MapIcon className="mr-1.5 h-4 w-4" aria-hidden="true" />
            Mapa
          </TabsTrigger>
          <TabsTrigger value="lista">
            <LayoutList className="mr-1.5 h-4 w-4" aria-hidden="true" />
            Lista
          </TabsTrigger>
        </TabsList>
      </Tabs>

      <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_340px]">
        <div className="min-w-0">
          {aguardando ? (
            <Skeleton className="h-[min(72vh,680px)] min-h-[420px] w-full rounded-xl" />
          ) : vista === "mapa" ? (
            <RoadmapCanvas
              nos={roadmap.nos}
              arestas={roadmap.arestas}
              estados={estados}
              selecionadoId={selecionadoId}
              onSelecionar={selecionar}
            />
          ) : (
            <RoadmapLista
              nos={roadmap.nos}
              arestas={roadmap.arestas}
              estados={estados}
              selecionadoId={selecionadoId}
              onSelecionar={selecionar}
            />
          )}
        </div>

        <aside aria-label="Detalhes do nó" className="lg:sticky lg:top-24 lg:self-start" aria-live="polite">
          {selecionado ? (
            <RoadmapDetalhe
              no={selecionado}
              estado={estadoSelecionado}
              pendentes={prerequisitosPendentes(selecionado.id, roadmap.nos, roadmap.arestas, estados)}
              loginHref={loginHref}
              concluindo={concluindo}
              onConcluir={concluir}
              onEnfrentarChefe={() => setQuizAberto(true)}
              onSelecionarNo={selecionar}
              onFechar={() => setSelecionadoId(null)}
            />
          ) : (
            <p className="rounded-xl border border-dashed border-border bg-card p-5 text-sm text-muted-foreground">
              Selecione um nó para ver o que estudar, os pré-requisitos e quanto XP ele vale.
            </p>
          )}
        </aside>
      </div>

      {selecionado?.quiz && (
        <Dialog open={quizAberto} onOpenChange={setQuizAberto}>
          <DialogContent className="max-h-[90vh] max-w-xl overflow-y-auto">
            <DialogHeader>
              <DialogTitle>{selecionado.titulo}</DialogTitle>
              <DialogDescription>
                Acerte pelo menos 70% para dominar o nó. Você pode tentar quantas vezes precisar.
              </DialogDescription>
            </DialogHeader>
            <ChefeQuiz
              key={selecionado.id}
              questoes={selecionado.quiz.questoes}
              onEnviar={enviarQuiz}
              onFechar={() => setQuizAberto(false)}
            />
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
}
