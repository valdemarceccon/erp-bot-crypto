# routers/users.py
from typing import Annotated
from typing import List

from fastapi import APIRouter
from fastapi import Depends
from fastapi import HTTPException
from fastapi import status
from sqlalchemy.orm import Session
from src.dependencies.database import get_db
from src.models.roles import PermissionEnum
from src.models.user import User
from src.repository import user as user_repo
from src.routers.auth import get_current_user
from src.routers.auth import has_permission
from src.schemas.user import ApiKeyRequest
from src.schemas.user import UserInfo
from src.schemas.user import UserList
from src.schemas.user import UserUpdateRequest

router = APIRouter(prefix="/users", tags=["users"])


# API endpoints
@router.patch("/", response_model=UserInfo)
def update_user(
    user: UserUpdateRequest,
    db: Annotated[Session, Depends(get_db)],
    user_id: Annotated[int, Depends(get_current_user)],
):
    db_user = user_repo.update_user(db, user_id, user)
    if db_user is None:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid token",
            headers={"WWW-Authenticate": "Bearer"},
        )

    return UserInfo(email=db_user.email, name=db_user.name, username=db_user.username)


@router.get("/", response_model=UserList)
def list_user(
    _: Annotated[int, Depends(has_permission(PermissionEnum.LIST_USERS))],
    db: Annotated[Session, Depends(get_db)],
):
    users: List[User] = user_repo.all(db)
    return UserList(
        users=[
            UserInfo(email=str(u.email), name=str(u.name), username=str(u.username))
            for u in users
        ]
    )


@router.get("/me", response_model=UserInfo)
def me(
    user_id: Annotated[int, Depends(get_current_user)],
    db: Annotated[Session, Depends(get_db)],
) -> UserInfo:
    user: User | None = user_repo.get_user(db, user_id)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid token. Please login again.",
            headers={"WWW-Authenticate": "Bearer"},
        )
    return UserInfo(email=user.email, name=user.name, username=user.username)


@router.post("/api_key", status_code=status.HTTP_201_CREATED)
def add_apikey(
    api_key_req: ApiKeyRequest,
    user_id: Annotated[int, Depends(get_current_user)],
    db: Annotated[Session, Depends(get_db)],
):
    user_repo.add_api_key(db, user_id, api_key_req)

    return {"message": "ok"}
