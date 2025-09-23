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
	ModPerms       int64 = discordgo.PermissionBanMembers
	GuildID              = flag.String("guild", genEnvGuild(), "Test guild ID. If not passed - bot registers commands globally")
	GlobalCommands       = []*discordgo.ApplicationCommand{
		// fun
		{
			Name:        "ping",
			Description: "Does something cool!",
		},
		{
			Name:        "dad-joke",
			Description: "Tells you a dad joke!",
		},
		// Util
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
			Name:        "charts",
			Description: "Displays the charts of an airport.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "icao",
					Description: "The ICAO code of the airport",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        "taf",
			Description: "Displays the taf of an airport.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "icao",
					Description: "The ICAO code of the airport",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
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

		// VATSIM
		{
			Name:        "hours",
			Description: "See how many hours you have on the network!",
		},
	}
	GuildCommands = []*discordgo.ApplicationCommand{
		// Admin
		{
			Name:                     "sop-post",
			Description:              "Allows admin to update SOP and other text in the about and SOP channel.",
			DefaultMemberPermissions: &AdminPerms,
		},
		// Charters
		{
			Name:        "charters-aircraft-request",
			Description: "Use this command to request an aircraft for use in OnAir.",
		},
		{
			Name:        "charters-ferry-request",
			Description: "Use this command to request an aircraft to be ferried to another location.",
		},
		{
			Name:        "charters-join",
			Description: "Use this command if you would like to join TPC Charters",
		},
		// FCP
		{
			Name:                     "add-log",
			Description:              "Add a FCP audit log to a member",
			DefaultMemberPermissions: &ModPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "member",
					Description: "The member you wish to add an audit log for.",
					Type:        discordgo.ApplicationCommandOptionUser,
					Required:    true,
				},
				{
					Name:        "entry",
					Description: "Log Content",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        "get-callsign",
			Description: "Get a member's TPC Callsign",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "member",
					Description: "The member you wish to get the callsign of",
					Type:        discordgo.ApplicationCommandOptionUser,
					Required:    true,
				},
			},
		},
		{
			Name:        "member-info",
			Description: "Get a member's FCP Info!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "member",
					Description: "The member you wish to get the info of",
					Type:        discordgo.ApplicationCommandOptionUser,
					Required:    true,
				},
			},
		},
		{
			Name:        "fcp-link",
			Description: "The link to the Flight Crew Portal",
		},

		// General
		{
			Name:        "leaderboard",
			Description: "The link to find our leaderboard!",
		},
		{
			Name:        "member-count",
			Description: "Displays Number of Members in the Club",
		},
		{
			Name:        "next-flight",
			Description: "The link to find out our next flight!",
		},
		{
			Name:        "server-commands",
			Description: "The link to get a list of server commands!",
		},
		// Giveaway
		{
			Name:                     "giveaway",
			Description:              "Picks a random member with the Giveaway Role!",
			DefaultMemberPermissions: &AdminPerms,
			Version:                  "Dep",
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
		// Levels
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

		// Suggestions

		//Training
		{
			Name:        "training-request",
			Description: "Use this command if you would like to request training!",
		},

		// VATSIM
		{
			Name:        "get-online-members",
			Description: "Gets the members who are online",
		},
		{
			Name:        "sync",
			Description: "Sync your VATSIM Ratings for TPC!",
		},
	}
)

func init() { flag.Parse() }
