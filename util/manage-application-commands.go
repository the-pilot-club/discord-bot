package util

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"log"
	"tpc-discord-bot/commands"
)

func HandleApplicationCommandUpdates(session *discordgo.Session) {
	registeredGlobalCommands := make([]*discordgo.ApplicationCommand, len(commands.GlobalCommands))
	registeredGuildCommands := make([]*discordgo.ApplicationCommand, len(commands.GuildCommands))
	log.Println("Updating Commands")
	for i, v := range commands.GlobalCommands {
		registeredGlobalCommands[i] = v
	}
	for i, v := range commands.GuildCommands {
		registeredGuildCommands[i] = v
	}
	_, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, "", registeredGlobalCommands)
	if err != nil {
		sentry.CaptureException(err)
		println(err.Error())
	}
	_, err = session.ApplicationCommandBulkOverwrite(session.State.User.ID, *commands.GuildID, registeredGuildCommands)
	if err != nil {
		sentry.CaptureException(err)
		println(err.Error())
	}
}
