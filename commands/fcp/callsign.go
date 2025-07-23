package fcp

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"github.com/the-pilot-club/tpcgo"
	"time"
	"tpc-discord-bot/internal/config"
)

func FCPSession() (session *tpcgo.Session, err error) {

	s, errr := tpcgo.NewSession(tpcgo.SessionConfig{
		config.FCPToken,
		config.FCPEnv,
		"",
		config.CoreAPIToken,
	})
	if errr != nil {
		sentry.CaptureException(err)
		return nil, err
	}
	return s, nil

}

func GetFcpCallsign(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.ApplicationCommandData().Options
	fcp, err := FCPSession()
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	member := options[0].UserValue(nil)
	call, err := fcp.GetFCPCallsign(member.ID)
	if err != nil {
		fmt.Print(err)
		sentry.CaptureException(err)
		return
	}

	m, err := s.GuildMember(i.GuildID, member.ID)
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    m.Nick,
			IconURL: m.User.AvatarURL(""),
		},
		Description: fmt.Sprintf("TPC Callsign: %v", call.Callsign),
		Color:       3651327,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Made by TPC Tech Team",
			IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}
