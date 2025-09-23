package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"time"
	"tpc-discord-bot/internal/config"
)

/*
s *discordgo.Session isn't being used currently,
but added for future expansion to access disc api
*/

func OnGuildMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	f, err := FCPSession()
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	_, err = f.DeleteFCPUser(m.User.ID)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	c, err := s.Channel(config.GetChannelId(m.GuildID, "Bot Dump"))
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	_, err = s.ChannelMessageSendComplex(c.ID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Member Left",
				Author: &discordgo.MessageEmbedAuthor{
					Name:    m.User.Username,
					IconURL: m.User.AvatarURL(""),
				},
				Description: fmt.Sprintf("<@%v>", m.Member.User.ID),
				Color:       16512948,
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
		fmt.Println(err)
		return
	}
}
