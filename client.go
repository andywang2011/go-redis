package redis

import (
	"context"
	"fmt"
	"sync"

	goredislib "github.com/redis/go-redis/v9"
)

type client struct {
	goredislib.Cmdable
	raw   *goredislib.Client
	close redisCloser
	retry *retryManager

	closeOnce sync.Once
	closeErr  error
}

var _ Client = (*client)(nil)

func newClientWithConfig(ctx context.Context, cfg Config) (Client, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	uni, raw, err := connectClient(cfg)
	if err != nil {
		return nil, err
	}
	rm := newRetryManager(cfg.RetryConfig)
	if err := rm.execute(ctx, func(ctx context.Context) error {
		return pingRedis(ctx, uni)
	}); err != nil {
		_ = uni.Close()
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return &client{Cmdable: uni, raw: raw, close: uni, retry: rm}, nil
}

// NewClient creates a client and pings Redis with startup retry.
// The context bounds the startup Ping retry loop.
func NewClient(ctx context.Context, cfg Config) (Client, error) {
	return newClientWithConfig(ctx, cfg)
}

// NewClientWithRetry is an alias for NewClient (startup Ping uses RetryConfig).
func NewClientWithRetry(ctx context.Context, cfg Config) (Client, error) {
	return newClientWithConfig(ctx, cfg)
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
	c.closeOnce.Do(func() {
		c.closeErr = c.close.Close()
	})
	return c.closeErr
}

// Redis returns the underlying *goredislib.Client for standalone/sentinel (nil for cluster).
func (c *client) Redis() *goredislib.Client {
	if c == nil {
		return nil
	}
	return c.raw
}
