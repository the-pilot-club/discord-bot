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
