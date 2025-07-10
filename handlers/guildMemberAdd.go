package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"github.com/the-pilot-club/tpcgo"
	"log"
	"tpc-discord-bot/internal/config"
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
	v, err := f.AddFCPUser(&tpcgo.FCPUserAdd{UserID: m.User.ID})
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	fmt.Println(v)

	log.Printf("added user %s to FCP", m.User.ID)
}
