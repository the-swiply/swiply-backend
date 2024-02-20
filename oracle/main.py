from app.service import oracle

from app.houston import config, runner
from app.server import grpc


if __name__ == "__main__":
    config.parse_yaml()

    oracle_service = oracle.OracleService()

    server = grpc.OracleServer(config.get("grpc").get("addr"), oracle_service)

    server.serve()

    runner.terminator().wait()

    server.stop(config.get("app").get("graceful_shutdown_timeout_seconds"))
