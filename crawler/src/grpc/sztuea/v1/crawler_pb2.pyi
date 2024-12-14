from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class GetCookieRequest(_message.Message):
    __slots__ = ("account", "password")
    ACCOUNT_FIELD_NUMBER: _ClassVar[int]
    PASSWORD_FIELD_NUMBER: _ClassVar[int]
    account: str
    password: str
    def __init__(self, account: _Optional[str] = ..., password: _Optional[str] = ...) -> None: ...

class GetCookieResponse(_message.Message):
    __slots__ = ("cookie",)
    COOKIE_FIELD_NUMBER: _ClassVar[int]
    cookie: str
    def __init__(self, cookie: _Optional[str] = ...) -> None: ...
