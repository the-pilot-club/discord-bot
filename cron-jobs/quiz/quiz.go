package quiz

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/the-pilot-club/tpcgo"

	controllers "tpc-discord-bot/controllers/levels"
	"tpc-discord-bot/internal/config"
	"tpc-discord-bot/internal/leveling"
	tpcclient "tpc-discord-bot/internal/tpc"
)

const quizChannelName = "Quiz Channel"
const quizCorrectXp = 25

type quizAnswerClient interface {
	GetCurrentQuizQuestions() (*tpcgo.QuizQuestion, error)
	GetQuizUserResponses(id, answer string) ([]*tpcgo.QuizUserResponse, error)
	ResetQuizUserResponses() (bool, error)
}

type quizAnswerDiscord interface {
	GuildMember(guildID, userID string, options ...discordgo.RequestOption) (*discordgo.Member, error)
	ChannelMessageSend(channelID, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
}

type quizAnswerDependencies struct {
	quizClient    quizAnswerClient
	discord       quizAnswerDiscord
	guildIDs      []string
	quizChannelID func(guildID string) string
	awardXp       func(guildID, userID string, amount int) error
}

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

	controller := &controllers.LeaderboardController{}
	return sendQuizAnswer(quizAnswerDependencies{
		quizClient: tpc,
		discord:    s,
		guildIDs:   config.AllGuildIDs(),
		quizChannelID: func(guildID string) string {
			return config.GetChannelId(guildID, quizChannelName)
		},
		awardXp: func(guildID, userID string, amount int) error {
			_, err := leveling.AwardXp(s, controller, guildID, userID, amount)
			return err
		},
	})
}

func sendQuizAnswer(deps quizAnswerDependencies) error {
	q, err := deps.quizClient.GetCurrentQuizQuestions()
	if err != nil || q == nil || q.ID == "" {
		log.Printf("GetCurrentQuizQuestions failed: %v", err)
		if err != nil {
			return fmt.Errorf("get current quiz question: %w", err)
		}
		return errors.New("current quiz question unavailable")
	}

	ans := strings.ToLower(strings.TrimSpace(q.CorrectAnswer))
	users, err := deps.quizClient.GetQuizUserResponses(q.ID, ans)
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

	winnerIDs := make([]string, 0, len(users))
	for _, u := range users {
		if u != nil && u.User != nil && u.User.User != nil && u.User.User.ID != "" {
			winnerIDs = append(winnerIDs, u.User.User.ID)
		}
	}
	answerText := buildAnswerText(header, winnerIDs)

	var errs []error
	for _, guildID := range deps.guildIDs {
		for _, userID := range winnerIDs {
			member, err := deps.discord.GuildMember(guildID, userID)
			if err != nil {
				if quizMemberNotFound(err) {
					continue
				}

				log.Printf("failed to verify quiz winner %s in guild %s: %v", userID, guildID, err)
				errs = append(errs, fmt.Errorf("verify quiz winner %s in guild %s: %w", userID, guildID, err))
				continue
			}
			if member == nil {
				continue
			}

			if err := deps.awardXp(guildID, userID, quizCorrectXp); err != nil {
				log.Printf("failed to award quiz XP to user %s in guild %s: %v", userID, guildID, err)
				errs = append(errs, fmt.Errorf("award quiz XP to user %s in guild %s: %w", userID, guildID, err))
			}
		}

		channelID := deps.quizChannelID(guildID)
		if channelID == "" {
			continue
		}
		if _, err := deps.discord.ChannelMessageSend(channelID, answerText); err != nil {
			log.Printf("failed to send quiz answer to guild %s: %v", guildID, err)
			errs = append(errs, fmt.Errorf("send quiz answer to guild %s: %w", guildID, err))
		}
	}

	if _, err := deps.quizClient.ResetQuizUserResponses(); err != nil {
		log.Printf("ResetQuizUserResponses error: %v", err)
		errs = append(errs, fmt.Errorf("reset quiz user responses: %w", err))
	}
	return errors.Join(errs...)
}

func buildAnswerText(header string, winnerIDs []string) string {
	if len(winnerIDs) == 0 {
		return header + "\n\nNo one got the correct answer today :("
	}

	var text strings.Builder
	text.WriteString(header)
	text.WriteString("\n\nThe following members were correct:\n")
	for _, userID := range winnerIDs {
		fmt.Fprintf(&text, "- <@%s> (+%d XP)\n", userID, quizCorrectXp)
	}
	return text.String()
}

func quizMemberNotFound(err error) bool {
	var restErr *discordgo.RESTError
	return errors.As(err, &restErr) &&
		restErr.Message != nil &&
		restErr.Message.Code == discordgo.ErrCodeUnknownMember
}
