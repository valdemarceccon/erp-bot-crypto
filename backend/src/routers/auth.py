import datetime
from typing import Annotated

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
from src.repository import user as user_repo
from src.repository.user import user_has_permission
from src.schemas.user import Token
from src.schemas.user import UserCreate
from src.schemas.user import UserLogin

router = APIRouter(prefix="/auth", tags=["auth"])
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")


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
        user_id: int = int(payload["id"])
        return user_id

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


@router.post("/token", response_model=Token)
def login(user: UserLogin, db: Annotated[Session, Depends(get_db)]):
    authenticated_user = user_repo.authenticate_user(db, user.username, user.password)
    if not authenticated_user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect email or password",
            headers={"WWW-Authenticate": "Bearer"},
        )
    access_token = create_access_token(data={"id": authenticated_user.id})
    return {"access_token": access_token, "token_type": "bearer"}


@router.post("/signup", response_model=Token)
def create_user(user: UserCreate, db: Annotated[Session, Depends(get_db)]):
    if user_repo.user_exists(db, user):
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Username or email already taken",
        )
    db_user = user_repo.create_user(db, user)
    access_token = create_access_token(data={"id": db_user.id})
    return {"access_token": access_token, "token_type": "bearer"}
