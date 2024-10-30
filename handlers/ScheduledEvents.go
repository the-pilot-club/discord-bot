package handlers

import (
	"tpc-discord-bot/internal/config"

	"github.com/bwmarrin/discordgo"
)

func HandleScheduledEventCreate(s *discordgo.Session, m *discordgo.GuildScheduledEventCreate) {
	cnl := config.GetChannelId(m.GuildID, "Admin Log Channel")
	// ensure channel found
	if cnl != "" {
		s.ChannelMessageSend(cnl, "New scheduled event created: "+m.GuildScheduledEvent.Name)
	}
}

func HandleScheduledEventUpdate(s *discordgo.Session, m *discordgo.GuildScheduledEventUpdate) {
	cnl := config.GetChannelId(m.GuildID, "Admin Log Channel")
	// ensure channel found
	if cnl != "" {
		s.ChannelMessageSend(cnl, "Scheduled event updated: "+m.GuildScheduledEvent.Name)
	}
}

func HandleScheduledEventDelete(s *discordgo.Session, m *discordgo.GuildScheduledEventDelete) {
	cnl := config.GetChannelId(m.GuildID, "Admin Log Channel")
	// ensure channel found
	if cnl != "" {
		s.ChannelMessageSend(cnl, "Scheduled event deleted: "+m.GuildScheduledEvent.Name)
	}
}
