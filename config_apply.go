package goredis

import (
	"crypto/tls"

	"github.com/redis/go-redis/v9"
)

func applySingleOptions(cfg Config, tlsConfig *tls.Config) *redis.Options {
	o := &redis.Options{
		Addr:         cfg.Address,
		Password:     cfg.Password,
		DB:           cfg.DB,
		TLSConfig:    tlsConfig,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		Protocol:     cfg.Protocol,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolTimeout:  cfg.PoolTimeout,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	}
	return o
}

func applyFailoverOptions(cfg Config, tlsConfig *tls.Config) *redis.FailoverOptions {
	return &redis.FailoverOptions{
		MasterName:       cfg.MasterName,
		Username:         cfg.Username,
		Password:         cfg.Password,
		DB:               cfg.DB,
		SentinelAddrs:    cfg.SentinelAddrs,
		SentinelUsername: cfg.SentinelUsername,
		SentinelPassword: cfg.SentinelPassword,
		TLSConfig:        tlsConfig,
		PoolSize:         cfg.PoolSize,
		MinIdleConns:     cfg.MinIdleConns,
		MaxRetries:       cfg.MaxRetries,
		DialTimeout:      cfg.DialTimeout,
		ReadTimeout:      cfg.ReadTimeout,
		WriteTimeout:     cfg.WriteTimeout,
		PoolTimeout:      cfg.PoolTimeout,
		ConnMaxIdleTime:  cfg.ConnMaxIdleTime,
		ConnMaxLifetime:  cfg.ConnMaxLifetime,
	}
}

func applyClusterOptions(cfg Config, tlsConfig *tls.Config) *redis.ClusterOptions {
	return &redis.ClusterOptions{
		Addrs:        cfg.ClusterAddrs,
		Username:     cfg.Username,
		Password:     cfg.Password,
		TLSConfig:    tlsConfig,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolTimeout:  cfg.PoolTimeout,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	}
}
