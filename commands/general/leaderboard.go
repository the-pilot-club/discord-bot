package general

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"

	controllers "tpc-discord-bot/controllers/levels"
)

const (
	leaderboardCommandLimit = 10
	leaderboardCustomID     = "leaderboard:"
)

type leaderboardEntry struct {
	UserID string
	Stats  *leaderboardStats
	Rank   int
}

func HandleLeaderboardCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	editLeaderboardPage(s, i, 0)
}

func HandleLeaderboardPagination(s *discordgo.Session, i *discordgo.InteractionCreate) {
	offset, err := leaderboardOffsetFromCustomID(i.MessageComponentData().CustomID)
	if err != nil {
		sentry.CaptureException(err)
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "I couldn't read that leaderboard page button.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			sentry.CaptureException(err)
		}
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	editLeaderboardPage(s, i, offset)
}

func IsLeaderboardPaginationButton(customID string) bool {
	return strings.HasPrefix(customID, leaderboardCustomID)
}

func editLeaderboardPage(s *discordgo.Session, i *discordgo.InteractionCreate, offset int) {
	controller := &controllers.LeaderboardController{}
	page, err := controller.FindLeaderboardUsers(i.GuildID, offset, leaderboardCommandLimit)
	if err != nil {
		sentry.CaptureException(err)
		errorEmbed := buildLeaderboardErrorEmbed("I couldn't fetch leaderboard data right now. Try again later.")
		editLeaderboardResponse(s, i, "", []*discordgo.MessageEmbed{errorEmbed}, []discordgo.MessageComponent{})
		return
	}

	embed := buildLeaderboardEmbed(page, offset)
	components := leaderboardPaginationComponents(page, offset)
	editLeaderboardResponse(s, i, "", []*discordgo.MessageEmbed{embed}, components)
}

func editLeaderboardResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string, embeds []*discordgo.MessageEmbed, components []discordgo.MessageComponent) {
	edit := &discordgo.WebhookEdit{
		Content: &content,
	}
	if embeds != nil {
		edit.Embeds = &embeds
	}
	if components != nil {
		edit.Components = &components
	}

	_, err := s.InteractionResponseEdit(i.Interaction, edit)
	if err != nil {
		sentry.CaptureException(err)
	}
}

func buildLeaderboardEmbed(page *controllers.LeaderboardPage, offset int) *discordgo.MessageEmbed {
	entries := leaderboardEntries(page, offset)
	description := "No leaderboard data is available yet."
	totalCount := 0
	if page != nil {
		totalCount = page.TotalCount
	}

	fields := []*discordgo.MessageEmbedField{}
	if len(entries) > 0 {
		description = leaderboardPageDescription(offset, len(entries), totalCount)
		fields = leaderboardTableFields(entries)
	}

	return &discordgo.MessageEmbed{
		Title:       "TPC XP Leaderboard",
		Description: description,
		Color:       tpcEmbedColor,
		Fields:      fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: leaderboardPageFooter(offset, len(entries), totalCount),
		},
	}
}

func buildLeaderboardErrorEmbed(description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "TPC XP Leaderboard",
		Description: description,
		Color:       xpAdjustmentErrorColor,
	}
}

func leaderboardPaginationComponents(page *controllers.LeaderboardPage, offset int) []discordgo.MessageComponent {
	if page == nil || page.TotalCount <= leaderboardCommandLimit {
		return nil
	}

	totalPages := leaderboardTotalPages(page.TotalCount)
	currentPage := leaderboardCurrentPage(offset, totalPages)

	previousOffset := offset - leaderboardCommandLimit
	if previousOffset < 0 {
		previousOffset = 0
	}

	nextOffset := offset + leaderboardCommandLimit
	lastPageOffset := (totalPages - 1) * leaderboardCommandLimit
	if nextOffset > lastPageOffset {
		nextOffset = lastPageOffset
	}

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Previous",
					Style:    discordgo.SecondaryButton,
					CustomID: leaderboardCustomIDForOffset(previousOffset),
					Disabled: currentPage <= 1,
				},
				discordgo.Button{
					Label:    fmt.Sprintf("Page %d/%d", currentPage, totalPages),
					Style:    discordgo.SecondaryButton,
					CustomID: leaderboardCustomIDForOffset(offset),
					Disabled: true,
				},
				discordgo.Button{
					Label:    "Next",
					Style:    discordgo.PrimaryButton,
					CustomID: leaderboardCustomIDForOffset(nextOffset),
					Disabled: currentPage >= totalPages,
				},
			},
		},
	}
}

