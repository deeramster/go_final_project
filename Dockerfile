FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o scheduler ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/scheduler .

COPY web ./web
COPY .env ./.env

CMD ["./scheduler"]
