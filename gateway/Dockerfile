FROM golang:latest

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

RUN go build -v -o /usr/local/bin/app ./cmd/api/main.go

EXPOSE 8080

CMD ["app"]