# routers/users.py
from typing import List

from fastapi import APIRouter
from fastapi import Depends
from fastapi import HTTPException
from fastapi import status
from fastapi.security import OAuth2PasswordBearer
from jose import ExpiredSignatureError
from jose import jwt
from jose import JWTError
from sqlalchemy.orm import Session
from src.config.settings import ALGORITHM
from src.config.settings import SECRET_KEY
from src.dependencies.database import get_db
from src.models.user import User
from src.repository import user as user_repo
from src.schemas.user import Token
from src.schemas.user import UserCreate
from src.schemas.user import UserInfo
from src.schemas.user import UserList
from src.schemas.user import UserLogin

router = APIRouter(prefix="/users", tags=["users"])
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")


def create_access_token(data: dict):
    return jwt.encode(data, SECRET_KEY, algorithm=ALGORITHM)


async def get_current_user(token: str = Depends(oauth2_scheme)):
    return decode_jwt_token(token)


# Create a function to decode a JWT token
def decode_jwt_token(token: str) -> UserInfo:
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=["HS256"])
        user_info: UserInfo = UserInfo(email=payload["email"], name=payload["name"])
        return user_info

    except ExpiredSignatureError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Token expired",
            headers={"WWW-Authenticate": "Bearer"},
        )
    except JWTError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid token",
            headers={"WWW-Authenticate": "Bearer"},
        )


# API endpoints
@router.post("/", response_model=Token)
def create_user_endpoint(user: UserCreate, db: Session = Depends(get_db)):
    db_user = user_repo.create_user(db, user)
    access_token = create_access_token(
        data={"email": db_user.email, "name": db_user.name}
    )
    return {"access_token": access_token, "token_type": "bearer"}


@router.get("/", response_model=UserList)
def list_user_endpoint(
    current_user: UserInfo = Depends(get_current_user), db: Session = Depends(get_db)
):
    db_user = user_repo.get_user(db, current_user.email)
    if db_user is None:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="User not found",
            headers={"WWW-Authenticate": "Bearer"},
        )
    if db_user.role.desc != "admin":
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="User not authorized",
            headers={"WWW-Authenticate": "Bearer"},
        )
    users: List[User] = user_repo.all(db)
    return UserList(
        users=[UserInfo(email=str(u.email), name=str(u.name)) for u in users]
    )


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


@router.get("/me", response_model=UserInfo)
def me(current_user: UserInfo = Depends(get_current_user)) -> UserInfo:
    return current_user
