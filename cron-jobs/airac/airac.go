package airac

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"tpc-discord-bot/internal/cache"
	"tpc-discord-bot/internal/config"
)

const referenceDate = "2024-01-25"
const cycleDays = 28

const tpcIconURL = "https://static1.squarespace.com/static/614689d3918044012d2ac1b4/t/616ff36761fabc72642806e3/1634726781251/TPC_FullColor_TransparentBg_1280x1024_72dpi.png"
const airacDescription = "AIRAC (Aeronautical Information Regulation and Control) cycles contain updated navigation data including waypoints, airways, procedures, and airport information."
const embedColor = 3651327

func AiracReminder(s *discordgo.Session) {
	now := time.Now().UTC().Truncate(24 * time.Hour)

	currentEffective, currentID := getCurrentCycle(now)
	nextEffective, nextID := getNextCycle(currentEffective)

	var key string
	var embed *discordgo.MessageEmbed

	if now.Equal(currentEffective) {
		key = currentID + "-dayof"
		embed = buildDayOfEmbed(currentID, currentEffective, nextID, nextEffective)
	} else if now.Equal(nextEffective.AddDate(0, 0, -1)) {
		key = nextID + "-daybefore"
		embed = buildDayBeforeEmbed(nextID)
	} else {
		return
	}

	ctx := context.Background()
	lastSentKey, _ := cache.Get(ctx, "airac:lastSentKey")
	if lastSentKey == key {
		return
	}

	for _, guild := range s.State.Guilds {
		channelID := config.GetChannelId(guild.ID, "Crew Chat")
		if channelID == "" {
			continue
		}
		_, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{embed},
		})
		if err != nil {
			sentry.CaptureException(err)
		}
	}

	cache.Set(ctx, "airac:lastSentKey", key, 0)
}

func getCurrentCycle(now time.Time) (effectiveDate time.Time, identifier string) {
	ref, _ := time.Parse("2006-01-02", referenceDate)
	effective := ref
	for effective.AddDate(0, 0, cycleDays).Compare(now) <= 0 {
		effective = effective.AddDate(0, 0, cycleDays)
	}
	return effective, cycleIdentifier(effective)
}

func getNextCycle(currentEffective time.Time) (effectiveDate time.Time, identifier string) {
	next := currentEffective.AddDate(0, 0, cycleDays)
	return next, cycleIdentifier(next)
}

func cycleIdentifier(effectiveDate time.Time) string {
	ref, _ := time.Parse("2006-01-02", referenceDate)
	year := effectiveDate.Year()

	// Find the first cycle whose effective date falls in this year
	firstOfYear := ref
	for firstOfYear.Year() < year {
		firstOfYear = firstOfYear.AddDate(0, 0, cycleDays)
	}

	// Count from firstOfYear to effectiveDate
	nn := 1
	d := firstOfYear
	for d.Before(effectiveDate) {
		d = d.AddDate(0, 0, cycleDays)
		nn++
	}

	return fmt.Sprintf("%02d%02d", year%100, nn)
}

func buildDayOfEmbed(currentID string, currentEffective time.Time, nextID string, nextEffective time.Time) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "\u2708\ufe0f AIRAC Update Reminder",
		Description: fmt.Sprintf("**AIRAC %s** is now effective! Time to update your aircraft navigation data.", currentID),
		Color:       embedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "\U0001f4c5 Current AIRAC Cycle",
				Value: fmt.Sprintf("**AIRAC %s** (Effective: %s)", currentID, currentEffective.Format("Jan 02, 2006")),
			},
			{
				Name:  "\U0001f504 Next AIRAC Cycle",
				Value: fmt.Sprintf("**AIRAC %s** (Effective: %s)", nextID, nextEffective.Format("Jan 02, 2006")),
			},
			{
				Name:  "\u2753 What is AIRAC?",
				Value: airacDescription,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "The Pilot Club Bot \u2022 AIRAC Reminder",
			IconURL: tpcIconURL,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func buildDayBeforeEmbed(tomorrowID string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "\u2708\ufe0f AIRAC Update Reminder",
		Description: fmt.Sprintf("**AIRAC %s** goes effective tomorrow. Time to prepare your aircraft navigation data.", tomorrowID),
		Color:       embedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "\u2753 What is AIRAC?",
				Value: airacDescription,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "The Pilot Club Bot \u2022 AIRAC Reminder",
			IconURL: tpcIconURL,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
