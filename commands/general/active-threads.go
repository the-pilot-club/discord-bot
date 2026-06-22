package general

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
)

const activeThreadsMessageLimit = 1900

type activeThreadEntry struct {
	ID       string
	GuildID  string
	Name     string
	ParentID string
}

func HandleActiveThreadsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	if i.GuildID == "" {
		editActiveThreadsResponse(s, i, "This command can only be used in a server.")
		return
	}

	threads, err := s.GuildThreadsActive(i.GuildID)
	if err != nil {
		sentry.CaptureException(err)
		editActiveThreadsResponse(s, i, "Could not load active threads. Please try again later.")
		return
	}

	var activeThreads []*discordgo.Channel
	if threads != nil {
		activeThreads = threads.Threads
	}

	messages := buildActiveThreadMessages(i.GuildID, activeThreads, time.Now())
	if len(messages) == 0 {
		editActiveThreadsResponse(s, i, "No active threads found.")
		return
	}

	editActiveThreadsResponse(s, i, messages[0])
	for _, msg := range messages[1:] {
		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			sentry.CaptureException(err)
			return
		}
	}
}

func editActiveThreadsResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	if err != nil {
		sentry.CaptureException(err)
	}
}

func buildActiveThreadMessages(guildID string, threads []*discordgo.Channel, generatedAt time.Time) []string {
	entries := activeThreadEntries(guildID, threads)
	header := fmt.Sprintf("**Active Threads (%d)**\nGenerated <t:%d:R>.\n\n", len(entries), generatedAt.Unix())
	if len(entries) == 0 {
		return []string{header + "No active threads found."}
	}

	lines := make([]string, 0, len(entries))
	for _, entry := range entries {
		lines = append(lines, entry.line())
	}

	return chunkActiveThreadLines(header, "**Active Threads (continued)**\n\n", lines, activeThreadsMessageLimit)
}

func activeThreadEntries(guildID string, threads []*discordgo.Channel) []activeThreadEntry {
	entries := make([]activeThreadEntry, 0, len(threads))

	for _, thread := range threads {
		if thread == nil || !thread.IsThread() {
			continue
		}
		if thread.ThreadMetadata != nil && thread.ThreadMetadata.Archived {
			continue
		}

		entries = append(entries, activeThreadEntry{
			ID:       thread.ID,
			GuildID:  guildID,
			Name:     thread.Name,
			ParentID: thread.ParentID,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		leftName := strings.ToLower(entries[i].Name)
		rightName := strings.ToLower(entries[j].Name)
		if leftName != rightName {
			return leftName < rightName
		}
		return entries[i].ID < entries[j].ID
	})

	return entries
}

func (entry activeThreadEntry) line() string {
	threadName := strings.TrimSpace(entry.Name)
	if threadName == "" {
		threadName = "Untitled thread"
	}

	line := fmt.Sprintf("- [%s](https://discord.com/channels/%s/%s)", escapeMarkdownLinkText(threadName), entry.GuildID, entry.ID)
	if entry.ParentID != "" {
		line += fmt.Sprintf(" in <#%s>", entry.ParentID)
	}
	return line
}

func chunkActiveThreadLines(firstHeader, continuedHeader string, lines []string, limit int) []string {
	messages := make([]string, 0, 1)
	current := firstHeader

	for _, line := range lines {
		line += "\n"
		if len(current)+len(line) > limit && strings.TrimSpace(current) != "" {
			messages = append(messages, strings.TrimRight(current, "\n"))
			current = continuedHeader
		}
		current += line
	}

	if strings.TrimSpace(current) != "" {
		messages = append(messages, strings.TrimRight(current, "\n"))
	}

	return messages
}

func escapeMarkdownLinkText(value string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"\n", " ",
		"\r", " ",
	)
	return replacer.Replace(value)
}
