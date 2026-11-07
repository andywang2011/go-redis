package redis

import (
	"context"
	"crypto/tls"
	"fmt"

	goredislib "github.com/redis/go-redis/v9"
)

type redisCloser interface {
	goredislib.Cmdable
	Close() error
}

var (
	newStandaloneClient = goredislib.NewClient
	newFailoverClient   = goredislib.NewFailoverClient
	newClusterClient    = goredislib.NewClusterClient
	pingRedis           = func(ctx context.Context, c goredislib.Cmdable) error { return c.Ping(ctx).Err() }
)

func buildClient(cfg Config, tlsConfig *tls.Config) (redisCloser, *goredislib.Client, error) {
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

func connectClient(cfg Config) (redisCloser, *goredislib.Client, error) {
	tlsConfig, err := cfg.BuildTLS()
	if err != nil {
		return nil, nil, fmt.Errorf("redis tls: %w", err)
	}
	return buildClient(cfg, tlsConfig)
}
