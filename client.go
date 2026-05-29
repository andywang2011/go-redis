package goredis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type client struct {
	redis.Cmdable
	raw   *redis.Client
	close redisCloser
	retry *retryManager
}

var _ Client = (*client)(nil)

func newClientWithConfig(cfg Config) (Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	uni, raw, err := connectClient(cfg)
	if err != nil {
		return nil, err
	}
	rm := newRetryManager(cfg.RetryConfig)
	ctx := context.Background()
	if err := rm.execute(ctx, func(ctx context.Context) error {
		return pingRedis(ctx, uni)
	}); err != nil {
		_ = uni.Close()
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return &client{Cmdable: uni, raw: raw, close: uni, retry: rm}, nil
}

// NewClient creates a client and pings Redis with startup retry.
func NewClient(cfg Config) (Client, error) {
	return newClientWithConfig(cfg)
}

// NewClientWithRetry is an alias for NewClient (startup Ping uses RetryConfig).
func NewClientWithRetry(cfg Config) (Client, error) {
	return newClientWithConfig(cfg)
}

// PingWithRetry pings Redis with exponential-backoff retry.
func (c *client) PingWithRetry(ctx context.Context) error {
	if c == nil || c.Cmdable == nil {
		return ErrNotInitialized
	}
	return c.retry.execute(ctx, func(ctx context.Context) error {
		return pingRedis(ctx, c.Cmdable)
	})
}

// Close closes the client; safe to call multiple times.
func (c *client) Close() error {
	if c == nil || c.close == nil {
		return nil
	}
	return c.close.Close()
}

// Redis returns the underlying *redis.Client for standalone/sentinel (nil for cluster).
func (c *client) Redis() *redis.Client {
	if c == nil {
		return nil
	}
	return c.raw
}
