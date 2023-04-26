from enum import IntEnum
from typing import List

from sqlalchemy import Enum
from sqlalchemy import ForeignKey
from sqlalchemy import String
from sqlalchemy.orm import Mapped
from sqlalchemy.orm import mapped_column
from sqlalchemy.orm import relationship
from src.models.base import Base


class Role(Base):
    __tablename__ = "roles"
    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
    name: Mapped[str] = mapped_column(
        String(255), nullable=False, unique=True, index=True
    )
    users: Mapped[List["User"]] = relationship(
        "User", secondary="user_roles", back_populates="roles"
    )
    permissions: Mapped[List["Permission"]] = relationship(
        "Permission", secondary="role_permissions", back_populates="roles"
    )


class Permission(Base):
    __tablename__ = "permissions"
    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
    name: Mapped[str] = mapped_column(String(255), nullable=False, unique=True)
    roles: Mapped[List["Role"]] = relationship(
        "Role", secondary="role_permissions", back_populates="permissions"
    )


class UserRole(Base):
    __tablename__ = "user_roles"
    user_id: Mapped[int] = mapped_column(ForeignKey("users.id"), primary_key=True)
    role_id: Mapped[int] = mapped_column(ForeignKey("roles.id"), primary_key=True)


class RolePermission(Base):
    __tablename__ = "role_permissions"
    role_id: Mapped[int] = mapped_column(ForeignKey("roles.id"), primary_key=True)
    permission_id: Mapped[int] = mapped_column(
        ForeignKey("permissions.id"), primary_key=True
    )


class User(Base):
    __tablename__ = "users"

    id: Mapped[int] = mapped_column(autoincrement=True, primary_key=True)
    email: Mapped[str] = mapped_column(
        String(255), index=True, unique=True, nullable=False
    )
    username: Mapped[str] = mapped_column(
        String(255), index=True, unique=True, nullable=False
    )
    name: Mapped[str] = mapped_column(String(255))
    hashed_password: Mapped[str] = mapped_column(String(255), nullable=False)

    api_keys: Mapped[List["ApiKey"]] = relationship(
        "ApiKey", back_populates="user"
    )  # Add apikeys relationship

    roles: Mapped[List[Role]] = relationship(
        "Role", secondary="user_roles", back_populates="users"
    )


class ApiKeyStatusEnum(IntEnum):
    INACTIVE = 0
    WAITING_ACTIVE = 1
    ACTIVE = 2
    WAITING_INATIVE = 3


class ApiKey(Base):
    __tablename__ = "api_key"

    id: Mapped[int] = mapped_column(primary_key=True, index=True, autoincrement=True)
    user_id: Mapped[int] = mapped_column(ForeignKey("users.id"), index=True)

    name: Mapped[str] = mapped_column(String(255), nullable=False)
    exchange: Mapped[str] = mapped_column(String(255), nullable=False)
    api_key: Mapped[str] = mapped_column(String(255), nullable=False)
    secret: Mapped[str] = mapped_column(String(255), nullable=False)
    status: Mapped[int] = mapped_column(default=ApiKeyStatusEnum.INACTIVE.value)

    user = relationship(
        "User",
        back_populates="api_keys",
    )

    class Config:
        orm_mode = True


class ClosedPnl(Base):
    __tablename__ = "closed_pnl"

    id: Mapped[int] = mapped_column(primary_key=True, index=True, autoincrement=True)
    user_id: Mapped[int] = mapped_column(ForeignKey("users.id"), index=True)
    api_key_id: Mapped[int] = mapped_column(ForeignKey("api_key.id"), index=True)

    symbol: Mapped[str] = mapped_column(String(100))
    orderId: Mapped[str] = mapped_column(String(100))
    side: Mapped[str] = mapped_column(String(100))
    qty: Mapped[str] = mapped_column(String(100))
    orderPrice: Mapped[str] = mapped_column(String(100))
    orderType: Mapped[str] = mapped_column(String(100))
    execType: Mapped[str] = mapped_column(String(100))
    closedSize: Mapped[str] = mapped_column(String(100))
    cumEntryValue: Mapped[str] = mapped_column(String(100))
    avgEntryPrice: Mapped[str] = mapped_column(String(100))
    cumExitValue: Mapped[str] = mapped_column(String(100))
    avgExitPrice: Mapped[str] = mapped_column(String(100))
    closedPnl: Mapped[str] = mapped_column(String(100))
    fillCount: Mapped[str] = mapped_column(String(100))
    leverage: Mapped[str] = mapped_column(String(100))
    createdTime: Mapped[str] = mapped_column(String(100))
    updatedTime: Mapped[str] = mapped_column(String(100))
