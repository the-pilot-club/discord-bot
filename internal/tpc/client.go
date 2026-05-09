package tpc

import (
	"log"
	"sync"

	"github.com/the-pilot-club/tpcgo"

	"tpc-discord-bot/internal/config"
)

var (
	session *tpcgo.Session
	mu      sync.Mutex
)

func Session() *tpcgo.Session {
	mu.Lock()
	defer mu.Unlock()

	if session != nil {
		return session
	}

	s, err := tpcgo.NewSession(tpcgo.SessionConfig{
		CoreApiKey: config.CoreAPIToken,
	})
	if err != nil {
		log.Printf("tpcgo.NewSession error: %v", err)
		return nil
	}

	session = s
	return s
}
