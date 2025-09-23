package charters

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"tpc-discord-bot/internal/config"
)

func SendChartersAircraftRequestModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "ChartersAircraftRequest",
			Title:    "Request an Aircraft for TPC Charters",
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "airline-code",
							Label:    "Whats your Airline code on OnAir?",
							Style:    discordgo.TextInputShort,
							Required: true,
						},
					},
				},
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "aircraft-type",
							Label:    "What is the aircraft type you would like?",
							Style:    discordgo.TextInputShort,
							Required: true,
						},
					},
				},
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "seating-config",
							Label:    "What is the preferred seating configuration?",
							Style:    discordgo.TextInputShort,
							Required: true,
						},
					},
				},
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "starting-location",
							Label:    "Where would you like the aircraft to be?",
							Style:    discordgo.TextInputShort,
							Required: true,
						},
					},
				},
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "tail-number",
							Label:    "What is the tail number?",
							Style:    discordgo.TextInputShort,
							Required: false,
						},
					},
				},
			},
		},
	}
	rolesmap := make(map[string]string)

	for _, v := range i.Member.Roles {
		rolesmap[v] = v
	}

	if rolesmap[config.GetRoleId(i.GuildID, "Charters Pilots")] != config.GetRoleId(i.GuildID, "Charters Pilots") {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You do not have the Charters Pilot Roles to complete this command",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		return
	}
	// Send the response
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
		return
	}

}
