from fastapi import APIRouter, UploadFile, File, HTTPException, Depends
from app.services.convert import to_text
from app.services.embedding import embed_text
from pathlib import Path
from app.models.schemas import QueryRequest, QueryResponse
from app.services.query_service import run_query_pipeline
import httpx
from sqlalchemy.orm import Session
from app.db.session import get_db_session

router = APIRouter()
ALLOWED = {'.txt', '.md', '.pdf'}

@router.post("/upload")
async def upload(file: UploadFile = File(...)):
    
    # if not file:
    #     raise ValueError("No file value present")
    
    ext = Path(file.filename).suffix.lower()
    if ext not in ALLOWED:
        raise HTTPException(400, f"Unsupported file type: {ext}")

    text = await to_text(file)
    embedding = await embed_text(text)
    payload = {
        "filename": file.filename,
        "content": text,
        "embedding": embedding
    }

    async with httpx.AsyncClient() as client:
        resp = await client.post("http://localhost:8080/files", json=payload)
        resp.raise_for_status()
        return resp.json()

@router.post("/query", response_model=QueryResponse)
async def query_route(req: QueryRequest, db: Session = Depends(get_db_session)):
    files, answer = await run_query_pipeline(req.prompt, db)
    return QueryResponse(matches=files, answer=answer)