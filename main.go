package main

import (
	"log"
	"os"
	"tpc-discord-bot/internal/bot"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get Discord token from environment variables
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("No Discord token provided")
	}

	// Create and start the bot
	client := bot.NewClient(token)
	client.Start()
}
