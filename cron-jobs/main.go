package cron_jobs

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"github.com/go-co-op/gocron/v2"
	"time"
	"tpc-discord-bot/cron-jobs/events"
	"tpc-discord-bot/handlers"
)

func TimeNYC() time.Time {
	l, err := time.LoadLocation("America/New_York")
	if err != nil {
		sentry.CaptureException(err)
	}
	t := time.Now().In(l)
	return t
}

func HandleCronJobs(ds *discordgo.Session) {
	s, err := gocron.NewScheduler()
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	defer func() { _ = s.Shutdown() }()

	// Event reminder - runs every minute
	_, err = s.NewJob(
		gocron.CronJob("* * * * *", false),
		gocron.NewTask(func() {
			go events.EventReminder(ds)
		}),
	)
	if err != nil {
		sentry.CaptureException(err)
	}

	// Send quiz answer at 8:52 AM EST
	_, err = s.NewJob(
		gocron.CronJob("52 8 * * *", false),
		gocron.NewTask(func() {
			handlers.SendQuizAnswer(ds)
		}),
	)
	if err != nil {
		sentry.CaptureException(err)
	}

	// Send quiz question at 9:00 AM EST
	_, err = s.NewJob(
		gocron.CronJob("0 9 * * *", false),
		gocron.NewTask(func() {
			handlers.SendQuizQuestion(ds)
		}),
	)
	if err != nil {
		sentry.CaptureException(err)
	}

	s.Start()
	select {}
}
