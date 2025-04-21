package general

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"

	controllers "tpc-discord-bot/controllers/levels"
	"tpc-discord-bot/internal/leveling"
)

func HandleGiveXpCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get command options
	options := i.ApplicationCommandData().Options
	targetUser := options[0].UserValue(s)
	xpAmount := options[1].IntValue()

	// Skip if target is a bot
	if targetUser.Bot {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Cannot give XP to a bot!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			sentry.CaptureException(err)
		}
		return
	}

	controller := &controllers.LeaderboardController{}

	// if user not found, tell the user that the user if not found and return
	user, err := controller.FindUser(targetUser.ID, i.GuildID)
	if err != nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You are kind, but the user not found!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			sentry.CaptureException(err)
		}
		return
	}

	// Check if user has noXp flag
	if noXp, ok := user["noXp"].(bool); ok && noXp {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This user has XP disabled!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			sentry.CaptureException(err)
		}
		return
	}

	// Get current user stats
	currentLevel := int(user["level"].(float64))
	currentXp := int(user["xp"].(float64))
	totalXp := int(user["totalXp"].(float64))
	messageCount := int(user["messageCount"].(float64))

	// Calculate new level
	newLevel, requiredXp := leveling.CalculateUserLevel(currentLevel, currentXp+int(xpAmount))

	if newLevel > currentLevel {
		// Level up
		content := fmt.Sprintf("Congrats <@%v>, you just advanced to TPC **level %v **!", targetUser.ID, newLevel)
		_, err = s.ChannelMessageSend(i.ChannelID, content)
		if err != nil {
			sentry.CaptureException(err)
			return
		}

		err = controller.UpdateUserLevel(targetUser.ID, newLevel, messageCount+1, 1, requiredXp, i.GuildID)
		if err != nil {
			sentry.CaptureException(err)
			return
		}

		// Check and assign role rewards
		leveling.CheckRoleRewards(s, i.Interaction.Member.GuildID, targetUser.ID, controller, newLevel)
		return
	}

	// Update points
	err = controller.UpdateUserPoints(
		targetUser.ID,
		newLevel,
		messageCount+1,
		currentXp+int(xpAmount),
		totalXp+int(xpAmount),
		requiredXp,
		time.Now().Add(time.Minute).UnixMilli(),
		i.GuildID,
	)
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	// Send success response
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Successfully gave %d XP to <@%s>!", xpAmount, targetUser.ID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}
