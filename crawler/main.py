import asyncio
from concurrent import futures
from loguru import logger
from grpc.aio import ServicerContext
from grpc.aio import server as grpc_aio_server

from src.config import GRPC_THREAD_POOL_WORKER, GRPC_SERVER_PORT
from src.interceptor import JwtInterceptor
from src.grpc.health.v1.health_pb2 import HealthCheckRequest, HealthCheckResponse
from src.grpc.sztuea.v1.crawler_pb2 import GetCookieRequest, GetCookieResponse
from src.grpc.health.v1.health_pb2_grpc import HealthServicer, add_HealthServicer_to_server
from src.grpc.sztuea.v1.crawler_pb2_grpc import CrawlerServicer, add_CrawlerServicer_to_server
from src.handler.get_cookie.handler import get_cookie_handler


class _CrawlerServicer(CrawlerServicer):
    async def GetCookie(self, request: GetCookieRequest, context: ServicerContext) -> GetCookieResponse:
        response = get_cookie_handler(request, context)
        return response


class _HealthServicer(HealthServicer):
    async def Check(self, request: HealthCheckRequest, context: ServicerContext) -> HealthCheckResponse:
        return HealthCheckResponse(status=HealthCheckResponse.SERVING)


async def start_server():
    server = grpc_aio_server(
        futures.ThreadPoolExecutor(max_workers=GRPC_THREAD_POOL_WORKER),
        interceptors=[JwtInterceptor()],
        options=[
            ("grpc.so_reuseport", 0),
            ("grpc.enable_retries", 1),
            ("grpc.enable_strict_message_length", False),
        ],
    )

    add_CrawlerServicer_to_server(_CrawlerServicer(), server)

    add_HealthServicer_to_server(_HealthServicer(), server)

    server.add_insecure_port(f"[::]:{GRPC_SERVER_PORT}")

    logger.info(f"[Start Server] Stared server on port: {GRPC_SERVER_PORT}")

    await server.start()

    # since server.start() will not block,
    # a sleep-loop is added to keep alive
    try:
        await server.wait_for_termination()
    except KeyboardInterrupt:
        await server.stop(None)

if __name__ == "__main__":
    asyncio.run(start_server())
