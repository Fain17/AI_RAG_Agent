from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import Column, String, Text
from pgvector.sqlalchemy import Vector

Base = declarative_base()

class File(Base):
    __tablename__ = "files"
    
    filename = Column(String, primary_key=True)
    content = Column(Text)
    embedding = Column(Vector(384))
