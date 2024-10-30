package handlers

import (
	"github.com/bwmarrin/discordgo"
	eventresponses "tpc-discord-bot/event-responses"
)

func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	eventresponses.GuildCommands(s, i)
	eventresponses.GlobalCommands(s, i)

}
