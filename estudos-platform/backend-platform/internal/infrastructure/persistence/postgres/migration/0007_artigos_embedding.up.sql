CREATE EXTENSION IF NOT EXISTS vector;

ALTER TABLE artigos ADD COLUMN IF NOT EXISTS embedding VECTOR(1536);

CREATE INDEX IF NOT EXISTS idx_artigos_embedding_hnsw
    ON artigos USING hnsw (embedding vector_cosine_ops);
