package handlers

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/the-pilot-club/tpcgo"

	"tpc-discord-bot/internal/config"
)

const quizChannelName = "Quiz Channel"

var tpcSession *tpcgo.Session

func getTPC() *tpcgo.Session {
	if tpcSession != nil {
		return tpcSession
	}
	s, err := tpcgo.NewSession(tpcgo.SessionConfig{
		CoreApiKey: config.CoreAPIToken,
	})
	if err != nil {
		log.Printf("tpcgo.NewSession error: %v", err)
		return nil
	}
	tpcSession = s
	return s
}

// SendQuizQuestion posts the question to the quiz channel in each guild that has one configured.
func SendQuizQuestion(s *discordgo.Session) {
	tpc := getTPC()
	if tpc == nil {
		return
	}

	q, err := tpc.GetNextQuizQuestion()
	if err != nil || q == nil || q.ID == "" {
		log.Printf("GetNextQuizQuestion failed (%v), falling back to current", err)
		q, err = tpc.GetCurrentQuizQuestions()
		if err != nil || q == nil || q.ID == "" {
			log.Printf("GetCurrentQuizQuestions failed: %v", err)
			return
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

	for _, g := range s.State.Guilds {
		channelID := config.GetChannelId(g.ID, quizChannelName)
		if channelID == "" {
			continue
		}

		msg, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
			Content:    content,
			Components: []discordgo.MessageComponent{row},
		})
		if err != nil {
			log.Printf("failed to send quiz question to guild %s (%s): %v", g.Name, g.ID, err)
			continue
		}

		if _, err := tpc.SetQuestionForResponse(msg.ID, q.ID); err != nil {
			log.Printf("SetQuestionForResponse error (guild %s): %v", g.ID, err)
		}
	}
}

// SendQuizAnswer posts the correct answer to the quiz channel in each guild that has one configured,
// then resets user responses.
func SendQuizAnswer(s *discordgo.Session) {
	tpc := getTPC()
	if tpc == nil {
		return
	}

	q, err := tpc.GetCurrentQuizQuestions()
	if err != nil || q == nil || q.ID == "" {
		log.Printf("GetCurrentQuizQuestions failed: %v", err)
		return
	}

	ans := strings.ToLower(strings.TrimSpace(q.CorrectAnswer))
	users, err := tpc.GetQuizUserResponses(q.ID, ans)
	if err != nil {
		log.Printf("GetQuizUserResponses failed (questionID=%s, ans=%q): %v", q.ID, ans, err)
		return
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

	for _, g := range s.State.Guilds {
		channelID := config.GetChannelId(g.ID, quizChannelName)
		if channelID == "" {
			continue
		}
		if _, err := s.ChannelMessageSend(channelID, answerText); err != nil {
			log.Printf("failed to send quiz answer to guild %s (%s): %v", g.Name, g.ID, err)
		}
	}

	if _, err := tpc.ResetQuizUserResponses(); err != nil {
		log.Printf("ResetQuizUserResponses error: %v", err)
	}
}
