DROP INDEX IF EXISTS idx_artigos_modulo;
DROP INDEX IF EXISTS idx_artigos_trilha;
ALTER TABLE artigos DROP COLUMN IF EXISTS modulo_id;
ALTER TABLE artigos DROP COLUMN IF EXISTS trilha_id;
