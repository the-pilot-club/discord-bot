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
	handleXpAdjustmentCommand(s, i, false)
}

func HandleRemoveXpCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	handleXpAdjustmentCommand(s, i, true)
}

func handleXpAdjustmentCommand(s *discordgo.Session, i *discordgo.InteractionCreate, remove bool) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	targetUser, xpAmount, err := xpAdjustmentOptions(s, i)
	if err != nil {
		editXpAdjustmentError(s, i, "XP Update Failed", err.Error())
		return
	}

	if targetUser.Bot {
		editXpAdjustmentWarning(s, i, "XP Not Updated", "Bots do not have leaderboard XP.")
		return
	}

	delta := xpAmount
	if remove {
		delta = -xpAmount
	}
	controller := &controllers.LeaderboardController{}
	if remove {
		_, found, err := controller.FindUserStats(targetUser.ID, i.GuildID)
		if err != nil {
			sentry.CaptureException(err)
			editXpAdjustmentError(s, i, "XP Update Failed", "I couldn't fetch leaderboard data right now. Try again later.")
			return
		}

		if !found {
			editXpAdjustmentWarning(
				s,
				i,
				"No XP Removed",
				fmt.Sprintf("<@%s> does not have leaderboard data yet, so there is no XP to remove.", targetUser.ID),
			)
			return
		}
	}

	change, err := leveling.AwardXp(s, controller, i.GuildID, targetUser.ID, delta)
	if err != nil {
		sentry.CaptureException(err)
		editXpAdjustmentError(s, i, "XP Update Failed", "I couldn't update leaderboard data right now. Try again later.")
		return
	}

	editXpAdjustmentEmbed(s, i, buildXpAdjustmentEmbed(i, targetUser, change, false, remove))
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

func editXpAdjustmentEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	embeds := []*discordgo.MessageEmbed{embed}
	emptyContent := ""

	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &emptyContent,
		Embeds:  &embeds,
	})
	if err != nil {
		sentry.CaptureException(err)
	}
}

func editXpAdjustmentWarning(s *discordgo.Session, i *discordgo.InteractionCreate, title string, description string) {
	editXpAdjustmentEmbed(s, i, buildXpAdjustmentStatusEmbed(title, description, colorWarning))
}

func editXpAdjustmentError(s *discordgo.Session, i *discordgo.InteractionCreate, title string, description string) {
	editXpAdjustmentEmbed(s, i, buildXpAdjustmentStatusEmbed(title, description, colorError))
}

func buildXpAdjustmentStatusEmbed(title string, description string, color int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
	}
}

func buildXpAdjustmentEmbed(i *discordgo.InteractionCreate, targetUser *discordgo.User, change leveling.UserProgressChange, createdRecord bool, removing bool) *discordgo.MessageEmbed {
	description := fmt.Sprintf("✅ %s XP has been given to <@%s>", formatXpAmount(change.AppliedDelta), targetUser.ID)
	if removing {
		description = fmt.Sprintf("✅ %s XP has been removed from <@%s>", formatXpAmount(abs(change.AppliedDelta)), targetUser.ID)
	}

	if removing && change.AppliedDelta != change.RequestedDelta {
		description = fmt.Sprintf("⚠️ %s XP was removed from <@%s> because that was their full available total.", formatXpAmount(abs(change.AppliedDelta)), targetUser.ID)
	}

	return &discordgo.MessageEmbed{
		Description: description,
		Color:       tpcEmbedColor,
	}
}

func leaderboardUserNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "404")
}

func formatXpAmount(value int) string {
	if value < 0 {
		value = -value
	}

	text := fmt.Sprintf("%d", value)
	if len(text) <= 3 {
		return text
	}

	var formatted []byte
	prefixLen := len(text) % 3
	if prefixLen == 0 {
		prefixLen = 3
	}

	formatted = append(formatted, text[:prefixLen]...)
	for idx := prefixLen; idx < len(text); idx += 3 {
		formatted = append(formatted, ',')
		formatted = append(formatted, text[idx:idx+3]...)
	}

	return string(formatted)
}

func abs(value int) int {
	if value < 0 {
		return -value
	}

	return value
}
