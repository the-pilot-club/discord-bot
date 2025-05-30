package util

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
)

func FetchAllMembers(s *discordgo.Session, guildId string) []*discordgo.Member {

	data := make(chan *discordgo.GuildMembersChunk)
	var completedMC *discordgo.GuildMembersChunk
	var members []*discordgo.Member
	err := s.RequestGuildMembers(guildId, "", 0, "", false)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	s.AddHandlerOnce(func(s *discordgo.Session, mc *discordgo.GuildMembersChunk) {
		completedMC = mc
		data <- completedMC
		close(data)
	})

	finshedData := <-data

	for _, m := range finshedData.Members {
		members = append(members, m)
	}

	return members
}
