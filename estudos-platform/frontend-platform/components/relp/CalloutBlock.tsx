import * as React from "react";
import { AlertTriangle, Info, Lightbulb } from "lucide-react";

import type { CalloutVariant } from "@/lib/types";
import { cn } from "@/lib/utils";

const ESTILO: Record<CalloutVariant, { box: string; icon: React.ElementType; titulo: string }> = {
  info: { box: "border-info/40 bg-info-muted text-info-fg", icon: Info, titulo: "Nota" },
  warning: { box: "border-warning/40 bg-warning-muted text-warning-fg", icon: AlertTriangle, titulo: "Atenção" },
  // "Dica" usa o verde-água de aprendizado, não o azul de marca.
  tip: { box: "border-learn/40 bg-learn-muted text-learn", icon: Lightbulb, titulo: "Dica" },
};

export type CalloutBlockProps = {
  variant?: CalloutVariant;
  title?: string;
  children: React.ReactNode;
  className?: string;
};

/** Destaque dentro do artigo (nota / atenção / dica). Conteúdo estático: sem role de alerta. */
export function CalloutBlock({ variant = "info", title, children, className }: CalloutBlockProps) {
  const { box, icon: Icon, titulo } = ESTILO[variant] ?? ESTILO.info;
  return (
    <aside className={cn("flex gap-3 rounded-lg border p-4 text-sm", box, className)} aria-label={title ?? titulo}>
      <Icon className="mt-0.5 h-4 w-4 shrink-0" aria-hidden="true" />
      <div className="min-w-0 space-y-1">
        <p className="font-semibold">{title ?? titulo}</p>
        <div className="leading-relaxed">{children}</div>
      </div>
    </aside>
  );
}
