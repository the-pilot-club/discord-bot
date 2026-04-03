package suggestions

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"

	"tpc-discord-bot/internal/config"
)

const (
	ideaAdminErrorMessage  = "Something went wrong updating the idea. Please try again later."
	archivedThreadPageSize = 100
)

func IdeaAdminCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	if len(i.ApplicationCommandData().Options) == 0 {
		reply(i, s, "Missing subcommand.")
		return
	}

	sub := i.ApplicationCommandData().Options[0]
	subName := sub.Name

	ideaNumber := 0
	reason := ""

	for _, opt := range sub.Options {
		switch opt.Name {
		case "idea-number":
			ideaNumber = int(opt.IntValue())
		case "reason":
			reason = opt.StringValue()
		}
	}

	reason = strings.TrimSpace(reason)

	if ideaNumber <= 0 {
		reply(i, s, "Please provide a valid idea number.")
		return
	}

	api := &ideaAPIClient{}
	actionUserID := i.Member.User.ID

	setStatus := func(status int) updateIdeaRequest {
		st := status
		req := updateIdeaRequest{
			StatusID:     &st,
			ActionUserID: actionUserID,
		}
		if reason != "" {
			req.Reason = reason
		}
		return req
	}

	switch subName {

	case "implement":
		body, err := api.patchIdea(i.GuildID, ideaNumber, setStatus(3))
		if err != nil {
			sentry.CaptureException(err)
		}
		if err != nil || body == nil || body.Detail != nil {
			reply(i, s, "That idea does not exist or the API returned an error.")
			return
		}

		if err := archiveFlow(s, i, api, ideaNumber, body, true, reason); err != nil {
			replyIdeaAdminError(i, s, err)
			return
		}

		reply(i, s, "Idea marked as implemented and moved to the archive.")

	case "deny":
		body, err := api.patchIdea(i.GuildID, ideaNumber, setStatus(4))
		if err != nil {
			sentry.CaptureException(err)
		}
		if err != nil || body == nil || body.Detail != nil {
			reply(i, s, "That idea does not exist or the API returned an error.")
			return
		}

		if err := archiveFlow(s, i, api, ideaNumber, body, false, reason); err != nil {
			replyIdeaAdminError(i, s, err)
			return
		}

		reply(i, s, "Idea marked as denied and moved to the archive.")

	case "consider":
		body, err := api.patchIdea(i.GuildID, ideaNumber, setStatus(1))
		if err != nil {
			sentry.CaptureException(err)
		}
		if err != nil || body == nil || body.Detail != nil {
			reply(i, s, "That idea does not exist or the API returned an error.")
			return
		}

		if err := updateMessageFlow(s, i, ideaNumber, body, reason, false); err != nil {
			replyIdeaAdminError(i, s, err)
			return
		}

		reply(i, s, "Idea marked as under review.")

	case "approve":
		body, err := api.patchIdea(i.GuildID, ideaNumber, setStatus(2))
		if err != nil {
			sentry.CaptureException(err)
		}
		if err != nil || body == nil || body.Detail != nil {
			reply(i, s, "That idea does not exist or the API returned an error.")
			return
		}

		if err := updateMessageFlow(s, i, ideaNumber, body, reason, false); err != nil {
			replyIdeaAdminError(i, s, err)
			return
		}

		reply(i, s, "Idea marked as approved.")

	case "edit-reason":
		if reason == "" {
			reply(i, s, "Please provide a reason.")
			return
		}

		body, err := api.patchIdea(i.GuildID, ideaNumber, updateIdeaRequest{
			Reason:       reason,
			ActionUserID: actionUserID,
		})
		if err != nil {
			sentry.CaptureException(err)
		}
		if err != nil || body == nil || body.Detail != nil {
			reply(i, s, "That idea does not exist or the API returned an error.")
			return
		}

		if err := updateMessageFlow(s, i, ideaNumber, body, reason, true); err != nil {
			replyIdeaAdminError(i, s, err)
			return
		}

		reply(i, s, "Idea staff note updated.")

	default:
		reply(i, s, "Unknown subcommand.")
	}
}

