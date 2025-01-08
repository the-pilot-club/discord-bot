package event_responses

import (
	"github.com/bwmarrin/discordgo"
	"tpc-discord-bot/commands/general"
	"tpc-discord-bot/commands/vatsim"
)

func GuildCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var GuildCommandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"member-count": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			general.HandleMemberCountCommand(s, i)
		},
		"sync": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			vatsim.SyncCommand(s, i)
		},
		"hours": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			vatsim.HoursCommand(s, i)
		},
		"leaderbaord": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			general.HandleLeaderboardCommand(s, i)
		},
	}
	if h, ok := GuildCommandHandler[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}
