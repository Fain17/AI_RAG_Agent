from sqlalchemy import text, bindparam
from pgvector.sqlalchemy import Vector
from app.db.session import SessionLocal  # adjust import as needed
from app.services.embedding import embed_prompt


async def find_similar_files(prompt: str, top_k: int = 5):
    embedding = await embed_prompt(prompt)

    with SessionLocal() as session:
        query = text("""
            SELECT filename, content, embedding <=> :embedding AS similarity
            FROM files
            ORDER BY embedding <=> :embedding
            LIMIT :top_k;
        """).bindparams(
            bindparam("embedding", type_=Vector),
            bindparam("top_k", type_=int)
        )

        result = session.execute(query, {
            "embedding": embedding,
            "top_k": top_k
        })

        return result.fetchall()