func archiveFlow(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	api *ideaAPIClient,
	ideaNumber int,
	body *ideaAPIModel,
	implemented bool,
	reason string,
) error {

	oldMsg, err := s.ChannelMessage(body.ChannelID, body.MessageID)
	if err != nil {
		return err
	}

	threadID := findThread(s, i.GuildID, ideaNumber, body.ChannelID)
	statusID := 4
	if implemented {
		statusID = 3
	}

	archiveID := config.GetChannelId(i.GuildID, "ideabox-archive")
	if archiveID == "" {
		return fmt.Errorf("missing ideabox-archive channel")
	}

	embed := ideaEmbed(
		body.ID,
		statusID,
		ideaAuthorName(oldMsg),
		ideaAuthorIcon(oldMsg),
		body.IdeaText,
		statusColor(statusID),
		"Staff Note",
		reason,
	)

	if err := createArchiveAndUpdateRecord(
		func() (*discordgo.Message, error) {
			return s.ChannelMessageSendEmbed(archiveID, embed)
		},
		func() error {
			if threadID == "" {
				sentry.CaptureMessage(fmt.Sprintf("suggestion thread not found for idea %d in guild %s", ideaNumber, i.GuildID))
				return nil
			}
			return sendIdeaThreadUpdate(s, threadID, body.DiscordID, statusThreadUpdate(statusID, reason, false))
		},
		func(archivedMsg *discordgo.Message) error {
			_, err := api.patchIdea(i.GuildID, ideaNumber, updateIdeaRequest{
				ChannelID: archivedMsg.ChannelID,
				MessageID: archivedMsg.ID,
			})
			return err
		},
		func(archivedMsg *discordgo.Message) error {
			return s.ChannelMessageDelete(archivedMsg.ChannelID, archivedMsg.ID)
		},
	); err != nil {
		return err
	}

	if err := s.ChannelMessageDelete(body.ChannelID, body.MessageID); err != nil {
		sentry.CaptureException(err)
	}

	return nil
}

func createArchiveAndUpdateRecord(
	createArchive func() (*discordgo.Message, error),
	notifyThread func() error,
	updateRecord func(*discordgo.Message) error,
	cleanupArchive func(*discordgo.Message) error,
) error {
	archivedMsg, err := createArchive()
	if err != nil {
		return err
	}

	if err := notifyThread(); err != nil {
		return rollbackArchiveMessage(archivedMsg, cleanupArchive, err)
	}

	if err := updateRecord(archivedMsg); err != nil {
		return rollbackArchiveMessage(archivedMsg, cleanupArchive, err)
	}

	return nil
}

func rollbackArchiveMessage(archivedMsg *discordgo.Message, cleanupArchive func(*discordgo.Message) error, cause error) error {
	if archivedMsg == nil || cleanupArchive == nil {
		return cause
	}

	if err := cleanupArchive(archivedMsg); err != nil {
		return fmt.Errorf("%w; failed to clean up archive message: %v", cause, err)
	}

	return cause
}

func updateMessageFlow(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	ideaNumber int,
	body *ideaAPIModel,
	reason string,
	reasonUpdated bool,
) error {

	oldMsg, err := s.ChannelMessage(body.ChannelID, body.MessageID)
	if err != nil {
		return err
	}

	threadID := findThread(s, i.GuildID, ideaNumber, body.ChannelID)
	embed := ideaEmbed(
		body.ID,
		body.StatusID,
		ideaAuthorName(oldMsg),
		ideaAuthorIcon(oldMsg),
		body.IdeaText,
		statusColor(body.StatusID),
		"Staff Note",
		reason,
	)

	if threadID == "" {
		sentry.CaptureMessage(fmt.Sprintf("suggestion thread not found for idea %d in guild %s", ideaNumber, i.GuildID))
	} else if err := sendIdeaThreadUpdate(s, threadID, body.DiscordID, statusThreadUpdate(body.StatusID, reason, reasonUpdated)); err != nil {
		return err
	}

	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      body.MessageID,
		Channel: body.ChannelID,
		Embeds:  &[]*discordgo.MessageEmbed{embed},
	})

	return err
}

func findThread(s *discordgo.Session, guildID string, ideaNumber int, parentID string) string {
	expected := fmt.Sprintf("Idea #%d", ideaNumber)
	parentIDs := ideaThreadParentIDs(parentID, config.GetChannelId(guildID, "idea-box"))

	threads, err := s.GuildThreadsActive(guildID)
	if err != nil {
		sentry.CaptureException(err)
	} else if threadID := matchIdeaThreadID(threads.Threads, expected, parentIDs); threadID != "" {
		return threadID
	}

	for _, candidateParentID := range parentIDs {
		threadID, err := findArchivedIdeaThreadID(func(before *time.Time) (*discordgo.ThreadsList, error) {
			return s.ThreadsArchived(candidateParentID, before, archivedThreadPageSize)
		}, candidateParentID, expected)
		if err != nil {
			sentry.CaptureException(err)
			continue
		}
		if threadID != "" {
			return threadID
		}
	}

	return ""
}

func ideaThreadParentIDs(parentIDs ...string) []string {
	seen := make(map[string]struct{}, len(parentIDs))
	result := make([]string, 0, len(parentIDs))

	for _, parentID := range parentIDs {
		parentID = strings.TrimSpace(parentID)
		if parentID == "" {
			continue
		}
		if _, ok := seen[parentID]; ok {
			continue
		}
		seen[parentID] = struct{}{}
		result = append(result, parentID)
	}

	return result
}

