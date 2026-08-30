
CREATE TABLE IF NOT EXISTS anotacoes (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      usuario_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
      artigo_id UUID NOT NULL REFERENCES artigos(id) ON DELETE CASCADE,
      conteudo JSONB NOT NULL DEFAULT '{}',
      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
      UNIQUE (usuario_id, artigo_id)
);

CREATE INDEX IF NOT EXISTS idx_anotacoes_artigo ON anotacoes(artigo_id);