package event_responses

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"tpc-discord-bot/internal/config"
)

func TrainingRequestModal(s *discordgo.Session, i *discordgo.InteractionCreate) {

	modalData := i.ModalSubmitData()

	name := modalData.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	cid := modalData.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	crs := modalData.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	time := modalData.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

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
				Value: fmt.Sprintf("> **Full Name:** %v \n> **VATSIM CID:** %v \n> **Course In Progress:** %v \n> **Availability Today:** %v", name, cid, crs, time),
			},
		},
	}

	message := &discordgo.MessageSend{Embeds: []*discordgo.MessageEmbed{embed}}

	_, err := s.ChannelMessageSendComplex(config.GetChannelId(i.GuildID, "Training Request"), message)
	if err != nil {
		return
	}

	content := fmt.Sprintf("Thank you for submitting an Ad Hoc training request for %v. Please note, requests may or may not be honored, and are deleted every 24h.", time)
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
