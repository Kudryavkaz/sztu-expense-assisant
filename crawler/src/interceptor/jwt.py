from typing import Awaitable, Callable

from grpc.aio import ServerInterceptor
from grpc.aio._metadata import Metadata

from grpc import HandlerCallDetails, RpcMethodHandler
from src.resource.jwt import parse_token, validate_payload

METHOD_WHITE_LIST = [
    "Check",
    "Watch",
]


class JwtInterceptor(ServerInterceptor):
    async def intercept_service(
        self, continuation: Callable[[HandlerCallDetails], Awaitable[RpcMethodHandler]], handler_call_details: HandlerCallDetails
    ) -> RpcMethodHandler:
        self._validate_request(handler_call_details)
        return await continuation(handler_call_details)

    def _validate_request(self, handler_call_details) -> None:
        metadata = Metadata(*handler_call_details.invocation_metadata)
        token = metadata["authorization"]
        payload = parse_token(token)
        validate_payload(payload)
