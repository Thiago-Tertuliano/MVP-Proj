export type Block =
  | { type: "h"; level: number; text: string }
  | { type: "p"; text: string }
  | { type: "code"; lang?: string; text: string };

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

export type Artigo = {
  id: string;
  slug: string;
  titulo: string;
  subtitulo?: string;
  status: "publicado" | "rascunho";
  trilhaSlug?: string;
  moduloSlug?: string;
  conteudo: { blocks: Block[] };
  metadados: ArtigoMetadados;
};

export type Modulo = {
  slug: string;
  titulo: string;
  descricao: string;
  artigos: { slug: string; titulo: string; concluido: boolean }[];
};

export type Trilha = {
  slug: string;
  titulo: string;
  descricao: string;
  publicada: boolean;
  progressoPct: number;
  modulos: Modulo[];
};
