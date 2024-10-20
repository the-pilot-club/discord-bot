package event_responses

import (
	"fmt"
	"tpc-discord-bot/internal/config"

	"github.com/bwmarrin/discordgo"
)

func ModeratorMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	RoleId := config.GetRoleId(m.GuildID, "Moderator")
	Message := fmt.Sprintf("A <@&%v> will be with you shortly", RoleId)
	s.ChannelMessageSendReply(m.ChannelID, Message, m.Reference())
}

func WhatIsVatsimMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	Message := "VATSIM is the Virtual Air Traffic Simulation network, connecting people from around the world flying online or acting as virtual Air Traffic Controllers.\n \n" +
		"This completely free network allows aviation enthusiasts the ultimate experience." +
		"Air Traffic Control (ATC) is available in our communities throughout the world, operating as close as possible to the real-life procedures and utilizing real-life weather, airport and route data." +
		"\n \nOn VATSIM you can join people on the other side of the planet to fly and control, with nothing more than a home computer! If you would like more information, please go to https://www.thepilotclub.org/resources#VATSIM"
	s.ChannelMessageSendReply(m.ChannelID, Message, m.Reference())
}

func TpcLiveriesMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					URL:   "https://www.thepilotclub.org/liveries",
					Label: "TPC Liveries",
					Style: discordgo.LinkButton,
				},
			},
		},
	}
	Message := &discordgo.MessageSend{
		Components: components,
		Content:    "Club liveries can be downloaded here:",
		Reference:  m.Reference(),
	}
	s.ChannelMessageSendComplex(m.ChannelID, Message)
}

func TpcCallsignMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					URL:   "https://flightcrew.thepilotclub.org",
					Label: "Get a Callsign Here!",
					Style: discordgo.LinkButton,
				},
			},
		},
	}
	embed := &discordgo.MessageEmbed{
		Title: "TPC Callsign",
		Color: 3651327,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "How to get a TPC Callsign",
				Value: "When flying group flights you get an extra 1000xp points for using a TPC callsign during the flight.",
			},
			{
				Name:  "\u200b",
				Value: "To get a TPC callsign you just need to register one that has not yet been taken. You can do so with the button below and fill in the blanks!",
			},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: "Made by TPC Dev Team"},
	}
	Message := &discordgo.MessageSend{
		Components: components,
		Embed:      embed,
		Content:    "Club liveries can be downloaded here:",
		Reference:  m.Reference(),
	}
	s.ChannelMessageSendComplex(m.ChannelID, Message)
}

func TpcThanksMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "The Pilot Club",
			IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
		},
		Color:       3651327,
		Description: "You're Welcome! Anytime!",
	}
	s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
}

func BoosterMessageContent(s *discordgo.Session, m *discordgo.MessageCreate) {
	var CrewChat = config.GetChannelId(m.GuildID, "Crew Chat")
	message := fmt.Sprintf("<@%v> Thank you for boosting the club!", m.Author.ID)
	s.ChannelMessageSend(CrewChat, message)
}

func BumpWarsMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := "**[BUMP WARS:](<discord://-/channels/830201397974663229/958549204073087086>)** \n \n" +
		"__Team 1: Hot Dogs__\n\n- Dylan | TPC1496 | DELA\n- Rich P | N7RP\n- Chris | TPC139 | ZNY" +
		"\n\n__Team 2: Big Guns__\n\n- Serge | TPC6\n- Caleb Y | TPC452\n- Kelvin | TPC1992 | SBxx" +
		" \n \n__Rules:__ \n" +
		"1: type `/bump` to bump the server on Disboard \n" +
		"2: Bumps are only possible once every 2 hours \n" +
		"3: If your bump is successful you must post the current score under your bump. Others can forfeit the point (nobody gets a point) if you don't post it until next bump. \n" +
		"4: Have fun! \n" +
		"5: This war starts on 11/02 0400z (00:00 ET) and ends on 12/01 0359z (11/30 23:59 ET) \n" +
		"\nImportant info: \n" +
		"* If there are more than one bump at a time, only those claimed will be valid, no matter how many there are. \n" +
		"* The team with the most bumps under their belt at the end of the month wins! \n" +
		"* Winning team members get 1000 TPC points and a shout-out during next town-hall \n" +
		"\nWhy are we doing this? \n\nBumping this server often helps to keep us at the top of the server list on Disboard." +
		" It gives our community a chance to grow and allows you to be involved in the process. Have fun!"

	s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
}

func FnoMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					URL:   "https://docs.google.com/document/d/1n2dorXXbRavCci0FqYMMDQngrYqnn3UXNDAiK95Kc98",
					Label: "Friday Night Operations Information",
					Style: discordgo.LinkButton,
				},
			},
		},
	}
	Message := &discordgo.MessageSend{
		Components: components,
		Content:    "FNO Stands for Friday Night Ops. You can find more information here!",
		Reference:  m.Reference(),
	}
	s.ChannelMessageSendComplex(m.ChannelID, Message)
}

func InviteLink(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendReply(m.ChannelID, "Please use this link when inviting somebody to the server: https://thepilotclub.org", m.Reference())
}

func Msfs2020Message(s *discordgo.Session, m *discordgo.MessageCreate) {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					URL:   "https://www.reddit.com/r/flightsim/wiki/msfsfaq",
					Label: "Microsoft Flight Simulator 2020 FAQ",
					Style: discordgo.LinkButton,
				},
			},
		},
	}
	Message := &discordgo.MessageSend{
		Components: components,
		Content:    "Check out MSFS2020 FAQ!",
		Reference:  m.Reference(),
	}
	s.ChannelMessageSendComplex(m.ChannelID, Message)
}

func RulesMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					URL:   "discord://-/channels/830201397974663229/833198809701679124/848232804282138644",
					Label: "The Pilot Club Rules",
					Style: discordgo.LinkButton,
				},
			},
		},
	}
	Message := &discordgo.MessageSend{
		Components: components,
		Content:    "You can find the club rules here!",
		Reference:  m.Reference(),
	}
	s.ChannelMessageSendComplex(m.ChannelID, Message)
}

func SupportMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					URL:   "https://support.thepilotclub.org/open.php",
					Label: "The Pilot Club Support",
					Style: discordgo.LinkButton,
				},
			},
		},
	}
	Message := &discordgo.MessageSend{
		Components: components,
		Content:    "To get support or submit feedback, click the button below! Thank you for being a valued member of The Pilot Club!!",
		Reference:  m.Reference(),
	}
	s.ChannelMessageSendComplex(m.ChannelID, Message)
}

func WorldTourMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					URL:   "https://discord.com/channels/830201397974663229/833198809701679124/848245312815497237",
					Label: "Get the World Tour Role",
					Style: discordgo.LinkButton,
				},
			},
		},
	}
	Message := &discordgo.MessageSend{
		Components: components,
		Content:    "Want to join the World Tour Flight? Proceed to this message and click the World Tour Logo!",
		Reference:  m.Reference(),
	}
	s.ChannelMessageSendComplex(m.ChannelID, Message)
}

func JoinVatsimMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					URL:   "https://www.thepilotclub.org/resources#VATSIM",
					Label: "Joining VATSIM",
					Style: discordgo.LinkButton,
				},
			},
		},
	}
	Message := &discordgo.MessageSend{
		Components: components,
		Content:    "To Join VATSIM you should go to this website and click register!",
		Reference:  m.Reference(),
	}
	s.ChannelMessageSendComplex(m.ChannelID, Message)
}

func WhatServerMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	Message := "We do not use the default Microsoft Flight Simulator 2020 Multiplayer Server here in The Pilot Club. We use VATSIM for all of our group flights!\n \n" +
		"VATSIM is the Virtual Air Traffic Simulation network, connecting people from around the world flying online or acting as virtual Air Traffic Controllers. " +
		"\n \nThis completely free network allows aviation enthusiasts the ultimate experience. Air Traffic Control (ATC) is available in our communities throughout the world, operating as close as possible to the real-life procedures and utilizing real-life weather, airport and route data. " +
		"\n \nOn VATSIM you can join people on the other side of the planet to fly and control, with nothing more than a home computer! If you would like more information, please go to https://vatsim.net"
	s.ChannelMessageSendReply(m.ChannelID, Message, m.Reference())
}
