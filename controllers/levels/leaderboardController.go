package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"tpc-discord-bot/internal/config"
)

type LeaderboardController struct{}

type UserCreate struct {
	GuildID         string `json:"guildId"`
	UserID          string `json:"userId"`
	MessageCount    int    `json:"messageCount"`
	Xp              int    `json:"xp"`
	TotalXp         int    `json:"totalXp"`
	LevelXp         int    `json:"levelXp"`
	Level           int    `json:"level"`
	Rank            int    `json:"rank"`
	NoXp            bool   `json:"noXp"`
	MessageLastSent int64  `json:"messageLastSent"`
}

type RoleReward struct {
	GuildID string `json:"guildId"`
	Level   int    `json:"level"`
	RoleID  string `json:"roleId"`
}

func (c *LeaderboardController) FindUser(id string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/discord/leaderboard/users/find/%s", config.QuizBaseUrl, id)
	resp, err := c.sendRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

func (c *LeaderboardController) AddUser(userId, guildId string) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/create", config.QuizBaseUrl)
	data := UserCreate{
		GuildID:         guildId,
		UserID:          userId,
		MessageCount:    1,
		Xp:              0,
		TotalXp:         0,
		LevelXp:         0,
		Level:           0,
		Rank:            0,
		NoXp:            false,
		MessageLastSent: time.Now().UnixMilli(),
	}
	_, err := c.sendRequest("POST", url, data)
	return err
}

func (c *LeaderboardController) UpdateUserRole(id, roleId string) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/%s", config.QuizBaseUrl, id)
	data := map[string]string{"roleId": roleId}
	_, err := c.sendRequest("PATCH", url, data)
	return err
}

func (c *LeaderboardController) UpdateUserLevel(id string, level, messageCount, xp, levelXp int) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/%s", config.QuizBaseUrl, id)
	data := map[string]interface{}{
		"level":        level,
		"userId":       id,
		"messageCount": messageCount,
		"xp":           xp,
		"levelXp":      levelXp,
	}
	_, err := c.sendRequest("PATCH", url, data)
	return err
}

func (c *LeaderboardController) UpdateUserPoints(id string, level, messageCount, xp, totalXp, levelXp int, messageLastSent int64) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/%s", config.QuizBaseUrl, id)
	data := map[string]interface{}{
		"level":           level,
		"messageCount":    messageCount,
		"xp":              xp,
		"totalXp":         totalXp,
		"levelXp":         levelXp,
		"messageLastSent": messageLastSent,
	}
	_, err := c.sendRequest("PATCH", url, data)
	return err
}

func (c *LeaderboardController) AddXp(id string, xp, totalXp, level, levelXp int) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/%s", config.QuizBaseUrl, id)
	data := map[string]interface{}{
		"level":   level,
		"xp":      xp,
		"totalXp": totalXp,
		"levelXp": levelXp,
	}
	_, err := c.sendRequest("PATCH", url, data)
	return err
}

func (c *LeaderboardController) NoUserXp(id string, xp bool) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/%s", config.QuizBaseUrl, id)
	data := map[string]bool{"noXp": xp}
	_, err := c.sendRequest("PATCH", url, data)
	return err
}

func (c *LeaderboardController) GetRoleRewards() []RoleReward {
	return []RoleReward{
		{
			GuildID: config.GuildId,
			RoleID:  config.CommuterRoleId,
			Level:   15,
		},
		{
			GuildID: config.GuildId,
			RoleID:  config.FfRoleId,
			Level:   4,
		},
		{
			GuildID: config.GuildId,
			RoleID:  config.VipRoleId,
			Level:   3,
		},
	}
}

// Helper function to handle HTTP requests
func (c *LeaderboardController) sendRequest(method, url string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "TPCDiscordBot")
	req.Header.Set("X-API-Key", config.QuizToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	return resp, nil
}