func matchIdeaThreadID(threads []*discordgo.Channel, expected string, parentIDs []string) string {
	allowedParents := make(map[string]struct{}, len(parentIDs))
	for _, parentID := range parentIDs {
		allowedParents[parentID] = struct{}{}
	}

	for _, th := range threads {
		if th == nil || th.Name != expected {
			continue
		}
		if len(allowedParents) == 0 {
			return th.ID
		}
		if _, ok := allowedParents[th.ParentID]; ok {
			return th.ID
		}
	}

	return ""
}

func findArchivedIdeaThreadID(fetch func(before *time.Time) (*discordgo.ThreadsList, error), parentID, expected string) (string, error) {
	var before *time.Time

	for {
		threads, err := fetch(before)
		if err != nil {
			return "", err
		}
		if threads == nil {
			return "", nil
		}
		if threadID := matchIdeaThreadID(threads.Threads, expected, []string{parentID}); threadID != "" {
			return threadID, nil
		}
		if !threads.HasMore {
			return "", nil
		}

		before = nextArchivedThreadSearchBefore(threads.Threads)
		if before == nil {
			return "", nil
		}
	}
}

func nextArchivedThreadSearchBefore(threads []*discordgo.Channel) *time.Time {
	for idx := len(threads) - 1; idx >= 0; idx-- {
		th := threads[idx]
		if th == nil || th.ThreadMetadata == nil || th.ThreadMetadata.ArchiveTimestamp.IsZero() {
			continue
		}

		before := th.ThreadMetadata.ArchiveTimestamp
		return &before
	}

	return nil
}

func replyIdeaAdminError(i *discordgo.InteractionCreate, s *discordgo.Session, err error) {
	if err != nil {
		sentry.CaptureException(err)
	}
	reply(i, s, ideaAdminErrorMessage)
}

func reply(i *discordgo.InteractionCreate, s *discordgo.Session, msg string) {
	_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}

func ideaAuthorName(msg *discordgo.Message) string {
	if msg == nil || len(msg.Embeds) == 0 || msg.Embeds[0].Author == nil || strings.TrimSpace(msg.Embeds[0].Author.Name) == "" {
		return "Unknown"
	}
	return strings.TrimPrefix(msg.Embeds[0].Author.Name, "Submitted by ")
}

func ideaAuthorIcon(msg *discordgo.Message) string {
	if msg == nil || len(msg.Embeds) == 0 || msg.Embeds[0].Author == nil {
		return ""
	}
	return msg.Embeds[0].Author.IconURL
}

func sendIdeaThreadUpdate(s *discordgo.Session, threadID, discordID, update string) error {
	if threadID == "" {
		return nil
	}

	content := fmt.Sprintf("<@%s> %s", discordID, update)
	_, err := s.ChannelMessageSend(threadID, content)
	return err
}

func statusThreadUpdate(statusID int, note string, noteUpdated bool) string {
	switch {
	case noteUpdated:
		return appendStaffNote(
			"Update on your suggestion: the staff note has been updated.",
			note,
		)
	case statusID == 1:
		return appendStaffNote(
			"Update on your suggestion: it is now under review. The team is actively discussing it and deciding whether to move forward.",
			note,
		)
	case statusID == 2:
		return appendStaffNote(
			"Update on your suggestion: it has been approved! We plan to move forward with it, although implementation may happen later.",
			note,
		)
	case statusID == 3:
		return appendStaffNote(
			"Update on your suggestion: it has been implemented. The original post has been moved to the archive.",
			note,
		)
	case statusID == 4:
		return appendStaffNote(
			"Update on your suggestion: it has been denied. The original post has been moved to the archive.",
			note,
		)
	default:
		return appendStaffNote(
			"Update on your suggestion: there is a new status update.",
			note,
		)
	}
}

func appendStaffNote(msg, note string) string {
	note = strings.TrimSpace(note)
	if note == "" {
		return msg
	}
	return msg + "\n\nStaff note:\n" + note
}

func statusSummary(id int) string {
	return map[int]string{
		0: "The idea has been submitted and is waiting for staff review. We may ask follow-up questions in the thread before taking action.",
		1: "The team is actively reviewing the idea and deciding whether to move forward. This status does not mean the idea has been approved yet.",
		2: "The team has decided to move forward with the idea. It may still take time before the work is fully implemented.",
		3: "The idea has been completed and moved to the archive.",
		4: "The team decided not to move forward with the idea, and the post has been moved to the archive.",
	}[id]
}

func statusText(id int) string {
	return map[int]string{
		0: "Submitted",
		1: "Under Review",
		2: "Approved",
		3: "Implemented",
		4: "Denied",
	}[id]
}

func statusColor(id int) int {
	return map[int]int{
		0: 0x37B6FF,
		1: 0xF8DE7E,
		2: 0x55F861,
		3: 0x37B6FF,
		4: 0xD50028,
	}[id]
}
