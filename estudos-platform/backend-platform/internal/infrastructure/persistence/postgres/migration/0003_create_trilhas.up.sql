CREATE TABLE IF NOT EXISTS trilhas (
    id UUID PRIMARY KEY,
    slug VARCHAR(200) NOT NULL UNIQUE,
    titulo VARCHAR(200) NOT NULL,
    descricao TEXT,
    capa_url TEXT,
    ordem INT NOT NULL DEFAULT 0,
    publicada BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS modulos (
    id UUID PRIMARY KEY,
    trilha_id UUID NOT NULL REFERENCES trilhas(id) ON DELETE CASCADE,
    slug VARCHAR(200) NOT NULL,
    titulo VARCHAR(200) NOT NULL,
    descricao TEXT,
    ordem INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (trilha_id, slug)
);

CREATE INDEX IF NOT EXISTS idx_trilhas_publicada ON trilhas(publicada);
CREATE INDEX IF NOT EXISTS idx_modulos_trilha ON modulos(trilha_id);
