from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, Session
from typing import Generator, Annotated
from fastapi import Depends


from contextlib import contextmanager
# Replace with your actual connection string
DATABASE_URL = "postgresql+psycopg2://postgres:test123@localhost:5432/ragDB"

engine = create_engine(
    DATABASE_URL,
    pool_pre_ping=True
)

SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


# Dependency for FastAPI routes/services
def get_db_session():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
        
