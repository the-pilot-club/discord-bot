package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func HandleScheduledEventCreate(s *discordgo.Session, m *discordgo.GuildScheduledEventCreate) {
	// Add your event handling logic here
	log.Printf("New scheduled event created: %s", m.GuildScheduledEvent.Name)
}

func HandleScheduledEventUpdate(s *discordgo.Session, m *discordgo.GuildScheduledEventUpdate) {
	log.Printf("Scheduled event updated: %s", m.GuildScheduledEvent.Name)
}

func HandleScheduledEventDelete(s *discordgo.Session, m *discordgo.GuildScheduledEventDelete) {
	log.Printf("Scheduled event deleted: %s", m.GuildScheduledEvent.Name)
}
