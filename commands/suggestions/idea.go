package suggestions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"

	"tpc-discord-bot/internal/config"
)

type createIdeaRequest struct {
	StatusID  int    `json:"statusId"`
	DiscordID string `json:"discordId"`
	IdeaText  string `json:"ideaText"`
}

type updateIdeaRequest struct {
	StatusID     *int   `json:"statusId,omitempty"`
	Reason       string `json:"reason,omitempty"`
	ActionUserID string `json:"actionUserId,omitempty"`
	ChannelID    string `json:"channelId,omitempty"`
	MessageID    string `json:"messageId,omitempty"`
}

type ideaAPIModel struct {
	ID        string `json:"id"`
	StatusID  int    `json:"statusId"`
	DiscordID string `json:"discordId"`
	IdeaText  string `json:"ideaText"`
	Reason    string `json:"reason"`
	ChannelID string `json:"channelId"`
	MessageID string `json:"messageId"`
	Detail    any    `json:"detail"`
}

type ideaAPIClient struct{}

func ideaIDToInt(id string) (int, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return 0, fmt.Errorf("empty idea id")
	}
	n, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("invalid idea id %q: %w", id, err)
	}
	return n, nil
}

func (c *ideaAPIClient) baseURL(guildID string) string {
	if v := strings.TrimSpace(config.GetBaseUrl(guildID, "Internal API")); v != "" {
		return strings.TrimRight(v, "/")
	}

	for _, k := range []string{"INTERNAL_API_BASE_URL", "CORE_API_BASE_URL", "INTERNAL_API_URL"} {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			return strings.TrimRight(v, "/")
		}
	}

	return ""
}

func (c *ideaAPIClient) sendJSON(method, url string, body any, guildID string, out any) error {
	if strings.TrimSpace(url) == "" || !strings.HasPrefix(url, "http") {
		return fmt.Errorf("internal api url looks wrong: %q", url)
	}

	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "TPCDiscordBot")

	apiKey := strings.TrimSpace(config.GetInternalApiKey(guildID))
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("INTERNAL_API_KEY"))
	}
	req.Header.Set("X-API-Key", apiKey)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("internal api %s %s failed: status=%d body=%s", method, url, resp.StatusCode, string(raw))
	}

	if out != nil {
		if err := json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("json decode error: %w body=%s", err, string(raw))
		}
	}

	return nil
}

func (c *ideaAPIClient) createIdea(guildID, discordID, text string) (*ideaAPIModel, error) {
	base := c.baseURL(guildID)
	if base == "" {
		return nil, fmt.Errorf("Internal API base URL is not configured (config baseurl \"Internal API\" or env INTERNAL_API_BASE_URL)")
	}

	url := fmt.Sprintf("%s/suggestions/new", base)
	req := createIdeaRequest{
		StatusID:  0,
		DiscordID: discordID,
		IdeaText:  text,
	}

	var out ideaAPIModel
	if err := c.sendJSON("POST", url, req, guildID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ideaAPIClient) patchIdea(guildID string, ideaNumber int, req updateIdeaRequest) (*ideaAPIModel, error) {
	base := c.baseURL(guildID)
	if base == "" {
		return nil, fmt.Errorf("Internal API base URL is not configured (config baseurl \"Internal API\" or env INTERNAL_API_BASE_URL)")
	}

	url := fmt.Sprintf("%s/suggestions/%d", base, ideaNumber)

	var out ideaAPIModel
	if err := c.sendJSON("PATCH", url, req, guildID, &out); err != nil {
		return &out, err
	}
	return &out, nil
}

// helper
func ptrString(v string) *string { return &v }

func displayName(m *discordgo.Member) string {
	if m == nil || m.User == nil {
		return "Unknown"
	}
	if strings.TrimSpace(m.Nick) != "" {
		return m.Nick
	}
	return m.User.Username
}

func ideaEmbed(id string, statusID int, authorName, authorIconURL, text string, color int, detailLabel, detail string) *discordgo.MessageEmbed {
	status := statusText(statusID)
	fields := []*discordgo.MessageEmbedField{{
		Name:   "Current Status",
		Value:  status,
		Inline: true,
	}}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "What This Means",
		Value: statusSummary(statusID),
	})

	if strings.TrimSpace(detail) != "" && strings.TrimSpace(detailLabel) != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  detailLabel,
			Value: detail,
		})
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Idea #%s", id),
		Description: text,
		Color:       color,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("Submitted by %s", authorName),
			IconURL: authorIconURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "TPC Suggestion Tracker",
			IconURL: "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png",
		},
		Fields: fields,
	}

	return embed
}

func startThreadFromMessage(s *discordgo.Session, channelID, messageID, name string) (*discordgo.Channel, error) {
	ts := &discordgo.ThreadStart{
		Name:                name,
		AutoArchiveDuration: 1440,
	}

	if thread, err := s.MessageThreadStartComplex(channelID, messageID, ts); err == nil {
		return thread, nil
	}
	return s.MessageThreadStart(channelID, messageID, name, 1440)
}

func IdeaCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	var text string
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "your-idea" {
			text = opt.StringValue()
			break
		}
	}
	text = strings.TrimSpace(text)

	if text == "" {
		_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: ptrString("Please include your idea text."),
		})
		return
	}

	api := &ideaAPIClient{}
	created, err := api.createIdea(i.GuildID, i.Member.User.ID, text)
	if err != nil {
		sentry.CaptureException(err)

		msg := "Submit failed: " + err.Error()
		if len(msg) > 1800 {
			msg = msg[:1800]
		}
		_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: ptrString(msg),
		})
		return
	}

	ideaBoxChannelID := strings.TrimSpace(config.GetChannelId(i.GuildID, "idea-box"))
	if ideaBoxChannelID == "" {
		_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: ptrString("I could not find the idea-box channel in the server config (channels: name: idea-box)."),
		})
		return
	}

	author := displayName(i.Member)
	authorIcon := ""
	if i.Member != nil && i.Member.User != nil {
		authorIcon = i.Member.User.AvatarURL("128")
	}

	embed := ideaEmbed(
		created.ID,
		0,
		author,
		authorIcon,
		text,
		0x37B6FF,
		"",
		"",
	)

	msg, err := s.ChannelMessageSendEmbed(ideaBoxChannelID, embed)
	if err != nil {
		sentry.CaptureException(err)
		_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: ptrString("Idea saved to the database, but I failed to post it in Discord. Check bot permissions for #idea-box."),
		})
		return
	}

	_ = s.MessageReactionAdd(msg.ChannelID, msg.ID, "⬆")
	_ = s.MessageReactionAdd(msg.ChannelID, msg.ID, "⬇")

	threadName := fmt.Sprintf("Idea #%s", created.ID)
	if thread, threadErr := startThreadFromMessage(s, msg.ChannelID, msg.ID, threadName); threadErr == nil {
		_, _ = s.ChannelMessageSend(
			thread.ID,
			fmt.Sprintf(
				"<@%s> Your idea has been submitted. We will use this thread for follow-up questions and status updates.",
				i.Member.User.ID,
			),
		)
	} else {
		sentry.CaptureException(threadErr)
	}

	ideaNum, convErr := ideaIDToInt(created.ID)
	if convErr != nil {
		sentry.CaptureException(convErr)
	} else {
		if _, err := api.patchIdea(i.GuildID, ideaNum, updateIdeaRequest{
			ChannelID: msg.ChannelID,
			MessageID: msg.ID,
		}); err != nil {
			sentry.CaptureException(err)
		}
	}

	_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: ptrString("Your idea has been posted!"),
	})
}
