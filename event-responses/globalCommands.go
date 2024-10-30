package event_responses

import (
	"github.com/bwmarrin/discordgo"
	"tpc-discord-bot/commands/fun"
)

func GlobalCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var GlobalCommandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			fun.HandlePingCommand(s, i)
		},
		"dad-joke": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			fun.HandleDadJokeCommand(s, i)
		},
	}

	if h, ok := GlobalCommandHandler[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}
