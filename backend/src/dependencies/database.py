import logging
import os

from fastapi.logger import logger
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

fallback_url = "sqlite:///fallback.db"

SQLALCHEMY_DATABASE_URL = os.getenv("POSTGRES_URL", fallback_url)

if SQLALCHEMY_DATABASE_URL == fallback_url:
    logger.info(f"Using {SQLALCHEMY_DATABASE_URL} as database")

engine = create_engine(SQLALCHEMY_DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


# Utility functions
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
