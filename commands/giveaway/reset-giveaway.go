package giveaway

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"

	"tpc-discord-bot/internal/config"
	"tpc-discord-bot/util"
)

func ResetGiveaway(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	roleID := config.GetRoleId(i.GuildID, "Giveaway")
	if roleID == "" {
		editGiveawayText(s, i, "The Giveaway role is not configured for this server.")
		return
	}

	members, err := util.FetchAllMembers(s, i.GuildID)
	if err != nil {
		sentry.CaptureException(err)
		editGiveawayText(s, i, "Could not fetch server members. Please try again later.")
		return
	}

	removed, failed := 0, 0
	for _, m := range members {
		for _, role := range m.Roles {
			if role != roleID {
				continue
			}
			if err := s.GuildMemberRoleRemove(i.GuildID, m.User.ID, roleID); err != nil {
				sentry.CaptureException(err)
				failed++
			} else {
				removed++
			}
			break
		}
	}

	msg := fmt.Sprintf("Giveaway reset: removed the Giveaway role from %d member(s).", removed)
	if failed > 0 {
		msg += fmt.Sprintf(" %d removal(s) failed; check the bot's role permissions.", failed)
	}
	editGiveawayText(s, i, msg)
}
