package general

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"

	controllers "tpc-discord-bot/controllers/levels"
	"tpc-discord-bot/internal/leveling"
)

func HandleSetUserNoXpCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	targetUser, enabled, err := noXpCommandOptions(s, i)
	if err != nil {
		editNoXpResponse(s, i, err.Error())
		return
	}

	if targetUser.Bot {
		editNoXpResponse(s, i, "Bots do not have leaderboard XP.")
		return
	}

	controller := &controllers.LeaderboardController{}
	userData, err := controller.FindUser(targetUser.ID, i.GuildID)
	userExists := err == nil

	if err != nil && !leaderboardUserNotFound(err) {
		sentry.CaptureException(err)
		editNoXpResponse(s, i, "I couldn't fetch leaderboard data right now. Try again later.")
		return
	}

	if !userExists {
		if !enabled {
			editNoXpResponse(s, i, fmt.Sprintf("<@%s> does not have leaderboard data yet, so there is no no-XP flag to clear.", targetUser.ID))
			return
		}

		initialProgress := leveling.ProgressFromTotalXp(0)
		createErr := controller.CreateUserRecord(controllers.UserCreate{
			GuildID:         i.GuildID,
			UserID:          targetUser.ID,
			MessageCount:    0,
			Xp:              initialProgress.CurrentXp,
			TotalXp:         initialProgress.TotalXp,
			LevelXp:         initialProgress.NextLevelXp,
			Level:           initialProgress.Level,
			Rank:            0,
			NoXp:            true,
			MessageLastSent: 0,
		}, i.GuildID)
		if createErr != nil {
			sentry.CaptureException(createErr)
			editNoXpResponse(s, i, "I couldn't create leaderboard data for that member.")
			return
		}

		editNoXpResponse(s, i, fmt.Sprintf("Marked <@%s> as no-XP. A new leaderboard record was created with XP gain disabled.", targetUser.ID))
		return
	}

	currentNoXp := leaderboardOptionalBool(userData, "noXp")
	if currentNoXp == enabled {
		editNoXpResponse(s, i, formatNoXpNoChangeResponse(targetUser.ID, enabled))
		return
	}

	err = controller.NoUserXp(targetUser.ID, enabled, i.GuildID)
	if err != nil {
		sentry.CaptureException(err)
		editNoXpResponse(s, i, "I couldn't update leaderboard data right now. Try again later.")
		return
	}

	editNoXpResponse(s, i, formatNoXpUpdatedResponse(targetUser.ID, enabled))
}

func noXpCommandOptions(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.User, bool, error) {
	var targetUser *discordgo.User
	var enabled bool
	var enabledSet bool

	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "user":
			targetUser = option.UserValue(s)
		case "enabled":
			enabled = option.BoolValue()
			enabledSet = true
		}
	}

	if targetUser == nil {
		return nil, false, fmt.Errorf("unable to resolve the requested member")
	}

	if !enabledSet {
		return nil, false, fmt.Errorf("unable to read the requested no-XP state")
	}

	return targetUser, enabled, nil
}

func editNoXpResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	if err != nil {
		sentry.CaptureException(err)
	}
}

func formatNoXpNoChangeResponse(userID string, enabled bool) string {
	if enabled {
		return fmt.Sprintf("<@%s> is already marked as no-XP.", userID)
	}

	return fmt.Sprintf("<@%s> is already eligible to earn XP.", userID)
}

func formatNoXpUpdatedResponse(userID string, enabled bool) string {
	if enabled {
		return fmt.Sprintf("Marked <@%s> as no-XP. Their current leaderboard progress was kept, but future message XP is now disabled.", userID)
	}

	return fmt.Sprintf("Removed the no-XP flag from <@%s>. They can earn message XP again.", userID)
}
