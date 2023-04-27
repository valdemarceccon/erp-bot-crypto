from enum import Enum


class PermissionEnum(Enum):
    LIST_USERS = "ListUsers"
    GET_USER_INFO = "GetUserInfo"
    LIST_API_KEYS = "ListApiKeys"
    WRITE_ALL_API_KEYS = "WriteAllApiKeys"
    READ_ALL_API_KEYS = "ReadAllApiKeys"
    # Add other permissions as needed
