package handlers

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"tpc-discord-bot/internal/config"
)

func HandleCLientReady(s *discordgo.Session) {
	s.UpdateGameStatus(0, "Microsoft Flight Simulator 2020")
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	guilds := s.State.Guilds
	for _, guild := range guilds {
		cnl := config.GetChannelId(guild.ID, "Git Channel")
		s.ChannelMessageSend(cnl, "https://tenor.com/view/b-25-pbj-1j-pbj-commmemorative-air-force-caf-socal-gif-2902325729065148155")
	}
}
