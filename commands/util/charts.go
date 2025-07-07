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

func ChartsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {

	o := i.ApplicationCommandData().Options[0]
	icao := o.StringValue()

	url := fmt.Sprintf("https://metar.vatsim.net/%v", icao)

	var j string
	err := requests.
		URL(url).
		Accept("text/plain").
		ToString(&j).
		Fetch(context.Background())
	if err != nil {
		fmt.Println(err)
		sentry.CaptureException(err)
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Airport",
					Description: fmt.Sprintf("Information about %v\n\n[Charts (SkyVector)](https://skyvector.com/api/airportSearch?query=%v)\n**METAR**\n%v", strings.ToUpper(icao), icao, j),
					Color:       3651327,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Made by TPC Tech Team",
						IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Content: j,
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}
