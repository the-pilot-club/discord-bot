package config

import "os"
import _ "github.com/joho/godotenv/autoload"

var DiscordToken = os.Getenv("BOT_TOKEN")
