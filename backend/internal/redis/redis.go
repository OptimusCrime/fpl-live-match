// Package redis provides a Redis client constructor configured from environment variables.
package redis

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

// New creates and pings a Redis client using REDIS_HOST, REDIS_PORT, and REDIS_PASSWORD env vars.
func New() (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6380"))

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis: failed to connect to %s: %w", addr, err)
	}

	slog.Info("connected to redis", "addr", addr)
	return client, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
