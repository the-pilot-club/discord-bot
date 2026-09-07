package handlers

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"

	controllers "tpc-discord-bot/controllers/levels"
	"tpc-discord-bot/internal/config"
	"tpc-discord-bot/internal/leveling"

	"github.com/bwmarrin/discordgo"
)

type LevelingConfig struct {
	XpRate float64
}

func HandleXpGive(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Skip if message is from bot or not in guild
	if m.Author.Bot || m.GuildID == "" || !config.GetXpGiveEnabled(m.GuildID) {
		return
	}

	// Skip if message is from Serge
	if m.Author.ID == "524567291128709140" {
		return
	}

	controller := &controllers.LeaderboardController{}
	levelingConfig := LevelingConfig{XpRate: 1.0}

	// Get channel
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		return
	}

	// Skip if the channel (or a category/channel it lives under) has XP disabled
	if !config.ValidXpChannel(m.GuildID, channel, xpChannelAncestors(s, channel)...) {
		return
	}

	// Find or create user
	user, err := controller.FindUser(m.Author.ID, m.GuildID)
	if err != nil {
		if !leaderboardUserNotFound(err) {
			return
		}

		xpPerMessage := rand.Intn(15) + 10
		initialProgress := leveling.ProgressFromTotalXp(int(float64(xpPerMessage) * levelingConfig.XpRate))

		err = controller.CreateUserRecord(controllers.UserCreate{
			GuildID:         m.GuildID,
			UserID:          m.Author.ID,
			MessageCount:    1,
			Xp:              initialProgress.CurrentXp,
			TotalXp:         initialProgress.TotalXp,
			LevelXp:         initialProgress.NextLevelXp,
			Level:           initialProgress.Level,
			Rank:            0,
			NoXp:            false,
			MessageLastSent: time.Now().Add(time.Minute).UnixMilli(),
		}, m.GuildID)
		if err != nil {
			sentry.CaptureException(err)
		}
		return
	}

	// Check if user has noXp flag
	if leaderboardOptionalBool(user["noXp"]) {
		return
	}

	//Add a check to see if member has sent a message in the last minute
	if processLastMessageSent(fmt.Sprint(user["messageLastSent"])) {
		return
	}

	// Calculate XP
	xpPerMessage := rand.Intn(15) + 10 // Random between 7-12
	xpToAssign := float64(xpPerMessage) * levelingConfig.XpRate

	// Get current user stats
	currentLevel, err := leaderboardInt(user, "level")
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	currentXp, err := leaderboardInt(user, "xp")
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	totalXp, err := leaderboardInt(user, "totalXp")
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	messageCount, err := leaderboardInt(user, "messageCount")
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	change := leveling.ApplyXpDelta(currentLevel, currentXp, totalXp, int(xpToAssign))

	err = controller.UpdateUserPoints(
		m.Author.ID,
		change.After.Level,
		messageCount+1,
		change.After.CurrentXp,
		change.After.TotalXp,
		change.After.NextLevelXp,
		time.Now().Add(time.Minute).UnixMilli(),
		m.GuildID,
	)
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	if change.After.Level > change.Before.Level {
		content := fmt.Sprintf("Congrats <@%v>, you just advanced to TPC **level %v**!", m.Author.ID, change.After.Level)
		_, err = s.ChannelMessageSend(m.ChannelID, content)
		if err != nil {
			sentry.CaptureException(err)
		}

		leveling.SyncRoleRewards(s, m.GuildID, m.Author.ID, controller, change.After.Level)
	}
}

// helper func to walk up from a channel to collect its parent channel and
// category, so ValidXpChannel can honour XP being disabled on a thread's parent
// channel or its category.
func xpChannelAncestors(s *discordgo.Session, channel *discordgo.Channel) []*discordgo.Channel {
	var ancestors []*discordgo.Channel
	current := channel
	for i := 0; i < 2 && current != nil && current.ParentID != ""; i++ {
		parent, err := s.Channel(current.ParentID)
		if err != nil || parent == nil {
			break
		}
		ancestors = append(ancestors, parent)
		current = parent
	}
	return ancestors
}

func processLastMessageSent(m string) bool {
	i, err := strconv.ParseInt(m, 10, 64)

	if err != nil {
		sentry.CaptureException(err)
	}

	if i > time.Now().UnixMilli() {
		return true
	}
	return false
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

func leaderboardOptionalBool(value interface{}) bool {
	if value == nil {
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

func leaderboardUserNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "404")
}
