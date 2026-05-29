package goredis

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type redisCloser interface {
	redis.Cmdable
	Close() error
}

var (
	newStandaloneClient = redis.NewClient
	newFailoverClient   = redis.NewFailoverClient
	newClusterClient    = redis.NewClusterClient
	pingRedis           = func(ctx context.Context, c redis.Cmdable) error { return c.Ping(ctx).Err() }
)

func buildClient(cfg Config, tlsConfig *tls.Config) (redisCloser, *redis.Client, error) {
	switch clientMode(cfg) {
	case "cluster":
		cc := newClusterClient(applyClusterOptions(cfg, tlsConfig))
		return cc, nil, nil
	case "sentinel":
		rc := newFailoverClient(applyFailoverOptions(cfg, tlsConfig))
		return rc, rc, nil
	default:
		rc := newStandaloneClient(applySingleOptions(cfg, tlsConfig))
		return rc, rc, nil
	}
}

func connectClient(cfg Config) (redisCloser, *redis.Client, error) {
	tlsConfig, err := cfg.BuildTLS()
	if err != nil {
		return nil, nil, fmt.Errorf("redis tls: %w", err)
	}
	return buildClient(cfg, tlsConfig)
}
