# src/schemas/user.py
from typing import List

from pydantic import BaseModel
from pydantic import EmailStr
from src.models.user import ApiKeyStatusEnum


# Pydantic models
class UserCreateRequest(BaseModel):
    email: EmailStr
    password: str
    name: str
    username: str


class UserUpdateRequest(BaseModel):
    password: str
    name: str


class PermissionResp(BaseModel):
    name: str

    class Config:
        orm_mode = True


class UserInfo(BaseModel):
    email: str
    name: str
    username: str
    id: int

    class Config:
        orm_mode = True


class UserInfoMe(BaseModel):
    email: str
    name: str
    username: str

    permissions: List[PermissionResp]

    class Config:
        orm_mode = True


class ApiKeyRequestBase(BaseModel):
    name: str
    api_key: str
    exchange: str
    status: ApiKeyStatusEnum = ApiKeyStatusEnum.INACTIVE


class ApiKeyRequestIn(ApiKeyRequestBase):
    api_secret: str


class ApiKeyRequestOut(ApiKeyRequestBase):
    id: int

    class Config:
        orm_mode = True


class ApiKeyAdminRequestOut(ApiKeyRequestOut):
    id: int
    user: UserInfo
    api_key: str
    secret: str

    class Config:
        orm_mode = True


class UserDetail(UserInfo):
    api_keys: List[ApiKeyRequestOut]


class UserLogin(BaseModel):
    username: str
    password: str


class Token(BaseModel):
    access_token: str
    token_type: str


class ApiKeyRequestUpdate(BaseModel):
    name: str | None
    api_key: str | None
    exchange: str | None
    status: ApiKeyStatusEnum | None = ApiKeyStatusEnum.INACTIVE
    api_secret: str | None
