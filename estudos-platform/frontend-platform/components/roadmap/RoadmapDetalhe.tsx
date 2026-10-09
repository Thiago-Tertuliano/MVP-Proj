import * as React from "react";
import Link from "next/link";
import { BookOpen, Lock, Swords, Trophy, X } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { descricaoEstado, ESTADO_ROTULO, TIPO_ROTULO } from "@/lib/roadmaps/graph";
import type { EstadoNo, RoadmapNo } from "@/lib/roadmaps/types";
import { cn } from "@/lib/utils";

const BADGE_ESTADO: Record<EstadoNo, "default" | "done" | "info" | "outline"> = {
  bloqueado: "outline",
  disponivel: "info",
  em_curso: "default",
  dominado: "done",
};

export type RoadmapDetalheProps = {
  no: RoadmapNo;
  /** `null` = visitante (sem estados). */
  estado: EstadoNo | null;
  /** Pré-requisitos `requer` ainda não dominados. */
  pendentes?: RoadmapNo[];
  /** Caminho para entrar e voltar a este roadmap. */
  loginHref?: string;
  concluindo?: boolean;
  onConcluir?: () => void;
  onEnfrentarChefe?: () => void;
  onSelecionarNo?: (id: string) => void;
  onFechar?: () => void;
  className?: string;
};

/** Painel de detalhe de um nó: explica o estado e oferece a ação certa para o modo de conclusão. */
export function RoadmapDetalhe({
  no,
  estado,
  pendentes = [],
  loginHref,
  concluindo,
  onConcluir,
  onEnfrentarChefe,
  onSelecionarNo,
  onFechar,
  className,
}: RoadmapDetalheProps) {
  const chefe = no.tipo === "chefe";
  const bloqueado = estado === "bloqueado";
  const dominado = estado === "dominado";
  const podeAgir = estado === "disponivel" || estado === "em_curso";
  const artigoHref = no.artigo_slug ? `/artigos/${no.artigo_slug}` : null;

  return (
    <section
      aria-labelledby="detalhe-no-titulo"
      className={cn("space-y-4 rounded-xl border border-border bg-card p-5", chefe && "border-notify/50", className)}
    >
      <header className="flex items-start gap-3">
        <div className="min-w-0 flex-1 space-y-2">
          <div className="flex flex-wrap items-center gap-2">
            <Badge variant={chefe ? "notify" : "outline"}>
              {chefe && <Swords className="mr-1 h-3 w-3" aria-hidden="true" />}
              {TIPO_ROTULO[no.tipo]}
            </Badge>
            {estado && <Badge variant={BADGE_ESTADO[estado]}>{ESTADO_ROTULO[estado]}</Badge>}
            {no.xp > 0 && <span className="text-xs font-semibold text-notify">+{no.xp} XP</span>}
          </div>
          <h2 id="detalhe-no-titulo" className="text-xl font-semibold leading-tight tracking-tight">
            {no.titulo}
          </h2>
        </div>
        {onFechar && (
          <Button type="button" variant="ghost" size="icon" onClick={onFechar} aria-label="Fechar detalhes">
            <X />
          </Button>
        )}
      </header>

      {no.descricao && <p className="text-sm leading-relaxed text-muted-foreground">{no.descricao}</p>}

      {estado && <p className="sr-only">{descricaoEstado(estado, no.tipo)}</p>}

      {bloqueado && (
        <div className="rounded-md border border-dashed border-border bg-muted p-3 text-sm">
          <p className="flex items-center gap-2 font-medium text-foreground">
            <Lock className="h-4 w-4" aria-hidden="true" />
            Conclua antes:
          </p>
          <ul className="mt-2 space-y-1">
            {pendentes.map((p) => (
              <li key={p.id}>
                <button
                  type="button"
                  onClick={() => onSelecionarNo?.(p.id)}
                  className="rounded-sm text-left text-primary underline-offset-4 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                >
                  {p.titulo}
                </button>
              </li>
            ))}
            {pendentes.length === 0 && <li className="text-muted-foreground">Pré-requisitos pendentes.</li>}
          </ul>
        </div>
      )}

      {estado === null && (
        <p className="rounded-md bg-info-muted p-3 text-sm text-info-fg">
          Para registrar seu progresso e ganhar XP,{" "}
          <Link href={loginHref ?? "/login"} className="font-medium underline underline-offset-4">
            entre na sua conta
          </Link>
          .
        </p>
      )}

      {dominado && (
        <p className="flex items-center gap-2 text-sm font-medium text-done-fg">
          <Trophy className="h-4 w-4" aria-hidden="true" />
          {chefe ? "Chefe derrotado!" : "Nó dominado."}
        </p>
      )}

      <div className="flex flex-wrap items-center gap-2">
        {podeAgir && no.conclusao === "manual" && (
          <Button type="button" variant="success" onClick={onConcluir} loading={concluindo}>
            Marcar como dominado
          </Button>
        )}
        {podeAgir && no.conclusao === "quiz" && (
          <Button type="button" onClick={onEnfrentarChefe}>
            <Swords />
            {chefe ? "Enfrentar chefe" : "Fazer o quiz"}
          </Button>
        )}
        {podeAgir && no.conclusao === "artigo" &&
          (artigoHref ? (
            <Button asChild>
              <Link href={artigoHref}>
                <BookOpen />
                Ler o artigo para dominar
              </Link>
            </Button>
          ) : (
            <p className="text-sm text-muted-foreground">O artigo deste nó não está disponível no momento.</p>
          ))}

        {artigoHref && !(podeAgir && no.conclusao === "artigo") && !bloqueado && (
          <Button asChild variant="outline">
            <Link href={artigoHref}>
              <BookOpen />
              {no.artigo_titulo ? `Ler: ${no.artigo_titulo}` : "Ler o artigo"}
            </Link>
          </Button>
        )}
      </div>

      {podeAgir && no.conclusao === "artigo" && (
        <p className="text-xs text-muted-foreground">
          O nó é dominado automaticamente quando você marca o artigo como lido.
        </p>
      )}
    </section>
  );
}
