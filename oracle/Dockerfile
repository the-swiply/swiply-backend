FROM python:3.11.7

ARG BASE_PATH=oracle
WORKDIR /app

COPY $BASE_PATH/requirements.txt .

RUN pip install --upgrade pip && pip install -r requirements.txt

COPY $BASE_PATH/main.py main.py
COPY $BASE_PATH/app app
COPY $BASE_PATH/configs configs


CMD ["python", "main.py"]
