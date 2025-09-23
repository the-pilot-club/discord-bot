package event_responses

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"tpc-discord-bot/internal/config"
)

func HandleButtonSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	buttonData := i.MessageComponentData()

	switch buttonData.CustomID {
	case "onboarded":
		onboardedRole(s, i)
	case "giveaway":
		selfAssignedRoles(s, i, "Giveaway")
	case "ga-pilots":
		selfAssignedRoles(s, i, "GA Flights")
	case "group-flights":
		selfAssignedRoles(s, i, "Group Flights")
	case "streamers":
		selfAssignedRoles(s, i, "Streamers")
	case "insiders":
		selfAssignedRoles(s, i, "Insiders")
	case "flight-school":
		selfAssignedRoles(s, i, "Flight School")
	case "heli":
		selfAssignedRoles(s, i, "Helicopter")
	case "workshops":
		selfAssignedRoles(s, i, "Workshops")
	case "msfs":
		selfAssignedRoles(s, i, "FS2020/24")
	case "xplane":
		selfAssignedRoles(s, i, "X-Plane")
	case "p3d":
		selfAssignedRoles(s, i, "P3D")
	case "fsx":
		selfAssignedRoles(s, i, "FSX")
	}
}

func onboardedRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	oldroles := make(map[string]string)
	newRoles := new([]string)
	for _, r := range i.Member.Roles {
		oldroles[r] = r
		*newRoles = append(*newRoles, r)
	}

	roleId := config.GetRoleId(i.GuildID, "Onboarded")
	if oldroles[roleId] == "" {
		*newRoles = append(*newRoles, roleId)
	} else {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You have already completed your onboarding steps!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		return
	}

	_, err := s.GuildMemberEdit(i.GuildID, i.Member.User.ID, &discordgo.GuildMemberParams{
		Roles: newRoles,
	})
	if err != nil {
		panic(err)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Thanks for completing the onboarding! You have been granted the role!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}

func selfAssignedRoles(s *discordgo.Session, i *discordgo.InteractionCreate, roleName string) {
	oldroles := make(map[string]string)
	newRoles := new([]string)
	var isAssigned bool
	roleId := config.GetRoleId(i.GuildID, roleName)
	for _, r := range i.Member.Roles {
		oldroles[r] = r
	}

	if oldroles[roleId] == "" {
		for _, v := range oldroles {
			*newRoles = append(*newRoles, v)
		}
		*newRoles = append(*newRoles, roleId)
		isAssigned = true
	} else {
		for _, v := range oldroles {
			if v != roleId {
				*newRoles = append(*newRoles, v)
			}
		}
		isAssigned = false
	}

	_, err := s.GuildMemberEdit(i.GuildID, i.Member.User.ID, &discordgo.GuildMemberParams{
		Roles: newRoles,
	})
	if err != nil {
		panic(err)
	}

	var content string

	if isAssigned {
		content = fmt.Sprintf("You have been asssigned the %v role!", roleName)
	} else {
		content = fmt.Sprintf("The %v role has been removed!", roleName)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}
