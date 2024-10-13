package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"tpc-discord-bot/internal/config"
)

func HandleCLientReady(s *discordgo.Session) {
	var Id string
	if config.Env == "dev" {
		Id = "1148307481085358190"
	} else {
		Id = "830201397974663229"
	}
	cnl := config.GetChannelId(Id, "Git Channel")
	fmt.Println(cnl)
	s.UpdateGameStatus(0, "Microsoft Flight Simulator 2020")
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	s.ChannelMessageSend(cnl, "https://tenor.com/view/b-25-pbj-1j-pbj-commmemorative-air-force-caf-socal-gif-2902325729065148155")
}
