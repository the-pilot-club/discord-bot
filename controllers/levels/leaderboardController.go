package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"

	"tpc-discord-bot/internal/config"
)

type LeaderboardController struct{}

type LeaderboardPage struct {
	TotalCount int                      `json:"totalCount"`
	PageCount  int                      `json:"pageCount"`
	Items      []map[string]interface{} `json:"items"`
}

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
	url := fmt.Sprintf("%s/discord/leaderboard/users/find/%s", config.GetBaseUrl(guildId, "Internal API"), id)
	resp, err := c.sendRequest("GET", url, nil, guildId)
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}
	defer closeResponseBody(resp)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

func (c *LeaderboardController) FindLeaderboardUsers(guildId string, offset, limit int) (*LeaderboardPage, error) {
	url := fmt.Sprintf("%s/discord/leaderboard/users?offset=%d&limit=%d", config.GetBaseUrl(guildId, "Internal API"), offset, limit)
	resp, err := c.sendRequest("GET", url, nil, guildId)
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}
	defer closeResponseBody(resp)

	var result LeaderboardPage
	err = json.NewDecoder(resp.Body).Decode(&result)
	return &result, err
}

func (c *LeaderboardController) AddUser(userId, guildId string) error {
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
	return c.CreateUserRecord(data, guildId)
}

func (c *LeaderboardController) CreateUserRecord(data UserCreate, guildId string) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/create", config.GetBaseUrl(guildId, "Internal API"))
	resp, err := c.sendRequest("POST", url, data, guildId)
	closeResponseBody(resp)
	return err
}

func (c *LeaderboardController) UpdateUserRole(id, roleId string, guildId string) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/%s", config.GetBaseUrl(guildId, "Internal API"), id)
	data := map[string]string{"roleId": roleId}
	resp, err := c.sendRequest("PATCH", url, data, guildId)
	closeResponseBody(resp)
	return err
}

func (c *LeaderboardController) UpdateUserLevel(id string, level, messageCount, xp, levelXp int, guildId string) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/%s", config.GetBaseUrl(guildId, "Internal API"), id)
	data := map[string]interface{}{
		"level":        level,
		"userId":       id,
		"messageCount": messageCount,
		"xp":           xp,
		"levelXp":      levelXp,
	}
	resp, err := c.sendRequest("PATCH", url, data, guildId)
	closeResponseBody(resp)
	return err
}

func (c *LeaderboardController) UpdateUserPoints(id string, level, messageCount, xp, totalXp, levelXp int, messageLastSent int64, guildId string) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/%s", config.GetBaseUrl(guildId, "Internal API"), id)
	data := map[string]interface{}{
		"level":           level,
		"messageCount":    messageCount,
		"xp":              xp,
		"totalXp":         totalXp,
		"levelXp":         levelXp,
		"messageLastSent": messageLastSent,
	}
	resp, err := c.sendRequest("PATCH", url, data, guildId)
	closeResponseBody(resp)
	return err
}

func (c *LeaderboardController) AddXp(id string, xp, totalXp, level, levelXp int, guildId string) error {
	return c.UpdateUserXpState(id, level, xp, totalXp, levelXp, guildId)
}

func (c *LeaderboardController) UpdateUserXpState(id string, level, xp, totalXp, levelXp int, guildId string) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/%s", config.GetBaseUrl(guildId, "Internal API"), id)
	data := map[string]interface{}{
		"level":   level,
		"xp":      xp,
		"totalXp": totalXp,
		"levelXp": levelXp,
	}
	resp, err := c.sendRequest("PATCH", url, data, guildId)
	closeResponseBody(resp)
	return err
}

func (c *LeaderboardController) NoUserXp(id string, xp bool, guildId string) error {
	url := fmt.Sprintf("%s/discord/leaderboard/users/%s", config.GetBaseUrl(guildId, "Internal API"), id)
	data := map[string]bool{"noXp": xp}
	resp, err := c.sendRequest("PATCH", url, data, guildId)
	closeResponseBody(resp)
	return err
}

func closeResponseBody(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
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
		closeResponseBody(resp)
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	return resp, nil
}
