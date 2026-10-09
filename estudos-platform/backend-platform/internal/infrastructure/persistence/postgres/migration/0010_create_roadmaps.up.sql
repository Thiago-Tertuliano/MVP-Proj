CREATE TABLE IF NOT EXISTS roadmaps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT NOT NULL UNIQUE,
    titulo TEXT NOT NULL,
    descricao TEXT,
    icone TEXT,
    publicado BOOLEAN NOT NULL DEFAULT false,
    autor_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_roadmaps_publicado ON roadmaps (publicado, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_roadmaps_autor ON roadmaps (autor_id);

-- Nós do grafo. O gabarito do quiz (coluna quiz) nunca é enviado ao aluno.
CREATE TABLE IF NOT EXISTS roadmap_nos (
    id UUID PRIMARY KEY,
    roadmap_id UUID NOT NULL REFERENCES roadmaps(id) ON DELETE CASCADE,
    tipo TEXT NOT NULL CHECK (tipo IN ('topico', 'chefe', 'marco')),
    titulo TEXT NOT NULL,
    descricao TEXT,
    artigo_id UUID REFERENCES artigos(id) ON DELETE SET NULL,
    pos_x DOUBLE PRECISION NOT NULL DEFAULT 0,
    pos_y DOUBLE PRECISION NOT NULL DEFAULT 0,
    xp INT NOT NULL DEFAULT 10 CHECK (xp >= 0 AND xp <= 1000),
    conclusao TEXT NOT NULL DEFAULT 'manual' CHECK (conclusao IN ('manual', 'artigo', 'quiz')),
    quiz JSONB
);

CREATE INDEX IF NOT EXISTS idx_roadmap_nos_roadmap ON roadmap_nos (roadmap_id);
CREATE INDEX IF NOT EXISTS idx_roadmap_nos_artigo ON roadmap_nos (artigo_id) WHERE artigo_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS roadmap_arestas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    roadmap_id UUID NOT NULL REFERENCES roadmaps(id) ON DELETE CASCADE,
    origem_id UUID NOT NULL REFERENCES roadmap_nos(id) ON DELETE CASCADE,
    destino_id UUID NOT NULL REFERENCES roadmap_nos(id) ON DELETE CASCADE,
    tipo TEXT NOT NULL DEFAULT 'requer' CHECK (tipo IN ('requer', 'opcional')),
    CHECK (origem_id <> destino_id),
    UNIQUE (origem_id, destino_id)
);

CREATE INDEX IF NOT EXISTS idx_roadmap_arestas_roadmap ON roadmap_arestas (roadmap_id);

-- Progresso do aluno: um registro por nó dominado.
CREATE TABLE IF NOT EXISTS roadmap_progresso (
    usuario_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    no_id UUID NOT NULL REFERENCES roadmap_nos(id) ON DELETE CASCADE,
    roadmap_id UUID NOT NULL REFERENCES roadmaps(id) ON DELETE CASCADE,
    concluido_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (usuario_id, no_id)
);

CREATE INDEX IF NOT EXISTS idx_roadmap_progresso_roadmap ON roadmap_progresso (usuario_id, roadmap_id);