func leaderboardPageDescription(offset, entryCount, totalCount int) string {
	if entryCount == 0 || totalCount == 0 {
		return "No leaderboard data is available yet."
	}

	start := offset + 1
	end := offset + entryCount
	return fmt.Sprintf("Ranks %d-%d of %s by total XP.", start, end, formatLeaderboardInt(totalCount))
}

func leaderboardPageFooter(offset, entryCount, totalCount int) string {
	if totalCount == 0 {
		return "No ranked pilots yet"
	}

	totalPages := leaderboardTotalPages(totalCount)
	currentPage := leaderboardCurrentPage(offset, totalPages)
	start := offset + 1
	end := offset + entryCount

	return fmt.Sprintf(
		"Page %d of %d | Showing %d-%d of %s ranked pilots",
		currentPage,
		totalPages,
		start,
		end,
		formatLeaderboardInt(totalCount),
	)
}

func leaderboardTotalPages(totalCount int) int {
	if totalCount <= 0 {
		return 1
	}

	return (totalCount + leaderboardCommandLimit - 1) / leaderboardCommandLimit
}

func leaderboardCurrentPage(offset, totalPages int) int {
	currentPage := (offset / leaderboardCommandLimit) + 1
	if currentPage < 1 {
		return 1
	}
	if currentPage > totalPages {
		return totalPages
	}

	return currentPage
}

func leaderboardTableFields(entries []leaderboardEntry) []*discordgo.MessageEmbedField {
	rows := make([]string, 0, len(entries))

	for _, entry := range entries {
		rows = append(rows, fmt.Sprintf(
			"`#%d` <@%s>  |  Level **%d**  |  **%s** XP",
			entry.Rank,
			entry.UserID,
			entry.Stats.Level,
			formatLeaderboardInt(entry.Stats.TotalXp),
		))
	}

	return []*discordgo.MessageEmbedField{
		{
			Name:  "Leaderboard",
			Value: strings.Join(rows, "\n"),
		},
	}
}

func leaderboardEntries(page *controllers.LeaderboardPage, offset int) []leaderboardEntry {
	if page == nil {
		return nil
	}

	entries := make([]leaderboardEntry, 0, len(page.Items))
	for _, item := range page.Items {
		userID, err := leaderboardUserID(item)
		if err != nil {
			sentry.CaptureException(err)
			continue
		}

		stats, err := parseLeaderboardStats(item)
		if err != nil {
			sentry.CaptureException(err)
			continue
		}

		rank := stats.Rank
		if rank <= 0 {
			rank = offset + len(entries) + 1
		}

		entries = append(entries, leaderboardEntry{
			UserID: userID,
			Stats:  stats,
			Rank:   rank,
		})
	}

	return entries
}

func leaderboardOffsetFromCustomID(customID string) (int, error) {
	offsetString := strings.TrimPrefix(customID, leaderboardCustomID)
	if offsetString == customID || offsetString == "" {
		return 0, fmt.Errorf("invalid leaderboard custom id %q", customID)
	}

	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		return 0, fmt.Errorf("invalid leaderboard offset %q: %w", offsetString, err)
	}
	if offset < 0 {
		return 0, nil
	}

	return offset, nil
}

func leaderboardCustomIDForOffset(offset int) string {
	if offset < 0 {
		offset = 0
	}

	return fmt.Sprintf("%s%d", leaderboardCustomID, offset)
}

func leaderboardUserID(user map[string]interface{}) (string, error) {
	value, ok := user["userId"]
	if !ok || value == nil {
		return "", fmt.Errorf("leaderboard field %q missing", "userId")
	}

	id, ok := value.(string)
	if !ok || id == "" {
		return "", fmt.Errorf("leaderboard field %q has unsupported type %T", "userId", value)
	}

	return id, nil
}

func formatLeaderboardInt(value int) string {
	str := strconv.Itoa(value)
	if len(str) <= 3 {
		return str
	}

	var builder strings.Builder
	prefixLen := len(str) % 3
	if prefixLen == 0 {
		prefixLen = 3
	}

	builder.WriteString(str[:prefixLen])
	for index := prefixLen; index < len(str); index += 3 {
		builder.WriteString(",")
		builder.WriteString(str[index : index+3])
	}

	return builder.String()
}

func leaderboardLinkComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "TPC Leaderboard",
					Style: discordgo.LinkButton,
					URL:   "https://mee6.xyz/thepilotclub",
				},
			},
		},
	}
}
