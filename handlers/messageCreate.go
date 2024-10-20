package handlers

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	eventresponses "tpc-discord-bot/event-responses"
	"tpc-discord-bot/internal/config"
)

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	chnl, _ := s.Channel(m.ChannelID)
	if chnl.Type == 1 || chnl.Type == 3 {
		return
	}

	switch strings.ToLower(m.Content) {
	case "bump wars":
		eventresponses.BumpWarsMessage(s, m)
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
		go eventresponses.ModeratorMessage(s, m)
		break
	case "msfs 2020 help":
		//function here
		break
	case "rules":
		// function here
		break
	case "support":
		//function here
		break
	case "tpc callsign":
		go eventresponses.TpcCallsignMessage(s, m)
		break
	case "tpc livery":
		go eventresponses.TpcLiveriesMessage(s, m)
		break
	}

	if strings.Contains(strings.ToLower(m.Content), "join vatsim") {
		// function here
	} else if strings.Contains(strings.ToLower(m.Content), "what server") {
		// function here
	} else if strings.Contains(strings.ToLower(m.Content), "thanks tpc") {
		go eventresponses.TpcThanksMessage(s, m)
	} else if strings.Contains(strings.ToLower(m.Content), "what is vatsim?") {
		go eventresponses.WhatIsVatsimMessage(s, m)
	}

	if m.Type == 8 || m.Type == 9 || m.Type == 10 || m.Type == 11 {
		go eventresponses.BoosterMessageContent(s, m)
	}

	var Contest = config.GetChannelId(m.GuildID, "Screenshot Contest")

	if m.ChannelID == Contest || strings.Contains(chnl.Name, "SCREENSHOT CONTEST") || chnl.ParentID == Contest {
		var ToReact bool
		for i := 0; i < len(m.Attachments); i++ {
			if strings.Contains(m.Attachments[i].ContentType, "image") {
				ToReact = true
			}
		}
		if ToReact == true {
			emoji := config.GetEmojiId(m.GuildID, "TPC Reaction")
			s.MessageReactionAdd(m.ChannelID, m.ID, emoji)
		}
		return
	}
	return
}
