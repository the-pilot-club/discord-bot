package general

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
)

func ServerCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Create button component
	button := discordgo.Button{
		Label: "TPC Server Commands",
		Style: discordgo.LinkButton,
		URL:   "https://vats.im/tpc-commands",
	}

	// Create action row with the button
	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{button},
	}

	// Create interaction response with the button
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Here is a full list of member friendly commands:",
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
