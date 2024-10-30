package event_responses

import (
	"github.com/bwmarrin/discordgo"
	"tpc-discord-bot/commands/general"
)

func GuildCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var GuildCommandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"member-count": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			general.HandleMemberCountCommand(s, i)
		},
	}
	if h, ok := GuildCommandHandler[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}
