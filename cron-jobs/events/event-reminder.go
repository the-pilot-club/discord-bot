package events

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"

	"tpc-discord-bot/internal/cache"
	"tpc-discord-bot/internal/config"
)

func EventReminder(s *discordgo.Session) error {
	var errs []error
	for _, guildID := range config.AllGuildIDs() {
		if config.EventRemindersEnabled(guildID) {
			if err := sendEventReminder(s, guildID); err != nil {
				errs = append(errs, fmt.Errorf("guild %s: %w", guildID, err))
			}
		}
	}
	return errors.Join(errs...)
}

func sendEventReminder(s *discordgo.Session, guildID string) error {
	events, err := s.GuildScheduledEvents(guildID, false)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	if len(events) == 0 {
		return nil
	}

	// Filter to only scheduled events (exclude active/completed/canceled)
	var scheduled []*discordgo.GuildScheduledEvent
	for _, e := range events {
		if e.Status == discordgo.GuildScheduledEventStatusScheduled {
			scheduled = append(scheduled, e)
		}
	}
	if len(scheduled) == 0 {
		return nil
	}

	// Sort by start time and get the nearest upcoming event
	sort.Slice(scheduled, func(i, j int) bool {
		return scheduled[i].ScheduledStartTime.Before(scheduled[j].ScheduledStartTime)
	})
	ne := scheduled[0]

	// Check if event starts within 1 hour
	now := time.Now().UTC()
	timeUntilEvent := ne.ScheduledStartTime.Sub(now)
	if timeUntilEvent < 0 || timeUntilEvent > time.Hour {
		return nil
	}

	// Check dedup cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("eventreminder:%s", ne.ID)
	existing, _ := cache.Get(ctx, cacheKey)
	if existing == "sent" {
		return nil
	}

	// Get target channel
	channelID := config.GetChannelId(guildID, "Crew Chat")
	if channelID == "" {
		log.Printf("No 'Crew Chat' channel configured for guild %s", guildID)
		return nil
	}

	// Group flight ping all days, ping GA flights on Tuesday and Wednesday
	groupFlights := config.GetRoleId(guildID, "Group Flights")
	gaFlights := config.GetRoleId(guildID, "GA Flights")
	day := ne.ScheduledStartTime.Weekday()

	var pings []string
	if groupFlights != "" {
		pings = append(pings, fmt.Sprintf("<@&%s>", groupFlights))
	}
	if gaFlights != "" && (day == time.Tuesday || day == time.Wednesday) {
		pings = append(pings, fmt.Sprintf("<@&%s>", gaFlights))
	}
	pingStr := strings.Join(pings, " ")

	// Resolve creator mention
	creatorMention := ""
	if ne.Creator != nil {
		creatorMention = fmt.Sprintf("<@%s>", ne.Creator.ID)
	} else if ne.CreatorID != "" {
		creatorMention = fmt.Sprintf("<@%s>", ne.CreatorID)
	}

	// Build message content
	eventURL := fmt.Sprintf("https://discord.com/events/%s/%s", guildID, ne.ID)
	var contentParts []string
	if pingStr != "" {
		contentParts = append(contentParts, pingStr)
	}
	contentParts = append(contentParts, "**The event is starting in 1 hour. See you there!**")
	if ne.Description != "" {
		contentParts = append(contentParts, ne.Description)
	}
	if creatorMention != "" {
		contentParts = append(contentParts, fmt.Sprintf("Hosted by %s", creatorMention))
	}
	contentParts = append(contentParts, eventURL)

	msg := &discordgo.MessageSend{
		Content: strings.Join(contentParts, "\n\n"),
	}

	// Add cover image as embed if available
	if ne.Image != "" {
		imageURL := fmt.Sprintf("https://cdn.discordapp.com/guild-events/%s/%s.png?size=4096", ne.ID, ne.Image)
		msg.Embeds = []*discordgo.MessageEmbed{
			{
				Image: &discordgo.MessageEmbedImage{
					URL: imageURL,
				},
			},
		}
	}

	// Send the reminder
	_, err = s.ChannelMessageSendComplex(channelID, msg)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Failed to send event reminder for '%s': %v", ne.Name, err)
		return err
	}

	// Mark as sent with 48h TTL
	_ = cache.Set(ctx, cacheKey, "sent", 48*time.Hour)
	log.Printf("Event reminder sent for '%s' in guild %s", ne.Name, guildID)
	return nil
}
