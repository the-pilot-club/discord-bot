package quiz

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

	"tpc-discord-bot/internal/config"
	tpcclient "tpc-discord-bot/internal/tpc"
)

const quizChannelName = "Quiz Channel"

// SendQuizQuestion posts the question to the quiz channel in each configured guild.
func SendQuizQuestion(s *discordgo.Session) error {
	tpc := tpcclient.Session()
	if tpc == nil {
		return errors.New("tpc session unavailable")
	}

	q, err := tpc.GetNextQuizQuestion()
	if err != nil || q == nil || q.ID == "" {
		nextErr := err
		log.Printf("GetNextQuizQuestion failed (%v), falling back to current", err)
		q, err = tpc.GetCurrentQuizQuestions()
		if err != nil || q == nil || q.ID == "" {
			log.Printf("GetCurrentQuizQuestions failed: %v", err)
			if err != nil {
				return fmt.Errorf("get current quiz question after next failed (%v): %w", nextErr, err)
			}
			return fmt.Errorf("current quiz question unavailable after next failed (%v)", nextErr)
		}
	}

	content := "<:training_team:895480894901592074> " + q.Question +
		"\n\n🇦 " + q.OptionA + "\n🇧 " + q.OptionB + "\n🇨 " + q.OptionC

	row := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{CustomID: "quiz-option-a", Style: discordgo.SecondaryButton, Emoji: &discordgo.ComponentEmoji{Name: "🇦"}},
			discordgo.Button{CustomID: "quiz-option-b", Style: discordgo.SecondaryButton, Emoji: &discordgo.ComponentEmoji{Name: "🇧"}},
			discordgo.Button{CustomID: "quiz-option-c", Style: discordgo.SecondaryButton, Emoji: &discordgo.ComponentEmoji{Name: "🇨"}},
		},
	}

	var errs []error
	for _, guildID := range config.AllGuildIDs() {
		channelID := config.GetChannelId(guildID, quizChannelName)
		if channelID == "" {
			continue
		}

		msg, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
			Content:    content,
			Components: []discordgo.MessageComponent{row},
		})
		if err != nil {
			log.Printf("failed to send quiz question to guild %s: %v", guildID, err)
			errs = append(errs, fmt.Errorf("send quiz question to guild %s: %w", guildID, err))
			continue
		}

		if _, err := tpc.SetQuestionForResponse(msg.ID, q.ID); err != nil {
			log.Printf("SetQuestionForResponse error (guild %s): %v", guildID, err)
			errs = append(errs, fmt.Errorf("set question response for guild %s: %w", guildID, err))
		}
	}
	return errors.Join(errs...)
}

// SendQuizAnswer posts the correct answer to the quiz channel in each configured guild,
// then resets user responses.
func SendQuizAnswer(s *discordgo.Session) error {
	tpc := tpcclient.Session()
	if tpc == nil {
		return errors.New("tpc session unavailable")
	}

	q, err := tpc.GetCurrentQuizQuestions()
	if err != nil || q == nil || q.ID == "" {
		log.Printf("GetCurrentQuizQuestions failed: %v", err)
		if err != nil {
			return fmt.Errorf("get current quiz question: %w", err)
		}
		return errors.New("current quiz question unavailable")
	}

	ans := strings.ToLower(strings.TrimSpace(q.CorrectAnswer))
	users, err := tpc.GetQuizUserResponses(q.ID, ans)
	if err != nil {
		log.Printf("GetQuizUserResponses failed (questionID=%s, ans=%q): %v", q.ID, ans, err)
		return fmt.Errorf("get quiz user responses: %w", err)
	}

	var header string
	switch q.CorrectAnswer {
	case "A":
		header = "<:training_team:895480894901592074> The correct answer is ||🇦 " + q.OptionA + "||"
	case "B":
		header = "<:training_team:895480894901592074> The correct answer is ||🇧 " + q.OptionB + "||"
	case "C":
		header = "<:training_team:895480894901592074> The correct answer is ||🇨 " + q.OptionC + "||"
	default:
		header = "<:training_team:895480894901592074> The correct answer is unavailable."
	}

	buildAnswerText := func() string {
		text := header
		if len(users) == 0 {
			return text + "\n\nNo one got the correct answer today :("
		}
		text += "\n\nThe following members were correct:\n"
		for _, u := range users {
			if u != nil && u.User != nil && u.User.User != nil {
				text += "- <@" + u.User.User.ID + ">\n"
			}
		}
		return text
	}

	answerText := buildAnswerText()

	var errs []error
	for _, guildID := range config.AllGuildIDs() {
		channelID := config.GetChannelId(guildID, quizChannelName)
		if channelID == "" {
			continue
		}
		if _, err := s.ChannelMessageSend(channelID, answerText); err != nil {
			log.Printf("failed to send quiz answer to guild %s: %v", guildID, err)
			errs = append(errs, fmt.Errorf("send quiz answer to guild %s: %w", guildID, err))
		}
	}

	if _, err := tpc.ResetQuizUserResponses(); err != nil {
		log.Printf("ResetQuizUserResponses error: %v", err)
		errs = append(errs, fmt.Errorf("reset quiz user responses: %w", err))
	}
	return errors.Join(errs...)
}
