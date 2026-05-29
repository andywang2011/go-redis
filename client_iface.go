package goredis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Client abstracts Redis access (Cmdable) for tests and application wiring.
type Client interface {
	redis.Cmdable
	PingWithRetry(ctx context.Context) error
	Close() error
	// Redis returns the underlying *redis.Client when using standalone or sentinel (nil for cluster).
	Redis() *redis.Client
}
