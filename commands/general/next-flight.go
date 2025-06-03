package general

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"sort"
	"time"
)

func NextFlight(s *discordgo.Session, i *discordgo.InteractionCreate) {

	events, err := s.GuildScheduledEvents(i.GuildID, false)

	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	if len(events) == 0 {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "We have a problem..... There is no events posted or I am not able to retrieve them",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].ScheduledStartTime.Before(events[j].ScheduledStartTime)
	})
	nextflight := events[0]
	var desc string

	if nextflight.Description != "" {
		desc = nextflight.Description
	} else {
		desc += "Not Provided"
	}

	embed := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: nextflight.Name,
		},
		Color:       3651327,
		Description: desc,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Start Time",
				Value: nextflight.ScheduledStartTime.Format(time.RFC3339),
			},
			{
				Name:  "Voice Channel",
				Value: fmt.Sprintf("<#%v>", nextflight.ChannelID),
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: fmt.Sprintf("https://cdn.discordapp.com/guild-events/%v/%v.png?size=4096", nextflight.ID, nextflight.Image),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Made by TPC Tech Team",
			IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
		},
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Next TPC Group Flight",
			Embeds:  []*discordgo.MessageEmbed{&embed},
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}
