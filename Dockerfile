FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY **/*.go ./

RUN go build -o scheduler .

FROM alpine:latest

WORKDIR /opt/

COPY --from=builder /app/scheduler .

ARG TODO_PORT
ARG TODO_DBFILE
ARG TODO_PASSWORD
ENV TODO_PORT=${{PORT}}
ENV TODO_DBFILE=${{DB_NAME}}
ENV TODO_PASSWORD=${{PASSWORD}}


COPY web ./web

CMD ["./scheduler"]
