from typing import List

from pydantic import BaseModel


class FileUploadResponse(BaseModel):
    id: int
    filename: str
    content: str
    embedding: List[float]
