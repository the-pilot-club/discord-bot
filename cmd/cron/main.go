package main

import (
	"fmt"
	"os"

	"github.com/getsentry/sentry-go"

	"tpc-discord-bot/cron-jobs/airac"
	"tpc-discord-bot/cron-jobs/events"
	"tpc-discord-bot/cron-jobs/quiz"
	"tpc-discord-bot/internal/bootstrap"
)

const usage = "usage: cron <event-reminder|airac-reminder|quiz-question|quiz-answer>"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	job := os.Args[1]
	if !isKnownJob(job) {
		fmt.Fprintf(os.Stderr, "unknown job: %s\n", job)
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	session, cleanup, err := bootstrap.InitForCron()
	if err != nil {
		fmt.Fprintf(os.Stderr, "bootstrap failed: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	var jobErr error
	switch job {
	case "event-reminder":
		jobErr = events.EventReminder(session)
	case "airac-reminder":
		jobErr = airac.AiracReminder(session)
	case "quiz-question":
		jobErr = quiz.SendQuizQuestion(session)
	case "quiz-answer":
		jobErr = quiz.SendQuizAnswer(session)
	}
	if jobErr != nil {
		sentry.CaptureException(jobErr)
		fmt.Fprintf(os.Stderr, "%s failed: %v\n", job, jobErr)
		cleanup()
		os.Exit(1)
	}
}

func isKnownJob(job string) bool {
	switch job {
	case "event-reminder", "airac-reminder", "quiz-question", "quiz-answer":
		return true
	default:
		return false
	}
}
