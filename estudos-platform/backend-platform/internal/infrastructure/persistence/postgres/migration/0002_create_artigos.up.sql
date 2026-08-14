CREATE TABLE IF NOT EXISTS artigos (
    id UUID PRIMARY KEY,
    slug VARCHAR(200) NOT NULL UNIQUE,
    titulo VARCHAR(300) NOT NULL,
    subtitulo TEXT,
    capa_url TEXT,
    conteudo JSONB NOT NULL DEFAULT '{}',
    metadados JSONB NOT NULL DEFAULT '{}',
    autor_id UUID NOT NULL REFERENCES usuarios(id),
    status VARCHAR(20) NOT NULL DEFAULT 'rascunho',
    publicado_em TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_artigos_status ON artigos(status);
CREATE INDEX IF NOT EXISTS idx_artigos_autor ON artigos(autor_id);
CREATE INDEX IF NOT EXISTS idx_artigos_conteudo_gin ON artigos USING GIN (conteudo);
CREATE INDEX IF NOT EXISTS idx_artigos_metadados_gin ON artigos USING GIN (metadados);
