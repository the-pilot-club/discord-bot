package fun

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
	"github.com/getsentry/sentry-go"
)

func HandleDadJokeCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var j string
	err := requests.
		URL("https://icanhazdadjoke.com/").
		Accept("text/plain").
		ToString(&j).
		Fetch(context.Background())
	if err != nil {
		fmt.Println(err)
		sentry.CaptureException(err)
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: j,
		},
	})
}
