package suggestions

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"

	"tpc-discord-bot/internal/config"
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
		if err != nil || body == nil || body.Detail != nil {
			reply(i, s, "That idea does not exist or the API returned an error.")
			return
		}

		if err := archiveFlow(s, i, api, ideaNumber, body, true, reason); err != nil {
			sentry.CaptureException(err)
			reply(i, s, err.Error())
			return
		}

		reply(i, s, "Idea marked as implemented and moved to the archive.")

	case "deny":
		body, err := api.patchIdea(i.GuildID, ideaNumber, setStatus(4))
		if err != nil || body == nil || body.Detail != nil {
			reply(i, s, "That idea does not exist or the API returned an error.")
			return
		}

		if err := archiveFlow(s, i, api, ideaNumber, body, false, reason); err != nil {
			sentry.CaptureException(err)
			reply(i, s, err.Error())
			return
		}

		reply(i, s, "Idea marked as denied and moved to the archive.")

	case "consider":
		body, err := api.patchIdea(i.GuildID, ideaNumber, setStatus(1))
		if err != nil || body == nil || body.Detail != nil {
			reply(i, s, "That idea does not exist or the API returned an error.")
			return
		}

		if err := updateMessageFlow(s, i, ideaNumber, body, reason, false); err != nil {
			sentry.CaptureException(err)
			reply(i, s, err.Error())
			return
		}

		reply(i, s, "Idea marked as under review.")

	case "approve":
		body, err := api.patchIdea(i.GuildID, ideaNumber, setStatus(2))
		if err != nil || body == nil || body.Detail != nil {
			reply(i, s, "That idea does not exist or the API returned an error.")
			return
		}

		if err := updateMessageFlow(s, i, ideaNumber, body, reason, false); err != nil {
			sentry.CaptureException(err)
			reply(i, s, err.Error())
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
		if err != nil || body == nil || body.Detail != nil {
			reply(i, s, "That idea does not exist or the API returned an error.")
			return
		}

		if err := updateMessageFlow(s, i, ideaNumber, body, reason, true); err != nil {
			sentry.CaptureException(err)
			reply(i, s, err.Error())
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

	_ = s.ChannelMessageDelete(body.ChannelID, body.MessageID)

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

	archivedMsg, err := s.ChannelMessageSendEmbed(archiveID, embed)
	if err != nil {
		return err
	}

	if err := sendIdeaThreadUpdate(s, threadID, body.DiscordID, statusThreadUpdate(statusID, reason, false)); err != nil {
		return err
	}

	_, err = api.patchIdea(i.GuildID, ideaNumber, updateIdeaRequest{
		ChannelID: archivedMsg.ChannelID,
		MessageID: archivedMsg.ID,
	})

	return err
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

	if err := sendIdeaThreadUpdate(s, threadID, body.DiscordID, statusThreadUpdate(body.StatusID, reason, reasonUpdated)); err != nil {
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

	threads, err := s.GuildThreadsActive(guildID)
	if err != nil {
		return ""
	}

	for _, th := range threads.Threads {
		if th.ParentID == parentID && th.Name == expected {
			return th.ID
		}
	}

	return ""
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
