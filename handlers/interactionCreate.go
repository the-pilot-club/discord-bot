package handlers

import (
	"github.com/bwmarrin/discordgo"
	"tpc-discord-bot/commands/fun"
)

func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var CommandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			fun.HandlePingCommand(s, i)
		},
	}

	if h, ok := CommandHandler[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}
