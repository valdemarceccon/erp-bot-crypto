from sqlalchemy import Column
from sqlalchemy import ForeignKey
from sqlalchemy import Integer
from sqlalchemy import String
from src.models.base import Base

# from sqlalchemy.orm import relationship


class User(Base):
    __tablename__ = "users"

    id = Column(Integer, primary_key=True, index=True, autoincrement=True)
    email = Column(String, primary_key=True, index=True)
    name = Column(String, index=True)
    hashed_password = Column(String)

    # apikeys = relationship('APIKey', backref='user')


class APIKey(Base):
    __tablename__ = "apikeys"

    id = Column(Integer, primary_key=True, index=True, autoincrement=True)
    user_id = Column(Integer, ForeignKey(User.id), primary_key=True)
    apikey = Column(String, nullable=False)

    # user = relationship('User', backref='apikeys')
