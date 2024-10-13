FROM ubuntu:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o /scheduler ./main.go

CMD ["/scheduler"]