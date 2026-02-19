package fcp

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
)

func SendStaffVacancies(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Create button component
	button := discordgo.Button{
		Label: "Staff Vacancies",
		Style: discordgo.LinkButton,
		URL:   "https://www.thepilotclub.org/staff-vacancies",
	}

	// Create action row with the button
	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{button},
	}

	// Create interaction response with the button
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Here is a link to Staff Vacancies:",
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