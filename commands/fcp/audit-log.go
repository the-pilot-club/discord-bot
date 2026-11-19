package fcp

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"github.com/the-pilot-club/tpcgo"
)

func AddAuditLogCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {

	o1 := i.ApplicationCommandData().Options[0]
	o2 := i.ApplicationCommandData().Options[1]
	u := o1.UserValue(nil)
	e := o2.StringValue()
	m, err := s.GuildMember(i.GuildID, u.ID)
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	fs, err := FCPSession()
	if err != nil {
		fmt.Println(err)
		sentry.CaptureException(err)
		return
	}
	re, err := fs.PostAuditLogEntry(u.ID, tpcgo.AuditLogEntry{
		UserId:      u.ID,
		SubmittedBy: i.Member.User.ID,
		Text:        e,
	})
	if err != nil {
		if errors.Is(err, tpcgo.ErrBadForm) {
			v, ferr := fs.AddFCPUser(&tpcgo.FCPUserAdd{UserID: u.ID})
			if ferr != nil {
				sentry.CaptureException(err)
				return
			}
			if v != nil {
				_, ferr := fs.PostAuditLogEntry(u.ID, tpcgo.AuditLogEntry{
					UserId:      u.ID,
					SubmittedBy: i.Member.User.ID,
					Text:        e,
				})
				if ferr != nil {
					sentry.CaptureException(err)
					return
				}
				response := &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
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
								Description: fmt.Sprintf("Audit log submitted for %s", func() string {
									if m.Nick != "" {
										return m.Nick
									} else {
										return m.User.Username
									}
								}()),
								Color: 3651327,
								Footer: &discordgo.MessageEmbedFooter{
									Text:    "Made by TPC Tech Team",
									IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
								},
								Timestamp: time.Now().Format(time.RFC3339),
							},
						},
					},
				}
				err = s.InteractionRespond(i.Interaction, response)
				if err != nil {
					sentry.CaptureException(err)
					return
				}
			}
			return
		}
		fmt.Println(err)
		sentry.CaptureException(err)
		return
	}
	fmt.Println(re)

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
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
					Description: fmt.Sprintf("Audit log submitted for %s", func() string {
						if m.Nick != "" {
							return m.Nick
						} else {
							return m.User.Username
						}
					}()),
					Color: 3651327,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Made by TPC Tech Team",
						IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		},
	}
	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}
