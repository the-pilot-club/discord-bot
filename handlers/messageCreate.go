package handlers

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	eventresponses "tpc-discord-bot/event-responses"
)

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	chnl, _ := s.Channel(m.ChannelID)
	if chnl.Type == 1 || chnl.Type == 3 {
		return
	}
	switch strings.ToLower(m.Content) {
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
		eventresponses.ModeratorMessage(s, m)
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
		// booster function
	}

	//TODO: Add Auto Reactions for Screenshots
}
