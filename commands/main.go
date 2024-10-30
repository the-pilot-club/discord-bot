package commands

import (
	"flag"
	"github.com/bwmarrin/discordgo"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	GlobalCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Does something cool!",
		},
	}
	GuildCommands = []*discordgo.ApplicationCommand{}
)

func init() { flag.Parse() }
