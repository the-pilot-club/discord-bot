package event_responses

import "github.com/bwmarrin/discordgo"

func ModeratorMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendReply(m.ChannelID, "A Mod will be with you shortly", m.Reference())
}
