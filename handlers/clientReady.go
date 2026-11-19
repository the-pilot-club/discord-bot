package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"log"
	"tpc-discord-bot/internal/config"
)

func HandleClientReady(s *discordgo.Session) {
	err := s.UpdateGameStatus(0, "Microsoft Flight Simulator 2024")
	if err != nil {
		sentry.CaptureException(err)
	}
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	guilds := s.State.Guilds
	for _, guild := range guilds {
		cnl := config.GetChannelId(guild.ID, "Git Channel")
		_, e := s.ChannelMessageSend(cnl, "https://tenor.com/view/b-25-pbj-1j-pbj-commmemorative-air-force-caf-socal-gif-2902325729065148155")
		if e != nil {
			sentry.CaptureException(err)
		}
	}
}
