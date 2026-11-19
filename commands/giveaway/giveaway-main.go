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

func GiveawayMain(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Defer first: fetching every member can take longer than Discord's
	// 3-second window for an initial response.
	if !deferGiveaway(s, i) {
		return
	}

	mem, err := util.FetchAllMembers(s, i.GuildID)
	if err != nil {
		sentry.CaptureException(err)
		editGiveawayText(s, i, "Error in the command. Please reach out to the Tech Team for further troubleshooting. \n\n ERROR AREA: COULD NOT FETCH MEMBERS")
		return
	}

	roleID := config.GetRoleId(i.GuildID, "Giveaway")
	var roledMembers []*discordgo.Member
	if roleID != "" {
		for _, m := range mem {
			for _, role := range m.Roles {
				if role == roleID {
					roledMembers = append(roledMembers, m)
				}
			}
		}
	}

	if len(roledMembers) == 0 {
		editGiveawayText(s, i, "No members have the Giveaway role, so there is no one to pick.")
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

// deferGiveaway sends the "thinking..." response and reports whether the
// handler should continue.
func deferGiveaway(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})
	if err != nil {
		sentry.CaptureException(err)
		return false
	}
	return true
}

func editGiveawayText(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	}); err != nil {
		sentry.CaptureException(err)
	}
}
