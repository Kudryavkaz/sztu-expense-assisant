# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: src/grpc/sztuea/v1/crawler.proto
# Protobuf Python Version: 5.28.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    28,
    1,
    '',
    'src/grpc/sztuea/v1/crawler.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n src/grpc/sztuea/v1/crawler.proto\x12!com.github.sztuea.grpc.crawler.v1\"5\n\x10GetCookieRequest\x12\x0f\n\x07\x61\x63\x63ount\x18\x01 \x01(\t\x12\x10\n\x08password\x18\x02 \x01(\t\"#\n\x11GetCookieResponse\x12\x0e\n\x06\x63ookie\x18\x01 \x01(\t2\x81\x01\n\x07\x43rawler\x12v\n\tGetCookie\x12\x33.com.github.sztuea.grpc.crawler.v1.GetCookieRequest\x1a\x34.com.github.sztuea.grpc.crawler.v1.GetCookieResponseB\nZ\x08./protosb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'src.grpc.sztuea.v1.crawler_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\010./protos'
  _globals['_GETCOOKIEREQUEST']._serialized_start=71
  _globals['_GETCOOKIEREQUEST']._serialized_end=124
  _globals['_GETCOOKIERESPONSE']._serialized_start=126
  _globals['_GETCOOKIERESPONSE']._serialized_end=161
  _globals['_CRAWLER']._serialized_start=164
  _globals['_CRAWLER']._serialized_end=293
# @@protoc_insertion_point(module_scope)