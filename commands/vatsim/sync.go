package vatsim

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
	"github.com/getsentry/sentry-go"
)

type CidResponse struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
}

type V2Response struct {
	Rating      int `json:"rating"`
	PilotRating int `json:"pilot_rating"`
}

func getCID(i *discordgo.InteractionCreate) (CidResponse, error) {

	var r CidResponse
	err := requests.
		URL("https://api.vatsim.net").
		Pathf("/v2/members/discord/%s", "0").
		CheckStatus(200).
		ToJSON(&r).
		Fetch(context.Background())
	return r, err
}

func getRatings(c string) (V2Response, error) {
	var u V2Response
	err := requests.
		URL("https://api.vatsim.net").
		Pathf("/v2/members/discord/%v", c).
		ToJSON(&u).
		Fetch(context.Background())
	return u, err

}

func SyncCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, err := getCID(i)
	if requests.HasStatusErr(err, 404) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please connect your VATSIM account to the VATSIM Community Hub!",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								URL:   "https://community.vatsim.net",
								Label: "Connect my Account!",
								Style: discordgo.LinkButton,
							},
						},
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	} else if err != nil {
		sentry.CaptureException(err)
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
