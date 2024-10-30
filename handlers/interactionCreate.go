package handlers

import (
	"github.com/bwmarrin/discordgo"
	"tpc-discord-bot/commands/fun"
	"tpc-discord-bot/commands/general"
)

func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var GlobalCommandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			fun.HandlePingCommand(s, i)
		},
	}

	if h, ok := GlobalCommandHandler[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}

	var GuildCommandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"member-count": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			general.HandleMemberCountCommand(s, i)
		},
	}
	if h, ok := GuildCommandHandler[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}

}
