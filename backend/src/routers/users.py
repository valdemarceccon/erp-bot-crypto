# routers/users.py
from fastapi import APIRouter
from fastapi import Depends
from fastapi import HTTPException
from fastapi import status
from jose import jwt
from sqlalchemy.orm import Session
from src.config.settings import ALGORITHM
from src.config.settings import SECRET_KEY
from src.dependencies.database import get_db
from src.repository import user as user_repo
from src.schemas.user import Token
from src.schemas.user import UserCreate
from src.schemas.user import UserLogin

router = APIRouter()


def create_access_token(data: dict):
    return jwt.encode(data, SECRET_KEY, algorithm=ALGORITHM)


# API endpoints
@router.post("/users", response_model=Token)
def create_user_endpoint(user: UserCreate, db: Session = Depends(get_db)):
    db_user = user_repo.create_user(db, user)
    access_token = create_access_token(data={"sub": db_user.email})
    return {"access_token": access_token, "token_type": "bearer"}


@router.post("/token", response_model=Token)
def login(user: UserLogin, db: Session = Depends(get_db)):
    authenticated_user = user_repo.authenticate_user(db, user.email, user.password)
    if not authenticated_user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect email or password",
            headers={"WWW-Authenticate": "Bearer"},
        )
    access_token = create_access_token(data={"sub": user.email})
    return {"access_token": access_token, "token_type": "bearer"}
