# src/schemas/user.py
from typing import List

from pydantic import BaseModel
from pydantic import EmailStr


# Pydantic models
class UserCreateRequest(BaseModel):
    email: EmailStr
    password: str
    name: str
    username: str


class UserUpdateRequest(BaseModel):
    password: str
    name: str


class UserInfo(BaseModel):
    email: str
    name: str
    username: str


class UserLogin(BaseModel):
    username: str
    password: str


class UserList(BaseModel):
    users: List[UserInfo]


class Token(BaseModel):
    access_token: str
    token_type: str


class ApiKeyRequest(BaseModel):
    name: str
    key: str
    api_secret: str
    exchange: str
