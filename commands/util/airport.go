package util

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
	"github.com/getsentry/sentry-go"
	"tpc-discord-bot/internal/config"
)

type Airport struct {
	Icao      string `json:"icao"`
	Iata      string `json:"iata"`
	Name      string `json:"name"`
	City      string `json:"city"`
	Region    string `json:"region"`
	Country   string `json:"country"`
	Elevation string `json:"elevation_ft"`
}

type Stations struct {
	Data StationsData `json:"data"`
}

type StationsData struct {
	Callsign  string `json:"callsign"`
	Name      string `json:"name"`
	Frequency string `json:"frequency"`
	Ctaf      bool   `json:"ctaf"`
}

func getWeather(i *discordgo.InteractionCreate) string {
	options := i.ApplicationCommandData().Options
	var w string
	err := requests.
		URL("https://aviationweather.gov/api/data/metar").
		Param("ids", options[0].StringValue()).
		Accept("text/plain").
		ToString(&w).
		Fetch(context.Background())
	if err != nil {
		fmt.Println(err)
		sentry.CaptureException(err)
	}
	return w
}

func getAirportInfo(i *discordgo.InteractionCreate) ([]Airport, error) {
	var a []Airport
	options := i.ApplicationCommandData().Options

	err := requests.
		URL("https://api.api-ninjas.com/v1/airports").
		Param("icao", options[0].StringValue()).
		Accept("application/json").
		Header("X-Api-Key", config.NinjaApiKey).
		ToJSON(&a).
		Fetch(context.Background())
	return a, err
}

func getStationData(i *discordgo.InteractionCreate) (Stations, error) {
	var s Stations
	options := i.ApplicationCommandData().Options

	err := requests.
		URL("https://my.vatsim.net/api/v2/aip").
		Pathf("/airports/%v/stations", options[0].StringValue()).
		Accept("application/json").
		UserAgent("TPCDiscordBotv3").
		ToJSON(&s).
		Fetch(context.Background())
	fmt.Println(s)
	return s, err
}

func AirportCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	weather := getWeather(i)
	//airport, _ := getAirportInfo(i)
	stations, _ := getStationData(i)
	//fmt.Println(airport[0])
	fmt.Println(stations)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: weather,
		},
	})
}
