import { Badge } from "@/components/ui/badge";
import type { ArtigoStatus } from "@/lib/types";

const MAPA: Record<ArtigoStatus, { label: string; variant: "rascunho" | "revisao" | "publicado" | "arquivado" }> = {
  rascunho: { label: "Rascunho", variant: "rascunho" },
  revisao: { label: "Em revisão", variant: "revisao" },
  publicado: { label: "Publicado", variant: "publicado" },
  arquivado: { label: "Arquivado", variant: "arquivado" },
};

/** Status de conteúdo (rascunho / revisão / publicado / arquivado) com os tokens `status-*`. */
export function StatusBadge({ status, className }: { status: ArtigoStatus; className?: string }) {
  const { label, variant } = MAPA[status] ?? MAPA.rascunho;
  return (
    <Badge variant={variant} className={className}>
      {label}
    </Badge>
  );
}
