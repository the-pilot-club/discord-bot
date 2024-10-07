package main

import (
	"github.com/getsentry/sentry-go"
	"log"
	"time"
	"tpc-discord-bot/internal/bot"
	"tpc-discord-bot/internal/config"
)

func main() {
	bot.Run()
	if config.Env != "dev" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:         config.SentryDSN,
			Debug:       false,
			Environment: config.Env,
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
		defer sentry.Flush(2 * time.Second)
	}
}
