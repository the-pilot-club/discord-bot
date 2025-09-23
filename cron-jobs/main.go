package cron_jobs

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"github.com/go-co-op/gocron/v2"
	"time"
	"tpc-discord-bot/cron-jobs/events"
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

	s, _ := gocron.NewScheduler()
	defer func() { _ = s.Shutdown() }()

	_, _ = s.NewJob(
		gocron.CronJob(
			"* * * * *",
			false,
		),
		gocron.NewTask(
			func() {
				go events.EventReminder(ds)
			},
		),
	)

	_, _ = s.NewJob(
		gocron.CronJob(
			"* * * * *",
			false,
		),
		gocron.NewTask(
			func() {
				t := TimeNYC()
				if t.Hour() == 8 && t.Minute() == 52 {
					//quiz answer from yesterday
				}
				if t.Hour() == 9 && t.Minute() == 0 {
					//quiz question
				}
			},
		),
	)

	s.Start()
	select {}

}
