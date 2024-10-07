package handlers

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	switch m.Content {
	case "bump wars":
		// function here
		break
	case "what is fno?":
		//function here
		break
	case "invite link":
		//function here
		break
	case "invite mrs bot":
		//function here
		break
	case "moderator":
		// function here
		break
	case "msfs 2020 help":
		//function here
		break

	case "help":
		//function here
		break

	}

	if strings.Contains(strings.ToLower(m.Content), "moderator") {
		log.Println(m.MessageReference)
		_, err := s.ChannelMessageSendReply(m.ChannelID, "Hello I am Here to Help!", m.Reference())
		if err != nil {
			log.Println(err)
		}
	}
}
