package vatsim

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
	"github.com/getsentry/sentry-go"
)

type Hours struct {
	Id    int     `json:"id"`
	ATC   float64 `json:"atc"`
	Pilot float64 `json:"pilot"`
	S1    float64 `json:"s1"`
	S2    float64 `json:"s2"`
	S3    float64 `json:"s3"`
	C1    float64 `json:"c1"`
	C3    float64 `json:"c3"`
	I1    float64 `json:"i1"`
	I3    float64 `json:"i3"`
	SUP   float64 `json:"sup"`
	ADM   float64 `json:"adm"`
}

var CID CidResponse
var HoursData Hours

func getCid(s *discordgo.Session, i *discordgo.InteractionCreate) {

	err := requests.
		URL("https://api.vatsim.net").
		Pathf("/v2/members/discord/%s", i.Member.User.ID).
		CheckStatus(200).
		ToJSON(&CID).
		Fetch(context.Background())

	if requests.HasStatusErr(err, 404) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please connect your VATSIM account to the VATSIM Community Hub!",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								URL:   "https://community.vatsim.net",
								Label: "Connect my Account!",
								Style: discordgo.LinkButton,
							},
						},
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	} else if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This operation could not be completed. Please contact a developer via the support channel. \n\n reference: sync-discord-id",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		sentry.CaptureException(err)
		return
	}
	return
}

func getRatingsHours(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := requests.
		URL("https://api.vatsim.net").
		Pathf("/v2/members/%s/stats", CID.UserId).
		ToJSON(&HoursData).
		CheckStatus(200).
		Fetch(context.Background())

	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This operation could not be completed. Please contact a developer via the support channel. \n\n reference: sync-vatsim-ratings",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		fmt.Println(err)
		sentry.CaptureException(err)
		return
	}
	return

}

func HoursCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	getCid(s, i)
	getRatingsHours(s, i)

	pilotHours := discordgo.MessageEmbedField{
		Name:  "Pilot Hours",
		Value: fmt.Sprintf("%v", HoursData.Pilot),
	}

	atcHours := discordgo.MessageEmbedField{
		Name:  "ATC Hours",
		Value: fmt.Sprintf("%v", HoursData.ATC),
	}

	supHours := discordgo.MessageEmbedField{
		Name:  "Supervisor / Administrator Hours",
		Value: fmt.Sprintf("%v", HoursData.SUP+HoursData.ADM),
	}

	if HoursData.S1 != 0 {

	}

	embed := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    i.Member.Nick,
			IconURL: i.Member.User.AvatarURL(""),
		},
		Title: "Your hours on VATSIM!",
		Color: 3651327,
		Fields: []*discordgo.MessageEmbedField{
			&pilotHours,
			&atcHours,
			&supHours,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Made by TPC Dev Team",
			IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				&embed,
			},
		},
	})

}
