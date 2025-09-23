package giveaway

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"tpc-discord-bot/internal/config"
	"tpc-discord-bot/util"
)

func ResetGiveaway(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Resting Giveaway. This command takes some time so the bot will end with this message!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	members := util.FetchAllMembers(s, i.GuildID)
	var giveawayMembers []*discordgo.Member
	for _, r := range members {
		oldroles := make(map[string]string)
		newroles := new([]string)
		for _, role := range r.Roles {
			oldroles[role] = role
		}
		roleId := config.GetRoleId(i.GuildID, "Giveaway")
		if oldroles[roleId] == roleId {
			delete(oldroles, roleId)
			for _, rl := range oldroles {
				*newroles = append(*newroles, rl)
			}
			giveawayMembers = append(giveawayMembers, &discordgo.Member{User: r.User, Roles: *newroles})
		}
	}
	for _, m := range giveawayMembers {
		_, err = s.GuildMemberEdit(i.GuildID, m.User.ID, &discordgo.GuildMemberParams{
			Roles: &m.Roles,
		})
		if err != nil {
			panic(err)
		}
	}
}
