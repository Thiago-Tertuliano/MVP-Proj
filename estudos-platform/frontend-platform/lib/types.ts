/**
 * Tipos do contrato com a API Go (`/api/v1`). Os nomes seguem o JSON do backend (snake_case).
 * Modelos de tela (roadmap etc.) ficam em `lib/roadmap.ts`.
 */

/* ------------------------------------------------------------------ conteúdo */

export type Block =
  | { type: "h"; level: number; text: string }
  | { type: "p"; text: string }
  | { type: "code"; lang?: string; text: string }
  | { type: "list"; items: string[]; ordered?: boolean }
  | { type: "link"; url: string; text?: string }
  | { type: "callout"; variant?: CalloutVariant; title?: string; text: string };

export type CalloutVariant = "info" | "warning" | "tip";

export type QuizOpcao = { id: string; texto: string };

export type QuizQuestao = {
  id: string;
  enunciado: string;
  opcoes: QuizOpcao[];
  correta: string;
  explicacao: string;
};

export type ArtigoMetadados = {
  tempo_leitura_min?: number;
  objetivo?: string;
  tags?: string[];
  origem?: string;
  fontes?: { titulo: string; url: string }[];
  quiz?: { questoes: QuizQuestao[] };
};

export type ArtigoStatus = "rascunho" | "revisao" | "publicado" | "arquivado";

export type Artigo = {
  id: string;
  slug: string;
  titulo: string;
  subtitulo?: string;
  capa_url?: string;
  trilha_id: string | null;
  modulo_id: string | null;
  conteudo: { blocks?: Block[] };
  metadados: ArtigoMetadados;
  autor_id: string;
  status: ArtigoStatus;
  publicado_em?: number;
  created_at: number;
  updated_at: number;
};

export type Modulo = {
  id: string;
  slug: string;
  titulo: string;
  descricao?: string;
  ordem: number;
};

export type Trilha = {
  id: string;
  slug: string;
  titulo: string;
  descricao?: string;
  capa_url?: string;
  ordem: number;
  publicada: boolean;
  modulos: Modulo[];
  created_at: number;
  updated_at: number;
};

export type ListaTrilhas = { itens: Trilha[] };
export type ListaArtigos = { itens: Artigo[] };

/* ------------------------------------------------------------------ aluno */

export type Usuario = {
  id: string;
  nome: string;
  email: string;
  /** Só `GET /auth/me` informa o papel; login/registro omitem. `editor` pode criar roadmaps. */
  papel?: "aluno" | "editor";
};

export type AuthResposta = {
  tokens: { access_token: string; refresh_token: string; expiracao_em: number };
  usuario: Usuario;
};

export type ProgressoTrilha = {
  trilha_id: string;
  concluidos: number;
  total: number;
  /** 0–100, pode vir com muitas casas decimais. */
  percentual: number;
  /** IDs dos artigos lidos — pinta os nós do mapa. */
  artigos_concluidos: string[];
};

export type ProgressoArtigo = { artigo_id: string; concluido: boolean };

export type ContinuarItem = {
  artigo_id: string;
  artigo_slug: string;
  artigo_titulo: string;
  trilha_id: string | null;
  trilha_slug: string | null;
  trilha_titulo: string | null;
  concluido: boolean;
  atualizado_em: number;
};

export type Continuar = { item: ContinuarItem | null };

export type Anotacao = {
  artigo_id: string;
  /** JSON livre; o front guarda `{ texto }`. */
  conteudo: { texto?: string } | null;
};

export type ResultadoBusca = { slug: string; titulo: string; similarity: number };
export type Busca = { itens: ResultadoBusca[] };
