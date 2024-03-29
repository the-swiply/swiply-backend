# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import app.api.oracle.oracle_pb2 as oracle__pb2


class OracleStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.RetrainLFMv1 = channel.unary_unary(
                '/swiply.oracle.Oracle/RetrainLFMv1',
                request_serializer=oracle__pb2.RetrainLFMv1Request.SerializeToString,
                response_deserializer=oracle__pb2.RetrainLFMv1Response.FromString,
                )
        self.GetTaskStatus = channel.unary_unary(
                '/swiply.oracle.Oracle/GetTaskStatus',
                request_serializer=oracle__pb2.GetTaskStatusRequest.SerializeToString,
                response_deserializer=oracle__pb2.GetTaskStatusResponse.FromString,
                )


class OracleServicer(object):
    """Missing associated documentation comment in .proto file."""

    def RetrainLFMv1(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetTaskStatus(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_OracleServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'RetrainLFMv1': grpc.unary_unary_rpc_method_handler(
                    servicer.RetrainLFMv1,
                    request_deserializer=oracle__pb2.RetrainLFMv1Request.FromString,
                    response_serializer=oracle__pb2.RetrainLFMv1Response.SerializeToString,
            ),
            'GetTaskStatus': grpc.unary_unary_rpc_method_handler(
                    servicer.GetTaskStatus,
                    request_deserializer=oracle__pb2.GetTaskStatusRequest.FromString,
                    response_serializer=oracle__pb2.GetTaskStatusResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'swiply.oracle.Oracle', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Oracle(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def RetrainLFMv1(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/swiply.oracle.Oracle/RetrainLFMv1',
            oracle__pb2.RetrainLFMv1Request.SerializeToString,
            oracle__pb2.RetrainLFMv1Response.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetTaskStatus(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/swiply.oracle.Oracle/GetTaskStatus',
            oracle__pb2.GetTaskStatusRequest.SerializeToString,
            oracle__pb2.GetTaskStatusResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
