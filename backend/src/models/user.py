from sqlalchemy import Column
from sqlalchemy import ForeignKey
from sqlalchemy import Integer
from sqlalchemy import String
from sqlalchemy.orm import relationship
from src.models.base import Base


class Role(Base):
    __tablename__ = "roles"
    id = Column(Integer, primary_key=True, autoincrement=True)
    name = Column(String, nullable=False, unique=True)
    users = relationship("User", secondary="user_roles", back_populates="roles")
    permissions = relationship(
        "Permission", secondary="role_permissions", back_populates="roles"
    )


class Permission(Base):
    __tablename__ = "permissions"
    id = Column(Integer, primary_key=True, autoincrement=True)
    name = Column(String, nullable=False, unique=True)
    roles = relationship(
        "Role", secondary="role_permissions", back_populates="permissions"
    )


class UserRole(Base):
    __tablename__ = "user_roles"
    user_id = Column(String, ForeignKey("users.email"), primary_key=True)
    role_id = Column(Integer, ForeignKey("roles.id"), primary_key=True)


class RolePermission(Base):
    __tablename__ = "role_permissions"
    role_id = Column(Integer, ForeignKey("roles.id"), primary_key=True)
    permission_id = Column(Integer, ForeignKey("permissions.id"), primary_key=True)


class User(Base):
    __tablename__ = "users"

    id = Column(Integer, autoincrement=True, primary_key=True)
    email = Column(String, index=True, unique=True, nullable=False)
    username = Column(String, index=True, unique=True, nullable=False)
    name = Column(String)
    hashed_password = Column(String, nullable=False)

    apikeys = relationship(
        "APIKey", back_populates="user", lazy="select"
    )  # Add apikeys relationship

    roles = relationship("Role", secondary="user_roles", back_populates="users")


class APIKey(Base):
    __tablename__ = "apikeys"

    id = Column(Integer, primary_key=True, index=True, autoincrement=True)
    user_id = Column(Integer, ForeignKey(User.id), primary_key=True)
    apikey = Column(String, nullable=False)

    user = relationship(
        "User",
        back_populates="apikeys",
    )  # Add user relationship
