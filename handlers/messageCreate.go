package handlers

import (
	"strings"
	eventresponses "tpc-discord-bot/event-responses"
	"tpc-discord-bot/internal/config"

	"github.com/bwmarrin/discordgo"
)

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	go HandleXpGive(s, m)

	chnl, _ := s.Channel(m.ChannelID)
	if chnl.Type == 1 || chnl.Type == 3 {
		return
	}

	switch strings.ToLower(m.Content) {
	case "bump wars":
		go eventresponses.BumpWarsMessage(s, m)
		return
	case "what is fno?":
		go eventresponses.FnoMessage(s, m)
		return
	case "invite link":
		go eventresponses.InviteLink(s, m)
		return
	case "invite link mrs bot":
		go eventresponses.InviteLink(s, m)
		return
	case "moderator":
		go eventresponses.ModeratorMessage(s, m)
		return
	case "msfs2020 help":
		go eventresponses.Msfs2020Message(s, m)
		return
	case "rules":
		go eventresponses.RulesMessage(s, m)
		return
	case "support":
		go eventresponses.SupportMessage(s, m)
		return
	case "tpc callsign":
		go eventresponses.TpcCallsignMessage(s, m)
		return
	case "tpc livery":
		go eventresponses.TpcLiveriesMessage(s, m)
		return
	case "world tour":
		go eventresponses.WorldTourMessage(s, m)
		return
	}

	if strings.Contains(strings.ToLower(m.Content), "join vatsim") {
		go eventresponses.JoinVatsimMessage(s, m)
	} else if strings.Contains(strings.ToLower(m.Content), "what server") {
		go eventresponses.WhatServerMessage(s, m)
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
		if ToReact {
			emoji := config.GetEmojiId(m.GuildID, "TPC Reaction")
			s.MessageReactionAdd(m.ChannelID, m.ID, emoji)
		}
		return
	}

}
