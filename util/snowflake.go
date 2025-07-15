package util

import (
	"fmt"
	"strconv"
	"time"
)

const DiscordEpoch int64 = 1420070400000

func SnowflakeToTime(snowflakeID string) (time.Time, error) {
	id, err := strconv.ParseInt(snowflakeID, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse snowflake ID: %w", err)
	}

	// Extract the timestamp by right-shifting 22 bits
	timestampMillis := (id >> 22) + DiscordEpoch

	// Convert Unix milliseconds to time.Time
	t := time.UnixMilli(timestampMillis)
	return t, nil
}
