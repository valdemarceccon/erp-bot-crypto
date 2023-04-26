import itertools
from typing import List
from typing import Type

from passlib.hash import bcrypt
from sqlalchemy import or_
from sqlalchemy.orm import Session
from src.models.user import ApiKey
from src.models.user import Role
from src.models.user import User
from src.schemas.user import ApiKeyRequestIn
from src.schemas.user import ApiKeyRequestUpdate
from src.schemas.user import UserCreateRequest
from src.schemas.user import UserUpdateRequest


def get_password_hash(password: str) -> str:
    return bcrypt.hash(password)


def verify_password(plain_password: str, hashed_password: str) -> bool:
    return bcrypt.verify(plain_password, hashed_password)


def create_user(db: Session, user: UserCreateRequest) -> User:
    hashed_password = get_password_hash(user.password)
    db_user = User(
        email=user.email,
        hashed_password=hashed_password,
        name=user.name,
        username=user.username,
    )
    db.add(db_user)
    db.commit()
    db.refresh(db_user)
    return db_user


def user_exists(db: Session, user: UserCreateRequest) -> bool:
    dbuser = (
        db.query(User)
        .filter(or_(User.email == user.email, User.username == user.username))
        .first()
    )
    if dbuser is None:
        return False
    return True


def update_user(db: Session, user_id: int, user: UserUpdateRequest) -> User | None:
    hashed_password = get_password_hash(user.password)
    db_user = db.get(User, user_id)
    if not db_user:
        return None

    db_user.hashed_password = hashed_password
    db.commit()
    db.refresh(db_user)
    return db_user


def authenticate_user(db: Session, username: str, password: str) -> User | None:
    user = get_user_by_username(db, username)

    if not user:
        return None

    if not verify_password(password, str(user.hashed_password)):
        return None
    return user


def get_user(db: Session, user_id: int) -> User | None:
    return db.query(User).filter(User.id == user_id).first()


def get_user_by_username(db: Session, username: str) -> User | None:
    user: User | None = db.query(User).filter(User.username == username).first()
    if not user:
        return None

    return user


def get_all(db: Session) -> List[User]:
    return db.query(User).all()


def user_has_permission(session: Session, user_id: int, permission_name: str) -> bool:
    user = session.query(User).filter(User.id == user_id).first()

    if not user:
        return False

    user_roles: List[Role] = user.roles
    permissions = [ur.permissions for ur in user_roles]
    permissions = itertools.chain(*permissions)

    return permission_name in [p.name for p in permissions]


def add_api_key(session: Session, user_id: int, api_key: ApiKeyRequestIn) -> ApiKey:
    api_key_db = ApiKey(
        name=api_key.name,
        api_key=api_key.api_key,
        secret=api_key.api_secret,
        exchange=api_key.exchange,
        user_id=user_id,
    )
    session.add(api_key_db)
    session.commit()
    session.refresh(api_key_db)
    return api_key_db


def get_user_api_keys(db: Session, user_id: int) -> List[ApiKey]:
    return db.query(ApiKey).filter(ApiKey.user_id == user_id).all()


def get_api_key(db: Session, user_id: int, api_key_id: int) -> ApiKey | None:
    return (
        db.query(ApiKey)
        .filter(ApiKey.user_id == user_id, ApiKey.id == api_key_id)
        .first()
    )


def list_api_keys(db: Session) -> List[ApiKey] | None:
    return db.query(ApiKey).all()


def delete_api_key(db: Session, user_id: int, api_key_id: int) -> None:
    api_keys = (
        db.query(ApiKey)
        .filter(ApiKey.user_id == user_id, ApiKey.id == api_key_id)
        .first()
    )
    if api_keys:
        db.delete(api_keys)
        db.commit()


def update_api_key(
    db: Session, user_id: int, api_key_id: int, api_key: ApiKeyRequestUpdate
) -> ApiKey | None:
    db_api_key = (
        db.query(ApiKey)
        .filter(ApiKey.user_id == user_id, ApiKey.id == api_key_id)
        .first()
    )
    if not db_api_key:
        return None

    db_api_key.name = api_key.name if api_key.name else db_api_key.name
    db_api_key.status = api_key.status.value if api_key.status else db_api_key.status
    db_api_key.exchange = api_key.exchange if api_key.exchange else db_api_key.exchange
    db_api_key.api_key = api_key.api_key if api_key.api_key else db_api_key.api_key
    db_api_key.secret = api_key.api_secret if api_key.api_secret else db_api_key.secret
    db.commit()
    return db_api_key
