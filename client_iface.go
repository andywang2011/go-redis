package redis

import (
	"context"

	goredislib "github.com/redis/go-redis/v9"
)

// Client abstracts Redis access (Cmdable) for tests and application wiring.
type Client interface {
	goredislib.Cmdable
	PingWithRetry(ctx context.Context) error
	Close() error
	// Redis returns the underlying *goredislib.Client when using standalone or sentinel (nil for cluster).
	Redis() *goredislib.Client
}
