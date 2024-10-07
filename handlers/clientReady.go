package handlers

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"tpc-discord-bot/internal/config"
)

func HandleCLientReady(s *discordgo.Session) {
	cnl := config.GetGitChannel()
	s.UpdateGameStatus(0, "Microsoft Flight Simulator 2020")
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	s.ChannelMessageSend(strconv.Itoa(cnl), "https://tenor.com/view/b-25-pbj-1j-pbj-commmemorative-air-force-caf-socal-gif-2902325729065148155")
}
