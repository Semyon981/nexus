FROM golang:latest

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

RUN go build -v -o /usr/local/bin/app ./cmd/main/main.go


CMD ["app"]