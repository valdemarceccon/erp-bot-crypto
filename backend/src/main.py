from typing import Optional

from fastapi import Depends
from fastapi import FastAPI
from fastapi import HTTPException
from fastapi import status
from fastapi.security import OAuth2PasswordBearer
from jose import jwt
from passlib.hash import bcrypt
from pydantic import BaseModel
from sqlalchemy import Column
from sqlalchemy import create_engine
from sqlalchemy import Integer
from sqlalchemy import String
from sqlalchemy import text
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import Session
from sqlalchemy.orm import sessionmaker
from src.repository.user import authenticate_user
from src.repository.user import create_user
from src.schemas.user import Token
from src.schemas.user import UserCreate
from src.schemas.user import UserLogin

from .models.base import engine
from .models.base import SessionLocal

# import httpx

app = FastAPI()


@app.get("/")
def read_root():
    return {"Hello": "World"}


app = FastAPI()


# EXTERNAL_API_URL = "https://example.com/api/health"


@app.get("/health")
async def health_check():
    health_status = {
        "database": "unknown",
        # "external_api": "unknown",
    }

    # Check the database connection
    try:
        with engine.connect() as connection:
            connection.execute(text("SELECT 1"))
        health_status["database"] = "OK"
    except Exception as e:
        health_status["database"] = "unavailable"

    # Check the external API
    # try:
    #     async with httpx.AsyncClient() as client:
    #         response = await client.get(EXTERNAL_API_URL)
    #         response.raise_for_status()
    #     health_status["external_api"] = "OK"
    # except Exception as e:
    #     health_status["external_api"] = "unavailable"

    if all(value == "OK" for value in health_status.values()):
        return health_status
    else:
        raise HTTPException(status_code=500, detail=health_status)


# OAuth2 setup
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")


# Utility functions
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


# JWT utility functions
SECRET_KEY = "your-secret-key"
ALGORITHM = "HS256"


def create_access_token(data: dict):
    return jwt.encode(data, SECRET_KEY, algorithm=ALGORITHM)


# API endpoints
@app.post("/users", response_model=Token)
def create_user_endpoint(user: UserCreate, db: Session = Depends(get_db)):
    db_user = create_user(db, user)
    access_token = create_access_token(data={"sub": db_user.email})
    return {"access_token": access_token, "token_type": "bearer"}


@app.post("/token", response_model=Token)
def login(user: UserLogin, db: Session = Depends(get_db)):
    authenticated_user = authenticate_user(db, user.email, user.password)
    if not authenticated_user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect email or password",
            headers={"WWW-Authenticate": "Bearer"},
        )
    access_token = create_access_token(data={"sub": user.email})
    return {"access_token": access_token, "token_type": "bearer"}


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
