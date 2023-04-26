from enum import Enum


class PermissionEnum(Enum):
    LIST_USERS = "ListUsers"
    GET_USER_INFO = "GetUserInfo"
    LIST_API_KEYS = "ListApiKeys"
    # Add other permissions as needed
