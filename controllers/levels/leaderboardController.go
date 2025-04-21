package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"io"
	"math/rand"
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

func (c *LeaderboardController) FindUser(id string, guildId string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/discord/leaderboard/users/find/%s", config.GetApiBaseUrl(guildId), id)
	resp, err := c.sendRequest("GET", url, nil, guildId)
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

func (c *LeaderboardController) AddUser(userId, guildId string) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/create", config.GetApiBaseUrl(guildId))
	xpPerMessage := rand.Intn(15) + 10
	data := UserCreate{
		GuildID:         guildId,
		UserID:          userId,
		MessageCount:    1,
		Xp:              xpPerMessage,
		TotalXp:         xpPerMessage,
		LevelXp:         0,
		Level:           0,
		Rank:            0,
		NoXp:            false,
		MessageLastSent: time.Now().UnixMilli(),
	}
	_, err := c.sendRequest("POST", url, data, guildId)
	return err
}

func (c *LeaderboardController) UpdateUserRole(id, roleId string, guildId string) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/%s", config.GetApiBaseUrl(guildId), id)
	data := map[string]string{"roleId": roleId}
	_, err := c.sendRequest("PATCH", url, data, guildId)
	return err
}

func (c *LeaderboardController) UpdateUserLevel(id string, level, messageCount, xp, levelXp int, guildId string) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/%s", config.GetApiBaseUrl(guildId), id)
	data := map[string]interface{}{
		"level":        level,
		"userId":       id,
		"messageCount": messageCount,
		"xp":           xp,
		"levelXp":      levelXp,
	}
	_, err := c.sendRequest("PATCH", url, data, guildId)
	return err
}

func (c *LeaderboardController) UpdateUserPoints(id string, level, messageCount, xp, totalXp, levelXp int, messageLastSent int64, guildId string) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/%s", config.GetApiBaseUrl(guildId), id)
	data := map[string]interface{}{
		"level":           level,
		"messageCount":    messageCount,
		"xp":              xp,
		"totalXp":         totalXp,
		"levelXp":         levelXp,
		"messageLastSent": messageLastSent,
	}
	_, err := c.sendRequest("PATCH", url, data, guildId)
	return err
}

func (c *LeaderboardController) AddXp(id string, xp, totalXp, level, levelXp int, guildId string) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/%s", config.GetApiBaseUrl(guildId), id)
	data := map[string]interface{}{
		"level":   level,
		"xp":      xp,
		"totalXp": totalXp,
		"levelXp": levelXp,
	}
	_, err := c.sendRequest("PATCH", url, data, guildId)
	return err
}

func (c *LeaderboardController) NoUserXp(id string, xp bool, guildId string) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/%s", config.GetApiBaseUrl(guildId), id)
	data := map[string]bool{"noXp": xp}
	_, err := c.sendRequest("PATCH", url, data, guildId)
	return err
}

// Helper function to handle HTTP requests
func (c *LeaderboardController) sendRequest(method, url string, body interface{}, guildId string) (*http.Response, error) {
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
	req.Header.Set("X-API-Key", config.GetInternalApiKey(guildId))

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
