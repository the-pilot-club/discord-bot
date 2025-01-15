package vatsim

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
	"github.com/getsentry/sentry-go"
	"strings"
	"tpc-discord-bot/internal/config"
)

type CidResponse struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
}

type V2Response struct {
	Rating      int `json:"rating"`
	PilotRating int `json:"pilotrating"`
}

var Cid CidResponse
var Ratings V2Response

func getCID(s *discordgo.Session, i *discordgo.InteractionCreate) {

	err := requests.
		URL("https://api.vatsim.net").
		Pathf("/v2/members/discord/%s", i.Member.User.ID).
		CheckStatus(200).
		ToJSON(&Cid).
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

func getRatings(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := requests.
		URL("https://api.vatsim.net").
		Pathf("/v2/members/%s", Cid.UserId).
		ToJSON(&Ratings).
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

func SyncCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	getCID(s, i)
	getRatings(s, i)

	ratingRoles := config.GetRatingsRoles(i.GuildID)
	pilotRatingRoles := config.GetPilotRatingsRoles(i.GuildID)
	atcRole := config.GetRoleId(i.GuildID, "ATC")
	var rolesEmbed []string
	newRoles := new([]string)

	CurrentRoles := make(map[string]string)

	for _, v := range i.Member.Roles {
		CurrentRoles[v] = v
	}
	delete(CurrentRoles, atcRole)
	for _, v := range ratingRoles {
		if CurrentRoles[v.Id] == v.Id {
			delete(CurrentRoles, v.Id)
		}
		if Ratings.Rating == v.RatingValue {
			CurrentRoles[v.Id] = v.Id
			role := fmt.Sprintf("<@&%s>", v.Id)
			rolesEmbed = append(rolesEmbed, role)
		}
	}
	for _, v := range pilotRatingRoles {
		if CurrentRoles[v.Id] == v.Id {
			delete(CurrentRoles, v.Id)
		}
		if Ratings.PilotRating == v.RatingValue {
			CurrentRoles[v.Id] = v.Id
			role := fmt.Sprintf("<@&%s>", v.Id)
			rolesEmbed = append(rolesEmbed, role)
		}
	}
	if Ratings.Rating > 1 && Ratings.Rating < 12 {
		CurrentRoles[atcRole] = atcRole
		role := fmt.Sprintf("<@&%s>", atcRole)
		rolesEmbed = append(rolesEmbed, role)
	}

	for k := range CurrentRoles {
		*newRoles = append(*newRoles, k)
	}

	_, err := s.GuildMemberEdit(i.GuildID, i.Member.User.ID, &discordgo.GuildMemberParams{
		Roles: newRoles,
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	embed := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    i.Member.Nick,
			IconURL: i.Member.User.AvatarURL(""),
		},
		Title:       "Your Roles Have Been Assigned!",
		Description: strings.Join(rolesEmbed, " "),
		Color:       3651327,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Made by TPC Dev Team",
			IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
		},
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				&embed,
			},
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}
