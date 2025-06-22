package event_responses

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"time"
	"tpc-discord-bot/internal/config"
)

func TrainingRequestModal(s *discordgo.Session, i *discordgo.InteractionCreate) {

	modalData := i.ModalSubmitData()

	name := modalData.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	cid := modalData.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	crs := modalData.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	timee := modalData.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    i.Member.Nick,
			IconURL: i.Member.User.AvatarURL("64"),
		},
		Title: "Adhoc Training Request",
		Color: 3651327,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Details",
				Value: fmt.Sprintf("> **Full Name:** %v \n> **VATSIM CID:** %v \n> **Course In Progress:** %v \n> **Availability Today:** %v", name, cid, crs, timee),
			},
		},
	}

	message := &discordgo.MessageSend{Embeds: []*discordgo.MessageEmbed{embed}}

	_, err := s.ChannelMessageSendComplex(config.GetChannelId(i.GuildID, "Training Request"), message)
	if err != nil {
		return
	}

	content := fmt.Sprintf("Thank you for submitting an Ad Hoc training request for %v. Please note, requests may or may not be honored, and are deleted every 24h.", timee)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		return
	}
}

func ChartersAircaftRequestModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	modalData := i.ModalSubmitData()

	AirlineCode := modalData.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	AircraftType := modalData.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	SeatingConfig := modalData.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	StartLocation := modalData.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	TailNumber := modalData.Components[4].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	var Details string

	if TailNumber != "" {
		Details += fmt.Sprintf("> **TPC Charters User:** <@%v>\n> **Airline Code:** %v\n> **Aircraft Type:** %v\n> **Tail Number:** %v\n> **Seating Configuration:** %v\n> **Starting Location:** %v", i.Member.User.ID, AirlineCode, AircraftType, TailNumber, SeatingConfig, StartLocation)

	} else {
		Details += fmt.Sprintf("> **TPC Charters User:** <@%v>\n> **Airline Code:** %v\n> **Aircraft Type:** %v\n> **Seating Configuration:** %v\n> **Starting Location:** %v", i.Member.User.ID, AirlineCode, AircraftType, SeatingConfig, StartLocation)

	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    i.Member.Nick,
			IconURL: i.Member.User.AvatarURL("64"),
		},
		Title: "New Aircraft Request",
		Color: 3651327,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Aircraft Details",
				Value: Details,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Made by the TPC Tech Team",
			IconURL: "https://cdn.thepilotclub.org/discord-bot/tpc-logo.png",
		},
	}

	message := &discordgo.MessageSend{Embeds: []*discordgo.MessageEmbed{embed}}

	_, err := s.ChannelMessageSendComplex(config.GetChannelId(i.GuildID, "Charters Requests"), message)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
		return
	}

	content := fmt.Sprintf("Thank you for submitting an aircraft requests for TPC Charters. You can view your request here: <#%v>", config.GetChannelId(i.GuildID, "Charters Requests"))
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		return
	}
}
