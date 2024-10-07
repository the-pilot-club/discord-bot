FROM golang:1.23.1-alpine3.20 AS build
WORKDIR /go/src/github.com/the-pilot-club/tpc-discord-bot
COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY functions ./functions
COPY internal ./internal
RUN go build -o bin/bot ./cmd/main.go