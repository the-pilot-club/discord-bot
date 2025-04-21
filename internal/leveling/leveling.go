package leveling

import (
	"math"

	controllers "tpc-discord-bot/controllers/levels"
	"tpc-discord-bot/internal/config"

	"github.com/bwmarrin/discordgo"
)

// CalculateUserLevel calculates the user's level and required XP for next level
func CalculateUserLevel(currentLevel, currentXp int) (level, requiredXp int) {
	xpForNextLevel := int(math.Pow(float64(currentLevel+1), 2) * 100)

	if currentXp >= xpForNextLevel {
		return currentLevel + 1, xpForNextLevel
	}

	return currentLevel, xpForNextLevel
}

// CheckRoleRewards checks and assigns role rewards based on level
func CheckRoleRewards(s *discordgo.Session, guildID string, userID string,
	controller *controllers.LeaderboardController, newLevel int) {

	roleRewards := config.GetRoleRewards(guildID)

	for _, reward := range roleRewards {
		if reward.Level == newLevel {
			err := s.GuildMemberRoleAdd(guildID, userID, reward.RoleID)
			if err != nil {
				continue
			}

			controller.UpdateUserRole(userID, reward.RoleID, guildID)
		}
	}
}
