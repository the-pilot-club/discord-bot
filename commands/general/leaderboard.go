package general

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
)

func HandleLeaderboardCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Check out our leaderboard!",
			Components: leaderboardLinkComponents(),
		},
	}

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}

func leaderboardLinkComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "TPC Leaderboard",
					Style: discordgo.LinkButton,
					URL:   "https://mee6.xyz/thepilotclub",
				},
			},
		},
	}
}
