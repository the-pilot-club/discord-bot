package commands

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"tpc-discord-bot/internal/config"
)

func genEnvGuild() string {
	if config.Env == "dev" {
		return "1148307481085358190"
	} else if config.Env == "prod" {
		return "830201397974663229"
	}
	return ""
}

var (
	AdminPerms     int64 = discordgo.PermissionAdministrator
	StaffPerms     int64 = discordgo.PermissionMentionEveryone
	GuildID              = flag.String("guild", genEnvGuild(), "Test guild ID. If not passed - bot registers commands globally")
	GlobalCommands       = []*discordgo.ApplicationCommand{
		{
			Name:        "airport",
			Description: "Displays Information about the selected airport",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "icao",
					Description: "What is the ICAO of the Airport?",
					Required:    true,
					MaxLength:   4,
				},
			},
		},
		{
			Name:        "ping",
			Description: "Does something cool!",
		},
		{
			Name:        "dad-joke",
			Description: "Tells you a dad joke!",
		},
		{
			Name:        "hours",
			Description: "See how many hours you have on the network!",
		},
		{
			Name:        "metar",
			Description: "Gives METAR for a Specific Airport",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "icao-code",
					Description: "use an ICAO to get a METAR request",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "icao",
							Description: "The ICAO of the Airport",
							Required:    true,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "iata-code",
					Description: "use an IATA to get a METAR request",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "iata",
							Description: "The IATA of the Airport",
							Required:    true,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
	}
	GuildCommands = []*discordgo.ApplicationCommand{
		{
			Name:                     "giveaway",
			Description:              "Picks a random member with the Giveaway Role!",
			DefaultMemberPermissions: &AdminPerms,
			Version:                  "Dep",
		},
		{
			Name:                     "sop-post",
			Description:              "Allows admin to update SOP and other text in the about and SOP channel.",
			DefaultMemberPermissions: &AdminPerms,
		},
		{
			Name:                     "perks-giveaway",
			Description:              "Picks a random member with the Company Perks Role(s)!",
			DefaultMemberPermissions: &AdminPerms,
		},
		{
			Name:        "reset-giveaway",
			Description: "Removes the giveaway roles from the users who have it",
		},
		{
			Name:        "get-online-members",
			Description: "Gets the members who are online",
		},
		{
			Name:        "server-commands",
			Description: "The link to get a list of server commands!",
		},
		{
			Name:        "training-request",
			Description: "Use this command if you would like to request training!",
		},
		{
			Name:        "member-count",
			Description: "Displays Number of Members in the Club",
		},
		{
			Name:        "sync",
			Description: "Sync your VATSIM Ratings for TPC!",
		},
		{
			Name:        "leaderboard",
			Description: "The link to find our leaderboard!",
		},
		{
			Name:                     "givexp",
			Description:              "Give XP to a user",
			DefaultMemberPermissions: &StaffPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to give XP to",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "amount",
					Description: "Amount of XP to give as a whole number - no partial gimmies here!",
					Required:    true,
				},
			},
		},
	}
)

func init() { flag.Parse() }
