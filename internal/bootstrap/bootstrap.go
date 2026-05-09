package bootstrap

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"

	"tpc-discord-bot/internal/cache"
	"tpc-discord-bot/internal/config"
)

// InitForCron initializes Sentry, Redis, and a REST-only Discord session.
// The returned cleanup function must be called before the process exits.
func InitForCron() (*discordgo.Session, func(), error) {
	cleanup := initSharedServices()

	session, err := newSession()
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	return session, cleanup, nil
}

// InitForBot initializes Sentry, Redis, and an opened gateway Discord session.
// registerHandlers is called before the gateway session is opened.
func InitForBot(registerHandlers func(*discordgo.Session)) (*discordgo.Session, func(), error) {
	cleanupShared := initSharedServices()

	session, err := newSession()
	if err != nil {
		cleanupShared()
		return nil, nil, err
	}

	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	if registerHandlers != nil {
		registerHandlers(session)
	}

	if err := session.Open(); err != nil {
		sentry.CaptureException(err)
		cleanupShared()
		return nil, nil, err
	}

	cleanup := func() {
		if err := session.Close(); err != nil {
			sentry.CaptureException(err)
		}
		cleanupShared()
	}

	return session, cleanup, nil
}

func initSharedServices() func() {
	if config.Env == "production" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:         config.SentryDSN,
			Debug:       false,
			Environment: config.Env,
		}); err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
	}

	cache.InitCache()

	return func() {
		cache.CloseCache()
		if config.Env == "production" {
			sentry.Flush(2 * time.Second)
		}
	}
}

func newSession() (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}
	return session, nil
}
