package goredis

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"
)

const (
	defaultPoolSize     = 50
	defaultMinIdleConns = 10
	defaultMaxRetries   = 2
	defaultProtocol     = 2
)

// Config holds Redis connection settings for standalone, Sentinel, or Cluster.
type Config struct {
	Address      string
	ClusterAddrs []string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	Cert, Key, CA string
	MasterName       string
	SentinelAddrs    []string
	Username         string
	SentinelUsername string
	SentinelPassword string
	DialTimeout      time.Duration
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	PoolTimeout      time.Duration
	ConnMaxIdleTime  time.Duration
	ConnMaxLifetime  time.Duration
	MaxRetries       int
	Protocol         int
	RetryConfig      RetryConfig
}

// RetryConfig controls exponential-backoff retry during startup Ping and PingWithRetry.
type RetryConfig struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// ConfigDefault returns production-oriented defaults (Address or Sentinel/Cluster must still be set).
func ConfigDefault() Config {
	return Config{
		PoolSize:     defaultPoolSize,
		MinIdleConns: defaultMinIdleConns,
		MaxRetries:   defaultMaxRetries,
		Protocol:     defaultProtocol,
		RetryConfig:  DefaultRetryConfig(),
	}
}

// DefaultRetryConfig returns startup retry defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
	}
}

// Validate checks config and applies defaults without connecting.
func (c *Config) Validate() error {
	modes := 0
	if c.Address != "" {
		modes++
	}
	if c.MasterName != "" {
		modes++
	}
	if len(c.ClusterAddrs) > 0 {
		modes++
	}
	if modes == 0 {
		return fmt.Errorf("%w: address, sentinel (master_name), or cluster_addrs is required", ErrInvalidConfig)
	}
	if modes > 1 {
		return fmt.Errorf("%w: only one of address, sentinel, or cluster_addrs may be set", ErrInvalidConfig)
	}
	if c.MasterName != "" && len(c.SentinelAddrs) == 0 {
		return fmt.Errorf("%w: sentineladdrs required for failover", ErrInvalidConfig)
	}
	if c.PoolSize <= 0 {
		c.PoolSize = defaultPoolSize
	}
	if c.MinIdleConns <= 0 {
		c.MinIdleConns = defaultMinIdleConns
	}
	if c.MaxRetries <= 0 {
		c.MaxRetries = defaultMaxRetries
	}
	if c.Protocol <= 0 {
		c.Protocol = defaultProtocol
	}
	if c.RetryConfig.MaxAttempts <= 0 {
		c.RetryConfig = DefaultRetryConfig()
	} else {
		normalizeRetry(&c.RetryConfig)
	}
	return nil
}

// BuildTLS builds TLS config when Cert, Key, and CA are all set.
func (c Config) BuildTLS() (*tls.Config, error) {
	if c.Cert == "" || c.Key == "" || c.CA == "" {
		return nil, nil
	}
	certificate, err := tls.X509KeyPair([]byte(c.Cert), []byte(c.Key))
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM([]byte(c.CA)) {
		return nil, fmt.Errorf("%w: failed to parse ca", ErrInvalidConfig)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      pool,
	}, nil
}

func clientMode(cfg Config) string {
	if len(cfg.ClusterAddrs) > 0 {
		return "cluster"
	}
	if cfg.MasterName != "" {
		return "sentinel"
	}
	return "standalone"
}
