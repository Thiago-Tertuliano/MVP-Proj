/**
 * Contrato dos Roadmaps gamificados com a API Go (`/api/v1/roadmaps`, `/gamificacao`).
 * Nomes seguem o JSON do backend (snake_case). Estados e XP são sempre decididos pelo servidor.
 */

export type TipoNo = "topico" | "chefe" | "marco";
export type ConclusaoNo = "manual" | "artigo" | "quiz";
export type TipoAresta = "requer" | "opcional";
export type EstadoNo = "bloqueado" | "disponivel" | "em_curso" | "dominado";

export type QuizOpcaoRoadmap = { id: string; texto: string };

/** `correta`/`explicacao` só existem na visão do editor; o aluno nunca recebe o gabarito. */
export type QuizQuestaoRoadmap = {
  id: string;
  enunciado: string;
  opcoes: QuizOpcaoRoadmap[];
  correta?: string;
  explicacao?: string;
};

export type QuizRoadmap = { questoes: QuizQuestaoRoadmap[] };

export type RoadmapNo = {
  id: string;
  tipo: TipoNo;
  titulo: string;
  descricao?: string;
  artigo_id?: string;
  artigo_slug?: string;
  artigo_titulo?: string;
  pos_x: number;
  pos_y: number;
  xp: number;
  conclusao: ConclusaoNo;
  quiz?: QuizRoadmap;
};

export type RoadmapAresta = { origem: string; destino: string; tipo: TipoAresta };

export type Roadmap = {
  id: string;
  slug: string;
  titulo: string;
  descricao: string;
  icone: string;
  publicado: boolean;
  nos: RoadmapNo[];
  arestas: RoadmapAresta[];
  xp_total: number;
  atualizado_em: number;
};

export type RoadmapResumo = {
  id: string;
  slug: string;
  titulo: string;
  descricao: string;
  icone: string;
  publicado: boolean;
  total_nos: number;
  total_chefes: number;
  xp_total: number;
  atualizado_em: number;
};

export type ListaRoadmaps = { itens: RoadmapResumo[] };

/* ------------------------------------------------------------------ progresso do aluno */

export type RoadmapProgresso = {
  slug: string;
  estados: Record<string, EstadoNo>;
  concluidos: number;
  total: number;
  percentual: number;
  xp_ganho: number;
  xp_total: number;
  completo: boolean;
};

export type RoadmapProgressoItem = { slug: string; concluidos: number; total: number; percentual: number };
export type ListaProgressoRoadmaps = { itens: RoadmapProgressoItem[] };

/* ------------------------------------------------------------------ gamificação */

export type Conquista = {
  codigo: string;
  nome: string;
  descricao: string;
  icone: string;
  conquistada: boolean;
  em?: number;
};

export type Gamificacao = {
  xp_total: number;
  nivel: number;
  xp_no_nivel: number;
  xp_para_proximo: number;
  streak_atual: number;
  streak_max: number;
  ativo_hoje: boolean;
  conquistas: Conquista[];
};

/** O que mudou após uma ação do aluno — alimenta toasts, barra de XP e level-up. */
export type ResultadoGamificacao = {
  xp_ganho: number;
  xp_total: number;
  nivel: number;
  nivel_anterior: number;
  subiu_de_nivel: boolean;
  streak: number;
  nos_concluidos: string[];
  nos_desbloqueados: string[];
  roadmaps_completos: string[];
  conquistas_novas: Conquista[];
};

export type ConclusaoNoResposta = {
  no_id: string;
  ja_concluido: boolean;
  resultado: ResultadoGamificacao;
  estados: Record<string, EstadoNo>;
  concluidos: number;
  total: number;
};

export type CorrecaoQuestao = { questao_id: string; correta: boolean; explicacao?: string };

export type QuizResultado = {
  aprovado: boolean;
  acertos: number;
  total: number;
  nota_minima_pct: number;
  /** Só vem quando aprovado (evita descobrir o gabarito por tentativa e erro). */
  correcao?: CorrecaoQuestao[];
  conclusao?: ConclusaoNoResposta;
};

/* ------------------------------------------------------------------ edição (editor) */

export type SalvarRoadmapPayload = {
  titulo: string;
  descricao: string;
  icone: string;
  nos: RoadmapNo[];
  arestas: RoadmapAresta[];
};
