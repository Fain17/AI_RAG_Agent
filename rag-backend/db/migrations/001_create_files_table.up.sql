-- 001_create_files_table.sql
CREATE TABLE IF NOT EXISTS files (
    id SERIAL PRIMARY KEY,
    filename TEXT NOT NULL,
    content TEXT NOT NULL,
    embedding VECTOR(1536) -- Adjust to match your embedding size
);