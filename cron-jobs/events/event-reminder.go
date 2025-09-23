package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"log"
	"sort"
	"time"
	"tpc-discord-bot/internal/config"
)

func EventReminder(s *discordgo.Session) {
	var guildsToRemind []string
	guilds := s.State.Guilds
	for _, guild := range guilds {
		if config.EventRemindersEnabled(guild.ID) {
			guildsToRemind = append(guildsToRemind, guild.ID)
		}
	}
	for _, guild := range guildsToRemind {
		go func() {
			g, err := s.GuildScheduledEvents(guild, false)
			if len(g) == 0 {
				log.Println("No events scheduled. Skipping.")
				return
			}
			if err != nil {
				sentry.CaptureException(err)
			}
			sort.Slice(g, func(i, j int) bool {
				return g[i].ScheduledStartTime.Before(g[j].ScheduledStartTime)
			})
			ne := g[0]
			if ne.ScheduledStartTime.Format(time.DateTime) == time.Now().UTC().Add(time.Hour*1).Format(time.DateTime) {
				// send the ping for events

			}
		}()

	}

}
