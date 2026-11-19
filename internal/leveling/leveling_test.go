package leveling

import (
	"errors"
	"reflect"
	"testing"

	controllers "tpc-discord-bot/controllers/levels"
)

type xpUpdateCall struct {
	userID  string
	level   int
	xp      int
	totalXp int
	levelXp int
	guildID string
}

type fakeLeaderboardXpStore struct {
	stats       *controllers.UserStats
	found       bool
	findErr     error
	createErr   error
	updateErr   error
	created     *controllers.UserCreate
	createGuild string
	updates     []xpUpdateCall
}

func (f *fakeLeaderboardXpStore) FindUserStats(_, _ string) (*controllers.UserStats, bool, error) {
	return f.stats, f.found, f.findErr
}

func (f *fakeLeaderboardXpStore) CreateUserRecord(data controllers.UserCreate, guildID string) error {
	f.created = &data
	f.createGuild = guildID
	return f.createErr
}

func (f *fakeLeaderboardXpStore) UpdateUserXpState(userID string, level, xp, totalXp, levelXp int, guildID string) error {
	f.updates = append(f.updates, xpUpdateCall{
		userID:  userID,
		level:   level,
		xp:      xp,
		totalXp: totalXp,
		levelXp: levelXp,
		guildID: guildID,
	})
	return f.updateErr
}

func TestAwardXpUpdatesExistingUserAndSyncsLevel(t *testing.T) {
	store := &fakeLeaderboardXpStore{
		stats: &controllers.UserStats{
			Level:     0,
			CurrentXp: 90,
			TotalXp:   90,
		},
		found: true,
	}
	var syncedLevels []int

	change, err := awardXp(store, "guild-1", "user-1", 25, func(level int) {
		syncedLevels = append(syncedLevels, level)
	})
	if err != nil {
		t.Fatalf("awardXp returned an error: %v", err)
	}

	if change.AppliedDelta != 25 || change.After.Level != 1 || change.After.CurrentXp != 15 || change.After.TotalXp != 115 || change.After.NextLevelXp != 400 {
		t.Fatalf("awardXp returned an unexpected change: %+v", change)
	}

	wantUpdate := []xpUpdateCall{{
		userID:  "user-1",
		level:   1,
		xp:      15,
		totalXp: 115,
		levelXp: 400,
		guildID: "guild-1",
	}}
	if !reflect.DeepEqual(store.updates, wantUpdate) {
		t.Fatalf("updates = %+v, want %+v", store.updates, wantUpdate)
	}
	if store.created != nil {
		t.Fatalf("awardXp created a record for an existing user: %+v", store.created)
	}
	if !reflect.DeepEqual(syncedLevels, []int{1}) {
		t.Fatalf("synced levels = %v, want [1]", syncedLevels)
	}
}

func TestAwardXpCreatesMissingUser(t *testing.T) {
	store := &fakeLeaderboardXpStore{}
	var syncedLevels []int

	change, err := awardXp(store, "guild-1", "user-1", 25, func(level int) {
		syncedLevels = append(syncedLevels, level)
	})
	if err != nil {
		t.Fatalf("awardXp returned an error: %v", err)
	}

	wantCreate := &controllers.UserCreate{
		GuildID:         "guild-1",
		UserID:          "user-1",
		MessageCount:    0,
		Xp:              25,
		TotalXp:         25,
		LevelXp:         100,
		Level:           0,
		Rank:            0,
		NoXp:            false,
		MessageLastSent: 0,
	}
	if !reflect.DeepEqual(store.created, wantCreate) {
		t.Fatalf("created record = %+v, want %+v", store.created, wantCreate)
	}
	if store.createGuild != "guild-1" {
		t.Fatalf("create guild = %q, want guild-1", store.createGuild)
	}
	if len(store.updates) != 0 {
		t.Fatalf("awardXp updated a missing user: %+v", store.updates)
	}
	if change.AppliedDelta != 25 || change.After.TotalXp != 25 {
		t.Fatalf("awardXp returned an unexpected change: %+v", change)
	}
	if !reflect.DeepEqual(syncedLevels, []int{0}) {
		t.Fatalf("synced levels = %v, want [0]", syncedLevels)
	}
}

func TestAwardXpClampsRemovalAtZero(t *testing.T) {
	store := &fakeLeaderboardXpStore{
		stats: &controllers.UserStats{CurrentXp: 10, TotalXp: 10},
		found: true,
	}

	change, err := awardXp(store, "guild-1", "user-1", -25, nil)
	if err != nil {
		t.Fatalf("awardXp returned an error: %v", err)
	}
	if change.AppliedDelta != -10 || change.After.TotalXp != 0 || change.After.CurrentXp != 0 {
		t.Fatalf("awardXp returned an unexpected clamped change: %+v", change)
	}
}

func TestAwardXpDoesNotSyncAfterPersistenceFailure(t *testing.T) {
	fetchErr := errors.New("fetch failed")
	createErr := errors.New("create failed")
	updateErr := errors.New("update failed")

	tests := []struct {
		name    string
		store   *fakeLeaderboardXpStore
		wantErr error
	}{
		{
			name:    "fetch",
			store:   &fakeLeaderboardXpStore{findErr: fetchErr},
			wantErr: fetchErr,
		},
		{
			name:    "create",
			store:   &fakeLeaderboardXpStore{createErr: createErr},
			wantErr: createErr,
		},
		{
			name: "update",
			store: &fakeLeaderboardXpStore{
				stats:     &controllers.UserStats{},
				found:     true,
				updateErr: updateErr,
			},
			wantErr: updateErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			syncCalls := 0
			_, err := awardXp(tt.store, "guild-1", "user-1", 25, func(int) {
				syncCalls++
			})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("awardXp error = %v, want %v", err, tt.wantErr)
			}
			if syncCalls != 0 {
				t.Fatalf("role sync called %d times after failure", syncCalls)
			}
		})
	}
}
