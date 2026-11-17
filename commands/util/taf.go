package util

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
)

func TafCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	code := codeOption(i)

	if !deferOrLog(s, i) {
		return
	}

	icao, userMsg, err := resolveICAO(code)
	if err != nil {
		respondFailure(s, i, err)
		return
	}
	if userMsg != "" {
		editResponseText(s, i, userMsg)
		return
	}

	var taf string
	err = requests.
		URL("https://aviationweather.gov/api/data/taf").
		Param("ids", icao).
		Accept("text/plain").
		ToString(&taf).
		CheckStatus(200, 204).
		Fetch(context.Background())
	if err != nil {
		respondFailure(s, i, err)
		return
	}

	taf = strings.TrimSpace(taf)
	if taf == "" {
		taf = fmt.Sprintf("TAF is not posted for %v", icao)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Weather Report",
		Description: icao,
		Color:       3651327,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "TAF",
				Value: truncate(taf, 1024),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Made by TPC Tech Team",
			IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	editResponse(s, i, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
}
