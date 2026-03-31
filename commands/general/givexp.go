package general

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"

	controllers "tpc-discord-bot/controllers/levels"
	"tpc-discord-bot/internal/leveling"
)

func HandleGiveXpCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	handleXpAdjustmentCommand(s, i, 1)
}

func HandleRemoveXpCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	handleXpAdjustmentCommand(s, i, -1)
}

func handleXpAdjustmentCommand(s *discordgo.Session, i *discordgo.InteractionCreate, direction int) {
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

	targetUser, xpAmount, err := xpAdjustmentOptions(s, i)
	if err != nil {
		editXpAdjustmentResponse(s, i, err.Error())
		return
	}

	if targetUser.Bot {
		editXpAdjustmentResponse(s, i, "Bots do not have leaderboard XP.")
		return
	}

	delta := direction * xpAmount
	controller := &controllers.LeaderboardController{}
	userData, err := controller.FindUser(targetUser.ID, i.GuildID)
	userExists := err == nil

	if err != nil && !leaderboardUserNotFound(err) {
		sentry.CaptureException(err)
		editXpAdjustmentResponse(s, i, "I couldn't fetch leaderboard data right now. Try again later.")
		return
	}

	if !userExists && direction < 0 {
		editXpAdjustmentResponse(s, i, fmt.Sprintf("<@%s> does not have leaderboard data yet, so there is no XP to remove.", targetUser.ID))
		return
	}

	stats := &leaderboardStats{}
	if userExists {
		stats, err = parseLeaderboardStats(userData)
		if err != nil {
			sentry.CaptureException(err)
			editXpAdjustmentResponse(s, i, "Leaderboard data for that member could not be read.")
			return
		}
	}

	change := leveling.ApplyXpDelta(stats.Level, stats.CurrentXp, stats.TotalXp, delta)

	if !userExists {
		createErr := controller.CreateUserRecord(controllers.UserCreate{
			GuildID:         i.GuildID,
			UserID:          targetUser.ID,
			MessageCount:    0,
			Xp:              change.After.CurrentXp,
			TotalXp:         change.After.TotalXp,
			LevelXp:         change.After.NextLevelXp,
			Level:           change.After.Level,
			Rank:            0,
			NoXp:            false,
			MessageLastSent: 0,
		}, i.GuildID)
		if createErr != nil {
			sentry.CaptureException(createErr)
			editXpAdjustmentResponse(s, i, "I couldn't create leaderboard data for that member.")
			return
		}
	} else {
		err = controller.UpdateUserXpState(
			targetUser.ID,
			change.After.Level,
			change.After.CurrentXp,
			change.After.TotalXp,
			change.After.NextLevelXp,
			i.GuildID,
		)
		if err != nil {
			sentry.CaptureException(err)
			editXpAdjustmentResponse(s, i, "I couldn't update leaderboard data right now. Try again later.")
			return
		}
	}

	leveling.SyncRoleRewards(s, i.GuildID, targetUser.ID, controller, change.After.Level)

	content := formatXpAdjustmentResponse(targetUser.ID, change, !userExists, direction < 0)
	editXpAdjustmentResponse(s, i, content)
}

func xpAdjustmentOptions(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.User, int, error) {
	var targetUser *discordgo.User
	var amount int

	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "user":
			targetUser = option.UserValue(s)
		case "amount":
			amount = int(option.IntValue())
		}
	}

	if targetUser == nil {
		return nil, 0, fmt.Errorf("Unable to resolve the requested member")
	}

	if amount <= 0 {
		return nil, 0, fmt.Errorf("XP amount must be greater than zero")
	}

	return targetUser, amount, nil
}

func editXpAdjustmentResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	if err != nil {
		sentry.CaptureException(err)
	}
}

func formatXpAdjustmentResponse(userID string, change leveling.UserProgressChange, createdRecord bool, removing bool) string {
	action := "Granted"
	if removing {
		action = "Removed"
	}

	lines := []string{
		fmt.Sprintf("%s %d XP %s <@%s>.", action, abs(change.AppliedDelta), xpDirectionPreposition(removing), userID),
		fmt.Sprintf("Level: %d -> %d", change.Before.Level, change.After.Level),
		fmt.Sprintf("Progress: %d / %d XP -> %d / %d XP", change.Before.CurrentXp, change.Before.NextLevelXp, change.After.CurrentXp, change.After.NextLevelXp),
		fmt.Sprintf("Total XP: %d -> %d", change.Before.TotalXp, change.After.TotalXp),
	}

	if createdRecord {
		lines = append(lines, "A new leaderboard record was created for this member.")
	}

	if removing && change.AppliedDelta != change.RequestedDelta {
		lines = append(lines, fmt.Sprintf("Requested removal was capped because the member only had %d total XP available.", change.Before.TotalXp))
	}

	levelDelta := change.After.Level - change.Before.Level
	switch {
	case levelDelta > 0:
		lines = append(lines, fmt.Sprintf("Level increase: +%d", levelDelta))
	case levelDelta < 0:
		lines = append(lines, fmt.Sprintf("Level decrease: %d", abs(levelDelta)))
	}

	return strings.Join(lines, "\n")
}

func xpDirectionPreposition(removing bool) string {
	if removing {
		return "from"
	}

	return "to"
}

func leaderboardUserNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "404")
}

func abs(value int) int {
	if value < 0 {
		return -value
	}

	return value
}
