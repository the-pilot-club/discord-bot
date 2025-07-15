package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"time"
	"tpc-discord-bot/internal/config"
)

func HandleGuildBanAdd(s *discordgo.Session, m *discordgo.GuildBanAdd) {
	_, err := s.ChannelMessageSendComplex(config.GetChannelId(m.GuildID, "Bot Dump"), &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Author: &discordgo.MessageEmbedAuthor{
					Name:    fmt.Sprintf("%v", m.User.Username),
					IconURL: m.User.AvatarURL(""),
				},
				Title:       "Member Banned",
				Description: fmt.Sprintf("<@%v>", m.User.ID),
				Color:       14115665,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    fmt.Sprintf("ID: %v", m.User.ID),
					IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
				},
				Timestamp: time.Now().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}
