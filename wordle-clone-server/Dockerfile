# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY src/go.mod ./
COPY src/go.sum ./

RUN go mod download

COPY src/db_client/db_client.go ./db_client/db_client.go

COPY src/main.go ./main.go

RUN go build -o /wordle-clone-server

EXPOSE 8080

CMD [ "/wordle-clone-server" ]