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

		reply(i, s, "Idea set as implemented and archived")

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

		reply(i, s, "Idea set as denied and archived")

	case "consider":
		body, err := api.patchIdea(i.GuildID, ideaNumber, setStatus(1))
		if err != nil || body == nil || body.Detail != nil {
			reply(i, s, "That idea does not exist or the API returned an error.")
			return
		}

		if err := updateMessageFlow(s, i, ideaNumber, body, "Considered", 0xF8DE7E, reason, true); err != nil {
			sentry.CaptureException(err)
			reply(i, s, err.Error())
			return
		}

		reply(i, s, "Idea set as considered")

	case "approve":
		body, err := api.patchIdea(i.GuildID, ideaNumber, setStatus(2))
		if err != nil || body == nil || body.Detail != nil {
			reply(i, s, "That idea does not exist or the API returned an error.")
			return
		}

		if err := updateMessageFlow(s, i, ideaNumber, body, "Approved", 0x55F861, reason, false); err != nil {
			sentry.CaptureException(err)
			reply(i, s, err.Error())
			return
		}

		reply(i, s, "Idea approved")

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

		if err := updateMessageFlow(s, i, ideaNumber, body, statusText(body.StatusID), statusColor(body.StatusID), reason, false); err != nil {
			sentry.CaptureException(err)
			reply(i, s, err.Error())
			return
		}

		reply(i, s, "Idea reason updated")

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
	if threadID != "" {
		msg := ""
		if implemented {
			msg = fmt.Sprintf("<@%s> we have implemented your idea!", body.DiscordID)
		} else {
			msg = fmt.Sprintf("<@%s> we have denied your idea.", body.DiscordID)
		}
		if reason != "" {
			msg += "\n\n``" + reason + "``"
		}
		_, _ = s.ChannelMessageSend(threadID, msg)
	}

	_ = s.ChannelMessageDelete(body.ChannelID, body.MessageID)

	archiveID := config.GetChannelId(i.GuildID, "ideabox-archive")
	if archiveID == "" {
		return fmt.Errorf("missing ideabox-archive channel")
	}

	title := "Denied"
	color := 0xD50028
	if implemented {
		title = "Implemented"
		color = 0x37B6FF
	}

	embed := ideaEmbed(
		fmt.Sprintf("Idea #%s **%s**", body.ID, title),
		oldMsg.Embeds[0].Author.Name,
		oldMsg.Embeds[0].Author.IconURL,
		body.IdeaText,
		color,
		"Reason",
		reason,
	)

	archivedMsg, err := s.ChannelMessageSendEmbed(archiveID, embed)
	if err != nil {
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
	status string,
	color int,
	reason string,
	consider bool,
) error {

	oldMsg, err := s.ChannelMessage(body.ChannelID, body.MessageID)
	if err != nil {
		return err
	}

	threadID := findThread(s, i.GuildID, ideaNumber, body.ChannelID)
	if threadID != "" {
		msg := fmt.Sprintf("<@%s> your idea has been %s.", body.DiscordID, strings.ToLower(status))
		if reason != "" {
			msg += "\n\n``" + reason + "``"
		}
		_, _ = s.ChannelMessageSend(threadID, msg)
	}

	embed := ideaEmbed(
		fmt.Sprintf("Idea #%s **%s**", body.ID, status),
		oldMsg.Embeds[0].Author.Name,
		oldMsg.Embeds[0].Author.IconURL,
		body.IdeaText,
		color,
		"Reason",
		reason,
	)

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

func statusText(id int) string {
	return map[int]string{
		0: "Submitted",
		1: "Considered",
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