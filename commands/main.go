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
	GuildID        = flag.String("guild", genEnvGuild(), "Test guild ID. If not passed - bot registers commands globally")
	GlobalCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "airport",
			Description: "Displays Informatinon about the selected airport",
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
	}
	GuildCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "member-count",
			Description: "Displays Number of Members in the Club",
		},
	}
)

func init() { flag.Parse() }
