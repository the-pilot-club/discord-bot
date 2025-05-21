package fcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"net/http"
	"tpc-discord-bot/internal/config"
)

// add a user to FCP
func AddUser(discordUserID string, guildID string) error {
	baseURL := config.GetBaseUrl(guildID, "FCP")
	if baseURL == "" {
		err := fmt.Errorf("FCP base URL not found for guild %s", guildID)
		sentry.CaptureException(err)
		return err
	}

	url := fmt.Sprintf("%s/api/users/new", baseURL)
	payload := map[string]string{"id": discordUserID}
	body, err := json.Marshal(payload)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.FCPToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201, 208:
		return nil
	case 401:
		err := fmt.Errorf("401 when adding user %s", discordUserID)
		sentry.CaptureException(err)
		return err
	case 422:
		err := fmt.Errorf("422 when adding user %s - User ID not found in Discord", discordUserID)
		sentry.CaptureException(err)
		return err
	default:
		err := fmt.Errorf("unexpected status %d for user %s", resp.StatusCode, discordUserID)
		sentry.CaptureException(err)
		return err
	}
}

// remove a user from FCP
func RemoveUser(discordUserID string, guildID string) error {
	baseURL := config.GetBaseUrl(guildID, "FCP")
	if baseURL == "" {
		err := fmt.Errorf("FCP base URL not found for guild %s", guildID)
		sentry.CaptureException(err)
		return err
	}

	url := fmt.Sprintf("%s/api/users/find/%s/delete", baseURL, discordUserID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.FCPToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201, 204, 404:
		return nil
	case 401:
		err := fmt.Errorf("401 when removing user %s", discordUserID)
		sentry.CaptureException(err)
		return err
	default:
		err := fmt.Errorf("unexpected status %d when removing user %s", resp.StatusCode, discordUserID)
		sentry.CaptureException(err)
		return err
	}
}
