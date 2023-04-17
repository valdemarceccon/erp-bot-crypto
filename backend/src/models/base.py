import logging
import os

from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker

logger = logging.getLogger(__name__)

fallback_url = "sqlite:///fallback.db"

SQLALCHEMY_DATABASE_URL = os.getenv("POSTGRES_URL", fallback_url)

if SQLALCHEMY_DATABASE_URL == fallback_url:
    logger.info(f"Using {SQLALCHEMY_DATABASE_URL} as database")

engine = create_engine(SQLALCHEMY_DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base = declarative_base()
