package admin

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"os"
	"tpc-discord-bot/internal/config"
	static_text "tpc-discord-bot/static-text"
)

func SOPCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	md, _ := os.ReadFile("./static-text/sop.md")

	embed := &discordgo.MessageEmbed{
		Description: string(md),
		Color:       3651327,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Made by TPC Tech Team",
			IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
		},
	}

	_, err := s.ChannelMessageSendComplex(config.GetChannelId(i.GuildID, "SOP"), &discordgo.MessageSend{
		Embed: embed,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: "onboarded",
						Label:    "Complete Onboarding",
						Style:    discordgo.SecondaryButton,
						Emoji: &discordgo.ComponentEmoji{
							ID:   "967794040592535622",
							Name: "onboarded",
						},
					},
				},
			},
		},
	})

	rolesEmbed := &discordgo.MessageEmbed{
		Description: static_text.SetSOPRolesText(i.GuildID),
		Color:       3651327,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Made by TPC Tech Team",
			IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
		},
	}

	_, err = s.ChannelMessageSendComplex(config.GetChannelId(i.GuildID, "SOP"), &discordgo.MessageSend{
		Embed: rolesEmbed,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: "giveaway",
						Label:    "Enter Giveaway",
						Style:    discordgo.SecondaryButton,
						Emoji: &discordgo.ComponentEmoji{
							ID:   "895480872243978280",
							Name: "giveaway",
						},
					},
					discordgo.Button{
						CustomID: "ga-pilots",
						Label:    "GA Crew",
						Style:    discordgo.SecondaryButton,
						Emoji: &discordgo.ComponentEmoji{
							ID:   "898533440532672512",
							Name: "GA_Gang",
						},
					},
					discordgo.Button{
						CustomID: "group-flights",
						Label:    "Group Flights",
						Style:    discordgo.SecondaryButton,
						Emoji: &discordgo.ComponentEmoji{
							ID:   "938262268972511373",
							Name: "group_flights",
						},
					},
					discordgo.Button{
						CustomID: "streamers",
						Label:    "Streamers",
						Style:    discordgo.SecondaryButton,
						Emoji: &discordgo.ComponentEmoji{
							ID:   "895480920960819211",
							Name: "streaming",
						},
					},
					discordgo.Button{
						CustomID: "insiders",
						Label:    "Daily Briefing",
						Style:    discordgo.SecondaryButton,
						Emoji: &discordgo.ComponentEmoji{
							ID:   "996214296503652383",
							Name: "Insiders",
						},
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: "flight-school",
						Label:    "TPC Flight School Student",
						Style:    discordgo.SecondaryButton,
						Emoji: &discordgo.ComponentEmoji{
							ID:   "895480894901592074",
							Name: "training",
						},
					},
					discordgo.Button{
						CustomID: "heli",
						Label:    "Helicopter Pilot",
						Style:    discordgo.SecondaryButton,
						Emoji: &discordgo.ComponentEmoji{
							ID:   "1167835494869110784",
							Name: "helicopter",
						},
					},
					discordgo.Button{
						CustomID: "workshops",
						Label:    "TPC Workshops",
						Style:    discordgo.SecondaryButton,
						Emoji: &discordgo.ComponentEmoji{
							ID:   "1241894339546972242",
							Name: "workshops",
						},
					},
				},
			},
		},
	})

	_, err = s.ChannelMessageSendComplex(config.GetChannelId(i.GuildID, "SOP"), &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Description: "📑 Pick your Simulator\n\nSelect your simulator:\n\n<:SimFS2020:1008378435304947792> MSFS 2020/2024\n<:Simxplane:1008374576763371520> X-Plane 11/12\n<:SimP3D:1008374664302702622> Prepar3D\n<:SimFSX:1008374722507047042> Flight Simulator X\n\n",
			Color:       3651327,
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Made by TPC Tech Team",
				IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: "msfs",
						Label:    "MSFS 2020/2024",
						Style:    discordgo.SecondaryButton,
						Emoji: &discordgo.ComponentEmoji{
							ID:   "1008378435304947792",
							Name: "SimFS2020",
						},
					},
					discordgo.Button{
						CustomID: "xplane",
						Label:    "XPLANE 11/12",
						Style:    discordgo.SecondaryButton,
						Emoji: &discordgo.ComponentEmoji{
							ID:   "1008374576763371520",
							Name: "Simxplane",
						},
					},
					discordgo.Button{
						CustomID: "p3d",
						Label:    "P3D",
						Style:    discordgo.SecondaryButton,
						Emoji: &discordgo.ComponentEmoji{
							ID:   "1008374664302702622",
							Name: "SimP3D",
						},
					},
					discordgo.Button{
						CustomID: "fsx",
						Label:    "FSX",
						Style:    discordgo.SecondaryButton,
						Emoji: &discordgo.ComponentEmoji{
							ID:   "1008374722507047042",
							Name: "SimFSX",
						},
					},
				},
			},
		},
	})

	if err != nil {
		sentry.CaptureException(err)
		panic(err)
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "SOP Message Sent",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
		return
	}

}
