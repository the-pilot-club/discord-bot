package fcp

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"github.com/the-pilot-club/tpcgo"
	"time"
)

func UserInfoFCP(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	o := i.ApplicationCommandData().Options[0]
	me := o.UserValue(nil)
	f, err := FCPSession()
	if err != nil {
		sentry.CaptureException(err)
		fmt.Println(err)
		return
	}
	u, ferr := f.GetFCPUser(me.ID)
	if ferr != nil {
		fmt.Println(ferr)
		if errors.As(ferr, &tpcgo.ErrNotFound) {
			var content string
			content += "Request could not be completed as dialed. Please try again later"
			_, errr := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
			})
			if errr != nil {
				sentry.CaptureException(errr)
				return
			}
			return
		} else {
			sentry.CaptureException(ferr)
			fmt.Println(ferr)
			return
		}
	}
	m, err := s.GuildMember(i.GuildID, me.ID)
	if err != nil {
		sentry.CaptureException(err)
		fmt.Println(err)
		return
	}

	embed := &discordgo.MessageEmbed{
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
		Title: fmt.Sprintf("%s's FCP Details", func() string {
			if m.Nick != "" {
				return m.Nick
			} else {
				return m.User.Username
			}
		}()),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "TPC Callsign",
				Value: fmt.Sprintf("%s", func() string {
					if u.Callsign != 0 {
						return fmt.Sprintf("TPC%v", u.Callsign)
					} else {
						return "Not Set"
					}
				}()),
			},
			{
				Name: "VATSIM CID",
				Value: fmt.Sprintf("%s", func() string {
					if u.VATSIMCid != 0 {
						return fmt.Sprintf("%v", u.VATSIMCid)
					} else {
						return "Not Set"
					}
				}()),
			},
			{
				Name: "Home Airport",
				Value: fmt.Sprintf("%s", func() string {
					if u.HomeAirport != "" {
						return fmt.Sprintf("%v", u.HomeAirport)
					} else {
						return "Not Set"
					}
				}()),
			},
			{
				Name: "Charters Code",
				Value: fmt.Sprintf("%s", func() string {
					if u.ChartersCode != "" {
						return fmt.Sprintf("%v", u.ChartersCode)
					} else {
						return "Not Set"
					}
				}()),
			},
			{
				Name: "Bio",
				Value: fmt.Sprintf("%s", func() string {
					if u.Bio != "" {
						return fmt.Sprintf("%v", u.Bio)
					} else {
						return "Not Set"
					}
				}()),
			},
			{
				Name: "Aircraft Hangar",
				Value: fmt.Sprintf("%s", func() string {
					if len(u.AircraftHangar) != 0 {
						var value string
						for _, v := range u.AircraftHangar {
							value += "- " + v.Name + "\n"
						}
						return value
					} else {
						return "No Aircraft in Hangar"
					}
				}()),
			},
		},
		Color: 3651327,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Made by TPC Tech Team",
			IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}
