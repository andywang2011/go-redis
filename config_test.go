package redis

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

const testCAPEM = `-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIQIRi6zePL6mKjOipn+dNuaTAKBggqhkjOPQQDAjASMRAw
DgYDVQQKEwdBY21lIENvMB4XDTE3MTAyMDE5NDMwNloXDTE4MTAyMDE5NDMwNlow
EjEQMA4GA1UEChMHQWNtZSBDbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABD0d
7VNhbWvZLWPuj/RtHFjvtJBEwOkhbN/BnnE8rnZR8+sbwnc/KhCk3FhnpHZnQz7B
5aETbbIgmuvewdjvSBSjYzBhMA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAKBggr
BgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1UdEQQiMCCCDmxvY2FsaG9zdDo1
NDUzgg4xMjcuMC4wLjE6NTQ1MzAKBggqhkjOPQQDAgNIADBFAiEA2zpJEPQyz6/l
Wf86aX6PepsntZv2GYlA5UpabfT2EZICICpJ5h/iI+i341gBmLiAFQOyTDT+/wQc
6MF9+Yw1Yy0t
-----END CERTIFICATE-----`

func TestConfig_BuildTLS_CAOnly(t *testing.T) {
	tlsCfg, err := (Config{CA: testCAPEM}).BuildTLS()
	if err != nil {
		t.Fatal(err)
	}
	if tlsCfg == nil || tlsCfg.RootCAs == nil || len(tlsCfg.Certificates) != 0 {
		t.Fatalf("unexpected tls config: %+v", tlsCfg)
	}
}

func TestConfig_BuildTLS_badCA(t *testing.T) {
	if _, err := (Config{CA: "bad"}).BuildTLS(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfig_BuildTLS_certKeyMismatch(t *testing.T) {
	if _, err := (Config{Cert: "x"}).BuildTLS(); err == nil {
		t.Fatal("expected error for cert without key")
	}
	if _, err := (Config{Key: "x", CA: testCAPEM}).BuildTLS(); err == nil {
		t.Fatal("expected error for key without cert")
	}
}

func TestConfig_BuildTLS_badKeyPair(t *testing.T) {
	if _, err := (Config{Cert: "bad", Key: "bad"}).BuildTLS(); err == nil {
		t.Fatal("expected error")
	}
}
