FROM golang:1.23.2 as build

WORKDIR /app
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=1

COPY . .

RUN go mod download

RUN go build -o /scheduler ./main.go

FROM ubuntu:latest

COPY --from=build /scheduler /scheduler

COPY web /web
COPY .env /

CMD ["/scheduler"]