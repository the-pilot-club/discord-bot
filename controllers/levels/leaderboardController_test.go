package controllers

import (
	"strings"
	"testing"
)

func TestParseUserStats(t *testing.T) {
	stats, err := parseUserStats(map[string]interface{}{
		"level":        float64(3),
		"xp":           "45",
		"totalXp":      int64(1445),
		"messageCount": 27,
		"noXp":         "true",
	})
	if err != nil {
		t.Fatalf("parseUserStats returned an error: %v", err)
	}

	if stats.Level != 3 || stats.CurrentXp != 45 || stats.TotalXp != 1445 || stats.MessageCount != 27 || !stats.NoXp {
		t.Fatalf("parseUserStats returned unexpected stats: %+v", stats)
	}
}

func TestParseUserStatsRequiresProgressFields(t *testing.T) {
	_, err := parseUserStats(map[string]interface{}{
		"level":        1,
		"xp":           2,
		"messageCount": 3,
	})
	if err == nil {
		t.Fatal("parseUserStats returned no error for a missing totalXp field")
	}
	if !strings.Contains(err.Error(), `"totalXp"`) {
		t.Fatalf("parseUserStats error did not identify totalXp: %v", err)
	}
}
