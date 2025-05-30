package event_responses

import (
	"tpc-discord-bot/commands/general"
	"tpc-discord-bot/commands/giveaway"
	"tpc-discord-bot/commands/vatsim"

	"github.com/bwmarrin/discordgo"
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
		"leaderboard": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			general.HandleLeaderboardCommand(s, i)
		},
		"givexp": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			general.HandleGiveXpCommand(s, i)
		},
		"giveaway": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			giveaway.GiveawayMain(s, i)
		},
		"perks-giveaway": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			giveaway.PekrsGiveaway(s, i)
		},
	}
	if h, ok := GuildCommandHandler[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}
