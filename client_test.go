package redis

import (
	"context"
	"errors"
	"testing"

	"github.com/alicebob/miniredis/v2"
	goredislib "github.com/redis/go-redis/v9"
)

func TestNewClient_miniredis(t *testing.T) {
	mr := miniredis.RunT(t)
	cfg := Config{Address: mr.Addr()}
	cli, err := NewClient(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	if err := cli.Set(context.Background(), "k", "v", 0).Err(); err != nil {
		t.Fatal(err)
	}
}

func TestNewClient_validateError(t *testing.T) {
	_, err := NewClient(context.Background(), Config{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestClient_Close_nil(t *testing.T) {
	var c *client
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestClient_PingWithRetry_uninitialized(t *testing.T) {
	var c *client
	if err := c.PingWithRetry(context.Background()); err != ErrNotInitialized {
		t.Fatal(err)
	}
}

func TestClient_PingWithRetry(t *testing.T) {
	mr := miniredis.RunT(t)
	cli, err := NewClient(context.Background(), Config{Address: mr.Addr()})
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	if err := cli.PingWithRetry(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestNewClient_pingFailure(t *testing.T) {
	old := pingRedis
	defer func() { pingRedis = old }()
	pingRedis = func(context.Context, goredislib.Cmdable) error { return errors.New("ping fail") }
	mr := miniredis.RunT(t)
	_, err := NewClient(context.Background(), Config{Address: mr.Addr()})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewClientWithRetry(t *testing.T) {
	mr := miniredis.RunT(t)
	cli, err := NewClientWithRetry(context.Background(), Config{Address: mr.Addr()})
	if err != nil {
		t.Fatal(err)
	}
	cli.Close()
}

func TestClient_Redis_nilForMissing(t *testing.T) {
	var c *client
	if c.Redis() != nil {
		t.Fatal("expected nil")
	}
}

func TestNewClient_clusterMode(t *testing.T) {
	mr := miniredis.RunT(t)
	oldC := newClusterClient
	oldP := pingRedis
	defer func() {
		newClusterClient = oldC
		pingRedis = oldP
	}()
	newClusterClient = func(*goredislib.ClusterOptions) *goredislib.ClusterClient {
		return goredislib.NewClusterClient(&goredislib.ClusterOptions{Addrs: []string{mr.Addr()}})
	}
	pingRedis = func(context.Context, goredislib.Cmdable) error { return nil }

	cfg := Config{ClusterAddrs: []string{mr.Addr()}}
	cli, err := NewClient(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	if cli.Redis() != nil {
		t.Fatal("cluster should not expose *goredislib.Client")
	}
}

func TestNewClient_sentinelHook(t *testing.T) {
	mr := miniredis.RunT(t)
	oldF := newFailoverClient
	oldP := pingRedis
	defer func() {
		newFailoverClient = oldF
		pingRedis = oldP
	}()
	newFailoverClient = func(*goredislib.FailoverOptions) *goredislib.Client {
		return goredislib.NewClient(&goredislib.Options{Addr: mr.Addr()})
	}
	pingRedis = func(context.Context, goredislib.Cmdable) error { return nil }

	cfg := Config{MasterName: "m", SentinelAddrs: []string{"127.0.0.1:26379"}}
	cli, err := NewClient(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	if cli.Redis() == nil {
		t.Fatal("expected raw client for sentinel hook")
	}
}

func Test_connectClient_tlsError(t *testing.T) {
	cfg := Config{
		Address: "127.0.0.1:6379",
		Cert:    "bad", Key: "bad", CA: "bad",
	}
	_, _, err := connectClient(cfg)
	if err == nil {
		t.Fatal("expected tls error")
	}
}

func TestClient_Close_idempotent(t *testing.T) {
	fc := &countingCloser{}
	c := &client{close: fc}
	for i := 0; i < 3; i++ {
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
	}
	if fc.n != 1 {
		t.Fatalf("Close called %d times, want 1", fc.n)
	}
}

func TestClient_Close_error_cached(t *testing.T) {
	want := errors.New("close fail")
	c := &client{close: &countingCloser{err: want}}
	if err := c.Close(); !errors.Is(err, want) {
		t.Fatal(err)
	}
	if err := c.Close(); !errors.Is(err, want) {
		t.Fatal(err)
	}
}

func TestNewClient_nilContext(t *testing.T) {
	mr := miniredis.RunT(t)
	//lint:ignore SA1012 intentional nil context in test
	cli, err := NewClient(nil, Config{Address: mr.Addr()})
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
}

type countingCloser struct {
	goredislib.Cmdable
	n   int
	err error
}

func (c *countingCloser) Close() error {
	c.n++
	return c.err
}
