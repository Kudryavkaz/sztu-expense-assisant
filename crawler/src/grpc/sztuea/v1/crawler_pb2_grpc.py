# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc
import warnings

from src.grpc.sztuea.v1 import crawler_pb2 as src_dot_grpc_dot_sztuea_dot_v1_dot_crawler__pb2

GRPC_GENERATED_VERSION = '1.68.1'
GRPC_VERSION = grpc.__version__
_version_not_supported = False

try:
    from grpc._utilities import first_version_is_lower
    _version_not_supported = first_version_is_lower(GRPC_VERSION, GRPC_GENERATED_VERSION)
except ImportError:
    _version_not_supported = True

if _version_not_supported:
    raise RuntimeError(
        f'The grpc package installed is at version {GRPC_VERSION},'
        + f' but the generated code in src/grpc/sztuea/v1/crawler_pb2_grpc.py depends on'
        + f' grpcio>={GRPC_GENERATED_VERSION}.'
        + f' Please upgrade your grpc module to grpcio>={GRPC_GENERATED_VERSION}'
        + f' or downgrade your generated code using grpcio-tools<={GRPC_VERSION}.'
    )


class CrawlerStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetCookie = channel.unary_unary(
                '/com.github.sztuea.grpc.crawler.v1.Crawler/GetCookie',
                request_serializer=src_dot_grpc_dot_sztuea_dot_v1_dot_crawler__pb2.GetCookieRequest.SerializeToString,
                response_deserializer=src_dot_grpc_dot_sztuea_dot_v1_dot_crawler__pb2.GetCookieResponse.FromString,
                _registered_method=True)


class CrawlerServicer(object):
    """Missing associated documentation comment in .proto file."""

    def GetCookie(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_CrawlerServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GetCookie': grpc.unary_unary_rpc_method_handler(
                    servicer.GetCookie,
                    request_deserializer=src_dot_grpc_dot_sztuea_dot_v1_dot_crawler__pb2.GetCookieRequest.FromString,
                    response_serializer=src_dot_grpc_dot_sztuea_dot_v1_dot_crawler__pb2.GetCookieResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'com.github.sztuea.grpc.crawler.v1.Crawler', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('com.github.sztuea.grpc.crawler.v1.Crawler', rpc_method_handlers)


 # This class is part of an EXPERIMENTAL API.
class Crawler(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def GetCookie(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/com.github.sztuea.grpc.crawler.v1.Crawler/GetCookie',
            src_dot_grpc_dot_sztuea_dot_v1_dot_crawler__pb2.GetCookieRequest.SerializeToString,
            src_dot_grpc_dot_sztuea_dot_v1_dot_crawler__pb2.GetCookieResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
