package giveaway

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"math/rand"
	"time"
	"tpc-discord-bot/internal/config"
	"tpc-discord-bot/util"
)

func GiveawayMain(s *discordgo.Session, i *discordgo.InteractionCreate) {

	mem := util.FetchAllMembers(s, i.GuildID)

	var roledMembers []*discordgo.Member

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

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	for _, m := range mem {
		for _, role := range m.Roles {
			if role == config.GetRoleId(i.GuildID, "Giveaway") {
				roledMembers = append(roledMembers, m)
			}
		}
	}

	// Seed the random number generator to ensure different results each run
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random index within the bounds of the array
	randomIndex := rand.Intn(len(roledMembers))

	// Access the element at the random index
	randomElement := roledMembers[randomIndex]

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "The Pilot Club",
			IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
		},
		Description: fmt.Sprintf("And the winner is <@%v> Congratulations!", randomElement.User.ID),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Made by the TPC Tech Team",
		},
	}

	var editedEmbed []*discordgo.MessageEmbed

	editedEmbed = append(editedEmbed, embed)

	embedString := fmt.Sprintf("Congrats <@%v>!", randomElement.User.ID)

	time.Sleep(3 * time.Second)

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &embedString,
		Embeds:  &editedEmbed,
	})
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

}
