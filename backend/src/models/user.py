from sqlalchemy import Column
from sqlalchemy import ForeignKey
from sqlalchemy import Integer
from sqlalchemy import String
from sqlalchemy.orm import relationship
from src.models.base import Base


class User(Base):
    __tablename__ = "users"

    email = Column(String, primary_key=True, index=True)
    name = Column(String, index=True)
    hashed_password = Column(String)

    apikeys = relationship(
        "APIKey", back_populates="user", lazy="select"
    )  # Add apikeys relationship


class APIKey(Base):
    __tablename__ = "apikeys"

    id = Column(Integer, primary_key=True, index=True, autoincrement=True)
    user_email = Column(String, ForeignKey(User.email), primary_key=True)
    apikey = Column(String, nullable=False)

    user = relationship("User", back_populates="apikeys")  # Add user relationship
