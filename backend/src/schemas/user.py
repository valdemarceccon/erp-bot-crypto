# src/schemas/user.py
from typing import List

from pydantic import BaseModel
from pydantic import EmailStr


# Pydantic models
class UserCreate(BaseModel):
    email: EmailStr
    password: str
    name: str


class UserUpdate(BaseModel):
    password: str
    name: str


class UserInfo(BaseModel):
    email: str
    name: str


class UserLogin(BaseModel):
    email: str
    password: str


class UserList(BaseModel):
    users: List[UserInfo]


class Token(BaseModel):
    access_token: str
    token_type: str
