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
from src.models.roles import PermissionEnum
from src.models.user import User
from src.repository import user as user_repo
from src.repository.user import user_has_permission
from src.schemas.user import Token
from src.schemas.user import UserCreate
from src.schemas.user import UserInfo
from src.schemas.user import UserList
from src.schemas.user import UserLogin
from src.schemas.user import UserUpdate


oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")


def get_current_user(token: str = Depends(oauth2_scheme)) -> int:
    return decode_jwt_token(token)


def has_permission(permission_name: PermissionEnum):
    def inner(
        session: Session = Depends(get_db),
        user_id: int = Depends(get_current_user),
    ) -> int:
        if not user_has_permission(session, user_id, permission_name.value):
            raise HTTPException(status_code=403, detail="Permission denied")
        return user_id

    return inner


router = APIRouter(prefix="/users", tags=["users"])

import datetime


def create_access_token(data: dict):
    data["timestamp"] = datetime.datetime.now().isoformat()
    return jwt.encode(data, SECRET_KEY, algorithm=ALGORITHM)


# Create a function to decode a JWT token
def decode_jwt_token(token: str) -> int:
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=["HS256"])
        if "id" not in token:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Invalid token",
                headers={"WWW-Authenticate": "Bearer"},
            )
        id: int = int(payload["id"])
        return id

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
@router.patch("/", response_model=UserInfo)
def update_user_endpoint(
    user: UserUpdate,
    db: Session = Depends(get_db),
    user_id: int = Depends(get_current_user),
):
    db_user = user_repo.update_user(db, user_id, user)
    if db_user is None:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid token",
            headers={"WWW-Authenticate": "Bearer"},
        )

    return UserInfo(email=db_user.email, name=db_user.name, username=db_user.username)


# API endpoints
@router.post("/", response_model=Token)
def create_user_endpoint(user: UserCreate, db: Session = Depends(get_db)):
    if user_repo.user_exists(db, user):
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Username or email already taken",
        )
    db_user = user_repo.create_user(db, user)
    access_token = create_access_token(data={"id": db_user.id})
    return {"access_token": access_token, "token_type": "bearer"}


@router.get("/", response_model=UserList)
def list_user_endpoint(
    user_id: int = Depends(has_permission(PermissionEnum.LIST_USERS)),
    db: Session = Depends(get_db),
):
    # db_user = user_repo.get_user(db, str(current_user.email)
    users: List[User] = user_repo.all(db)
    return UserList(
        users=[
            UserInfo(email=str(u.email), name=str(u.name), username=str(u.username))
            for u in users
        ]
    )


@router.post("/token", response_model=Token)
def login(user: UserLogin, db: Session = Depends(get_db)):
    authenticated_user = user_repo.authenticate_user(db, user.username, user.password)
    if not authenticated_user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect email or password",
            headers={"WWW-Authenticate": "Bearer"},
        )
    access_token = create_access_token(data={"id": authenticated_user.id})
    return {"access_token": access_token, "token_type": "bearer"}


@router.get("/me", response_model=UserInfo)
def me(
    user_id: int = Depends(get_current_user), db: Session = Depends(get_db)
) -> UserInfo:
    user: User | None = user_repo.get_user(db, user_id)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid token. Please login again.",
            headers={"WWW-Authenticate": "Bearer"},
        )
    return UserInfo(email=user.email, name=user.name, username=user.username)
