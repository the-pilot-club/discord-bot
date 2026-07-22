package quiz

import (
	"errors"
	"reflect"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/the-pilot-club/tpcgo"
)

type fakeQuizAnswerClient struct {
	question           *tpcgo.QuizQuestion
	questionErr        error
	responses          []*tpcgo.QuizUserResponse
	responsesErr       error
	responseQuestionID string
	responseAnswer     string
	resetErr           error
	resetCalls         int
}

func (f *fakeQuizAnswerClient) GetCurrentQuizQuestions() (*tpcgo.QuizQuestion, error) {
	return f.question, f.questionErr
}

func (f *fakeQuizAnswerClient) GetQuizUserResponses(id, answer string) ([]*tpcgo.QuizUserResponse, error) {
	f.responseQuestionID = id
	f.responseAnswer = answer
	return f.responses, f.responsesErr
}

func (f *fakeQuizAnswerClient) ResetQuizUserResponses() (bool, error) {
	f.resetCalls++
	return f.resetErr == nil, f.resetErr
}

type quizMemberResult struct {
	member *discordgo.Member
	err    error
}

type sentQuizMessage struct {
	channelID string
	content   string
}

type fakeQuizAnswerDiscord struct {
	members     map[string]quizMemberResult
	memberCalls []string
	sendErrs    map[string]error
	messages    []sentQuizMessage
}

func (f *fakeQuizAnswerDiscord) GuildMember(guildID, userID string, _ ...discordgo.RequestOption) (*discordgo.Member, error) {
	key := guildID + "/" + userID
	f.memberCalls = append(f.memberCalls, key)
	result := f.members[key]
	return result.member, result.err
}

func (f *fakeQuizAnswerDiscord) ChannelMessageSend(channelID, content string, _ ...discordgo.RequestOption) (*discordgo.Message, error) {
	f.messages = append(f.messages, sentQuizMessage{channelID: channelID, content: content})
	if err := f.sendErrs[channelID]; err != nil {
		return nil, err
	}
	return &discordgo.Message{ID: "message-" + channelID}, nil
}

type quizAwardCall struct {
	guildID string
	userID  string
	amount  int
}

func quizResponse(userID string) *tpcgo.QuizUserResponse {
	return &tpcgo.QuizUserResponse{
		User: &discordgo.Member{User: &discordgo.User{ID: userID}},
	}
}

func unknownMemberError() error {
	return &discordgo.RESTError{
		Message: &discordgo.APIErrorMessage{Code: discordgo.ErrCodeUnknownMember},
	}
}

func TestBuildAnswerTextIncludesQuizXp(t *testing.T) {
	got := buildAnswerText("The correct answer is A", []string{"user-1", "user-2"})
	want := "The correct answer is A\n\nThe following members were correct:\n" +
		"- <@user-1> (+25 XP)\n" +
		"- <@user-2> (+25 XP)\n"

	if got != want {
		t.Fatalf("buildAnswerText() = %q, want %q", got, want)
	}
}

func TestBuildAnswerTextWithNoWinners(t *testing.T) {
	got := buildAnswerText("The correct answer is A", nil)
	want := "The correct answer is A\n\nNo one got the correct answer today :("

	if got != want {
		t.Fatalf("buildAnswerText() = %q, want %q", got, want)
	}
}

func TestQuizMemberNotFound(t *testing.T) {
	notFound := unknownMemberError()
	if !quizMemberNotFound(notFound) {
		t.Fatal("quizMemberNotFound returned false for Discord's unknown-member error")
	}

	if quizMemberNotFound(errors.New("temporary Discord failure")) {
		t.Fatal("quizMemberNotFound returned true for an unrelated error")
	}
}

func TestSendQuizAnswerAwardsMembersAcrossGuildsOffline(t *testing.T) {
	quizClient := &fakeQuizAnswerClient{
		question: &tpcgo.QuizQuestion{
			ID:            "question-1",
			CorrectAnswer: "A",
			OptionA:       "Answer A",
		},
		responses: []*tpcgo.QuizUserResponse{
			quizResponse("user-1"),
			quizResponse("user-2"),
		},
	}
	discord := &fakeQuizAnswerDiscord{
		members: map[string]quizMemberResult{
			"guild-1/user-1": {member: &discordgo.Member{User: &discordgo.User{ID: "user-1"}}},
			"guild-1/user-2": {err: unknownMemberError()},
			"guild-2/user-1": {},
			"guild-2/user-2": {member: &discordgo.Member{User: &discordgo.User{ID: "user-2"}}},
		},
	}
	var awards []quizAwardCall

	err := sendQuizAnswer(quizAnswerDependencies{
		quizClient: quizClient,
		discord:    discord,
		guildIDs:   []string{"guild-1", "guild-2"},
		quizChannelID: func(guildID string) string {
			return "channel-" + guildID
		},
		awardXp: func(guildID, userID string, amount int) error {
			awards = append(awards, quizAwardCall{guildID: guildID, userID: userID, amount: amount})
			return nil
		},
	})
	if err != nil {
		t.Fatalf("sendQuizAnswer returned an error: %v", err)
	}

	wantAwards := []quizAwardCall{
		{guildID: "guild-1", userID: "user-1", amount: 25},
		{guildID: "guild-2", userID: "user-2", amount: 25},
	}
	if !reflect.DeepEqual(awards, wantAwards) {
		t.Fatalf("awards = %+v, want %+v", awards, wantAwards)
	}
	wantMemberCalls := []string{
		"guild-1/user-1",
		"guild-1/user-2",
		"guild-2/user-1",
		"guild-2/user-2",
	}
	if !reflect.DeepEqual(discord.memberCalls, wantMemberCalls) {
		t.Fatalf("member calls = %v, want %v", discord.memberCalls, wantMemberCalls)
	}
	if len(discord.messages) != 2 {
		t.Fatalf("sent %d messages, want 2", len(discord.messages))
	}
	for _, message := range discord.messages {
		wantContent := buildAnswerText(
			"<:training_team:895480894901592074> The correct answer is ||🇦 Answer A||",
			[]string{"user-1", "user-2"},
		)
		if message.content != wantContent {
			t.Fatalf("message content = %q, want %q", message.content, wantContent)
		}
	}
	if quizClient.responseQuestionID != "question-1" || quizClient.responseAnswer != "a" {
		t.Fatalf("response query = (%q, %q), want (question-1, a)", quizClient.responseQuestionID, quizClient.responseAnswer)
	}
	if quizClient.resetCalls != 1 {
		t.Fatalf("reset calls = %d, want 1", quizClient.resetCalls)
	}
}

