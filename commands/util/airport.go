package util

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
	"github.com/getsentry/sentry-go"
	"strings"
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
	StationData []StationsData `json:"data"`
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

func getStationData(i *discordgo.InteractionCreate) ([]StationsData, error) {
	var s Stations
	options := i.ApplicationCommandData().Options

	err := requests.
		URL("https://my.vatsim.net").
		Pathf("/api/v2/aip/airports/%s/stations", options[0].StringValue()).
		Accept("application/json").
		UserAgent("TPCDiscordBotv3").
		ToJSON(&s).
		Fetch(context.Background())

	return s.StationData, err
}

func getFreq(s []StationsData) (*discordgo.MessageEmbedField, error) {

	var freq string

	for _, f := range s {
		if f.Ctaf == true {
			var ctaf = fmt.Sprintf("- %v (%v): **CTAF Frequency: Yes**\n  - %v\n", f.Name, f.Callsign, f.Frequency)
			freq += ctaf
		} else {
			var nonCtaf = fmt.Sprintf("- %v (%v)\n  - %v\n", f.Name, f.Callsign, f.Frequency)
			freq += nonCtaf
		}
	}

	fi := &discordgo.MessageEmbedField{
		Name:  "Frequencies",
		Value: freq,
	}

	return fi, nil
}

func AirportCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	weather := getWeather(i)
	airport, err := getAirportInfo(i)
	if err != nil || len(airport) == 0 {
		msg := fmt.Sprintf("%v has no airport information provided. [SkyVector may be able to provide more information.](https://skyvector.com/api/airportSearch?query=%v)", strings.ToUpper(options[0].StringValue()), options[0].StringValue())
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		return
	}
	stations, _ := getStationData(i)

	des := fmt.Sprintf("Information about %s (Elevation: %s)", airport[0].Name, airport[0].Elevation)

	freqs, ferr := getFreq(stations)

	var Embed *discordgo.MessageEmbed

	var IataF *discordgo.MessageEmbedField

	if len(airport[0].Iata) > 0 {
		IataF = &discordgo.MessageEmbedField{
			Name:  "IATA",
			Value: airport[0].Iata,
		}
	} else {
		IataF = &discordgo.MessageEmbedField{
			Name:  "IATA",
			Value: "N/A",
		}
	}

	if ferr == nil {
		Embed = &discordgo.MessageEmbed{
			Title:       "Airport",
			Description: des,
			Color:       3651327,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "ICAO",
					Value: airport[0].Icao,
				},
				IataF,
				{
					Name:  "Region",
					Value: airport[0].Region,
				},
				{
					Name:  "Charts",
					Value: fmt.Sprintf("[View Charts Here](https://skyvector.com/api/airportSearch?query=%v)", options[0].StringValue()),
				},
				{
					Name:  "METAR",
					Value: weather,
				},
				freqs,
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Made by TPC Dev Team",
				IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
			},
		}
	}
	if ferr != nil {
		Embed = &discordgo.MessageEmbed{
			Title:       "Airport",
			Description: des,
			Color:       3651327,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "ICAO",
					Value: airport[0].Icao,
				},
				IataF,
				{
					Name:  "Region",
					Value: airport[0].Region,
				},
				{
					Name:  "Charts",
					Value: fmt.Sprintf("[View Charts Here](https://skyvector.com/api/airportSearch?query=%v)", options[0].StringValue()),
				},
				{
					Name:  "METAR",
					Value: weather,
				},
				{
					Name:  "Frequencies",
					Value: "None Found on VATSIM",
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Made by TPC Dev Team",
				IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
			},
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				Embed,
			},
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		fmt.Println(err)
	}
}
