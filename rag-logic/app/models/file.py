from pydantic import BaseModel
from typing import List

class FileUploadResponse(BaseModel):
    id: int
    filename: str
    content: str
    embedding: List[float]
