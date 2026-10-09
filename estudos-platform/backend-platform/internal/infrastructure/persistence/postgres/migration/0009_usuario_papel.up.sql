-- Papel do usuário: 'aluno' (padrão) ou 'editor' (autoria de roadmaps).
ALTER TABLE usuarios
    ADD COLUMN IF NOT EXISTS papel TEXT NOT NULL DEFAULT 'aluno';

ALTER TABLE usuarios
    DROP CONSTRAINT IF EXISTS usuarios_papel_check;

ALTER TABLE usuarios
    ADD CONSTRAINT usuarios_papel_check CHECK (papel IN ('aluno', 'editor'));
