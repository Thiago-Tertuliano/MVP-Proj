CREATE TABLE IF NOT EXISTS progresso_estudo (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    usuario_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    artigo_id UUID NOT NULL REFERENCES artigos(id) ON DELETE CASCADE,
    trilha_id UUID REFERENCES trilhas(id) ON DELETE CASCADE,
    concluido BOOLEAN NOT NULL DEFAULT false,
    percentual INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (usuario_id, artigo_id)
);