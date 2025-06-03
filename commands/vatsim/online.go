package vatsim

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"strings"
	controllers "tpc-discord-bot/controllers/vatsim"
)

func GetOnlineMembers(s *discordgo.Session, i *discordgo.InteractionCreate) {
	controller := controllers.VATSIMController{}

	var onlineMembersBoth []*controllers.Pilot
	var onlineMembersCallsign []*controllers.Pilot
	var onlineMembersRemarks []*controllers.Pilot
	var onlineMembersNoFlightPlan []*controllers.Pilot

	var callsigns string

	data, err := controller.GetDataFeed()
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	for _, v := range data.Pilots {
		if strings.HasPrefix(v.Callsign, "TPC") {
			if v.FlightPlan != nil {
				if strings.Contains(v.FlightPlan.Remarks, "THEPILOTCLUB.ORG") {
					onlineMembersBoth = append(onlineMembersBoth, &v)
				} else {
					onlineMembersCallsign = append(onlineMembersCallsign, &v)
				}
			} else {
				onlineMembersNoFlightPlan = append(onlineMembersNoFlightPlan, &v)
			}
		}
		if v.FlightPlan != nil {
			if strings.Contains(v.FlightPlan.Remarks, "THEPILOTCLUB.ORG") {
				onlineMembersRemarks = append(onlineMembersRemarks, &v)
			}
		}
	}

	if len(onlineMembersBoth) == 0 && len(onlineMembersCallsign) == 0 && len(onlineMembersNoFlightPlan) == 0 && len(onlineMembersRemarks) == 0 {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Current Online TPC Members",
						Description: "No members are currently online 😦",
						Color:       3651327,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Made by TPC Tech Team",
							IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
						},
					},
				},
			},
		})
		return
	}

	if len(onlineMembersBoth) > 0 {
		callsigns += "**Correct Remarks with TPC Callsign:**\n"
		for _, v := range onlineMembersBoth {
			callsigns += fmt.Sprintf("- %v - %v - %v\n", v.Callsign, v.Name, v.CID)
		}
	} else {
		callsigns += "**Correct Remarks with TPC Callsign:**\n- None\n"
	}
	if len(onlineMembersCallsign) > 0 {
		callsigns += "**TPC Callsign Without Remarks Set Correctly:**\n"
		for _, v := range onlineMembersCallsign {
			callsigns += fmt.Sprintf("- %v - %v - %v\n", v.Callsign, v.Name, v.CID)
		}
	} else {
		callsigns += "**TPC Callsign Without Remarks Set Correctly:**\n- None\n"
	}
	if len(onlineMembersRemarks) > 0 {
		callsigns += "**Remarks Set Correctly:**\n"
		for _, v := range onlineMembersRemarks {
			callsigns += fmt.Sprintf("- %v - %v - %v\n", v.Callsign, v.Name, v.CID)
		}
	} else {
		callsigns += "**Remarks Set Correctly:**\n- None"
	}
	if len(onlineMembersNoFlightPlan) > 0 {
		callsigns += "**TPC Callsign With No Flight Plan on File:**\n"
		for _, v := range onlineMembersNoFlightPlan {
			callsigns += fmt.Sprintf("- %v - %v - %v\n", v.Callsign, v.Name, v.CID)
		}
	} else {
		callsigns += "**TPC Callsign With No Flight Plan on File:**\n- None\n"
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Current Online TPC Members",
					Description: callsigns,
					Color:       3651327,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Made by TPC Tech Team",
						IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
					},
				},
			},
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}
