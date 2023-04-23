import itertools
from typing import List

from fastapi import Depends
from fastapi import HTTPException
from passlib.hash import bcrypt
from sqlalchemy import or_
from sqlalchemy import select
from sqlalchemy.orm import Session
from src.dependencies.database import get_db
from src.models.roles import PermissionEnum
from src.models.user import Permission
from src.models.user import Role
from src.models.user import RolePermission
from src.models.user import User
from src.models.user import UserRole
from src.schemas.user import UserCreate
from src.schemas.user import UserUpdate


def get_password_hash(password: str) -> str:
    return bcrypt.hash(password)


def verify_password(plain_password: str, hashed_password: str) -> bool:
    return bcrypt.verify(plain_password, hashed_password)


def create_user(db: Session, user: UserCreate) -> User:
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


def user_exists(db: Session, user: UserCreate) -> bool:
    dbuser = (
        db.query(User)
        .filter(or_(User.email == user.email, User.username == user.username))
        .first()
    )
    if dbuser is None:
        return False
    return True


def update_user(db: Session, user_id: int, user: UserUpdate) -> User | None:
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
    user: User | None = db.query(User).filter(User.id == user_id).first()
    if not user:
        return None

    return user


def get_user_by_username(db: Session, username: str) -> User | None:
    user: User | None = db.query(User).filter(User.username == username).first()
    if not user:
        return None

    return user


def all(db: Session) -> List[User]:
    return db.query(User).all()


# Add this code to main.py
from typing import List


def user_has_permission(session: Session, user_id: int, permission_name: str) -> bool:
    user = session.query(User).filter(User.id == user_id).first()

    if not user:
        return False

    user_roles: List[Role] = user.roles
    permissions = [ur.permissions for ur in user_roles]
    permissions = itertools.chain(*permissions)

    return permission_name in [p.name for p in permissions]
