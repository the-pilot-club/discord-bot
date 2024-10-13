package event_responses

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"tpc-discord-bot/internal/config"
)

func ModeratorMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	RoleId := config.GetRoleId(m.GuildID, "Moderator")
	Message := fmt.Sprintf("A <@&%v> will be with you shortly", RoleId)
	s.ChannelMessageSendReply(m.ChannelID, Message, m.Reference())
}
