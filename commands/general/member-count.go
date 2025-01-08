package general

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func HandleMemberCountCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild, _ := s.State.Guild(i.GuildID)
	var Message = fmt.Sprintf("Number of pilots in The Pilot Club: %v", guild.MemberCount)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: Message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
