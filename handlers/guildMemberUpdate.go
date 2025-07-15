package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"strings"
	"time"
	"tpc-discord-bot/internal/config"
)

func HandleGuildMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	before := m.BeforeUpdate
	after := m
	var RoleEmbed []string
	roles := make(map[string]string)

	if len(before.Roles) > len(after.Roles) {
		for _, role := range before.Roles {
			roles[role] = role
		}
		for _, role := range after.Roles {
			delete(roles, role)
		}
		for _, role := range roles {
			proper := fmt.Sprintf("<@&%s>", role)
			RoleEmbed = append(RoleEmbed, proper)
		}
		_, err := s.ChannelMessageSendComplex(config.GetChannelId(m.GuildID, "Bot Dump"), &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Author: &discordgo.MessageEmbedAuthor{
						Name: fmt.Sprintf("%s", func() string {
							if m.Nick != "" {
								return m.Nick
							} else {
								return m.User.Username
							}
						}()),
						IconURL: m.User.AvatarURL(""),
					},
					Title:       "Role Removed",
					Description: strings.Join(RoleEmbed, " "),
					Color:       10066329,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    fmt.Sprintf("ID: %v", m.User.ID),
						IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		})
		if err != nil {
			sentry.CaptureException(err)
			return
		}
	}
	if len(before.Roles) < len(after.Roles) {
		for _, role := range after.Roles {
			roles[role] = role
		}
		for _, role := range before.Roles {
			delete(roles, role)
		}
		for _, role := range roles {
			proper := fmt.Sprintf("<@&%s>", role)
			RoleEmbed = append(RoleEmbed, proper)
		}
		_, err := s.ChannelMessageSendComplex(config.GetChannelId(m.GuildID, "Bot Dump"), &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Author: &discordgo.MessageEmbedAuthor{
						Name: fmt.Sprintf("%s", func() string {
							if m.Nick != "" {
								return m.Nick
							} else {
								return m.User.Username
							}
						}()),
						IconURL: m.User.AvatarURL(""),
					},
					Title:       "Role(s) Added",
					Description: strings.Join(RoleEmbed, " "),
					Color:       5793266,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    fmt.Sprintf("ID: %v", m.User.ID),
						IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		})
		if err != nil {
			sentry.CaptureException(err)
			return
		}
	}
	if before.Nick != after.Nick {
		_, err := s.ChannelMessageSendComplex(config.GetChannelId(m.GuildID, "Bot Dump"), &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Author: &discordgo.MessageEmbedAuthor{
						Name: fmt.Sprintf("%s", func() string {
							if m.Nick != "" {
								return m.Nick
							} else {
								return m.User.Username
							}
						}()),
						IconURL: m.User.AvatarURL(""),
					},
					Title: "Display Name Change",
					Description: fmt.Sprintf("**Before:** %s\n**After:** %s", func() string {
						if before.Nick == "" {
							return "``No Display Name Set``"
						} else {
							return before.Nick
						}

					}(), func() string {
						if after.Nick == "" {
							return "``Removed Display Name``"
						} else {
							return after.Nick
						}
					}()),
					Color: 5793266,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    fmt.Sprintf("ID: %v", m.User.ID),
						IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		})
		if err != nil {
			sentry.CaptureException(err)
			return
		}
	}
}
