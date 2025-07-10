package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"log"
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
	v, err := f.DeleteFCPUser(m.User.ID)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	fmt.Println(v)
	fmt.Println()
	log.Printf("removed user %s from FCP", m.User.ID)
}
