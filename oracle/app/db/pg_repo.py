import psycopg2


class OracleRepository:
    def __init__(self, host, port, user, password, db, ssl):
        self.__conn = psycopg2.connect(user=user, password=password, host=host, port=port, dbname=db, sslmode=ssl)

    def close(self):
        self.__conn.close()
