# routers/users.py
from typing import Annotated
from typing import Any
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
from src.schemas.user import ApiKeyRequestIn
from src.schemas.user import ApiKeyRequestOut
from src.schemas.user import ApiKeyRequestUpdate
from src.schemas.user import UserDetail
from src.schemas.user import UserInfo
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


@router.get("/", response_model=List[UserInfo])
def list_user(
    _: Annotated[int, Depends(has_permission(PermissionEnum.LIST_USERS))],
    db: Annotated[Session, Depends(get_db)],
) -> Any:
    return user_repo.get_all(db)


@router.get("/{user_id}", response_model=UserDetail)
def get_user_detail(
    user_id: int,
    _: Annotated[int, Depends(has_permission(PermissionEnum.GET_USER_INFO))],
    db: Annotated[Session, Depends(get_db)],
):
    user_detail: User | None = user_repo.get_user(db, user_id=user_id)
    if not user_detail:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="User not found"
        )

    return user_detail


@router.post(
    "/api_key", status_code=status.HTTP_201_CREATED, response_model=ApiKeyRequestOut
)
def add_apikey(
    api_key_req: ApiKeyRequestIn,
    user_id: Annotated[int, Depends(get_current_user)],
    db: Annotated[Session, Depends(get_db)],
):
    return user_repo.add_api_key(db, user_id, api_key_req)


@router.get("/api_key", response_model=List[ApiKeyRequestOut])
def get_apikey(
    user_id: Annotated[int, Depends(get_current_user)],
    db: Annotated[Session, Depends(get_db)],
) -> Any:
    api_key_list = user_repo.get_api_key(db, user_id)
    return api_key_list


@router.get("/api_key/{api_key_id}", response_model=List[ApiKeyRequestOut])
def get_apikey_id(
    api_key_id: int,
    user_id: Annotated[int, Depends(get_current_user)],
    db: Annotated[Session, Depends(get_db)],
) -> Any:
    api_key_list = user_repo.get_api_key(db, user_id, api_key_id)
    return api_key_list


@router.delete("/api_key/{api_key_id}", status_code=status.HTTP_202_ACCEPTED)
def delete_api_key(
    api_key_id: int,
    user_id: Annotated[int, Depends(get_current_user)],
    db: Annotated[Session, Depends(get_db)],
) -> Any:
    user_repo.delete_api_key(db, user_id, api_key_id)
    return {"message": "ok"}


@router.patch(
    "/api_key/{api_key_id}",
    response_model=ApiKeyRequestOut,
    status_code=status.HTTP_202_ACCEPTED,
)
def update_api_key(
    api_key_id: int,
    values: ApiKeyRequestUpdate,
    user_id: Annotated[int, Depends(get_current_user)],
    db: Annotated[Session, Depends(get_db)],
) -> Any:
    updated_api_key = user_repo.update_api_key(db, user_id, api_key_id, values)

    if not updated_api_key:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="Api key not found"
        )

    return updated_api_key
