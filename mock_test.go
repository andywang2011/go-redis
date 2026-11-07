package redis

import (
	"context"
	"testing"

	goredislib "github.com/redis/go-redis/v9"
)

type fakeRedis struct {
	goredislib.Cmdable
	closed bool
}

func (f *fakeRedis) PingWithRetry(ctx context.Context) error { return nil }
func (f *fakeRedis) Close() error {
	f.closed = true
	return nil
}
func (f *fakeRedis) Redis() *goredislib.Client { return nil }

var _ Client = (*fakeRedis)(nil)

func TestFakeRedis_Close(t *testing.T) {
	f := &fakeRedis{}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if !f.closed {
		t.Fatal("expected closed")
	}
}
