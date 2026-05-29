package goredis

import (
	"testing"
	"time"
)

func TestConfig_Validate_standalone(t *testing.T) {
	cfg := Config{Address: "127.0.0.1:6379"}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	if cfg.PoolSize != defaultPoolSize {
		t.Fatalf("PoolSize=%d", cfg.PoolSize)
	}
}

func TestConfig_Validate_errors(t *testing.T) {
	cfgEmpty := Config{}
	if err := cfgEmpty.Validate(); err == nil {
		t.Fatal("expected error")
	}
	cfg := Config{MasterName: "m"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected sentinel error")
	}
	cfg2 := Config{Address: "a", ClusterAddrs: []string{"c"}}
	if err := cfg2.Validate(); err == nil {
		t.Fatal("expected exclusive mode error")
	}
}

func TestConfig_Validate_cluster(t *testing.T) {
	cfg := Config{ClusterAddrs: []string{"127.0.0.1:7000"}}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestConfig_Validate_sentinel(t *testing.T) {
	cfg := Config{MasterName: "m", SentinelAddrs: []string{"127.0.0.1:26379"}}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestConfig_BuildTLS_Empty(t *testing.T) {
	tlsCfg, err := (Config{}).BuildTLS()
	if err != nil || tlsCfg != nil {
		t.Fatal(tlsCfg, err)
	}
}

func TestDefaultRetryConfig(t *testing.T) {
	c := DefaultRetryConfig()
	if c.MaxAttempts != 3 {
		t.Fatal(c)
	}
}

func TestConfigDefault(t *testing.T) {
	c := ConfigDefault()
	if c.PoolSize != defaultPoolSize {
		t.Fatal(c.PoolSize)
	}
}

func TestConfig_Validate_retry(t *testing.T) {
	cfg := Config{
		Address: "127.0.0.1:6379",
		RetryConfig: RetryConfig{
			MaxAttempts: 2, InitialDelay: time.Millisecond, MaxDelay: time.Second, Multiplier: 2,
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
}

func Test_clientMode(t *testing.T) {
	if clientMode(Config{ClusterAddrs: []string{"a"}}) != "cluster" {
		t.Fatal()
	}
	if clientMode(Config{MasterName: "m"}) != "sentinel" {
		t.Fatal()
	}
	if clientMode(Config{Address: "a"}) != "standalone" {
		t.Fatal()
	}
}
