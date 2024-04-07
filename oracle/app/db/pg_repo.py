import csv
from io import StringIO

import psycopg2


class OracleRepository:
    def __init__(self, host, port, user, password, db, ssl):
        self.__user = user
        self.__password = password
        self.__host = host
        self.__port = port
        self.__db = db
        self.__ssl = ssl

    def get_lfmv1_train_data(self):
        # TODO: query
        conn = psycopg2.connect(user=self.__user, password=self.__password, host=self.__host, port=self.__port, dbname=self.__db, sslmode=self.__ssl)
        cursor = conn.cursor()
        cursor.execute('''SELECT i."from" AS "from", i."to" AS "to", i.positive AS "positive", p1.updated_at AS updated_at__from, p2.updated_at  AS updated_at__to
FROM profile AS p1
LEFT JOIN interaction AS i ON i."from" = p1.id
LEFT JOIN  profile AS p2 ON i."to" = p2.id;''')
        data = cursor.fetchall()

        return data

    def update_lfmv1_results(self, results):
        buf = StringIO()

        writer = csv.writer(buf)

        for r in results:
            writer.writerow(r)

        buf.seek(0)

        conn = psycopg2.connect(user=self.__user, password=self.__password, host=self.__host, port=self.__port, dbname=self.__db, sslmode=self.__ssl)
        cursor = conn.cursor()
        cursor.execute('''TRUNCATE TABLE oracle_prediction''')

        cursor.copy_from(buf, 'oracle_prediction', sep=',')

        conn.commit()
