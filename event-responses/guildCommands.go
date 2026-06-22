package event_responses

import (
	"github.com/bwmarrin/discordgo"
	"tpc-discord-bot/commands/admin"
	"tpc-discord-bot/commands/charters"
	"tpc-discord-bot/commands/fcp"
	"tpc-discord-bot/commands/general"
	"tpc-discord-bot/commands/giveaway"
	"tpc-discord-bot/commands/suggestions"
	"tpc-discord-bot/commands/training"
	"tpc-discord-bot/commands/vatsim"
)

func GuildCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var GuildCommandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"charters-aircraft-request": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go charters.SendChartersAircraftRequestModal(s, i)
		},
		"charters-ferry-request": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go charters.SendFerryRequestModal(s, i)
		},
		"charters-join": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go charters.SendChartersJoinRequest(s, i)
		},
		"add-log": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go fcp.AddAuditLogCommand(s, i)
		},
		"fcp-link": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go fcp.SendFCPLink(s, i)
		},
		"staff-vacancies": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go fcp.SendStaffVacancies(s, i)
		},
		"get-callsign": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go fcp.GetFcpCallsign(s, i)
		},
		"member-info": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go fcp.UserInfoFCP(s, i)
		},
		"member-count": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go general.HandleMemberCountCommand(s, i)
		},
		"sync": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go vatsim.SyncCommand(s, i)
		},
		"hours": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go vatsim.HoursCommand(s, i)
		},
		"leaderboard": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go general.HandleLeaderboardCommand(s, i)
		},
		"rank": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go general.HandleRankCommand(s, i)
		},
		"get-online-members": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go vatsim.GetOnlineMembers(s, i)
		},
		"givexp": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go general.HandleGiveXpCommand(s, i)
		},
		"removexp": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go general.HandleRemoveXpCommand(s, i)
		},
		"set-user-no-xp": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go general.HandleSetUserNoXpCommand(s, i)
		},
		"giveaway": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go giveaway.GiveawayMain(s, i)
		},
		"perks-giveaway": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go giveaway.PekrsGiveaway(s, i)
		},
		"next-flight": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go general.NextFlight(s, i)
		},
		"reset-giveaway": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go giveaway.ResetGiveaway(s, i)
		},
		"server-commands": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go general.ServerCommands(s, i)
		},
		"sop-post": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go admin.SOPCommand(s, i)
		},
		"training-request": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go training.TrainingRequest(s, i)
		},
		"training-faq": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go general.HandleTrainingFAQ(s, i)
		},
		"active-threads": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go general.HandleActiveThreadsCommand(s, i)
		},
		"idea": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go suggestions.IdeaCommand(s, i)
		},
		"idea-admin": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go suggestions.IdeaAdminCommand(s, i)
		},
	}
	if h, ok := GuildCommandHandler[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}
