FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY **/*.go ./

RUN go build -o scheduler .

FROM alpine:latest

WORKDIR /opt/

COPY --from=builder /app/scheduler .

COPY web ./web

CMD ["./scheduler"]
