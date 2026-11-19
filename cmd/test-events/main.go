// Command test-events fires synthetic MessageCreate events against the
// configured Discord server. It is intended for local/dev testing only.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"tpc-discord-bot/handlers"
	"tpc-discord-bot/internal/config"
)

func main() {
	if strings.TrimSpace(config.DiscordToken) == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	if len(os.Args) < 2 {
		usage()
	}

	guildID := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	userID := guildID.String("user", "999999999999999999", "synthetic user ID")
	messageID := guildID.String("message", "", "existing message ID (required for screenshot)")
	guildID.Parse(os.Args[2:])

	session, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		log.Fatal(err)
	}

	devGuild := config.AllGuildIDs()
	if len(devGuild) == 0 {
		log.Fatal("no guild configuration loaded; set CONFIG_PATH")
	}

	switch os.Args[1] {
	case "boost":
		channelID := config.GetChannelId(devGuild[0], "Crew Chat")
		fireBoost(session, devGuild[0], channelID, *userID)
	case "screenshot":
		if *messageID == "" {
			log.Fatal("screenshot requires -message <existing-message-id>")
		}
		channelID := config.GetChannelId(devGuild[0], "Screenshot Contest")
		fireScreenshot(session, devGuild[0], channelID, *messageID, *userID)
	default:
		usage()
	}

	// Handlers are registered asynchronously by the bot, so give the REST
	// request time to complete before the short-lived CLI exits.
	time.Sleep(2 * time.Second)
}

func fireBoost(s *discordgo.Session, guildID, channelID, userID string) {
	if channelID == "" {
		log.Fatal("Crew Chat is not configured for the selected guild")
	}
	handlers.MessageCreateHandler(s, &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ID:        "synthetic-boost-message",
			GuildID:   guildID,
			ChannelID: channelID,
			Type:      discordgo.MessageType(8),
			Author:    &discordgo.User{ID: userID, Username: "Synthetic Booster"},
		},
	})
	fmt.Println("Synthetic boost event fired")
}

func fireScreenshot(s *discordgo.Session, guildID, channelID, messageID, userID string) {
	if channelID == "" {
		log.Fatal("Screenshot Contest is not configured for the selected guild")
	}
	handlers.MessageCreateHandler(s, &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ID:        messageID,
			GuildID:   guildID,
			ChannelID: channelID,
			Author:    &discordgo.User{ID: userID, Username: "Synthetic Photographer"},
			Attachments: []*discordgo.MessageAttachment{{
				ContentType: "image/png",
				URL:         "https://example.com/test.png",
			}},
		},
	})
	fmt.Println("Synthetic screenshot event fired")
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/test-events boost [-user USER_ID]")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/test-events screenshot -message MESSAGE_ID [-user USER_ID]")
	os.Exit(2)
}
