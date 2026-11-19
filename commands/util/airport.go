package util

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
)

type Airport struct {
	Icao      string      `json:"icao"`
	Iata      string      `json:"iata"`
	Name      string      `json:"name"`
	City      string      `json:"city"`
	Region    string      `json:"region"`
	Country   string      `json:"country"`
	Elevation json.Number `json:"elevation_ft"`
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

func getStationData(icao string) ([]StationsData, error) {
	var s Stations
	err := requests.
		URL("https://my.vatsim.net").
		Pathf("/api/v2/aip/airports/%s/stations", icao).
		Accept("application/json").
		UserAgent("TPCDiscordBotv3").
		ToJSON(&s).
		Fetch(context.Background())

	return s.StationData, err
}

func getFreq(s []StationsData) *discordgo.MessageEmbedField {
	var freq string

	for _, f := range s {
		if f.Ctaf {
			freq += fmt.Sprintf("- %v (%v): **CTAF Frequency: Yes**\n  - %v\n", f.Name, f.Callsign, f.Frequency)
		} else {
			freq += fmt.Sprintf("- %v (%v)\n  - %v\n", f.Name, f.Callsign, f.Frequency)
		}
	}

	if freq == "" {
		freq = "None Found on VATSIM"
	}

	return &discordgo.MessageEmbedField{
		Name:  "Frequencies",
		Value: truncate(freq, 1024),
	}
}

func AirportCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	code := codeOption(i)

	if !deferOrLog(s, i) {
		return
	}

	if len(code) != 3 && len(code) != 4 {
		editResponseText(s, i, invalidCodeMessage(code))
		return
	}

	airport, err := lookupAirport(code)
	if err != nil {
		respondFailure(s, i, err)
		return
	}
	if len(airport) == 0 || airport[0].Icao == "" {
		msg := fmt.Sprintf("%v has no airport information provided. [SkyVector may be able to provide more information.](https://skyvector.com/api/airportSearch?query=%v)", code, code)
		editResponseText(s, i, msg)
		return
	}

	icao := strings.ToUpper(airport[0].Icao)

	weather, werr := getWeatherMetar(icao)
	if werr != nil || weather == "" {
		weather = "Not Available"
	}

	var freqs *discordgo.MessageEmbedField
	if stations, serr := getStationData(icao); serr != nil {
		freqs = &discordgo.MessageEmbedField{Name: "Frequencies", Value: "None Found on VATSIM"}
	} else {
		freqs = getFreq(stations)
	}

	iata := airport[0].Iata
	if iata == "" {
		iata = "N/A"
	}
	region := airport[0].Region
	if region == "" {
		region = "N/A"
	}
	elevation := "N/A"
	if airport[0].Elevation != "" {
		elevation = airport[0].Elevation.String() + " ft"
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Airport",
		Description: fmt.Sprintf("Information about %s (Elevation: %s)", airport[0].Name, elevation),
		Color:       3651327,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "ICAO", Value: icao},
			{Name: "IATA", Value: iata},
			{Name: "Region", Value: region},
			{
				Name:  "Charts",
				Value: fmt.Sprintf("[View Charts Here](https://skyvector.com/api/airportSearch?query=%v)", icao),
			},
			{Name: "METAR", Value: truncate(weather, 1024)},
			freqs,
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