func TestSendQuizAnswerContinuesAndResetsAfterPerGuildFailures(t *testing.T) {
	membershipErr := errors.New("membership failed")
	awardErr := errors.New("award failed")
	sendErr := errors.New("send failed")
	resetErr := errors.New("reset failed")
	quizClient := &fakeQuizAnswerClient{
		question: &tpcgo.QuizQuestion{ID: "question-1", CorrectAnswer: "B", OptionB: "Answer B"},
		responses: []*tpcgo.QuizUserResponse{
			quizResponse("user-1"),
			quizResponse("user-2"),
			quizResponse("user-3"),
		},
		resetErr: resetErr,
	}
	discord := &fakeQuizAnswerDiscord{
		members: map[string]quizMemberResult{
			"guild-1/user-1": {err: membershipErr},
			"guild-1/user-2": {member: &discordgo.Member{User: &discordgo.User{ID: "user-2"}}},
			"guild-1/user-3": {member: &discordgo.Member{User: &discordgo.User{ID: "user-3"}}},
		},
		sendErrs: map[string]error{"channel-1": sendErr},
	}
	var awards []quizAwardCall

	err := sendQuizAnswer(quizAnswerDependencies{
		quizClient: quizClient,
		discord:    discord,
		guildIDs:   []string{"guild-1"},
		quizChannelID: func(string) string {
			return "channel-1"
		},
		awardXp: func(guildID, userID string, amount int) error {
			awards = append(awards, quizAwardCall{guildID: guildID, userID: userID, amount: amount})
			if userID == "user-2" {
				return awardErr
			}
			return nil
		},
	})

	for _, wantErr := range []error{membershipErr, awardErr, sendErr, resetErr} {
		if !errors.Is(err, wantErr) {
			t.Errorf("sendQuizAnswer error %v does not contain %v", err, wantErr)
		}
	}
	wantAwards := []quizAwardCall{
		{guildID: "guild-1", userID: "user-2", amount: 25},
		{guildID: "guild-1", userID: "user-3", amount: 25},
	}
	if !reflect.DeepEqual(awards, wantAwards) {
		t.Fatalf("awards = %+v, want %+v", awards, wantAwards)
	}
	if len(discord.messages) != 1 {
		t.Fatalf("message attempts = %d, want 1", len(discord.messages))
	}
	if quizClient.resetCalls != 1 {
		t.Fatalf("reset calls = %d, want 1", quizClient.resetCalls)
	}
}

func TestSendQuizAnswerStopsBeforeResetOnRetrievalFailure(t *testing.T) {
	questionErr := errors.New("question failed")
	responsesErr := errors.New("responses failed")

	tests := []struct {
		name     string
		client   *fakeQuizAnswerClient
		wantErr  error
		wantText string
	}{
		{
			name:    "question error",
			client:  &fakeQuizAnswerClient{questionErr: questionErr},
			wantErr: questionErr,
		},
		{
			name:     "question unavailable",
			client:   &fakeQuizAnswerClient{},
			wantText: "current quiz question unavailable",
		},
		{
			name: "responses error",
			client: &fakeQuizAnswerClient{
				question:     &tpcgo.QuizQuestion{ID: "question-1", CorrectAnswer: "A"},
				responsesErr: responsesErr,
			},
			wantErr: responsesErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discord := &fakeQuizAnswerDiscord{}
			err := sendQuizAnswer(quizAnswerDependencies{
				quizClient: tt.client,
				discord:    discord,
				quizChannelID: func(string) string {
					return "channel-1"
				},
				awardXp: func(string, string, int) error { return nil },
			})

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("sendQuizAnswer error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantText != "" && (err == nil || err.Error() != tt.wantText) {
				t.Fatalf("sendQuizAnswer error = %v, want %q", err, tt.wantText)
			}
			if tt.client.resetCalls != 0 {
				t.Fatalf("reset calls = %d, want 0", tt.client.resetCalls)
			}
			if len(discord.messages) != 0 {
				t.Fatalf("message attempts = %d, want 0", len(discord.messages))
			}
		})
	}
}
