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
from src.models.user import ApiKeyStatusEnum
from src.models.user import Permission
from src.models.user import User
from src.repository import user as user_repo
from src.routers.auth import get_current_user
from src.routers.auth import has_permission
from src.schemas.user import ApiKeyRequestIn
from src.schemas.user import ApiKeyRequestOut
from src.schemas.user import ApiKeyRequestUpdate
from src.schemas.user import PermissionResp
from src.schemas.user import UserDetail
from src.schemas.user import UserInfo
from src.schemas.user import UserUpdateRequest

router = APIRouter(prefix="/users", tags=["users"])


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


@router.get("/", response_model=List[UserInfo])
def list_user(
    _: Annotated[int, Depends(has_permission(PermissionEnum.LIST_USERS))],
    db: Annotated[Session, Depends(get_db)],
) -> Any:
    return user_repo.get_all(db)


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


@router.get("/permissions", response_model=List[PermissionResp])
def get_user_permissions(
    user_id: Annotated[int, Depends(get_current_user)],
    db: Annotated[Session, Depends(get_db)],
):
    user_permissions: List[Permission] | None = user_repo.list_user_permissions(
        db, user_id=user_id
    )
    if user_permissions is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="User not found"
        )

    return user_permissions


@router.get("/{user_id}", response_model=UserDetail)
def get_user_detail(
    user_id: Annotated[int, Depends(has_permission(PermissionEnum.GET_USER_INFO))],
    db: Annotated[Session, Depends(get_db)],
):
    user_detail: User | None = user_repo.get_user(db, user_id=user_id)
    if not user_detail:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="User not found"
        )

    return user_detail


@router.get("/api_keys/all", response_model=List[ApiKeyRequestOut])
def get_apikey_all_users(
    user_id: Annotated[int, Depends(has_permission(PermissionEnum.LIST_API_KEYS))],
    db: Annotated[Session, Depends(get_db)],
) -> Any:
    api_key_list = user_repo.list_api_keys(db=db)
    return api_key_list


@router.get("/api_key/{api_key_id}", response_model=ApiKeyRequestOut)
def get_apikey_id(
    api_key_id: int,
    user_id: Annotated[int, Depends(get_current_user)],
    db: Annotated[Session, Depends(get_db)],
) -> Any:
    api_key_detail = user_repo.get_api_key(
        db=db, api_key_id=api_key_id, user_id=user_id
    )
    if not api_key_detail:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Api key with id {api_key_id} not found",
        )
    return api_key_detail


@router.get("/api_keys/", response_model=List[ApiKeyRequestOut])
def get_apikey(
    user_id: Annotated[int, Depends(get_current_user)],
    db: Annotated[Session, Depends(get_db)],
) -> Any:
    api_key_list = user_repo.get_user_api_keys(db, user_id)
    return api_key_list


@router.post(
    "/api_key", status_code=status.HTTP_201_CREATED, response_model=ApiKeyRequestOut
)
def add_apikey(
    api_key_req: ApiKeyRequestIn,
    user_id: Annotated[int, Depends(get_current_user)],
    db: Annotated[Session, Depends(get_db)],
):
    return user_repo.add_api_key(db, user_id, api_key_req)


@router.delete("/api_key/{api_key_id}", status_code=status.HTTP_202_ACCEPTED)
def delete_api_key(
    api_key_id: int,
    user_id: Annotated[int, Depends(get_current_user)],
    db: Annotated[Session, Depends(get_db)],
) -> Any:
    user_repo.delete_api_key(db, user_id, api_key_id)
    return {"message": "ok"}


def get_toggle_status_client(status: int) -> ApiKeyStatusEnum | None:
    if status == ApiKeyStatusEnum.ACTIVE:
        return ApiKeyStatusEnum.WAITING_INATIVE

    if status == ApiKeyStatusEnum.INACTIVE:
        return ApiKeyStatusEnum.WAITING_ACTIVE

    if status == ApiKeyStatusEnum.WAITING_ACTIVE:
        return ApiKeyStatusEnum.INACTIVE

    if status == ApiKeyStatusEnum.WAITING_INATIVE:
        return ApiKeyStatusEnum.ACTIVE

    return None


@router.patch(
    "/api_key/client-toggle/{api_key_id}", status_code=status.HTTP_202_ACCEPTED
)
def client_api_key_toggle(
    api_key_id: int,
    user_id: Annotated[int, Depends(get_current_user)],
    db: Annotated[Session, Depends(get_db)],
) -> Any:
    current_api_key = user_repo.get_api_key(db, user_id=user_id, api_key_id=api_key_id)
    if not current_api_key:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Api key with id {api_key_id} not found",
        )
    old_status = current_api_key.status
    new_status = get_toggle_status_client(old_status)

    if new_status is None:
        raise HTTPException(
            status_code=status.HTTP_422_UNPROCESSABLE_ENTITY,
            detail=f"Invalid status {old_status}",
        )

    user_repo.update_api_key(
        db=db,
        user_id=user_id,
        api_key_id=api_key_id,
        api_key=ApiKeyRequestUpdate(
            status=new_status, api_key=None, api_secret=None, exchange=None, name=None
        ),
    )
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
