package redis

import (
	"testing"
	"time"
)

func Test_applySingleOptions(t *testing.T) {
	cfg := Config{
		Address:     "127.0.0.1:6379",
		DialTimeout: 2 * time.Second,
		MaxRetries:  3,
	}
	_ = cfg.Validate()
	o := applySingleOptions(cfg, nil)
	if o.Addr != cfg.Address || o.DialTimeout != 2*time.Second || o.MaxRetries != 3 {
		t.Fatal(o)
	}
}

func Test_applyClusterOptions(t *testing.T) {
	cfg := Config{ClusterAddrs: []string{"a:7000", "b:7000"}}
	_ = cfg.Validate()
	o := applyClusterOptions(cfg, nil)
	if len(o.Addrs) != 2 {
		t.Fatal(o.Addrs)
	}
}
