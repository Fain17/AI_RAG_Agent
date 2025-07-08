from typing import List

from pydantic import BaseModel


class QueryRequest(BaseModel):
    prompt: str


class FileData(BaseModel):
    filename: str
    content: str
    similarity: float


class QueryResponse(BaseModel):
    matches: List[FileData]
    answer: str
