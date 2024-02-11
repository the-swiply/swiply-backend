FROM python:3.12

ARG BASE_PATH=oracle
WORKDIR /app

COPY $BASE_PATH/requirements.txt .

RUN pip install --upgrade pip && pip install -r requirements.txt

COPY $BASE_PATH/app app
COPY $BASE_PATH/configs configs


CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "80"]
