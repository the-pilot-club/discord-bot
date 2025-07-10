package util

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
	"github.com/getsentry/sentry-go"
	"strings"
	"time"
)

func TafCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var TAF string

	options := i.ApplicationCommandData().Options
	icao := options[0].StringValue()

	err := requests.
		URL("https://aviationweather.gov/api/data/taf").
		Param("ids", icao).
		Accept("text/plain").
		ToString(&TAF).
		CheckStatus(200).
		Fetch(context.Background())
	if err != nil {
		ierr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command could not be completed as dailed. Please try again later",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if ierr != nil {
			fmt.Println(err)
			sentry.CaptureException(err)
			return
		}
		fmt.Println(err)
		sentry.CaptureException(err)
		return
	}
	fmt.Println(TAF)
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Weather Report",
					Description: strings.ToUpper(icao),
					Color:       3651327,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name: "TAF",
							Value: fmt.Sprintf("%s", func() string {
								if TAF != "" {
									return TAF
								}
								return fmt.Sprintf("TAF is not posted for %v", strings.ToUpper(icao))
							}()),
						},
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Made by TPC Tech Team",
						IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}

}
