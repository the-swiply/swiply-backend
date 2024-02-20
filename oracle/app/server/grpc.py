import logging
from concurrent import futures
from typing import NoReturn

import grpc
from app.api.oracle import oracle_pb2_grpc


class OracleServer:
    def __init__(self, addr: str, oracle_service) -> None:
        self.__address = addr
        self.__server = grpc.server(futures.ThreadPoolExecutor(max_workers=2))
        self.__server.add_insecure_port(self.__address)

        oracle_pb2_grpc.add_OracleServicer_to_server(oracle_service, self.__server)

    def serve(self) -> NoReturn:
        self.__server.start()
        logging.info(f'server started on {self.__address}')

    def stop(self, grace_timeout_seconds: int) -> NoReturn:
        self.__server.stop(grace=grace_timeout_seconds)
        logging.info('server stopped')
