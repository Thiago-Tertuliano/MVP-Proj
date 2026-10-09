import type { ResultadoGamificacao } from "@/lib/roadmaps/types";

export type Celebracao = {
  /** Mensagem principal do toast de XP (ausente quando só há conquista/level-up). */
  titulo: string | null;
  descricao?: string;
  conquistas: { codigo: string; nome: string; descricao: string }[];
  subiuDeNivel: number | null;
  roadmapsCompletos: string[];
};

/** Traduz o resultado de uma ação no que mostrar ao aluno (puro: separa texto de efeitos). */
export function celebracaoDe(r: ResultadoGamificacao): Celebracao {
  const partes: string[] = [];
  if (r.nos_concluidos.length > 1) partes.push(`${r.nos_concluidos.length} nós dominados`);
  if (r.nos_desbloqueados.length > 0) {
    partes.push(r.nos_desbloqueados.length === 1 ? "1 nó desbloqueado" : `${r.nos_desbloqueados.length} nós desbloqueados`);
  }
  if (r.streak >= 2) partes.push(`sequência de ${r.streak} dias`);

  return {
    titulo: r.xp_ganho > 0 ? `+${r.xp_ganho} XP` : null,
    descricao: partes.length ? partes.join(" · ") : undefined,
    conquistas: r.conquistas_novas.map((c) => ({ codigo: c.codigo, nome: c.nome, descricao: c.descricao })),
    subiuDeNivel: r.subiu_de_nivel ? r.nivel : null,
    roadmapsCompletos: r.roadmaps_completos,
  };
}
