package tpc

import (
	"log"
	"strings"
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
		CoreAPIEnv: resolveCoreAPIEnvironment(config.CoreAPIEnv),
	})
	if err != nil {
		log.Printf("tpcgo.NewSession error: %v", err)
		return nil
	}

	session = s
	return s
}

func resolveCoreAPIEnvironment(value string) tpcgo.Environment {
	if strings.EqualFold(strings.TrimSpace(value), string(tpcgo.EnvBeta)) {
		return tpcgo.EnvBeta
	}

	return tpcgo.EnvProduction
}
