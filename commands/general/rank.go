package general

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"

	controllers "tpc-discord-bot/controllers/levels"
	"tpc-discord-bot/internal/leveling"
)

const tpcEmbedColor = 3651327

type leaderboardStats struct {
	Rank         int
	Level        int
	CurrentXp    int
	TotalXp      int
	MessageCount int
	NoXp         bool
}

func HandleRankCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	var targetUser *discordgo.User
	if i.Member != nil {
		targetUser = i.Member.User
	}

	options := i.ApplicationCommandData().Options
	if len(options) > 0 {
		targetUser = options[0].UserValue(s)
	}

	if targetUser == nil {
		editRankError(s, i, "Unable to resolve the requested member.")
		return
	}

	if targetUser.Bot {
		editRankError(s, i, "Bots do not have leaderboard ranks.")
		return
	}

	controller := &controllers.LeaderboardController{}
	userData, err := controller.FindUser(targetUser.ID, i.GuildID)
	if err != nil {
		if !strings.Contains(err.Error(), "404") {
			sentry.CaptureException(err)
			editRankError(s, i, "I couldn't fetch leaderboard data right now. Try again later.")
			return
		}

		editRankError(s, i, fmt.Sprintf("<@%s> does not have leaderboard data yet. Send a message in an XP-enabled channel first.", targetUser.ID))
		return
	}

	stats, err := parseLeaderboardStats(userData)
	if err != nil {
		sentry.CaptureException(err)
		editRankError(s, i, "Leaderboard data for that member could not be read.")
		return
	}

	embed := buildRankEmbed(i, targetUser, stats)
	components := leaderboardLinkComponents()
	embeds := []*discordgo.MessageEmbed{embed}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds:     &embeds,
		Components: &components,
	})
	if err != nil {
		sentry.CaptureException(err)
	}
}

func editRankError(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	if err != nil {
		sentry.CaptureException(err)
	}
}

func parseLeaderboardStats(user map[string]interface{}) (*leaderboardStats, error) {
	level, err := leaderboardInt(user, "level")
	if err != nil {
		return nil, err
	}

	currentXp, err := leaderboardInt(user, "xp")
	if err != nil {
		return nil, err
	}

	totalXp, err := leaderboardInt(user, "totalXp")
	if err != nil {
		return nil, err
	}

	messageCount, err := leaderboardInt(user, "messageCount")
	if err != nil {
		return nil, err
	}

	stats := &leaderboardStats{
		Rank:         leaderboardOptionalInt(user, "rank"),
		Level:        level,
		CurrentXp:    currentXp,
		TotalXp:      totalXp,
		MessageCount: messageCount,
		NoXp:         leaderboardOptionalBool(user, "noXp"),
	}

	return stats, nil
}

func buildRankEmbed(i *discordgo.InteractionCreate, user *discordgo.User, stats *leaderboardStats) *discordgo.MessageEmbed {
	displayName := rankDisplayName(i, user)
	nextLevelXp := leveling.XpForNextLevel(stats.Level)
	progressPct := 0
	if nextLevelXp > 0 {
		progressPct = (stats.CurrentXp * 100) / nextLevelXp
		if progressPct > 100 {
			progressPct = 100
		}
	}

	rankValue := "Unranked"
	if stats.Rank > 0 {
		rankValue = fmt.Sprintf("#%d", stats.Rank)
	}

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Rank",
			Value:  rankValue,
			Inline: true,
		},
		{
			Name:   "Level",
			Value:  strconv.Itoa(stats.Level),
			Inline: true,
		},
		{
			Name:   "Messages",
			Value:  strconv.Itoa(stats.MessageCount),
			Inline: true,
		},
		{
			Name:   "Progress",
			Value:  fmt.Sprintf("%d / %d XP (%d%%)", stats.CurrentXp, nextLevelXp, progressPct),
			Inline: true,
		},
		{
			Name:   "Total XP",
			Value:  strconv.Itoa(stats.TotalXp),
			Inline: true,
		},
	}

	if stats.NoXp {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "XP Status",
			Value:  "Disabled",
			Inline: true,
		})
	}

	return &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    displayName,
			IconURL: user.AvatarURL(""),
		},
		Title:       fmt.Sprintf("%s's Leaderboard Rank", displayName),
		Description: "Use the button below to open the full leaderboard.",
		Fields:      fields,
		Color:       tpcEmbedColor,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Made by TPC Tech Team",
			IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
		},
	}
}

func rankDisplayName(i *discordgo.InteractionCreate, user *discordgo.User) string {
	if i.Member != nil && i.Member.User != nil && i.Member.User.ID == user.ID && i.Member.Nick != "" {
		return i.Member.Nick
	}

	if resolved := i.ApplicationCommandData().Resolved; resolved != nil {
		if member, ok := resolved.Members[user.ID]; ok && member.Nick != "" {
			return member.Nick
		}
	}

	return user.Username
}

func leaderboardInt(user map[string]interface{}, key string) (int, error) {
	value, ok := user[key]
	if !ok || value == nil {
		return 0, fmt.Errorf("leaderboard field %q missing", key)
	}

	switch v := value.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("leaderboard field %q is not numeric: %w", key, err)
		}
		return i, nil
	default:
		return 0, fmt.Errorf("leaderboard field %q has unsupported type %T", key, value)
	}
}

func leaderboardOptionalInt(user map[string]interface{}, key string) int {
	value, err := leaderboardInt(user, key)
	if err != nil {
		return 0
	}

	return value
}

func leaderboardOptionalBool(user map[string]interface{}, key string) bool {
	value, ok := user[key]
	if !ok || value == nil {
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case string:
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return false
		}
		return parsed
	default:
		return false
	}
}
