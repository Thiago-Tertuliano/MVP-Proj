-- Totais de gamificação por usuário (derivados de xp_eventos + atividade diária).
CREATE TABLE IF NOT EXISTS gamificacao (
    usuario_id UUID PRIMARY KEY REFERENCES usuarios(id) ON DELETE CASCADE,
    xp_total INT NOT NULL DEFAULT 0 CHECK (xp_total >= 0),
    streak_atual INT NOT NULL DEFAULT 0,
    streak_max INT NOT NULL DEFAULT 0,
    ultimo_dia DATE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Livro-razão de XP. A UNIQUE torna a concessão idempotente (sem XP duplicado).
CREATE TABLE IF NOT EXISTS xp_eventos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    usuario_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    origem TEXT NOT NULL,
    ref_id TEXT NOT NULL,
    xp INT NOT NULL CHECK (xp >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (usuario_id, origem, ref_id)
);

CREATE INDEX IF NOT EXISTS idx_xp_eventos_usuario ON xp_eventos (usuario_id, created_at DESC);

CREATE TABLE IF NOT EXISTS usuario_conquistas (
    usuario_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    codigo TEXT NOT NULL,
    conquistada_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (usuario_id, codigo)
);
