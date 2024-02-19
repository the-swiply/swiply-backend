from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class TaskStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNKNOWN: _ClassVar[TaskStatus]
    IN_PROGRESS: _ClassVar[TaskStatus]
    SUCCESS: _ClassVar[TaskStatus]
    ERROR: _ClassVar[TaskStatus]
UNKNOWN: TaskStatus
IN_PROGRESS: TaskStatus
SUCCESS: TaskStatus
ERROR: TaskStatus

class RetrainLFMv1Request(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class RetrainLFMv1Response(_message.Message):
    __slots__ = ("task_id",)
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    task_id: str
    def __init__(self, task_id: _Optional[str] = ...) -> None: ...

class GetTaskStatusRequest(_message.Message):
    __slots__ = ("task_id",)
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    task_id: str
    def __init__(self, task_id: _Optional[str] = ...) -> None: ...

class GetTaskStatusResponse(_message.Message):
    __slots__ = ("status", "details")
    STATUS_FIELD_NUMBER: _ClassVar[int]
    DETAILS_FIELD_NUMBER: _ClassVar[int]
    status: TaskStatus
    details: str
    def __init__(self, status: _Optional[_Union[TaskStatus, str]] = ..., details: _Optional[str] = ...) -> None: ...
