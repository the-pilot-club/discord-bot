FROM golang:1.24.1-alpine AS build
WORKDIR /go/src/github.com/the-pilot-club/tpc-discord-bot
COPY go.mod ./
COPY go.sum ./
COPY commands ./commands
COPY controllers ./controllers
COPY cmd ./cmd
COPY cron-jobs ./cron-jobs
COPY event-responses ./event-responses
COPY handlers ./handlers
COPY internal ./internal
COPY static-text ./static-text
COPY util ./util
RUN go build -o bin/bot ./cmd/bot
RUN go build -o bin/cron ./cmd/cron
ENTRYPOINT ["bin/bot"]
