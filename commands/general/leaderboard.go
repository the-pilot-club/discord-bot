package general

import (
	"github.com/bwmarrin/discordgo"
)

func HandleLeaderboardCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Create button component
	button := discordgo.Button{
		Label: "TPC Leaderboard",
		Style: discordgo.LinkButton,
		URL:   "https://mee6.xyz/thepilotclub",
	}

	// Create action row with the button
	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{button},
	}

	// Create interaction response with the button
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Check out our leaderboard!",
			Components: []discordgo.MessageComponent{actionRow},
		},
	}

	// Send the response
	s.InteractionRespond(i.Interaction, response)
}
