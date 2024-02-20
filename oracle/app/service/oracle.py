import logging

from app.api.oracle import oracle_pb2_grpc, oracle_pb2


class OracleService(oracle_pb2_grpc.OracleServicer):
    def __init__(self) -> None:
        pass

    def RetrainLFMv1(self, request, context):
        logging.info("start retrain LFMv1")

        return oracle_pb2.RetrainLFMv1Response()
