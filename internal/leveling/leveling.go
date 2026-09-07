package leveling

import (
	"math"

	controllers "tpc-discord-bot/controllers/levels"
	"tpc-discord-bot/internal/config"

	"github.com/bwmarrin/discordgo"
)

func XpForNextLevel(currentLevel int) int {
	return int(math.Pow(float64(currentLevel+1), 2) * 100)
}

type UserProgress struct {
	Level       int
	CurrentXp   int
	TotalXp     int
	NextLevelXp int
}

type UserProgressChange struct {
	Before         UserProgress
	After          UserProgress
	RequestedDelta int
	AppliedDelta   int
}

type leaderboardXpStore interface {
	FindUserStats(id, guildID string) (*controllers.UserStats, bool, error)
	CreateUserRecord(data controllers.UserCreate, guildID string) error
	UpdateUserXpState(id string, level, xp, totalXp, levelXp int, guildID string) error
}

func AwardXp(s *discordgo.Session, controller *controllers.LeaderboardController, guildID, userID string, amount int) (UserProgressChange, error) {
	return awardXp(controller, guildID, userID, amount, func(newLevel int) {
		SyncRoleRewards(s, guildID, userID, controller, newLevel)
	})
}

func awardXp(store leaderboardXpStore, guildID, userID string, amount int, syncRoleRewards func(newLevel int)) (UserProgressChange, error) {
	stats, found, err := store.FindUserStats(userID, guildID)
	if err != nil {
		return UserProgressChange{}, err
	}

	if stats == nil {
		stats = &controllers.UserStats{}
	}

	change := ApplyXpDelta(stats.Level, stats.CurrentXp, stats.TotalXp, amount)
	if !found {
		err = store.CreateUserRecord(controllers.UserCreate{
			GuildID:         guildID,
			UserID:          userID,
			MessageCount:    0,
			Xp:              change.After.CurrentXp,
			TotalXp:         change.After.TotalXp,
			LevelXp:         change.After.NextLevelXp,
			Level:           change.After.Level,
			Rank:            0,
			NoXp:            false,
			MessageLastSent: 0,
		}, guildID)
	} else {
		err = store.UpdateUserXpState(
			userID,
			change.After.Level,
			change.After.CurrentXp,
			change.After.TotalXp,
			change.After.NextLevelXp,
			guildID,
		)
	}
	if err != nil {
		return UserProgressChange{}, err
	}

	if syncRoleRewards != nil {
		syncRoleRewards(change.After.Level)
	}
	return change, nil
}

func TotalXpForLevel(level int) int {
	if level <= 0 {
		return 0
	}

	return 100 * level * (level + 1) * (2*level + 1) / 6
}

func ProgressFromTotalXp(totalXp int) UserProgress {
	if totalXp < 0 {
		totalXp = 0
	}

	level := 0
	currentXp := totalXp

	for {
		nextLevelXp := XpForNextLevel(level)
		if currentXp < nextLevelXp {
			return UserProgress{
				Level:       level,
				CurrentXp:   currentXp,
				TotalXp:     totalXp,
				NextLevelXp: nextLevelXp,
			}
		}

		currentXp -= nextLevelXp
		level++
	}
}

func NormalizeProgress(level, currentXp, totalXp int) UserProgress {
	if level < 0 {
		level = 0
	}

	if currentXp < 0 {
		currentXp = 0
	}

	if totalXp < 0 {
		totalXp = 0
	}

	minimumTotalXp := TotalXpForLevel(level) + currentXp
	if minimumTotalXp > totalXp {
		totalXp = minimumTotalXp
	}

	return ProgressFromTotalXp(totalXp)
}

func ApplyXpDelta(level, currentXp, totalXp, delta int) UserProgressChange {
	before := NormalizeProgress(level, currentXp, totalXp)

	afterTotalXp := before.TotalXp + delta
	if afterTotalXp < 0 {
		afterTotalXp = 0
	}

	after := ProgressFromTotalXp(afterTotalXp)

	return UserProgressChange{
		Before:         before,
		After:          after,
		RequestedDelta: delta,
		AppliedDelta:   after.TotalXp - before.TotalXp,
	}
}

// CalculateUserLevel calculates the user's level and required XP for next level
func CalculateUserLevel(currentLevel, currentXp int) (level, requiredXp int) {
	xpForNextLevel := XpForNextLevel(currentLevel)

	if currentXp >= xpForNextLevel {
		return currentLevel + 1, xpForNextLevel
	}

	return currentLevel, xpForNextLevel
}

func SyncRoleRewards(s *discordgo.Session, guildID string, userID string,
	controller *controllers.LeaderboardController, newLevel int) {

	roleRewards := config.GetRoleRewards(guildID)
	if len(roleRewards) == 0 {
		return
	}

	member, err := s.State.Member(guildID, userID)
	if err != nil || member == nil {
		member, _ = s.GuildMember(guildID, userID)
	}

	currentRoles := make(map[string]struct{}, len(roleRewards))
	knownRoles := member != nil
	if member != nil {
		for _, roleID := range member.Roles {
			currentRoles[roleID] = struct{}{}
		}
	}

	highestRewardLevel := -1
	highestRewardRoleID := ""

	for _, reward := range roleRewards {
		if reward.Level <= newLevel {
			if _, hasRole := currentRoles[reward.RoleID]; !hasRole || !knownRoles {
				if err := s.GuildMemberRoleAdd(guildID, userID, reward.RoleID); err == nil {
					currentRoles[reward.RoleID] = struct{}{}
				}
			}

			if reward.Level >= highestRewardLevel {
				highestRewardLevel = reward.Level
				highestRewardRoleID = reward.RoleID
			}

			continue
		}

		if _, hasRole := currentRoles[reward.RoleID]; hasRole || !knownRoles {
			if err := s.GuildMemberRoleRemove(guildID, userID, reward.RoleID); err == nil {
				delete(currentRoles, reward.RoleID)
			}
		}
	}

	if highestRewardRoleID != "" {
		_ = controller.UpdateUserRole(userID, highestRewardRoleID, guildID)
	}
}

// CheckRoleRewards checks and assigns role rewards based on level.
func CheckRoleRewards(s *discordgo.Session, guildID string, userID string,
	controller *controllers.LeaderboardController, newLevel int) {

	SyncRoleRewards(s, guildID, userID, controller, newLevel)
}
