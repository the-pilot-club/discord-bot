package util

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
)

func getWeatherMetar(icao string) (string, error) {
	var weather string
	err := requests.
		URL("https://aviationweather.gov/api/data/metar").
		Param("ids", icao).
		Accept("text/plain").
		ToString(&weather).
		CheckStatus(200, 204).
		Fetch(context.Background())
	return strings.TrimSpace(weather), err
}

func MetarCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	weather, err := getWeatherMetar(icao)
	if err != nil {
		respondFailure(s, i, err)
		return
	}
	if weather == "" {
		weather = "Not Available"
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Weather Report",
		Description: icao,
		Color:       3651327,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "METAR",
				Value: truncate(weather, 1024),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Made by TPC Dev Team",
			IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
		},
	}

	editResponse(s, i, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
}
