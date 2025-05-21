package giveaway

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"tpc-discord-bot/util"
)

func GiveawayMain(s *discordgo.Session, i *discordgo.InteractionCreate) {

	mem := util.FetchAllMembers(s, i.GuildID)

	if len(mem) == 0 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error in the command. Please reach out to the Tech Team for further troubleshooting. \n\n ERROR AREA: NO MEMBERS FOUND",
			},
		})
		if err != nil {
			sentry.CaptureException(err)
			panic(err)
		}
	}

	for _, m := range mem {
		fmt.Println(m.User.Username)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "I did it!",
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

}
