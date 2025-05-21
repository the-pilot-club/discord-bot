package util

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"time"
)

var GMC *discordgo.GuildMembersChunk

func processGMChunk(s *discordgo.Session, mc *discordgo.GuildMembersChunk) {
	GMC = mc
}

func FetchAllMembers(s *discordgo.Session, guildId string) []*discordgo.Member {

	var members []*discordgo.Member
	err := s.RequestGuildMembers(guildId, "", 0, "", false)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	s.AddHandlerOnce(processGMChunk)

	time.Sleep(500 * time.Millisecond)

	for _, m := range GMC.Members {
		members = append(members, m)
	}

	return members
}
