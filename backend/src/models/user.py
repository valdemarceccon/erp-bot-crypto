from sqlalchemy import Column
from sqlalchemy import ForeignKey
from sqlalchemy import Integer
from sqlalchemy import String
from sqlalchemy.orm import Relationship
from sqlalchemy.orm import relationship
from src.models.base import Base


class Role(Base):
    __tablename__ = "users_role"
    id = Column(Integer, primary_key=True, index=True, autoincrement=True)
    name = Column(String, nullable=False, index=True)

    users = relationship("User", back_populates="role")


class User(Base):
    __tablename__ = "users"

    email = Column(String, primary_key=True, index=True)
    name = Column(String, index=True)
    hashed_password = Column(String)
    role_id = Column(Integer, ForeignKey(Role.id), nullable=True)

    apikeys = relationship(
        "APIKey", back_populates="user", lazy="select"
    )  # Add apikeys relationship

    role: Relationship[Role] = relationship("Role")


class APIKey(Base):
    __tablename__ = "apikeys"

    id = Column(Integer, primary_key=True, index=True, autoincrement=True)
    user_email = Column(String, ForeignKey(User.email), primary_key=True)
    apikey = Column(String, nullable=False)

    user = relationship(
        "User",
        back_populates="apikeys",
    )  # Add user relationship
