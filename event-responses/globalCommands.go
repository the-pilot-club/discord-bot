package event_responses

import (
	"github.com/bwmarrin/discordgo"
	"tpc-discord-bot/commands/fun"
	"tpc-discord-bot/commands/util"
)

func GlobalCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var GlobalCommandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go fun.HandlePingCommand(s, i)
		},
		"dad-joke": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go fun.HandleDadJokeCommand(s, i)
		},
		"airport": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go util.AirportCommand(s, i)
		},
		"metar": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go util.MetarCommand(s, i)
		},
	}

	if h, ok := GlobalCommandHandler[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}
