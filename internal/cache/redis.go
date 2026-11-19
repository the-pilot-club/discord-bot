package cache

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/redis/go-redis/v9"
	"tpc-discord-bot/internal/config"
)

var (
	redisClient *redis.Client
	memoryStore map[string]string
	memoryMu    sync.RWMutex
	usingRedis  bool
)

func InitCache() {
	url := config.RedisURL
	if url == "" {
		log.Println("REDIS_URL not set, using in-memory cache")
		memoryStore = make(map[string]string)
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr: url,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		sentry.CaptureException(err)
		log.Printf("Failed to connect to Redis (%s), falling back to in-memory cache: %v", url, err)
		memoryStore = make(map[string]string)
		return
	}

	redisClient = client
	usingRedis = true
	log.Println("Connected to Redis")
}

func CloseCache() {
	if usingRedis && redisClient != nil {
		redisClient.Close()
	}
}

func Get(ctx context.Context, key string) (string, error) {
	if usingRedis {
		return redisClient.Get(ctx, key).Result()
	}
	memoryMu.RLock()
	defer memoryMu.RUnlock()
	val, ok := memoryStore[key]
	if !ok {
		return "", redis.Nil
	}
	return val, nil
}

func Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	if usingRedis {
		return redisClient.Set(ctx, key, value, expiration).Err()
	}
	memoryMu.Lock()
	defer memoryMu.Unlock()
	memoryStore[key] = value
	return nil
}
