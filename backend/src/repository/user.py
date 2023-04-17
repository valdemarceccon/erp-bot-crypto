from typing import List

from fastapi import Depends
from fastapi import HTTPException
from passlib.hash import bcrypt
from sqlalchemy import select
from sqlalchemy.orm import Session
from src.dependencies.database import get_db
from src.models.user import Permission
from src.models.user import Role
from src.models.user import RolePermission
from src.models.user import User
from src.models.user import UserRole
from src.schemas.user import UserCreate


def get_password_hash(password: str) -> str:
    return bcrypt.hash(password)


def verify_password(plain_password: str, hashed_password: str) -> bool:
    return bcrypt.verify(plain_password, hashed_password)


def create_user(db: Session, user: UserCreate) -> User:
    hashed_password = get_password_hash(user.password)
    db_user = User(email=user.email, hashed_password=hashed_password, name=user.name)
    db.add(db_user)
    db.commit()
    db.refresh(db_user)
    return db_user


def authenticate_user(db: Session, email: str, password: str) -> User | None:
    user = get_user(db, email)

    if not user:
        return None

    if not verify_password(password, str(user.hashed_password)):
        return None
    return user


def get_user(db: Session, email: str) -> User | None:
    user: User = db.query(User).filter(User.email == email).first()
    if not user:
        return None

    return user


def all(db: Session) -> List[User]:
    return db.query(User).all()


# Add this code to main.py
from typing import List


def user_has_permission(session: Session, user: User, permission_name: str) -> bool:
    stmt = (
        select(RolePermission)
        .join(Role, RolePermission.role_id == Role.id)
        .join(UserRole, UserRole.role_id == Role.id)
        .join(Permission, RolePermission.permission_id == Permission.id)
        .where(UserRole.user_id == user.id, Permission.name == permission_name)
    )

    result = session.execute(stmt)
    return result.scalar_one_or_none() is not None
