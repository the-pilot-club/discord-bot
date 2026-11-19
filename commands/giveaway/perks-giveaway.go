package giveaway

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"

	"tpc-discord-bot/internal/config"
	"tpc-discord-bot/util"
)

func PekrsGiveaway(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !deferGiveaway(s, i) {
		return
	}

	mem, err := util.FetchAllMembers(s, i.GuildID)
	if err != nil {
		sentry.CaptureException(err)
		editGiveawayText(s, i, "Error in the command. Please reach out to the Tech Team for further troubleshooting. \n\n ERROR AREA: COULD NOT FETCH MEMBERS")
		return
	}

	// Entries per role: VIP gets 3, Frequent Flyer 2, the rest 1. Unset roles
	// resolve to "" and are skipped so they can't match anything.
	weights := map[string]int{}
	addWeight := func(name string, n int) {
		if id := config.GetRoleId(i.GuildID, name); id != "" {
			weights[id] = n
		}
	}
	addWeight("VIP", 3)
	addWeight("Frequent Flyer", 2)
	addWeight("Commuter", 1)
	addWeight("Early Adopters", 1)
	addWeight("Charters Top Gun", 1)

	var roledMembers []*discordgo.Member
	for _, m := range mem {
		for _, role := range m.Roles {
			for n := 0; n < weights[role]; n++ {
				roledMembers = append(roledMembers, m)
			}
		}
	}

	if len(roledMembers) == 0 {
		editGiveawayText(s, i, "No members have an eligible perks role, so there is no one to pick.")
		return
	}

	randomElement := roledMembers[rand.Intn(len(roledMembers))]

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

	embedString := fmt.Sprintf("Congrats <@%v>!", randomElement.User.ID)

	time.Sleep(3 * time.Second)

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &embedString,
		Embeds:  &[]*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		sentry.CaptureException(err)
	}
}
