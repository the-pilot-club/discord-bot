package handlers

import (
	"fmt"
	"math/rand"
	"strconv"
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

	// Check if channel is in XP channels list
	channelAllowed := config.ValidXpChannel(m.GuildID, channel)
	if !channelAllowed {
		return
	}

	// Find or create user
	user, err := controller.FindUser(m.Author.ID, m.GuildID)
	if err != nil {
		// Add user if not found
		err = controller.AddUser(m.Author.ID, m.GuildID)
		if err != nil {
			return
		}
		return
	}

	// Check if user has noXp flag
	if noXp, ok := user["noXp"].(bool); ok && noXp {
		return
	}

	//Add a check to see if member has sent a message in the last minute
	if processLastMessageSent(user["messageLastSent"].(string)) {
		return
	}

	// Calculate XP
	xpPerMessage := rand.Intn(15) + 10 // Random between 7-12
	xpToAssign := float64(xpPerMessage) * levelingConfig.XpRate

	// Get current user stats
	currentLevel := int(user["level"].(float64))
	currentXp := int(user["xp"].(float64))
	totalXp := int(user["totalXp"].(float64))
	messageCount := int(user["messageCount"].(float64))

	// Calculate new level
	newLevel, requiredXp := leveling.CalculateUserLevel(currentLevel, currentXp+int(xpToAssign))

	if newLevel > currentLevel {
		// Level up
		content := fmt.Sprintf("Congrats <@%v>, you just advanced to TPC **level %v **!", m.Author.ID, newLevel)
		_, err = s.ChannelMessageSend(m.ChannelID, content)
		if err != nil {
			return
		}

		err = controller.UpdateUserLevel(m.Author.ID, newLevel, messageCount+1, 1, requiredXp, m.GuildID)
		if err != nil {
			return
		}

		// Check and assign role rewards
		leveling.CheckRoleRewards(s, m.GuildID, m.Author.ID, controller, newLevel)
		return
	}

	// Update points
	err = controller.UpdateUserPoints(
		m.Author.ID,
		newLevel,
		messageCount+1,
		currentXp+int(xpToAssign),
		totalXp+int(xpToAssign),
		requiredXp,
		time.Now().Add(time.Minute).UnixMilli(),
		m.GuildID,
	)
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
