package config

import "os"
import _ "github.com/joho/godotenv/autoload"

var DiscordToken = os.Getenv("BOT_TOKEN")
var SentryDSN = os.Getenv("SENTRY_DSN")
var Env = os.Getenv("GO_ENV")

func GetGitChannel() int {
	if Env == "dev" {
		return 1148307485904601172
	} else {
		return 1
	}
}
