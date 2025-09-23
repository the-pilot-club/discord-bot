package charters

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"tpc-discord-bot/internal/config"
)

func SendFerryRequestModal(s *discordgo.Session, i *discordgo.InteractionCreate) {

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "ChartersFerryRequest",
			Title:    "Ferry a TPC Charters Aircraft",
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "tail-number",
							Label:    "What is the Tail Number?",
							Style:    discordgo.TextInputShort,
							Required: true,
						},
					},
				},
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "starting-location",
							Label:    "What is your starting location?",
							Style:    discordgo.TextInputShort,
							Required: true,
						},
					},
				},
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "ending-location",
							Label:    "What is your ending location?",
							Style:    discordgo.TextInputShort,
							Required: true,
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
