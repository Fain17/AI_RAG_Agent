import httpx
from typing import cast

EMBEDDING_API_URL = "http://localhost:8001/embed"


async def embed_text(text: str) -> list[float]:
    async with httpx.AsyncClient() as client:
        response = await client.post(EMBEDDING_API_URL, json={"text": text})
        response.raise_for_status()
        return cast(list[float], response.json()["embedding"])


async def get_embedding(prompt: str) -> list[float]:
    async with httpx.AsyncClient() as client:
        response = await client.post(EMBEDDING_API_URL, json={"text": prompt})
        response.raise_for_status()
        return cast(list[float], response.json()["embedding"])
