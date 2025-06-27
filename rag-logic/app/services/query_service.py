from app.services.embedding import get_embedding
from app.services.retriever import fetch_similar_files_pgvector
from app.services.llm_chain import chain
from app.models.schemas import FileData
from sqlalchemy.orm import Session

async def run_query_pipeline(prompt: str, db: Session) -> tuple[list[FileData], str]:
    embedding = await get_embedding(prompt)
    files = await fetch_similar_files_pgvector(embedding, db)

    context = "\n\n".join(f"{f.filename}:\n{f.content}" for f in files)
    answer = chain.invoke({"context": context, "question": prompt})

    return files, answer