package event_responses

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"tpc-discord-bot/internal/config"
)

func ModeratorMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	RoleId := config.GetRoleId(m.GuildID, "Moderator")
	Message := fmt.Sprintf("A <@&%v> will be with you shortly", RoleId)
	s.ChannelMessageSendReply(m.ChannelID, Message, m.Reference())
}

func WhatIsVatsimMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	Message := "VATSIM is the Virtual Air Traffic Simulation network, connecting people from around the world flying online or acting as virtual Air Traffic Controllers.\n \n" +
		"This completely free network allows aviation enthusiasts the ultimate experience." +
		"Air Traffic Control (ATC) is available in our communities throughout the world, operating as close as possible to the real-life procedures and utilizing real-life weather, airport and route data." +
		"\n \nOn VATSIM you can join people on the other side of the planet to fly and control, with nothing more than a home computer! If you would like more information, please go to https://www.thepilotclub.org/resources#VATSIM"
	s.ChannelMessageSendReply(m.ChannelID, Message, m.Reference())
}
func TpcLiveriesMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					URL:   "https://www.thepilotclub.org/liveries",
					Label: "TPC Liveries",
					Style: discordgo.LinkButton,
				},
			},
		},
	}
	Message := &discordgo.MessageSend{
		Components: components,
		Content:    "Club liveries can be downloaded here:",
		Reference:  m.Reference(),
	}
	s.ChannelMessageSendComplex(m.ChannelID, Message)
}

func TpcCallsignMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					URL:   "https://flightcrew.thepilotclub.org",
					Label: "Get a Callsign Here!",
					Style: discordgo.LinkButton,
				},
			},
		},
	}
	embed := &discordgo.MessageEmbed{
		Title: "TPC Callsign",
		Color: 3651327,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "How to get a TPC Callsign",
				Value: "When flying group flights you get an extra 1000xp points for using a TPC callsign during the flight.",
			},
			{
				Name:  "\u200b",
				Value: "To get a TPC callsign you just need to register one that has not yet been taken. You can do so with the button below and fill in the blanks!",
			},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: "Made by TPC Dev Team"},
	}
	Message := &discordgo.MessageSend{
		Components: components,
		Embed:      embed,
		Content:    "Club liveries can be downloaded here:",
		Reference:  m.Reference(),
	}
	s.ChannelMessageSendComplex(m.ChannelID, Message)
}

func TpcThanksMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "The Pilot Club",
			IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
		},
		Color:       3651327,
		Description: "You're Welcome! Anytime!",
	}
	s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
}
