import * as React from "react";
import { Check, Flag, Lock, Play, Swords, type LucideIcon } from "lucide-react";

import { descricaoEstado, ESTADO_ROTULO, LAYOUT, TIPO_ROTULO } from "@/lib/roadmaps/graph";
import type { EstadoNo, TipoNo } from "@/lib/roadmaps/types";
import { cn } from "@/lib/utils";

const ICONE_ESTADO: Record<EstadoNo, LucideIcon> = {
  bloqueado: Lock,
  disponivel: Play,
  em_curso: Play,
  dominado: Check,
};

const ICONE_TIPO: Partial<Record<TipoNo, LucideIcon>> = { chefe: Swords, marco: Flag };

export type RoadmapNodeCardProps = {
  titulo: string;
  tipo: TipoNo;
  /** Visitante não tem estados: o card aparece "neutro" (como disponível, sem rótulo de estado). */
  estado?: EstadoNo | null;
  xp: number;
  selecionado?: boolean;
  onSelecionar?: () => void;
  /** Largura fixa no canvas; `auto` na lista/preview. */
  fixo?: boolean;
  className?: string;
};

/**
 * Cartão de um nó do roadmap. Puro: cor, ícone e texto vêm do estado. O estado nunca é só cor —
 * há ícone + rótulo (`sr-only`) para daltonismo e leitor de tela.
 */
export function RoadmapNodeCard({
  titulo,
  tipo,
  estado,
  xp,
  selecionado,
  onSelecionar,
  fixo = true,
  className,
}: RoadmapNodeCardProps) {
  const efetivo: EstadoNo = estado ?? "disponivel";
  const IconeBase = estado ? ICONE_ESTADO[efetivo] : (ICONE_TIPO[tipo] ?? Play);
  // Chefe/marco mantêm o ícone do tipo enquanto não dominados (identidade visual do "chefe").
  const Icone = efetivo !== "dominado" && efetivo !== "bloqueado" && ICONE_TIPO[tipo] ? ICONE_TIPO[tipo]! : IconeBase;
  const chefe = tipo === "chefe";

  return (
    <button
      type="button"
      onClick={onSelecionar}
      aria-pressed={selecionado}
      aria-label={`${titulo}. ${TIPO_ROTULO[tipo]}. ${estado ? descricaoEstado(estado, tipo) : "Roadmap"}. ${xp} XP.`}
      data-estado={estado ?? "visitante"}
      data-tipo={tipo}
      style={fixo ? { width: LAYOUT.larguraNo, minHeight: LAYOUT.alturaNo } : undefined}
      className={cn(
        "group flex items-center gap-3 rounded-xl border-2 bg-card p-3 text-left shadow-sm transition-[box-shadow,transform,border-color] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background motion-safe:hover:-translate-y-0.5",
        !fixo && "w-full",
        chefe && "border-[3px]",
        efetivo === "bloqueado" && "border-dashed border-border bg-muted text-muted-foreground",
        efetivo === "disponivel" && "border-primary/60 hover:border-primary",
        efetivo === "em_curso" && "border-primary bg-primary-muted ring-4 ring-primary/20",
        efetivo === "dominado" && "border-done-fg bg-done-muted",
        selecionado && "ring-4 ring-ring/40",
        className,
      )}
    >
      <span
        aria-hidden="true"
        className={cn(
          "flex h-10 w-10 shrink-0 items-center justify-center rounded-full border-2",
          efetivo === "bloqueado" && "border-border bg-card text-muted-foreground",
          efetivo === "disponivel" && "border-primary bg-card text-primary",
          efetivo === "em_curso" && "border-primary bg-primary text-primary-foreground",
          efetivo === "dominado" && "border-done-fg bg-done-fg text-primary-foreground",
        )}
      >
        <Icone className="h-5 w-5" />
      </span>

      <span className="min-w-0 flex-1 space-y-1">
        <span
          className={cn(
            "line-clamp-2 block text-sm font-semibold leading-snug",
            efetivo === "bloqueado" ? "text-muted-foreground" : "text-foreground",
          )}
        >
          {titulo}
        </span>
        <span className="flex flex-wrap items-center gap-x-2 gap-y-0.5 text-[11px] font-medium uppercase tracking-wide text-muted-foreground">
          <span>{TIPO_ROTULO[tipo]}</span>
          {xp > 0 && <span className="text-notify">+{xp} XP</span>}
        </span>
        {estado && <span className="sr-only">{ESTADO_ROTULO[estado]}</span>}
      </span>
    </button>
  );
}
