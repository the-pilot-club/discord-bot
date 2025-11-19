package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/the-pilot-club/tpcgo"
	"log"
	"strings"
)

// HandleQuizButton handles quiz option selection and makes sure user hasnt already answered
func HandleQuizButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}
	cid := i.MessageComponentData().CustomID
	if !strings.HasPrefix(cid, "quiz-option-") {
		return
	}

	var answer string
	switch cid {
	case "quiz-option-a":
		answer = "A"
	case "quiz-option-b":
		answer = "B"
	case "quiz-option-c":
		answer = "C"
	default:
		return
	}

	tpc := getTPC()
	if tpc == nil {
		return
	}

	// check if user already anbswered
	already, err := tpc.CheckUserQuizResponse(i.Member.User.ID)
	if err != nil {
		log.Printf("CheckUserQuizResponse error: %v", err)
	}
	if already {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You have already answered the question for the day. Please wait until tomorrow.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Record the user's response
	ok, err := tpc.SetQuizUserResponse(&tpcgo.QuizUserResponseSet{
		MessageID: i.Message.ID,
		UserID:    i.Member.User.ID,
		Answer:    answer,
		User:      i.Member,
	})
	if err != nil || !ok {
		log.Printf("SetQuizUserResponse error: %v", err)
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Sorry, we couldn't record your answer. Please try again.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Thanks! Your answer has been recorded. Check back here tomorrow for the answer!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
