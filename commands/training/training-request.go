package training

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
)

func TrainingRequest(s *discordgo.Session, i *discordgo.InteractionCreate) {

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "TrainingRequest",
			Title:    "Request a Training Session",
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "name",
							Label:    "Whats your Full Name?",
							Style:    discordgo.TextInputShort,
						},
					},
				},
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "cid",
							Label:    "Whats your VATSIM CID?",
							Style:    discordgo.TextInputShort,
						},
					},
				},
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "course",
							Label:    "What course are you taking?",
							Style:    discordgo.TextInputShort,
						},
					},
				},
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "time",
							Label:    "What's your availability today?",
							Style:    discordgo.TextInputShort,
						},
					},
				},
			},
		},
	}
	// Send the response
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}
