package util

import (
	"context"
	"fmt"
	"strings"
	"tpc-discord-bot/internal/config"

	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
	"github.com/getsentry/sentry-go"
)

const commandFailureMessage = "This command could not be completed as dialed. Please try again later"

// codeOption reads the airport code option, trimmed and upper-cased.
func codeOption(i *discordgo.InteractionCreate) string {
	return strings.ToUpper(strings.TrimSpace(i.ApplicationCommandData().Options[0].StringValue()))
}

// lookupAirport queries API Ninjas by ICAO (4 chars) or IATA (3 chars).
// It returns an empty slice when no airport matches.
func lookupAirport(code string) ([]Airport, error) {
	param := "icao"
	if len(code) == 3 {
		param = "iata"
	}

	var airports []Airport
	err := requests.
		URL("https://api.api-ninjas.com/v1/airports").
		Param(param, code).
		Accept("application/json").
		Header("X-Api-Key", config.NinjaApiKey).
		ToJSON(&airports).
		CheckStatus(200).
		Fetch(context.Background())
	return airports, err
}

// resolveICAO turns an ICAO or IATA code into an ICAO code. A 4-character
// code is used as-is. A 3-character code is looked up via API Ninjas.
// When userMsg is non-empty the caller should reply with it and stop.
func resolveICAO(code string) (icao string, userMsg string, err error) {
	switch len(code) {
	case 4:
		return code, "", nil
	case 3:
		airports, err := lookupAirport(code)
		if err != nil {
			return "", "", err
		}
		if len(airports) == 0 || airports[0].Icao == "" {
			return "", fmt.Sprintf("Could not find an airport with the IATA code %s.", code), nil
		}
		return strings.ToUpper(airports[0].Icao), "", nil
	default:
		return "", invalidCodeMessage(code), nil
	}
}

func invalidCodeMessage(code string) string {
	return fmt.Sprintf("%s is not a valid code. Use a 4-letter ICAO (e.g. KJFK) or 3-letter IATA (e.g. JFK).", code)
}

func deferResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

func editResponse(s *discordgo.Session, i *discordgo.InteractionCreate, edit *discordgo.WebhookEdit) {
	if _, err := s.InteractionResponseEdit(i.Interaction, edit); err != nil {
		fmt.Println(err)
		sentry.CaptureException(err)
	}
}

func editResponseText(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	editResponse(s, i, &discordgo.WebhookEdit{Content: &msg})
}

// respondFailure logs the cause and replaces the deferred response with a generic failure message.
func respondFailure(s *discordgo.Session, i *discordgo.InteractionCreate, cause error) {
	fmt.Println(cause)
	sentry.CaptureException(cause)
	editResponseText(s, i, commandFailureMessage)
}

// deferOrLog defers the interaction and reports whether the handler should continue.
func deferOrLog(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	if err := deferResponse(s, i); err != nil {
		fmt.Println(err)
		sentry.CaptureException(err)
		return false
	}
	return true
}

// truncate shortens s to at most n runes so it fits Discord embed limits.
func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}
