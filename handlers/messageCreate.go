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
		eventresponses.FnoMessage(s, m)
		break
	case "invite link":
		eventresponses.InviteLink(s, m)
		break
	case "invite link mrs bot":
		eventresponses.InviteLink(s, m)
		break
	case "moderator":
		eventresponses.ModeratorMessage(s, m)
		break
	case "msfs2020 help":
		eventresponses.Msfs2020Message(s, m)
		break
	case "rules":
		eventresponses.RulesMessage(s, m)
		break
	case "support":
		//function here
		break
	case "tpc callsign":
		eventresponses.TpcCallsignMessage(s, m)
		break
	case "tpc livery":
		eventresponses.TpcLiveriesMessage(s, m)
		break
	}

	if strings.Contains(strings.ToLower(m.Content), "join vatsim") {
		// function here
	} else if strings.Contains(strings.ToLower(m.Content), "what server") {
		// function here
	} else if strings.Contains(strings.ToLower(m.Content), "thanks tpc") {
		eventresponses.TpcThanksMessage(s, m)
	} else if strings.Contains(strings.ToLower(m.Content), "what is vatsim?") {
		eventresponses.WhatIsVatsimMessage(s, m)
	}

	if m.Type == 8 || m.Type == 9 || m.Type == 10 || m.Type == 11 {
		eventresponses.BoosterMessageContent(s, m)
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
