import os

from app.db.pg_repo import OracleRepository
from app.service import oracle

from app.houston import config, runner, loggy
from app.server import grpc

if __name__ == "__main__":
    config.parse_yaml()
    loggy.init(config.get("app").get("name"))

    oracle_repo = OracleRepository(
        config.get("postgres").get("host"),
        config.get("postgres").get("port"),
        config.get("postgres").get("username"),
        os.environ.get("POSTGRES_PASSWORD"),
        config.get("postgres").get("db_name"),
        config.get("postgres").get("ssl_mode"),
    )

    oracle_service = oracle.OracleService(oracle_repo)
    oracle_service.RetrainLFMv1(None, None)

    server = grpc.OracleServer(config.get("grpc").get("addr"), oracle_service)
    server.serve()
    runner.terminator().wait()

    server.stop(config.get("app").get("graceful_shutdown_timeout_seconds"))
