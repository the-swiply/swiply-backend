FROM golang:1.21 as builder

ARG BASE_PATH=chat
WORKDIR /app

COPY pkg ../pkg

COPY $BASE_PATH/go.mod .
COPY $BASE_PATH/go.sum .

RUN go mod download

COPY $BASE_PATH/cmd cmd
COPY $BASE_PATH/internal internal

ENV CGO_ENABLED=0
RUN go build -o $BASE_PATH ./cmd/$BASE_PATH/main.go

FROM alpine:3.19.0

RUN apk update

ARG BASE_PATH=chat

COPY $BASE_PATH/configs configs


COPY --from=builder /app/$BASE_PATH $BASE_PATH

CMD ["./chat"]
