package handlers

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"github.com/the-pilot-club/tpcgo"
	"log"
	"time"
	"tpc-discord-bot/internal/config"
	"tpc-discord-bot/util"
)

func FCPSession() (s *tpcgo.FCPSession, err error) {

	s, err = tpcgo.NewFCPSession("Bearer "+config.FCPToken, config.FCPEnv)
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}
	return s, nil

}

func OnGuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	f, err := FCPSession()
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	_, err = f.AddFCPUser(&tpcgo.FCPUserAdd{UserID: m.User.ID})
	if err != nil {
		if !errors.As(err, &tpcgo.ErrAlreadyReported) {
			sentry.CaptureException(err)
			fmt.Println(err)
			return
		}
	}

	g, err := s.State.Guild(m.GuildID)
	if err != nil {
		sentry.CaptureException(err)
		fmt.Println(err)
		return
	}

	c, err := s.Channel(config.GetChannelId(m.GuildID, "Bot Dump"))
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	created, err := util.SnowflakeToTime(m.User.ID)
	if err != nil {
		sentry.CaptureException(err)
		fmt.Println(err)
		return
	}

	err = s.GuildMemberRoleAdd(m.GuildID, m.User.ID, config.GetRoleId(m.GuildID, "Pilots"))
	if err != nil {
		sentry.CaptureException(err)
		fmt.Println(err)
		return
	}
	_, err = s.ChannelMessageSendComplex(c.ID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Member Joined",
				Author: &discordgo.MessageEmbedAuthor{
					Name:    m.User.Username,
					IconURL: m.User.AvatarURL(""),
				},
				Description: fmt.Sprintf("<@%v> %vth to join\nCreated <t:%v:R>", m.Member.User.ID, g.MemberCount, created.Unix()),
				Color:       5498804,
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

	log.Printf("added user %s to FCP", m.User.ID)
}
