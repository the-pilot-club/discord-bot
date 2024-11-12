package util

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
	"github.com/getsentry/sentry-go"
	"strings"
)

var Weather string

func getWeatherMetar(m string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := requests.
		URL("https://aviationweather.gov/api/data/metar").
		Param("ids", m).
		Accept("text/plain").
		ToString(&Weather).
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
}

func MetarCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var AirportCode string

	options := i.ApplicationCommandData().Options

	if options[0].Name == "icao-code" {
		code := options[0].Options
		getWeatherMetar(code[0].StringValue(), s, i)
		if Weather == "" {
			Weather = "Not Available"
		}
		AirportCode = strings.ToUpper(code[0].StringValue())

	}

	if options[0].Name == "iata-code" {

	}

	embed := &discordgo.MessageEmbed{
		Title:       "Weather Report",
		Description: AirportCode,
		Color:       3651327,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "METAR",
				Value: Weather,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Made by TPC Dev Team",
			IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
		},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				embed,
			},
		},
	})
	if err != nil {
		fmt.Println("From Bottom Function")
		fmt.Println(err)
	}
}
