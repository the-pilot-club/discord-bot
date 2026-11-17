package util

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func ChartsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	metar, err := getWeatherMetar(icao)
	if err != nil || metar == "" {
		metar = "Not Available"
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Airport",
		Description: fmt.Sprintf("Information about %v\n\n[Charts (SkyVector)](https://skyvector.com/api/airportSearch?query=%v)\n**METAR**\n%v", icao, icao, metar),
		Color:       3651327,
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
