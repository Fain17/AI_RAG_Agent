CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    filename TEXT NOT NULL,
    content TEXT NOT NULL,
    embedding VECTOR(384) NOT NULL
);

CREATE INDEX idx_files_embedding ON files USING ivfflat (embedding vector_cosine_ops) WITH (lists= 100);

