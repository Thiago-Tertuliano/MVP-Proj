DROP INDEX IF EXISTS idx_artigos_embedding_hnsw;
ALTER TABLE artigos DROP COLUMN IF EXISTS embedding;
