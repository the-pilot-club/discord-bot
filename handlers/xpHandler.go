package handlers

import (
	"math"
	"math/rand"
	"strings"
	"time"

	controllers "tpc-discord-bot/controllers/levels"
	"tpc-discord-bot/internal/config"

	"github.com/bwmarrin/discordgo"
)

type LevelingConfig struct {
	XpRate float64
}

func HandleXpGive(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Skip if message is from bot or not in guild
	if m.Author.Bot || m.GuildID == "" {
		return
	}

	controller := &controllers.LeaderboardController{}
	levelingConfig := LevelingConfig{XpRate: 1.0}

	// List of XP-enabled channels
	xpChannels := []string{
		"crew-chat", "fly-with-me", "Event Comms", "Flight Voice 1",
		"Flight Voice 2", "Shared Cockpit", "Magneto Lounge", "1 on 1",
		"general-aviation-ops", "streams-and-videos", "screenshots",
		"Charters Comms", "community-help-forum", "pilot-training-chat",
		"charters-chat", "explorer-missions", "irl-pilots", "atc-lounge",
		"company-perks", "early-adopters", "airliner-and-bizjet-ops",
		"rotary-ops", "your-setups",
	}

	// Get channel
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		return
	}

	// Check if channel is in XP channels list
	channelAllowed := false
	for _, xpChannel := range xpChannels {
		if strings.EqualFold(channel.Name, xpChannel) {
			channelAllowed = true
			break
		}
	}
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

	// Calculate XP
	xpPerMessage := rand.Intn(5) + 7 // Random between 7-12
	xpToAssign := float64(xpPerMessage) * levelingConfig.XpRate

	// Get current user stats
	currentLevel := int(user["level"].(float64))
	currentXp := int(user["xp"].(float64))
	totalXp := int(user["totalXp"].(float64))
	messageCount := int(user["messageCount"].(float64))

	// Calculate new level
	newLevel, requiredXp := calculateUserLevel(currentLevel, currentXp+int(xpToAssign))

	if newLevel > currentLevel {
		// Level up
		_, err = s.ChannelMessageSend(m.ChannelID,
			"Congrats <@"+m.Author.ID+">, you just advanced to TPC **level "+
				string(rune(newLevel))+"**!")
		if err != nil {
			return
		}

		err = controller.UpdateUserLevel(m.Author.ID, newLevel, messageCount+1, 1, requiredXp, m.GuildID)
		if err != nil {
			return
		}

		// Check and assign role rewards
		checkRoleRewards(s, m, controller, newLevel)
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

// calculate user level and required xp
func calculateUserLevel(currentLevel, currentXp int) (level, requiredXp int) {
	xpForNextLevel := int(math.Pow(float64(currentLevel+1), 2) * 100)

	if currentXp >= xpForNextLevel {
		return currentLevel + 1, xpForNextLevel
	}

	return currentLevel, xpForNextLevel
}

// see if user has reached a new level and assign role rewards
func checkRoleRewards(s *discordgo.Session, m *discordgo.MessageCreate,
	controller *controllers.LeaderboardController, newLevel int) {

	roleRewards := config.GetRoleRewards(m.GuildID)

	for _, reward := range roleRewards {
		if reward.Level == newLevel {
			err := s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, reward.RoleID)
			if err != nil {
				continue
			}

			controller.UpdateUserRole(m.Author.ID, reward.RoleID, m.GuildID)
		}
	}
}
