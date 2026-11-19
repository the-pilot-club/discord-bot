package util

import (
	"github.com/bwmarrin/discordgo"
)

// FetchAllMembers returns every member of the guild, paging through the REST
// API 1000 members at a time.
//
// It previously asked the gateway for member chunks, which only ever read the
// first chunk (so guilds over 1000 members were truncated), could miss the
// reply entirely because the handler was registered after the request, and
// then blocked forever.
func FetchAllMembers(s *discordgo.Session, guildID string) ([]*discordgo.Member, error) {
	var members []*discordgo.Member
	after := ""
	for {
		page, err := s.GuildMembers(guildID, after, 1000)
		if err != nil {
			return nil, err
		}
		members = append(members, page...)
		if len(page) < 1000 {
			return members, nil
		}
		after = page[len(page)-1].User.ID
	}
}
