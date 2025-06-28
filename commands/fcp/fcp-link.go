package fcp

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
)

func SendFCPLink(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Create button component
	button := discordgo.Button{
		Label: "Flight Crew Portal",
		Style: discordgo.LinkButton,
		URL:   "https://flightcrew.thepilotclub.org/",
	}

	// Create action row with the button
	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{button},
	}

	// Create interaction response with the button
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Here is a link to the Flight Crew Portal:",
			Components: []discordgo.MessageComponent{actionRow},
		},
	}

	// Send the response
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}
