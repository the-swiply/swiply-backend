FROM python:3.11.7

ARG BASE_PATH=backoffice
WORKDIR /app

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

COPY $BASE_PATH/requirements.txt .

RUN pip install --upgrade pip && pip install -r requirements.txt

COPY $BASE_PATH/manage.py manage.py
COPY $BASE_PATH/config config
COPY $BASE_PATH/static static
COPY $BASE_PATH/apps apps


COPY $BASE_PATH/entrypoint.sh entrypoint.sh
RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
