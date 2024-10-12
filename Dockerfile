FROM ubuntu:latest AS builder

RUN apt-get update && apt-get install -y \
    golang \
    gcc \
    make \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o scheduler ./main.go

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/scheduler /app/scheduler
COPY web /app/web
COPY .env /app/.env

EXPOSE 7540

CMD ["./scheduler"]
