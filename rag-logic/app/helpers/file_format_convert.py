from pathlib import Path
from typing import cast

from fastapi import UploadFile
from pdfminer.high_level import extract_text


async def to_text(file: UploadFile) -> str:
    if file.filename is None:
        raise ValueError("File must have a filename")

    ext = Path(file.filename).suffix.lower()
    data = await file.read()

    if ext == ".pdf":
        tmp = f"/tmp/{file.filename}"
        with open(tmp, "wb") as f:
            f.write(data)
        return cast(str, extract_text(tmp))

    return cast(str, data.decode("utf-8"))
